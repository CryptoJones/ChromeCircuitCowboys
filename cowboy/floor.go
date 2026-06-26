package cowboy

import (
	"sort"
	"strconv"
	"strings"
)

// Floor items: DROP leaves gear in the room for anyone (esp. a crewmate) to GET.
// In-memory like corpses; not persisted across a server restart.

// floorIn returns the item pile lying in a room (may be nil).
func (w *World) floorIn(roomID string) map[string]int { return w.floor[roomID] }

// floorList renders the room's floor pile for LOOK, or "" if empty.
func (w *World) floorList(roomID string) string {
	pile := w.floor[roomID]
	if len(pile) == 0 {
		return ""
	}
	names := make([]string, 0, len(pile))
	for n, q := range pile {
		if q > 0 {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	var parts []string
	for _, n := range names {
		parts = append(parts, itoa(pile[n])+"x "+n)
	}
	return strings.Join(parts, ", ")
}

// drop puts item(s) from the pack onto the room floor. DROP ALL drops everything.
func (w *World) drop(p *Player, arg string) {
	fields := strings.Fields(strings.ToLower(strings.TrimSpace(arg)))
	if len(fields) == 0 {
		p.send(style(dim, "Drop what? (DROP <item> [qty] / DROP ALL)") + crlf)
		return
	}
	if fields[0] == "all" {
		if len(p.Inv) == 0 {
			p.send(style(dim, "You've nothing to drop.") + crlf)
			return
		}
		n := 0
		for name, q := range p.Inv {
			if q > 0 {
				w.addFloor(p.RoomID, name, q)
				n += q
			}
		}
		p.Inv = map[string]int{}
		p.send(style(green, "You dump everything ("+itoa(n)+" items) on the floor.") + crlf)
		w.broadcast(p.RoomID, p, style(dim, p.Name+" dumps a pile of gear on the floor.")+crlf)
		return
	}
	name := fields[0]
	if p.Inv[name] <= 0 {
		p.send(style(dim, "You're not carrying any "+name+".") + crlf)
		return
	}
	qty := 1
	if len(fields) >= 2 {
		if q, err := strconv.Atoi(fields[1]); err == nil && q >= 1 {
			qty = q
		}
	}
	if qty > p.Inv[name] {
		qty = p.Inv[name]
	}
	p.Inv[name] -= qty
	if p.Inv[name] <= 0 {
		delete(p.Inv, name)
	}
	w.addFloor(p.RoomID, name, qty)
	p.send(style(green, "You drop "+itoa(qty)+"x "+name+".") + crlf)
	w.broadcast(p.RoomID, p, style(dim, p.Name+" drops "+itoa(qty)+"x "+name+".")+crlf)
}

// pickUp takes item(s) off the room floor into the pack. GET ALL sweeps it all.
func (w *World) pickUp(p *Player, arg string) {
	pile := w.floor[p.RoomID]
	if len(pile) == 0 {
		p.send(style(dim, "There's nothing on the floor here.") + crlf)
		return
	}
	fields := strings.Fields(strings.ToLower(strings.TrimSpace(arg)))
	if len(fields) == 0 || fields[0] == "all" {
		n := 0
		for name, q := range pile {
			if q > 0 {
				p.Inv[name] += q
				n += q
			}
		}
		delete(w.floor, p.RoomID)
		p.send(style(green, "You scoop up everything ("+itoa(n)+" items).") + crlf)
		w.broadcast(p.RoomID, p, style(dim, p.Name+" scoops up the floor pile.")+crlf)
		return
	}
	name := fields[0]
	if pile[name] <= 0 {
		p.send(style(dim, "There's no "+name+" on the floor.") + crlf)
		return
	}
	qty := pile[name]
	if len(fields) >= 2 {
		if q, err := strconv.Atoi(fields[1]); err == nil && q >= 1 && q < qty {
			qty = q
		}
	}
	pile[name] -= qty
	if pile[name] <= 0 {
		delete(pile, name)
	}
	if len(pile) == 0 {
		delete(w.floor, p.RoomID)
	}
	p.Inv[name] += qty
	p.send(style(green, "You pick up "+itoa(qty)+"x "+name+".") + crlf)
	w.broadcast(p.RoomID, p, style(dim, p.Name+" picks up "+itoa(qty)+"x "+name+".")+crlf)
}

// addFloor adds qty of an item to a room's floor pile.
func (w *World) addFloor(roomID, name string, qty int) {
	if w.floor[roomID] == nil {
		w.floor[roomID] = map[string]int{}
	}
	w.floor[roomID][name] += qty
}
