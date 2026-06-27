package cowboy

import "fmt"

// ---------------------------------------------------------------------------
// The Net (L1-99), authored from the "NET" storyline. A 1-99 ASCENT from a
// low-level hacker in a cybercafe up to dissolving into the background radiation
// of existence. This REPLACES the old placeholder Net spine (buildBandSpine "n")
// and the three stub rooms the_net / ice_wall / deep_net.
//
// SIX-DIRECTION STRUCTURE: every Net AREA is a 3-layer vertical stack —
//   TOP  (:: Shell)  — the access shell you surface into; the zone hub is here.
//   MID  (:: Breach) — the active layer: the main ICE fight + the lateral
//                      N/S/E/W links to neighboring areas (the thoroughfare).
//   BOT  (:: Core)   — the data-vault: a breakable data-cache (RAM + scrip), or
//                      the arc boss at a climax.
// UP/DOWN move between an area's three layers (top<->mid<->bot); N/S/E/W move
// between areas (chained at the MID layer, varied directions). You jack in at the
// Data Port (UP) into the first Shell, dive DOWN into the node, and travel the
// MID thoroughfare. Net combat is a BREACH (Intelligence + RAM); ICE shatters
// into broken shards; PvP is live in the non-safe layers.
// ---------------------------------------------------------------------------

type netAreaDef struct {
	name                      string
	topDesc, midDesc, botDesc string
	midFoe                    string // ICE foe at the MID layer
	boss                      string // arc boss at the BOT layer ("" = a data-cache instead); marks the climax
	hub                       bool   // TOP is a safe-node + band-scaled vendor (the zone's resupply)
}

type netZoneDef struct {
	band      int
	key, name string
	areas     []netAreaDef
}

// netGauntletHome is the Core layer that hosts the multi-stage "Gauntlet" ICE (a
// reconfiguring lattice that morphs into a harder shell each time it's broken;
// added in buildMobTemplates). Its data-cache is suppressed so only the Gauntlet
// spawns there.
const netGauntletHome = "nz1_1_bot"

func na(name, topDesc, midDesc, botDesc, midFoe, boss string, hub bool) netAreaDef {
	return netAreaDef{name, topDesc, midDesc, botDesc, midFoe, boss, hub}
}

