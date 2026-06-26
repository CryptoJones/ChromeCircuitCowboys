package cowboy

import (
	"sort"
	"strconv"
	"strings"
)

// inNet reports whether the player is inside the Net — where attacks are netrun
// BREACHes driven by Intelligence (and spend RAM), not meatspace strikes driven
// by Body. Every authored Net room carries the Net flag (see netzones.go).
func (w *World) inNet(p *Player) bool {
	r := w.room(p.RoomID)
	return r != nil && r.Net
}

// effAttack is the player's damage this round, route-dependent (breach vs melee).
func (w *World) effAttack(p *Player) int {
	if w.inNet(p) {
		return 3 + p.Intelligence/2 + p.Level + p.WeaponBonus
	}
	return p.attack()
}

var dirAliases = map[string]string{
	"n": "north", "s": "south", "e": "east", "w": "west", "u": "up", "d": "down",
	"north": "north", "south": "south", "east": "east", "west": "west", "up": "up", "down": "down",
	"out": "out", "o": "out", // capsule pod -> street
	"in": "in", // street -> your capsule pod
}

// Command parses and executes a single input line for player p. It returns true
// if the player asked to quit (the server then disconnects them).
func (w *World) Command(p *Player, line string) (quit bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		w.sendPrompt(p)
		return false
	}
	// ":<action>" is RP-emote shorthand (": waves" or ":waves").
	if strings.HasPrefix(line, ":") {
		w.emote(p, strings.TrimSpace(line[1:]))
		w.sendPrompt(p)
		return false
	}
	// ";<msg>" is the quick party-chat shortcut (like GSAY) — straight to crew.
	if strings.HasPrefix(line, ";") {
		w.groupChat(p, strings.TrimSpace(line[1:]))
		w.sendPrompt(p)
		return false
	}
	fields := strings.Fields(line)
	cmd := strings.ToLower(fields[0])
	arg := strings.TrimSpace(strings.TrimPrefix(line, fields[0]))

	if dir, ok := dirAliases[cmd]; ok {
		w.move(p, dir)
		w.sendPrompt(p)
		return false
	}

	// A bare number is a hack guess while a run is live, otherwise an inventory
	// quick-use slot (the server emits it on a single digit keypress).
	if n, err := strconv.Atoi(cmd); err == nil && arg == "" {
		if p.hack != nil {
			w.hackGuess(p, n)
		} else {
			w.quickUse(p, n)
		}
		w.sendPrompt(p)
		return false
	}

	switch cmd {
	case "look", "l":
		if strings.TrimSpace(arg) != "" {
			w.examine(p, arg)
		} else {
			w.lookText(p)
		}
	case "roomid", "whereami", "#id":
		w.showRoomID(p) // hidden dev/SysOp command (not in HELP)
	case "open":
		w.openContainer(p, arg)
	case "talk", "ask", "speak":
		w.talk(p, arg)
	case "pay", "hire":
		w.payJoytoy(p, arg)
	case "hack":
		w.startHack(p)
	case "send":
		w.sendMail(p, arg)
	case "wire", "transfer":
		w.wireScrip(p, arg)
	case "mail", "messages":
		w.readMail(p)
	case "trade":
		w.trade(p, arg)
	case "offer":
		w.tradeOffer(p, arg)
	case "confirm":
		w.tradeConfirm(p)
	case "cancel":
		w.tradeCancel(p)
	case "map", "m":
		w.showMap(p)
	case "say", "'":
		w.say(p, arg)
	case "emote", "me", "em":
		w.emote(p, arg)
	case "who":
		w.who(p)
	case "score", "stats", "st", "sc":
		w.score(p)
	case "spend":
		w.spend(p, arg)
	case "attack", "a", "kill", "k", "breach":
		w.engage(p, arg)
	case "flee", "jackout", "disconnect":
		w.flee(p)
	case "list", "shop":
		w.list(p)
	case "buy":
		w.buy(p, arg)
	case "sell":
		w.sell(p, arg)
	case "use":
		w.use(p, arg)
	case "drop":
		w.drop(p, arg)
	case "get", "pickup":
		w.pickUp(p, arg)
	case "loot", "lo", "salvage":
		w.loot(p)
	case "install", "ripper":
		w.install(p, arg)
	case "give", "hand":
		w.give(p, arg)
	case "inventory", "inv", "i":
		w.inventory(p)
	case "stash":
		w.stash(p, arg)
	case "grab", "unstash", "withdraw":
		w.grab(p, arg)
	case "quests", "missions", "bounties":
		w.showQuests(p)
	case "accept", "take":
		// A bare ACCEPT takes a pending crew invite; ACCEPT <#> claims a bounty.
		if arg == "" && p.partyInvite != nil {
			w.acceptInvite(p)
		} else {
			w.accept(p, arg)
		}
	case "decline":
		w.declineInvite(p)
	case "claim", "turnin":
		w.claim(p)
	case "run", "exec":
		w.run(p, arg)
	case "programs", "demons":
		w.listPrograms(p)
	case "group", "crew":
		w.group(p, arg)
	case "invite":
		w.invite(p, arg)
	case "home", "rest":
		w.goHome(p)
	case "leave", "ungroup":
		w.leaveParty(p)
	case "gsay", "crewchat", "party":
		w.groupChat(p, arg)
	case "leaderboard", "top", "rankings":
		w.leaderboard(p)
	case "help", "?", "commands":
		p.send(helpText())
	case "quit", "logout", "exit":
		p.send(style(neon, "Jacking out. The grid forgets you... for now.") + crlf)
		return true
	default:
		p.send(style(dim, "Unknown command. Type HELP.") + crlf)
	}
	w.sendPrompt(p)
	return false
}

