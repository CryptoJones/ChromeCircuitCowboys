package cowboy

import (
	"strings"
	"testing"
)

func TestInviteAllRecruitsRoomRunners(t *testing.T) {
	w := NewWorld(NewMemStore())
	w.EnableBots(3)
	bots := botsIn(w)
	boss, _ := newTestPlayer(w, "Boss", "ic_1")
	bots[0].RoomID = "ic_1"
	bots[1].RoomID = "ic_1"
	bots[2].RoomID = "sb_3" // not in the room — must be skipped

	w.inviteAll(boss)

	if boss.party == nil {
		t.Fatal("inviteAll should form a crew")
	}
	if bots[0].party != boss.party || bots[1].party != boss.party {
		t.Error("free bots in the room should fall in")
	}
	if bots[2].party != nil {
		t.Error("a bot in another room must not be recruited")
	}

	// A bot already in someone else's crew can't be poached by inviteAll.
	other := &Party{Leader: bots[2]}
	other.add(bots[2])
	bots[2].RoomID = "ic_1"
	boss2, _ := newTestPlayer(w, "Boss2", "ic_1")
	w.inviteAll(boss2)
	if bots[2].party == boss2.party {
		t.Error("a crewed bot must not be poached")
	}
}

func TestRallyAndPartyCombatLog(t *testing.T) {
	w := NewWorld(NewMemStore())
	w.EnableBots(1)
	bot := botsIn(w)[0]
	boss, drain := newTestPlayer(w, "Boss", "ic_1")
	w.invite(boss, bot.Name) // bot auto-joins and warps to ic_1
	drain()                  // clear the join chatter

	mob := mobInRoom(w, "ic_1")
	w.rallyCrewBots(boss, mob)
	if bot.fighting != mob {
		t.Errorf("rally should immediately engage the crewed bot, got %v", bot.fighting)
	}

	// The leader sees a crewmate's action relayed into their log.
	w.partyCombatLog(bot, "RIKO-HIT\r\n")
	if out := drain(); !strings.Contains(out, "RIKO-HIT") {
		t.Error("leader should see a crewmate's combat action")
	}
}

func TestInstallByNumberAndQuantity(t *testing.T) {
	w := NewWorld(NewMemStore())
	boss, _ := newTestPlayer(w, "Boss", "market")  // Night Market is a medic room
	boss.Inv = map[string]int{"reflex-booster": 8} // sole item → INVENTORY #1; +2 Reflexes each
	before := boss.Reflexes

	w.install(boss, "1 8") // install 8 of item #1

	if got := boss.Reflexes - before; got != 16 {
		t.Errorf("INSTALL 1 8 should add 8×2 = 16 Reflexes, added %d", got)
	}
	if boss.Inv["reflex-booster"] != 0 {
		t.Errorf("all 8 should be consumed, %d left", boss.Inv["reflex-booster"])
	}

	// Quantity caps at what you carry.
	boss.Inv = map[string]int{"reflex-booster": 2}
	before = boss.Reflexes
	w.install(boss, "1 100")
	if got := boss.Reflexes - before; got != 4 {
		t.Errorf("qty should cap at 2 carried (=+4 Reflexes), added %d", got)
	}
}
