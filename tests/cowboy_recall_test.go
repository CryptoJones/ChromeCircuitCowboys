package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestRecallHomeCompletes: HOME jacks a timed recall from anywhere that, left
// uninterrupted, phases the runner back to their Re-Clone Bay (issue #4).
func TestRecallHomeCompletes(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, _ := sink()
	p := w.Connect("Case", out)
	p.RoomID = "the_sprawl" // a quiet, mob-free street room

	w.Command(p, "home")
	if p.RoomID != "the_sprawl" {
		t.Fatalf("recall should not be instant; moved to %s", p.RoomID)
	}
	// Hold still through the cast — it lands at the Re-Clone Bay.
	for i := 0; i < 6 && p.RoomID != "capsule"; i++ {
		w.Tick()
	}
	if p.RoomID != "capsule" {
		t.Fatalf("recall should phase the runner home, at %s", p.RoomID)
	}
}

// TestRecallBreaksOnHit: a mob attack during the cast shatters the recall — you
// stay put and have to deal with the fight (issue #4).
func TestRecallBreaksOnHit(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(alwaysHit)
	out, buf := sink()
	p := w.Connect("Case", out)
	p.MaxHP, p.HP = 1000, 1000
	p.RoomID = "back_alley" // an aggressive ganger lives here

	w.Command(p, "home") // homing is engine-internal; we assert behaviorally
	w.Tick()             // the ganger aggros and hits -> recall shatters
	if p.RoomID != "back_alley" {
		t.Fatalf("an interrupted recall must NOT teleport; at %s", p.RoomID)
	}
	if !strings.Contains(buf.String(), "recall shatters") {
		t.Errorf("expected the recall to be interrupted by the attack; got:\n%s", lastLines(buf.String()))
	}
	// And it does not belatedly fire on later ticks.
	for i := 0; i < 6; i++ {
		w.Tick()
	}
	if p.RoomID == "capsule" {
		t.Fatal("a broken recall must not complete later")
	}
}
