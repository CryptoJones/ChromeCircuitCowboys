package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestDroneLeavesWreckNotBody checks that a machine foe (a drone) destroyed and
// looted reads as wreckage, never a "body"/"corpse" (#8).
func TestDroneLeavesWreckNotBody(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(alwaysHit)
	out, buf := sink()
	p := w.Connect("Case", out)

	// Stand with the Kurokawa security drone on the Perimeter Fenceline (z1_13).
	p.RoomID = "z1_13"
	p.HP, p.MaxHP = 400, 400 // survive the fight; we only care about the wording
	buf.Reset()
	w.Command(p, "attack drone")
	for i := 0; i < 40; i++ {
		w.Tick()
	}
	w.Command(p, "look")
	w.Command(p, "loot")

	got := buf.String()
	if !strings.Contains(got, "wreck") {
		t.Errorf("a destroyed drone should leave wreckage; got:\n%s", got)
	}
	if strings.Contains(got, "body") || strings.Contains(got, "corpse") {
		t.Errorf("a bodiless drone must not be called a body/corpse; got:\n%s", got)
	}
}
