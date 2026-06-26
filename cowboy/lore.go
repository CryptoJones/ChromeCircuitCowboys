package cowboy

// lore.go embeds and applies the Violet Lotus-generated content overlay
// (lore.json). A BBS door game runs as an isolated process with no network
// access, so the generation happens DEV-TIME (tools/violet-gen, driven by the
// local Violet Lotus model) and the result is baked into the binary here. The
// structural world (exits, flags, mobs, item stats) stays authored in Go; this
// overlay only enriches PROSE: room descriptions (#35), per-zone TALK lore and
// flavor-NPC lines (#39), and item examine-lore (#54).
//
// Regenerate after a content change:
//   LORE_DUMP=manifest.json go test ./cowboy -run TestDumpLoreManifest
//   python3 tools/violet-gen/generate.py --manifest manifest.json --out cowboy/lore.json

import (
	_ "embed"
	"encoding/json"
)

//go:embed lore.json
var loreJSON []byte

type loreOverlay struct {
	Rooms    map[string]string   `json:"rooms"`
	ZoneLore map[string][]string `json:"zone_lore"`
	RoomNPC  map[string][]string `json:"room_npc"`
	Items    map[string]string   `json:"items"`
}

var lore loreOverlay

// init applies the overlay to the package-level content that is already
// initialized by the time init runs: per-zone TALK lore, flavor-NPC line sets,
// and item examine-lore. Room descriptions are applied later, per built map, by
// applyRoomLore (rooms don't exist until buildRooms runs).
func init() {
	if len(loreJSON) == 0 {
		return
	}
	if err := json.Unmarshal(loreJSON, &lore); err != nil {
		// A malformed overlay must not break the game — fall back to the authored
		// prose by leaving everything as-is.
		lore = loreOverlay{}
		return
	}
	// Zone TALK lore and flavor-NPC lines are ADDITIVE: the authored canon lines
	// stay first (so the iconic openers survive and TALK is stable), with the
	// Violet Lotus-generated lines appended for variety (#39).
	for key, lines := range lore.ZoneLore {
		if len(lines) > 0 {
			zoneLore[key] = append(append([]string{}, zoneLore[key]...), lines...)
		}
	}
	for id, lines := range lore.RoomNPC {
		if len(lines) == 0 {
			continue
		}
		if npc, ok := roomNPC[id]; ok {
			npc.lines = append(append([]string{}, npc.lines...), lines...)
			roomNPC[id] = npc
		}
	}
	for i := range shopWares {
		if l, ok := lore.Items[shopWares[i].name]; ok && l != "" {
			shopWares[i].lore = l
		}
	}
}

// applyRoomLore overlays the generated, more-vivid descriptions onto a freshly
// built room map, wrapped to the same width as the authored rooms.
func applyRoomLore(m map[string]*Room) {
	for id, desc := range lore.Rooms {
		if desc == "" {
			continue
		}
		if r, ok := m[id]; ok {
			r.Desc = wrapText(desc, 76)
		}
	}
}
