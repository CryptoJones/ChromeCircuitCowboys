package cowboy

import (
	"strconv"
	"strings"
)

// classBlocked returns a player-facing reason if this ware is restricted to a
// class the player isn't, or "" if they may buy/use it. The reason names the
// requirement (e.g. "This is Hacker gear — your Enforcer can't use it.").
func classBlocked(p *Player, x ware) string {
	if x.forClass == "" || strings.EqualFold(x.forClass, p.Class) {
		return ""
	}
	mine := p.Class
	if mine == "" {
		mine = "build"
	}
	return "This is " + strings.Title(x.forClass) + " gear — your " + strings.Title(mine) + " can't use it."
}

func (w *World) atVendor(p *Player) bool {
	r := w.room(p.RoomID)
	return r != nil && r.Vendor
}

// carryCap is how many items a runner can carry. It scales with LEVEL (not a
// stat), so it's class-neutral and fair to every build. Overflow lives in your
// Re-Clone Bay stash.
func carryCap(p *Player) int { return 10 + 2*p.Level }

// invCount is the total number of items carried (counts toward carryCap).
func invCount(p *Player) int {
	n := 0
	for _, q := range p.Inv {
		n += q
	}
	return n
}

func (w *World) list(p *Player) {
	if !w.atVendor(p) {
		p.send(style(dim, "There's no vendor here.") + crlf)
		return
	}
	wares := waresForRoom(p.RoomID)
	p.send(style(neon, "-- Vendor wares (BUY <#> or <item>) --") + crlf)
	for i, x := range wares {
		p.send("  " + style(gold, itoa(i+1)+")") + " " + style(gold, "€$"+itoa(x.price)) + "  " + x.name + style(dim, " — "+x.desc) + crlf)
	}
	p.send(style(dim, "You have €$"+itoa(p.Eddies)+".") + crlf)
}

func (w *World) buy(p *Player, arg string) {
	if !w.atVendor(p) {
		p.send(style(dim, "There's no vendor here.") + crlf)
		return
	}
	// Syntax: BUY <#|name> [qty]  — e.g. "buy 3 4" = four of item 3. Default qty 1.
	fields := strings.Fields(arg)
	if len(fields) == 0 {
		p.send(style(dim, "Buy what? Type LIST.") + crlf)
		return
	}
	wares := waresForRoom(p.RoomID)
	var x ware
	var ok bool
	if n, err := strconv.Atoi(fields[0]); err == nil {
		if n >= 1 && n <= len(wares) {
			x, ok = wares[n-1], true
		}
	} else {
		nm := strings.ToLower(fields[0])
		for _, ww := range wares {
			if ww.name == nm {
				x, ok = ww, true
				break
			}
		}
	}
	if !ok {
		p.send(style(dim, "No such item. Type LIST.") + crlf)
		return
	}
	if reason := classBlocked(p, x); reason != "" {
		// Don't let them waste scrip on gear they can't use.
		p.send(style(dim, reason+" The vendor won't sell it to you.") + crlf)
		return
	}
	qty := 1
	if len(fields) >= 2 {
		q, err := strconv.Atoi(fields[1])
		if err != nil || q < 1 {
			p.send(style(dim, "Quantity must be a positive number.") + crlf)
			return
		}
		qty = q
	}

	// Weapons and cyberdecks are permanent one-time UPGRADES — quantity doesn't
	// apply; always a single purchase.
	if x.bonus > 0 {
		if x.bonus <= p.WeaponBonus {
			p.send(style(dim, "Your current weapon is already better.") + crlf)
			return
		}
		if p.Eddies < x.price {
			p.send(style(red, "Not enough scrip (need €$"+itoa(x.price)+").") + crlf)
			return
		}
		p.Eddies -= x.price
		p.WeaponBonus, p.WeaponName = x.bonus, x.name
		p.send(style(green, "You jack in the "+x.name+". Attack +"+itoa(x.bonus)+".") + crlf)
		return
	}
	if x.deck > 0 {
		if x.deck <= p.DeckBonus {
			p.send(style(dim, "Your current deck is already as good.") + crlf)
			return
		}
		if p.Eddies < x.price {
			p.send(style(red, "Not enough scrip (need €$"+itoa(x.price)+").") + crlf)
			return
		}
		p.Eddies -= x.price
		p.DeckBonus = x.deck
		p.RAM = maxRAM(p) // fresh deck boots with full RAM
		p.send(style(green, "You install the "+x.name+". Max RAM is now "+itoa(maxRAM(p))+".") + crlf)
		return
	}

	// Consumables: buy qty of them.
	cost := x.price * qty
	if p.Eddies < cost {
		p.send(style(red, "Not enough scrip ("+itoa(qty)+"x "+x.name+" = €$"+itoa(cost)+").") + crlf)
		return
	}
	if invCount(p)+qty > carryCap(p) {
		p.send(style(dim, "That won't fit: "+itoa(invCount(p))+"/"+itoa(carryCap(p))+" carried. STASH something at your Re-Clone Bay first.") + crlf)
		return
	}
	p.Eddies -= cost
	p.Inv[x.name] += qty
	p.send(style(green, "Bought "+itoa(qty)+"x "+x.name+" for €$"+itoa(cost)+". You have "+itoa(p.Inv[x.name])+".") + crlf)
}

