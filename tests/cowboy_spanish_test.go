package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestSpanishNPC checks a Spanish-speaking local greets you in the rings (#21).
func TestSpanishNPC(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(func(n int) int { return 0 })
	out, buf := sink()
	p := w.Connect("Case", out)

	p.RoomID = "ic_1" // Inner Circuit :: Neon Gate — Rosa the flower-cart vendor
	buf.Reset()
	w.Command(p, "talk")
	got := buf.String()
	if !strings.Contains(got, "Rosa") {
		t.Errorf("TALK at ic_1 should reach Rosa; got:\n%s", got)
	}
	if !strings.Contains(got, "Ciudad de la Noche") {
		t.Errorf("Rosa should greet in Spanish; got:\n%s", got)
	}
}
