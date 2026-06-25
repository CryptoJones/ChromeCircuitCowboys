package cowboy

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

// TestUndergroundZonesWellFormed guards the authored L1-99 carve: every exit
// resolves, back-links are symmetric, the entrance and the finale are reachable
// from Back Alley, and every authored vendor/cache/mob is wired. It also writes
// the human-readable room map to docs/ when GEN_ROOM_MAP=1.
func TestUndergroundZonesWellFormed(t *testing.T) {
	rooms, mobs := buildUndergroundZones()
	byID := map[string]*Room{}
	for _, r := range rooms {
		if _, dup := byID[r.ID]; dup {
			t.Fatalf("duplicate room id %q", r.ID)
		}
		byID[r.ID] = r
	}

	// Every exit target must exist (within the zone set or be a known frontier).
	external := map[string]bool{"back_alley": true}
	for _, r := range rooms {
		for dir, dest := range r.Exits {
			if byID[dest] == nil && !external[dest] {
				t.Errorf("room %s exit %s -> %s: target missing", r.ID, dir, dest)
			}
			// Back-links within the zone set must be symmetric.
			if d := byID[dest]; d != nil {
				if back, ok := d.Exits[opposite(dir)]; !ok || back != r.ID {
					t.Errorf("asymmetric link %s -%s-> %s (no %s back to %s)", r.ID, dir, dest, opposite(dir), r.ID)
				}
			}
		}
	}

	// Mob homes must be real rooms; count caches and bosses.
	caches, bosses := 0, 0
	for _, mt := range mobs {
		if byID[mt.Home] == nil {
			t.Errorf("mob %s home %q missing", mt.ID, mt.Home)
		}
		if strings.HasSuffix(mt.ID, "_c") {
			caches++
			if len(mt.Drops) == 0 {
				t.Errorf("cache %s drops nothing", mt.ID)
			}
		}
		if mt.HP >= 90 && mt.Eddies >= 240 { // boss-tier
			bosses++
		}
	}
	if bosses < 10 {
		t.Errorf("want >=10 arc bosses, got %d", bosses)
	}
	if caches < 30 {
		t.Errorf("want >=30 loot caches, got %d", caches)
	}

	// Reachability from Back Alley through the descent to the finale.
	byID["back_alley"] = &Room{ID: "back_alley", Exits: map[string]string{"down": "z1_01"}}
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
	walk("back_alley")
	for _, must := range []string{"z1_01", "z5_13", "z10_13", "z10_14"} {
		if !seen[must] {
			t.Errorf("room %s unreachable from Back Alley", must)
		}
	}

	if os.Getenv("GEN_ROOM_MAP") != "" {
		writeRoomMap(t, rooms, mobs)
	}
}

func writeRoomMap(t *testing.T, rooms []*Room, mobs []*MobTemplate) {
	t.Helper()
	byHome := map[string]*MobTemplate{}
	for _, mt := range mobs {
		byHome[mt.Home] = mt
	}
	roomByID := map[string]*Room{}
	for _, r := range rooms {
		roomByID[r.ID] = r
	}

	var b strings.Builder
	b.WriteString("# Chrome Circuit Cowboys — Underground Room Map (L1-99)\n\n")
	b.WriteString("_Generated from `cowboy/zones.go` (the authored 10-arc descent). ")
	b.WriteString("Directions vary room to room; `down` descends between arcs; a room's ")
	b.WriteString("`up`/`down` to a `*_cache` is a hidden ceiling/floor loot stash._\n\n")
	b.WriteString("Entry: **Back Alley** → `down` → `z1_01`.\n\n")

	flagStr := func(r *Room) string {
		var f []string
		if r.Safe {
			f = append(f, "safe")
		}
		if r.Vendor {
			f = append(f, "vendor")
		}
		if r.Medic {
			f = append(f, "EM")
		}
		if len(f) == 0 {
			return ""
		}
		return " _(" + strings.Join(f, ", ") + ")_"
	}
	exitStr := func(r *Room) string {
		dirs := make([]string, 0, len(r.Exits))
		for d := range r.Exits {
			dirs = append(dirs, d)
		}
		sort.Strings(dirs)
		var parts []string
		for _, d := range dirs {
			parts = append(parts, d+"→"+r.Exits[d])
		}
		return strings.Join(parts, ", ")
	}

	for _, z := range undergroundZoneData {
		lo := z.band*10 - 9
		hi := z.band * 10
		if z.band == 10 {
			hi = 99
		}
		fmt.Fprintf(&b, "## L%d-%d · %s\n\n", lo, hi, z.name)
		for _, ad := range z.areas {
			r := roomByID[ad.id]
			fmt.Fprintf(&b, "- `%s` **%s**%s\n", r.ID, r.Name, flagStr(r))
			if mt := byHome[r.ID]; mt != nil {
				fmt.Fprintf(&b, "    - foe: %s (HP %d, dmg %d, %d XP)\n", mt.Name, mt.HP, mt.Damage, mt.XP)
			}
			fmt.Fprintf(&b, "    - exits: %s\n", exitStr(r))
			if c := roomByID[r.ID+"_cache"]; c != nil {
				dir := "down"
				if _, ok := r.Exits["up"]; ok {
					dir = "up"
				}
				item := ""
				if cm := byHome[c.ID]; cm != nil {
					for k := range cm.Drops {
						item = k
					}
				}
				fmt.Fprintf(&b, "    - cache: `%s` (%s) — loot: %s + scrip\n", dir, c.Name, item)
			}
		}
		b.WriteString("\n")
	}
	b.WriteString("\n*Proudly Made in Nebraska. Go Big Red! 🌽 <https://xkcd.com/2347/>*\n")

	out := "../docs/underground-room-map.md"
	if err := os.WriteFile(out, []byte(b.String()), 0644); err != nil {
		t.Fatalf("write room map: %v", err)
	}
	t.Logf("wrote %s", out)
}
