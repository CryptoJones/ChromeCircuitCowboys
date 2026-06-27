// Package cowboy is the engine for "Chrome Circuit Cowboys", a multiplayer
// cyberpunk MUD in the MajorMUD/Worldgroup tradition. It runs as a persistent
// resident door: one shared world, many simultaneous players, bridged in by
// AdmiralBBS. The engine itself is single-threaded and network-free — the
// server (cmd/cowboy) serializes all access on one goroutine and owns I/O, so
// the engine is deterministic and unit-testable.
package cowboy

// Player is a connected cowboy (runner).
type Player struct {
	ID           int
	Name         string
	Class        string
	IsBot        bool   // an AI "runner" that wanders + chatters to keep the world lively (#37); never saved or attacked
	Clan         string // clan tag (""=none); clanmates partying together earn bonus rewards
	Theme        string // color scheme: ""=default, "cbdark"/"cblight"=colorblind-friendly
	RoomID       string
	HP, MaxHP    int
	Level, XP    int
	Eddies       int
	Body         int            // melee/breach damage
	Reflexes     int            // dodge / damage reduction
	Intelligence int            // (flavor + future deck mechanics)
	StatPoints   int            // unspent character points (earned on level-up) — SPEND to raise a stat
	Done         map[string]int // completed one-time bounty ids (ring rumors are exempt/repeatable)
	WeaponBonus  int            // from a purchased weapon (e.g. ICE-breaker)
	WeaponName   string
	RAM          int            // netrun resource: powers breaches in the Net; regenerates out of combat
	DeckBonus    int            // bonus MaxRAM from a purchased cyberdeck (permanent)
	Inv          map[string]int // item name -> qty (carried; capped by level)
	Stash        map[string]int // item name -> qty stored at your Re-Clone Bay (uncapped, persists)
	Quests       map[string]int // active questID -> kills so far (>= Count means ready to claim)
	fighting     *Mob           // current mob target (nil = not in mob combat)
	pvpTarget    *Player        // current PvP target in the Net (nil = not duelling)
	party        *Party         // co-op crew (nil = solo)
	partyInvite  *Player        // pending crew invite from this leader (nil = none); consent before joining
	trade        *tradeSession  // active face-to-face trade (nil = not trading)
	hack         *hackGame      // active terminal hacking mini-game (nil = not hacking)
	downed       bool           // knocked out in a gym team-spar (out until the match resets)
	passwordHash string         // bcrypt hash of the character's password (""=unset; legacy)
	shieldTicks  int            // remaining ticks of the mirror program's damage shield
	shieldAmt    int            // flat damage reduction while shielded
	homing       int            // ticks left on a recall-home cast (0 = not recalling)
	confirmExit  string         // room id of a pending one-way departure awaiting YES/NO (""=none)
	dirty        bool           // output was sent since the last prompt (gates tick re-prompts)
	out          func(string)   // content output sink (set by the server; nil-safe via send)
	prompter     func(string)   // optional dedicated prompt sink (managed-prompt I/O); falls back to out
}

// attack is the player's deterministic damage per round.
func (p *Player) attack() int { return 3 + p.Body/2 + p.Level + p.WeaponBonus }

// defense reduces incoming damage (floored to 1 by the caller).
func (p *Player) defense() int { return p.Reflexes / 4 }

func (p *Player) send(s string) {
	if p.out != nil {
		p.out(recolor(p.Theme, s)) // remap the palette to the player's chosen theme
	}
	p.dirty = true // output emitted; a fresh prompt is owed (cleared by sendPrompt)
}

// MobTemplate is the static definition of a hostile program/NPC.
type MobTemplate struct {
	ID         string
	Name       string
	HP         int
	Damage     int // attack power
	AC         int // armor class (to-hit difficulty + light damage soak)
	XP         int
	Eddies     int
	Aggressive bool // attacks players on sight
	ICE        bool // a Net construct: shatters into "broken shards" (not a body) and "regenerates"
	Container  bool // an inert breakable (supply/data cache): cracks open into salvage, never a "body"
	Mechanical bool // a machine (drone/turret/mech): leaves "wreckage", not a flatlined body
	Home       string
	Next       string         // multi-stage ICE: on "death" it morphs into this template instead of dying
	Drops      map[string]int // item drops seeded into the corpse on death (e.g. loot-cache consumables)
}

// Mob is a live instance of a MobTemplate in the world.
type Mob struct {
	tmpl         *MobTemplate // current template (may be a later stage after a morph)
	origin       *MobTemplate // the spawn template — restored on respawn so multi-stage ICE resets
	HP           int
	RoomID       string
	target       *Player
	respawnIn    int // ticks until this dead mob respawns (0 = alive)
	dead         bool
	awaitingLoot bool // dead body not yet looted — respawn stays gated until it is
}

