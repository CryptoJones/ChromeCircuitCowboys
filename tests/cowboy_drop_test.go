package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestDropGet checks dropping items to the floor and another player picking them
// up, plus DROP ALL (#51).
func TestDropGet(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	o1, _ := sink()
	o2, b2 := sink()
	a := w.Connect("Rook", o1)
	b := w.Connect("Jett", o2)
	a.RoomID, b.RoomID = "neon_alley", "neon_alley"
	b.Inv["stimpak"] = 0 // fresh clones start with one; clear for a clean count

	a.Inv["stimpak"] = 3
	w.Command(a, "drop stimpak 2")
	if a.Inv["stimpak"] != 1 {
		t.Fatalf("drop should leave 1 stimpak: have %d", a.Inv["stimpak"])
	}

	// LOOK (by Jett) shows the floor pile.
	b2.Reset()
	w.Command(b, "look")
	if !strings.Contains(b2.String(), "On the floor") {
		t.Errorf("look should show the floor pile; got:\n%s", b2.String())
	}

	// The other player picks it up.
	w.Command(b, "get stimpak")
	if b.Inv["stimpak"] != 2 {
		t.Errorf("Jett should pick up the 2 dropped stimpaks: have %d", b.Inv["stimpak"])
	}

	// DROP ALL dumps the whole pack.
	a.Inv["ram-chip"] = 1
	w.Command(a, "drop all")
	if len(a.Inv) != 0 {
		t.Errorf("drop all should empty the pack: %v", a.Inv)
	}
	w.Command(a, "get all")
	if a.Inv["stimpak"] != 1 || a.Inv["ram-chip"] != 1 {
		t.Errorf("get all should recover the pile: %v", a.Inv)
	}
}
