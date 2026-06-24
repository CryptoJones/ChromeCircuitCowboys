<p align="center"><em>Proudly Made in Nebraska. Go Big Red! 🌽 <a href="https://xkcd.com/2347/">https://xkcd.com/2347/</a></em></p>

# Chrome Circuit Cowboys (C³)

A multiplayer, generic-cyberpunk MUD **door game** — jack into the Net, level up,
breach the ICE, run the streets, and duel other cowboys. One persistent world,
many simultaneous players, played over a BBS as a **resident door**.

It speaks the resident-door bridge protocol, so it drops into any compatible
host. Originally bundled with **AdmiralBBS**; now its own project.

## Runs as a door for

- **AdmiralBBS** — the clean-room ANSI BBS it was built for:
  <https://github.com/CryptoJones/AdmiralBBS>
- Any host implementing the **ABBS Door Specification**:
  <https://github.com/CryptoJones/ABBS-Door-Specification>

## What it is

- **Classes:** Hacker / Enforcer / Operator / Mechanic (no two builds play the same).
- **Two arenas:** meatspace combat (Body) and the Net, where ATTACK *breaches* ICE
  with Intelligence and spends RAM.
- **Death = re-clone:** your stack is backed up, so you wake in a fresh full-HP
  clone at your private Re-Clone Bay (spawn-safe). Your old body drops where you
  fell as a lootable corpse holding your gear + cyberware — recover it (or a
  crewmate does) and re-install cyberware at an Emergency Medic.
- **Open PvP** everywhere except the safe street outside the clone bays, where a
  security drone flatlines anyone who draws first.
- **Crews** (consent-based invites + shared XP), **bounties**, a **leaderboard**,
  netrun **programs**, a per-runner **stash** + level-scaled carry cap, and RP
  **emotes**.

## Build & run

```sh
go build -o ccc-server .
./ccc-server -addr 127.0.0.1:4000 -db cowboy.db -tick 2s
```

The BBS bridges callers to the listen address. Flags:

| Flag | Default | Purpose |
|------|---------|---------|
| `-addr` | `127.0.0.1:4000` | TCP listen address for the BBS bridge |
| `-db` | `cowboy.db` | SQLite character database path |
| `-tick` | `2s` | World/combat tick interval |
| `-update-url` | `$CCC_UPDATE_URL` | Forge `releases/latest` JSON endpoint to check for updates (GitHub/Codeberg/Forgejo shape). Empty = no check — **no forge is hardcoded.** |
| `-version` | | Print version and exit |

## Update checks

Forge-agnostic: point `-update-url` (or `CCC_UPDATE_URL`) at your forge's
`.../repos/<owner>/<repo>/releases/latest` endpoint. On startup the server
compares the running version to the latest release tag and logs a notice if a
newer one exists — works the same on GitHub, Codeberg, or any Forgejo. If unset,
no check runs.
