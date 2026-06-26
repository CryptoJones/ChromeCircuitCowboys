package cowboy

// A quick terminal hacking mini-game (#33): crack a hidden access code at any
// data terminal. It's a bounded high/low search — guess the code (1–100), the
// daemon says "higher"/"lower"; crack it within the tries for a scrip + XP
// payout. Burn all your tries and the trace costs you a little RAM. Numbers you
// type while a hack is live go to the daemon, not inventory quick-use.

type hackGame struct {
	secret int
	tries  int
	reward int
}

const hackTries = 7 // log2(100) ≈ 7 — always crackable with good guesses

// startHack launches the mini-game at a terminal.
func (w *World) startHack(p *Player) {
	if !w.atTerminal(p) {
		p.send(style(dim, "You need a data terminal to jack a system (any vendor, medic, or the Data Port).") + crlf)
		return
	}
	if p.hack != nil {
		p.send(style(dim, "You're already in a run — type a number to guess, or HACK again to abort.") + crlf)
		p.hack = nil
		p.send(style(dim, "Run aborted.") + crlf)
		return
	}
	p.hack = &hackGame{
		secret: w.roll(100) + 1, // 1..100
		tries:  hackTries,
		reward: 40 + p.Level*8, // scales with level
	}
	p.send(style(neon, "** ICE handshake open ** ") + style(green, "Crack the access code (1-100). "+itoa(hackTries)+" tries.") + crlf)
	p.send(style(dim, "Type a number to guess.") + crlf)
}

// hackGuess processes a numeric guess against the active hack.
func (w *World) hackGuess(p *Player, n int) {
	g := p.hack
	if g == nil {
		return
	}
	if n < 1 || n > 100 {
		p.send(style(dim, "The code is between 1 and 100.") + crlf)
		return
	}
	if n == g.secret {
		xp := g.reward / 4
		p.Eddies += g.reward
		p.XP += xp
		p.hack = nil
		p.send(style(gold, "*** CRACKED. The system rolls over. +€$"+itoa(g.reward)+", +"+itoa(xp)+"xp. ***") + crlf)
		w.checkLevelUp(p)
		return
	}
	g.tries--
	if g.tries <= 0 {
		p.hack = nil
		ram := maxRAM(p) / 4
		if ram > 0 {
			p.RAM -= ram
			if p.RAM < 0 {
				p.RAM = 0
			}
		}
		p.send(style(red, "*** TRACE LOCKED — the daemon boots you. The code was "+itoa(g.secret)+". (-"+itoa(ram)+" RAM) ***") + crlf)
		return
	}
	hint := "higher"
	if n > g.secret {
		hint = "lower"
	}
	p.send(style(dim, itoa(n)+" — go "+hint+". ("+itoa(g.tries)+" tries left)") + crlf)
}
