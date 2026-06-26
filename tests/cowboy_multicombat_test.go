package tests

import (
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestMultiAttacker checks two runners can damage the same mob in one fight, and
// the mob's retaliation can land on either of them (#41).
func TestMultiAttacker(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(alwaysHit)
	o1, _ := sink()
	o2, _ := sink()
	a := w.Connect("Rook", o1)
	b := w.Connect("Jett", o2)
	a.RoomID, b.RoomID = "back_alley", "back_alley"
	a.HP, a.MaxHP = 500, 500
	b.HP, b.MaxHP = 500, 500
	a.Body, b.Body = 0, 0 // weak hits so the ganger survives a few ticks and fights back

	// Both pile onto the ganger.
	w.Command(a, "attack ganger")
	w.Command(b, "attack ganger")

	// Track the lowest HP each runner hits during the fight (passive regen
	// refills it afterward, which would otherwise mask the damage taken).
	minA, minB := a.HP, b.HP
	for i := 0; i < 30; i++ {
		w.Tick()
		if a.HP < minA {
			minA = a.HP
		}
		if b.HP < minB {
			minB = b.HP
		}
	}

	// While both were attacking the same mob, the mob retaliated on its
	// attackers — at least one of them took a hit (shared threat, #41).
	if minA >= 500 && minB >= 500 {
		t.Errorf("a group fight should land hits on the attackers: minA=%d minB=%d", minA, minB)
	}
}
