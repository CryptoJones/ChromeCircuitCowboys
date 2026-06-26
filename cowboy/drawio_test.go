package cowboy

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestGenLevelDrawio regenerates the draw.io level diagrams from the live zone
// data, so docs/*.drawio never drift from the authored world. It only writes
// when GEN_ROOM_MAP=1 (same gate as the SVG/markdown maps); otherwise it's a
// no-op guard that keeps the builders compiling.
func TestGenLevelDrawio(t *testing.T) {
	ugRooms, ugMobs := buildUndergroundZones()
	netRooms, netMobs := buildNetZones()

	ugByID := map[string]*Room{}
	for _, r := range ugRooms {
		ugByID[r.ID] = r
	}
	ugByHome := map[string]*MobTemplate{}
	for _, m := range ugMobs {
		ugByHome[m.Home] = m
	}
	netByID := map[string]*Room{}
	for _, r := range netRooms {
		netByID[r.ID] = r
	}
	netByHome := map[string]*MobTemplate{}
	for _, m := range netMobs {
		netByHome[m.Home] = m
	}

	ringByID := map[string]*Room{}
	for _, r := range buildRingRooms() {
		ringByID[r.ID] = r
	}

	// Always build (compile + smoke), only persist under the env gate.
	ug := buildUndergroundDrawio(ugByID, ugByHome)
	net := buildNetDrawio(netByID, netByHome)
	rings := buildRingDrawio(ringByID)
	for name, xml := range map[string]string{"underground": ug, "net": net, "rings": rings} {
		if !strings.Contains(xml, "<mxfile") || !strings.Contains(xml, "</mxfile>") {
			t.Fatalf("%s drawio is not a well-formed mxfile", name)
		}
	}

	if os.Getenv("GEN_ROOM_MAP") == "" {
		return
	}
	for path, xml := range map[string]string{
		"../docs/underground-descent.drawio": ug,
		"../docs/net-ascent.drawio":          net,
		"../docs/rp-rings.drawio":            rings,
	} {
		if err := os.WriteFile(path, []byte(xml), 0644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
		t.Logf("wrote %s", path)
	}
}

// --- shared drawio helpers ---------------------------------------------------

const (
	dioBG     = "#07090f"
	dioCyan   = "#27d4ff"
	dioGreen  = "#55ff99"
	dioAmber  = "#ffb000"
	dioRed    = "#ff4f4f"
	dioOrange = "#ff8a3d"
	dioGrey   = "#5a6678"
	dioInk    = "#0e1320"
)

func dioHeader(b *strings.Builder, name, id string, w, h int) {
	fmt.Fprintf(b, `<mxfile host="app.diagrams.net" type="device">`+"\n")
	fmt.Fprintf(b, `  <diagram name="%s" id="%s">`+"\n", xmlEsc(name), id)
	fmt.Fprintf(b, `    <mxGraphModel dx="900" dy="700" grid="0" guides="1" tooltips="1" connect="1" arrows="1" fold="1" page="1" pageScale="1" pageWidth="%d" pageHeight="%d" math="0" background="%s" shadow="0">`+"\n", w, h, dioBG)
	b.WriteString("      <root>\n        <mxCell id=\"0\" />\n        <mxCell id=\"1\" parent=\"0\" />\n")
}

func dioFooter(b *strings.Builder) {
	b.WriteString("      </root>\n    </mxGraphModel>\n  </diagram>\n</mxfile>\n")
}

func dioText(b *strings.Builder, id, html string, x, y, w, h, size int, color string, bold bool) {
	fs := ""
	if bold {
		fs = "fontStyle=1;"
	}
	fmt.Fprintf(b, `        <mxCell id="%s" value="%s" style="text;html=1;fontColor=%s;fontSize=%d;%sfontFamily=Menlo;align=left;verticalAlign=top;whiteSpace=wrap;" vertex="1" parent="1"><mxGeometry x="%d" y="%d" width="%d" height="%d" as="geometry"/></mxCell>`+"\n",
		id, xmlEsc(html), color, size, fs, x, y, w, h)
}

func dioBox(b *strings.Builder, id, html string, x, y, w, h int, fill, stroke string) {
	fmt.Fprintf(b, `        <mxCell id="%s" value="%s" style="rounded=1;arcSize=10;html=1;whiteSpace=wrap;fillColor=%s;strokeColor=%s;fontColor=%s;fontSize=12;fontFamily=Menlo;align=left;verticalAlign=top;spacingLeft=6;spacingTop=4;" vertex="1" parent="1"><mxGeometry x="%d" y="%d" width="%d" height="%d" as="geometry"/></mxCell>`+"\n",
		id, xmlEsc(html), fill, stroke, stroke, x, y, w, h)
}

func dioEdge(b *strings.Builder, id, src, dst, label, color string, dashed bool) {
	dash := ""
	if dashed {
		dash = "dashed=1;"
	}
	fmt.Fprintf(b, `        <mxCell id="%s" value="%s" style="edgeStyle=orthogonalEdgeStyle;rounded=0;html=1;endArrow=block;strokeColor=%s;%sfontColor=%s;fontFamily=Menlo;fontSize=11;" edge="1" parent="1" source="%s" target="%s"><mxGeometry relative="1" as="geometry"/></mxCell>`+"\n",
		id, xmlEsc(label), color, dash, color, src, dst)
}

// roomKindColors picks the box fill/stroke from the mob kind + room flags,
// matching the SVG legend (boss=red, elite=orange, vendor/EM/safe=green/amber).
func roomKindColors(r *Room, kind string) (fill, stroke string) {
	switch {
	case kind == "b":
		return "#1a0f0f", dioRed
	case r.Safe && (r.Vendor || r.Medic):
		return "#0c1810", dioGreen
	case r.Vendor || r.Medic:
		return dioInk, dioAmber
	case r.Safe:
		return dioInk, dioGreen
	case kind == "e":
		return dioInk, dioOrange
	default:
		return dioInk, dioCyan
	}
}

// --- underground descent -----------------------------------------------------

func buildUndergroundDrawio(byID map[string]*Room, byHome map[string]*MobTemplate) string {
	const (
		x0       = 40
		colW     = 360
		boxW     = 316
		boxH     = 66
		vgap     = 34
		hdrY     = 92
		rowStart = 124
	)
	maxRows := 0
	for _, z := range undergroundZoneData {
		if len(z.areas) > maxRows {
			maxRows = len(z.areas)
		}
	}
	width := x0*2 + len(undergroundZoneData)*colW
	height := rowStart + maxRows*(boxH+vgap) + 96

	var b strings.Builder
	dioHeader(&b, "C3 Underground Descent (L1-99)", "ccc-underground", width, height)
	dioText(&b, "title", "CHROME CIRCUIT COWBOYS — UNDERGROUND DESCENT (L1-99)", 40, 24, width-80, 34, 26, dioCyan, true)
	dioText(&b, "subtitle", "Each column = one arc, rooms in descent order. Letters on edges = the actual exit direction. Entry: Back Alley —D→ z1_01.", 40, 60, width-80, 22, 13, dioGrey, false)

	for zi, z := range undergroundZoneData {
		colX := x0 + zi*colW
		lo, hi := z.band*10-9, z.band*10
		if z.band == 10 {
			hi = 99
		}
		dioText(&b, fmt.Sprintf("zh%d", zi), fmt.Sprintf("L%d-%d · %s", lo, hi, trunc(z.name, 24)), colX, hdrY, colW-30, 22, 15, dioCyan, true)
		for ai, ad := range z.areas {
			r := byID[ad.id]
			boxY := rowStart + ai*(boxH+vgap)
			kind, mobName := "", ""
			if ad.mob != "" {
				kind, mobName = splitMob(ad.mob)
			}
			fill, stroke := roomKindColors(r, kind)
			var marks []string
			switch kind {
			case "b":
				marks = append(marks, "BOSS: "+trunc(mobName, 22))
			case "e":
				marks = append(marks, "elite: "+trunc(mobName, 20))
			case "c":
				marks = append(marks, "foe: "+trunc(mobName, 22))
			}
			if r.Vendor {
				marks = append(marks, "vendor")
			}
			if r.Medic {
				marks = append(marks, "EM")
			}
			if r.Safe {
				marks = append(marks, "safe")
			}
			if ad.cache == "up" {
				marks = append(marks, "[^ loot cache]")
			} else if ad.cache == "down" {
				marks = append(marks, "[v loot cache]")
			}
			html := fmt.Sprintf("<b>%s</b>  %s<br><font color='#9fb0c4'>%s</font>",
				r.ID, trunc(r.Name, 26), trunc(strings.Join(marks, " · "), 48))
			dioBox(&b, r.ID, html, colX, boxY, boxW, boxH, fill, stroke)

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
				dioEdge(&b, "e_"+r.ID+"_"+nextID, r.ID, nextID, fwd, dioGrey, false)
			}
		}
		// Inter-arc descent edge: this arc's last room down into the next arc's first.
		if zi+1 < len(undergroundZoneData) {
			last := z.areas[len(z.areas)-1].id
			next := undergroundZoneData[zi+1].areas[0].id
			dioEdge(&b, "descend_"+last+"_"+next, last, next, "D ▼ descend", dioAmber, true)
		}
	}

	dioText(&b, "legend", "Legend:  <font color='#55ff99'>green</font>=safe/hub  <font color='#ffb000'>amber</font>=vendor/EM  <font color='#27d4ff'>cyan</font>=combat  <font color='#ff8a3d'>orange</font>=elite  <font color='#ff4f4f'>red</font>=arc boss  ·  N/S/E/W/U/D = exit dir  ·  amber edge = descent to the next arc",
		40, height-30, width-80, 22, 13, "#9fb0c4", false)
	dioFooter(&b)
	return b.String()
}

