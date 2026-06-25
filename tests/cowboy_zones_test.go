package tests

import (
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestUndergroundDescentAndLootCache smoke-tests the authored L1-99 carve: you
// descend from the street into the Neon Wasteland, then break open a hidden
// ceiling loot-cache and recover its consumable + scrip.
func TestUndergroundDescentAndLootCache(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(func(n int) int { return 0 }) // deterministic to-hit

	p := w.Connect("descender", func(string) {})

	// Back Alley drops down into the first authored zone.
	p.RoomID = "back_alley"
	w.Command(p, "down")
	if p.RoomID != "z1_01" {
		t.Fatalf("descent failed: at %q, want z1_01", p.RoomID)
	}

	// A fresh clone starts with one stimpak; record scrip before looting.
	startStim := p.Inv["stimpak"]
	startScrip := p.Eddies

	// Step into the hidden ceiling cache off The Sodium Strip and crack it open.
	p.RoomID = "z1_02_cache"
	w.Command(p, "attack cache")
	for i := 0; i < 3 && p.Inv["stimpak"] == startStim; i++ {
		w.Tick()
		w.Command(p, "loot")
	}

	if p.Inv["stimpak"] <= startStim {
		t.Fatalf("loot cache yielded no consumable (stimpak %d -> %d)", startStim, p.Inv["stimpak"])
	}
	if p.Eddies <= startScrip {
		t.Fatalf("loot cache yielded no scrip (%d -> %d)", startScrip, p.Eddies)
	}
}
