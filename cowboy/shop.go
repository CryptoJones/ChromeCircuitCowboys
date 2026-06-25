package cowboy

import (
	"strconv"
	"strings"
)

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
	p.send(style(neon, "-- Vendor wares (BUY <#> or <item>) --") + crlf)
	for i, x := range shopWares {
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
	var x ware
	var ok bool
	if n, err := strconv.Atoi(fields[0]); err == nil {
		if n >= 1 && n <= len(shopWares) {
			x, ok = shopWares[n-1], true
		}
	} else {
		x, ok = findWare(strings.ToLower(fields[0]))
	}
	if !ok {
		p.send(style(dim, "No such item. Type LIST.") + crlf)
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
	if !ok || (x.bonus <= 0 && x.deck <= 0) {
		p.send(style(dim, "That's not cyberware a Emergency Medic can install.") + crlf)
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
	if !ok || (x.heal <= 0 && x.ram <= 0) {
		p.send(style(dim, "You can't use that.") + crlf)
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
