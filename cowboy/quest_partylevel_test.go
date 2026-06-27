package cowboy

import "testing"

func TestPartyWaivesQuestLevel(t *testing.T) {
	q, ok := questByID("ug1_heist") // Giver z1_11, MinLevel 7
	if !ok {
		t.Skip("quest fixture missing")
	}

	// Solo + under-level: ACCEPT is refused, quest not taken.
	w := NewWorld(NewMemStore())
	solo, _ := newTestPlayer(w, "Solo", q.Giver)
	solo.Level = 1
	w.accept(solo, "all")
	if _, took := solo.Quests[q.ID]; took {
		t.Error("under-level solo runner should NOT be able to accept the bounty")
	}

	// In a crew + under-level: the level requirement is waived, quest is taken.
	w2 := NewWorld(NewMemStore())
	w2.EnableBots(1)
	bot := botsIn(w2)[0]
	boss, _ := newTestPlayer(w2, "Boss", q.Giver)
	boss.Level = 1
	w2.invite(boss, bot.Name) // forms a crew (boss + bot)
	if boss.party == nil {
		t.Fatal("crew should have formed")
	}
	w2.accept(boss, "all")
	if _, took := boss.Quests[q.ID]; !took {
		t.Error("a crewed under-level runner should be able to accept the bounty")
	}
}
