package cowboy

// The red-light strip along the Sprawlbelt: joytoys you can PAY for company. The
// transaction is fade-to-black (no explicit content) — you pay scrip, take an
// hour off the grind, and walk out "unwound": fully restored (HP + RAM) and
// shaking off the night. A tasteful, genre-standard cyberpunk beat.

const joytoyFee = 75

// joytoyRooms maps a red-light room to the joytoy working it.
var joytoyRooms = map[string]string{
	"ic_5": "Velvet, a joytoy at the Rolling Rose",
	"sb_5": "Lux, a joytoy at the Fortune Stall",
}

// payJoytoy handles PAY/HIRE in a red-light room.
func (w *World) payJoytoy(p *Player, arg string) {
	name, ok := joytoyRooms[p.RoomID]
	if !ok {
		p.send(style(dim, "There's no one here to pay for that.") + crlf)
		return
	}
	if p.Eddies < joytoyFee {
		p.send(style(dim, name+" looks you over. \"Come back when you've got €$"+itoa(joytoyFee)+", choom.\"") + crlf)
		return
	}
	p.Eddies -= joytoyFee
	p.HP = p.MaxHP
	p.RAM = maxRAM(p)
	p.send(style(neon, name+" takes your hand and the booth curtain falls...") + crlf)
	p.send(style(dim, "  (an hour you won't be billing anyone for)") + crlf)
	p.send(style(green, "You walk out unwound — fully restored ("+itoa(p.HP)+"/"+itoa(p.MaxHP)+" HP), the night's edge gone. (-€$"+itoa(joytoyFee)+")") + crlf)
}
