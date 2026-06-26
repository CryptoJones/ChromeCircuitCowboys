package cowboy

import (
	"strconv"
	"strings"
)

// Quest is a story bounty: kill Count of a target mob, then CLAIM at any broker
// (vendor room) for the reward. A quest with a Giver is offered only in that room
// (its quest-giver NPC); a quest with Giver "" is a generic street bounty offered
// at the city brokers (Chrome Rose / Night Market). Bounties are repeatable.
type Quest struct {
	ID        string
	Name      string
	Desc      string
	Target    string // mob template ID
	Count     int
	XP        int
	Eddies    int
	MinLevel  int
	Giver     string // room id where this quest is offered ("" = generic street brokers)
	GiverName string // the quest-giver NPC, for flavor
	Pool      string // "ring" = a roving RP-ring rumor (scattered across ring givers, see assignRingQuests)
}

var quests = []Quest{
	// --- Generic street bounties (offered at the Chrome Rose / Night Market) ---
	{ID: "clear_alley", Name: "Clear the Back Alley", Target: "ganger", Count: 3, XP: 120, Eddies: 80, MinLevel: 1,
		Desc: "Gangers are taxing the block. Drop 3 street gangers."},
	{ID: "corp_sabotage", Name: "Corporate Sabotage", Target: "drone", Count: 2, XP: 200, Eddies: 150, MinLevel: 2,
		Desc: "A rival corp wants deniable chaos. Wreck 2 security drones in Corporate Plaza."},
	{ID: "break_ice", Name: "Break the Ice", Target: "nz1_1_mid_m", Count: 3, XP: 400, Eddies: 250, MinLevel: 3,
		Desc: "Prove you can run. Destroy 3 ICE constructs in the Net's underbelly."},
	{ID: "ghost_machine", Name: "Ghost in the Machine", Target: "gauntlet3", Count: 1, XP: 1000, Eddies: 800, MinLevel: 5,
		Desc: "Beat the reconfiguring Gauntlet ICE down to its lethal lock and shatter it."},

	// ===================== MEATSPACE — the underground descent =====================
	{ID: "ug1_snatch", Name: "The Snatch & Grab", Target: "z1_05_m", Count: 1, XP: 150, Eddies: 120, MinLevel: 1,
		Giver: "z1_01", GiverName: "Marcus the Fixer",
		Desc: "The Scrap-Hounds boosted a high-end cyberdeck. Recover it from their warboss, Razorback Kane."},
	{ID: "ug1_hounds", Name: "Thin the Pack", Target: "z1_02_m", Count: 3, XP: 180, Eddies: 100, MinLevel: 1,
		Giver: "z1_09", GiverName: "Doc 'Stitches' Vance",
		Desc: "Kurokawa's sweep runs on Scrap-Hound muscle. Drop 3 packs on the Sodium Strip."},
	{ID: "ug1_heist", Name: "The Data Heist", Target: "z1_17_m", Count: 1, XP: 500, Eddies: 350, MinLevel: 7,
		Giver: "z1_11", GiverName: "Jax",
		Desc: "The full EREBUS manifesto sits in a Kurokawa logistics hub. Jack the mainframe — put down Warden Sato's combat-frame."},
	{ID: "ug2_thorne", Name: "The Cyberspace Duel", Target: "z2_14_m", Count: 1, XP: 900, Eddies: 600, MinLevel: 11,
		Giver: "z2_01", GiverName: "Cipher",
		Desc: "Sever EREBUS at its core — face Dr. Aris Thorne in the sub-zero Net duel."},
	{ID: "ug3_praetor", Name: "The Black Site", Target: "z3_15_m", Count: 1, XP: 1350, Eddies: 900, MinLevel: 21,
		Giver: "z3_03", GiverName: "Silas",
		Desc: "Hit the off-grid Black Site and expose the culling — kill the enforcer Praetor-9."},
	{ID: "ug4_siege", Name: "Siege of the Archive", Target: "z4_14_m", Count: 1, XP: 1800, Eddies: 1200, MinLevel: 31,
		Giver: "z4_04", GiverName: "Dr. Evelyn Vance",
		Desc: "Hold the bunker and broadcast the Ascension Protocol — break the corporate Heavy-Mech Commander."},
	{ID: "ug5_overlord", Name: "The Executive Core", Target: "z5_13_m", Count: 1, XP: 2250, Eddies: 1500, MinLevel: 41,
		Giver: "z5_02", GiverName: "the Undercrew quartermaster",
		Desc: "Crack Tartarus and entomb the elite — dismantle the Kurokawa CEO's Overlord mech."},
	{ID: "ug6_meltdown", Name: "The Core Meltdown", Target: "z6_14_m", Count: 1, XP: 2700, Eddies: 1800, MinLevel: 51,
		Giver: "z6_01", GiverName: "Silas & Dr. Vance",
		Desc: "The Tartarus Loyalists are melting the core — put down their Commander on the magma catwalk."},
	{ID: "ug7_platform", Name: "The Apex Broadcast", Target: "z7_13_m", Count: 1, XP: 3150, Eddies: 2100, MinLevel: 61,
		Giver: "z7_02", GiverName: "Old Pelle",
		Desc: "Reach the surface array — down the gunship Tempest-Actual on Platform 09."},
	{ID: "ug8_god", Name: "God in the Machine", Target: "z8_13_m", Count: 1, XP: 3600, Eddies: 2400, MinLevel: 71,
		Giver: "z8_01", GiverName: "a ghost-signal fixer",
		Desc: "EREBUS has gone singular — shatter the God in the Machine with the paradox virus."},
	{ID: "ug9_overlord", Name: "The Decapitation Strike", Target: "z9_14_m", Count: 1, XP: 4050, Eddies: 2700, MinLevel: 81,
		Giver: "z9_01", GiverName: "the Coalition quartermaster",
		Desc: "Decapitate corporate command — destroy the Iron Overlord's neural bridge."},
	{ID: "ug10_loom", Name: "Welding the Sky", Target: "z10_13_m", Count: 1, XP: 4500, Eddies: 3000, MinLevel: 91,
		Giver: "z10_01", GiverName: "Wraith",
		Desc: "Topple the last pillar and weld the sky shut — bring down the Loom Masterframe."},

	// ========================= NETSPACE — the Net ascent =========================
	{ID: "net1_trace", Name: "The GigaMesh Ledger", Target: "nz1_5_bot_m", Count: 1, XP: 500, Eddies: 350, MinLevel: 1,
		Giver: "nz1_1_top", GiverName: "Fixer-7",
		Desc: "Seize the Syndicate's full ledger — burn the Active-ICE warden Tracewright in the Black Spire."},
	{ID: "net2_arbiter", Name: "Pick a Patron", Target: "nz2_5_bot_m", Count: 1, XP: 900, Eddies: 600, MinLevel: 11,
		Giver: "nz2_1_top", GiverName: "Mr. Lattice",
		Desc: "Settle the proxy war — breach the Sundered Arbiter and seal your alliance."},
	{ID: "net3_warden", Name: "The Deep Infrastructure", Target: "nz3_5_bot_m", Count: 1, XP: 1350, Eddies: 900, MinLevel: 21,
		Giver: "nz3_1_top", GiverName: "Ravel",
		Desc: "Seize the foundation of the deep Net — shatter WARDEN-PRIME and decide Echo-9's fate."},
	{ID: "net4_overseer", Name: "The Silent Throne", Target: "nz4_5_bot_m", Count: 1, XP: 1800, Eddies: 1200, MinLevel: 31,
		Giver: "nz4_1_top", GiverName: "a counter-intel fixer",
		Desc: "Stop the blackout — rewrite the Rogue Overseer's core before two worlds go dark."},
	{ID: "net5_catalyst", Name: "The Catalyst", Target: "nz5_5_bot_m", Count: 1, XP: 2250, Eddies: 1500, MinLevel: 41,
		Giver: "nz5_1_top", GiverName: "a First Network echo",
		Desc: "Reach the Master Protocol — defeat the Prime Architect at the Catalyst Core."},
	{ID: "net6_architects", Name: "The Architect's Trial", Target: "nz6_5_bot_m", Count: 1, XP: 2700, Eddies: 1800, MinLevel: 51,
		Giver: "nz6_1_top", GiverName: "an Architect-cipher defector",
		Desc: "Face the makers of the Net — overwrite the Genesis Protocol Architects."},
	{ID: "net7_entropy", Name: "The Multiversal Leak", Target: "nz7_5_bot_m", Count: 1, XP: 3150, Eddies: 2100, MinLevel: 61,
		Giver: "nz7_1_top", GiverName: "the Cosmic Arbiter relay",
		Desc: "Seal the rift between living and dead universes — weave the firewall through the Entropy-Titan."},
	{ID: "net8_ancient", Name: "The Grand Strategy", Target: "nz8_5_bot_m", Count: 1, XP: 3600, Eddies: 2400, MinLevel: 71,
		Giver: "nz8_1_top", GiverName: "the Last Cartographer",
		Desc: "Keep the multiverse infinite — out-breach the Reconciled Ancient on the game-board."},
	{ID: "net9_unmaking", Name: "Siege of the Forge", Target: "nz9_5_bot_m", Count: 1, XP: 4050, Eddies: 2700, MinLevel: 81,
		Giver: "nz9_1_top", GiverName: "the Eternal Mentor",
		Desc: "Defend every universe you authored — seal The Great Unmaking at the Forge's core."},
	{ID: "net10_final", Name: "The Ultimate Gift", Target: "nz10_5_bot_m", Count: 1, XP: 4500, Eddies: 3000, MinLevel: 91,
		Giver: "nz10_1_top", GiverName: "the Absolute Codex",
		Desc: "Finish the ascent — overwrite the Final Compilation, your own source code, and broadcast the Codex."},

	// ===================== RP RINGS — roving "rumor" bounties =====================
	// Low-level, standalone, repeatable; Pool "ring" so they're scattered randomly
	// across the ring givers each session. Each carries an easter egg nodding to one
	// of the two PvE paths.
	{ID: "ring_busker", Pool: "ring", Name: "Busker's Lost Loop", Target: "ring_scavver", Count: 1, XP: 60, Eddies: 40, MinLevel: 1,
		Desc: "A scrap-drone swallowed the busker's favorite data-loop on the belt. Smash it open and get it back — the jingle glitches into a Kurokawa logo, which is just weird."},
	{ID: "ring_tagwar", Pool: "ring", Name: "Tag War Truce", Target: "ring_ganger", Count: 3, XP: 90, Eddies: 55, MinLevel: 2,
		Desc: "Taggers are carving up the Sprawlbelt. Drop 3 to hold the truce. One swears the real muscle is 'down in the Wasteland, running with the Scrap-Hounds.'"},
	{ID: "ring_wireheads", Pool: "ring", Name: "Wirehead Roundup", Target: "ring_junkie", Count: 3, XP: 70, Eddies: 40, MinLevel: 1,
		Desc: "Clear 3 wireheads blissed out on a leaked braindance loop. The medic says the footage is stolen surveillance — branded EREBUS."},
	{ID: "ring_ghostwire", Pool: "ring", Name: "Ghost on the Wire", Target: "ring_scavver", Count: 2, XP: 110, Eddies: 70, MinLevel: 3,
		Desc: "Something jacked out of a data-port and rode the ring as a scrap-drone swarm. Purge 2. The bartender mutters it 'smells like a GigaMesh job.'"},
	{ID: "ring_palm", Pool: "ring", Name: "The Palm Reading", Target: "ring_hustler", Count: 1, XP: 60, Eddies: 45, MinLevel: 1,
		Desc: "A hustler swiped the fortune-teller's 'lucky' cortex-chip. Get it back and she'll read your fate: 'you will descend, and you will ascend.'"},
	{ID: "ring_noodle", Pool: "ring", Name: "Noodle Run", Target: "ring_ganger", Count: 2, XP: 85, Eddies: 50, MinLevel: 2,
		Desc: "The cook's broth-cache got jacked by taggers. Recover it — drop 2. He grumbles the Night Market's gone soft since the Strip got fancy."},
	{ID: "ring_cores", Pool: "ring", Name: "Core Salvage", Target: "ring_scavver", Count: 4, XP: 150, Eddies: 100, MinLevel: 4,
		Desc: "The scrap-trader wants 4 drone-cores. He swears these models 'crawled up out of the Sump' — pull them and cash in."},
	{ID: "ring_maglev", Pool: "ring", Name: "The Maglev Ghost", Target: "ring_junkie", Count: 2, XP: 80, Eddies: 45, MinLevel: 2,
		Desc: "Two wireheads keep spooking riders: a runner jacked in at the data-port and 'never came back to his body.' Quiet them down."},
	{ID: "ring_promenade", Pool: "ring", Name: "Peace on the Promenade", Target: "ring_hustler", Count: 2, XP: 65, Eddies: 45, MinLevel: 1,
		Desc: "The Rolling Rose's bouncer wants 2 hustlers off the belt before they fleece the regulars blind."},
	{ID: "ring_lostfound", Pool: "ring", Name: "Lost & Found", Target: "ring_scavver", Count: 1, XP: 60, Eddies: 40, MinLevel: 1,
		Desc: "A kid's drone-toy got swept into a scrap-drone. Recover it — it's a tiny replica of some giant 'sky-loom' the kid swears she saw in a dream."},
	{ID: "ring_tab", Pool: "ring", Name: "The Bartender's Tab", Target: "ring_ganger", Count: 1, XP: 80, Eddies: 60, MinLevel: 2,
		Desc: "Collect the Rolling Rose's tab from a tagger who skipped out. He blames 'a bad run against some reconfiguring Gauntlet ICE.'"},
	{ID: "ring_patrol", Pool: "ring", Name: "Belt Patrol", Target: "ring_junkie", Count: 3, XP: 100, Eddies: 60, MinLevel: 3,
		Desc: "Keep the loop calm — clear 3 belt strays so the RP crowd can breathe and the rumors keep flowing."},
}

