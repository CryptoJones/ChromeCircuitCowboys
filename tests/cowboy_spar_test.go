package tests

import (
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestSparKnockout checks that a PvP "kill" in the gym is a non-lethal knockout:
// the loser keeps their scrip, stays put, and wakes at full HP (#19).
func TestSparKnockout(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(alwaysHit)
	o1, _ := sink()
	o2, _ := sink()
	winner := w.Connect("Rook", o1)
	loser := w.Connect("Jett", o2)

	winner.RoomID = "sb_gym"
	loser.RoomID = "sb_gym"
	loser.HP = 1
	loser.Eddies = 500

	w.Command(winner, "attack Jett")
	for i := 0; i < 30 && loser.RoomID == "sb_gym" && loser.HP <= 1; i++ {
		w.Tick()
	}

	if loser.RoomID != "sb_gym" {
		t.Fatalf("a sparring loss must not re-sleeve: loser at %s", loser.RoomID)
	}
	if loser.HP != loser.MaxHP {
		t.Errorf("knocked-out sparrer should wake at full HP: %d/%d", loser.HP, loser.MaxHP)
	}
	if loser.Eddies != 500 {
		t.Errorf("sparring must not siphon scrip: have %d, want 500", loser.Eddies)
	}
}
