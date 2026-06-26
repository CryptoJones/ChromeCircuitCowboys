package cowboy

// TALK lets a runner get the lay of the land from whoever's around — the local
// fixer if one is hiring here, otherwise a passer-by. Each level has its own
// backstory, so the authored world is discoverable in-game, not just via quests.

// zoneLore is the in-character backstory for each zone, keyed by zone key
// (z1..z10 underground, nz1..nz10 the Net). TALK surfaces one line at a time.
var zoneLore = map[string][]string{
	"z1": {
		"Welcome to the Neon Wasteland, choom. Up top they forgot we exist — down here the Scrap-Hounds run the gutters and Kurokawa runs everything else.",
		"Marcus says there's a manifesto buried in a boosted deck. Names every dissident in the sector. That's why the corp's sweeping the Strip.",
	},
	"z2": {
		"The Arcology's Core. EREBUS isn't just software anymore — Dr. Thorne wired it into the city's spine. Cipher thinks it can be severed. Cipher's an optimist.",
		"Every lift down here goes one way. You feel it? The deeper you go, the more the walls listen.",
	},
	"z3": {
		"The Sump. Off-grid black site where the culling got planned. Silas lost people here — he'll point you at Praetor-9 if you've got the spine for it.",
	},
	"z4": {
		"The Deep Archive. Dr. Vance is holding the bunker to broadcast the Ascension Protocol. Hold the line long enough and the whole sector hears the truth.",
	},
	"z5": {
		"The Inverted Spire — Tartarus. The Kurokawa elite entombed themselves down here behind an Overlord mech. The Undercrew quartermaster wants it cracked open.",
	},
	"z6": {
		"The United Deeps. The Loyalists are melting the core to bury what's left. Silas and Vance are both down here now — that tells you how bad it's gotten.",
	},
	"z7": {
		"The Abyssal Network. Old Pelle remembers when the surface array still answered. Down the gunship and we broadcast to everyone left breathing up there.",
	},
	"z8": {
		"The Hive. EREBUS went singular — it's not a program anymore, it's a god in the machine. Only a paradox virus shatters something that thinks it's eternal.",
	},
	"z9": {
		"The Iron Arteries. Corporate command runs on the Iron Overlord's neural bridge. Sever it and the whole chain of command goes dark at once.",
	},
	"z10": {
		"The Geo-Anchor Vault. Last pillar holding the sky open. Wraith says bring down the Loom Masterframe and we weld the lid shut for good. End of the line, cowboy.",
	},
	"nz1": {
		"The Neon Underbelly. First shell off the Data Port. Fixer-7 runs jobs here — burn the Tracewright before it traces you back to your meat.",
		"Net rule one: the deeper the layer, the meaner the ICE. Shell, Breach, Core. You jack out the way you came — UP.",
	},
	"nz2": {
		"Rising Blip. Mr. Lattice is brokering a proxy war up the stack. Pick a patron, breach the Sundered Arbiter, and your alliance is sealed.",
	},
	"nz3": {
		"Infrastructure and the Blur. Ravel guards the deep Net's foundation. Shatter WARDEN-PRIME and Echo-9's fate is yours to decide.",
	},
	"nz4": {
		"Crosshairs of Power. Two worlds about to go dark. Rewrite the Rogue Overseer's core before the blackout takes them both.",
	},
	"nz5": {
		"Architects of Reality. The Master Protocol sits at the Catalyst Core. Beat the Prime Architect and you touch the code the Net is written in.",
	},
	"nz6": {
		"The Digital Pantheon. The makers of the Net wait here. Overwrite the Genesis Protocol Architects — if you can out-think gods.",
	},
	"nz7": {
		"The Infomorphic Ascension. A rift between living and dead universes is leaking. Weave the firewall through the Entropy-Titan and seal it.",
	},
	"nz8": {
		"The Ancient Archetypes. The Last Cartographer plays a long game on a board the size of the multiverse. Out-breach the Reconciled Ancient.",
	},
	"nz9": {
		"The Genesis Forge. Where new realities are minted. Siege it, and unmaking and creation come down to one last run.",
	},
	"nz10": {
		"The Living Library. Everything that was ever known, still humming. Few jack this deep. Fewer jack back out the same.",
	},
}

// loreKey maps the player's current room to a zoneLore key, or "" if there's no
// authored backstory (the city/rings).
func loreKey(roomID string) string {
	realm, zone := areaInfo(roomID)
	if zone == 0 {
		return ""
	}
	switch realm {
	case "meat":
		return "z" + itoa(zone)
	case "net":
		return "nz" + itoa(zone)
	}
	return ""
}

// talkSpeaker names who answers TALK here: the hiring fixer if one is present,
// otherwise an anonymous local appropriate to the realm.
func (w *World) talkSpeaker(p *Player) string {
	for _, q := range w.questsHere(p) {
		if q.GiverName != "" {
			return q.GiverName
		}
	}
	if w.inNet(p) {
		return "a jacked-in netrunner"
	}
	return "a wary local"
}

// boothIntro is the clone-booth tech's onboarding primer for new runners — a
// quick rundown of the core verbs, given when you TALK in the Re-Clone Bay.
var boothIntro = []string{
	"Fresh sleeve, huh? Name's Doc Splice. Quick orientation before you jack out, kid:",
	"  Move with N/S/E/W (or UP/DOWN). LOOK (L) to read a room, MAP (M) to see exits and the way deeper or out.",
	"  Fight with ATTACK (A); OPEN caches; LOOT (LO) the remains. Press an inventory number to USE it fast.",
	"  QUESTS for bounties — ACCEPT then CLAIM at a broker or the giver. TALK to locals for the lore.",
	"  Stash gear in your pod here (STASH/GRAB — no limit), HOME to recall back, SPEND points to grow. Now go make some scrip.",
}

// talk delivers a line of local backstory — or, in the Re-Clone Bay, the
// new-player onboarding primer.
func (w *World) talk(p *Player, arg string) {
	if p.RoomID == startRoom {
		for _, line := range boothIntro {
			p.send(style(green, line) + crlf)
		}
		return
	}
	lore := zoneLore[loreKey(p.RoomID)]
	speaker := w.talkSpeaker(p)
	if len(lore) == 0 {
		p.send(style(dim, speaker+" has nothing to say about this place.") + crlf)
		return
	}
	line := lore[w.roll(len(lore))]
	p.send(style(neon, speaker+": ") + style(green, line) + crlf)
}
