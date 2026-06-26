package cowboy

import "strings"

// Data terminals (every vendor/medic room, plus the Data Port) let runners reach
// each other remotely: SEND a message or WIRE scrip, delivered even if the
// recipient is offline.

// atTerminal reports whether the player is at a data terminal.
func (w *World) atTerminal(p *Player) bool {
	r := w.room(p.RoomID)
	return r != nil && (r.Term || r.Vendor || r.Medic)
}

// resolveName returns a recipient's canonical character name (online or saved),
// and whether they exist at all.
func (w *World) resolveName(name string) (string, bool) {
	if o := w.onlineByName(name); o != nil {
		return o.Name, true
	}
	if sp, ok, _ := w.store.Load(name); ok {
		return sp.Name, true
	}
	return "", false
}

// deliverMail pops and prints any queued messages for p.
func (w *World) deliverMail(p *Player) {
	msgs, _ := w.store.PopMail(p.Name)
	if len(msgs) == 0 {
		return
	}
	p.send(style(neon, "-- You have "+itoa(len(msgs))+" message(s) --") + crlf)
	for _, m := range msgs {
		p.send(style(gold, "  from "+m.From+": ") + style(green, m.Body) + crlf)
	}
}

// sendMail: SEND <player> <message> — queue a message at a terminal.
func (w *World) sendMail(p *Player, arg string) {
	if !w.atTerminal(p) {
		p.send(style(dim, "You need a data terminal for that (any vendor, medic, or the Data Port).") + crlf)
		return
	}
	fields := strings.Fields(arg)
	if len(fields) < 2 {
		p.send(style(dim, "SEND <runner> <message>") + crlf)
		return
	}
	to, ok := w.resolveName(fields[0])
	if !ok {
		p.send(style(dim, "No runner named \""+fields[0]+"\" on record.") + crlf)
		return
	}
	body := strings.TrimSpace(strings.TrimPrefix(arg, fields[0]))
	if err := w.store.PushMail(to, p.Name, body); err != nil {
		p.send(style(red, "The terminal rejects the message.") + crlf)
		return
	}
	p.send(style(green, "Message sent to "+to+".") + crlf)
	if o := w.onlineByName(to); o != nil {
		o.send(style(neon, "** A message arrives from "+p.Name+" — type MAIL. **") + crlf)
	}
}

// wireScrip: WIRE <player> <amount> — transfer scrip at a terminal (recipient
// credited even if offline).
func (w *World) wireScrip(p *Player, arg string) {
	if !w.atTerminal(p) {
		p.send(style(dim, "You need a data terminal to wire scrip (any vendor, medic, or the Data Port).") + crlf)
		return
	}
	fields := strings.Fields(arg)
	if len(fields) < 2 {
		p.send(style(dim, "WIRE <runner> <amount>") + crlf)
		return
	}
	amt := atoiSafe(fields[1])
	if amt <= 0 {
		p.send(style(dim, "Wire how much? (a positive amount)") + crlf)
		return
	}
	if p.Eddies < amt {
		p.send(style(dim, "You don't have €$"+itoa(amt)+" to wire.") + crlf)
		return
	}
	to, ok := w.resolveName(fields[0])
	if !ok {
		p.send(style(dim, "No runner named \""+fields[0]+"\" on record.") + crlf)
		return
	}
	if strings.EqualFold(to, p.Name) {
		p.send(style(dim, "Wiring scrip to yourself accomplishes nothing.") + crlf)
		return
	}
	p.Eddies -= amt
	if o := w.onlineByName(to); o != nil {
		o.Eddies += amt
		w.save(o)
		o.send(style(gold, "** "+p.Name+" wired you €$"+itoa(amt)+". **") + crlf)
	} else if sp, found, _ := w.store.Load(to); found {
		sp.Eddies += amt
		_ = w.store.Save(sp)
	}
	w.save(p)
	p.send(style(green, "Wired €$"+itoa(amt)+" to "+to+".") + crlf)
}

// readMail: MAIL — read and clear your messages.
func (w *World) readMail(p *Player) {
	msgs, _ := w.store.PopMail(p.Name)
	if len(msgs) == 0 {
		p.send(style(dim, "No messages.") + crlf)
		return
	}
	p.send(style(neon, "-- Messages ("+itoa(len(msgs))+") --") + crlf)
	for _, m := range msgs {
		p.send(style(gold, "  from "+m.From+": ") + style(green, m.Body) + crlf)
	}
}

// atoiSafe parses a non-negative int, returning 0 on garbage.
func atoiSafe(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}
