package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestDixieEasterEgg checks the laughing ROM construct in the rings calls you
// "Boy" (the oblique homage, #23).
func TestDixieEasterEgg(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(func(n int) int { return 0 })
	out, buf := sink()
	p := w.Connect("Case", out)

	p.RoomID = "sb_7"
	buf.Reset()
	w.Command(p, "talk")
	got := buf.String()
	if !strings.Contains(got, "ROM construct") {
		t.Errorf("sb_7 should host the ROM construct; got:\n%s", got)
	}
	if !strings.Contains(got, "Boy") {
		t.Errorf("the construct should call you Boy; got:\n%s", got)
	}
}