// Prompt re-displays the player's status prompt (used by the server right after
// a player joins).
func (w *World) Prompt(p *Player) { w.sendPrompt(p) }

// PromptIfDirty re-displays the prompt ONLY if the player received output since
// their last prompt. The server calls this after each world tick so a player who
// saw combat/chat/room output gets a fresh prompt — but an IDLE player does not
// get the prompt re-printed every tick (which would spam it while they read).
func (w *World) PromptIfDirty(p *Player) {
	if p.dirty {
		w.sendPrompt(p)
	}
}

func (w *World) sendPrompt(p *Player) {
	hpColor := green
	if p.HP*3 < p.MaxHP {
		hpColor = red
	}
	mode := "MEAT"
	ram := ""
	if w.inNet(p) {
		mode = "NET"
		ram = style(neon, " ["+itoa(p.RAM)+"/"+itoa(maxRAM(p))+"ram]")
	}
	promptStr := style(hpColor, "["+itoa(p.HP)+"/"+itoa(p.MaxHP)+"hp]") + ram +
		style(dim, " ["+mode+"] ") + style(green, "> ")
	if p.prompter != nil {
		p.prompter(promptStr) // managed-prompt sink (redraws around async output)
	} else {
		p.send(promptStr)
	}
	p.dirty = false // prompt now shown; nothing owed until the next output
}

func (w *World) lookText(p *Player) {
	r := w.room(p.RoomID)
	if r == nil {
		p.send(style(red, "You are nowhere. (corrupted location)") + crlf)
		return
	}
	p.send(crlf + style(neon, r.Name) + crlf + r.Desc + crlf)
	if r.Vendor {
		p.send(style(gold, "A vendor terminal hums here. Type LIST.") + crlf)
	}
	if r.Medic {
		p.send(style(gold, "A Emergency Medic's chair waits here. INSTALL salvaged cyberware.") + crlf)
	}
	if len(w.questsHere(p)) > 0 {
		p.send(style(gold, "Someone here is hiring — type QUESTS.") + crlf)
	}
	if _, hasNPC := roomNPC[p.RoomID]; hasNPC || loreKey(p.RoomID) != "" {
		p.send(style(dim, "There's someone here to TALK to.") + crlf)
	}
	if p.RoomID == startRoom {
		p.send(style(dim, "The clone-booth tech is here — TALK for a quick orientation.") + crlf)
	}
	if pile := w.floorList(p.RoomID); pile != "" {
		p.send(style(gold, "On the floor: ") + pile + style(dim, " (GET <item> / GET ALL)") + crlf)
	}
	if jt, ok := joytoyRooms[p.RoomID]; ok {
		p.send(style(dim, jt+" is working here — PAY for some company (€$"+itoa(joytoyFee)+").") + crlf)
	}
	// Exits.
	var dirs []string
	for _, d := range []string{"north", "south", "east", "west", "up", "down", "in", "out"} {
		if _, ok := r.Exits[d]; ok {
			dirs = append(dirs, d)
		}
	}
	p.send(style(dim, "Exits: "+strings.Join(dirs, ", ")) + crlf)
	// Other players.
	for _, other := range w.playersIn(p.RoomID, p) {
		p.send(style(green, other.Name+" is here.") + crlf)
	}
	// Mobs.
	for _, m := range w.liveMobsIn(p.RoomID) {
		p.send(style(hot, m.tmpl.Name+" is here.") + crlf)
	}
	// Flatlined bodies / shattered ICE / cracked-open caches waiting to be looted.
	for _, c := range w.corpsesIn(p.RoomID) {
		switch {
		case c.IsBox:
			p.send(style(dim, "A cracked-open "+strings.TrimPrefix(c.Owner, "a sealed ")+" lies spilled here. (LOOT)") + crlf)
		case c.IsICE:
			p.send(style(dim, "Broken shards of "+c.Owner+" glitter here. (LOOT)") + crlf)
		case c.IsMech:
			p.send(style(dim, "The wreck of "+c.Owner+" lies here. (LOOT)") + crlf)
		default:
			p.send(style(dim, c.Owner+"'s flatlined body lies here. (LOOT)") + crlf)
		}
	}
}

