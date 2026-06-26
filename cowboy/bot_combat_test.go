package cowboy

import "testing"

// liveMobAnywhere grabs any spawned live mob and relocates it to room, for tests.
func mobInRoom(w *World, room string) *Mob {
	for _, m := range w.mobs {
		if !m.dead {
			m.RoomID = room
			return m
		}
	}
	return nil
}

func TestBotAssistsAndDisengages(t *testing.T) {
	w, boss, bot := crewBot(t, "ic_1")
	mob := mobInRoom(w, boss.RoomID)
	if mob == nil {
		t.Fatal("no live mob to test with")
	}

	// Human crewmate is fighting the mob → bot joins it.
	boss.fighting = mob
	w.botAssist()
	if bot.fighting != mob {
		t.Errorf("bot should assist the crewmate's fight; bot.fighting=%v", bot.fighting)
	}

	// Crewmate stops fighting → bot disengages.
	boss.fighting = nil
	w.botAssist()
	if bot.fighting != nil {
		t.Errorf("bot should disengage when no crewmate is fighting; bot.fighting=%v", bot.fighting)
	}

	// Bot never initiates on a passive mob (no crewmate fighting).
	w.botAssist()
	if bot.fighting != nil {
		t.Error("bot must not initiate combat on its own")
	}
}

func TestKillCreditGoesToHuman(t *testing.T) {
	w, boss, bot := crewBot(t, "ic_1")
	mob := mobInRoom(w, boss.RoomID)
	boss.fighting = mob // boss is a swinging attacker of the mob

	if got := w.killCredit(bot, mob); got != boss {
		t.Errorf("a bot's kill should credit the human crewmate, got %v", got)
	}
	if got := w.killCredit(boss, mob); got != boss {
		t.Errorf("a human's kill credits themselves, got %v", got)
	}
}

func TestBotNeverFlatlines(t *testing.T) {
	w, boss, bot := crewBot(t, "ic_1")
	w.SetRoll(func(n int) int { return n - 1 }) // bias rolls high so the mob lands its hit
	mob := mobInRoom(w, boss.RoomID)
	if mob.tmpl.Damage < 1 {
		mob.tmpl = &MobTemplate{Name: "Test Brute", Damage: 999, AC: 0, HP: 999}
		mob.HP = 999
	}
	boss.fighting = mob
	bot.fighting = mob // bot is an attacker, so the mob can hit it
	bot.HP = 1

	w.resolveCombat()

	if bot.HP < 1 {
		t.Errorf("crewed bot should never drop below 1 HP, got %d", bot.HP)
	}
	if bot.RoomID == startRoom {
		t.Error("crewed bot should not be re-cloned to the start room")
	}
	if mob.target == bot {
		t.Error("a mob should never persistently target a bot")
	}
}

func TestCrewedBotCatchesUpToLeader(t *testing.T) {
	w, boss, bot := crewBot(t, "ic_1")
	// Separate the bot from the leader, idle.
	bot.RoomID = "sb_3"
	bot.fighting = nil
	w.tickBots()
	if bot.RoomID != boss.RoomID {
		t.Errorf("idle crewed bot should regroup on the leader (%s), is in %s", boss.RoomID, bot.RoomID)
	}
}
