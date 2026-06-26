# Changelog

All notable changes to **Chrome Circuit Cowboys** are documented here.

Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
versioning: [SemVer](https://semver.org/spec/v2.0.0.html).

## [2.3.0] - 2026-06-26

Second backlog batch: party/co-op play, economy, QoL, and polish. (Items land
incrementally; see the issue links.)

### Added
- **`LOOK <item>` examines an item** — name, flavor, mechanical effect, and
  market value — instead of always describing the room. (#53)
- **`DROP` / `GET` floor items** — drop loot on the room floor (incl. `DROP ALL`)
  for crewmates to `GET` (or `GET ALL`); LOOK shows what's lying there. (#51)
- **Unsafe-logout penalty** — logging out somewhere unsafe means you got jumped
  offline: lose 5%% HP (never below 1) with a flavor line on return; safe-room
  logouts are free. (#48)
- **MAP points back to your crew** when separated — a "TO YOUR CREW: go <DIR>"
  pointer toward the nearest crewmate. (#52)
- **Party follow** — when the crew leader moves, members in the room follow
  along (those mid-combat stay to finish the fight). (#42)
- **`SELL` items for scrip** at vendors — offload unwanted gear for 50%%
  catalog buyback (`SELL <item> [qty]`). (#40)
- **Quick party-chat shortcut `;`** — `;<msg>` goes straight to crew chat
  (same as GSAY), like `'`=say and `:`=emote. (#43)
- **Hidden `roomid` command** (also `whereami`) prints the current room's
  internal id + exits — for building/debugging; not listed in HELP. (#28)
- **Character sheet spells out "Armor Class"** instead of the cryptic "AC". (#30)

### Changed
- **Character points always shown on the sheet** (0 when none), with a bold
  "You have character points to spend." call-out when you have some. (#31)
- **Renamed the "Back Alley" room to "Underground Entrance"** (the descent into
  the Undercity). (#29)

### Fixed
- **Containers no longer get "lunged at."** Opening/attacking a cache uses a
  randomized container verb (pry / crack / jimmy / force / …) instead of the
  combat "You lunge at …" line. (#46)
- **Lighter dim/ambience text** in the default scheme — the dark-grey hint text
  is now a more readable light grey. (#47)

## [2.2.0] - 2026-06-26

Backlog batch: QoL, content, and combat-flow improvements. (Items land
incrementally; see the issue links.)

### Added
- **`OPEN` command for caches.** Crack open a supply/data cache with `OPEN`
  instead of attacking it — the intuitive verb for an inert container. (#16)
- **Combat shortcuts.** `A` = attack, `LO` = loot (`L` still = look). (#14)
- **Numbered inventory + quick-use.** `INVENTORY` is numbered; pressing a digit
  at the prompt instantly USEs that slot (no Enter — fast for combat). (#13)
- **`TALK` to locals for lore.** Ask whoever's around (the hiring fixer, or a
  passer-by) about the level you're on — every underground arc and Net zone has
  its own backstory, surfaced in-game. (#12)
- **Clone-booth onboarding.** TALK in the Re-Clone Bay for a quick primer on
  the core commands from the booth tech, Doc Splice. (#20)
- **Spanish-speaking locals.** Flavor NPCs around the rings greet you in
  Spanish (¡Viva la Ciudad de la Noche!). (#21)
- **Chinese-speaking locals.** Some ring NPCs trash-talk fresh clones in
  Mandarin, with a dim English gloss. (#22)
- **Easter egg.** A laughing ROM construct hides in the Sprawlbelt — TALK to
  it for cynical hacker wisdom; it calls you "Boy". (#23)
- **Red-light strip.** Joytoys along the Sprawlbelt (the Rolling Rose, the
  Fortune Stall) — PAY for company: a fade-to-black hour that leaves you
  fully restored (HP + RAM). (#27)
- **Sparring gym.** The Iron Temple on the Sprawlbelt (north off the Stripped
  Lot) hosts non-lethal PvP: ATTACK another runner to spar — a downed sparrer is
  knocked out, keeps all their gear/scrip, and wakes at full HP. (#19)
- **Data terminals.** At any vendor/medic room or the Data Port: `SEND <runner>
  <msg>` to mail another player (delivered on their next login), `WIRE <runner>
  <scrip>` to transfer credits (even to offline runners), and `MAIL` to read. (#24)
- **Player trading.** `TRADE <runner>` opens a face-to-face swap: each side
  OFFERs items/scrip, and nothing moves until BOTH CONFIRM (any change re-opens
  both confirmations). Atomic, validated, CANCEL anytime. (#25)
- **Batch quest accept.** `ACCEPT 1 2 3` takes several bounties at once and
  `ACCEPT ALL` takes every eligible one; each pick is guarded independently. (#17)
- **Character points.** Each level banks spendable character points (shown on
  the score sheet when you have any); `SPEND <body|reflexes|intelligence>` raises
  a stat (Body also lifts MaxHP). Persisted across logout. (#11)
- **Stat-implant cyberware.** Each level band's vendors now stock three implants
  (Body / Reflexes / Intelligence), so every class can buy a grind aid and
  `INSTALL` it at a medic for a permanent stat boost. (#15)
- **Fresh-sleeve adrenaline.** Re-sleeving after a flatline now wakes you with a
  +15% HP buffer (overheal above max) that bleeds off as you take damage. (#26)

### Changed
- **Quest board shows your state.** On a fixer's board, a bounty you've already
  accepted is greyed out (with progress), and one that's complete and ready to
  turn in is flagged RED `[READY — turn in]`. (#10)
- **Story bounties are one-time.** Once you've claimed a story/street bounty you
  can't re-accept it; the low-level RP-ring rumors stay repeatable. Persisted. (#18)

### Fixed
- **MAP now always shows the way out.** The `▲ WAY OUT` arrow follows realm
  boundaries toward the surface (zone-1 Undercity → street, the Net → the Data
  Port → street), so even the first zone and Net rooms show how to get out, not
  just how to go deeper. (#9)
- **Machines leave wreckage, not a "corpse".** Drones, turrets, mechs and the
  like are now flagged mechanical — destroying/looting one reads as wreckage
  ("its frame sparks and goes dark", "you strip the wreck"), never a flatlined
  body. (#8)

## [2.1.0] - 2026-06-25

A quality-of-life release that works the backlog: a navigation map, fairer
quest turn-ins, and flavor fixes for bodiless foes.

### Added
- **`MAP` command (alias `M`).** A CP437 "you-are-here" panel: every exit
  labelled with where it leads (and whether it goes deeper/harder, back/easier,
  to a cache, shop, medic, or safe spot), plus the single move that takes you
  onward to the next harder area (`▼ PROCEED`) or back out (`▲ WAY OUT`). Works
  in the Undercity, the Net, and on the surface. (#6)

### Changed
- **Quests redeem at the giver too.** A completed bounty can now be `CLAIM`ed
  back with the quest-giver who offered it (the fixer's room, or — for roving
  RP-ring rumors — wherever it was scattered this session), not only at a
  broker. (#7)
- **Bodiless foes no longer leave a "corpse".** Loot crates (supply/data caches)
  now *crack open* into a cracked-open cache and Net constructs *shatter into
  shards*; only a real kill leaves a flatlined body. (#5)

### Docs
- Generated `docs/underground-descent.drawio` and `docs/net-ascent.drawio`
  (with PNGs) straight from the live zone data.

## [2.0.0] - 2026-06-25

A massive content release: both 1–99 progression paths are now hand-authored
worlds, with story quests, a roleplay zone, and a recall.

### Added
- **The underground descent (meatspace L1-99).** Replaced the placeholder band
  spine with ~143 authored rooms across 10 story arcs — the Neon Wasteland down to
  the Geo-Anchor Vault — with per-band foes, an arc boss each (Razorback Kane …
  the Loom Masterframe), varied room directions, and hidden ceiling/floor **loot
  caches** (`up`/`down`) that refill on cooldown.
- **The Net ascent (cyberspace L1-99).** Replaced the placeholder Net spine with
  ~150 authored rooms — 50 areas × 3 layers (Shell / Breach / Core) — so netspace
  moves in all **six directions** (`up`/`down` between layers, N/S/E/W between
  areas). 10 arcs from the Neon Underbelly up to the Living Library, with ICE
  foes, data-caches (RAM + scrip), and band-scaled deck/RAM vendors.
- **Story quest-givers & bounties** for every level section: the plot NPCs
  (Marcus, Cipher, Silas, Dr. Vance, Fixer-7, Mr. Lattice, Ravel …) hire you to
  take down their arc's boss. Giver-gated `ACCEPT`; `CLAIM` at any broker.
- **The RP transit rings** — a fast, RP-safe **Inner Circuit** and a longer
  **Sprawlbelt** loop (off Neon Alley, `north`) — plus 12 standalone, randomized,
  repeatable **rumor bounties** scattered across the ring NPCs fresh each session
  (no linear progression), each an easter egg to one of the PvE paths.
- **HOME recall.** `home` / `rest` now jacks a ~10-second teleport to your
  Re-Clone Bay from anywhere, broken by a hostile hit or by moving. Closes #4.

### Fixed
- **Netrunners can reach the Net.** The Data Port jacks straight into the first
  Net node, closing the long-standing "can't reach the Net" gap.

### Changed
- The world map changed, so a returning character whose saved location no longer
  exists simply wakes in their Re-Clone Bay — stats, gear, and progress are kept.
- The multi-stage Gauntlet ICE was preserved and re-homed in the first Net node's
  core; the old placeholder generators and safehouse tiers were retired.

[2.0.0]: https://github.com/CryptoJones/ChromeCircuitCowboys/releases/tag/v2.0.0

## [1.0.5] - 2026-06-25

### Changed
- **Per-area, level-scaled vendor stock.** Vendors no longer all sell the same
  global list: the Chrome Rose carries street stims + a starter blade, the Night
  Market carries better gear + cyberware, and the band safehouses (L30/50/70/90)
  stock progressively stronger weapons/decks/stims scaled to their depth. The
  master catalog gained those higher tiers (so looted gear still installs). Closes #1.

## [1.0.4] - 2026-06-25

### Added
- **Arrow-key movement** — Up/Down/Left/Right map to N/S/W/E when you're not
  mid-typing a command (great on a phone). Closes #2.

## [1.0.3] - 2026-06-24

### Fixed
- **Character creation no longer silently defaults on bad input.** A non-numeric
  or out-of-range class no longer becomes a Hacker, and letters at the
  skill-point prompts no longer count as 0 (which dumped your points into
  INTELLIGENCE). Both **re-prompt** now.

### Added
- **Type `Q` (or `QUIT`) at any creation prompt — or the handle prompt — to jack
  out** cleanly.

## [1.0.2] - 2026-06-24

### Added
- **Default runner name = your BBS handle.** The door advertises `caps=handle`
  in its handshake; the BBS pushes your handle back, so the name prompt shows
  `Handle [YourHandle] (Enter to use):` — just hit Enter to use it.
- **Numbered vendor list + quantity buys.** `LIST` numbers each ware, and
  `BUY <#|name> [qty]` lets you buy several at once (`buy 3 4` = four of item 3;
  default 1). Weapons/decks stay one-time upgrades.

### Docs
- `docs/world-map.drawio` + `world-map.png` — Cyberdeck-dark level-area map of
  the dual meatspace/Net paths (L1–99) for theming reference.

## [1.0.1] - 2026-06-24

### Added
- **Resident-door version handshake** — advertises the game version to the host
  as the first bytes on connect (OSC `ESC ] ABBS;version=<v> BEL`), so the BBS
  shows it on the launch line (ABBS Door Spec §2.2).
- **Dual-path world skeleton (L11–99)** — placeholder meatspace + Net progression
  bands hang off Back Alley and the Deep Net; level cap raised **50 → 99**. Rooms
  are TODO placeholders (no monsters yet) pending per-band theming.

### Changed
- **Randomized kill rewards** — XP and scrip now roll ~75–125% of each mob's base.
- **Mobs drop a lootable body; respawn is gated until it's looted** — no spawning
  over an unlooted corpse. ICE constructs shatter into **broken shards** you
  *salvage* (not a body), with the same loot-gated regeneration.

[1.0.1]: https://github.com/CryptoJones/ChromeCircuitCowboys/releases/tag/v1.0.1

## [1.0.0] - 2026-06-24

First standalone release — carved out of AdmiralBBS into its own repo (the door
was bundled with the BBS through AdmiralBBS v1.x; from here it ships and versions
on its own). Generic cyberpunk, tied to no single franchise.

### Features
- Persistent multiplayer world played as a resident door (single world goroutine;
  lock-free, deterministic engine).
- Classes: Hacker / Enforcer / Operator / Mechanic.
- Meatspace + Net combat; netrun programs (scalpel/hammer/leech/mirror/medic),
  RAM economy, multi-stage morphing ICE.
- Re-clone death model: full-HP clone at a private, spawn-safe Re-Clone Bay;
  10% clone fee, no XP loss; your old body drops as a lootable corpse with your
  gear + cyberware; cyberware re-installs at an Emergency Medic; `give` to return
  recovered gear.
- Open PvP everywhere except the safe street zones (a security drone flatlines
  aggressors there).
- Consent-based crews (invite/accept, leader-only, succession on death), shared
  XP, crew radio.
- Bounties, leaderboard, vendors, per-runner stash + level-scaled carry cap,
  and RP emotes (`me` / `emote` / `:`).
- Forge-agnostic update check (`-update-url` / `CCC_UPDATE_URL`).

[1.0.0]: https://github.com/CryptoJones/ChromeCircuitCowboys/releases/tag/v1.0.0
