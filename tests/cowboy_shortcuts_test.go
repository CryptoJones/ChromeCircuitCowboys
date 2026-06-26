package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestAttackLootShortcuts checks A = attack and LO = loot, with L still = look (#14).
func TestAttackLootShortcuts(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(alwaysHit)
	out, buf := sink()
	p := w.Connect("Case", out)

	// A = attack: engaging the back-alley ganger.
	p.RoomID = "back_alley"
	p.HP, p.MaxHP = 400, 400
	buf.Reset()
	w.Command(p, "a ganger")
	if !strings.Contains(buf.String(), "ganger") {
		t.Errorf("A should attack; got:\n%s", buf.String())
	}
	for i := 0; i < 30; i++ {
		w.Tick()
	}

	// LO = loot the remains.
	buf.Reset()
	w.Command(p, "lo")
	if strings.Contains(buf.String(), "Unknown command") {
		t.Errorf("LO should loot, not be unknown; got:\n%s", buf.String())
	}

	// L still = look (unchanged).
	buf.Reset()
	w.Command(p, "l")
	if strings.Contains(buf.String(), "Unknown command") {
		t.Errorf("L should still look; got:\n%s", buf.String())
	}
}
