package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestTalkGivesLore checks TALK delivers level backstory, spoken by the local
// fixer where one is hiring (#12).
func TestTalkGivesLore(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(func(n int) int { return 0 }) // first lore line, deterministically
	out, buf := sink()
	p := w.Connect("Case", out)

	// z1_01 hosts Marcus the Fixer — TALK should name him and give z1 lore.
	p.RoomID = "z1_01"
	buf.Reset()
	w.Command(p, "talk")
	got := buf.String()
	if !strings.Contains(got, "Marcus") {
		t.Errorf("TALK at z1_01 should be spoken by Marcus the Fixer; got:\n%s", got)
	}
	if !strings.Contains(got, "Neon Wasteland") {
		t.Errorf("TALK should deliver z1 backstory; got:\n%s", got)
	}

	// On the surface (no authored zone) there's nothing to say.
	p.RoomID = "neon_alley"
	buf.Reset()
	w.Command(p, "talk")
	if !strings.Contains(buf.String(), "nothing to say") {
		t.Errorf("TALK on the surface should have no lore; got:\n%s", buf.String())
	}
}