// questsHere returns the bounties offered to p in their current room: the
// quest-giver NPCs homed here, plus the generic street bounties when standing at
// a city broker (the Chrome Rose / Night Market). Order matches the displayed
// board, so ACCEPT <#> indexes this same list.
func (w *World) questsHere(p *Player) []Quest {
	room := p.RoomID
	streetBroker := room == "chrome_bar" || room == "market"
	var out []Quest
	// Roving RP-ring rumors scattered to this giver this session.
	for _, idx := range w.ringOffer[room] {
		out = append(out, quests[idx])
	}
	// Story givers + generic street bounties (ring-pool quests are surfaced only
	// via ringOffer above, never here).
	for _, q := range quests {
		if q.Pool == "ring" {
			continue
		}
		if q.Giver == room || (q.Giver == "" && streetBroker) {
			out = append(out, q)
		}
	}
	return out
}

func questByID(id string) (Quest, bool) {
	for _, q := range quests {
		if q.ID == id {
			return q, true
		}
	}
	return Quest{}, false
}

// quests command: show the bounties offered here (by the local fixer/NPC) plus
// your active progress.
func (w *World) showQuests(p *Player) {
	offered := w.questsHere(p)
	if len(offered) > 0 {
		header := "== FIXER BOUNTY BOARD =="
		if g := offered[0].GiverName; g != "" {
			header = "== " + g + " has work =="
		}
		p.send(crlf + style(neon, header+"  ") + style(dim, "(ACCEPT <#>)") + crlf)
		for i, q := range offered {
			lvl := ""
			if p.Level < q.MinLevel {
				lvl = style(red, "  [needs level "+itoa(q.MinLevel)+"]")
			}
			// Colour by the player's state with this bounty: greyed if already
			// accepted and in progress, RED when it's complete and ready to turn
			// in here, normal (green) if not yet taken.
			nameColor, tag := green, ""
			if got, active := p.Quests[q.ID]; active {
				if got >= q.Count {
					nameColor = red
					tag = style(red, "  [READY — turn in]")
				} else {
					nameColor = dim
					tag = style(dim, "  [accepted "+itoa(got)+"/"+itoa(q.Count)+"]")
				}
			}
			p.send("  " + style(gold, itoa(i+1)+")") + " " + style(nameColor, q.Name) +
				style(dim, " — "+q.Desc) + " " + style(gold, "(+"+itoa(q.XP)+"xp, €$"+itoa(q.Eddies)+")") + lvl + tag + crlf)
		}
	}
	if len(p.Quests) == 0 {
		p.send(style(dim, "You have no active bounties.") + crlf)
		return
	}
	p.send(style(neon, "-- Active bounties --") + crlf)
	for id, got := range p.Quests {
		q, ok := questByID(id)
		if !ok {
			continue
		}
		state := itoa(got) + "/" + itoa(q.Count)
		if got >= q.Count {
			state = style(gold, "READY — CLAIM at a broker or its giver")
		}
		p.send("  " + style(green, q.Name) + style(dim, " ["+q.Target+"] ") + state + crlf)
	}
}

