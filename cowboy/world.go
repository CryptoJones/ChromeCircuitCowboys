package cowboy

import (
	"math"
	"math/rand"
	"strings"
)

// respawnTicks is how many world ticks a slain mob stays down before it
// respawns in its home room (MajorMUD-style room cooldown, not instant).
const defaultRespawnTicks = 20

// World is the single shared game state. It is NOT safe for concurrent use; the
// server drives every method from one goroutine (commands + ticks serialized).
type World struct {
	rooms        map[string]*Room
	tmpls        map[string]*MobTemplate
	mobs         []*Mob
	players      map[int]*Player
	byName       map[string]*Player
	corpses      []*Corpse // dropped bodys awaiting recovery (in-memory; not persisted)
	floor        map[string]map[string]int // roomID -> item -> qty dropped on the floor (in-memory)
	nextID       int
	store        Persistence
	roll         func(n int) int // returns 0..n-1; injectable for tests
	respawnTicks int
	ringOffer    map[string][]int // ring-giver room -> indices into quests (randomized per session)
}

// corpsesIn returns the bodys lying in a room.
func (w *World) corpsesIn(roomID string) []*Corpse {
	var out []*Corpse
	for _, c := range w.corpses {
		if c.RoomID == roomID {
			out = append(out, c)
		}
	}
	return out
}

// removeCorpsesIn drops all corpses in a room from the world (after they're looted).
func (w *World) removeCorpsesIn(roomID string) {
	out := w.corpses[:0]
	for _, c := range w.corpses {
		if c.RoomID != roomID {
			out = append(out, c)
		}
	}
	w.corpses = out
}

// NewWorld builds the world from the static content and spawns initial mobs.
// store persists characters; pass NewMemStore() for ephemeral/test worlds.
func NewWorld(store Persistence) *World {
	w := &World{
		rooms:        buildRooms(),
		tmpls:        buildMobTemplates(),
		players:      map[int]*Player{},
		byName:       map[string]*Player{},
		store:        store,
		roll:         rand.Intn,
		respawnTicks: defaultRespawnTicks,
		floor:        map[string]map[string]int{},
	}
	for _, t := range w.tmpls {
		if t.Home != "" { // morph-only stages (no Home) are never spawned directly
			w.spawn(t)
		}
	}
	w.assignRingQuests() // scatter the RP ring "rumor" bounties across the ring givers
	return w
}

// maxRAM is a player's RAM ceiling — Intelligence-derived, plus any cyberdeck.
func maxRAM(p *Player) int { return 5 + p.Intelligence/2 + p.DeckBonus }

// SetRoll overrides the RNG (tests use this to make combat deterministic).
func (w *World) SetRoll(f func(n int) int) { w.roll = f }

// SetPrompter routes a player's status prompt to a dedicated sink (the server's
// managed-prompt writer), so prompts can be redrawn around async output without
// garbling the caller's in-progress input. nil = prompts go through the content
// sink (the default/test behavior).
func (w *World) SetPrompter(p *Player, fn func(string)) { p.prompter = fn }

func (w *World) spawn(t *MobTemplate) {
	w.mobs = append(w.mobs, &Mob{tmpl: t, origin: t, HP: t.HP, RoomID: t.Home})
}

// ---- accessors used by commands.go and tests ----

func (w *World) room(id string) *Room { return w.rooms[id] }

func (w *World) playersIn(roomID string, except *Player) []*Player {
	// A private room is a per-runner capsule pod: occupants never see, hear, or
	// can be targeted by anyone else, even though they share the room id.
	if r := w.rooms[roomID]; r != nil && r.Private {
		return nil
	}
	var out []*Player
	for _, p := range w.players {
		if p.RoomID == roomID && p != except {
			out = append(out, p)
		}
	}
	return out
}

func (w *World) liveMobsIn(roomID string) []*Mob {
	var out []*Mob
	for _, m := range w.mobs {
		if !m.dead && m.RoomID == roomID {
			out = append(out, m)
		}
	}
	return out
}

func (w *World) broadcast(roomID string, except *Player, msg string) {
	for _, p := range w.playersIn(roomID, except) {
		p.send(msg)
	}
}

// ---- session lifecycle ----

// Online reports whether a character is already connected (one session/name).
func (w *World) Online(name string) bool { _, ok := w.byName[strings.ToLower(name)]; return ok }

// HasCharacter reports whether a saved character exists (returning vs new) — the
// server uses this to decide whether to run character creation.
func (w *World) HasCharacter(name string) bool {
	_, ok, _ := w.store.Load(name)
	return ok
}

