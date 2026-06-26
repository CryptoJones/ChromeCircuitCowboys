package cowboy

import "strings"

// breakRecall interrupts a runner's HOME recall when they take a hit.
func (w *World) breakRecall(p *Player) {
	if p.homing > 0 {
		p.homing = 0
		p.send(style(red, "Your recall shatters — you're under fire!") + crlf)
	}
}

// engage targets a hostile (a mob, or — in the Net — another runner for PvP) and
// starts a fight. Rounds resolve on Tick, MajorMUD-style.
func (w *World) engage(p *Player, arg string) {
	arg = strings.ToLower(strings.TrimSpace(arg))
	if p.homing > 0 { // throwing the first punch breaks your own recall
		p.homing = 0
		p.send(style(dim, "Your recall breaks as you move on a target.") + crlf)
	}

	// Targeting another runner? PvP is live everywhere EXCEPT the safe zone
	// outside the clone pods — draw there and a security drone flatlines you.
	if arg != "" {
		if target := w.playerInRoomByName(p.RoomID, arg, p); target != nil {
			if w.pvpAllowed(p) {
				w.engagePvP(p, target)
			} else {
				p.send(style(red, "You move on "+target.Name+" — but this is a no-violence zone.") + crlf)
				w.broadcast(p.RoomID, p, style(hot, "A security drone locks onto "+p.Name+" for assault and opens fire!")+crlf)
				target.send(style(dim, "A security drone drops "+p.Name+" before they reach you.") + crlf)
				w.securityKill(p)
			}
			return
		}
	}

	mobs := w.liveMobsIn(p.RoomID)
	if len(mobs) == 0 {
		p.send(style(dim, "Nothing here to fight.") + crlf)
		return
	}
	var target *Mob
	if arg == "" {
		target = mobs[0]
	} else {
		for _, m := range mobs {
			if strings.Contains(strings.ToLower(m.tmpl.Name), arg) || strings.Contains(m.tmpl.ID, arg) {
				target = m
				break
			}
		}
	}
	if target == nil {
		p.send(style(dim, "You don't see '"+arg+"' here.") + crlf)
		return
	}
	p.pvpTarget = nil // mob combat supersedes any duel
	p.fighting = target
	if target.target == nil {
		target.target = p
	}
	verb := "You lunge at "
	switch {
	case target.tmpl.Container:
		// You don't "lunge at" a crate — pry/crack/jimmy it, varied each time.
		verb = containerVerbs[w.roll(len(containerVerbs))]
	case w.inNet(p):
		verb = "You jack a breach protocol into "
	}
	p.send(style(hot, verb+target.tmpl.Name+"!") + crlf)
	w.broadcast(p.RoomID, p, style(dim, p.Name+" attacks "+target.tmpl.Name+".")+crlf)
}

// containerVerbs are the (randomized) ways you go at an inert cache — a
// gentleman pries, he does not "lunge at" a crate. (#46)
var containerVerbs = []string{
	"You pry at ", "You crack open ", "You jimmy ", "You force ",
	"You bash at ", "You wrench at ", "You crowbar ", "You lever open ",
}

func (w *World) engagePvP(p, target *Player) {
	if target == p {
		p.send(style(dim, "You can't jack yourself.") + crlf)
		return
	}
	p.fighting = nil
	p.pvpTarget = target
	p.send(style(hot, "You jack into "+target.Name+"'s deck — netrun duel!") + crlf)
	target.send(style(hot, p.Name+" jacks into your deck — defend yourself!") + crlf)
	w.broadcast(p.RoomID, p, style(dim, p.Name+" and "+target.Name+" are duelling in the Net.")+crlf)
}

func (w *World) playerInRoomByName(roomID, name string, except *Player) *Player {
	name = strings.ToLower(name)
	for _, o := range w.playersIn(roomID, except) {
		if strings.ToLower(o.Name) == name || strings.HasPrefix(strings.ToLower(o.Name), name) {
			return o
		}
	}
	return nil
}