// --- net ascent --------------------------------------------------------------

func buildNetDrawio(byID map[string]*Room, byHome map[string]*MobTemplate) string {
	const (
		x0       = 40
		colW     = 320
		boxW     = 280
		boxH     = 84
		vgap     = 30
		hdrY     = 92
		rowStart = 124
	)
	foe := func(id string) string {
		if mt := byHome[id]; mt != nil {
			if strings.HasSuffix(mt.ID, "_c") {
				return "data-cache (RAM+scrip)"
			}
			return mt.Name
		}
		return "—"
	}
	maxRows := 0
	for _, z := range netZoneData {
		if len(z.areas) > maxRows {
			maxRows = len(z.areas)
		}
	}
	width := x0*2 + len(netZoneData)*colW
	height := rowStart + maxRows*(boxH+vgap) + 96

	var b strings.Builder
	dioHeader(&b, "C3 Net Ascent (L1-99)", "ccc-net", width, height)
	dioText(&b, "title", "CHROME CIRCUIT COWBOYS — THE NET ASCENT (L1-99)", 40, 24, width-80, 34, 26, dioCyan, true)
	dioText(&b, "subtitle", "Each column = one Net zone; each box = an area's 3-layer stack (Shell → Breach → Core). Jack in: Data Port —U→ nz1_1_top. Amber edge = ascend to the next zone.", 40, 60, width-80, 22, 13, dioGrey, false)

	for zi, z := range netZoneData {
		colX := x0 + zi*colW
		lo, hi := z.band*10-9, z.band*10
		if z.band == 10 {
			hi = 99
		}
		dioText(&b, fmt.Sprintf("nzh%d", zi), fmt.Sprintf("L%d-%d · %s", lo, hi, trunc(z.name, 22)), colX, hdrY, colW-30, 22, 15, dioCyan, true)
		for ai, ar := range z.areas {
			base := fmt.Sprintf("%s_%d", z.key, ai+1)
			top := byID[base+"_top"]
			boxY := rowStart + ai*(boxH+vgap)
			// Colour by the strongest thing in the stack: boss in Core => red.
			kind := "c"
			if bm := byHome[base+"_bot_m"]; bm != nil && bm.HP > 80 {
				kind = "b"
			}
			stack := &Room{Safe: top != nil && top.Safe, Vendor: top != nil && top.Vendor}
			fill, stroke := roomKindColors(stack, kind)
			html := fmt.Sprintf("<b>%s</b><br><font color='#9fb0c4'>Shell: %s<br>Breach: %s<br>Core: %s</font>",
				trunc(ar.name, 26),
				trunc(foe(base+"_top_m"), 22),
				trunc(foe(base+"_mid_m"), 22),
				trunc(foe(base+"_bot_m"), 22))
			dioBox(&b, base, html, colX, boxY, boxW, boxH, fill, stroke)
			if ai+1 < len(z.areas) {
				next := fmt.Sprintf("%s_%d", z.key, ai+2)
				dioEdge(&b, "e_"+base+"_"+next, base, next, "", dioGrey, false)
			}
		}
		if zi+1 < len(netZoneData) {
			last := fmt.Sprintf("%s_%d", z.key, len(z.areas))
			next := fmt.Sprintf("%s_1", netZoneData[zi+1].key)
			dioEdge(&b, "ascend_"+last+"_"+next, last, next, "▲ ascend", dioAmber, true)
		}
	}

	dioText(&b, "legend", "Legend:  <font color='#27d4ff'>cyan</font>=area  <font color='#ff4f4f'>red</font>=boss-Core  <font color='#55ff99'>green</font>=safe/hub  ·  Shell/Breach/Core = the 3 layers  ·  amber edge = ascend to the next zone",
		40, height-30, width-80, 22, 13, "#9fb0c4", false)
	dioFooter(&b)
	return b.String()
}

