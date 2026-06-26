package cowboy

import (
	"sort"
	"strings"
)

// AI runners (#37): lightweight bots that look like players in WHO and in the
// room, wander the surface, and bark ambient lines when a real runner is around
// — so someone jacking into a quiet server isn't alone. They never fight, are
// never targeted by mobs, and are never written to the character store. They
// stay on the city surface (realm "city") so they never wander into Net PvP or
// Undercity ICE, sidestepping every combat/persistence edge case.

// botProfile is one AI runner's handle + class flavor (shown in WHO).
type botProfile struct {
	name  string
	class string
}

// botRoster is the pool of AI runners. EnableBots takes the first N. Handles are
// deliberately distinctive so they don't collide with real character names.
var botRoster = []botProfile{
	{"Riko-Vex", "Netrunner"},
	{"Static-Jane", "Solo"},
	{"Coil", "Techie"},
	{"Mirrorface", "Fixer"},
	{"Ozone", "Solo"},
	{"Dr-Patch", "Medtech"},
	{"Lowkey", "Netrunner"},
	{"Brick-Tao", "Solo"},
	{"Vega-Cruz", "Solo"},
	{"Nyx", "Netrunner"},
	{"Halcyon", "Fixer"},
	{"Rustbucket", "Techie"},
	{"Saint-Iggy", "Medtech"},
	{"Decker-Mo", "Netrunner"},
	{"Glitch-Marie", "Netrunner"},
	{"Bonesaw", "Solo"},
	{"Pixel", "Techie"},
	{"Mama-Voltage", "Fixer"},
	{"Reyes", "Solo"},
	{"Quill", "Media"},
	{"Dredge", "Solo"},
	{"Cinder-Lou", "Rocker"},
	{"Hex", "Netrunner"},
	{"Tox", "Techie"},
	{"Goldtooth", "Fixer"},
	{"Patchwork", "Medtech"},
	{"Switchblade-Su", "Solo"},
	{"Echo-Naught", "Netrunner"},
	{"Grit", "Nomad"},
	{"Wire-Mother", "Techie"},
	{"Calaca", "Solo"},
	{"Zero-Bahn", "Netrunner"},
}

// botLines are the one-liners an AI runner says when a real runner shares the room.
var botLines = []string{
	"Watch the alleys tonight, choom — Kurokawa's sweeping again.",
	"You holding? I'll trade scrip for anything that breaks ICE.",
	"Heard the Undercity's hot. Good loot, worse odds.",
	"Stay frosty. The grid eats the careless.",
	"Anybody seen a fixer hiring? My pockets are running on fumes.",
	"Fresh sleeve, fresh start. Same old Night City.",
	"Don't trust the medics on the outer ring. Trust me.",
	"They say the deeper Net layers think now. I say jack out while you can.",
}

// botEmotes are ambient actions an AI runner performs (rendered like a player emote).
var botEmotes = []string{
	"lights a cigarette and watches the street.",
	"checks a battered wrist-deck, frowns, moves on.",
	"leans against the wall, scanning faces.",
	"flicks a coin and pockets it.",
	"mutters into a dead comm-line.",
}

// EnableBots seeds n AI runners across the surface. Call once after NewWorld
// (the server does; tests opt in). No-op for n<=0; capped at the roster size.
// Names already taken by a real session are skipped.
func (w *World) EnableBots(n int) {
	if n <= 0 {
		return
	}
	if n > len(botRoster) {
		n = len(botRoster)
	}
	spots := w.citySpots()
	if len(spots) == 0 {
		return
	}
	for i := 0; i < n; i++ {
		prof := botRoster[i]
		if _, taken := w.byName[strings.ToLower(prof.name)]; taken {
			continue
		}
		b := w.newPlayer(prof.name, func(string) {}) // output discarded
		b.IsBot = true
		b.Class = prof.class
		b.Level = 2 + i%6 // a little spread in WHO
		b.HP, b.MaxHP = 30, 30
		b.RoomID = spots[(i*7)%len(spots)] // fan them out across the surface
		w.players[b.ID] = b
		w.byName[strings.ToLower(b.Name)] = b
	}
}

// citySpots is the sorted set of public surface rooms bots may occupy/wander —
// never private pods, the Undercity, or the Net.
func (w *World) citySpots() []string {
	var spots []string
	for id, r := range w.rooms {
		if r.Private {
			continue
		}
		if realm, _ := areaInfo(id); realm != "city" {
			continue
		}
		spots = append(spots, id)
	}
	sort.Strings(spots) // deterministic order (for seeded tests)
	return spots
}

// tickBots drives every AI runner one step: mostly idle, sometimes wander,
// sometimes chatter. Called from Tick.
func (w *World) tickBots() {
	for _, b := range w.players {
		if !b.IsBot {
			continue
		}
		if b.party != nil { // crewed up: shut up and follow — no chatter, no wandering off
			continue
		}
		switch w.roll(6) { // ~2/3 of ticks the bot just idles
		case 0:
			w.botWander(b)
		case 1:
			w.botChatter(b)
		}
	}
}

// botWander steps a bot to a random adjacent surface room (broadcasting the
// arrival/departure that real runners in those rooms see).
func (w *World) botWander(b *Player) {
	r := w.room(b.RoomID)
	if r == nil {
		return
	}
	var dirs []string
	for dir, dest := range r.Exits {
		if realm, _ := areaInfo(dest); realm != "city" {
			continue // stay on the surface
		}
		if dr := w.room(dest); dr == nil || dr.Private {
			continue
		}
		dirs = append(dirs, dir)
	}
	if len(dirs) == 0 {
		return
	}
	sort.Strings(dirs) // deterministic before the seeded pick
	w.move(b, dirs[w.roll(len(dirs))])
}

// botChatter emits a bark — but only when a real runner is in the room, so the
// flavor lands where it's meant to (relieving the loneliness) and quiet rooms
// stay quiet.
func (w *World) botChatter(b *Player) {
	if !w.realRunnerNear(b) {
		return
	}
	if w.roll(3) == 0 { // a third emotes, the rest speak
		w.broadcast(b.RoomID, b, style(neon, b.Name+" "+botEmotes[w.roll(len(botEmotes))])+crlf)
		return
	}
	w.broadcast(b.RoomID, b, style(green, b.Name+" says: ")+botLines[w.roll(len(botLines))]+crlf)
}

// realRunnerNear reports whether a human player shares the bot's room.
func (w *World) realRunnerNear(b *Player) bool {
	for _, o := range w.playersIn(b.RoomID, b) {
		if !o.IsBot {
			return true
		}
	}
	return false
}
