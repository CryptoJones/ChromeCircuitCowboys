package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestLookItem checks LOOK <item> examines an item (not the room) (#53).
func TestLookItem(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out)

	buf.Reset()
	w.Command(p, "look subdermal-plating")
	got := buf.String()
	if !strings.Contains(got, "subdermal-plating") {
		t.Errorf("LOOK <item> should examine the item; got:\n%s", got)
	}
	if !strings.Contains(got, "Body implant") {
		t.Errorf("examine should describe the effect; got:\n%s", got)
	}

	// Bare LOOK still describes the room.
	buf.Reset()
	w.Command(p, "look")
	if !strings.Contains(buf.String(), "Exits:") {
		t.Errorf("bare LOOK should still describe the room; got:\n%s", buf.String())
	}
}
