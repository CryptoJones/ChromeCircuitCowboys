package cowboy

// startRoom is where new and respawning cowboys appear — a PRIVATE capsule pod,
// so a fresh jack-in or a respawn can never be spawn-camped. You step OUT into
// the street (Neon Alley) under your own power.
const startRoom = "capsule"

// buildRooms returns the Chrome Circuit Cowboys world map — a slice of the city
// and the Net beyond the jack-in port.
func buildRooms() map[string]*Room {
	rooms := []*Room{
		{ID: "capsule", Name: "Re-Clone Bay :: Your Booth", Private: true, Safe: true,
			Desc:  "A private booth in the clone clinic. You come to in a fresh clone, your\r\nmind restored from its realtime backup — calm, whole, and a few scrip lighter\r\nfor the new body. The clinic doors slide OUT into the street.",
			Exits: map[string]string{"out": "neon_alley"}},
		{ID: "neon_alley", Name: "Neon Alley", Safe: true,
			Desc:  "Rain hisses on hot neon. Holo-ads for synth-ramen and combat clinics\r\nflicker across puddles. The Strip roars to the east; a battered door to\r\nthe south leads into the Chrome Rose; the re-clone clinic is just IN off\r\nthe street. Security drones hum overhead — draw on another runner here and they\r\nflatline you on the spot. (no-violence zone)",
			Exits: map[string]string{"east": "the_sprawl", "south": "chrome_bar", "in": "capsule"}},
		{ID: "chrome_bar", Name: "The Chrome Rose", Vendor: true,
			Desc:  "A runner dive. Chrome-plated regulars jack into the bar's local node\r\nwhile an augmented bartender slings stims and gear. A vendor terminal glows\r\nhere (type LIST).",
			Exits: map[string]string{"north": "neon_alley"}},
		{ID: "the_sprawl", Name: "The Strip",
			Desc:  "Endless arcologies stacked into the smog. Crowds churn between street\r\nstalls. A black alley opens north; corporate spires gleam east; the Night\r\nMarket is south.",
			Exits: map[string]string{"west": "neon_alley", "north": "back_alley", "east": "corpo_plaza", "south": "market"}},
		{ID: "back_alley", Name: "Back Alley",
			Desc:  "A dead-end choked with dumpsters and busted drones. Gangers tag the walls\r\nin UV paint and don't like tourists.",
			Exits: map[string]string{"south": "the_sprawl"}},
		{ID: "market", Name: "Night Market", Vendor: true, Medic: true,
			Desc:  "Stalls of grey-market cyberware and noodle carts under string lights. A\r\nbroker runs a vendor stall (LIST) and a back-room Emergency Medic wires in chrome\r\n(INSTALL salvaged cyberware).",
			Exits: map[string]string{"north": "the_sprawl"}},
		{ID: "corpo_plaza", Name: "Corporate Plaza",
			Desc:  "Glass and gun-metal. Security drones sweep the concourse and corpo-sec in\r\nmirror visors watch everything. A guarded data port hums to the east.",
			Exits: map[string]string{"west": "the_sprawl", "east": "data_port"}},
		{ID: "data_port", Name: "Data Port", Term: true,
			Desc:  "A jack-in cradle wired to the city grid. Jacking in (UP) drops your\r\nconsciousness into the Net — the seedy underbelly of cyberspace.",
			Exits: map[string]string{"west": "corpo_plaza", "up": "nz1_1_top"}},
	}
	// Meatspace main quest (L1-99): the authored 10-arc underground descent
	// (zones.go), hanging off Back Alley. The Net (L1-99): the authored 10-arc
	// ascent of 3-layer nodes (netzones.go), jacked into from the Data Port.
	zoneRooms, _ := buildUndergroundZones()
	rooms = append(rooms, zoneRooms...)
	netRooms, _ := buildNetZones()
	rooms = append(rooms, netRooms...)
	rooms = append(rooms, buildRingRooms()...) // street-level RP transit rings

	m := make(map[string]*Room, len(rooms))
	for _, r := range rooms {
		m[r.ID] = r
	}
	// Back Alley drops DOWN into the Neon Wasteland; the Data Port jacks UP into
	// the first Net node's access shell (and UP from there jacks back out).
	m["back_alley"].Exits["down"] = "z1_01"
	m["z1_01"].Exits["up"] = "back_alley"
	m["nz1_1_top"].Exits["up"] = "data_port"
	// Neon Alley steps NORTH up onto the Inner Circuit (the RP transit rings).
	m["neon_alley"].Exits["north"] = "ic_1"
	return m
}

