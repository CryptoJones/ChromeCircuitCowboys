package tests

import (
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestScatterTerminals checks ~1/4 of surface ring rooms become terminals, and
// the seeded placement is identical across builds (#34).
func TestScatterTerminals(t *testing.T) {
	ring := []string{"ic_1", "ic_2", "ic_3", "ic_4", "ic_5", "ic_6",
		"sb_1", "sb_2", "sb_3", "sb_4", "sb_5", "sb_6", "sb_7", "sb_8", "sb_9", "sb_10"}

	count := func(w *cowboy.World) (n int, set string) {
		for _, id := range ring {
			if w.RoomIsTerminal(id) {
				n++
				set += id + ","
			}
		}
		return
	}

	w1 := cowboy.NewWorld(cowboy.NewMemStore())
	n1, s1 := count(w1)
	w2 := cowboy.NewWorld(cowboy.NewMemStore())
	n2, s2 := count(w2)

	if n1 == 0 {
		t.Fatal("expected some ring rooms to be scattered terminals, got none")
	}
	if s1 != s2 {
		t.Errorf("terminal placement must be deterministic across builds:\n %s\n vs %s", s1, s2)
	}
	_ = n2
	// Roughly a quarter (loose bound — some ring rooms are already vendor/medic).
	if n1 > len(ring)/2 {
		t.Errorf("too many ring terminals (%d of %d) — should be ~1/4", n1, len(ring))
	}
}
