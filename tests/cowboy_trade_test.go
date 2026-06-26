package tests

import (
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestTradeSwap checks a two-sided, confirm-locked item+scrip swap (#25).
func TestTradeSwap(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	o1, _ := sink()
	o2, _ := sink()
	a := w.Connect("Rook", o1)
	b := w.Connect("Jett", o2)

	a.RoomID, b.RoomID = "neon_alley", "neon_alley"
	a.Inv["stimpak"] = 2
	a.Eddies = 100
	b.Inv["stimpak"] = 0 // fresh clones start with one; clear for a clean count
	b.Inv["ram-chip"] = 1
	b.Eddies = 100

	w.Command(a, "trade Jett")
	w.Command(a, "offer stimpak 2")
	w.Command(a, "offer scrip 50")
	w.Command(b, "offer ram-chip 1")

	// One-sided confirm does nothing yet.
	w.Command(a, "confirm")
	if a.Inv["stimpak"] != 2 {
		t.Fatalf("trade must not execute on one confirm")
	}

	// Both confirm → atomic swap.
	w.Command(b, "confirm")
	if a.Inv["stimpak"] != 0 || b.Inv["stimpak"] != 2 {
		t.Errorf("stimpaks should move to Jett: a=%d b=%d", a.Inv["stimpak"], b.Inv["stimpak"])
	}
	if b.Inv["ram-chip"] != 0 || a.Inv["ram-chip"] != 1 {
		t.Errorf("ram-chip should move to Rook: a=%d b=%d", a.Inv["ram-chip"], b.Inv["ram-chip"])
	}
	if a.Eddies != 50 || b.Eddies != 150 {
		t.Errorf("50 scrip should move Rook->Jett: a=%d b=%d", a.Eddies, b.Eddies)
	}
}

// TestTradeOfferResetsConfirm checks changing an offer clears confirmations (#25).
func TestTradeOfferResetsConfirm(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	o1, _ := sink()
	o2, _ := sink()
	a := w.Connect("Rook", o1)
	b := w.Connect("Jett", o2)
	a.RoomID, b.RoomID = "neon_alley", "neon_alley"
	a.Inv["stimpak"] = 3
	b.Inv["ram-chip"] = 2

	w.Command(a, "trade Jett")
	w.Command(a, "offer stimpak 1")
	w.Command(b, "offer ram-chip 1")
	w.Command(a, "confirm")
	w.Command(b, "confirm") // would execute...
	// ...it did execute (both confirmed). Re-open a fresh trade and verify a
	// mid-stream offer change blocks a stale confirm.
	w.Command(a, "trade Jett")
	w.Command(a, "offer stimpak 1")
	w.Command(a, "confirm")
	w.Command(a, "offer stimpak 2") // change after confirming -> resets both
	w.Command(b, "confirm")         // only b confirmed now; should NOT execute
	if a.Inv["stimpak"] != 2 {      // 3 - 1 (first trade) = 2; second trade not executed
		t.Errorf("a stale confirm after an offer change must not execute: stimpak=%d", a.Inv["stimpak"])
	}
}
