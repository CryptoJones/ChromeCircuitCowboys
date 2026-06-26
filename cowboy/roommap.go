package cowboy

import "strings"

// The MAP command draws a quick CP437/ANSI "you-are-here" of the current area:
// every exit labelled with where it leads (and whether it goes DEEPER into a
// harder zone, BACK toward easier ground, or off to a cache/shop), plus the one
// move that takes you onward to the next harder area or back out. It works the
// same in the Undercity, the Net, and the surface — it reads the live room
// graph, so it never drifts from the authored world.

var mapDirs = []string{"north", "south", "east", "west", "up", "down", "in", "out"}

// areaInfo classifies a room id into a realm ("meat", "net", or "city") and a
// zone number within that realm (0 for the surface). Difficulty rises with the
// zone number inside a realm, so neighbours can be tagged harder/back.
func areaInfo(id string) (realm string, zone int) {
	switch {
	case strings.HasPrefix(id, "nz"):
		return "net", leadingInt(id[2:])
	case len(id) > 1 && id[0] == 'z' && id[1] >= '0' && id[1] <= '9':
		return "meat", leadingInt(id[1:])
	default:
		return "city", 0
	}
}

// leadingInt reads the run of digits at the start of s (stopping at the first
// non-digit, e.g. the "_" in "z12_03"), returning 0 if there are none.
func leadingInt(s string) int {
	n := 0
	for i := 0; i < len(s) && s[i] >= '0' && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
	}
	return n
}

// levelRange turns a zone number into its 10-level band label (zone 10 caps at 99).
func levelRange(zone int) string {
	lo, hi := zone*10-9, zone*10
	if zone >= 10 {
		hi = 99
	}
	return "L" + itoa(lo) + "-" + itoa(hi)
}

// areaLabel is the headline for the current area.
func areaLabel(realm string, zone int) string {
	switch realm {
	case "meat":
		name := ""
		if zone >= 1 && zone <= len(undergroundZoneData) {
			name = undergroundZoneData[zone-1].name
		}
		return "THE UNDERCITY " + levelRange(zone) + " — " + name
	case "net":
		name := ""
		if zone >= 1 && zone <= len(netZoneData) {
			name = netZoneData[zone-1].name
		}
		return "THE NET " + levelRange(zone) + " — " + name
	default:
		return "NIGHT CITY — the surface"
	}
}

// exitLabel describes where a single exit leads, from the player's vantage.
func (w *World) exitLabel(here, dest string) string {
	if strings.HasSuffix(dest, "_cache") {
		return style(green, "a sealed cache") + style(dim, " — break it for loot")
	}
	dr := w.room(dest)
	name := dest
	if dr != nil && dr.Name != "" {
		name = dr.Name
	}
	hRealm, hZone := areaInfo(here)
	dRealm, dZone := areaInfo(dest)
	tag := ""
	switch {
	case dRealm != hRealm:
		switch dRealm {
		case "net":
			tag = style(neon, "  ⤓ jack into the Net")
		case "meat":
			tag = style(gold, "  ↓ down into the Undercity")
		case "city":
			tag = style(dim, "  ↑ up to the surface")
		}
	case dZone > hZone:
		tag = style(gold, "  ▼ deeper — HARDER")
	case dZone < hZone:
		tag = style(dim, "  ▲ back — easier")
	}
	flags := ""
	if dr != nil {
		if dr.Vendor {
			flags += style(gold, " $shop")
		}
		if dr.Medic {
			flags += style(gold, " +medic")
		}
		if dr.Safe {
			flags += style(green, " ·safe")
		}
	}
	return style(green, name) + tag + flags
}

// onwardStep finds the single move that starts you toward the next area: the
// nearest room in this zone with an exit crossing into a harder zone (harder
// true) or an easier one (harder false). It BFS-walks only rooms of the current
// realm+zone (skipping dead-end caches), so the first step it returns is always
// a real direction you can type now. ok is false when there's no such frontier
// (e.g. the deepest zone has no "deeper", the surface has no "back").
func (w *World) onwardStep(start string, harder bool) (dir string, ok bool) {
	sRealm, sZone := areaInfo(start)
	type node struct{ id, first string }
	seen := map[string]bool{start: true}
	queue := []node{{start, ""}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		r := w.room(cur.id)
		if r == nil {
			continue
		}
		for _, d := range mapDirs {
			dest, has := r.Exits[d]
			if !has {
				continue
			}
			first := cur.first
			if first == "" {
				first = d
			}
			dRealm, dZone := areaInfo(dest)
			if dRealm == sRealm && ((harder && dZone > sZone) || (!harder && dZone < sZone)) {
				return first, true
			}
			if seen[dest] || strings.HasSuffix(dest, "_cache") {
				continue
			}
			if dRealm == sRealm && dZone == sZone {
				seen[dest] = true
				queue = append(queue, node{dest, first})
			}
		}
	}
	return "", false
}

// showMap renders the MAP command.
func (w *World) showMap(p *Player) {
	r := w.room(p.RoomID)
	if r == nil {
		p.send(style(red, "No map signal — your location is corrupted.") + crlf)
		return
	}
	realm, zone := areaInfo(p.RoomID)
	const rule = "  ════════════════════════════════════════"

	var b strings.Builder
	b.WriteString(crlf + style(neon, "  ╔══════════════ LOCAL MAP ══════════════╗") + crlf)
	b.WriteString("  " + style(gold, areaLabel(realm, zone)) + crlf)
	b.WriteString("  " + style(dim, "you are here ▸ ") + style(hot, r.Name) + crlf)
	b.WriteString(style(neon, rule) + crlf)

	any := false
	for _, d := range mapDirs {
		dest, has := r.Exits[d]
		if !has {
			continue
		}
		any = true
		b.WriteString("   " + style(hot, padRight(strings.ToUpper(d), 6)) + style(dim, "→ ") + w.exitLabel(p.RoomID, dest) + crlf)
	}
	if !any {
		b.WriteString(style(dim, "   (no exits — you're boxed in)") + crlf)
	}
	b.WriteString(style(neon, rule) + crlf)

	if dir, ok := w.onwardStep(p.RoomID, true); ok {
		b.WriteString("  " + style(gold, "▼ PROCEED to the next harder area: go "+strings.ToUpper(dir)) + crlf)
	}
	if dir, ok := w.onwardStep(p.RoomID, false); ok {
		b.WriteString("  " + style(dim, "▲ WAY OUT toward easier ground: go "+strings.ToUpper(dir)) + crlf)
	}
	b.WriteString(style(neon, "  ╚════════════════════════════════════════╝") + crlf)
	p.send(b.String())
}

// padRight pads s with spaces to width w (never truncates).
func padRight(s string, w int) string {
	for len(s) < w {
		s += " "
	}
	return s
}
