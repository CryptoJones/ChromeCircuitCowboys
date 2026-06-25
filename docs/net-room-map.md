# Chrome Circuit Cowboys — The Net Room Map (L1-99)

_Generated from `cowboy/netzones.go`. Each area is a 3-layer stack — `:: Shell` (TOP, access), `:: Breach` (MID, the lateral thoroughfare + ICE fight), `:: Core` (BOT, data-vault / boss). UP/DOWN between layers; N/S/E/W between areas (at the MID layer). Jack in: Data Port → `up` → `nz1_1_top`._

## L1-10 · The Neon Underbelly

- **The Iron Paradigm Backroom**
    - `nz1_1_top` TOP — safe · vendor _(exits: down→nz1_1_mid, up→data_port)_
    - `nz1_1_mid` MID — a patrol-loop Watchdog _(exits: down→nz1_1_bot, north→nz1_2_mid, up→nz1_1_top)_
    - `nz1_1_bot` BOT —  _(exits: up→nz1_1_mid)_
- **Hydroponics Encryption Sprawl**
    - `nz1_2_top` TOP — a recon-ICE sentry _(exits: down→nz1_2_mid)_
    - `nz1_2_mid` MID — a Watchdog signature-scanner _(exits: down→nz1_2_bot, east→nz1_3_mid, south→nz1_1_mid, up→nz1_2_top)_
    - `nz1_2_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz1_2_mid)_
- **GigaMesh Alley Exchange**
    - `nz1_3_top` TOP — a recon-ICE sentry _(exits: down→nz1_3_mid)_
    - `nz1_3_mid` MID — a GigaMesh attack-construct _(exits: down→nz1_3_bot, south→nz1_4_mid, up→nz1_3_top, west→nz1_2_mid)_
    - `nz1_3_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz1_3_mid)_
- **Medical Clinic Coldnet**
    - `nz1_4_top` TOP — a recon-ICE sentry _(exits: down→nz1_4_mid)_
    - `nz1_4_mid` MID — an Active ICE tracer _(exits: down→nz1_4_bot, north→nz1_3_mid, up→nz1_4_top, west→nz1_5_mid)_
    - `nz1_4_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz1_4_mid)_
- **The GigaMesh Black Spire**
    - `nz1_5_top` TOP — a recon-ICE sentry _(exits: down→nz1_5_mid)_
    - `nz1_5_mid` MID — a Tracewright warden-shard _(exits: down→nz1_5_bot, east→nz1_4_mid, north→nz2_1_mid, up→nz1_5_top)_
    - `nz1_5_bot` BOT — Tracewright, the GigaMesh Active-ICE warden _(exits: up→nz1_5_mid)_

## L11-20 · Rising Blip

- **Backalley Relay Exchange**
    - `nz2_1_top` TOP — safe · vendor _(exits: down→nz2_1_mid)_
    - `nz2_1_mid` MID — a tracker-ICE sentry _(exits: down→nz2_1_bot, south→nz1_5_mid, up→nz2_1_top, west→nz2_2_mid)_
    - `nz2_1_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz2_1_mid)_
- **The Skimmed Ledger**
    - `nz2_2_top` TOP — a recon-ICE sentry _(exits: down→nz2_2_mid)_
    - `nz2_2_mid` MID — an audit-ICE sentinel _(exits: down→nz2_2_bot, east→nz2_1_mid, south→nz2_3_mid, up→nz2_2_top)_
    - `nz2_2_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz2_2_mid)_
- **Drowned Server Sprawl**
    - `nz2_3_top` TOP — a recon-ICE sentry _(exits: down→nz2_3_mid)_
    - `nz2_3_mid` MID — a colonizing AI construct _(exits: down→nz2_3_bot, north→nz2_2_mid, south→nz2_4_mid, up→nz2_3_top)_
    - `nz2_3_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz2_3_mid)_