// Room is one location in the city/net.
type Room struct {
	ID      string
	Name    string
	Desc    string
	Exits   map[string]string // direction -> room id
	Vendor  bool              // a shop operates here
	Medic   bool              // a Emergency Medic operates here — re-install salvaged cyberware
	Private bool              // a per-runner capsule pod — occupants are isolated (no one shares it)
	Safe    bool              // no-violence zone (outside the clone pods): a security drone flatlines PvP aggressors
	Net     bool              // inside the Net: combat is a BREACH (Intelligence + RAM), ICE shatters into shards
	Spar    bool              // a sparring gym: PvP is non-lethal — a downed runner is knocked out, keeps everything
	Term    bool              // a data terminal: SEND mail / WIRE scrip to other runners (vendors/medics also count)
}

// Mail is a stored message from one runner to another, delivered whenever the
// recipient next reads MAIL (works even if they were offline).
type Mail struct {
	From string
	Body string
}

// Corpse is a dead runner's old body, left where they flatlined. It holds the
// gear (consumables + cyberware by ware name) the new clone woke up without.
// It persists until looted; anyone may loot it (open recovery + risk).
type Corpse struct {
	Owner  string
	RoomID string
	Loot   map[string]int // ware name -> qty (consumables usable on loot; cyberware needs Emergency Medic INSTALL)
	Scrip  int            // scrip carried by the body (mob loot); 0 for runner corpses
	mob    *Mob           // the slain mob this body belongs to (nil for runner corpses); looting it ungates respawn
	IsICE  bool           // ICE construct's "broken shards" (salvage), not a flatlined body
	IsBox  bool           // inert container (supply/data cache): "a cracked-open cache", never a body
	IsMech bool           // a machine (drone/turret/mech): "wreckage", never a flatlined body
}

// SavedPlayer is the persisted slice of a Player (progress survives logout).
type SavedPlayer struct {
	Name                         string
	Class                        string
	Clan                         string
	Theme                        string
	PasswordHash                 string
	Level, XP, Eddies, HP, MaxHP int
	Body, Reflexes, Intelligence int
	StatPoints                   int
	WeaponBonus                  int
	WeaponName                   string
	RAM                          int
	DeckBonus                    int
	Room                         string
	Inv                          map[string]int
	Stash                        map[string]int
	Quests                       map[string]int
	Done                         map[string]int
}

// Persistence stores character progress between sessions. The server backs it
// with SQLite; tests use an in-memory implementation.
type Persistence interface {
	Load(name string) (*SavedPlayer, bool, error)
	Save(sp *SavedPlayer) error
	// Top returns up to n characters ranked by level then XP (for the leaderboard).
	Top(n int) ([]SavedPlayer, error)
	// PushMail queues a message for a recipient; PopMail returns and clears theirs.
	PushMail(to, from, body string) error
	PopMail(to string) ([]Mail, error)
}

// MemStore is an in-memory Persistence for tests and ephemeral runs.
type MemStore struct {
	m    map[string]*SavedPlayer
	mail map[string][]Mail
}

// NewMemStore builds an empty in-memory store.
func NewMemStore() *MemStore {
	return &MemStore{m: map[string]*SavedPlayer{}, mail: map[string][]Mail{}}
}

// PushMail queues a message for to.
func (s *MemStore) PushMail(to, from, body string) error {
	s.mail[to] = append(s.mail[to], Mail{From: from, Body: body})
	return nil
}

// PopMail returns and clears the recipient's queued mail.
func (s *MemStore) PopMail(to string) ([]Mail, error) {
	m := s.mail[to]
	delete(s.mail, to)
	return m, nil
}

// Load returns a saved character by name.
func (s *MemStore) Load(name string) (*SavedPlayer, bool, error) {
	sp, ok := s.m[name]
	return sp, ok, nil
}

// Save upserts a character.
func (s *MemStore) Save(sp *SavedPlayer) error {
	cp := *sp
	s.m[sp.Name] = &cp
	return nil
}

// Top returns up to n characters ranked by level then XP.
func (s *MemStore) Top(n int) ([]SavedPlayer, error) {
	out := make([]SavedPlayer, 0, len(s.m))
	for _, sp := range s.m {
		out = append(out, *sp)
	}
	sortSavedByRank(out)
	if len(out) > n {
		out = out[:n]
	}
	return out, nil
}