// --- RP transit rings --------------------------------------------------------

// ringColors picks a box fill/stroke from a ring room's flags.
func ringColors(r *Room) (fill, stroke string) {
	switch {
	case r.Spar:
		return "#1a0f0f", dioRed
	case r.Vendor && r.Medic:
		return "#0c1810", dioGreen
	case r.Vendor:
		return dioInk, dioAmber
	case r.Medic:
		return dioInk, dioGreen
	case r.Safe:
		return dioInk, dioCyan
	default:
		return dioInk, dioGrey
	}
}

// ringFlags renders a ring room's flags as a short tag string.
func ringFlags(r *Room) string {
	var f []string
	if r.Spar {
		f = append(f, "sparring gym")
	}
	if r.Safe {
		f = append(f, "safe")
	}
	if r.Vendor {
		f = append(f, "vendor")
	}
	if r.Medic {
		f = append(f, "medic")
	}
	return strings.Join(f, " · ")
}

// buildRingDrawio renders the carved RP transit rings (Inner Circuit + the outer
// Sprawlbelt, plus the Iron Temple sparring gym) straight from buildRingRooms,
// so the ring map stays in sync with the world. Boxes note off-ring exits.
func buildRingDrawio(byID map[string]*Room) string {
	const (
		x0    = 40
		colW  = 178
		boxW  = 152
		boxH  = 60
		rowIC = 122
		rowSB = 248
		rowGym = 374
	)
	var ic, sb []string
	for i := 1; i <= 6; i++ {
		ic = append(ic, "ic_"+itoa(i))
	}
	for i := 1; i <= 10; i++ {
		sb = append(sb, "sb_"+itoa(i))
	}
	ringSet := map[string]bool{"sb_gym": true}
	for _, id := range ic {
		ringSet[id] = true
	}
	for _, id := range sb {
		ringSet[id] = true
	}

	width := x0*2 + 10*colW
	height := rowGym + boxH + 96

	var b strings.Builder
	dioHeader(&b, "C3 RP Transit Rings", "ccc-rings", width, height)
	dioText(&b, "title", "CHROME CIRCUIT COWBOYS — RP TRANSIT RINGS", 40, 24, width-80, 34, 26, dioCyan, true)
	dioText(&b, "subtitle", "Inner Circuit (express loop) + Sprawlbelt (outer beltway) off Neon Alley. Edges = exits, letters = direction; grey notes = off-ring exits.", 40, 60, width-80, 22, 13, dioGrey, false)
	dioText(&b, "ich", "Inner Circuit", x0, rowIC-26, 360, 20, 15, dioCyan, true)
	dioText(&b, "sbh", "Sprawlbelt", x0, rowSB-26, 360, 20, 15, dioCyan, true)

	pos := map[string][2]int{}
	for i, id := range ic {
		pos[id] = [2]int{x0 + i*colW, rowIC}
	}
	for i, id := range sb {
		pos[id] = [2]int{x0 + i*colW, rowSB}
	}
	pos["sb_gym"] = [2]int{pos["sb_9"][0], rowGym}

	order := append(append([]string{}, ic...), sb...)
	order = append(order, "sb_gym")
	for _, id := range order {
		r := byID[id]
		if r == nil {
			continue
		}
		p := pos[id]
		fill, stroke := ringColors(r)
		name := r.Name
		for _, pfx := range []string{"Inner Circuit :: ", "Sprawlbelt :: "} {
			name = strings.TrimPrefix(name, pfx)
		}
		var ext []string
		for _, d := range []string{"north", "south", "east", "west", "up", "down", "in", "out"} {
			dest, ok := r.Exits[d]
			if !ok || ringSet[dest] {
				continue
			}
			letter := dirLetter[d]
			if letter == "" {
				letter = strings.ToUpper(d[:1])
			}
			ext = append(ext, letter+"→"+dest)
		}
		val := "<b>" + id + "</b> " + trunc(name, 18)
		if f := ringFlags(r); f != "" {
			val += "<br><font color='#9fb0c4'>" + f + "</font>"
		}
		if len(ext) > 0 {
			val += "<br><font color='#5a6678'>" + strings.Join(ext, "  ") + "</font>"
		}
		dioBox(&b, id, val, p[0], p[1], boxW, boxH, fill, stroke)
	}

	seen := map[string]bool{}
	for _, id := range order {
		r := byID[id]
		if r == nil {
			continue
		}
		for _, d := range []string{"north", "south", "east", "west", "up", "down"} {
			dest, ok := r.Exits[d]
			if !ok || !ringSet[dest] {
				continue
			}
			if seen[id+"|"+dest] || seen[dest+"|"+id] {
				continue
			}
			seen[id+"|"+dest] = true
			dioEdge(&b, "e_"+id+"_"+dest, id, dest, dirLetter[d], dioGrey, false)
		}
	}

	dioText(&b, "legend", "Legend:  <font color='#27d4ff'>cyan</font>=safe  <font color='#ffb000'>amber</font>=vendor  <font color='#55ff99'>green</font>=medic  <font color='#ff4f4f'>red</font>=sparring gym  ·  N/S/E/W = exit dir  ·  grey = off-ring exit",
		40, height-30, width-80, 22, 13, "#9fb0c4", false)
	dioFooter(&b)
	return b.String()
}