// emote broadcasts a freeform third-person action to the room, for the RP crowd:
// EMOTE / ME / ":" + an action -> "Wintermute lights a cigarette."
func (w *World) emote(p *Player, action string) {
	action = strings.TrimSpace(action)
	if action == "" {
		p.send(style(dim, "Emote what? e.g. ME lights a cigarette (or :leans on the wall).") + crlf)
		return
	}
	line := style(neon, p.Name+" "+action) + crlf
	p.send(line)
	w.broadcast(p.RoomID, p, line)
}

func (w *World) say(p *Player, msg string) {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		p.send(style(dim, "Say what?") + crlf)
		return
	}
	p.send(style(green, "You say: ") + msg + crlf)
	w.broadcast(p.RoomID, p, style(green, p.Name+" says: ")+msg+crlf)
}

func (w *World) who(p *Player) {
	p.send(style(neon, "-- Jacked in right now --") + crlf)
	for _, o := range w.players {
		cls := o.Class
		if cls != "" {
			cls = " " + cls
		}
		p.send("  " + style(green, o.Name) + style(dim, "  (level "+itoa(o.Level)+cls+")") + crlf)
	}
}

func (w *World) score(p *Player) {
	class := p.Class
	if class == "" {
		class = "cowboy"
	}
	p.send(crlf + style(neon, "== "+p.Name+" :: "+class+" ==") + crlf)
	xpLine := "  Level " + itoa(p.Level) + "   XP " + itoa(p.XP) + "/" + itoa(xpToNext(p.Level))
	if p.Level >= MaxLevel {
		xpLine = "  Level " + itoa(p.Level) + " " + style(gold, "(MAX)")
	}
	p.send(xpLine + crlf)
	p.send("  HP " + itoa(p.HP) + "/" + itoa(p.MaxHP) + "   RAM " + itoa(p.RAM) + "/" + itoa(maxRAM(p)) + "   Armor Class " + itoa(playerAC(p)) + crlf)
	p.send("  Body " + itoa(p.Body) + "   Reflexes " + itoa(p.Reflexes) + "   Intelligence " + itoa(p.Intelligence) + crlf)
	// Always show the points line (0 when none); when there are some, add a bold
	// call-to-action so it's obvious you can spend them.
	if p.StatPoints > 0 {
		p.send(style(gold, "  Character points: "+itoa(p.StatPoints)) + style(dim, " — type SPEND <body|reflexes|intelligence>") + crlf)
		p.send(style(hot, "  *** You have character points to spend. ***") + crlf)
	} else {
		p.send("  Character points: 0" + crlf)
	}
	weapon := "bare fists"
	if p.WeaponName != "" {
		weapon = p.WeaponName + " (+" + itoa(p.WeaponBonus) + " atk)"
	}
	p.send("  Weapon: " + weapon + crlf)
	deck := "stock deck"
	if p.DeckBonus > 0 {
		deck = "cyberdeck (+" + itoa(p.DeckBonus) + " max RAM)"
	}
	p.send("  Deck: " + deck + crlf)
	p.send(style(gold, "  €$ "+itoa(p.Eddies)+" scrip") + crlf)
	if p.shieldTicks > 0 {
		p.send(style(dim, "  Mirror shield: -"+itoa(p.shieldAmt)+" dmg for "+itoa(p.shieldTicks)+" more round(s)") + crlf)
	}
	if p.party != nil && len(p.party.Members) > 1 {
		p.send(style(dim, "  Crew: "+itoa(len(p.party.Members))+" members (GROUP to view)") + crlf)
	}
	p.send(style(dim, "  Programs: RUN <name> — see PROGRAMS") + crlf)
}

