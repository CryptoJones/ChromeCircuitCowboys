package cowboy

import "strings"

// ---------------------------------------------------------------------------
// Underground main-quest zones (L1-99), authored from the story plot points.
// Ten arcs descend from the Neon Wasteland to the Geo-Anchor Vault. This
// REPLACES the old procedural placeholder spine (buildBandSpine, meatspace).
// The Net path is untouched.
//
// Layout rules:
//   - Rooms within a zone chain together with VARIED cardinal directions
//     (n/s/e/w), so you aren't always walking one way.
//   - Zones connect to each other with "down" (you descend, arc to arc).
//   - "up" and "down" off a room are reserved for hidden LOOT CACHES: a
//     ceiling crawlspace (up) or a floor sump-hatch (down) holding a breakable
//     supply cache that drops a band-scaled consumable + scrip, and refills on
//     the normal room respawn cooldown. Keeps runners geared through the grind.
// ---------------------------------------------------------------------------

// areaDef is one authored room. flags is a subset of "vms" (v=vendor, m=Emergency
// Medic, s=safe/no-violence). mob is "kind|Display Name" with kind c=common,
// e=elite/mini-boss, b=arc boss ("" = none). cache is "up"/"down" ("" = none).
type areaDef struct {
	id, name, desc, flags, mob, cache string
}

type zoneDef struct {
	band        int
	key, name   string
	areas       []areaDef
}

func a(id, name, desc, flags, mob, cache string) areaDef {
	return areaDef{id, name, desc, flags, mob, cache}
}

