package cowboy

import (
	"strings"
	"testing"
)

func TestAreaLootWares(t *testing.T) {
	// Surface mobs drop street gear.
	if got := areaLootWares("ic_1"); len(got) == 0 {
		t.Error("city area should have street loot wares")
	}
	// A zone-1 underground room maps to band-1 vendor stock.
	z1 := areaLootWares("z1_05")
	if len(z1) == 0 {
		t.Fatal("zone-1 area should have loot wares")
	}
	// A deep zone maps to a higher band with different (better) gear.
	z8 := areaLootWares("z8_13")
	names := func(ws []ware) string {
		var b strings.Builder
		for _, w := range ws {
			b.WriteString(w.name + " ")
		}
		return b.String()
	}
	if names(z1) == names(z8) {
		t.Error("deep-zone loot should differ from shallow-zone loot")
	}
}

func TestPartyLootSkipsBots(t *testing.T) {
	w := NewWorld(NewMemStore())
	w.EnableBots(1)
	bot := botsIn(w)[0]
	boss, _ := newTestPlayer(w, "Boss", "z1_05")
	w.invite(boss, bot.Name) // bot joins, warps to boss's room
	boss.Class = "hacker"
	bot.Class = "enforcer" // a different class — would add its own item if counted

	loot := map[string]int{}
	w.addPartyLoot(boss, loot)

	// Exactly one class item (the human's), never the bot's.
	total := 0
	for _, n := range loot {
		total += n
	}
	if total != 1 {
		t.Errorf("party loot should include only the human's class item (1), got %d items: %v", total, loot)
	}
}

func TestGsayAllBotsRespondInVoice(t *testing.T) {
	w := NewWorld(NewMemStore())
	w.SetRoll(func(n int) int { return 0 }) // deterministic line pick
	w.EnableBots(5)
	bots := botsIn(w)
	boss, drain := newTestPlayer(w, "Boss", "ic_1")
	for _, b := range bots { // recruit all five into the crew
		w.invite(boss, b.Name)
	}
	drain() // clear join noise

	w.groupChat(boss, "form up on me")
	out := drain()

	// EVERY crewed bot in the room answers (not just two).
	for _, b := range bots {
		if !strings.Contains(out, b.Name+": ") {
			t.Errorf("crewed bot %q should answer GSAY, missing from: %q", b.Name, out)
		}
	}
	// And the reply is in that bot's class voice (e.g. a netrunner's first line).
	for _, b := range bots {
		if pool := botVoiceReplies[strings.ToLower(b.Class)]; len(pool) > 0 {
			if !strings.Contains(out, pool[0]) {
				t.Errorf("bot %q (%s) should answer in class voice %q", b.Name, b.Class, pool[0])
			}
		}
	}
}