// spend raises a stat with banked character points (1 point = +1). Raising Body
// also lifts MaxHP (and heals the new headroom). Usage: SPEND <stat> [amount].
func (w *World) spend(p *Player, arg string) {
	if p.StatPoints <= 0 {
		p.send(style(dim, "You have no character points to spend. Level up to earn them.") + crlf)
		return
	}
	fields := strings.Fields(strings.ToLower(strings.TrimSpace(arg)))
	if len(fields) == 0 {
		p.send(style(gold, "You have "+itoa(p.StatPoints)+" character point(s).") + crlf)
		p.send(style(dim, "SPEND <body|reflexes|intelligence> [amount] — 1 point = +1.") + crlf)
		return
	}
	n := 1
	if len(fields) >= 2 {
		if v, err := strconv.Atoi(fields[1]); err == nil {
			n = v
		}
	}
	if n < 1 {
		n = 1
	}
	if n > p.StatPoints {
		n = p.StatPoints
	}
	switch fields[0] {
	case "body", "bod", "b":
		p.Body += n
		old := p.MaxHP
		p.MaxHP = maxHPFor(p)
		p.HP += p.MaxHP - old
	case "reflexes", "ref", "r":
		p.Reflexes += n
	case "intelligence", "int", "i":
		p.Intelligence += n
	default:
		p.send(style(dim, "Spend on what? body, reflexes, or intelligence.") + crlf)
		return
	}
	p.StatPoints -= n
	p.send(style(green, "Spent "+itoa(n)+" point(s) on "+fields[0]+". ") + style(dim, itoa(p.StatPoints)+" left.") + crlf)
	w.save(p)
}

// sortedInv returns the carried item names (qty > 0) in a stable order, so the
// INVENTORY numbering and the digit quick-use map to the same slots.
func sortedInv(p *Player) []string {
	names := make([]string, 0, len(p.Inv))
	for k, q := range p.Inv {
		if q > 0 {
			names = append(names, k)
		}
	}
	sort.Strings(names)
	return names
}

func (w *World) inventory(p *Player) {
	p.send(style(neon, "-- Inventory ("+itoa(invCount(p))+"/"+itoa(carryCap(p))+") --") + crlf)
	p.send(style(gold, "  €$ "+itoa(p.Eddies)+" scrip") + crlf)
	items := sortedInv(p)
	if len(items) == 0 {
		p.send(style(dim, "  (no items)") + crlf)
		return
	}
	for i, name := range items {
		p.send("  " + style(gold, itoa(i+1)+")") + " " + name + " x" + itoa(p.Inv[name]) + crlf)
	}
	p.send(style(dim, "  press a number to USE that item (no Enter), or USE <name>") + crlf)
}

// quickUse uses the Nth inventory item (as numbered by INVENTORY) — driven by a
// single digit keypress at the prompt, so it's fast mid-combat.
func (w *World) quickUse(p *Player, n int) {
	items := sortedInv(p)
	if n < 1 || n > len(items) {
		p.send(style(dim, "No item #"+itoa(n)+" in your pack. (type I to list)") + crlf)
		return
	}
	w.use(p, items[n-1])
}

// atStash reports whether the runner is at their Re-Clone Bay, where their
// personal stash lives.
func (w *World) atStash(p *Player) bool { return p.RoomID == startRoom }