// accept takes bounties offered in the current room (by the local fixer/NPC).
// Accepts one ("accept 2"), several ("accept 1 2 3"), or every eligible job
// ("accept all"). Each pick is guarded independently (level / already-on-it);
// a bad token is skipped with a note rather than aborting the batch.
func (w *World) accept(p *Player, arg string) {
	offered := w.questsHere(p)
	if len(offered) == 0 {
		p.send(style(dim, "No one here is hiring. Find a fixer or a broker (a vendor room).") + crlf)
		return
	}
	arg = strings.ToLower(strings.TrimSpace(arg))
	if arg == "" {
		p.send(style(dim, "Accept which? See QUESTS for the numbered board (or ACCEPT ALL).") + crlf)
		return
	}
	var picks []int
	if arg == "all" {
		for i := range offered {
			picks = append(picks, i+1)
		}
	} else {
		for _, tok := range strings.Fields(arg) {
			n, err := strconv.Atoi(tok)
			if err != nil || n < 1 || n > len(offered) {
				p.send(style(dim, "Ignoring \""+tok+"\" — not a bounty number on the board.") + crlf)
				continue
			}
			picks = append(picks, n)
		}
	}
	accepted := 0
	seen := map[int]bool{}
	for _, n := range picks {
		if seen[n] {
			continue
		}
		seen[n] = true
		q := offered[n-1]
		if p.Level < q.MinLevel {
			p.send(style(red, q.Name+": you need level "+itoa(q.MinLevel)+" for that job.") + crlf)
			continue
		}
		if _, active := p.Quests[q.ID]; active {
			p.send(style(dim, "Already on "+q.Name+".") + crlf)
			continue
		}
		// One-time bounties can't be repeated; the RP-ring rumors are exempt.
		if q.Pool != "ring" && p.Done[q.ID] > 0 {
			p.send(style(dim, q.Name+": you've already done that job.") + crlf)
			continue
		}
		p.Quests[q.ID] = 0
		accepted++
		p.send(style(green, "Bounty accepted: ") + q.Name + style(dim, " — "+q.Desc) + crlf)
	}
	if accepted > 1 {
		p.send(style(gold, "Took on "+itoa(accepted)+" bounties.") + crlf)
	}
}

