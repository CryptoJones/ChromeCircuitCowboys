package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestMapCommand checks the MAP command labels the current area and points the
// player onward to the next harder zone (and back out from a deeper one).
func TestMapCommand(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Mapper", out)

	// The first Undercity zone: the map names the band and shows the way deeper.
	p.RoomID = "z1_01"
	buf.Reset()
	w.Command(p, "map")
	got := buf.String()
	if !strings.Contains(got, "THE UNDERCITY L1-10") {
		t.Errorf("map should label the Undercity L1-10 band; got:\n%s", got)
	}
	if !strings.Contains(got, "PROCEED to the next harder area") {
		t.Errorf("map should point the way deeper; got:\n%s", got)
	}

	// A deeper zone: the map offers a way back out toward easier ground.
	p.RoomID = "z2_01"
	buf.Reset()
	w.Command(p, "map")
	got = buf.String()
	if !strings.Contains(got, "THE UNDERCITY L11-20") {
		t.Errorf("map should label the L11-20 band; got:\n%s", got)
	}
	if !strings.Contains(got, "WAY OUT") {
		t.Errorf("a deeper zone should show a way out; got:\n%s", got)
	}
}
