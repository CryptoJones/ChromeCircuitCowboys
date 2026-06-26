package cowboy

import "testing"

// crewBot invites the first AI runner into a fresh human's crew and returns both.
func crewBot(t *testing.T, room string) (*World, *Player, *Player) {
	t.Helper()
	w := NewWorld(NewMemStore())
	w.EnableBots(1)
	bot := botsIn(w)[0]
	boss, _ := newTestPlayer(w, "Boss", room)
	w.invite(boss, bot.Name)
	return w, boss, bot
}

func TestBotAutoAcceptsAndWarpsToLeader(t *testing.T) {
	_, boss, bot := crewBot(t, "ic_1")
	if bot.party == nil || bot.party != boss.party {
		t.Fatalf("bot should share the boss's crew after GROUP <bot>")
	}
	if bot.RoomID != boss.RoomID {
		t.Errorf("bot should warp to the leader's room %q, is in %q", boss.RoomID, bot.RoomID)
	}
	if bot.partyInvite != nil {
		t.Error("bot should auto-join, not leave a pending invite")
	}
	if len(boss.party.Members) != 2 || boss.party.Leader != boss {
		t.Errorf("crew should be {Boss(leader), bot}, got %d members", len(boss.party.Members))
	}
}

func TestCrewedBotFollowsLeader(t *testing.T) {
	w, boss, bot := crewBot(t, "ic_1")
	origin := boss.RoomID
	dest := "sb_2" // any valid room; partyFollow just relocates members from origin
	boss.RoomID = dest
	w.partyFollow(boss, origin, dest, "south")
	if bot.RoomID != dest {
		t.Errorf("crewed bot should follow the leader to %q, stayed in %q", dest, bot.RoomID)
	}
}

func TestCrewedBotDoesNotWander(t *testing.T) {
	w, _, bot := crewBot(t, "ic_1")
	w.SetRoll(func(n int) int { return 0 }) // would force a wander/chatter if not crewed
	before := bot.RoomID
	for i := 0; i < 10; i++ {
		w.tickBots()
	}
	if bot.RoomID != before {
		t.Errorf("crewed bot wandered off to %q (should stay put and follow)", bot.RoomID)
	}
}

func TestLeavingFreesBotsAndDissolvesAllBotCrew(t *testing.T) {
	// 1 human + 1 bot: human leaves → crew < 2 → bot freed.
	w, boss, bot := crewBot(t, "ic_1")
	w.leaveParty(boss)
	if bot.party != nil {
		t.Errorf("bot should be freed when the crew drops below two")
	}

	// 1 human + 2 bots: human leaves → all-bot remainder dissolves, both freed.
	w2 := NewWorld(NewMemStore())
	w2.EnableBots(2)
	bots := botsIn(w2)
	boss2, _ := newTestPlayer(w2, "Boss2", "ic_1")
	w2.invite(boss2, bots[0].Name)
	w2.invite(boss2, bots[1].Name)
	if len(boss2.party.Members) != 3 {
		t.Fatalf("expected a 3-member crew, got %d", len(boss2.party.Members))
	}
	w2.leaveParty(boss2)
	for _, b := range bots {
		if b.party != nil {
			t.Errorf("bot %q should be freed when no human remains in the crew", b.Name)
		}
	}
}