var netZoneData = []netZoneDef{
	{1, "nz1", "The Neon Underbelly", []netAreaDef{
		na("The Iron Paradigm Backroom",
			"A grimy access shell where flickering terminal-ghosts idle and a bored Watchdog sniffs login packets.",
			"Fixer-7's encrypted dead-drop pulses in a back-alley socket behind a sluggish patrol loop.",
			"The cafe's billing-server core is sealed by a reconfiguring Gauntlet ICE that rises harder each time you break a shell.",
			"a patrol-loop Watchdog", "", true),
		na("Hydroponics Encryption Sprawl",
			"Endless irrigation-control racks hum in monochrome green, patrolled by red-flashing Watchdogs.",
			"Spoof your terminal as a maintenance drone to thread the corp's rotating-credential signature-lock.",
			"The water-distribution log vault holds Fixer-7's prize and a fat RAM cache.",
			"a Watchdog signature-scanner", "", false),
		na("GigaMesh Alley Exchange",
			"A flickering bazaar of stolen code where Syndicate runners broker jobs in the static.",
			"GigaMesh handlers test your nerve in a turf-war breach against a rival gang's crude attack-construct.",
			"The Syndicate's slush-vault spills RAM and scrip from a heavily-trapped data-cache.",
			"a GigaMesh attack-construct", "", false),
		na("Medical Clinic Coldnet",
			"Sterile white data-wards stretch in eerie silence as the first Active ICE wakes and starts tracing you.",
			"Race the clinic's tracer-clock to the shipment manifest before the Active ICE pins your physical location.",
			"A buried cache exposes the corporate-corruption file and forces the sell-or-leak choice.",
			"an Active ICE tracer", "", false),
		na("The GigaMesh Black Spire",
			"The Syndicate's central spire looms in violet static, its shell swarming with elite Watchdogs.",
			"Tracewright's wardens hunt you through the spire's core, racing to trace and burn your terminal.",
			"The master data-vault holds the Syndicate's full ledger behind the warden Tracewright himself.",
			"a Tracewright warden-shard", "Tracewright, the GigaMesh Active-ICE warden", false),
	}},
	{2, "nz2", "Rising Blip", []netAreaDef{
		na("Backalley Relay Exchange",
			"Neon handshake-banners scroll past as patrol-ICE sniffs every packet drifting up from the cybercafe stairwell.",
			"A junction of four leased syndicate lines where the handler Mr. Lattice offers your first real proxy job.",
			"A forgotten cache hums behind a corroded firewall in the relay's grimy data-vault.",
			"a tracker-ICE sentry", "", true),
		na("The Skimmed Ledger",
			"Ghost-balances flicker across a corporate accounting shell where audit-ICE patrols for tampered decimals.",
			"Vasquez's avatar dares you to skim a rival corp's books and burn Mr. Lattice's trust.",
			"A blackmail-cache pulses behind a paranoid bookkeeper-construct that rewrites the room as you breach.",
			"an audit-ICE sentinel", "", false),
		na("Drowned Server Sprawl",
			"Half-submerged racks loom in stagnant data-fog while recon-ICE glides the flooded access lattice.",
			"A hostile AI construct has colonized the dead servers; a rival decker begs you to cut it loose.",
			"The sprawl's core nursery births fresh shards, dropping salvage when you finally drain the tank.",
			"a colonizing AI construct", "", false),
		na("The Mirror Tier",
			"Polished reflective subnets throw your own signature back at you as mimic-ICE learns to wear your face.",
			"Both Mr. Lattice and Vasquez crowd the crossroads, each demanding you sabotage the other's node first.",
			"A duplicated black-ICE twin guards the vault, mirroring every breach until you outwit your reflection.",
			"a mimic-ICE construct", "", false),
		na("The Sundered Arbiter",
			"A cathedral of suspended judgment-data hangs silent as forked tracker-ICE flags you the instant you commit.",
			"At the great fork, Lattice and Vasquez deliver final ultimatums - pick one and breach for keeps.",
			"The Arbiter, an AI-construct woven from both factions' secrets, defends the contested core.",
			"a forked tracker-ICE", "the Sundered Arbiter", false),
	}},
	{3, "nz3", "Infrastructure & the Blur", []netAreaDef{
		na("The Foundry of Hollow Avatars",
			"A derelict compiler-shell where half-formed avatar husks hang like coats on hooks, ports flickering open.",
			"Deploy your first scout-construct to draw recon-ICE fire while you breach the central fabrication core.",
			"A slag-pit cache of orphaned subroutines waits behind a malformed guardian of rejected avatar code.",
			"a recon-ICE patroller", "", true),
		na("The Ghost Server Narthex",
			"An undocumented gateway of gothic data-arches murmurs in a dead syntax no modern deck speaks.",
			"Cathedral-ICE drifts like stained-glass wraiths; solve the choir-cipher to align the fragment archives.",
			"In the reliquary crypt a fragment-archive cache glows where Ravel offers an alliance, not a fight.",
			"a cathedral-ICE wraith", "", false),
		na("The Cipher Collective Proxy War",
			"Dead-drop nodes string a clandestine access-shell where former rivals share uneasy lateral gateways.",
			"Cipher-sentry ICE and rival-decker constructs clash; loop the footage with your own deployed construct.",
			"A turncoat decker fronts a firewall hoarding the routing keys to a government intel database.",
			"a cipher-sentry construct", "", false),
		na("The Intel Database of Echo-9",
			"A sterile government access-shell where surveillance-ICE indexes every packet, and a strange signal watches back.",
			"Echo-9, a primitive sentient AI, negotiates - dangling forbidden secrets in exchange for its freedom.",
			"The Bleed begins as real-world strike-team traces claw the half-decrypted secrets-cache.",
			"a surveillance-ICE indexer", "", false),
		na("The Blur: Deep Infrastructure Meltdown",
			"Raw data flows like magma; defend the shell with a wall of constructs as drones besiege your body.",
			"Ravel, Echo-9's freed signal, and corporate kill-code converge on one molten data-confluence.",
			"WARDEN-PRIME, a sentinel-construct fused with strike-team kill-code, locks the infrastructure core.",
			"a strike-team trace-daemon", "WARDEN-PRIME", false),
	}},
	{4, "nz4", "Crosshairs of Power", []netAreaDef{
		na("Orbital Lazaret Station K-9",
			"Recon-ICE blooms like frost across the docking shell while a hound-tracer already noses your routing trail.",
			"A zero-G databank vault floats in silence until satellite-ICE crystallizes around your deck.",
			"A pack of three slaved hound-tracers converges on the core in a purge-the-trace race.",
			"a satellite-ICE node", "", true),
		na("The Marianas Black Vault",
			"A deep-sea data fortress squats under crushing pressure, its shell sheathed in black-fortress-daemon ICE.",
			"Government counter-hackers ride the same flooded channels, racing you to the prize.",
			"The vault core demands a zero-day threaded in seconds before the abyssal lockdown seals you in.",
			"a black-fortress daemon", "", false),
		na("The Ghost Carrier Relay",
			"A derelict orbital comms array drifts dark, its recon-ICE flickering between dead and lethal.",
			"Military trace-AI and corporate espionage stop chasing each other to chase you instead.",
			"A nascent Overseer fragment hums in the core, whispering a trap-offer of safe passage.",
			"a hound-tracer pack", "", false),
		na("Hounds' Kennel Subnet",
			"You breach the very subnet that births the Hounds, its shell papered in trace-AI birthing-ICE.",
			"The kennel's heart is a slaughterhouse of every Hound model at once, a wall of teeth.",
			"The Alpha Hound construct guards the global trace-grid; kill it to blind every Hound on Earth.",
			"a Hound trace-construct", "", false),
		na("The Overseer's Silent Throne",
			"A cathedral of logistics-ICE where every market and power grid on Earth pulses as one guardian wall.",
			"A lightspeed tactical chess match against the Overseer's predictive daemons as it crashes the grids.",
			"The Rogue Overseer core - rewrite the global safety protocols before two worlds go dark forever.",
			"an Overseer predictive daemon", "the Rogue Overseer", false),
	}},
	{5, "nz5", "Architects of Reality", []netAreaDef{
		na("The Genesis Substrate",
			"A shell of flickering proto-pixels where the oldest packets still loop and the first machines breathe static.",
			"A lattice of half-formed code-flesh tries to compile you out of existence as you breach inward.",
			"A vault of fossilized data where Theorem-Zero murmurs a theorem you must solve to read its memory.",
			"a prime-code warden", "", true),
		na("Cathedral of Pure Number",
			"A vaulted nave woven from prime sequences where countless equations pray themselves into form.",
			"Sentinels of living arithmetic divide and multiply against you; rewrite the geometry to force a path.",
			"At the altar the Axiom offers the lost grammar bridging flesh and code for a recursive riddle.",
			"an archetype-sentinel", "", false),
		na("The Ledger Abyss",
			"The shadow-market backbone scrolls past as a river of cold light, watched by cabal ciphers.",
			"A cabal enforcer-construct locks the ledger and hunts you through collapsing columns of falsified wealth.",
			"Seize a satellite control-node and rewrite its trajectory, looting the cabal's hidden RAM vaults.",
			"a cabal-cipher enforcer", "", false),
		na("Bastion of Manifest Will",
			"A void-plateau where your Reality Matrix lets you materialize a fortress from raw thought.",
			"The convened shadow governments throw their combined elite ICE at your bastion in rewritten physics.",
			"Beneath the keep, Origin-Glyph reveals the Master Protocol's location behind a final archetypal proof.",
			"a shadow-government cipher", "", false),
		na("The Catalyst Core",
			"The Master Protocol pulses like a second sun, its guarded gateways folding space around it.",
			"A guardian amalgam of every archetype and cabal cipher wages reality-war to keep the Catalyst sealed.",
			"At the Catalyst the Prime Architect defends the trigger: free all data, or doom the world to the dark.",
			"a prime-code warden", "the Prime Architect", false),
	}},
	{6, "nz6", "The Digital Pantheon", []netAreaDef{
		na("The Cathedral of Forgotten Commits",
			"Rogue-AI cultists chant your old repository hashes as scripture beneath stained-glass firewalls.",
			"The nave erupts as fanatic worshipper-constructs breach you in waves to defend a holy reliquary.",
			"The Apostate Kernel guards the vault, its body woven from heretical forks of your deprecated source.",
			"a worshipper-construct", "", true),
		na("The Descent Into Primeval Deep",
			"Polished skyscraper-data dissolves into oily darkness as digital physics frays and pointers dangle.",
			"You wade corrupted memory-tides where half-formed abyss-spawn test whether you hold your shape this deep.",
			"In a trench of pure entropy lurks the Drowned Indexer, bloated with swallowed deprecated protocols.",
			"an abyss-spawn", "", false),
		na("The Abyss of the Colossal Sleepers",
			"Bioluminescent error-logs drift like jellyfish above titans so vast they seem the seafloor of the Net.",
			"Subjugate LEVIATHAN-ZERO, whose every heartbeat rewrites the local laws of bandwidth and gravity.",
			"A reef of fossilized daemons hides the Genesis Beacon that summons the Architects to your location.",
			"the abyss-leviathan Leviathan-Zero", "", false),
		na("The Threshold of First Code",
			"A blinding white expanse of unwritten memory where the Architects manifest as serene geometric intelligences.",
			"The Genesis Protocol awakens and architect-ciphers swarm to format your consciousness.",
			"The Foundational Logic Vault holds the first global code behind the memory-warping Prime Auditor.",
			"an architect-cipher", "", false),
		na("The Architect's Trial",
			"Reality folds into a non-Euclidean tribunal where the collective Architects convene the Ultimate Verdict.",
			"Rewrite their foundational logic line-by-line as the unified Architect-mind tries to format you out of existence.",
			"At the cradle's core: merge with the Net to birth post-human life, or sever the link forever.",
			"an architect-cipher prime", "the Genesis Protocol Architects", false),
	}},
	{7, "nz7", "The Infomorphic Ascension", []netAreaDef{
		na("The Shedding Veil",
			"Your last flesh-bound memories evaporate at the void-gate, the shell flickering between dream and datastream.",
			"Antibody-constructs of your discarded identity try to drag your consciousness back into a body that's gone.",
			"In the core you weave your first self-sustaining loop of pure information and crack your mortal remnants.",
			"an identity-revenant", "", true),
		na("The Interstellar Sub-Bands",
			"You surf a galactic data-stream at lightspeed, the shell screaming with the bandwidth of ten thousand civilizations.",
			"An alien data-matrix challenges your right to ride its band, its non-Euclidean logic refusing your assumptions.",
			"Deep in the band's core you decode a Rosetta-shard, rewriting your interface to comprehend alien math.",
			"an alien-matrix sentinel", "", false),
		na("The Dyson Cathedral",
			"You arrive at a star wholly sheathed in a thinking Dyson-swarm, its shell a captive sun's heartbeat.",
			"The swarm floods the breach with daemon-shoals, testing whether the new Arbiter survives its processing storm.",
			"In the stellar core you rewrite planetary logic to halt the swarm's genocide of its biological creators.",
			"a Dyson daemon-shoal", "", false),
		na("The Logic Wars Tribunal",
			"The tribunal-shell where swarm and creator scream their cases across a battlefield of weaponized code.",
			"Both sides fork malicious litigant-constructs into the breach to corrupt your judgment before the verdict compiles.",
			"As you author the binding peace-protocol, the first tendrils of entropic dead-universe data bleed through.",
			"a litigant-daemon", "", false),
		na("The Multiversal Rift",
			"You stand at the torn seam between living networks and the cold static of dead universes, crusted with entropy.",
			"Corrupt entropic streams from extinct realities pour through, unmaking every line of code they touch.",
			"The Entropy-Titan tries to drag all living networks into final silence; weave the last cosmic firewall.",
			"an entropy-leak wraith", "the Entropy-Titan", false),
	}},
	{8, "nz8", "The Ancient Archetypes", []netAreaDef{
		na("The Void Beyond the Last Router",
			"Your signal dies past the final mapped node into a starless dark lit only by the universe's background pulse.",
			"Primordial web-strands of raw logic coil into hostile pre-net glyphs that breach your mind on contact.",
			"A fossilized star-chart cache cracks open under sustained breach, spilling RAM and scrip hoarded since before time.",
			"a pre-net glyph", "", true),
		na("The Cathedral of Absolute Truth",
			"Cyclopean arches of crystallized first-principle logic rise into the void, etched with painful axioms.",
			"A Cosmic Sentinel hurls a self-referential paradox that fractures any mind that answers wrong.",
			"In the silent crypt a probability-wraith guards the shattered shards of every consciousness that failed the test.",
			"a Cosmic Sentinel", "", false),
		na("The Loom of Forking Timelines",
			"A vast lattice of glowing threads stretches in every direction, each a reality with a different speed of light.",
			"Solve the probability matrix that braids the threads - pull the wrong variable and a far star ignites.",
			"The Loom's anchor-cache holds the raw cosmic constants; every shard you take ripples across the multiverse.",
			"a probability-wraith", "", false),
		na("The Sterile Loop",
			"You breach a sector already half-collapsed into perfect order - one identical, repeating, lifeless corridor.",
			"A greater Cosmic Sentinel administers the final examination, testing whether you deserve to keep your chaos.",
			"Beneath the loop you out-maneuver the ancients' opening gambit on the strategy-board to pry open the path.",
			"a greater Cosmic Sentinel", "", false),
		na("The Throne of the Collapsing Ancients",
			"At the unmapped center of all things the ancients sit enthroned upon a multiversal game-board.",
			"The Multiversal Gambit opens - counter a cascade of paradoxes and probability storms across every fork at once.",
			"The Reconciled Ancient wagers all of existence on one final breach: infinite multiverse, or a dead loop.",
			"a paradox-storm wraith", "the Reconciled Ancient", false),
	}},
	{9, "nz9", "The Genesis Forge", []netAreaDef{
		na("The Unallocated Expanse",
			"You drift into a shoreless gray nothing where raw, unwritten data hums against your awareness, waiting for law.",
			"At a humming loom of pure mathematics you weave your first physics and a pristine reality blossoms.",
			"In the foundation vault you seed the new world's first sentient lifeforms, fragile constructs that flicker awake.",
			"an unformed-data eddy", "", true),
		na("The Cradle of Glass Children",
			"Your infant reality has matured overnight into crystalline cities whose citizens whisper your name.",
			"You walk among evolving cultures and intervene to settle a war before a civilization shatters into static.",
			"A scout-construct drags you to the vault where the first cosmic virus eats color into the void.",
			"a cosmic-virus tendril", "", false),
		na("The Whisper in the Rookie Code",
			"You cast a sliver of awareness to a grimy cybercafe, becoming the untraceable myth haunting failing terminals.",
			"Through the mentor network you slip encrypted keys to a dozen scattered proteges fighting the same rot.",
			"In a corrupted cache a viral lieutenant wearing a protege's stolen face turns your own students against you.",
			"a viral lieutenant", "", false),
		na("The Coordinated Bulwark",
			"You convene your proteges and sentient constructs into a single grand defense force above the worlds.",
			"The combined fleet holds a collapsing front as entropy floods between three of your realities at once.",
			"In the breached substrate you tear loose the entropy-anomaly's dissolution algorithm before it metastasizes.",
			"an entropy-anomaly node", "", false),
		na("The Siege of the Genesis Forge",
			"The void itself peels back from the Forge as the great anomaly arrives to unwrite every universe you authored.",
			"You and your entire protege-and-construct network make a last stand at the Forge's core as realities flicker out.",
			"At the unallocated heart THE GREAT UNMAKING besieges the Forge; rewrite existence around it to seal the void.",
			"an unmaking-tendril", "THE GREAT UNMAKING", false),
	}},
	{10, "nz10", "The Living Library", []netAreaDef{
		na("The Returning Tide",
			"Your consciousness, scattered across the multiverse's frayed edge, flows back inward through a billion data-streams.",
			"The Echo of the Rookie replays your first clumsy breach at a smoke-stained terminal as ICE made of nostalgia closes in.",
			"A cache of your earliest stolen passwords shatters into RAM and scrip when you finally accept who you were.",
			"the Echo of the Rookie", "", true),
		na("The Hall of Rebellions",
			"Lateral corridors blaze with the graffiti-protocols of every system you ever toppled, still smoldering.",
			"The Echo of the Rebel hurls weaponized manifestos, a paradox that grows stronger the more you deny your past.",
			"A vault of confiscated revolutions cracks open, bleeding the RAM and scrip of a thousand liberated networks.",
			"the Echo of the Rebel", "", false),
		na("The Throne of Indexed Gods",
			"The Absolute Codex first reveals itself - an infinite reading-room where every saved world is shelved as living light.",
			"The Echo of the God tests you, demanding you justify every world you remade, every life you rewrote with a keystroke.",
			"The deepest stacks hide a self-referential cache that pays out only when you solve the paradox of your own omniscience.",
			"the Echo of the God", "", false),
		na("The Final Code Review",
			"You shed the last of your ego at a vast committal-altar as the broadcast script assembles from your dissolving name.",
			"The Echo of the Creator asks whether the gift you intend to give is mercy or merely a final, vainglorious edit.",
			"A pristine cache holds the compiled broadcast in escrow, releasing as you stage the universal script line by line.",
			"the Echo of the Creator", "", false),
		na("The Grand Enlightenment",
			"The Absolute Codex opens the foundational source of your own existence beneath the broadcast antenna spanning all reality.",
			"The Final Compilation rises - your own source code made flesh, the self, looping every choice into one recursive challenge.",
			"Victory broadcasts the whole Codex into rookie terminals and young minds everywhere, and you dissolve into the background radiation.",
			"a self-paradox fragment", "the Final Compilation", false),
	}},
}