// stash with no arg shows your bay stash; STASH <item> deposits all of an item
// from your pack into the (uncapped) stash. Only usable at your Re-Clone Bay.
func (w *World) stash(p *Player, arg string) {
	if !w.atStash(p) {
		p.send(style(dim, "Your stash is back at your Re-Clone Bay — go HOME to reach it.") + crlf)
		return
	}
	arg = strings.ToLower(strings.TrimSpace(arg))
	if arg == "" {
		p.send(style(neon, "-- Stash :: Re-Clone Bay --") + crlf)
		if len(p.Stash) == 0 {
			p.send(style(dim, "  (empty) — STASH <item> to store, GRAB <item> to withdraw") + crlf)
			return
		}
		for name, qty := range p.Stash {
			p.send("  " + name + " x" + itoa(qty) + crlf)
		}
		return
	}
	n := p.Inv[arg]
	if n <= 0 {
		p.send(style(dim, "You're not carrying any "+arg+".") + crlf)
		return
	}
	p.Stash[arg] += n
	delete(p.Inv, arg)
	p.send(style(green, "Stashed "+itoa(n)+"x "+arg+". (pack "+itoa(invCount(p))+"/"+itoa(carryCap(p))+")") + crlf)
}

// grab withdraws an item from your bay stash back into your pack, up to your
// carry cap. Only usable at your Re-Clone Bay.
func (w *World) grab(p *Player, arg string) {
	if !w.atStash(p) {
		p.send(style(dim, "Your stash is back at your Re-Clone Bay — go HOME to reach it.") + crlf)
		return
	}
	arg = strings.ToLower(strings.TrimSpace(arg))
	if arg == "" {
		p.send(style(dim, "Grab what? (see STASH)") + crlf)
		return
	}
	have := p.Stash[arg]
	if have <= 0 {
		p.send(style(dim, "No "+arg+" in your stash.") + crlf)
		return
	}
	room := carryCap(p) - invCount(p)
	if room <= 0 {
		p.send(style(dim, "Your pack is full ("+itoa(invCount(p))+"/"+itoa(carryCap(p))+").") + crlf)
		return
	}
	take := have
	if take > room {
		take = room
	}
	p.Stash[arg] -= take
	if p.Stash[arg] <= 0 {
		delete(p.Stash, arg)
	}
	p.Inv[arg] += take
	msg := "Grabbed " + itoa(take) + "x " + arg + ". (pack " + itoa(invCount(p)) + "/" + itoa(carryCap(p)) + ")"
	if take < have {
		msg += " — " + itoa(have-take) + " left in the stash (pack full)"
	}
	p.send(style(green, msg) + crlf)
}

func (w *World) move(p *Player, dir string) {
	if p.fighting != nil {
		p.send(style(hot, "You're in combat! Break the connection with FLEE first.") + crlf)
		return
	}
	if p.homing > 0 { // moving breaks a recall cast
		p.homing = 0
		p.send(style(dim, "Your recall fizzles as you move.") + crlf)
	}
	r := w.room(p.RoomID)
	dest, ok := r.Exits[dir]
	if !ok {
		p.send(style(dim, "You can't go "+dir+".") + crlf)
		return
	}
	origin := p.RoomID
	w.broadcast(p.RoomID, p, style(dim, p.Name+" heads "+dir+".")+crlf)
	p.RoomID = dest
	w.broadcast(p.RoomID, p, style(dim, p.Name+" arrives.")+crlf)
	w.lookText(p)
	w.partyFollow(p, origin, dest, dir)
}

// partyFollow pulls the leader's crew along when the leader moves: any member in
// the room the leader just left (and not mid-combat) follows to the destination.
func (w *World) partyFollow(leader *Player, origin, dest, dir string) {
	if leader.party == nil || leader.party.Leader != leader {
		return
	}
	for _, m := range leader.party.Members {
		if m == leader || m.RoomID != origin {
			continue
		}
		if m.fighting != nil || m.pvpTarget != nil {
			m.send(style(dim, leader.Name+" heads "+dir+" — you stay and finish the fight.") + crlf)
			continue
		}
		m.homing = 0 // moving cancels any recall they were casting
		m.RoomID = dest
		w.broadcast(dest, m, style(dim, m.Name+" follows "+leader.Name+" in.")+crlf)
		m.send(style(dim, "You follow "+leader.Name+" "+dir+".") + crlf)
		w.lookText(m)
	}
}

// recallTicks is the cast time of a HOME recall — ~10 seconds at the default 2s
// world tick. The recall completes on tickRecall (combat.go); it is broken if the
// runner is hit by a hostile or moves before it lands.
const recallTicks = 5