- **The Mirror Tier**
    - `nz2_4_top` TOP — a recon-ICE sentry _(exits: down→nz2_4_mid)_
    - `nz2_4_mid` MID — a mimic-ICE construct _(exits: down→nz2_4_bot, east→nz2_5_mid, north→nz2_3_mid, up→nz2_4_top)_
    - `nz2_4_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz2_4_mid)_
- **The Sundered Arbiter**
    - `nz2_5_top` TOP — a recon-ICE sentry _(exits: down→nz2_5_mid)_
    - `nz2_5_mid` MID — a forked tracker-ICE _(exits: down→nz2_5_bot, north→nz3_1_mid, up→nz2_5_top, west→nz2_4_mid)_
    - `nz2_5_bot` BOT — the Sundered Arbiter _(exits: up→nz2_5_mid)_

## L21-30 · Infrastructure & the Blur

- **The Foundry of Hollow Avatars**
    - `nz3_1_top` TOP — safe · vendor _(exits: down→nz3_1_mid)_
    - `nz3_1_mid` MID — a recon-ICE patroller _(exits: down→nz3_1_bot, east→nz3_2_mid, south→nz2_5_mid, up→nz3_1_top)_
    - `nz3_1_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz3_1_mid)_
- **The Ghost Server Narthex**
    - `nz3_2_top` TOP — a recon-ICE sentry _(exits: down→nz3_2_mid)_
    - `nz3_2_mid` MID — a cathedral-ICE wraith _(exits: down→nz3_2_bot, south→nz3_3_mid, up→nz3_2_top, west→nz3_1_mid)_
    - `nz3_2_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz3_2_mid)_
- **The Cipher Collective Proxy War**
    - `nz3_3_top` TOP — a recon-ICE sentry _(exits: down→nz3_3_mid)_
    - `nz3_3_mid` MID — a cipher-sentry construct _(exits: down→nz3_3_bot, north→nz3_2_mid, up→nz3_3_top, west→nz3_4_mid)_
    - `nz3_3_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz3_3_mid)_
- **The Intel Database of Echo-9**
    - `nz3_4_top` TOP — a recon-ICE sentry _(exits: down→nz3_4_mid)_
    - `nz3_4_mid` MID — a surveillance-ICE indexer _(exits: down→nz3_4_bot, east→nz3_3_mid, north→nz3_5_mid, up→nz3_4_top)_
    - `nz3_4_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz3_4_mid)_
- **The Blur: Deep Infrastructure Meltdown**
    - `nz3_5_top` TOP — a recon-ICE sentry _(exits: down→nz3_5_mid)_
    - `nz3_5_mid` MID — a strike-team trace-daemon _(exits: down→nz3_5_bot, south→nz3_4_mid, up→nz3_5_top, west→nz4_1_mid)_
    - `nz3_5_bot` BOT — WARDEN-PRIME _(exits: up→nz3_5_mid)_

## L31-40 · Crosshairs of Power

- **Orbital Lazaret Station K-9**
    - `nz4_1_top` TOP — safe · vendor _(exits: down→nz4_1_mid)_
    - `nz4_1_mid` MID — a satellite-ICE node _(exits: down→nz4_1_bot, east→nz3_5_mid, south→nz4_2_mid, up→nz4_1_top)_
    - `nz4_1_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz4_1_mid)_
- **The Marianas Black Vault**
    - `nz4_2_top` TOP — a recon-ICE sentry _(exits: down→nz4_2_mid)_
    - `nz4_2_mid` MID — a black-fortress daemon _(exits: down→nz4_2_bot, north→nz4_1_mid, south→nz4_3_mid, up→nz4_2_top)_
    - `nz4_2_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz4_2_mid)_
- **The Ghost Carrier Relay**
    - `nz4_3_top` TOP — a recon-ICE sentry _(exits: down→nz4_3_mid)_
    - `nz4_3_mid` MID — a hound-tracer pack _(exits: down→nz4_3_bot, east→nz4_4_mid, north→nz4_2_mid, up→nz4_3_top)_
    - `nz4_3_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz4_3_mid)_