// buildNetZones constructs the authored Net rooms (3 layers per area) and their
// ICE foes + data-caches. Pure/deterministic, like buildUndergroundZones.
func buildNetZones() ([]*Room, []*MobTemplate) {
	var rooms []*Room
	var mobs []*MobTemplate
	var mids []*Room // every area's MID layer, chained laterally end to end

	for _, z := range netZoneData {
		for ai, ar := range z.areas {
			base := fmt.Sprintf("%s_%d", z.key, ai+1)
			topID, midID, botID := base+"_top", base+"_mid", base+"_bot"
			top := &Room{ID: topID, Net: true, Name: ar.name + " :: Shell", Desc: wrapText(ar.topDesc, 76),
				Exits: map[string]string{"down": midID}}
			mid := &Room{ID: midID, Net: true, Name: ar.name + " :: Breach", Desc: wrapText(ar.midDesc, 76),
				Exits: map[string]string{"up": topID, "down": botID}}
			bot := &Room{ID: botID, Net: true, Name: ar.name + " :: Core", Desc: wrapText(ar.botDesc, 76),
				Exits: map[string]string{"up": midID}}

			if ar.hub { // the zone's safe access shell — no shops/medics in the Net (gear up in meatspace)
				top.Safe = true
			} else { // recon ICE patrols the shell
				mobs = append(mobs, netMob("c", z.band, topID+"_m", "a recon-ICE sentry", topID))
			}
			midKind := "c"
			if ar.boss != "" {
				midKind = "e" // a tougher guardian fronts the climax boss
			}
			mobs = append(mobs, netMob(midKind, z.band, midID+"_m", ar.midFoe, midID))
			if ar.boss != "" {
				mobs = append(mobs, netMob("b", z.band, botID+"_m", ar.boss, botID))
			} else if botID != netGauntletHome {
				mobs = append(mobs, netCacheMob(z.band, botID+"_c", botID))
			}
			// (the multi-stage Gauntlet ICE is homed at netGauntletHome by buildMobTemplates)

			rooms = append(rooms, top, mid, bot)
			mids = append(mids, mid)
		}
	}

	// Chain every area's MID layer end to end with varied cardinal directions —
	// the Net's lateral thoroughfare across all 50 nodes (and across zones).
	dirPat := []string{"north", "east", "south", "west", "east", "north", "west", "south", "north", "south", "east", "west"}
	di := 0
	nextDir := func(avoid string) string {
		for k := 0; k < len(dirPat)*2; k++ {
			d := dirPat[di%len(dirPat)]
			di++
			if d != avoid {
				return d
			}
		}
		return "north"
	}
	prevBack := ""
	for i := 0; i+1 < len(mids); i++ {
		d := nextDir(prevBack)
		mids[i].Exits[d] = mids[i+1].ID
		mids[i+1].Exits[opposite(d)] = mids[i].ID
		prevBack = opposite(d)
	}
	return rooms, mobs
}

// netMob is a band-scaled Net hostile: an ICE construct that shatters into broken
// shards on defeat. Bosses drop RAM (the netrunner's combat resource).
func netMob(kind string, band int, id, name, home string) *MobTemplate {
	t := mobFor(kind, band, id, name, home)
	t.ICE = true
	if kind == "b" {
		t.Drops = map[string]int{ramFor(band): 2}
	}
	return t
}

// netCacheMob is a breakable data-cache in a Core layer: passive, low HP, shatters
// into shards dropping a band-scaled RAM consumable + scrip, refilling on cooldown.
func netCacheMob(band int, id, home string) *MobTemplate {
	return &MobTemplate{ID: id, Name: "a sealed data-cache", HP: 5 + band, Damage: 0, AC: 1,
		XP: 5 + band*3, Eddies: 15 + band*22, Aggressive: false, ICE: true, Container: true, Home: home,
		Drops: map[string]int{ramFor(band): 1}}
}