// goHome jacks a RECALL protocol: a timed teleport back to your Re-Clone Bay from
// ANYWHERE. It takes recallTicks to land and is interrupted by a mob/PvP hit or by
// moving — so it's an escape you have to survive, not an instant bail. (To just
// step into your pod from the street, go IN at Neon Alley.)
func (w *World) goHome(p *Player) {
	if p.RoomID == startRoom {
		p.send(style(dim, "You're already in your Re-Clone Bay.") + crlf)
		return
	}
	if p.fighting != nil || p.pvpTarget != nil {
		p.send(style(hot, "You can't focus a recall mid-fight — break it with FLEE first.") + crlf)
		return
	}
	if p.homing > 0 {
		p.send(style(dim, "You're already jacking a recall — hold still.") + crlf)
		return
	}
	p.homing = recallTicks
	p.send(style(neon, "You jack a recall protocol. Hold still (~10s) and you'll phase home to your Re-Clone Bay — a hit or a move breaks it.") + crlf)
	w.broadcast(p.RoomID, p, style(dim, p.Name+" flickers — phasing out.")+crlf)
}

// examine gives a detailed look at an item — `LOOK <item>`. Resolves from the
// catalog (so any known ware can be examined) and prints its flavor + mechanical
// effects. (Richer authored lore is a follow-on, #54.)
func (w *World) examine(p *Player, arg string) {
	name := strings.ToLower(strings.TrimSpace(arg))
	x, ok := findWare(name)
	if !ok {
		p.send(style(dim, "You don't see '"+arg+"' to examine. (try an item in your INVENTORY)") + crlf)
		return
	}
	carried := ""
	if p.Inv[name] > 0 {
		carried = style(dim, "  (you carry "+itoa(p.Inv[name])+")")
	}
	p.send(crlf + style(neon, x.name) + carried + crlf)
	p.send(style(green, "  "+x.desc) + crlf)
	var effects []string
	if x.heal > 0 {
		effects = append(effects, "restores "+itoa(x.heal)+" HP")
	}
	if x.ram > 0 {
		effects = append(effects, "restores "+itoa(x.ram)+" RAM")
	}
	if x.bonus > 0 {
		effects = append(effects, "+"+itoa(x.bonus)+" attack (install at a medic)")
	}
	if x.deck > 0 {
		effects = append(effects, "+"+itoa(x.deck)+" max RAM (install at a medic)")
	}
	if x.body > 0 {
		effects = append(effects, "+"+itoa(x.body)+" Body implant")
	}
	if x.refl > 0 {
		effects = append(effects, "+"+itoa(x.refl)+" Reflexes implant")
	}
	if x.intel > 0 {
		effects = append(effects, "+"+itoa(x.intel)+" Intelligence implant")
	}
	if len(effects) > 0 {
		p.send(style(gold, "  Effect: ") + strings.Join(effects, ", ") + crlf)
	}
	if x.price > 0 {
		p.send(style(dim, "  Market value: €$"+itoa(x.price)+" (sells for ~€$"+itoa(x.price*sellBuyback/100)+")") + crlf)
	}
}

// showRoomID prints the current room's internal ID (plus name + exits) — a
// hidden command (kept out of HELP) for building/debugging, so a room can be
// referred to by its actual id.
func (w *World) showRoomID(p *Player) {
	r := w.room(p.RoomID)
	name := p.RoomID
	if r != nil && r.Name != "" {
		name = r.Name
	}
	p.send(style(neon, "roomid: ") + style(gold, p.RoomID) + style(dim, "  ("+name+")") + crlf)
	if r != nil {
		var dirs []string
		for _, d := range []string{"north", "south", "east", "west", "up", "down", "in", "out"} {
			if dest, ok := r.Exits[d]; ok {
				dirs = append(dirs, d+"→"+dest)
			}
		}
		p.send(style(dim, "  exits: "+strings.Join(dirs, ", ")) + crlf)
	}
}

// openContainer breaks open a supply/data cache in the room — the intuitive
// verb for an inert container (you don't "attack" a crate). It engages the
// cache like a passive mob (Damage 0), so it cracks open and drops loot to LOOT.
func (w *World) openContainer(p *Player, arg string) {
	arg = strings.ToLower(strings.TrimSpace(arg))
	var target *Mob
	for _, m := range w.liveMobsIn(p.RoomID) {
		if !m.tmpl.Container {
			continue
		}
		if arg == "" || strings.Contains(strings.ToLower(m.tmpl.Name), arg) {
			target = m
			break
		}
	}
	if target == nil {
		p.send(style(dim, "There's nothing to open here. (Caches OPEN; everything else you ATTACK.)") + crlf)
		return
	}
	w.engage(p, target.tmpl.Name)
}