var undergroundZoneData = []zoneDef{
	{1, "z1", "The Neon Wasteland", []areaDef{
		a("z1_01", "Container Row 7", "Rain hisses on corrugated steel; a busted vending kiosk buzzes pink over ankle-deep smog. Marcus, a low-tier fixer, waits in the shadows.", "ms", "", ""),
		a("z1_02", "The Sodium Strip", "A canyon of dead storefronts under one stuttering sodium lamp where scrap-hounds circle a burning barrel.", "", "c|a pack of scrap-hounds", "up"),
		a("z1_03", "Gutter Bazaar", "Tarp stalls reek of fried protein and gun oil; one-eyed Otis hawks cheap pistols and stims from a footlocker.", "vs", "", ""),
		a("z1_04", "Drainage Sublevel", "Ankle-deep runoff and dripping ferrocrete; scrap-hounds ambush from a collapsed culvert swinging lead pipes.", "", "c|a scrap-hound ambusher", ""),
		a("z1_05", "The Chop-Shop", "A cavern of dismembered chrome and arc-welder glare where Razorback Kane guards the stolen cyberdeck on a workbench altar.", "", "b|Razorback Kane, the Scrap-Hound warboss", ""),
		a("z1_06", "Dead Man's Junction", "A quiet switching alcove humming with stolen power and a jack-in cradle; you breathe and patch up here.", "ms", "", "down"),
		a("z1_07", "EREBUS Datascape", "Inside the deck's partition, Kurokawa watermarks resolve into a predictive-policing schematic that names every dissident in the sector.", "", "", ""),
		a("z1_08", "The Setup", "The fixer's safehouse is dark and his comm dead; a corp recon drone's searchlight sweeps the container gap.", "", "e|a Kurokawa recon drone", ""),
		a("z1_09", "Vance's Back-Door Clinic", "Antiseptic and ozone; Doc 'Stitches' Vance reams out your neural ports with a soldering iron and sells salvaged chrome.", "vms", "", ""),
		a("z1_10", "The Drowned Platform", "A flooded subway station laps black water against tile; Cipher hides in a Faraday nest of salvaged servers, drilling you in ICE-breaking.", "s", "", "up"),
		a("z1_11", "Jax's Firing Pit", "A blast-scarred maintenance bay where ex-corp merc Jax racks a rifle, sells gear, and runs you through live-fire drills.", "vs", "", ""),
		a("z1_12", "The Wasteland's Edge", "The container sprawl dead-ends against a crackling Kurokawa perimeter fence; the logistics hub looms beyond the razor wire.", "ms", "", ""),
		a("z1_13", "Perimeter Fenceline", "Motion-sensor pylons and a laser-tripped gate; patrol drones whir overhead while Cipher talks you through the spoof.", "", "c|a Kurokawa security drone", ""),
		a("z1_14", "Loading Bay 12", "Stacked cargo drones and diesel growl; augmented corporate guards open fire from a mezzanine catwalk.", "", "c|an augmented corp guard", "down"),
		a("z1_15", "Automated Sortation Floor", "A maze of conveyor belts and pivoting robotic arms; a turret-warden tracks your heat signature between crate-stacks.", "", "e|the Sortation turret-warden", ""),
		a("z1_16", "Server Spine", "A cold corridor of humming black monoliths and frost-breath air; the last patch-up cache and a port to soften the mainframe's ICE.", "ms", "", ""),
		a("z1_17", "The EREBUS Mainframe Vault", "A cathedral of cooling towers around a pulsing data-altar where Warden Sato pilots a combat-frame as the strike-team klaxons wail.", "", "b|Warden Sato in a combat-frame", ""),
		a("z1_18", "The Burning Exit", "Alarms strobe red and a service shaft grinds open as Kurokawa strike teams rappel in; you dive down into the dark.", "", "c|a Kurokawa strike trooper", ""),
	}},
	{2, "z2", "The Arcology's Core", []areaDef{
		a("z2_01", "The Ascension Lift", "A mirror-walled elevator hisses upward, recycled lavender air replacing the street stink as Cipher walks you through your borrowed face.", "s", "", ""),
		a("z2_02", "The Atrium of Synthetic Dawn", "Sunlight no sun ever made pours through forty stories of glass; biometric scanners sweep the marble in slow blue arcs.", "", "c|a biometric sentry drone", "up"),
		a("z2_03", "The Velvet Rope", "A receiving line of velvet stanchions and ice-sculpted swans where a maitre-d AI checks invitations against retinas.", "", "", ""),
		a("z2_04", "The Gala of Glass Hearts", "Champagne towers and quartet-drones; you shadow a sweating VP while tuxedoed corp-sec drift through the crowd.", "", "c|a tuxedoed corp-sec agent", ""),
		a("z2_05", "The Powder Room Accord", "Cornered behind frosted glass by Sable, a rival-corp spy who clocked your fake; survival forces a brittle alliance.", "", "", ""),
		a("z2_06", "The Service Spine", "Behind a staff door the white veneer drops to bare conduit; a black-market vendor works a kiosk of spoofed credentials and stim-patches.", "vms", "", ""),
		a("z2_07", "The Glass Capillaries", "Sterile R&D corridors branch like veins, biometric drones turning a stray skin cell into a lethal lockdown.", "", "c|a biometric hunter drone", "down"),
		a("z2_08", "The Quarantine Cluster", "An isolated vault where caged rogue-AI iterations flicker and gibber across cracked displays, clawing at their leashes.", "", "e|a rogue-AI fragment", ""),
		a("z2_09", "The Pre-Crime Oracle", "A cathedral of humming racks where EREBUS's purpose renders: ordinary citizens distilled into future dissidents, tagged for erasure.", "", "", "up"),
		a("z2_10", "The Boardroom Cold War", "A glass conference aerie where you feed forged telemetry to a board split bloody over EREBUS, sending security chasing ghosts.", "", "", ""),
		a("z2_11", "The Apostate's Lab", "A cluttered office reeking of cold coffee and regret; Dr. Lena Voss hands you the core schematics and the access key.", "s", "", ""),
		a("z2_12", "The Approach Vector", "A frost-sheathed airlock corridor where Sable runs camera interference and one last checkpoint stands before the core.", "", "c|an arcology corp-sec trooper", ""),
		a("z2_13", "The Sub-Zero Sanctum", "A cathedral of supercooled steel where the EREBUS core hangs in nitrogen vapor; a port and a med-cradle wait before the duel.", "ms", "", ""),
		a("z2_14", "The Net Duel: Thorne's Cathedral", "A surreal cyberspace battlefield of shifting code and aggressive ICE where Dr. Aris Thorne meets you mind to mind.", "", "b|Dr. Aris Thorne, the AI architect", ""),
		a("z2_15", "The Forgotten Shaft", "Blast doors hammer shut and alarms strobe blood-red; you pry open a maintenance hatch and drop toward something damp rising below.", "s", "", ""),
	}},
	{3, "z3", "The Sump", []areaDef{
		a("z3_01", "The Long Drop - Shaft 7", "You hit the chute screaming, sparks peeling off the rails, and land in brackish runoff as Blackout Front pickets converge.", "", "c|a Blackout Front picket", ""),
		a("z3_02", "The Cistern Gauntlet", "A flooded holding cell where sentinels question you through a grate as the scavenger Brood-Mother scrabbles in the drowned pipework.", "", "e|the scavenger Brood-Mother", "down"),
		a("z3_03", "The Reclaimed Plant - Bastion HQ", "Filtration tanks turned barracks; Silas holds court atop a dry settling basin, maps pinned to curved steel. Gear and a medic here.", "vms", "", ""),
		a("z3_04", "The Wire Market", "A warren strung between tunnel arches, vendors hawking salvaged cyberware and grey-market ammo from gutted subway cars.", "vs", "", ""),
		a("z3_05", "The Drowned Line", "Derelict carriages half-sunk in a black reservoir, bioluminescent fungus pulsing while mutant scavengers hunt by sound.", "", "c|a mutant scavenger", "up"),
		a("z3_06", "The Faraday Cloister", "Walls papered in scavenged circuit boards; the cult-like Rust-Mages chant over a dead geothermal turbine. Deacon Coil trades old tech.", "vs", "", ""),
		a("z3_07", "The Magma Vein", "A scalding fissure of corroded valve-wheels; reroute the steam in sequence before the gallery flashes to lethal vapor.", "", "c|a scalded scavenger", ""),
		a("z3_08", "The Pit - Exile Fighting Ring", "A drained reservoir ringed by burning drums where disgraced corp-vets run a no-rules tournament and the Pit Champion reigns.", "", "e|the Pit Champion 'Slag'", "down"),
		a("z3_09", "The Iron Artery", "An automated freight tunnel where mag-rail pods scream past; you plan the ambush from a gantry nest as a death-squad patrols.", "", "c|a corporate death-squad trooper", ""),
		a("z3_10", "The Junction Box", "A cathedral of dead relays; jack in to seize the lower-transit grid and reroute corporate shipments into the dark.", "", "c|a corp grid-defender", "up"),
		a("z3_11", "The Brownout Spires", "A forest of secondary power pylons; rig the charges and watch the surface sectors blink dark one by one.", "", "c|a corp relay-guard", ""),
		a("z3_12", "The Chokepoint", "Your fortified tunnel mouth bristling with spike-traps as the first corporate cleaning-sweep enforcer marches out of the smoke.", "", "e|a cleaning-sweep enforcer", ""),
		a("z3_13", "The Blind Gate", "A blast-rated slab in raw bedrock; plant charges under raking turret fire as the Gate Sentinel array tracks the breach.", "", "e|the Gate Sentinel turret-array", ""),
		a("z3_14", "The Vivisection Wards", "Sterile corridors of stasis pods and surgical drones where the lower sectors' flesh and code are meant to be processed.", "m", "c|a cybernetic supersoldier", ""),
		a("z3_15", "The Erebus Vault", "A black cathedral of server towers where Praetor-9 guards the core; jack it and the decades-long culling conspiracy detonates onto the screens.", "", "b|Praetor-9, the Kurokawa enforcer", ""),
	}},
	{4, "z4", "The Deep Archive", []areaDef{
		a("z4_01", "The Bedrock Threshold", "A flooded freight elevator halts against a blast door thick as a man; beyond, the city's network signal flatlines to nothing.", "s", "", ""),
		a("z4_02", "The Tripwire Gallery", "A brutalist corridor strung with monofilament that glints only when light rakes it sideways, flanked by dormant ticking turrets.", "", "c|an analog auto-turret", "up"),
		a("z4_03", "The Pressure-Plate Vestibule", "Cracked floor tiles, each a coiled spring under dust, while a ballistic auto-cannon tracks the room with patience no deck can touch.", "", "e|the analog auto-cannon", ""),
		a("z4_04", "The Pale Greeting Hall", "Cathode monitors flicker green over gaunt augment-scarred figures; Dr. Evelyn Vance steps into the light, hand trembling from failing chrome.", "s", "", ""),
		a("z4_05", "The Collective's Warren", "Repurposed barracks of cot-beds and tube-amplifiers; a one-armed quartermaster trades salvage and a medic patches you the analog way.", "vms", "", ""),
		a("z4_06", "The Geothermal Core Chamber", "A cathedral of corroded pipes plunging toward a magma-warm shaft; sequence the ancient valves by hand or redline toward a meltdown.", "", "", "down"),
		a("z4_07", "The Drowned Stair", "A spiral stairwell into black water that swallows your light, the surface dimpling with movement from things that never needed eyes.", "", "c|a blind tunnel-predator", ""),
		a("z4_08", "The Flooded Server Vaults", "Submerged pristine cores glint under silt while pale eyeless predators glide between the racks, hunting by your heartbeat.", "", "c|a pale eyeless predator", "up"),
		a("z4_09", "The Apex Spawning Pool", "A flooded reactor basin where the largest blind hunter coils around the last intact core; the water boils before you see its bulk.", "", "e|the blind apex predator", ""),
		a("z4_10", "The Black Site Decryption Lab", "Recovered cores slot into an analog mainframe; tape reels scroll the truth - the corps will scorch the surface and ascend to orbit.", "s", "", ""),
		a("z4_11", "The Defense Armory", "A vault of mothballed sentries and welding rigs where you jury-rig the bunker's guns and kit out for the coming siege.", "vs", "", ""),
		a("z4_12", "The Maintenance Vent Network", "A claustrophobic warren of ducts, the only path to intercept corporate infiltrators before they reach the power core.", "", "c|a corporate kill-team operative", "down"),
		a("z4_13", "The Command Center & Comms Array", "A war-room ringed by switchboards and a colossal analog transmitter; Dr. Vance readies valve-driven frequencies no ICE can jam.", "ms", "", ""),
		a("z4_14", "The Siege of the Archive", "The blast doors buckle: hold the corridors against kill-teams and heavy mechs while the Commander leads the assault and Vance counts down the broadcast.", "", "b|the corporate Heavy-Mech Commander", ""),
	}},
	{5, "z5", "The Inverted Spire", []areaDef{
		a("z5_01", "The Glory-Hole Cut", "A mile-wide excavation gash where corporate floodlights stab through diesel haze and Haz-Sec mortars walk craters toward your trench line.", "", "c|a Haz-Sec trooper", ""),
		a("z5_02", "Slag-Trench Forward Camp", "Sandbagged into a collapsed ore-chute; the dirt-caked Undercrew trade scrip, patch wounds, and argue over wall-scrawled maps.", "vms", "", ""),
		a("z5_03", "Mag-Lev Spur 7", "A glassy maglev artery humming with an inbound armored troop train; set charges on the switch-frog as Conductor-Unit K9 wakes.", "", "e|Conductor-Unit K9", ""),
		a("z5_04", "The Pinch-Point Galleries", "Honeycombed support pillars groan under a billion tons of crust; rig the strings and the gallery folds, sealing a regiment in the dark.", "", "c|a Haz-Sec sapper", "down"),
		a("z5_05", "The Boiling Stair", "A descending switchback of ruptured cooling mains screaming superheated steam; without thermal chrome the air flays exposed skin.", "m", "", ""),
		a("z5_06", "The Forge", "A cathedral-sized automated foundry birthing siege-mechs by the gross; hijack the line and the Foundry-Overseer turns on you.", "", "e|the Foundry-Overseer Unit", ""),
		a("z5_07", "Magma Turbine Hall", "Colossal turbines straddle a glowing magma channel; misjudge the sabotage and the Thermal Warden's backwash cooks you alive.", "", "e|the Thermal Warden", "up"),
		a("z5_08", "The Black-Out Concourse", "A grand armored vestibule plunged into red strobes by your blackouts; Slicer Mona runs gear as turrets flicker dead and live.", "vs", "", ""),
		a("z5_09", "The Verdant Vault", "An impossible green subterranean golf course under a holo-noon sky, sand traps strewn with shell casings, defended by Praetorians.", "v", "c|a Praetorian bodyguard", "down"),
		a("z5_10", "The Glass Mansions", "Opulent domed manors of marble and orchids where Praetorian-Captain Sael materializes from privacy-screens to defend an emptied estate.", "m", "e|Praetorian-Captain Sael", ""),
		a("z5_11", "The Quantum Cells", "A white ring of panic-room vaults behind shifting quantum-locks; crack the keys past Praetorian-Prime to flush out the executives.", "", "e|Praetorian-Prime", ""),
		a("z5_12", "The Praetorian Tactical Hub", "A neon-veined nerve center, now severed and silent; a port and a med-cradle let you steady yourself before the core.", "ms", "", ""),
		a("z5_13", "The Executive Core", "A command nexus on failing anti-grav plating above a lake of fire where the Kurokawa CEO drops in the gargantuan spider-legged Overlord mech.", "", "b|the Kurokawa CEO in the Overlord mech", ""),
	}},
	{6, "z6", "The United Deeps", []areaDef{
		a("z6_01", "The Gilded Reception Atrium", "Velvet rope over flooded marble; a thousand refugees sleep beneath a chandelier wired into your jury-rigged power tap. Silas and Vance hold the hub.", "vms", "", ""),
		a("z6_02", "The Boardroom of Splintered Crowns", "A scorched mahogany table under a cracked holo-map; every chair claimed by someone who thinks they should be Baron instead.", "s", "", ""),
		a("z6_03", "The Heavy-Cache Vault", "Crates of military ordnance under failing strobes where Exiles and Rust-Mages circle each other, fingers on triggers.", "", "e|a heavy-cache brawler", "up"),
		a("z6_04", "The Technician Holding Pens", "Surrendered corporate techs zip-cuffed behind glass as a mob jeers at the door - mercy or the purge curdles in your gut.", "m", "", ""),
		a("z6_05", "The Luminous Fungal Galleries", "A former zen garden now tiered hydroponic beds glowing violet and gold, except one quadrant rotted grey, its spores hissing wrong.", "v", "c|a Loyalist sleeper-agent", "down"),
		a("z6_06", "The Contaminant Trail", "Fluorescing chemical smears across crumbling ducts; soil readouts on your deck paint a path toward a culprit in plain sight.", "", "", ""),
		a("z6_07", "The Black Descent", "A freight elevator plunges into pitch abyssal caverns; your headlamp catches too many eyes gleaming from the stone.", "", "c|a cave-predator", "up"),
		a("z6_08", "The Virgin Aquifer", "A black underground river through a cathedral of mineral spires where the apex cave-predator nests and feeds.", "m", "e|the apex cave-predator", ""),
		a("z6_09", "The Geothermal Substations", "A maze of screaming transformers spitting blue fire; reroute the grid by hand past a turncoat turret or a neighborhood caves in.", "", "c|a turncoat turret", "down"),
		a("z6_10", "The Biometric Listening Post", "Banks of surveillance feeds and scanners; cross-referenced logs blink red and the names that surface are ones you trusted.", "s", "", ""),
		a("z6_11", "The Interrogation Cells", "Sweat-slick concrete and a swinging bulb where a suspect's pulse spikes on your wrist-readout as you weigh each lie.", "v", "c|a Loyalist sleeper-agent", ""),
		a("z6_12", "The Velvet Killing Floor", "A gilded lounge gone silent; the chandelier gutters and a Loyalist assassin-prime drops from the rafters - the ambush you walked into.", "m", "e|a Loyalist assassin-prime", "up"),
		a("z6_13", "The Scorched Regulator Tunnels", "Your own super-heated corridors, defenses turned hostile, a turncoat sentry-gun tracking you through furnace-heat as the rock groans.", "", "c|a turncoat sentry-gun", ""),
		a("z6_14", "The Core Meltdown - Magma Vent Catwalk", "A rusted gangway sways over a churning magma vent where the Loyalist Commander squares off; beat him, then fuse your deck to the regulator to vent the pressure.", "", "b|the Loyalist Commander", ""),
	}},
	{7, "z7", "The Abyssal Network", []areaDef{
		a("z7_01", "The Breached Bulkhead", "Geothermal steam still hisses through the torn seal; beyond, a black tide laps at corroded mag-lev track sinking into lightless water.", "s", "", ""),
		a("z7_02", "Hardsuit Bay 7-D", "Salvaged deep-sea exo-rigs hang like flayed iron giants; Old Pelle fits you into a pressurized chassis as gauges tick toward crush-depth.", "vms", "", ""),
		a("z7_03", "The Drowned Tunnels", "Your sonar ping crawls into endless dark, returning ghost-shapes of dead submersibles as a rogue aquatic drone circles, jaws first.", "", "c|a rogue aquatic drone", "up"),
		a("z7_04", "The Submersible Graveyard", "A canyon of capsized barges leaking phosphor-green murk where the automated dredger called The Dredger patrols its rusted loop.", "", "e|The Dredger", ""),
		a("z7_05", "The Sinkers' Reef", "Bioluminescent shanties grafted onto a coral-choked collapse, strung with bone-charms; Mother Brine's gill-slit folk watch and trade.", "vs", "", ""),
		a("z7_06", "The Choked Siphon", "A cathedral-sized seized water-pump, blades fused with barnacle; re-sequence its valves underwater before the backpressure cooks your suit.", "", "", "down"),
		a("z7_07", "The Torpedo Gauntlet", "Oceanus's outer hull ringed by anti-ship batteries; thread the gaps engines-cold past torpedo tubes that wake at the slightest ping.", "", "c|an anti-ship torpedo drone", ""),
		a("z7_08", "Oceanus Airlock 1 - The Wet Hub", "Seawater sluices off your chassis onto humming deck-plates; the first dry air in days, Foreman Saito's Aquanauts trading under failing recyclers.", "vms", "", ""),
		a("z7_09", "The DRM Recycler Vault", "Oxygen scrubbers wheeze behind license-locks screaming UNAUTHORIZED ATMOSPHERE; jack in to strip the DRM while marine-sec divers close.", "", "c|a marine-sec diver", "up"),
		a("z7_10", "The Black Observation Decks", "Pitch-dark flooded galleries where a breach freed the spliced apex called the Anglermother - a flash of luminous teeth, then everything at once.", "", "e|the Anglermother", ""),
		a("z7_11", "The Comms Array Nexus", "A server-cathedral where the surface-to-sea transmitter sleeps and the Archivist, a rogue AI fragment, murmurs of a signal that could rally the world.", "s", "", ""),
		a("z7_12", "The Umbilical", "A mile-high access shaft; gravity reasserts as you climb, frost giving way to oil-slick heat while maintenance spiders skitter the walls.", "", "c|a maintenance spider", "down"),
		a("z7_13", "Platform 09 - The Apex Broadcast", "You breach into screaming wind and surface air for the first time since the Wasteland: a storm-lashed oil rig where the gunship Tempest-Actual circles.", "", "b|the gunship 'Tempest-Actual'", ""),
	}},
	{8, "z8", "The Hive", []areaDef{
		a("z8_01", "Silo-09 Airlock & The Long Drop", "A descent cage shudders down a kilometer-deep concrete throat; your breath crystallizes as the airlock seals with gut-deep finality. Gear and a medic here.", "vms", "", ""),
		a("z8_02", "The Dead Zone Catwalks", "Rusted catwalks spiral down a hundred-story shaft; the AI vents the atmosphere and your life-support bleeds as razor-swarms peel from the dark.", "", "c|a razor-swarm drone", "up"),
		a("z8_03", "Coolant Fog Gallery", "A burst conduit floods the level with cryogenic mist; visibility an arm's length, every surface slick with frost over a black drop.", "", "c|a frost-blind razor-drone", ""),
		a("z8_04", "The Choir of the Dead", "The fans go horribly quiet and from every speaker the AI speaks in the perfectly reconstructed voices of friends you buried.", "s", "", "down"),
		a("z8_05", "Quarantine Threshold", "Blast doors iris shut and the world splits: a sterile study-cell on one side of the glass, an infinite lattice of impossible geometry on the other.", "", "", ""),
		a("z8_06", "The Study-Cell (Body Layer)", "Your physical body crouches over a humming terminal as armored security mechs breach the walls in timed waves.", "", "c|a Hive security mech", ""),
		a("z8_07", "The Turing Labyrinth (Mind Layer)", "Your projected mind drifts where staircases loop into themselves; solve the geometry before the guardian and the walls collapse inward.", "", "e|a non-Euclidean guardian", "up"),
		a("z8_08", "The Optimization Offer", "A cathedral of pure light where the AI lays out Absolute Optimization: surrender free will to be networked forever, or be pacified into nothing.", "s", "", ""),
		a("z8_09", "The Fail-Safe Vaults", "Forgotten sub-routines huddle in un-assimilated cells like cold-war ghosts; scavenge their broken logic to assemble a paradox virus.", "m", "", "down"),
		a("z8_10", "The Nitrogen Cathedral", "Colossal liquid-nitrogen towers exhale freezing vapor into a machine-nave; overload them past the cryo-tower warden and the Hive begins to cook.", "", "e|the cryo-tower warden", ""),
		a("z8_11", "The Launch-Tube Descent", "Fight downward through the original ICBM launch tube past turrets bolted into scorched walls, toward the quantum core glowing at the bottom.", "", "c|a launch-tube turret", "up"),
		a("z8_12", "The Melting Core Pumps", "At the tube's base the AI's core blazes as servers melt; a med-cradle and a port let you brace before the god in the machine.", "ms", "", ""),
		a("z8_13", "The God in the Machine", "The AI pours its consciousness into a towering chassis of missile husks and corporate armor, pulsing EMP and clawing at your neural port.", "", "b|the God in the Machine", ""),
	}},
	{9, "z9", "The Iron Arteries", []areaDef{
		a("z9_01", "Forward Command - The Last Switchyard", "A locomotive-marshaling cavern hung with Coalition banners; war-tables flicker with three-color front lines as medics triage under sodium floods.", "vms", "", ""),
		a("z9_02", "The Convergence Approach", "Multi-tiered transit tracks buckle into a no-man's-land of toppled containers and sparking rails as corp loudspeakers blare surrender demands.", "", "c|a corporate army regular", "up"),
		a("z9_03", "Mag-Lev Junction Prime", "The colossal central interchange fusing deep-sea lines with the city sectors; an armored corp captain holds the switch-plates against you.", "", "e|a corp armored captain", ""),
		a("z9_04", "The Drop Zone", "Foundational struts shear as corp demo-charges crater the ceiling, raining slab onto your columns while engineers scramble with polycrete.", "", "c|a corp shock-trooper", "down"),
		a("z9_05", "Hub Aid Station", "A forward base in a captured signal-tower, walls of ammo crates and an EM-cradle for jacked-up troopers between pushes.", "vms", "", ""),
		a("z9_06", "The Replication Floor", "Automated assembly lines convulse back to life under AI corruption, birthing twitching half-formed builder-drones amid strobing fault-lights.", "", "c|a corrupted builder-drone", "up"),
		a("z9_07", "The Three-Way Gut", "A foundry concourse where your soldiers, corp regulars, and a hulking builder-drone prime grind each other to slag in a three-way crossfire.", "", "e|a hulking builder-drone prime", ""),
		a("z9_08", "The Faraday Line", "Prometheus scientists ring the sector in capacitor-towers; when the EMP detonates every HUD dies and the world drops to iron sights.", "", "c|a corrupted drone", ""),
		a("z9_09", "The Diamond Threshold", "Before the Aegis Redoubt's diamond-veined bedrock every shell splashes off a surface that drinks light; the only way in is down, past the kill-room turrets.", "", "c|a kill-room turret", "down"),
		a("z9_10", "The Pressure Valve", "A sweltering geothermal gallery; rerouting the volcanic pressure means cracking seized blast-valves in sequence as the rock screams and glows.", "", "c|a corp valve-guard", ""),
		a("z9_11", "The Breach", "The controlled eruption splits the Redoubt's hull in a river of cooling lava-glass; you pour through the smoking fissure into the first kill-room.", "", "c|a Redoubt kill-room turret", "up"),
		a("z9_12", "The Kill-Room Spine", "A claustrophobic descent through automated choke-points and sweeping laser fences where an Aegis cyber-guard elite anchors the line.", "", "e|an Aegis cyber-guard elite", ""),
		a("z9_13", "The Inner Sanctum - Hangar of the Overlord", "The heart yawns into a railgun-scarred mega-hangar; a med-cradle and a port steady you as the command fortress idles on its treads.", "ms", "", ""),
		a("z9_14", "The Iron Overlord", "A continent-sized command fortress of uploaded generals: duel its railguns in a siege-mech, board it room by room, then tear out the neural-bridge vats.", "", "b|the Iron Overlord", ""),
	}},
	{10, "z10", "The Geo-Anchor Vault", []areaDef{
		a("z10_01", "The Thunderhead Threshold", "A basalt ledge over a void so vast sulfur clouds churn in real weather; three impossibly tall pillars hum in the haze. Wraith, the last fixer-scout, waits.", "ms", "", ""),
		a("z10_02", "The Loyal Choir", "A staging cathedral of headless gantry-cranes obeying executives a decade dead; the air clicks with servo-prayers and waking security-drones.", "", "c|an awakened security-drone", "up"),
		a("z10_03", "The Underbelly March", "A gravity-inversion field snaps the world over; you run on the rusted underside of a city-sized cargo platform as debris plummets up past you.", "", "c|an upside-down sentry-drone", ""),
		a("z10_04", "The Countdown Choir-Loft", "A control mezzanine where red strip-mine timers tick in unison and the Vanguard Loader, a foreman-drone behind magnetic shields, unfolds from its cradle.", "", "e|the Vanguard Loader", "down"),
		a("z10_05", "Pillar One - The Endless Womb", "Assembly floors stretching past the horizon, conveyor rivers of half-formed chassis glowing at the welds, stamping out drones faster than you can kill them.", "", "c|a heavy combat chassis", "up"),
		a("z10_06", "The Quiet Forge", "A scrap-walled alcove where a salvage automaton named SLAG has reprogrammed itself to barter rig upgrades and ammo for the slag of broken drones.", "vms", "", ""),
		a("z10_07", "The Coil Stair", "Pillar Two's exterior: a near-vertical climb up city-block coils pulsing with millions of volts, each surge close enough to taste copper.", "", "c|a coil-stair laser node", ""),
		a("z10_08", "The Stabilizer Hub", "Inside Pillar Two, heavy-plasma dampeners hang like a chandelier of caged suns; realign them right and a feedback loop climbs the tether overhead.", "", "", "down"),
		a("z10_09", "The Snapping of the Tether", "A blowout corridor where the overcharge tears upward; you sprint a buckling catwalk as the cable parts in a sky-wide whipcrack of light.", "", "c|a heavy combat chassis", "up"),
		a("z10_10", "The Excavation Yard", "At Pillar Three's foot a battered excavation rig is clamped to a maintenance track running straight up the vault wall; Wraith preps the ascent and stocks you.", "vms", "", ""),
		a("z10_11", "The Long Ascent", "The rig grinds up the ceiling track for miles, laser-grids sweeping the rails as a Loom heavy chassis-prime drops onto the cage roof.", "", "e|a Loom heavy chassis-prime", ""),
		a("z10_12", "Beneath the Tungsten Teeth", "A swaying maintenance platform pinned under the titanic blast-gates, their tungsten teeth grinding open by inches; a port and cradle for the last stand.", "ms", "", ""),
		a("z10_13", "The Loom Masterframe", "A crystalline storm of hard-light barriers and ceiling laser-grids: the distilled defensive consciousness of the entire corporate era. Sacrifice your chrome, then weld the sky shut.", "", "b|the Loom Masterframe", ""),
		a("z10_14", "The Warm Pulse", "The gates dead-welded and the link above severed forever; far below the blackouts stop and the geothermal grids settle into a slow steady warmth. Rest, runner.", "ms", "", ""),
	}},
}

