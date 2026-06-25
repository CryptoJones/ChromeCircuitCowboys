# Changelog

All notable changes to **Chrome Circuit Cowboys** are documented here.

Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
versioning: [SemVer](https://semver.org/spec/v2.0.0.html).

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
