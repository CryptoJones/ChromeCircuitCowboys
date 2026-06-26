package tests

import (
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestPartyFollow checks crew members follow the leader when they move (#42).
func TestPartyFollow(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	o1, _ := sink()
	o2, _ := sink()
	leader := w.Connect("Rook", o1)
	member := w.Connect("Jett", o2)

	leader.RoomID, member.RoomID = "neon_alley", "neon_alley"
	w.Command(leader, "invite Jett")
	w.Command(member, "accept")

	// Leader moves east; the member should be pulled along.
	w.Command(leader, "east")
	if leader.RoomID != "the_sprawl" {
		t.Fatalf("leader move failed: at %s", leader.RoomID)
	}
	if member.RoomID != "the_sprawl" {
		t.Errorf("crew member should follow the leader: at %s, want the_sprawl", member.RoomID)
	}
}