// loot strips every flatlined body in the room into your pack. Items are
// usable immediately; salvaged cyberware must be re-installed at a Emergency Medic.
// Open recovery: anyone can loot any body (recover for a crewmate — or swipe it).
func (w *World) loot(p *Player) {
	cs := w.corpsesIn(p.RoomID)
	if len(cs) == 0 {
		p.send(style(dim, "There's nothing here to loot.") + crlf)
		return
	}
	total := 0
	scrip := 0
	// Classify what's lying here so the wording fits: an inert container (cache),
	// a Net construct's shards, or an actual flatlined body. A body in the pile
	// wins the phrasing (the most general, never-wrong noun).
	box, ice, mech, body := false, false, false, false
	var cyber []string
	for _, c := range cs {
		switch {
		case c.IsBox:
			box = true
		case c.IsICE:
			ice = true
		case c.IsMech:
			mech = true
		default:
			body = true
		}
		for name, qty := range c.Loot {
			if qty <= 0 {
				continue
			}
			p.Inv[name] += qty
			total += qty
			if x, ok := findWare(name); ok && (x.bonus > 0 || x.deck > 0) {
				cyber = append(cyber, name)
			}
		}
		scrip += c.Scrip
		if c.mob != nil {
			// Looting a slain mob ungates its respawn — meat bodies and ICE shards
			// alike; the area refills after the normal cooldown, never before.
			c.mob.awaitingLoot = false
			c.mob.respawnIn = w.respawnTicks
		}
	}
	w.removeCorpsesIn(p.RoomID)
	// A body in the pile wins the phrasing; otherwise a cache, a machine wreck,
	// or Net shards.
	useBox := !body && box
	useMech := !body && !box && mech
	useICE := !body && !box && !mech && ice
	if total == 0 && scrip == 0 {
		switch {
		case useBox:
			p.send(style(dim, "The cache is already cleaned out.") + crlf)
		case useMech:
			p.send(style(dim, "The wreck is already stripped — nothing but scrap.") + crlf)
		case useICE:
			p.send(style(dim, "The shards are inert — nothing to salvage.") + crlf)
		default:
			p.send(style(dim, "The body is already stripped bare.") + crlf)
		}
		return
	}
	if scrip > 0 {
		p.Eddies += scrip
		switch {
		case useBox:
			p.send(style(gold, "You scoop €$"+itoa(scrip)+" scrip from the cracked-open cache.") + crlf)
		case useMech:
			p.send(style(gold, "You pull €$"+itoa(scrip)+" scrip from the wreckage.") + crlf)
		case useICE:
			p.send(style(gold, "You salvage €$"+itoa(scrip)+" scrip from the broken shards.") + crlf)
		default:
			p.send(style(gold, "You recover €$"+itoa(scrip)+" scrip from the body.") + crlf)
		}
	}
	if total > 0 {
		switch {
		case useBox:
			p.send(style(green, "You clear out the cache — its contents are now in your pack.") + crlf)
		case useMech:
			p.send(style(green, "You strip the wreck — its gear is now in your pack.") + crlf)
		case useICE:
			p.send(style(green, "You pick the shards clean — the salvage is now in your pack.") + crlf)
		default:
			p.send(style(green, "You strip the body — its gear is now in your pack.") + crlf)
		}
		if len(cyber) > 0 {
			p.send(style(neon, "Salvaged cyberware: ") + strings.Join(cyber, ", ") +
				style(dim, " — INSTALL it at a Emergency Medic to use it again.") + crlf)
		}
	}
	switch {
	case useBox:
		w.broadcast(p.RoomID, p, style(dim, p.Name+" cracks open a cache.")+crlf)
	case useMech:
		w.broadcast(p.RoomID, p, style(dim, p.Name+" picks through a wreck.")+crlf)
	case useICE:
		w.broadcast(p.RoomID, p, style(dim, p.Name+" picks through the broken shards.")+crlf)
	default:
		w.broadcast(p.RoomID, p, style(dim, p.Name+" loots a flatlined body.")+crlf)
	}
}

