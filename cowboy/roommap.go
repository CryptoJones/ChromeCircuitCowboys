package cowboy

import "strings"

// The MAP command draws a quick CP437/ANSI "you-are-here" of the current area:
// every exit labelled with where it leads (and whether it goes DEEPER into a
// harder zone, BACK toward easier ground, or off to a cache/shop), plus the one
// move that takes you onward to the next harder area or back out. It works the
// same in the Undercity, the Net, and the surface — it reads the live room
// graph, so it never drifts from the authored world.

var mapDirs = []string{"north", "south", "east", "west", "up", "down", "northeast", "southeast", "southwest", "northwest", "in", "out"}

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
		return "NOCHE CITY — the surface"
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

// areaDepth is how far a room is from the surface: 0 at street level (Night
// City / the Data Port jack point), rising with each Undercity arc and each Net
// zone. The way out is always toward a lower depth — crossing realm boundaries
// (Net → Data Port → surface) when that's the route home.
func areaDepth(id string) int {
	_, zone := areaInfo(id) // realm is irrelevant to "how deep"; only the band is
	return zone             // city/surface (incl. the Data Port) = 0
}

// onwardStep finds the single move that starts you toward the next area.
//
// harder=true → the next harder zone: BFS the current realm+zone for an exit
// crossing UP the difficulty band, return the first step. ok=false in the
// deepest zone.
//
// harder=false → the way OUT toward the surface: BFS toward any exit that
// lowers areaDepth, following realm boundaries (zone-1 Undercity → street, the
// Net → the Data Port → street). ok=false once you're already at the surface.
func (w *World) onwardStep(start string, harder bool) (dir string, ok bool) {
	type node struct{ id, first string }
	seen := map[string]bool{start: true}
	queue := []node{{start, ""}}

	if !harder {
		startDepth := areaDepth(start)
		if startDepth == 0 {
			return "", false // already at the surface
		}
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
				if areaDepth(dest) < startDepth {
					return first, true // an exit toward the surface
				}
				if seen[dest] || strings.HasSuffix(dest, "_cache") {
					continue
				}
				if areaDepth(dest) <= startDepth { // don't wander deeper while seeking the exit
					seen[dest] = true
					queue = append(queue, node{dest, first})
				}
			}
		}
		return "", false
	}

	sRealm, sZone := areaInfo(start)
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
			if dRealm == sRealm && dZone > sZone {
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

// zoneObjectiveStep finds the way to the current zone's objective — the home of a
// quest-target mob in this realm+zone — for when there's no harder area to point
// to (the deepest band). dir == "" with ok means the objective is this very room.
func (w *World) zoneObjectiveStep(start string) (dir string, ok bool) {
	realm, zone := areaInfo(start)
	if zone == 0 {
		return "", false // surface has no single objective
	}
	goals := map[string]bool{}
	for _, q := range quests {
		t, has := w.tmpls[q.Target]
		if !has || t.Home == "" {
			continue
		}
		if tr, tz := areaInfo(t.Home); tr == realm && tz == zone {
			goals[t.Home] = true
		}
	}
	return w.stepToRooms(start, goals)
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
	if g, ok := departGates[p.RoomID]; ok { // gated one-way exit (not in r.Exits)
		any = true
		b.WriteString("   " + style(hot, padRight(strings.ToUpper(g.dir), 6)) + style(dim, "→ ") + style(gold, "a way out — one-way (confirm)") + crlf)
	}
	if !any {
		b.WriteString(style(dim, "   (no exits — you're boxed in)") + crlf)
	}
	b.WriteString(style(neon, rule) + crlf)

	if dir, ok := w.onwardStep(p.RoomID, true); ok {
		b.WriteString("  " + style(gold, "▼ PROCEED to the next harder area: go "+strings.ToUpper(dir)) + crlf)
	} else if dir, ok := w.zoneObjectiveStep(p.RoomID); ok {
		// Deepest band — nothing harder to head for, so point at the final objective.
		if dir == "" {
			b.WriteString("  " + style(gold, "▼ FINAL OBJECTIVE is HERE — end of the line, cowboy.") + crlf)
		} else {
			b.WriteString("  " + style(gold, "▼ FINAL OBJECTIVE: go "+strings.ToUpper(dir)) + crlf)
		}
	}
	if dir, ok := w.onwardStep(p.RoomID, false); ok {
		b.WriteString("  " + style(dim, "▲ WAY OUT toward easier ground: go "+strings.ToUpper(dir)) + crlf)
	}
	// If you're in a crew and got separated, point the way back to them.
	if p.party != nil {
		if dir, who, ok := w.stepToParty(p); ok {
			b.WriteString("  " + style(green, "▸ TO YOUR CREW: go "+strings.ToUpper(dir)+" (toward "+who+")") + crlf)
		} else if w.partyElsewhere(p) {
			b.WriteString("  " + style(dim, "▸ your crew is out of reach from here (try jacking in/out)") + crlf)
		}
	}
	b.WriteString(style(neon, "  ╚════════════════════════════════════════╝") + crlf)
	p.send(b.String())
}

// partyElsewhere reports whether any crewmate is in a different room.
func (w *World) partyElsewhere(p *Player) bool {
	if p.party == nil {
		return false
	}
	for _, m := range p.party.Members {
		if m != p && m.RoomID != p.RoomID {
			return true
		}
	}
	return false
}

// dirShort abbreviates a direction word for compact display (north → N, etc.).
var dirShort = map[string]string{
	"north": "N", "south": "S", "east": "E", "west": "W",
	"up": "U", "down": "D", "in": "IN", "out": "OUT",
	"northeast": "NE", "southeast": "SE", "southwest": "SW", "northwest": "NW",
}

// stepToRooms BFS-walks the room graph from start to the nearest room in goals
// and returns the first move toward it. dir == "" with ok == true means start is
// already one of the goal rooms ("you're here"). ok == false means unreachable.
func (w *World) stepToRooms(start string, goals map[string]bool) (dir string, ok bool) {
	if len(goals) == 0 {
		return "", false
	}
	if goals[start] {
		return "", true
	}
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
			if goals[dest] {
				return first, true
			}
			if seen[dest] {
				continue
			}
			seen[dest] = true
			queue = append(queue, node{dest, first})
		}
	}
	return "", false
}

// stepToParty BFS-walks the room graph to the nearest crewmate not in this room
// and returns the first move toward them (and whose room it is).
func (w *World) stepToParty(p *Player) (dir, who string, ok bool) {
	if p.party == nil {
		return "", "", false
	}
	targets := map[string]string{}
	for _, m := range p.party.Members {
		if m == p || m.RoomID == p.RoomID {
			continue
		}
		if _, exists := targets[m.RoomID]; !exists {
			targets[m.RoomID] = m.Name
		}
	}
	if len(targets) == 0 {
		return "", "", false
	}
	type node struct{ id, first string }
	seen := map[string]bool{p.RoomID: true}
	queue := []node{{p.RoomID, ""}}
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
			if name, isTarget := targets[dest]; isTarget {
				return first, name, true
			}
			if seen[dest] {
				continue
			}
			seen[dest] = true
			queue = append(queue, node{dest, first})
		}
	}
	return "", "", false
}

// padRight pads s with spaces to width w (never truncates).
func padRight(s string, w int) string {
	for len(s) < w {
		s += " "
	}
	return s
}