// sellBuyback is the fraction of an item's catalog price a vendor pays for it.
const sellBuyback = 50 // percent

// sell offloads unwanted items from the pack for scrip at a vendor (the opposite
// of BUY). Pays a fraction of the catalog price; items with no price can't be
// sold. Syntax: SELL <item> [qty].
func (w *World) sell(p *Player, arg string) {
	if !w.atVendor(p) {
		p.send(style(dim, "There's no vendor here to sell to.") + crlf)
		return
	}
	fields := strings.Fields(strings.ToLower(strings.TrimSpace(arg)))
	if len(fields) == 0 {
		p.send(style(dim, "Sell what? (SELL <item> [qty])") + crlf)
		return
	}
	name := fields[0]
	if p.Inv[name] <= 0 {
		p.send(style(dim, "You're not carrying any "+name+".") + crlf)
		return
	}
	qty := 1
	if len(fields) >= 2 {
		q, err := strconv.Atoi(fields[1])
		if err != nil || q < 1 {
			p.send(style(dim, "Quantity must be a positive number.") + crlf)
			return
		}
		qty = q
	}
	if qty > p.Inv[name] {
		qty = p.Inv[name]
	}
	x, ok := findWare(name)
	if !ok || x.price <= 0 {
		p.send(style(dim, "No one here will pay for "+name+".") + crlf)
		return
	}
	unit := x.price * sellBuyback / 100
	if unit < 1 {
		unit = 1
	}
	total := unit * qty
	p.Inv[name] -= qty
	if p.Inv[name] <= 0 {
		delete(p.Inv, name)
	}
	p.Eddies += total
	p.send(style(green, "Sold "+itoa(qty)+"x "+name+" for €$"+itoa(total)+" ("+itoa(unit)+" each).") + crlf)
}

// consumeInv removes one of an item, deleting the key when it hits zero.
func (w *World) consumeInv(p *Player, name string) {
	p.Inv[name]--
	if p.Inv[name] <= 0 {
		delete(p.Inv, name)
	}
}

