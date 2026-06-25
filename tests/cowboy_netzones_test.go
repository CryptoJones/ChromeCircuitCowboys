package tests

import (
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestNetAscentDiveAndDataCache smoke-tests the authored Net carve: you jack in
// at the first node's shell, dive DOWN into the breach layer (a netrun BREACH
// that spends RAM), travel the lateral MID thoroughfare, then dive to a Core
// data-cache and salvage its RAM + scrip.
func TestNetAscentDiveAndDataCache(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(alwaysHit) // deterministic hits against ICE

	p := w.Connect("netrunner", func(string) {})
	p.MaxHP, p.HP, p.RAM = 2000, 2000, 40

	// Jack in: Data Port -> up -> the first Net shell.
	p.RoomID = "data_port"
	w.Command(p, "up")
	if p.RoomID != "nz1_1_top" {
		t.Fatalf("jack-in failed: at %q, want nz1_1_top", p.RoomID)
	}

	// Dive into the breach layer; attacking here is a netrun breach (spends RAM).
	w.Command(p, "down")
	ramBefore := p.RAM
	w.Command(p, "attack")
	w.Tick()
	if p.RAM >= ramBefore {
		t.Fatalf("a breach in the Net should spend RAM (%d -> %d)", ramBefore, p.RAM)
	}

	// Travel to the second node (lateral, at the MID layer) and dive to its
	// data-cache core, then crack it for loot.
	p.RoomID = "nz1_2_bot"
	startRAMitems := p.Inv["ram-chip"]
	startScrip := p.Eddies
	w.Command(p, "attack cache")
	for i := 0; i < 3 && p.Inv["ram-chip"] == startRAMitems; i++ {
		w.Tick()
		w.Command(p, "loot")
	}
	if p.Inv["ram-chip"] <= startRAMitems {
		t.Fatalf("data-cache yielded no RAM (ram-chip %d -> %d)", startRAMitems, p.Inv["ram-chip"])
	}
	if p.Eddies <= startScrip {
		t.Fatalf("data-cache yielded no scrip (%d -> %d)", startScrip, p.Eddies)
	}
}