// flee attempts to break combat (mob or duel) and bolt to a random exit.
func (w *World) flee(p *Player) {
	if p.fighting == nil && p.pvpTarget == nil {
		p.send(style(dim, "You're not in combat.") + crlf)
		return
	}
	if w.roll(2) != 0 {
		p.send(style(red, "You can't break the connection — the fight holds you!") + crlf)
		return
	}
	if p.fighting != nil && p.fighting.target == p {
		p.fighting.target = nil
	}
	p.fighting = nil
	p.pvpTarget = nil
	r := w.room(p.RoomID)
	var exits []string
	for _, d := range []string{"north", "south", "east", "west", "up", "down"} {
		if _, ok := r.Exits[d]; ok {
			exits = append(exits, d)
		}
	}
	p.send(style(green, "You rip free and bolt!") + crlf)
	if len(exits) > 0 {
		dir := exits[w.roll(len(exits))]
		w.broadcast(p.RoomID, p, style(dim, p.Name+" flees "+dir+".")+crlf)
		p.RoomID = r.Exits[dir]
		w.broadcast(p.RoomID, p, style(dim, p.Name+" skids in, breathless.")+crlf)
		w.lookText(p)
	}
}

// Tick advances the world one round: aggro, mob fights, PvP duels, deaths,
// respawns, and out-of-combat regen (HP and RAM).
func (w *World) Tick() {
	w.aggro()
	w.resolveCombat()
	w.resolvePvP()
	w.tickRecall() // after combat, so a hit this tick interrupts before the recall lands
	w.respawnDead()
	w.expireShields()
	w.regen()
}

// tickRecall counts down HOME recalls and, when one completes, phases the runner
// back to their Re-Clone Bay. A hit (resolveCombat/resolvePvP) or a move zeroes
// homing before this runs, so only an uninterrupted cast lands.
func (w *World) tickRecall() {
	for _, p := range w.players {
		if p.homing <= 0 {
			continue
		}
		p.homing--
		if p.homing <= 0 {
			w.broadcast(p.RoomID, p, style(dim, p.Name+" phases out in a wash of static.")+crlf)
			p.RoomID = startRoom
			p.send(style(neon, "*** Recall complete — you blink into your Re-Clone Bay. ***") + crlf)
			w.lookText(p)
		}
	}
}

// expireShields counts down the Mirror program's damage shield.
func (w *World) expireShields() {
	for _, p := range w.players {
		if p.shieldTicks > 0 {
			p.shieldTicks--
			if p.shieldTicks == 0 {
				p.shieldAmt = 0
				p.send(style(dim, "Your Mirror deflector fades.") + crlf)
			}
		}
	}
}

func (w *World) aggro() {
	for _, m := range w.mobs {
		if m.dead || !m.tmpl.Aggressive || m.target != nil {
			continue
		}
		victims := w.playersIn(m.RoomID, nil)
		if len(victims) == 0 {
			continue
		}
		v := victims[w.roll(len(victims))]
		m.target = v
		if v.fighting == nil && v.pvpTarget == nil {
			v.fighting = m
		}
		v.send(style(hot, m.tmpl.Name+" locks onto you and attacks!") + crlf)
	}
}

// playerSwing returns the player's damage for one round, route-aware. In the Net
// a breach spends 1 RAM for full Intelligence-powered damage; with no RAM the
// breach sputters at half strength.
func (w *World) playerSwing(p *Player) int {
	atk := w.effAttack(p)
	if w.inNet(p) {
		if p.RAM > 0 {
			p.RAM--
		} else {
			atk /= 2
			if atk < 1 {
				atk = 1
			}
			p.send(style(dim, "Low RAM — your breach sputters at half power.") + crlf)
		}
	}
	return atk
}

func (w *World) resolveCombat() {
	// Pass 1: every engaged runner swings at the mob they're fighting. Several
	// runners can pile on the same mob — each lands their own hit this tick.
	for _, p := range w.players {
		m := p.fighting
		if m == nil {
			continue
		}
		if m.dead || m.RoomID != p.RoomID {
			p.fighting = nil
			continue
		}
		if w.toHit(p.Reflexes, m.tmpl.AC) {
			d := dmg(w.playerSwing(p), m.tmpl.AC)
			m.HP -= d
			p.send(style(green, "You hit "+m.tmpl.Name+" for "+itoa(d)+".") + crlf)
		} else {
			p.send(style(dim, "You miss "+m.tmpl.Name+".") + crlf)
		}
		if m.HP <= 0 {
			w.killMob(p, m)
		}
	}
	// Pass 2: every live mob that's being fought hits back at ONE of its
	// attackers, picked at random each tick — so a group fight spreads the
	// danger around instead of pinning one runner.
	for _, m := range w.mobs {
		if m.dead {
			continue
		}
		attackers := w.attackersOf(m)
		if len(attackers) == 0 {
			continue
		}
		target := attackers[w.roll(len(attackers))]
		m.target = target
		if w.toHit(m.tmpl.Damage/2, playerAC(target)) {
			d := applyShield(target, dmg(m.tmpl.Damage, playerAC(target)))
			target.HP -= d
			w.breakRecall(target)
			target.send(style(red, m.tmpl.Name+" hits you for "+itoa(d)+".") + crlf)
			if target.HP <= 0 {
				w.flatline(target, m)
			}
		} else {
			target.send(style(dim, m.tmpl.Name+" misses you.") + crlf)
		}
	}
}

