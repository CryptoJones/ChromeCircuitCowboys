package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestOpenCache checks the OPEN verb cracks open a supply cache (no ATTACK
// needed) and yields its loot (#16).
func TestOpenCache(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(func(n int) int { return 0 })
	out, buf := sink()
	p := w.Connect("Case", out)

	startStim := p.Inv["stimpak"]
	p.RoomID = "z1_02_cache"
	buf.Reset()
	w.Command(p, "open")
	for i := 0; i < 3 && p.Inv["stimpak"] == startStim; i++ {
		w.Tick()
		w.Command(p, "loot")
	}

	got := buf.String()
	if !strings.Contains(got, "cracks open") {
		t.Errorf("OPEN should crack the cache open; got:\n%s", got)
	}
	if p.Inv["stimpak"] <= startStim {
		t.Fatalf("OPEN+loot yielded no consumable (stimpak %d -> %d)", startStim, p.Inv["stimpak"])
	}

	// OPEN where there's no container should say so, not crack anything.
	p.RoomID = "z1_01"
	buf.Reset()
	w.Command(p, "open")
	if !strings.Contains(buf.String(), "nothing to open") {
		t.Errorf("OPEN with no cache should report nothing to open; got:\n%s", buf.String())
	}
}
