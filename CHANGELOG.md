# Changelog

All notable changes to **Chrome Circuit Cowboys** are documented here.

Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
versioning: [SemVer](https://semver.org/spec/v2.0.0.html).

## [2.2.0] - 2026-06-26

Backlog batch: QoL, content, and combat-flow improvements. (Items land
incrementally; see the issue links.)

### Added
- **`OPEN` command for caches.** Crack open a supply/data cache with `OPEN`
  instead of attacking it — the intuitive verb for an inert container. (#16)
- **Combat shortcuts.** `A` = attack, `LO` = loot (`L` still = look). (#14)
- **Batch quest accept.** `ACCEPT 1 2 3` takes several bounties at once and
  `ACCEPT ALL` takes every eligible one; each pick is guarded independently. (#17)
- **Character points.** Each level banks spendable character points (shown on
  the score sheet when you have any); `SPEND <body|reflexes|intelligence>` raises
  a stat (Body also lifts MaxHP). Persisted across logout. (#11)

### Changed

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
