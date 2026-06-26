package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestSell checks selling an item at a vendor pays buyback scrip and removes it (#40).
func TestSell(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out)

	p.RoomID = "market" // Night Market vendor
	p.Eddies = 0
	p.Inv["stimpak"] = 2 // catalog price 20 → 50% = 10 each

	w.Command(p, "sell stimpak 2")
	if p.Inv["stimpak"] != 0 {
		t.Errorf("sold stimpaks should leave the pack: have %d", p.Inv["stimpak"])
	}
	if p.Eddies != 20 { // 2 × (20×50%)
		t.Errorf("sell should pay 50%% buyback (expected 20): have %d", p.Eddies)
	}

	// Selling away from a vendor is refused.
	p.RoomID = "back_alley"
	p.Inv["stimpak"] = 1
	buf.Reset()
	w.Command(p, "sell stimpak")
	if !strings.Contains(buf.String(), "no vendor") {
		t.Errorf("sell off a vendor should be refused; got:\n%s", buf.String())
	}
}