// attackersOf returns every runner currently fighting mob m in its room.
func (w *World) attackersOf(m *Mob) []*Player {
	var out []*Player
	for _, p := range w.players {
		if p.fighting == m && p.RoomID == m.RoomID {
			out = append(out, p)
		}
	}
	return out
}

// resolvePvP runs one round of every active netrun duel. Both runners swing in
// the same tick (each processed on their turn).
func (w *World) resolvePvP() {
	for _, p := range w.players {
		d := p.pvpTarget
		if d == nil {
			continue
		}
		if !w.pvpAllowed(p) || d.RoomID != p.RoomID || w.players[d.ID] == nil {
			p.pvpTarget = nil
			p.send(style(dim, "Your duel target is gone.") + crlf)
			continue
		}
		if w.toHit(p.Reflexes, playerAC(d)) {
			hit := applyShield(d, dmg(w.playerSwing(p), playerAC(d)))
			d.HP -= hit
			w.breakRecall(d)
			p.send(style(green, "You breach "+d.Name+"'s deck for "+itoa(hit)+".") + crlf)
			d.send(style(red, p.Name+" breaches your deck for "+itoa(hit)+".") + crlf)
			if d.HP <= 0 {
				w.pvpFlatline(p, d)
			}
		} else {
			p.send(style(dim, d.Name+" slips your breach.") + crlf)
		}
	}
}

func (w *World) killMob(p *Player, m *Mob) {
	// Multi-stage ICE: morph into the next, harder form instead of dying.
	if m.tmpl.Next != "" {
		if next, ok := w.tmpls[m.tmpl.Next]; ok {
			old := m.tmpl.Name
			m.tmpl = next
			m.HP = next.HP
			m.target = p
			p.fighting = m
			p.send(style(hot, "The ICE reconfigures! "+old+" collapses into "+next.Name+"!") + crlf)
			w.broadcast(p.RoomID, p, style(dim, old+" reconfigures into "+next.Name+".")+crlf)
			return
		}
	}
	m.dead = true
	m.HP = 0
	m.awaitingLoot = true // respawn stays GATED until someone loots the body
	if m.target != nil {
		m.target.fighting = nil
		m.target = nil
	}
	p.fighting = nil

	// Rewards are randomized (75%-125% of the template). XP is awarded on the
	// kill; the scrip drops with the body and must be LOOTed — and looting is
	// what lets the area respawn the mob.
	xp := w.vary(m.tmpl.XP)
	scrip := w.vary(m.tmpl.Eddies)
	if pct := w.rewardPct(p); pct > 100 { // party / clan reward bonus (#44)
		xp = xp * pct / 100
		scrip = scrip * pct / 100
		p.send(style(neon, "(crew bonus: rewards ×"+itoa(pct)+"%)") + crlf)
	}
	ice := m.tmpl.ICE
	box := m.tmpl.Container
	mech := m.tmpl.Mechanical
	loot := map[string]int{}
	for k, v := range m.tmpl.Drops { // seeded item drops (loot caches, boss kit)
		loot[k] = v
	}
	if !box { // caches are their own loot; for real kills, cut the crew in
		w.addPartyLoot(p, loot)
	}
	w.corpses = append(w.corpses, &Corpse{
		Owner: m.tmpl.Name, RoomID: m.RoomID, Loot: loot, Scrip: scrip, mob: m, IsICE: ice, IsBox: box, IsMech: mech,
	})
	p.send(style(hot, "*** "+m.tmpl.Name+" is destroyed! ***") + crlf)
	drop := "Its body drops €$" + itoa(scrip) + " scrip - LOOT it."
	switch {
	case box:
		drop = "It cracks open, spilling €$" + itoa(scrip) + " scrip - LOOT it."
	case ice:
		drop = "It shatters into broken shards holding €$" + itoa(scrip) + " scrip - LOOT them."
	case mech:
		drop = "Its frame sparks and goes dark, scattering €$" + itoa(scrip) + " scrip - LOOT it."
	}
	p.send(style(gold, "You gain "+itoa(xp)+" XP. "+drop) + crlf)
	w.broadcast(p.RoomID, p, style(dim, p.Name+" destroys "+m.tmpl.Name+".")+crlf)
	w.creditQuestKill(p, m.tmpl.ID)
	w.awardXP(p, xp) // XP shared with crew in the room; handles level-ups
}