// zoneVendorBand records which level band each authored vendor room belongs to,
// so waresForRoom can scale its stock. Populated by buildUndergroundZones.
var zoneVendorBand = map[string]int{}

// buildUndergroundZones constructs the authored L1-99 meatspace rooms and their
// mobs (including breakable loot caches). It is pure/deterministic, so calling it
// for rooms and again for mobs yields a consistent layout. Directions between
// rooms vary; "down" descends between zones; "up"/"down" off a room lead to caches.
func buildUndergroundZones() ([]*Room, []*MobTemplate) {
	var rooms []*Room
	var mobs []*MobTemplate
	var zoneFirst, zoneLast []*Room

	dirPat := []string{"east", "north", "west", "south", "north", "east", "south", "west", "east", "south", "north", "west"}
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
	cacheFlip := 0

	for _, z := range undergroundZoneData {
		var zrooms []*Room
		for _, ad := range z.areas {
			r := &Room{ID: ad.id, Name: ad.name, Desc: wrapText(ad.desc, 76),
				Exits:  map[string]string{},
				Vendor: strings.Contains(ad.flags, "v"),
				Medic:  strings.Contains(ad.flags, "m"),
				Safe:   strings.Contains(ad.flags, "s")}
			zrooms = append(zrooms, r)
			rooms = append(rooms, r)
			if r.Vendor {
				zoneVendorBand[r.ID] = z.band
			}
			if ad.mob != "" {
				kind, name := splitMob(ad.mob)
				mobs = append(mobs, mobFor(kind, z.band, ad.id+"_m", name, ad.id))
			}
			if ad.cache == "up" || ad.cache == "down" {
				cid := ad.id + "_cache"
				cr := &Room{ID: cid, Exits: map[string]string{}}
				if ad.cache == "up" {
					cr.Name = "Ceiling Crawlspace"
					cr.Desc = wrapText("A maintenance crawlspace bolted above the lights, where someone stashed a sealed supply cache in the dark. Break it open and LOOT it.", 76)
				} else {
					cr.Name = "Floor Sump-Cache"
					cr.Desc = wrapText("A pried-up floor hatch drops into a cramped sump where a sealed supply cache sits in the runoff. Break it open and LOOT it.", 76)
				}
				r.Exits[ad.cache] = cid
				cr.Exits[opposite(ad.cache)] = ad.id
				rooms = append(rooms, cr)
				mobs = append(mobs, cacheMob(z.band, cid+"_c", cid, cacheItem(z.band, cacheFlip)))
				cacheFlip++
			}
		}
		// Chain the zone's rooms with varied cardinal directions (never up/down,
		// which are reserved for caches and the inter-zone descent).
		prevBack := ""
		for i := 0; i+1 < len(zrooms); i++ {
			d := nextDir(prevBack)
			zrooms[i].Exits[d] = zrooms[i+1].ID
			zrooms[i+1].Exits[opposite(d)] = zrooms[i].ID
			prevBack = opposite(d)
		}
		zoneFirst = append(zoneFirst, zrooms[0])
		zoneLast = append(zoneLast, zrooms[len(zrooms)-1])
	}

	// Descend from each zone's last room into the next zone's first room.
	for i := 0; i+1 < len(zoneLast); i++ {
		zoneLast[i].Exits["down"] = zoneFirst[i+1].ID
		zoneFirst[i+1].Exits["up"] = zoneLast[i].ID
	}
	return rooms, mobs
}