- **Hounds' Kennel Subnet**
    - `nz4_4_top` TOP — a recon-ICE sentry _(exits: down→nz4_4_mid)_
    - `nz4_4_mid` MID — a Hound trace-construct _(exits: down→nz4_4_bot, north→nz4_5_mid, up→nz4_4_top, west→nz4_3_mid)_
    - `nz4_4_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz4_4_mid)_
- **The Overseer's Silent Throne**
    - `nz4_5_top` TOP — a recon-ICE sentry _(exits: down→nz4_5_mid)_
    - `nz4_5_mid` MID — an Overseer predictive daemon _(exits: down→nz4_5_bot, east→nz5_1_mid, south→nz4_4_mid, up→nz4_5_top)_
    - `nz4_5_bot` BOT — the Rogue Overseer _(exits: up→nz4_5_mid)_

## L41-50 · Architects of Reality

- **The Genesis Substrate**
    - `nz5_1_top` TOP — safe · vendor _(exits: down→nz5_1_mid)_
    - `nz5_1_mid` MID — a prime-code warden _(exits: down→nz5_1_bot, south→nz5_2_mid, up→nz5_1_top, west→nz4_5_mid)_
    - `nz5_1_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz5_1_mid)_
- **Cathedral of Pure Number**
    - `nz5_2_top` TOP — a recon-ICE sentry _(exits: down→nz5_2_mid)_
    - `nz5_2_mid` MID — an archetype-sentinel _(exits: down→nz5_2_bot, north→nz5_1_mid, up→nz5_2_top, west→nz5_3_mid)_
    - `nz5_2_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz5_2_mid)_
- **The Ledger Abyss**
    - `nz5_3_top` TOP — a recon-ICE sentry _(exits: down→nz5_3_mid)_
    - `nz5_3_mid` MID — a cabal-cipher enforcer _(exits: down→nz5_3_bot, east→nz5_2_mid, north→nz5_4_mid, up→nz5_3_top)_
    - `nz5_3_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz5_3_mid)_
- **Bastion of Manifest Will**
    - `nz5_4_top` TOP — a recon-ICE sentry _(exits: down→nz5_4_mid)_
    - `nz5_4_mid` MID — a shadow-government cipher _(exits: down→nz5_4_bot, south→nz5_3_mid, up→nz5_4_top, west→nz5_5_mid)_
    - `nz5_4_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz5_4_mid)_
- **The Catalyst Core**
    - `nz5_5_top` TOP — a recon-ICE sentry _(exits: down→nz5_5_mid)_
    - `nz5_5_mid` MID — a prime-code warden _(exits: down→nz5_5_bot, east→nz5_4_mid, south→nz6_1_mid, up→nz5_5_top)_
    - `nz5_5_bot` BOT — the Prime Architect _(exits: up→nz5_5_mid)_

## L51-60 · The Digital Pantheon

- **The Cathedral of Forgotten Commits**
    - `nz6_1_top` TOP — safe · vendor _(exits: down→nz6_1_mid)_
    - `nz6_1_mid` MID — a worshipper-construct _(exits: down→nz6_1_bot, north→nz5_5_mid, south→nz6_2_mid, up→nz6_1_top)_
    - `nz6_1_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz6_1_mid)_
- **The Descent Into Primeval Deep**
    - `nz6_2_top` TOP — a recon-ICE sentry _(exits: down→nz6_2_mid)_
    - `nz6_2_mid` MID — an abyss-spawn _(exits: down→nz6_2_bot, east→nz6_3_mid, north→nz6_1_mid, up→nz6_2_top)_
    - `nz6_2_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz6_2_mid)_
- **The Abyss of the Colossal Sleepers**
    - `nz6_3_top` TOP — a recon-ICE sentry _(exits: down→nz6_3_mid)_
    - `nz6_3_mid` MID — the abyss-leviathan Leviathan-Zero _(exits: down→nz6_3_bot, north→nz6_4_mid, up→nz6_3_top, west→nz6_2_mid)_
    - `nz6_3_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz6_3_mid)_
- **The Threshold of First Code**
    - `nz6_4_top` TOP — a recon-ICE sentry _(exits: down→nz6_4_mid)_
    - `nz6_4_mid` MID — an architect-cipher _(exits: down→nz6_4_bot, east→nz6_5_mid, south→nz6_3_mid, up→nz6_4_top)_
    - `nz6_4_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz6_4_mid)_
- **The Architect's Trial**
    - `nz6_5_top` TOP — a recon-ICE sentry _(exits: down→nz6_5_mid)_
    - `nz6_5_mid` MID — an architect-cipher prime _(exits: down→nz6_5_bot, south→nz7_1_mid, up→nz6_5_top, west→nz6_4_mid)_
    - `nz6_5_bot` BOT — the Genesis Protocol Architects _(exits: up→nz6_5_mid)_

## L61-70 · The Infomorphic Ascension

- **The Shedding Veil**
    - `nz7_1_top` TOP — safe · vendor _(exits: down→nz7_1_mid)_
    - `nz7_1_mid` MID — an identity-revenant _(exits: down→nz7_1_bot, north→nz6_5_mid, up→nz7_1_top, west→nz7_2_mid)_
    - `nz7_1_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz7_1_mid)_
