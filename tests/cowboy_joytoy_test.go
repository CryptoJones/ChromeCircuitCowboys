package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestJoytoyRestores checks PAYing a joytoy in the red-light strip costs scrip
// and fully restores you (fade-to-black, #27).
func TestJoytoyRestores(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out)

	p.RoomID = "ic_5" // the Rolling Rose
	p.Eddies = 200
	p.HP = 1
	buf.Reset()
	w.Command(p, "pay")
	got := buf.String()
	if p.HP != p.MaxHP {
		t.Errorf("PAY should fully restore HP: %d/%d", p.HP, p.MaxHP)
	}
	if p.Eddies != 125 {
		t.Errorf("PAY should cost 75 scrip: have %d", p.Eddies)
	}
	if !strings.Contains(got, "unwound") {
		t.Errorf("PAY should be a fade-to-black restore; got:\n%s", got)
	}

	// Not enough scrip → refused.
	p.Eddies = 10
	p.HP = 1
	buf.Reset()
	w.Command(p, "pay")
	if p.HP != 1 {
		t.Errorf("PAY with too little scrip should not restore")
	}

	// Nowhere to pay outside the strip.
	p.RoomID = "ic_1"
	p.Eddies = 200
	buf.Reset()
	w.Command(p, "pay")
	if !strings.Contains(buf.String(), "no one here") {
		t.Errorf("PAY off the strip should report no one; got:\n%s", buf.String())
	}
}
