package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestQuickUseBySlot checks a bare number uses that inventory slot, and the
// inventory is numbered (#13).
func TestQuickUseBySlot(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out)

	// Numbered inventory listing.
	buf.Reset()
	w.Command(p, "inventory")
	if !strings.Contains(buf.String(), "1)") {
		t.Errorf("inventory should be numbered; got:\n%s", buf.String())
	}

	// A fresh runner starts with a stimpak (slot 1, alphabetically). Take damage
	// then quick-use slot 1 to heal.
	stims := p.Inv["stimpak"]
	if stims < 1 {
		t.Fatalf("expected a starting stimpak, inv=%v", p.Inv)
	}
	p.HP = 1
	w.Command(p, "1") // quick-use slot 1
	if p.Inv["stimpak"] != stims-1 {
		t.Errorf("quick-use should consume the stimpak: %d -> %d", stims, p.Inv["stimpak"])
	}
	if p.HP <= 1 {
		t.Errorf("quick-using a stimpak should heal: HP=%d", p.HP)
	}

	// An out-of-range slot is handled gracefully.
	buf.Reset()
	w.Command(p, "9")
	if !strings.Contains(buf.String(), "No item #9") {
		t.Errorf("out-of-range quick-use should report no such slot; got:\n%s", buf.String())
	}
}