// claim turns in any completed bounties for rewards. A bounty can be redeemed
// at a broker (a vendor room) OR back with the quest-giver who offered it (the
// fixer's room, or — for roving RP-ring rumors — wherever it was scattered this
// session).
func (w *World) claim(p *Player) {
	atVendor := w.atVendor(p)
	claimed := 0
	for id, got := range p.Quests {
		q, ok := questByID(id)
		if !ok || got < q.Count {
			continue
		}
		if !atVendor && !w.questOfferedHere(p, q) {
			continue // not at a broker, and not at this bounty's giver
		}
		delete(p.Quests, id)
		p.XP += q.XP
		p.Eddies += q.Eddies
		claimed++
		// Story/street bounties are one-time; mark them done so they can't be
		// re-accepted. The RP-ring rumors stay repeatable.
		if q.Pool != "ring" {
			if p.Done == nil {
				p.Done = map[string]int{}
			}
			p.Done[q.ID] = 1
		}
		p.send(style(gold, "*** Bounty paid: "+q.Name+" — +"+itoa(q.XP)+"xp, €$"+itoa(q.Eddies)+" ***") + crlf)
	}
	if claimed == 0 {
		p.send(style(dim, "No completed bounties to claim here. Return to a broker or the bounty's giver.") + crlf)
		return
	}
	w.checkLevelUp(p)
}

// questOfferedHere reports whether the player is standing with the quest-giver
// who offers q — its fixed Giver room, or (for ring rumors) a ring-giver room
// this quest was scattered to this session.
func (w *World) questOfferedHere(p *Player, q Quest) bool {
	if q.Giver != "" && p.RoomID == q.Giver {
		return true
	}
	for _, idx := range w.ringOffer[p.RoomID] {
		if quests[idx].ID == q.ID {
			return true
		}
	}
	return false
}

// creditQuestKill advances any active bounty whose target matches the slain mob.
func (w *World) creditQuestKill(p *Player, mobID string) {
	for id, got := range p.Quests {
		q, ok := questByID(id)
		if !ok || q.Target != mobID || got >= q.Count {
			continue
		}
		p.Quests[id] = got + 1
		if p.Quests[id] >= q.Count {
			p.send(style(gold, "Bounty objective complete: "+q.Name+" — CLAIM at a broker or its giver.") + crlf)
		} else {
			p.send(style(dim, "Bounty progress: "+q.Name+" "+itoa(p.Quests[id])+"/"+itoa(q.Count)) + crlf)
		}
	}
}
