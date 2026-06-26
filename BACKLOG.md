# Backlog

Ideas and tasks not yet scheduled.

- [ ] Secret (HELP-hidden) command to display the current room's ID ([#28](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/28))
- [ ] Rename the `back_alley` room's display name "Back Alley" → "Underground Entrance" ([#29](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/29))
- [ ] Character sheet: spell out "AC" as "Armor Class" ([#30](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/30))
- [ ] Character sheet: always show character points (0 if none) + bold "You have character points to spend." note ([#31](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/31))
- [ ] Some quests award a character point, class-flavored ([#32](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/32))
- [ ] Hacking mini-game at data terminals everywhere in meatspace ([#33](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/33))
- [ ] Scatter data terminals across ~1-in-4 surface rooms, randomized (not a predictable pattern) ([#34](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/34))
- [ ] Rewrite every room description (more vivid) via Violet Lotus on ronin28, with per-area plot context from the Obsidian Vault ([#35](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/35)) — _do via plan mode_
- [ ] `TALK` takes input (e.g. `talk "hi"`) with random responses — modular for future growth ([#36](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/36)) — _do via plan mode_
- [ ] AI characters that act as players to populate empty servers (so solo players aren't lonely) ([#37](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/37)) — _do via plan mode_
- [ ] Switchable colorblind-friendly palette (Claude Code's), light + dark — accessibility ([#38](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/38))
- [ ] Expand all NPCs' default responses via Violet Lotus (they're too terse) ([#39](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/39)) — _do via plan mode_
- [ ] Sell unwanted items for scrip at vendors (`SELL`) ([#40](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/40))
- [ ] Multiple players can attack the same mob at once (shared combat + reward model) ([#41](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/41))
- [ ] Party follow: members move with the leader ([#42](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/42))
- [ ] Quick party-chat shortcut (single-char prefix) that propagates to all crew ([#43](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/43))
- [ ] Clan system + party/clan reward bonuses (1.8x for clanmates, on top of a party bonus) ([#44](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/44))
- [ ] Party-vs-party sparring in the gym (team duels, non-lethal) ([#45](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/45))
- [ ] Container text: stop saying "You lunge at the container" — randomized container verbs (pry/crack/jimmy/…) ([#46](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/46))
- [ ] Lighten the dark-grey (dim) text in the default color scheme ([#47](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/47))
- [ ] Unsafe logout penalty: lose 5% HP + "X did Y while you were sleeping for Z damage" on return ([#48](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/48))
- [ ] Block BUYING unusable items (but allow looting + selling them) + always explain why USE fails ([#49](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/49))
- [ ] Party loot drops include something for each class in the party ([#50](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/50))

- [x] Make any key press continue from the MOTD ([#3](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/3)) — _AdmiralBBS `ShowMOTD` reads any key (or Enter); shipped in AdmiralBBS v2.0.7, live on pluto_
- [x] Drones have no body — reword the loot text when you loot them ([#8](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/8)) — _shipped v2.2.0: machines leave wreckage, never a corpse_
- [x] MAP should show the way back to an easier area / the exit, not just deeper ([#9](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/9)) — _shipped v2.2.0: WAY OUT now resolves to the surface across realms_
- [x] Quest board: grey out already-accepted bounties, show RED when ready to turn in ([#10](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/10)) — _shipped v2.2.0: board greys accepted, flags READY in red_
- [x] Character points: show available points on the character sheet + a "spend character points" option when you have any ([#11](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/11)) — _shipped v2.2.0: character points on the sheet + SPEND_
- [x] TALK to NPCs — every NPC gives backstory/lore about the level you're on ([#12](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/12)) — _shipped v2.2.0: TALK gives per-zone backstory_
- [x] Numbered inventory — press the number to USE that item, no Enter needed (fast for combat) ([#13](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/13)) — _shipped v2.2.0: numbered inventory, digit = quick-use_
- [x] Shortcuts: `A` = attack, `LO` = loot (`L` stays look) ([#14](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/14)) — _shipped v2.2.0: A=attack, LO=loot, L stays look_
- [x] Per-area cyberware that grants character points — 2-3 per band so each class can buy a grind aid ([#15](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/15)) — _shipped v2.2.0: 3 stat implants per band, install at a medic_
- [x] `OPEN` command for supply/data caches instead of attacking them ([#16](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/16)) — _shipped v2.2.0: OPEN cracks open caches_
- [x] `ACCEPT` multiple quests — `accept 1 2 3 4` and `accept all` ([#17](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/17)) — _shipped v2.2.0: accept 1 2 3 / accept all_
- [x] Completed quests can't be re-accepted — except the low-level RP ring rumors ([#18](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/18)) — _shipped v2.2.0: completed story bounties are one-time; rings exempt_
- [x] Fitness center on the outer ring (Sprawlbelt) — non-lethal PvP sparring, no death or loot/cyberware loss ([#19](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/19)) — _shipped v2.2.0: Iron Temple sparring gym, non-lethal_
- [x] New-player NPC in the clone booth — TALK for a brief intro to commands, actions, etc. ([#20](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/20)) — _shipped v2.2.0: TALK booth onboarding primer_
- [x] Spanish-speaking NPCs — greet in Spanish / "¡Viva la Ciudad de la Noche!" ([#21](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/21)) — _shipped v2.2.0: Spanish ring NPCs_
- [x] Chinese-speaking NPCs that trash-talk you in Chinese ([#22](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/22)) — _shipped v2.2.0: Chinese trash-talk NPCs w/ gloss_
- [x] Easter egg: a Dixie Flatline homage somewhere in the rings — calls you "Boy" (oblique, IP-scrub-safe) ([#23](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/23)) — _shipped v2.2.0: laughing ROM construct, calls you Boy_
- [x] Data terminals across meatspace — send messages + credits to other players ([#24](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/24)) — _shipped v2.2.0: SEND mail / WIRE scrip at terminals_
- [x] Player-to-player trading — two-sided confirm, items + scrip, atomic swap ([#25](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/25)) — _shipped v2.2.0: two-sided confirm-locked TRADE_
- [x] New-sleeve buff: +15% health when you get a new sleeve (re-sleeve) ([#26](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/26)) — _shipped v2.2.0: +15% HP buffer on re-sleeve_
- [x] Red-light district joytoy NPCs — pay scrip, fade-to-black, grants a buff ([#27](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/27)) — _shipped v2.2.0: PAY a joytoy for a fade-to-black restore_
- [x] Loot crates & net mobs shouldn't leave a "corpse" — they have no bodies; reword the death/loot text for them ([#5](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/5)) — _shipped in v2.1.0: caches crack open, ICE shatters into shards_
- [x] 16-bit map feature — render a CP437/ANSI map showing how to exit the "level" or proceed to the next harder area ([#6](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/6)) — _`MAP`/`M`: labelled exits + PROCEED/WAY-OUT arrows, works in Undercity, Net, and surface_
- [x] Quests redeemable at the giver, not just at brokers ([#7](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/7)) — _claim now pays out at the fixer/ring-giver room or a broker_
- [x] Teleport-home has a 10-second delay, interrupted if attacked by a mob ([#4](https://github.com/CryptoJones/ChromeCircuitCowboys/issues/4)) — _shipped in v2.0.0 (HOME recall)_

---

*Proudly Made in Nebraska. Go Big Red! 🌽 <https://xkcd.com/2347/>*