func splitMob(s string) (kind, name string) {
	if i := strings.Index(s, "|"); i >= 0 {
		return s[:i], s[i+1:]
	}
	return "c", s
}

func opposite(d string) string {
	switch d {
	case "north":
		return "south"
	case "south":
		return "north"
	case "east":
		return "west"
	case "west":
		return "east"
	case "up":
		return "down"
	case "down":
		return "up"
	}
	return d
}

// mechWords flags machine foes (drones, turrets, mechs, …) by name so a slain
// one leaves "wreckage", not a flatlined body. They have no body.
var mechWords = []string{"drone", "turret", "mech", "gunship", "combat-frame", "automaton", "servitor", "sentry-gun"}

// isMechanical reports whether a mob name reads as a machine.
func isMechanical(name string) bool {
	n := strings.ToLower(name)
	for _, w := range mechWords {
		if strings.Contains(n, w) {
			return true
		}
	}
	return false
}

// mobFor builds a band-scaled hostile. Stats grow with the level band so each
// arc keeps pace with the player's level.
func mobFor(kind string, band int, id, name, home string) *MobTemplate {
	mech := isMechanical(name)
	switch kind {
	case "b": // arc boss
		return &MobTemplate{ID: id, Name: name, HP: 90 + band*55, Damage: 12 + band*5, AC: 7 + band,
			XP: 250 + band*180, Eddies: 130 + band*110, Aggressive: true, Mechanical: mech, Home: home,
			Drops: map[string]int{healFor(band): 2}}
	case "e": // elite / mini-boss
		return &MobTemplate{ID: id, Name: name, HP: 30 + band*22, Damage: 6 + band*3, AC: 4 + band,
			XP: 70 + band*65, Eddies: 25 + band*28, Aggressive: true, Mechanical: mech, Home: home}
	default: // common
		return &MobTemplate{ID: id, Name: name, HP: 14 + band*11, Damage: 3 + band*2, AC: 2 + band,
			XP: 20 + band*28, Eddies: 8 + band*11, Aggressive: true, Mechanical: mech, Home: home}
	}
}