// vary returns base scaled by a random ~75%-125% so kill rewards aren't fixed.
// It uses the injectable roll, so tests stay deterministic. Non-positive bases
// (e.g. the non-final Gauntlet stages) pass straight through.
func (w *World) vary(base int) int {
	if base <= 0 {
		return base
	}
	return base - base/4 + w.roll(base/2+1)
}

// flatline handles player death by a mob: half HP, respawn at the start, and a
// credit/XP penalty — never permadeath.
func (w *World) flatline(p *Player, killer *Mob) {
	fee := p.Eddies / 10
	p.send(style(red, "*** FLATLINE — "+killer.tmpl.Name+" wastes your body. ***") + crlf)
	p.send(style(neon, "Your mind restores into a fresh clone at the clinic. ") +
		style(gold, "Clone fee: €$"+itoa(fee)) + style(neon, ".") + crlf)
	if killer.target == p {
		killer.target = nil
	}
	w.reClone(p, fee)
}

// pvpAllowed reports whether one runner may attack another in this room. PvP is
// live everywhere except no-violence safe zones (the street outside the pods)
// and private capsule bays.
func (w *World) pvpAllowed(p *Player) bool {
	r := w.room(p.RoomID)
	return r != nil && !r.Safe && !r.Private
}

// securityKill is the safe-zone enforcer: a runner who draws on another in a
// no-violence zone is flatlined on the spot by a city security drone (drops their body
// and pays the clone fee, like any death).
func (w *World) securityKill(p *Player) {
	fee := p.Eddies / 10
	p.send(style(red, "*** A City Security drone flatlines you for assault in a no-violence zone. ***") + crlf)
	p.send(style(neon, "Your mind restores into a fresh clone. ") +
		style(gold, "Clone fee: €$"+itoa(fee)) + style(neon, ".") + crlf)
	w.reClone(p, fee)
}

// reClone runs the shared death sequence: clear combat, drop the old body as
// a corpse, hand off crew leadership, and respawn the fresh clone (charging fee).
func (w *World) reClone(p *Player, fee int) {
	p.fighting = nil
	p.pvpTarget = nil
	w.dropCorpse(p)
	w.passLeadershipOnDeath(p)
	w.respawnPlayer(p, fee)
	// Fresh-sleeve adrenaline: a new clone wakes up with a +15% health buffer
	// (overheal above max) to take the sting off re-sleeving — it bleeds off as
	// you take damage and never regenerates back above max.
	if bonus := p.MaxHP * 15 / 100; bonus > 0 {
		p.HP += bonus
		p.send(style(gold, "Fresh-sleeve adrenaline: +"+itoa(bonus)+" HP ("+itoa(p.HP)+"/"+itoa(p.MaxHP)+") to shake off the re-sleeve.") + crlf)
	}
}

// dropCorpse leaves the dead runner's old body where they fell, holding all
// the gear the fresh clone woke up without: every inventory item PLUS their
// cyberware (weapon + deck), which is stripped from the clone. The corpse stays
// until someone loots it. Must run before respawnPlayer (which moves the room).
func (w *World) dropCorpse(p *Player) {
	loot := map[string]int{}
	for k, v := range p.Inv {
		loot[k] = v
	}
	p.Inv = map[string]int{}
	if p.WeaponName != "" && p.WeaponBonus > 0 {
		loot[p.WeaponName]++ // weapon implant stays with the old body
		p.WeaponName, p.WeaponBonus = "", 0
	}
	if p.DeckBonus > 0 {
		loot["cyberdeck"]++ // deck implant stays with the old body
		p.DeckBonus = 0
		if p.RAM > maxRAM(p) {
			p.RAM = maxRAM(p)
		}
	}
	if len(loot) == 0 {
		return
	}
	w.corpses = append(w.corpses, &Corpse{Owner: p.Name, RoomID: p.RoomID, Loot: loot})
	w.broadcast(p.RoomID, nil, style(dim, p.Name+"'s flatlined body crumples to the ground, gear and all. (LOOT)")+crlf)
}