// Connect logs a character in, loading saved progress or, if none exists,
// creating a default one. Returning players use this; brand-new players use
// CreateCharacter (after the creation screen). out receives the player's text.
func (w *World) Connect(name string, out func(string)) *Player {
	p := w.newPlayer(name, out)
	if sp, ok, _ := w.store.Load(name); ok {
		applySave(p, sp)
	} else {
		newCharacter(p)
	}
	w.enter(p)
	w.deliverMail(p) // hand over any messages that arrived while they were away
	w.unsafeLogoutPenalty(p)
	return p
}

// unsafeLogoutPenalty docks a returning runner who logged out somewhere unsafe:
// they got jumped while offline. Costs 5% of max HP (never below 1) with a flavor
// line. Logging out in a Safe room is free.
func (w *World) unsafeLogoutPenalty(p *Player) {
	if r := w.room(p.RoomID); r == nil || r.Safe {
		return
	}
	dmg := p.MaxHP * 5 / 100
	if dmg < 1 {
		dmg = 1
	}
	if dmg >= p.HP {
		dmg = p.HP - 1 // wake hurt, not dead
	}
	if dmg < 1 {
		return
	}
	p.HP -= dmg
	who := unsafeAttackers[w.roll(len(unsafeAttackers))]
	verb := unsafeVerbs[w.roll(len(unsafeVerbs))]
	if w.inNet(p) {
		who = "a lurking ICE fragment"
	}
	p.send(style(red, "*** "+who+" "+verb+" while you were logged off here — "+itoa(dmg)+" damage. ***") + crlf)
}

var unsafeAttackers = []string{
	"a street ganger", "a scrap-scavver", "a strung-out wirehead",
	"a pickpocket", "a stray drone", "a back-alley cutthroat",
}
var unsafeVerbs = []string{
	"shanked you", "rifled your pockets", "kicked you awake",
	"jacked your ports", "worked you over", "lifted what they could",
}

// onlineByName returns a connected player by (case-insensitive) name, or nil.
func (w *World) onlineByName(name string) *Player {
	for _, o := range w.players {
		if strings.EqualFold(o.Name, name) {
			return o
		}
	}
	return nil
}

// CreateCharacter brings a brand-new player into the world using the loadout
// chosen on the creation screen.
func (w *World) CreateCharacter(name string, spec CharSpec, out func(string)) *Player {
	p := w.newPlayer(name, out)
	newCharacter(p)
	if c, ok := classByID(spec.ClassID); ok {
		p.Class = c.Name
	}
	if spec.Body > 0 {
		p.Body = spec.Body
	}
	if spec.Reflexes > 0 {
		p.Reflexes = spec.Reflexes
	}
	if spec.Intelligence > 0 {
		p.Intelligence = spec.Intelligence
	}
	p.MaxHP = maxHPFor(p)
	p.HP = p.MaxHP
	p.RAM = maxRAM(p)
	w.save(p) // persist immediately so a fresh character survives a crash
	w.enter(p)
	return p
}

func (w *World) newPlayer(name string, out func(string)) *Player {
	w.nextID++
	return &Player{ID: w.nextID, Name: name, Inv: map[string]int{}, Stash: map[string]int{}, Quests: map[string]int{}, Done: map[string]int{}, out: out}
}

// enter registers a fully-built player in the world and greets them.
func (w *World) enter(p *Player) {
	if w.rooms[p.RoomID] == nil {
		p.RoomID = startRoom
	}
	w.players[p.ID] = p
	w.byName[strings.ToLower(p.Name)] = p

	p.send(banner())
	greet := "You jack in as " + p.Name
	if p.Class != "" {
		greet += " the " + p.Class
	}
	p.send(style(neon, greet+". The grid accepts your signature.") + crlf)
	w.lookText(p)
	w.broadcast(p.RoomID, p, style(dim, p.Name+" materializes in a wash of static.")+crlf)
}

// Disconnect saves progress, announces the exit, and removes the player.
func (w *World) Disconnect(p *Player) {
	if p == nil {
		return
	}
	w.save(p)
	if p.fighting != nil && p.fighting.target == p {
		p.fighting.target = nil
	}
	p.pvpTarget = nil
	p.partyInvite = nil
	if p.trade != nil { // any open trade dies with the leaver
		w.cancelTrade(p.trade, p.Name+" jacked out")
	}
	for _, other := range w.players { // anyone duelling the leaver disengages
		if other.pvpTarget == p {
			other.pvpTarget = nil
			other.send(style(dim, p.Name+" jacked out — your duel ends.") + crlf)
		}
		if other.partyInvite == p { // pending invites from the leaver expire
			other.partyInvite = nil
			other.send(style(dim, p.Name+"'s crew invite expires.") + crlf)
		}
	}
	w.dropFromParty(p)
	w.broadcast(p.RoomID, p, style(dim, p.Name+" flatlines from the grid.")+crlf)
	delete(w.players, p.ID)
	delete(w.byName, strings.ToLower(p.Name))
}

