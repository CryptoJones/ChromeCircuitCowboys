# Changelog

All notable changes to **Chrome Circuit Cowboys** are documented here.

Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
versioning: [SemVer](https://semver.org/spec/v2.0.0.html).

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
