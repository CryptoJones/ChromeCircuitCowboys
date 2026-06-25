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

	svg := buildRoomMapSVG(roomByID, byHome)
	svgOut := "../docs/underground-room-map.svg"
	if err := os.WriteFile(svgOut, []byte(svg), 0644); err != nil {
		t.Fatalf("write room map svg: %v", err)
	}
	t.Logf("wrote %s", svgOut)
}

func xmlEsc(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func trunc(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "."
}

var dirLetter = map[string]string{"north": "N", "south": "S", "east": "E", "west": "W", "up": "U", "down": "D"}

// buildRoomMapSVG renders the carved underground as a Cyberdeck-dark node map:
// one column per arc, rooms stacked in descent order with the actual varied
// exit direction labelled between them, plus foe/boss/vendor/cache markers.
func buildRoomMapSVG(roomByID map[string]*Room, byHome map[string]*MobTemplate) string {
	const (
		x0       = 40
		colW     = 384
		boxW     = 336
		boxH     = 58
		vgap     = 30
		hdrY     = 96
		rowStart = 132
	)
	maxRows := 0
	for _, z := range undergroundZoneData {
		if len(z.areas) > maxRows {
			maxRows = len(z.areas)
		}
	}
	width := x0*2 + len(undergroundZoneData)*colW
	height := rowStart + maxRows*(boxH+vgap) + 90

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="Menlo, Consolas, monospace">`, width, height, width, height)
	fmt.Fprintf(&b, `<rect x="0" y="0" width="%d" height="%d" fill="#07090f"/>`, width, height)
	b.WriteString(`<text x="40" y="40" fill="#27d4ff" font-size="26" font-weight="bold">CHROME CIRCUIT COWBOYS — UNDERGROUND ROOM MAP (L1-99, carved)</text>`)
	b.WriteString(`<text x="40" y="70" fill="#8a97ab" font-size="15">Each column = one arc, rooms in descent order. Letters between rooms = the actual exit direction (varies). [^]/[v] = ceiling/floor loot cache. Entry: Back Alley D-&gt; z1_01.</text>`)

	for zi, z := range undergroundZoneData {
		colX := x0 + zi*colW
		lo, hi := z.band*10-9, z.band*10
		if z.band == 10 {
			hi = 99
		}
		fmt.Fprintf(&b, `<text x="%d" y="%d" fill="#27d4ff" font-size="16" font-weight="bold">L%d-%d  %s</text>`, colX, hdrY, lo, hi, xmlEsc(trunc(z.name, 26)))
		for ai, ad := range z.areas {
			r := roomByID[ad.id]
			boxY := rowStart + ai*(boxH+vgap)
			kind, mobName := "", ""
			if ad.mob != "" {
				kind, mobName = splitMob(ad.mob)
			}
			stroke, fill := "#5a6678", "#0e1320"
			switch {
			case kind == "b":
				stroke, fill = "#ff4f4f", "#1a0f0f"
			case r.Vendor || r.Medic:
				stroke = "#ffb000"
				if r.Safe {
					stroke, fill = "#55ff99", "#0c1810"
				}
			case r.Safe:
				stroke = "#55ff99"
			case kind == "e":
				stroke = "#ff8a3d"
			case kind == "c":
				stroke = "#27d4ff"
			}
			fmt.Fprintf(&b, `<rect x="%d" y="%d" width="%d" height="%d" rx="7" fill="%s" stroke="%s" stroke-width="2"/>`, colX, boxY, boxW, boxH, fill, stroke)
			fmt.Fprintf(&b, `<text x="%d" y="%d" fill="%s" font-size="13" font-weight="bold">%s  %s</text>`, colX+8, boxY+20, stroke, ad.id, xmlEsc(trunc(r.Name, 28)))
			var marks []string
			if kind == "b" {
				marks = append(marks, "BOSS: "+trunc(mobName, 22))
			} else if mobName != "" {
				marks = append(marks, "foe: "+trunc(mobName, 22))
			}
			if r.Safe {
				marks = append(marks, "safe")
			}
			if r.Vendor {
				marks = append(marks, "vendor")
			}
			if r.Medic {
				marks = append(marks, "EM")
			}
			if ad.cache == "up" {
				marks = append(marks, "[^loot]")
			} else if ad.cache == "down" {
				marks = append(marks, "[v loot]")
			}
			fmt.Fprintf(&b, `<text x="%d" y="%d" fill="#9fb0c4" font-size="11">%s</text>`, colX+8, boxY+40, xmlEsc(trunc(strings.Join(marks, " · "), 46)))
			if ai+1 < len(z.areas) {
				nextID := z.areas[ai+1].id
				fwd := "?"
				for d, dest := range r.Exits {
					if dest == nextID {
						if l, ok := dirLetter[d]; ok {
							fwd = l
						}
						break
					}
				}
				cx := colX + boxW/2
				ly := boxY + boxH + vgap/2 + 4
				fmt.Fprintf(&b, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a4658" stroke-width="2"/>`, cx, boxY+boxH, cx, boxY+boxH+vgap)
				fmt.Fprintf(&b, `<circle cx="%d" cy="%d" r="10" fill="#0e1320" stroke="#5a6678" stroke-width="1.5"/>`, cx, ly-4)
				fmt.Fprintf(&b, `<text x="%d" y="%d" fill="#d7e0ee" font-size="12" font-weight="bold" text-anchor="middle">%s</text>`, cx, ly, fwd)
			}
		}
		if zi+1 < len(undergroundZoneData) {
			lastY := rowStart + (len(z.areas)-1)*(boxH+vgap) + boxH + 22
			hiNext := z.band*10 + 10
			if z.band == 9 {
				hiNext = 99
			}
			fmt.Fprintf(&b, `<text x="%d" y="%d" fill="#ffb000" font-size="12" font-weight="bold">D-&gt; descend to L%d-%d</text>`, colX+8, lastY, z.band*10+1, hiNext)
		}
	}

	ly := height - 26
	fmt.Fprintf(&b, `<text x="40" y="%d" fill="#9fb0c4" font-size="13">Legend:  `, ly)
	b.WriteString(`<tspan fill="#55ff99">green</tspan>=safe/hub  <tspan fill="#ffb000">amber</tspan>=vendor/EM  <tspan fill="#27d4ff">cyan</tspan>=combat  <tspan fill="#ff8a3d">orange</tspan>=elite  <tspan fill="#ff4f4f">red</tspan>=arc boss  ·  N/S/E/W = exit direction  ·  [^loot]/[v loot] = ceiling/floor cache</text>`)
	b.WriteString(`</svg>`)
	return b.String()
}
