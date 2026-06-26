package tests

import (
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestPartyLoot checks a party kill drops something for each class present (#50).
func TestPartyLoot(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(alwaysHit)
	o1, _ := sink()
	o2, _ := sink()
	leader := w.Connect("Rook", o1)
	member := w.Connect("Jett", o2)
	leader.Class, member.Class = "enforcer", "hacker"

	// Form a crew, both in the back alley with the ganger.
	leader.RoomID, member.RoomID = "neon_alley", "neon_alley"
	w.Command(leader, "invite Jett")
	w.Command(member, "accept")
	leader.RoomID, member.RoomID = "back_alley", "back_alley"
	leader.HP, leader.MaxHP = 400, 400

	// Leader kills the ganger; loot it and check both a heal (enforcer) and RAM
	// (hacker) ended up on the corpse.
	startStim := leader.Inv["stimpak"]
	startRam := leader.Inv["ram-chip"]
	w.Command(leader, "attack ganger")
	for i := 0; i < 40; i++ {
		w.Tick()
	}
	w.Command(leader, "loot")

	if leader.Inv["stimpak"] <= startStim {
		t.Errorf("party loot should include a heal for the Enforcer; stimpak %d -> %d", startStim, leader.Inv["stimpak"])
	}
	if leader.Inv["ram-chip"] <= startRam {
		t.Errorf("party loot should include RAM for the Hacker; ram-chip %d -> %d", startRam, leader.Inv["ram-chip"])
	}
}