// SaveAll persists every connected player. Called on the world goroutine (the
// server's periodic autosave + save-on-shutdown), so progress survives a server
// restart/crash, not only a clean per-player disconnect.
func (w *World) SaveAll() {
	for _, p := range w.players {
		w.save(p)
	}
}

func (w *World) save(p *Player) {
	_ = w.store.Save(&SavedPlayer{
		Name: p.Name, Class: p.Class, Level: p.Level, XP: p.XP, Eddies: p.Eddies,
		HP: p.HP, MaxHP: p.MaxHP, Body: p.Body, Reflexes: p.Reflexes,
		Intelligence: p.Intelligence, StatPoints: p.StatPoints, WeaponBonus: p.WeaponBonus,
		WeaponName: p.WeaponName, RAM: p.RAM, DeckBonus: p.DeckBonus,
		Room: p.RoomID, Inv: p.Inv, Stash: p.Stash, Quests: p.Quests, Done: p.Done,
	})
}

func newCharacter(p *Player) {
	p.Level, p.XP, p.Eddies = 1, 0, 50
	p.Body, p.Reflexes, p.Intelligence = 10, 10, 10
	p.RoomID = startRoom
	p.MaxHP = maxHPFor(p)
	p.HP = p.MaxHP
	p.RAM = maxRAM(p)
	p.Inv["stimpak"] = 1
}

func applySave(p *Player, sp *SavedPlayer) {
	p.Class = sp.Class
	p.Level, p.XP, p.Eddies = sp.Level, sp.XP, sp.Eddies
	p.HP, p.MaxHP = sp.HP, sp.MaxHP
	p.Body, p.Reflexes, p.Intelligence = sp.Body, sp.Reflexes, sp.Intelligence
	p.StatPoints = sp.StatPoints
	p.WeaponBonus, p.WeaponName = sp.WeaponBonus, sp.WeaponName
	p.DeckBonus = sp.DeckBonus
	p.RAM = sp.RAM
	if p.RAM <= 0 || p.RAM > maxRAM(p) {
		p.RAM = maxRAM(p)
	}
	p.RoomID = sp.Room
	if sp.Inv != nil {
		p.Inv = copyIntMap(sp.Inv)
	}
	if sp.Stash != nil {
		p.Stash = copyIntMap(sp.Stash)
	}
	if sp.Quests != nil {
		p.Quests = copyIntMap(sp.Quests)
	}
	if sp.Done != nil {
		p.Done = copyIntMap(sp.Done)
	}
	if p.Done == nil {
		p.Done = map[string]int{}
	}
	if p.MaxHP <= 0 {
		p.MaxHP = maxHPFor(p)
	}
	if p.HP <= 0 {
		p.HP = p.MaxHP
	}
}

// ---- progression math (MajorMUD-style: exponential XP, linear HP) ----

func copyIntMap(m map[string]int) map[string]int {
	out := make(map[string]int, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func maxHPFor(p *Player) int { return 20 + 2*p.Body + (p.Level-1)*8 }

// xpToNext is the XP needed to advance from the player's current level. Roughly
// 100 * level^1.4 — between LORD's N^1.5 and MajorMUD's N^1.3.
func xpToNext(level int) int { return int(100 * math.Pow(float64(level), 1.4)) }

// MaxLevel is the level cap (classic MUDs cap progression). At the cap, XP stops
// accumulating and excess is discarded.
const MaxLevel = 99

// pointsPerLevel is how many spendable character points a level-up banks (on top
// of the small automatic stat growth). SPEND them to customize your build.
const pointsPerLevel = 2

func (w *World) checkLevelUp(p *Player) {
	for p.Level < MaxLevel && p.XP >= xpToNext(p.Level) {
		p.XP -= xpToNext(p.Level)
		p.Level++
		p.Body += 2
		p.Reflexes++
		p.Intelligence++
		p.StatPoints += pointsPerLevel
		p.MaxHP = maxHPFor(p)
		p.HP = p.MaxHP
		p.send(style(neon, "*** UPLOAD COMPLETE — you reach level "+itoa(p.Level)+"! Stats boosted, HP restored. ***") + crlf)
		p.send(style(gold, "    +"+itoa(pointsPerLevel)+" character points to SPEND ("+itoa(p.StatPoints)+" banked).") + crlf)
		if p.Level == MaxLevel {
			p.send(style(gold, "*** You are MAXED — level "+itoa(MaxLevel)+", an elite cowboy. The grid is yours. ***") + crlf)
		}
	}
	if p.Level >= MaxLevel {
		p.XP = 0 // no more to gain at the cap
	}
}

// ---- combat math ----

func playerAC(p *Player) int { return 2 + p.Reflexes/5 }

func (w *World) toHit(dex, targetAC int) bool { return w.roll(20)+1+dex >= 10+targetAC }

func dmg(attack, targetAC int) int {
	d := attack - targetAC/2
	if d < 1 {
		d = 1
	}
	return d
}