// mobTemplates defines the hostiles and where they live. The Home field is set
// from the map key for respawn placement.
func buildMobTemplates() map[string]*MobTemplate {
	defs := []*MobTemplate{
		{ID: "ganger", Name: "a street ganger", HP: 18, Damage: 4, AC: 2, XP: 25, Eddies: 10, Aggressive: true, Home: "back_alley"},
		{ID: "drone", Name: "a security drone", HP: 30, Damage: 7, AC: 5, XP: 50, Eddies: 25, Aggressive: false, Home: "corpo_plaza"},
		{ID: "corposec", Name: "a corpo-sec officer", HP: 45, Damage: 10, AC: 6, XP: 80, Eddies: 40, Aggressive: false, Home: "corpo_plaza"},
	}
	// Authored underground hostiles + loot caches (L1-99 meatspace zones).
	_, zoneMobs := buildUndergroundZones()
	defs = append(defs, zoneMobs...)
	// Authored Net hostiles + data-caches (L1-99 Net ascent).
	_, netMobs := buildNetZones()
	defs = append(defs, netMobs...)
	// Light strays on the RP transit rings.
	defs = append(defs, buildRingMobs()...)
	// The multi-stage "Gauntlet" ICE: a reconfiguring lattice in the first Net
	// node's core. Only the white shell spawns (Home set); on "death" each stage
	// morphs into the next, harder one, and only the final lethal lock pays out.
	defs = append(defs,
		&MobTemplate{ID: "gauntlet1", Name: "the Gauntlet ICE [white shell]", HP: 40, Damage: 10, AC: 5, Aggressive: true, ICE: true, Home: netGauntletHome, Next: "gauntlet2"},
		&MobTemplate{ID: "gauntlet2", Name: "the Gauntlet ICE [black core]", HP: 70, Damage: 16, AC: 8, Aggressive: true, ICE: true, Next: "gauntlet3"},
		&MobTemplate{ID: "gauntlet3", Name: "the Gauntlet ICE [lethal lock]", HP: 110, Damage: 24, AC: 11, XP: 700, Eddies: 600, Aggressive: true, ICE: true},
	)
	m := make(map[string]*MobTemplate, len(defs))
	for _, t := range defs {
		m[t.ID] = t
	}
	return m
}

// ware is a purchasable item.
type ware struct {
	name  string
	price int
	heal  int // stimpak: HP restored on use
	ram   int // ram-chip: RAM restored on use
	bonus int // weapon: attack bonus granted on purchase (permanent)
	deck  int // cyberdeck: MaxRAM bonus granted on purchase (permanent)
	body  int // implant: permanent Body boost on install
	refl  int // implant: permanent Reflexes boost on install
	intel int // implant: permanent Intelligence boost on install
	desc  string
}

// isImplant reports whether a ware is a stat implant (installed at a medic for a
// permanent Body/Reflexes/Intelligence boost — the per-band class grind aids).
func (x ware) isImplant() bool { return x.body > 0 || x.refl > 0 || x.intel > 0 }

