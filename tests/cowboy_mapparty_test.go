package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestMapToParty checks the MAP points toward a separated crewmate (#52).
func TestMapToParty(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	o1, b1 := sink()
	o2, _ := sink()
	leader := w.Connect("Rook", o1)
	member := w.Connect("Jett", o2)

	leader.RoomID, member.RoomID = "neon_alley", "neon_alley"
	w.Command(leader, "invite Jett")
	w.Command(member, "accept")

	// Separate them: member walks east to the_sprawl, leader stays.
	member.RoomID = "the_sprawl"

	b1.Reset()
	w.Command(leader, "map")
	got := b1.String()
	if !strings.Contains(got, "TO YOUR CREW") {
		t.Errorf("map should point to the separated crewmate; got:\n%s", got)
	}
	if !strings.Contains(got, "Jett") {
		t.Errorf("map should name the crewmate; got:\n%s", got)
	}
}
