package cowboy

// ---------------------------------------------------------------------------
// The street-level RP transit rings: the Inner Circuit (a fast, RP-safe neon
// maglev loop) and the Sprawlbelt (a longer ground-level beltway with a few
// light L1-11 strays). Low-level hangout space for role-players, hung off Neon
// Alley; the two rings join at ic_4 <-> sb_1.
//
// Quest-givers around the rings hand out standalone "rumor" bounties — no linear
// progression: the pool (Pool == "ring") is shuffled and scattered across the
// givers fresh each session (assignRingQuests), and every rumor carries an
// easter-egg nod to one of the two PvE paths (the underground / the Net).
// ---------------------------------------------------------------------------

// ringGiverRooms are the NPCs around the rings who hand out rumor bounties.
var ringGiverRooms = []string{"ic_2", "ic_5", "sb_2", "sb_5", "sb_8", "sb_10"}

func ring(id, name, desc, flags string, exits map[string]string) *Room {
	r := &Room{ID: id, Name: name, Desc: wrapText(desc, 76), Exits: exits}
	for _, c := range flags {
		switch c {
		case 's':
			r.Safe = true
		case 'v':
			r.Vendor = true
		case 'm':
			r.Medic = true
		}
	}
	return r
}

func buildRingRooms() []*Room {
	return []*Room{
		// ---- Inner Circuit (RP-safe express loop) ----
		ring("ic_1", "Inner Circuit :: Neon Gate", "An elevated maglev platform where the Inner Circuit hums in a ring of cold blue light; the Strip glows far below and Neon Alley is a step south.", "s",
			map[string]string{"south": "neon_alley", "east": "ic_2", "west": "ic_6"}),
		ring("ic_2", "Inner Circuit :: Busker's Span", "A glass span over the avenues where a chrome-armed busker saws a synth-violin for scattered scrip and trades rumors with anyone who'll listen.", "s",
			map[string]string{"east": "ic_3", "west": "ic_1"}),
		ring("ic_3", "Inner Circuit :: Mirrorglass Curve", "The express ring curves through a tunnel of mirrored ad-panels, your own reflection sliding alongside you at speed.", "s",
			map[string]string{"east": "ic_4", "west": "ic_2"}),
		ring("ic_4", "Inner Circuit :: The Junction", "A spoked interchange where the fast Inner Circuit drops north onto the long Sprawlbelt; signage flickers directions in six languages at once.", "s",
			map[string]string{"east": "ic_5", "west": "ic_3", "north": "sb_1"}),
		ring("ic_5", "Inner Circuit :: The Rolling Rose", "A bar-car welded to the ring; the Rolling Rose pours cold synth-gin for commuters while a chrome-jawed bartender works the taps and the gossip.", "sv",
			map[string]string{"east": "ic_6", "west": "ic_4"}),
		ring("ic_6", "Inner Circuit :: Lantern Loop", "Paper-and-LED lanterns sway over the rail as the ring curves back toward the Neon Gate, warm light pooling on the wet platform.", "s",
			map[string]string{"east": "ic_1", "west": "ic_5"}),

		// ---- Sprawlbelt (outer beltway; a few light strays) ----
		ring("sb_1", "Sprawlbelt :: On-Ramp", "The ground-level beltway begins here, an off-ramp dropping from the Inner Circuit into the city's lived-in edge.", "s",
			map[string]string{"south": "ic_4", "east": "sb_2", "west": "sb_10"}),
		ring("sb_2", "Sprawlbelt :: Noodle Row Stop", "A steam-wreathed belt station packed with noodle carts; a one-armed cook ladles broth and hires odd hands between orders.", "sv",
			map[string]string{"east": "sb_3", "west": "sb_1"}),
		ring("sb_3", "Sprawlbelt :: Tagged Underpass", "A concrete underpass layered in UV gang-tags; a turf-tagger eyes you and reaches for a length of pipe.", "",
			map[string]string{"east": "sb_4", "west": "sb_2"}),
		ring("sb_4", "Sprawlbelt :: The Hustle Corner", "Folding tables and rigged dice under a busted streetlamp, where a fast-talking hustler works the belt crowd.", "",
			map[string]string{"east": "sb_5", "west": "sb_3"}),
		ring("sb_5", "Sprawlbelt :: The Fortune Stall", "A beaded stall of incense and chrome where a fortune-teller reads palms and cortex-chips alike for a coin.", "s",
			map[string]string{"east": "sb_6", "west": "sb_4"}),
		ring("sb_6", "Sprawlbelt :: Wirehead Hollow", "A dim belt alcove where a strung-out wirehead twitches through a looped braindance, jacks trailing from their temple.", "",
			map[string]string{"east": "sb_7", "west": "sb_5"}),
		ring("sb_7", "Sprawlbelt :: The Long Awning", "A covered belt promenade, dry beneath the drumming rain, barred windows modeling last season's chrome.", "s",
			map[string]string{"east": "sb_8", "west": "sb_6"}),
		ring("sb_8", "Sprawlbelt :: Scrap Exchange", "A cage-fronted stall heaped with salvaged parts where a scrap-trader buys drone-cores by the gram, no questions asked.", "sv",
			map[string]string{"east": "sb_9", "west": "sb_7"}),
		ring("sb_9", "Sprawlbelt :: The Stripped Lot", "A chain-link belt yard of gutted vehicles on cinderblocks where a scrap-scrounger drone picks the husks for copper.", "",
			map[string]string{"east": "sb_10", "west": "sb_8"}),
		ring("sb_10", "Sprawlbelt :: The Mending Stop", "A folding-cot belt clinic where a street medic patches the beltway's walking wounded and trades word of who's hiring.", "sm",
			map[string]string{"east": "sb_1", "west": "sb_9"}),
	}
}

// buildRingMobs places the rings' light L1-11 strays (homed in the non-safe belt
// rooms). They are weak and most are non-aggressive — the rings stay RP-friendly.
func buildRingMobs() []*MobTemplate {
	return []*MobTemplate{
		{ID: "ring_ganger", Name: "a turf-tagger", HP: 12, Damage: 3, AC: 2, XP: 18, Eddies: 8, Aggressive: true, Home: "sb_3"},
		{ID: "ring_hustler", Name: "a three-card hustler", HP: 8, Damage: 2, AC: 1, XP: 12, Eddies: 6, Aggressive: false, Home: "sb_4"},
		{ID: "ring_junkie", Name: "a strung-out wirehead", HP: 10, Damage: 2, AC: 1, XP: 14, Eddies: 5, Aggressive: false, Home: "sb_6"},
		{ID: "ring_scavver", Name: "a scrap-scrounger drone", HP: 14, Damage: 4, AC: 3, XP: 22, Eddies: 12, Aggressive: false, Home: "sb_9"},
	}
}

// assignRingQuests shuffles the ring "rumor" pool and round-robins it across the
// ring givers, so the rumors you find on the rings are scattered fresh each
// session (no fixed order, no progression). Uses the injectable roll.
func (w *World) assignRingQuests() {
	w.ringOffer = map[string][]int{}
	var pool []int
	for i := range quests {
		if quests[i].Pool == "ring" {
			pool = append(pool, i)
		}
	}
	for i := len(pool) - 1; i > 0; i-- { // Fisher-Yates via the world RNG
		j := w.roll(i + 1)
		pool[i], pool[j] = pool[j], pool[i]
	}
	for k, idx := range pool {
		room := ringGiverRooms[k%len(ringGiverRooms)]
		w.ringOffer[room] = append(w.ringOffer[room], idx)
	}
}