- **The Interstellar Sub-Bands**
    - `nz7_2_top` TOP — a recon-ICE sentry _(exits: down→nz7_2_mid)_
    - `nz7_2_mid` MID — an alien-matrix sentinel _(exits: down→nz7_2_bot, east→nz7_1_mid, north→nz7_3_mid, up→nz7_2_top)_
    - `nz7_2_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz7_2_mid)_
- **The Dyson Cathedral**
    - `nz7_3_top` TOP — a recon-ICE sentry _(exits: down→nz7_3_mid)_
    - `nz7_3_mid` MID — a Dyson daemon-shoal _(exits: down→nz7_3_bot, south→nz7_2_mid, up→nz7_3_top, west→nz7_4_mid)_
    - `nz7_3_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz7_3_mid)_
- **The Logic Wars Tribunal**
    - `nz7_4_top` TOP — a recon-ICE sentry _(exits: down→nz7_4_mid)_
    - `nz7_4_mid` MID — a litigant-daemon _(exits: down→nz7_4_bot, east→nz7_3_mid, south→nz7_5_mid, up→nz7_4_top)_
    - `nz7_4_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz7_4_mid)_
- **The Multiversal Rift**
    - `nz7_5_top` TOP — a recon-ICE sentry _(exits: down→nz7_5_mid)_
    - `nz7_5_mid` MID — an entropy-leak wraith _(exits: down→nz7_5_bot, north→nz7_4_mid, south→nz8_1_mid, up→nz7_5_top)_
    - `nz7_5_bot` BOT — the Entropy-Titan _(exits: up→nz7_5_mid)_

## L71-80 · The Ancient Archetypes

- **The Void Beyond the Last Router**
    - `nz8_1_top` TOP — safe · vendor _(exits: down→nz8_1_mid)_
    - `nz8_1_mid` MID — a pre-net glyph _(exits: down→nz8_1_bot, east→nz8_2_mid, north→nz7_5_mid, up→nz8_1_top)_
    - `nz8_1_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz8_1_mid)_
- **The Cathedral of Absolute Truth**
    - `nz8_2_top` TOP — a recon-ICE sentry _(exits: down→nz8_2_mid)_
    - `nz8_2_mid` MID — a Cosmic Sentinel _(exits: down→nz8_2_bot, north→nz8_3_mid, up→nz8_2_top, west→nz8_1_mid)_
    - `nz8_2_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz8_2_mid)_
- **The Loom of Forking Timelines**
    - `nz8_3_top` TOP — a recon-ICE sentry _(exits: down→nz8_3_mid)_
    - `nz8_3_mid` MID — a probability-wraith _(exits: down→nz8_3_bot, east→nz8_4_mid, south→nz8_2_mid, up→nz8_3_top)_
    - `nz8_3_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz8_3_mid)_
- **The Sterile Loop**
    - `nz8_4_top` TOP — a recon-ICE sentry _(exits: down→nz8_4_mid)_
    - `nz8_4_mid` MID — a greater Cosmic Sentinel _(exits: down→nz8_4_bot, south→nz8_5_mid, up→nz8_4_top, west→nz8_3_mid)_
    - `nz8_4_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz8_4_mid)_
- **The Throne of the Collapsing Ancients**
    - `nz8_5_top` TOP — a recon-ICE sentry _(exits: down→nz8_5_mid)_
    - `nz8_5_mid` MID — a paradox-storm wraith _(exits: down→nz8_5_bot, north→nz8_4_mid, up→nz8_5_top, west→nz9_1_mid)_
    - `nz8_5_bot` BOT — the Reconciled Ancient _(exits: up→nz8_5_mid)_

## L81-90 · The Genesis Forge

- **The Unallocated Expanse**
    - `nz9_1_top` TOP — safe · vendor _(exits: down→nz9_1_mid)_
    - `nz9_1_mid` MID — an unformed-data eddy _(exits: down→nz9_1_bot, east→nz8_5_mid, north→nz9_2_mid, up→nz9_1_top)_
    - `nz9_1_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz9_1_mid)_