// give hands an inventory item to another runner in the room (e.g. returning a
// crewmate's recovered gear). Syntax: GIVE <item> <runner>.
func (w *World) give(p *Player, arg string) {
	fields := strings.Fields(arg)
	if len(fields) < 2 {
		p.send(style(dim, "Give what to whom? GIVE <item> <runner>.") + crlf)
		return
	}
	targetName := fields[len(fields)-1]
	item := strings.ToLower(strings.Join(fields[:len(fields)-1], " "))
	target := w.playerInRoomByName(p.RoomID, targetName, p)
	if target == nil {
		p.send(style(dim, "No runner named '"+targetName+"' is here.") + crlf)
		return
	}
	if p.Inv[item] <= 0 {
		p.send(style(dim, "You don't have a "+item+".") + crlf)
		return
	}
	w.consumeInv(p, item)
	target.Inv[item]++
	p.send(style(green, "You hand "+target.Name+" the "+item+".") + crlf)
	target.send(style(green, p.Name+" hands you a "+item+".") + crlf)
}

func helpText() string {
	return crlf + style(neon, "== Chrome Circuit Cowboys — commands ==") + crlf +
		"  Movement : N S E W U D  (or north/south/... or the arrow keys)\r\n" +
		"  home / rest    — RECALL to your Re-Clone Bay (~10s cast; a hit or a move breaks it)\r\n" +
		"  in / out        — step into/out of your capsule pod from Neon Alley\r\n" +
		"  look (l)        — examine your location\r\n" +
		"  map (m)         — local map: exits, and the way deeper or out\r\n" +
		"  talk            — ask a local about this level (lore/backstory)\r\n" +
		"  hack            — at a terminal: crack the access code for scrip + XP\r\n" +
		"  send / wire      — at a terminal: SEND <runner> <msg> / WIRE <runner> <scrip>; MAIL to read\r\n" +
		"  trade <runner>   — face-to-face swap: OFFER <item>/<scrip>, CONFIRM, CANCEL\r\n" +
		"  spend <stat>    — spend character points to raise body/reflexes/intelligence\r\n" +
		"  attack <foe> (a) — engage a hostile (alias kill/breach)\r\n" +
		"  flee            — try to break a fight and bolt\r\n" +
		"  say <msg>       — talk to others in the room\r\n" +
		"  me / emote / :<action> — roleplay an action (\"Wintermute lights a cig\")\r\n" +
		"  who             — who's jacked in\r\n" +
		"  score (st)      — your character sheet\r\n" +
		"  list / buy <x>  — vendor (at shops); sell <item> [qty]; use <item> to consume\r\n" +
		"  drop / get      — drop <item>/all on the floor; get <item>/all (share loot)\r\n" +
		"  open            — crack open a supply/data cache in the room\r\n" +
		"  loot (lo)       — strip a flatlined body, ICE shards, or a cache of its gear\r\n" +
		"  install <cyber> — Emergency Medic re-installs salvaged cyberware (at the Night Market)\r\n" +
		"  give <item> <runner> — hand recovered gear back to a crewmate\r\n" +
		"  inventory (i)   — what you're carrying (cap grows with level)\r\n" +
		"  stash / grab <x> — store/withdraw gear at your Re-Clone Bay (uncapped)\r\n" +
		"  quests          — fixer bounty board; accept <#> / accept 1 2 3 / accept all / claim\r\n" +
		"  programs / run <name> — netrun demons (scalpel/hammer/leech/mirror/medic)\r\n" +
		"  invite <runner> — invite to your crew (leader only); they ACCEPT/DECLINE\r\n" +
		"  group / crew     — show your crew (shared XP in-room); gsay <msg> (or ;<msg>); leave\r\n" +
		"  leaderboard     — top runners by level\r\n" +
		"  quit            — jack out\r\n" +
		style(dim, "  In the Net, ATTACK breaches ICE using Intelligence and spends RAM\r\n"+
			"  (buy a cyberdeck for more, ram-chips to refill). PvP is LIVE EVERYWHERE\r\n"+
			"  except the street outside the clone pods — draw on a runner there and a\r\n"+
			"  security drone flatlines you. Die and your body drops with your gear;\r\n"+
			"  LOOT it, re-INSTALL cyberware at a Emergency Medic. Some ICE morphs when broken.") + crlf
}
