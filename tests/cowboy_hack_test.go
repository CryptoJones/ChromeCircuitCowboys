package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestHackMinigame checks HACK at a terminal runs a crackable high/low game that
// pays out on success (#33).
func TestHackMinigame(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	// roll(100) returns 99 → secret 100; guesses below say "higher".
	w.SetRoll(func(n int) int { return n - 1 })
	out, buf := sink()
	p := w.Connect("Case", out)

	// Off a terminal, HACK is refused.
	p.RoomID = "back_alley" // not a terminal
	buf.Reset()
	w.Command(p, "hack")
	if !strings.Contains(buf.String(), "need a data terminal") {
		t.Errorf("HACK off a terminal should be refused; got:\n%s", buf.String())
	}

	// At a terminal (Night Market), start a run and crack it.
	p.RoomID = "market"
	p.Eddies = 0
	w.Command(p, "hack")
	buf.Reset()
	w.Command(p, "100") // secret is 100 (roll returns max)
	if !strings.Contains(buf.String(), "CRACKED") {
		t.Errorf("guessing the secret should crack the system; got:\n%s", buf.String())
	}
	if p.Eddies <= 0 {
		t.Errorf("cracking should pay scrip; have %d", p.Eddies)
	}
}