// shopWares is the MASTER catalog — every purchasable/loadable item across all
// tiers. findWare searches it, so any looted item can be used/installed anywhere.
// What a given vendor actually STOCKS is a per-area subset (see waresForRoom).
var shopWares = []ware{
	// Tier 1 — the street (Chrome Rose / Night Market, L1–10).
	{name: "stimpak", price: 20, heal: 25, desc: "single-use trauma stim, restores 25 HP"},
	{name: "ram-chip", price: 30, ram: 8, desc: "single-use RAM chip, restores 8 RAM for netruns"},
	{name: "ice-breaker", price: 150, bonus: 5, desc: "intrusion blade, +5 attack (permanent)"},
	{name: "mono-katana", price: 400, bonus: 12, desc: "monomolecular katana, +12 attack (permanent)"},
	{name: "cyberdeck", price: 250, deck: 8, desc: "upgraded deck, +8 max RAM (permanent)"},
	// Deeper tiers — sold at the band safehouses (better, pricier gear).
	{name: "trauma-kit", price: 120, heal: 60, desc: "field trauma kit, restores 60 HP"},
	{name: "mega-stim", price: 400, heal: 120, desc: "military stim, restores 120 HP"},
	{name: "ram-bank", price: 150, ram: 20, desc: "RAM bank, restores 20 RAM for netruns"},
	{name: "war-axe", price: 1200, bonus: 20, desc: "powered war-axe, +20 attack (permanent)"},
	{name: "rail-blade", price: 3000, bonus: 30, desc: "rail-driven blade, +30 attack (permanent)"},
	{name: "monowire", price: 8000, bonus: 45, desc: "monomolecular wire, +45 attack (permanent)"},
	{name: "quantum-deck", price: 1500, deck: 16, desc: "quantum deck, +16 max RAM (permanent)"},
	{name: "neural-deck", price: 6000, deck: 28, desc: "neural-lace deck, +28 max RAM (permanent)"},

	// Stat implants — per-band class grind aids (install at a Emergency Medic for
	// a permanent stat boost). Three per tier: Body / Reflexes / Intelligence.
	{name: "subdermal-plating", price: 110, body: 2, desc: "subdermal armor weave, +2 Body (permanent)"},
	{name: "reflex-booster", price: 110, refl: 2, desc: "reflex booster chip, +2 Reflexes (permanent)"},
	{name: "neural-coprocessor", price: 110, intel: 2, desc: "neural coprocessor, +2 Intelligence (permanent)"},

	{name: "titanium-weave", price: 320, body: 3, desc: "titanium muscle weave, +3 Body (permanent)"},
	{name: "kerenzikov", price: 320, refl: 3, desc: "kerenzikov nerve wiring, +3 Reflexes (permanent)"},
	{name: "cortex-bridge", price: 320, intel: 3, desc: "cortex bridge, +3 Intelligence (permanent)"},

	{name: "myomer-bundle", price: 850, body: 4, desc: "myomer muscle bundle, +4 Body (permanent)"},
	{name: "synaptic-amp", price: 850, refl: 4, desc: "synaptic amplifier, +4 Reflexes (permanent)"},
	{name: "mnemonic-array", price: 850, intel: 4, desc: "mnemonic array, +4 Intelligence (permanent)"},

	{name: "juggernaut-frame", price: 2000, body: 5, desc: "juggernaut subframe, +5 Body (permanent)"},
	{name: "sandevistan", price: 2000, refl: 5, desc: "sandevistan rig, +5 Reflexes (permanent)"},
	{name: "quantum-cortex", price: 2000, intel: 5, desc: "quantum cortex, +5 Intelligence (permanent)"},

	{name: "goliath-chassis", price: 4500, body: 6, desc: "goliath chassis, +6 Body (permanent)"},
	{name: "hyper-reflex", price: 4500, refl: 6, desc: "hyper-reflex lattice, +6 Reflexes (permanent)"},
	{name: "ascendant-mind", price: 4500, intel: 6, desc: "ascendant mind-core, +6 Intelligence (permanent)"},
}

func findWare(name string) (ware, bool) {
	for _, w := range shopWares {
		if w.name == name {
			return w, true
		}
	}
	return ware{}, false
}

// pickWares pulls a curated subset of the master catalog by name.
func pickWares(names ...string) []ware {
	out := make([]ware, 0, len(names))
	for _, n := range names {
		if w, ok := findWare(n); ok {
			out = append(out, w)
		}
	}
	return out
}

// Per-area vendor stock. The fixed city vendors carry tier-1 street gear; the
// authored zone vendors (underground safehouses + Net access shells) scale their
// stock to their level band (see zoneVendorBand + waresForBand).
var (
	streetWares      = pickWares("stimpak", "ram-chip", "ice-breaker")              // Chrome Rose
	nightMarketWares = pickWares("stimpak", "ram-chip", "mono-katana", "cyberdeck") // Night Market (+ Emergency Medic)
)

// waresForRoom returns the gear a vendor in roomID stocks. The two city vendors
// carry curated street stock; authored zone vendors scale with their band; any
// other vendor falls back to the full catalog.
func waresForRoom(roomID string) []ware {
	switch roomID {
	case "chrome_bar":
		return streetWares
	case "market":
		return nightMarketWares
	}
	if b, ok := zoneVendorBand[roomID]; ok { // authored zone vendors (underground + Net)
		return waresForBand(b)
	}
	return shopWares
}
