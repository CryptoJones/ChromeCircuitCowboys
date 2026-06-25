package cowboy

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

// TestNetZonesWellFormed guards the authored Net carve: every area is a clean
// 3-layer stack (top<->mid<->bot via up/down), the MID layers chain end to end,
// every room is flagged Net, the jack-in reaches the first shell and the finale
// is reachable, and the bosses + data-caches are placed. With GEN_ROOM_MAP=1 it
// regenerates docs/net-room-map.md.
func TestNetZonesWellFormed(t *testing.T) {
	rooms, mobs := buildNetZones()
	byID := map[string]*Room{}
	for _, r := range rooms {
		if _, dup := byID[r.ID]; dup {
			t.Fatalf("duplicate net room id %q", r.ID)
		}
		if !r.Net {
			t.Errorf("net room %s missing Net flag", r.ID)
		}
		byID[r.ID] = r
	}

	// Each area is top -down-> mid -down-> bot, with symmetric up links.
	for _, z := range netZoneData {
		for ai := range z.areas {
			base := fmt.Sprintf("%s_%d", z.key, ai+1)
			top, mid, bot := byID[base+"_top"], byID[base+"_mid"], byID[base+"_bot"]
			if top == nil || mid == nil || bot == nil {
				t.Fatalf("area %s missing a layer", base)
			}
			if top.Exits["down"] != mid.ID || mid.Exits["up"] != top.ID {
				t.Errorf("%s: top<->mid link broken", base)
			}
			if mid.Exits["down"] != bot.ID || bot.Exits["up"] != mid.ID {
				t.Errorf("%s: mid<->bot link broken", base)
			}
		}
	}

	// MID layers chain symmetrically (the lateral thoroughfare).
	for _, r := range rooms {
		if !strings.HasSuffix(r.ID, "_mid") {
			continue
		}
		for dir, dest := range r.Exits {
			if dir == "up" || dir == "down" {
				continue
			}
			d := byID[dest]
			if d == nil || d.Exits[opposite(dir)] != r.ID {
				t.Errorf("asymmetric mid link %s -%s-> %s", r.ID, dir, dest)
			}
		}
	}

	// Reachability from the Data Port jack-in through the whole ascent.
	byID["data_port"] = &Room{ID: "data_port", Exits: map[string]string{"up": "nz1_1_top"}}
	byID["nz1_1_top"].Exits["up"] = "data_port"
	seen := map[string]bool{}
	var walk func(id string)
	walk = func(id string) {
		if seen[id] || byID[id] == nil {
			return
		}
		seen[id] = true
		for _, dest := range byID[id].Exits {
			walk(dest)
		}
	}
	walk("data_port")
	for _, must := range []string{"nz1_1_top", "nz1_1_mid", "nz5_1_mid", "nz10_5_mid", "nz10_5_bot"} {
		if !seen[must] {
			t.Errorf("net room %s unreachable from the Data Port", must)
		}
	}

	bosses, caches := 0, 0
	for _, mt := range mobs {
		if strings.HasSuffix(mt.ID, "_c") {
			caches++
			if len(mt.Drops) == 0 {
				t.Errorf("data-cache %s drops nothing", mt.ID)
			}
		}
		if mt.HP >= 90+55 && mt.Eddies >= 240 { // boss-tier (band>=1)
			bosses++
		}
	}
	if bosses < 10 {
		t.Errorf("want >=10 arc bosses, got %d", bosses)
	}
	if caches < 35 {
		t.Errorf("want >=35 data-caches, got %d", caches)
	}

	if os.Getenv("GEN_ROOM_MAP") != "" {
		writeNetMap(t, byID, mobs)
	}
}

func writeNetMap(t *testing.T, byID map[string]*Room, mobs []*MobTemplate) {
	t.Helper()
	byHome := map[string]*MobTemplate{}
	for _, mt := range mobs {
		byHome[mt.Home] = mt
	}
	var b strings.Builder
	b.WriteString("# Chrome Circuit Cowboys — The Net Room Map (L1-99)\n\n")
	b.WriteString("_Generated from `cowboy/netzones.go`. Each area is a 3-layer stack — ")
	b.WriteString("`:: Shell` (TOP, access), `:: Breach` (MID, the lateral thoroughfare + ICE fight), ")
	b.WriteString("`:: Core` (BOT, data-vault / boss). UP/DOWN between layers; N/S/E/W between areas ")
	b.WriteString("(at the MID layer). Jack in: Data Port → `up` → `nz1_1_top`._\n\n")

	foe := func(id string) string {
		if mt := byHome[id]; mt != nil {
			if strings.HasSuffix(mt.ID, "_c") {
				return "data-cache (RAM + scrip)"
			}
			return mt.Name
		}
		return ""
	}
	exits := func(r *Room) string {
		ds := make([]string, 0, len(r.Exits))
		for d := range r.Exits {
			ds = append(ds, d)
		}
		sort.Strings(ds)
		var ps []string
		for _, d := range ds {
			ps = append(ps, d+"→"+r.Exits[d])
		}
		return strings.Join(ps, ", ")
	}

	for _, z := range netZoneData {
		lo, hi := z.band*10-9, z.band*10
		if z.band == 10 {
			hi = 99
		}
		fmt.Fprintf(&b, "## L%d-%d · %s\n\n", lo, hi, z.name)
		for ai, ar := range z.areas {
			base := fmt.Sprintf("%s_%d", z.key, ai+1)
			fmt.Fprintf(&b, "- **%s**\n", ar.name)
			for _, layer := range []string{"top", "mid", "bot"} {
				r := byID[base+"_"+layer]
				tag := strings.ToUpper(layer)
				marks := []string{}
				if r.Safe {
					marks = append(marks, "safe")
				}
				if r.Vendor {
					marks = append(marks, "vendor")
				}
				if f := foe(r.ID); f != "" {
					marks = append(marks, f)
				}
				fmt.Fprintf(&b, "    - `%s` %s — %s _(exits: %s)_\n", r.ID, tag, strings.Join(marks, " · "), exits(r))
			}
		}
		b.WriteString("\n")
	}
	b.WriteString("Note: `nz1_1_bot` also hosts the multi-stage **Gauntlet ICE** (added separately).\n\n")
	b.WriteString("*Proudly Made in Nebraska. Go Big Red! 🌽 <https://xkcd.com/2347/>*\n")

	out := "../docs/net-room-map.md"
	if err := os.WriteFile(out, []byte(b.String()), 0644); err != nil {
		t.Fatalf("write net room map: %v", err)
	}
	t.Logf("wrote %s", out)
}
