package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestCharacterPointsAndSpend checks that leveling banks character points, the
// sheet shows them, and SPEND raises a stat and decrements the pool (#11).
func TestCharacterPointsAndSpend(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(alwaysHit)
	out, buf := sink()
	p := w.Connect("Case", out)

	// Force a level-up to bank points.
	p.XP = 1_000_000
	p.RoomID = "back_alley"
	for i := 0; i < 4 && p.StatPoints == 0; i++ {
		w.Command(p, "attack ganger")
		w.Tick()
	}
	if p.StatPoints <= 0 {
		t.Fatalf("level-up should bank character points, got %d", p.StatPoints)
	}

	// The score sheet advertises them.
	buf.Reset()
	w.Command(p, "score")
	if !strings.Contains(buf.String(), "Character points") {
		t.Errorf("score should show available character points; got:\n%s", buf.String())
	}

	// SPEND raises the stat and decrements the pool.
	before := p.Reflexes
	pts := p.StatPoints
	w.Command(p, "spend reflexes 1")
	if p.Reflexes != before+1 {
		t.Errorf("spend reflexes should +1: %d -> %d", before, p.Reflexes)
	}
	if p.StatPoints != pts-1 {
		t.Errorf("spend should decrement points: %d -> %d", pts, p.StatPoints)
	}

	// Spending Body lifts MaxHP.
	hpBefore := p.MaxHP
	if p.StatPoints > 0 {
		w.Command(p, "spend body 1")
		if p.MaxHP <= hpBefore {
			t.Errorf("spending Body should raise MaxHP: %d -> %d", hpBefore, p.MaxHP)
		}
	}
}