// passLeadershipOnDeath hands the crew to the longest-tenured surviving member
// when its leader flatlines — a dead runner doesn't keep leading. Members is in
// join order, so the longest-tenured survivor is the first member that isn't p.
func (w *World) passLeadershipOnDeath(p *Player) {
	if p.party == nil || p.party.Leader != p {
		return
	}
	for _, m := range p.party.Members {
		if m != p {
			p.party.Leader = m
			p.party.broadcast(style(dim, p.Name+" flatlined — "+m.Name+" now leads the crew.") + crlf)
			return
		}
	}
}

// pvpFlatline handles losing a netrun duel: the winner siphons a cut of the
// loser's scrip (data theft) and the loser respawns in meatspace.
func (w *World) pvpFlatline(winner, loser *Player) {
	// In a sparring gym a "kill" is just a knockout: the loser yields, wakes at
	// full HP, and keeps everything. No re-sleeve, no corpse, no siphon.
	if r := w.room(loser.RoomID); r != nil && r.Spar {
		loser.HP = loser.MaxHP
		loser.pvpTarget = nil
		winner.pvpTarget = nil
		loser.send(style(hot, "*** You're dropped — you tap out. Just a spar; you wake up sore but whole. ***") + crlf)
		winner.send(style(gold, "*** You win the spar! "+loser.Name+" taps out. ***") + crlf)
		w.broadcast(loser.RoomID, nil, style(dim, winner.Name+" wins a sparring match against "+loser.Name+".")+crlf)
		return
	}
	loot := loser.Eddies / 10
	loser.send(style(red, "*** YOUR DECK IS FRIED — "+winner.Name+" flatlines you and siphons €$"+itoa(loot)+". ***") + crlf)
	winner.send(style(gold, "*** You fry "+loser.Name+"'s deck and siphon €$"+itoa(loot)+"! ***") + crlf)
	winner.Eddies += loot
	winner.pvpTarget = nil
	w.reClone(loser, loot)
}

// respawnPlayer re-clones a defeated runner: the neural backup restores from
// its realtime backup into a FRESH, full-HP clone at the clone facility. The
// only cost is the clone-body fee (`fee` scrip, ~10% of credits) — no XP or
// skill loss, since the stack is intact. (Cyberware staying with the old body
// is handled separately by the corpse system.)
func (w *World) respawnPlayer(p *Player, fee int) {
	p.Eddies -= fee
	if p.Eddies < 0 {
		p.Eddies = 0
	}
	p.HP = p.MaxHP // fresh clone, full health
	p.RoomID = startRoom
	w.lookText(p)
}

func (w *World) respawnDead() {
	for _, m := range w.mobs {
		if !m.dead {
			continue
		}
		if m.awaitingLoot {
			continue // body not yet looted — hold the respawn (set by loot)
		}
		m.respawnIn--
		if m.respawnIn <= 0 {
			m.dead = false
			m.tmpl = m.origin // reset multi-stage ICE back to its first form
			m.HP = m.tmpl.HP
			m.RoomID = m.tmpl.Home
			m.target = nil
			w.broadcast(m.RoomID, nil, style(dim, m.tmpl.Name+" reinitializes.")+crlf)
		}
	}
}

func (w *World) regen() {
	for _, p := range w.players {
		inCombat := p.fighting != nil || p.pvpTarget != nil
		if inCombat {
			continue
		}
		if p.HP < p.MaxHP {
			heal := p.MaxHP / 20
			if heal < 1 {
				heal = 1
			}
			p.HP += heal
			if p.HP > p.MaxHP {
				p.HP = p.MaxHP
			}
		}
		if mr := maxRAM(p); p.RAM < mr {
			p.RAM++
		}
	}
}