// install wires salvaged cyberware (looted from a body) back into your body —
// only at a Emergency Medic. Reuses the same upgrade rules as buying it new.
func (w *World) install(p *Player, arg string) {
	r := w.room(p.RoomID)
	if r == nil || !r.Medic {
		p.send(style(dim, "You need a Emergency Medic to wire in cyberware. (try the Night Market)") + crlf)
		return
	}
	name := strings.ToLower(strings.TrimSpace(arg))
	if name == "" {
		p.send(style(dim, "Install what? (salvaged cyberware sits in your INVENTORY)") + crlf)
		return
	}
	if p.Inv[name] <= 0 {
		p.send(style(dim, "You're not carrying "+name+" to install.") + crlf)
		return
	}
	x, ok := findWare(name)
	if !ok || (x.bonus <= 0 && x.deck <= 0 && !x.isImplant()) {
		p.send(style(dim, "That's not cyberware a Emergency Medic can install.") + crlf)
		return
	}
	if x.isImplant() {
		p.Body += x.body
		p.Reflexes += x.refl
		p.Intelligence += x.intel
		if x.body > 0 { // Body raises max HP — heal the new headroom
			old := p.MaxHP
			p.MaxHP = maxHPFor(p)
			p.HP += p.MaxHP - old
		}
		w.consumeInv(p, name)
		var parts []string
		if x.body > 0 {
			parts = append(parts, "Body +"+itoa(x.body))
		}
		if x.refl > 0 {
			parts = append(parts, "Reflexes +"+itoa(x.refl))
		}
		if x.intel > 0 {
			parts = append(parts, "Intelligence +"+itoa(x.intel))
		}
		w.save(p)
		p.send(style(green, "The Emergency Medic wires in the "+name+". "+strings.Join(parts, ", ")+".") + crlf)
		return
	}
	if x.bonus > 0 {
		if x.bonus <= p.WeaponBonus {
			p.send(style(dim, "Your current weapon is already better — keep the "+name+" or GIVE it away.") + crlf)
			return
		}
		p.WeaponBonus, p.WeaponName = x.bonus, x.name
		w.consumeInv(p, name)
		p.send(style(green, "The Emergency Medic wires in the "+name+". Attack +"+itoa(x.bonus)+".") + crlf)
		return
	}
	if x.deck <= p.DeckBonus {
		p.send(style(dim, "Your current deck is already as good — keep the "+name+" or GIVE it away.") + crlf)
		return
	}
	p.DeckBonus = x.deck
	p.RAM = maxRAM(p) // freshly installed deck boots with full RAM
	w.consumeInv(p, name)
	p.send(style(green, "The Emergency Medic installs the "+name+". Max RAM is now "+itoa(maxRAM(p))+".") + crlf)
}

func (w *World) use(p *Player, arg string) {
	name := strings.ToLower(strings.TrimSpace(arg))
	if name == "" {
		p.send(style(dim, "Use what? (see INVENTORY)") + crlf)
		return
	}
	if p.Inv[name] <= 0 {
		p.send(style(dim, "You don't have a "+name+".") + crlf)
		return
	}
	x, ok := findWare(name)
	if !ok {
		p.send(style(dim, "That's not something you can use.") + crlf)
		return
	}
	if reason := classBlocked(p, x); reason != "" {
		p.send(style(dim, reason) + crlf)
		return
	}
	if x.heal <= 0 && x.ram <= 0 {
		// Not a consumable — explain WHY, don't just refuse.
		if x.isImplant() || x.bonus > 0 || x.deck > 0 {
			p.send(style(dim, "That's cyberware — INSTALL it at a Emergency Medic, don't USE it.") + crlf)
		} else {
			p.send(style(dim, "The "+name+" does nothing on its own — it's not a consumable.") + crlf)
		}
		return
	}
	// Don't waste a single-use consumable when it would have no effect.
	if x.ram <= 0 && x.heal > 0 && p.HP >= p.MaxHP {
		p.send(style(dim, "Your HP is already full — save the "+name+".") + crlf)
		return
	}
	if x.heal <= 0 && x.ram > 0 && p.RAM >= maxRAM(p) {
		p.send(style(dim, "Your RAM is already full — save the "+name+".") + crlf)
		return
	}
	p.Inv[name]--
	if p.Inv[name] == 0 {
		delete(p.Inv, name)
	}
	if x.heal > 0 {
		p.HP += x.heal
		if p.HP > p.MaxHP {
			p.HP = p.MaxHP
		}
		p.send(style(green, "You slot the "+name+" — HP now "+itoa(p.HP)+"/"+itoa(p.MaxHP)+".") + crlf)
	}
	if x.ram > 0 {
		p.RAM += x.ram
		if p.RAM > maxRAM(p) {
			p.RAM = maxRAM(p)
		}
		p.send(style(green, "You burn the "+name+" — RAM now "+itoa(p.RAM)+"/"+itoa(maxRAM(p))+".") + crlf)
	}
}
