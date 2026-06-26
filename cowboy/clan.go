package cowboy

import "strings"

// Clans are a simple persistent name-tag: runners who share a clan and party up
// together earn bonus rewards. A party of 2+ in a room gets a party bonus on
// XP/scrip; if 2+ of them share a clan, that bonus is multiplied by 1.8 (#44).

const (
	partyBonusPct = 115 // +15% for grouping
	clanBonusPct  = 180 // ×1.8 when clanmates party together
)

// RewardPct exposes rewardPct for tooling/tests.
func (w *World) RewardPct(killer *Player) int { return w.rewardPct(killer) }

// rewardPct returns the kill-reward multiplier (percent; 100 = no bonus) for a
// killer based on their in-room crew and shared clan membership.
func (w *World) rewardPct(killer *Player) int {
	if killer.party == nil {
		return 100
	}
	var inRoom []*Player
	for _, m := range killer.party.Members {
		if m.RoomID == killer.RoomID {
			inRoom = append(inRoom, m)
		}
	}
	if len(inRoom) < 2 {
		return 100
	}
	pct := partyBonusPct
	clans := map[string]int{}
	for _, m := range inRoom {
		if m.Clan != "" {
			clans[strings.ToLower(m.Clan)]++
		}
	}
	for _, n := range clans {
		if n >= 2 { // two+ crewmates of the same clan present
			pct = pct * clanBonusPct / 100
			break
		}
	}
	return pct
}

// clan handles the CLAN command: create/join/leave/list, or show your clan.
func (w *World) clan(p *Player, arg string) {
	fields := strings.Fields(strings.TrimSpace(arg))
	if len(fields) == 0 {
		if p.Clan == "" {
			p.send(style(dim, "You're not in a clan. CLAN CREATE <name> or CLAN JOIN <name>.") + crlf)
			return
		}
		p.send(style(neon, "Clan: ") + style(gold, p.Clan) + crlf)
		var mates []string
		for _, o := range w.players {
			if o != p && strings.EqualFold(o.Clan, p.Clan) {
				mates = append(mates, o.Name)
			}
		}
		if len(mates) > 0 {
			p.send(style(dim, "  online clanmates: "+strings.Join(mates, ", ")) + crlf)
		} else {
			p.send(style(dim, "  no clanmates online.") + crlf)
		}
		return
	}
	sub := strings.ToLower(fields[0])
	name := strings.TrimSpace(strings.TrimPrefix(arg, fields[0]))
	switch sub {
	case "create", "join":
		if p.Clan != "" {
			p.send(style(dim, "You're already in clan "+p.Clan+". CLAN LEAVE first.") + crlf)
			return
		}
		if name == "" {
			p.send(style(dim, "CLAN "+sub+" <name>") + crlf)
			return
		}
		if len(name) > 24 {
			name = name[:24]
		}
		p.Clan = name
		w.save(p)
		p.send(style(green, "You're now flying the "+name+" colors. Party with clanmates for bonus rewards.") + crlf)
	case "leave":
		if p.Clan == "" {
			p.send(style(dim, "You're not in a clan.") + crlf)
			return
		}
		old := p.Clan
		p.Clan = ""
		w.save(p)
		p.send(style(green, "You've left clan "+old+".") + crlf)
	case "list":
		seen := map[string]bool{}
		var names []string
		for _, o := range w.players {
			if o.Clan != "" && !seen[strings.ToLower(o.Clan)] {
				seen[strings.ToLower(o.Clan)] = true
				names = append(names, o.Clan)
			}
		}
		if len(names) == 0 {
			p.send(style(dim, "No clans active right now.") + crlf)
			return
		}
		p.send(style(neon, "Active clans: ") + strings.Join(names, ", ") + crlf)
	default:
		p.send(style(dim, "CLAN [create|join|leave|list] <name>") + crlf)
	}
}