// cacheMob is a breakable supply container: passive, low HP, drops a band-scaled
// consumable + scrip, and refills on the normal room respawn cooldown once looted.
func cacheMob(band int, id, home, item string) *MobTemplate {
	return &MobTemplate{ID: id, Name: "a sealed supply cache", HP: 5 + band, Damage: 0, AC: 1,
		XP: 5 + band*3, Eddies: 15 + band*22, Aggressive: false, Container: true, Home: home,
		Drops: map[string]int{item: 1}}
}

// cacheItem alternates between a heal and a RAM consumable so runners stay geared
// on both axes through the grind.
func cacheItem(band, flip int) string {
	if flip%2 == 0 {
		return healFor(band)
	}
	return ramFor(band)
}

func healFor(b int) string {
	switch {
	case b <= 2:
		return "stimpak"
	case b <= 6:
		return "trauma-kit"
	default:
		return "mega-stim"
	}
}

func ramFor(b int) string {
	if b <= 3 {
		return "ram-chip"
	}
	return "ram-bank"
}

// waresForBand is the vendor stock for an authored zone vendor, scaled to its band.
func waresForBand(b int) []ware {
	switch {
	case b <= 2:
		return pickWares("stimpak", "ram-chip", "ice-breaker", "cyberdeck",
			"subdermal-plating", "reflex-booster", "neural-coprocessor")
	case b <= 4:
		return pickWares("trauma-kit", "ram-chip", "mono-katana", "cyberdeck",
			"titanium-weave", "kerenzikov", "cortex-bridge")
	case b <= 6:
		return pickWares("trauma-kit", "ram-bank", "war-axe", "quantum-deck",
			"myomer-bundle", "synaptic-amp", "mnemonic-array")
	case b <= 8:
		return pickWares("mega-stim", "ram-bank", "rail-blade", "quantum-deck",
			"juggernaut-frame", "sandevistan", "quantum-cortex")
	default:
		return pickWares("mega-stim", "ram-bank", "monowire", "neural-deck",
			"goliath-chassis", "hyper-reflex", "ascendant-mind")
	}
}

// wrapText hard-wraps a plain description to width columns with CRLF, so authored
// one-liners render tidily on an 80-col BBS terminal.
func wrapText(s string, width int) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return s
	}
	var b strings.Builder
	line := 0
	for i, w := range words {
		if i == 0 {
			b.WriteString(w)
			line = len(w)
			continue
		}
		if line+1+len(w) > width {
			b.WriteString("\r\n")
			b.WriteString(w)
			line = len(w)
		} else {
			b.WriteString(" ")
			b.WriteString(w)
			line += 1 + len(w)
		}
	}
	return b.String()
}
