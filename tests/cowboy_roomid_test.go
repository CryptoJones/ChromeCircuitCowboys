package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestRoomIDCommand checks the hidden roomid command prints the room's id (#28).
func TestRoomIDCommand(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out)

	p.RoomID = "z3_07"
	buf.Reset()
	w.Command(p, "roomid")
	if !strings.Contains(buf.String(), "z3_07") {
		t.Errorf("roomid should print the room id; got:\n%s", buf.String())
	}

	// It must stay OUT of the HELP listing.
	buf.Reset()
	w.Command(p, "help")
	if strings.Contains(buf.String(), "roomid") {
		t.Errorf("roomid is a hidden command and must not appear in HELP")
	}
}
