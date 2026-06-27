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

- **Classes:** Hacker / Enforcer / Operator / Mechanic (no two builds play the same),
  with **character points** to raise Body / Reflexes / Intelligence.
- **Two arenas:** meatspace combat (Body) and the Net, where ATTACK *breaches* ICE
  with Intelligence and spends RAM.
- **A 1–99 authored world:** the Undercity descent and the Net ascent, each with its
  own zone lore, plus the Noche City surface and its transit rings.
- **Death = re-clone:** your stack is backed up, so you wake in a fresh full-HP
  clone at your private Re-Clone Bay (spawn-safe). Your old body drops where you
  fell as a lootable corpse holding your gear + cyberware — recover it (or a
  crewmate does) and re-install cyberware at an Emergency Medic.
- **Open PvP** everywhere except the safe street outside the clone bays, where a
  security drone flatlines anyone who draws first. Non-lethal **sparring gym** for
  solo and crew-vs-crew duels.
- **Crews:** consent-based invites, shared XP, party-follow, party loot for every
  class present, a crew-chat shortcut, and **clans** with reward bonuses.
- **AI runners:** the streets are populated by autonomous "runners" that wander and
  banter. `GROUP` one and it joins your crew, follows you everywhere, and fights any
  mob that attacks the crew — kills, XP, and loot stay yours.
- **Living world:** TALK to locals for lore or just to chat, a CP437 **MAP** with the
  way deeper or out, a **HACK** terminal mini-game, data terminals to **SEND** mail
  and **WIRE** scrip, face-to-face **TRADE**, vendors (BUY/**SELL**), bounties, a
  **leaderboard**, netrun **programs**, a per-runner **stash** + level-scaled carry
  cap, RP **emotes**, and switchable **colorblind-friendly** color themes.
- **Accounts:** per-character **password auth** (set at creation; legacy characters
  are prompted to set one on first login).

## Commands (highlights)

- **Move:** `N S E W U D` · `IN`/`OUT` of your pod · `MAP`/`M`
- **Look:** `LOOK`/`L` (room) · `LOOK <item|#>` (examine) · `WHO` · `SCORE`
- **Fight:** `ATTACK`/`A` · `FLEE` · `LOOT`/`LO` · `OPEN` (caches)
- **Items:** `INVENTORY`/`I` (numbered) — a number quick-USEs it, or use it as a
  target: `SELL 2`, `LOOK 3`, `GIVE 1 <runner>`, `DROP 1`, `USE 2`
- **Crew:** `GROUP <runner>` (invite / recruit a bot) · `ACCEPT` · `LEAVE` · `;` crew chat
- **World:** `TALK [words]` · `HACK` · `QUESTS`/`ACCEPT`/`CLAIM` · `SEND`/`WIRE`/`TRADE` · `HOME`

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
| `-bots` | `32` | Number of AI runners that populate the world (`0` = none; capped at the roster size) |
| `-update-url` | `$CCC_UPDATE_URL` | Forge `releases/latest` JSON endpoint to check for updates (GitHub/Codeberg/Forgejo shape). Empty = no check — **no forge is hardcoded.** |
| `-version` | | Print version and exit |

## Recent releases

- **v2.5.8** — `MAP` no longer goes quiet in the deepest band (L91-99): with no harder area left, it now points you to the zone's **FINAL OBJECTIVE** instead of only showing the way back.
- **v2.5.7** — **Randomized loot**: mobs now drop varied gear from their zone's shop stock (no more endless stimpaks); AI crewmates no longer inflate party loot. **`GSAY`** crew calls now get answered by your bot crew.
- **v2.5.6** — Quest **level requirements are waived while you're in a crew** — run higher-tier bounties with a party backing you up.
- **v2.5.5** — `QUESTS` now shows a **direction to head** for each active bounty (toward the target while hunting, toward the giver/broker once it's READY to claim).
- **v2.5.4** — Your crew piles onto a mob the instant **you** attack; crewmates' hits (and hits on them) now show in your **combat log**; `GROUP ALL` recruits every free runner in the room; `INSTALL <item|#> [qty]` installs by inventory number and in bulk (e.g. `INSTALL 1 8`).
- **v2.5.3** — Crewed AI runners follow you and **fight** any mob attacking the crew; kills, XP, quest credit, and loot stay with you (bots soak hits but never die or steal aggro).
- **v2.5.2** — **Recruit** AI runners into your crew with `GROUP <runner>` — they auto-join, warp to you, and tag along silently.
- **v2.5.1** — **32** AI runners by default.
- **v2.5.0** — `TALK <words>` gets a reply routed by tone; use inventory **numbers as command targets** (`SELL 2`, `LOOK 3`, `GIVE 1`, …); **AI runners** populate the streets so solo players aren't alone.
- **v2.4.0** — "Violet Lotus" content pass: vivid room descriptions, expanded NPC lore, per-item flavor on `LOOK`.
- **v2.3.1** — Playtest fixes: masked password entry, room re-shown after setup, distinct colorblind themes (map included).
- **v2.3.0** — Password auth, **HACK** terminal mini-game, **SELL** to vendors, **clans** + reward bonuses, party-vs-party **sparring**, group combat + party-follow, **DROP/GET** floor loot, `LOOK <item>` examine, colorblind themes, scattered data terminals, character-point spend, unsafe-logout penalty, and more.
- **v2.2.0** — Authored worlds + story quests, **MAP** command, TALK lore, numbered inventory quick-use, **OPEN** caches, multi-`ACCEPT` quests, the fitness/spar gym, data terminals (**SEND** mail / **WIRE** scrip), player **TRADE**, re-sleeve buff, joytoy buffs, multilingual NPCs, and a quiet Dixie Flatline homage.

## Update checks

Forge-agnostic: point `-update-url` (or `CCC_UPDATE_URL`) at your forge's
`.../repos/<owner>/<repo>/releases/latest` endpoint. On startup the server
compares the running version to the latest release tag and logs a notice if a
newer one exists — works the same on GitHub, Codeberg, or any Forgejo. If unset,
no check runs.
