package cowboy

// One-way departure gates. Moving into a gated exit doesn't move you — it prompts
// a YES/NO confirmation, and only YES performs the (non-returnable) transition.
// Used for the Undercity finale: The Warm Pulse (z10_14) emerges UP into Parque
// Central at the heart of Noche City, sealing the depths behind you.

type departGate struct {
	dir    string // the direction (from the gated room) that triggers it
	to     string // destination room
	warn   string // confirmation prompt body
	leave  string // appended to the runner's name, broadcast in the room they leave
	arrive string // flavor line shown to the runner on arrival
}

var departGates = map[string]departGate{
	"z10_14": {
		dir:    "up",
		to:     "parque_central",
		warn:   "A warm thermal updraft breathes UP through a fissure in the welded sky. Climb it and you'll surface at the heart of the city — but The Warm Pulse seals behind you. To stand here again you'd have to descend the whole Undercity from the top.",
		leave:  " climbs into the updraft and is gone.",
		arrive: "You ride the warm updraft up through the dead-welded sky and break the surface into open air. Behind and below, The Warm Pulse seals shut.",
	},
}

// tryDepartGate intercepts a move into a one-way gated exit: it arms the pending
// confirmation and returns true (consuming the move). Returns false if dir isn't
// gated from the player's current room.
func (w *World) tryDepartGate(p *Player, dir string) bool {
	g, ok := departGates[p.RoomID]
	if !ok || g.dir != dir {
		return false
	}
	p.confirmExit = p.RoomID
	p.send(style(hot, g.warn) + crlf)
	p.send(style(gold, "Type YES to climb out, or NO to stay.") + crlf)
	return true
}

// confirmDepart completes a pending one-way departure (the YES branch).
func (w *World) confirmDepart(p *Player) {
	g, ok := departGates[p.confirmExit]
	if !ok || p.RoomID != p.confirmExit {
		p.send(style(dim, "There's nothing here to confirm.") + crlf)
		p.confirmExit = ""
		return
	}
	p.confirmExit = ""
	w.broadcast(p.RoomID, p, style(dim, p.Name+g.leave)+crlf)
	p.RoomID = g.to
	p.send(style(neon, g.arrive) + crlf)
	w.broadcast(p.RoomID, p, style(dim, p.Name+" surfaces into the light.")+crlf)
	w.lookText(p)
}

// cancelDepart aborts a pending one-way departure (the NO branch).
func (w *World) cancelDepart(p *Player) {
	if p.confirmExit == "" {
		p.send(style(dim, "Nothing to cancel.") + crlf)
		return
	}
	p.confirmExit = ""
	p.send(style(green, "You stay in the warmth a while longer.") + crlf)
}
