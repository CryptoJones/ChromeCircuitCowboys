package cowboy

import "strings"

// Face-to-face player trading: a two-sided, confirm-locked swap of items and
// scrip between two runners in the same room. Nothing changes hands until BOTH
// confirm, and any change to an offer clears both confirmations — so neither
// side can confirm and then sneak the deal.

type tradeSession struct {
	a, b           *Player
	aItems, bItems map[string]int
	aScrip, bScrip int
	aOK, bOK       bool
}

// side returns the caller's offer (items, scrip ptr, ok ptr) and the other party.
func (ts *tradeSession) side(p *Player) (items map[string]int, scrip *int, ok *bool, other *Player) {
	if p == ts.a {
		return ts.aItems, &ts.aScrip, &ts.aOK, ts.b
	}
	return ts.bItems, &ts.bScrip, &ts.bOK, ts.a
}

// trade opens a trade with a named runner in the room, or (no arg) shows the
// current offer.
func (w *World) trade(p *Player, arg string) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		w.showTrade(p)
		return
	}
	if p.trade != nil {
		p.send(style(dim, "You're already trading. CONFIRM, or CANCEL first.") + crlf)
		return
	}
	target := w.playerInRoomByName(p.RoomID, arg, p)
	if target == nil {
		p.send(style(dim, "There's no one here by that name to trade with.") + crlf)
		return
	}
	if target.trade != nil {
		p.send(style(dim, target.Name+" is already in a trade.") + crlf)
		return
	}
	ts := &tradeSession{a: p, b: target, aItems: map[string]int{}, bItems: map[string]int{}}
	p.trade, target.trade = ts, ts
	p.send(style(green, "Trade opened with "+target.Name+".") + style(dim, " OFFER <item> [n] / OFFER scrip <n>, CONFIRM, CANCEL.") + crlf)
	target.send(style(green, p.Name+" opens a trade with you.") + style(dim, " OFFER ..., CONFIRM, CANCEL.") + crlf)
}

// tradeOffer adds an item or scrip to the caller's side of the open trade.
func (w *World) tradeOffer(p *Player, arg string) {
	ts := p.trade
	if ts == nil {
		p.send(style(dim, "You're not trading. TRADE <runner> first.") + crlf)
		return
	}
	items, scrip, _, other := ts.side(p)
	fields := strings.Fields(strings.ToLower(strings.TrimSpace(arg)))
	if len(fields) == 0 {
		p.send(style(dim, "OFFER <item> [n]  or  OFFER scrip <n>") + crlf)
		return
	}
	if fields[0] == "scrip" || fields[0] == "eddies" || fields[0] == "credits" {
		amt := 1
		if len(fields) >= 2 {
			amt = atoiSafe(fields[1])
		}
		if amt <= 0 || p.Eddies < amt {
			p.send(style(dim, "You don't have €$"+itoa(amt)+" to offer.") + crlf)
			return
		}
		*scrip = amt
	} else {
		name := fields[0]
		n := 1
		if len(fields) >= 2 {
			n = atoiSafe(fields[1])
		}
		if n <= 0 {
			n = 1
		}
		if p.Inv[name] < n {
			p.send(style(dim, "You're not carrying "+itoa(n)+"x "+name+".") + crlf)
			return
		}
		items[name] = n
	}
	ts.aOK, ts.bOK = false, false // any change re-opens both confirmations
	p.send(style(green, "Offer updated.") + crlf)
	other.send(style(neon, p.Name+" changed their offer — review and CONFIRM.") + crlf)
	w.showTrade(p)
}

// showTrade prints both sides of the open offer.
func (w *World) showTrade(p *Player) {
	ts := p.trade
	if ts == nil {
		p.send(style(dim, "You're not trading.") + crlf)
		return
	}
	mine, myScrip, myOK, other := ts.side(p)
	_, theirScrip, theirOK, _ := ts.side(other)
	theirs := ts.bItems
	if other == ts.a {
		theirs = ts.aItems
	}
	p.send(style(neon, "-- Trade with "+other.Name+" --") + crlf)
	p.send(style(gold, "  You offer:   ") + offerLine(mine, *myScrip) + tick(*myOK) + crlf)
	p.send(style(gold, "  They offer:  ") + offerLine(theirs, *theirScrip) + tick(*theirOK) + crlf)
	p.send(style(dim, "  CONFIRM when ready, or CANCEL.") + crlf)
}

func offerLine(items map[string]int, scrip int) string {
	var parts []string
	for name, q := range items {
		parts = append(parts, itoa(q)+"x "+name)
	}
	if scrip > 0 {
		parts = append(parts, "€$"+itoa(scrip))
	}
	if len(parts) == 0 {
		return "(nothing)"
	}
	return strings.Join(parts, ", ")
}

func tick(ok bool) string {
	if ok {
		return style(green, "  [CONFIRMED]")
	}
	return ""
}

// tradeConfirm marks the caller's side ready and executes when both agree.
func (w *World) tradeConfirm(p *Player) {
	ts := p.trade
	if ts == nil {
		p.send(style(dim, "You're not trading.") + crlf)
		return
	}
	if ts.a.RoomID != ts.b.RoomID {
		w.cancelTrade(ts, "the other runner left")
		return
	}
	_, _, myOK, other := ts.side(p)
	*myOK = true
	p.send(style(green, "You confirm the trade.") + crlf)
	other.send(style(neon, p.Name+" confirmed — CONFIRM to seal it.") + crlf)
	if ts.aOK && ts.bOK {
		w.execTrade(ts)
	}
}

// execTrade validates and performs the atomic swap, then closes the session.
func (w *World) execTrade(ts *tradeSession) {
	if !hasAll(ts.a, ts.aItems) || ts.a.Eddies < ts.aScrip ||
		!hasAll(ts.b, ts.bItems) || ts.b.Eddies < ts.bScrip {
		w.cancelTrade(ts, "someone no longer has what they offered")
		return
	}
	moveItems(ts.a, ts.b, ts.aItems)
	moveItems(ts.b, ts.a, ts.bItems)
	ts.a.Eddies += ts.bScrip - ts.aScrip
	ts.b.Eddies += ts.aScrip - ts.bScrip
	ts.a.trade, ts.b.trade = nil, nil
	w.save(ts.a)
	w.save(ts.b)
	ts.a.send(style(gold, "*** Trade complete. ***") + crlf)
	ts.b.send(style(gold, "*** Trade complete. ***") + crlf)
}

// tradeCancel tears down the caller's trade.
func (w *World) tradeCancel(p *Player) {
	if p.trade == nil {
		p.send(style(dim, "You're not trading.") + crlf)
		return
	}
	w.cancelTrade(p.trade, p.Name+" called it off")
}

func (w *World) cancelTrade(ts *tradeSession, why string) {
	a, b := ts.a, ts.b
	if a != nil {
		a.trade = nil
		a.send(style(dim, "Trade cancelled — "+why+".") + crlf)
	}
	if b != nil {
		b.trade = nil
		b.send(style(dim, "Trade cancelled — "+why+".") + crlf)
	}
}

// hasAll reports whether p still carries every offered item in the given qty.
func hasAll(p *Player, items map[string]int) bool {
	for name, q := range items {
		if p.Inv[name] < q {
			return false
		}
	}
	return true
}

// moveItems transfers offered items from one runner's pack to another's.
func moveItems(from, to *Player, items map[string]int) {
	for name, q := range items {
		from.Inv[name] -= q
		if from.Inv[name] <= 0 {
			delete(from.Inv, name)
		}
		to.Inv[name] += q
	}
}
