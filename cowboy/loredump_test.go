package cowboy

// loredump_test.go is a DEV-TIME helper (not a real test): it serializes the
// live, assembled game content — every room, item, zone-lore line and flavor NPC,
// with the context a writer needs — to a JSON manifest. The Violet Lotus content
// generator (tools/violet-gen) reads that manifest, rewrites the prose, and emits
// cowboy/lore.json, which the game embeds and overlays at build time (lore.go).
//
// It only runs when LORE_DUMP points at an output path:
//   LORE_DUMP=/tmp/manifest.json go test ./cowboy -run TestDumpLoreManifest
// so a normal `go test` skips it.

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

type dumpRoom struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	ZoneKey  string `json:"zone_key"`
	ZoneName string `json:"zone_name"`
}

type dumpItem struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Effect string `json:"effect"`
}

type dumpNPC struct {
	Speaker string   `json:"speaker"`
	Lines   []string `json:"lines"`
}

type loreManifest struct {
	Rooms     []dumpRoom          `json:"rooms"`
	Items     []dumpItem          `json:"items"`
	ZoneLore  map[string][]string `json:"zone_lore"`
	ZoneNames map[string]string   `json:"zone_names"`
	RoomNPC   map[string]dumpNPC  `json:"room_npc"`
	Booth     []string            `json:"booth"`
}

// zoneNameByKey builds key->display-name for every authored zone, plus the rings
// and the city surface, so the generator can anchor prose to the right place.
func zoneNameByKey() map[string]string {
	names := map[string]string{
		"":   "Noche City (the neon surface streets)",
		"ic": "The Inner Circuit (RP-safe neon maglev loop)",
		"sb": "The Sprawlbelt (ground-level beltway)",
	}
	for _, z := range undergroundZoneData {
		names[z.key] = z.name
	}
	for _, z := range netZoneData {
		names[z.key] = z.name
	}
	return names
}

// zoneKeyForRoom derives the zone key from a room id (z1_05 -> z1, nz3_2_mid ->
// nz3, ic_1 -> ic, sb_gym -> sb). City rooms have no zone ("").
func zoneKeyForRoom(id string) string {
	switch {
	case strings.HasPrefix(id, "nz"):
		if i := strings.IndexByte(id, '_'); i > 0 {
			return id[:i]
		}
	case strings.HasPrefix(id, "z"):
		if i := strings.IndexByte(id, '_'); i > 0 {
			return id[:i]
		}
	case strings.HasPrefix(id, "ic_"):
		return "ic"
	case strings.HasPrefix(id, "sb_"):
		return "sb"
	}
	return ""
}

func itemEffect(x ware) string {
	var fx []string
	add := func(s string) { fx = append(fx, s) }
	if x.heal > 0 {
		add("restores " + itoa(x.heal) + " HP")
	}
	if x.ram > 0 {
		add("restores " + itoa(x.ram) + " RAM")
	}
	if x.bonus > 0 {
		add("+" + itoa(x.bonus) + " attack")
	}
	if x.deck > 0 {
		add("+" + itoa(x.deck) + " max RAM")
	}
	if x.body > 0 {
		add("+" + itoa(x.body) + " Body implant")
	}
	if x.refl > 0 {
		add("+" + itoa(x.refl) + " Reflexes implant")
	}
	if x.intel > 0 {
		add("+" + itoa(x.intel) + " Intelligence implant")
	}
	if x.forClass != "" {
		add(x.forClass + "-only")
	}
	return strings.Join(fx, ", ")
}

func TestDumpLoreManifest(t *testing.T) {
	out := os.Getenv("LORE_DUMP")
	if out == "" {
		t.Skip("set LORE_DUMP=<path> to dump the lore manifest")
	}
	names := zoneNameByKey()
	man := loreManifest{
		ZoneLore:  zoneLore,
		ZoneNames: names,
		RoomNPC:   map[string]dumpNPC{},
		Booth:     boothIntro,
	}
	for id, r := range buildRooms() {
		zk := zoneKeyForRoom(id)
		man.Rooms = append(man.Rooms, dumpRoom{
			ID: id, Name: r.Name, Desc: r.Desc, ZoneKey: zk, ZoneName: names[zk],
		})
	}
	for _, x := range shopWares {
		man.Items = append(man.Items, dumpItem{Name: x.name, Desc: x.desc, Effect: itemEffect(x)})
	}
	for id, n := range roomNPC {
		man.RoomNPC[id] = dumpNPC{Speaker: n.speaker, Lines: n.lines}
	}
	b, err := json.MarshalIndent(man, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(out, b, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("wrote %d rooms, %d items to %s", len(man.Rooms), len(man.Items), out)
}
