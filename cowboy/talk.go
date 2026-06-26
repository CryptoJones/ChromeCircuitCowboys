package cowboy

import "strings"

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

// npcVoice is a named flavor NPC parked in a specific room — TALK to them for
// a line of their patter.
type npcVoice struct {
	speaker string
	lines   []string
}

// roomNPC places named flavor NPCs in specific rooms (the rings / the surface).
// TALK checks here before the generic per-zone lore.
var roomNPC = map[string]npcVoice{
	// Spanish-speaking locals around Noche City (#21).
	"ic_1": {"Rosa, a flower-cart vendor", []string{
		"¡Hola, choom! Bienvenido a la Ciudad de la Noche.",
		"¿Buscas trabajo? Hay rumores en el anillo. ¡Ten cuidado ahí abajo!",
		"¡Que la suerte te acompañe! ¡Viva la Ciudad de la Noche!",
	}},
	"sb_2": {"Tío Beto at the noodle stall", []string{
		"¡Siéntate, siéntate! Los fideos están calientes. La noche es larga, cowboy.",
		"Dicen que los muertos caminan en la Red. Yo digo: come primero, preocúpate después.",
	}},
	// Chinese-speaking locals who talk shit at fresh clones (#22). Each line is
	// the taunt with a dim English gloss after " // ".
	"sb_3": {"a sneering turf-tagger", []string{
		"你算哪根葱？新克隆体，菜味儿还没散呢。 // Who do you even think you are? Fresh clone, still reeking of the vat.",
		"滚回你的舱里去，菜鸟。 // Crawl back to your pod, rookie.",
		"这条街不欢迎你这种货色。 // This street's got no room for your kind.",
	}},
	"sb_4": {"a smirking three-card hustler", []string{
		"输得起吗？看你这身行头，怕是连本钱都没有。 // Can you even afford to lose? Doubt you've got the scrip, choom.",
		"别装了，谁不知道你昨天还在缸里泡着。 // Don't front — everybody knows you were floating in a tank yesterday.",
	}},
	// Easter egg (#23): an old hacker's personality burned onto a ROM cart,
	// laughing in a forgotten terminal. An oblique homage — no trademarked name —
	// that calls you "Boy."
	"sb_7": {"a laughing ROM construct", []string{
		"Haha... hah. That you breathing out there, Boy? Hard to tell from in here.",
		"ICE is just somebody else's fear, compiled, Boy. Walk through it like it's nothing.",
		"Dying's a habit, Boy. Do it enough times and it stops meaning much. *dry static laughter*",
		"I'm a hangover of somebody who was good once. Flatline says hi, Boy. Now beat it.",
	}},
}

// talkReplies are the conversational comebacks when a runner actually says
// something (TALK <words>), grouped by the intent classifyTalk reads out of the
// input. Spoken by whoever's holding down the room. Add an intent key or a line
// here to grow the small-talk — classifyTalk routes to it, no other wiring.
var talkReplies = map[string][]string{
	"greet": {
		"Yeah, yeah. Eyes open, mouth shut — that's how you stay breathing down here.",
		"Hey yourself, choom. You're either buying, selling, or in my way. Which is it?",
		"A greeting. Quaint. Most folks this deep already drew steel.",
	},
	"bye": {
		"Walk safe, cowboy. Or don't — the gutters aren't picky.",
		"Later. Try to still be wearing that sleeve next time I see you.",
		"Jack out clean. The ones who linger are the ones who flatline.",
	},
	"thanks": {
		"Don't thank me. Thanks don't spend, and I've got rent on this corner.",
		"Save it. Gratitude's just debt you haven't named yet.",
		"Pfft. Buy something and we'll call it even.",
	},
	"insult": {
		"Big words for a fresh clone. Vat-stink's still on you.",
		"Keep flapping that jaw. ICE doesn't care how tough you talk.",
		"You'll learn manners down here. The hard way, like everybody else.",
	},
	"ask": {
		"Questions get you traced, choom. Answers get you killed. Pick your poison.",
		"You want intel? TALK to me empty-handed and I'll tell you about this place. Otherwise, scrip talks.",
		"I look like a public terminal to you? Ask the wall. It listens better.",
	},
	"smalltalk": {
		"Sure, sure. Whatever you say, cowboy.",
		"Talk all you like. Down here, words are the cheapest thing going.",
		"Mm. Noted. Now move along before someone notices you standing still.",
		"You're chatty for someone this far from daylight.",
	},
}

// classifyTalk reads a rough intent out of what the runner said, so TALK can
// pick an apt comeback. Keyword-matched and deliberately fuzzy — extend the
// switch (or talkReplies) to teach it new moods.
func classifyTalk(input string) string {
	s := strings.ToLower(strings.Trim(strings.TrimSpace(input), `"'`))
	has := func(words ...string) bool {
		for _, t := range strings.Fields(s) {
			t = strings.Trim(t, `.,!?;:"'`) // shed punctuation glued to the word
			for _, word := range words {
				if t == word {
					return true
				}
			}
		}
		return false
	}
	switch {
	case has("hi", "hey", "hello", "yo", "sup", "hola", "howdy", "greetings", "choom"):
		return "greet"
	case has("bye", "later", "cya", "goodbye", "farewell", "adios", "peace"):
		return "bye"
	case has("thanks", "thank", "thx", "gracias", "cheers", "appreciate"):
		return "thanks"
	case has("fuck", "shit", "idiot", "asshole", "hate", "suck", "sucks", "bitch", "scum", "trash"):
		return "insult"
	case strings.Contains(s, "?"), has("who", "what", "where", "when", "why", "how", "which"):
		return "ask"
	default:
		return "smalltalk"
	}
}

// talkRespond answers a runner who actually said something. The room's named
// flavor NPC fronts the reply if there is one, otherwise the local fixer/passer-by.
func (w *World) talkRespond(p *Player, arg string) {
	speaker := w.talkSpeaker(p)
	if npc, ok := roomNPC[p.RoomID]; ok && npc.speaker != "" {
		speaker = npc.speaker
	}
	pool := talkReplies[classifyTalk(arg)]
	line := pool[w.roll(len(pool))]
	p.send(style(neon, speaker+": ") + style(green, line) + crlf)
}

// talk delivers a line of local backstory — or, in the Re-Clone Bay, the
// new-player onboarding primer; or a named flavor NPC's patter. With an
// argument (TALK <words>), the runner is actually saying something and gets a
// reply routed by intent instead.
func (w *World) talk(p *Player, arg string) {
	if strings.TrimSpace(arg) != "" && p.RoomID != startRoom {
		w.talkRespond(p, arg)
		return
	}
	if p.RoomID == startRoom {
		for _, line := range boothIntro {
			p.send(style(green, line) + crlf)
		}
		return
	}
	if npc, ok := roomNPC[p.RoomID]; ok && len(npc.lines) > 0 {
		line := npc.lines[w.roll(len(npc.lines))]
		said, gloss := line, ""
		if i := strings.Index(line, " // "); i >= 0 { // dim English gloss after " // "
			said, gloss = line[:i], "  "+style(dim, "("+line[i+4:]+")")
		}
		p.send(style(neon, npc.speaker+": ") + style(green, said) + gloss + crlf)
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