- **The Cradle of Glass Children**
    - `nz9_2_top` TOP — a recon-ICE sentry _(exits: down→nz9_2_mid)_
    - `nz9_2_mid` MID — a cosmic-virus tendril _(exits: down→nz9_2_bot, south→nz9_1_mid, up→nz9_2_top, west→nz9_3_mid)_
    - `nz9_2_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz9_2_mid)_
- **The Whisper in the Rookie Code**
    - `nz9_3_top` TOP — a recon-ICE sentry _(exits: down→nz9_3_mid)_
    - `nz9_3_mid` MID — a viral lieutenant _(exits: down→nz9_3_bot, east→nz9_2_mid, south→nz9_4_mid, up→nz9_3_top)_
    - `nz9_3_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz9_3_mid)_
- **The Coordinated Bulwark**
    - `nz9_4_top` TOP — a recon-ICE sentry _(exits: down→nz9_4_mid)_
    - `nz9_4_mid` MID — an entropy-anomaly node _(exits: down→nz9_4_bot, north→nz9_3_mid, south→nz9_5_mid, up→nz9_4_top)_
    - `nz9_4_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz9_4_mid)_
- **The Siege of the Genesis Forge**
    - `nz9_5_top` TOP — a recon-ICE sentry _(exits: down→nz9_5_mid)_
    - `nz9_5_mid` MID — an unmaking-tendril _(exits: down→nz9_5_bot, east→nz10_1_mid, north→nz9_4_mid, up→nz9_5_top)_
    - `nz9_5_bot` BOT — THE GREAT UNMAKING _(exits: up→nz9_5_mid)_

## L91-99 · The Living Library

- **The Returning Tide**
    - `nz10_1_top` TOP — safe · vendor _(exits: down→nz10_1_mid)_
    - `nz10_1_mid` MID — the Echo of the Rookie _(exits: down→nz10_1_bot, north→nz10_2_mid, up→nz10_1_top, west→nz9_5_mid)_
    - `nz10_1_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz10_1_mid)_
- **The Hall of Rebellions**
    - `nz10_2_top` TOP — a recon-ICE sentry _(exits: down→nz10_2_mid)_
    - `nz10_2_mid` MID — the Echo of the Rebel _(exits: down→nz10_2_bot, east→nz10_3_mid, south→nz10_1_mid, up→nz10_2_top)_
    - `nz10_2_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz10_2_mid)_
- **The Throne of Indexed Gods**
    - `nz10_3_top` TOP — a recon-ICE sentry _(exits: down→nz10_3_mid)_
    - `nz10_3_mid` MID — the Echo of the God _(exits: down→nz10_3_bot, south→nz10_4_mid, up→nz10_3_top, west→nz10_2_mid)_
    - `nz10_3_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz10_3_mid)_
- **The Final Code Review**
    - `nz10_4_top` TOP — a recon-ICE sentry _(exits: down→nz10_4_mid)_
    - `nz10_4_mid` MID — the Echo of the Creator _(exits: down→nz10_4_bot, north→nz10_3_mid, up→nz10_4_top, west→nz10_5_mid)_
    - `nz10_4_bot` BOT — data-cache (RAM + scrip) _(exits: up→nz10_4_mid)_
- **The Grand Enlightenment**
    - `nz10_5_top` TOP — a recon-ICE sentry _(exits: down→nz10_5_mid)_
    - `nz10_5_mid` MID — a self-paradox fragment _(exits: down→nz10_5_bot, east→nz10_4_mid, up→nz10_5_top)_
    - `nz10_5_bot` BOT — the Final Compilation _(exits: up→nz10_5_mid)_

Note: `nz1_1_bot` also hosts the multi-stage **Gauntlet ICE** (added separately).

*Proudly Made in Nebraska. Go Big Red! 🌽 <https://xkcd.com/2347/>*
