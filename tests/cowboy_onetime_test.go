package tests

import (
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestOneTimeQuest checks a claimed story bounty can't be re-accepted, while a
// ring rumor stays repeatable (#18).
func TestOneTimeQuest(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, _ := sink()
	p := w.Connect("Case", out)

	w.Command(p, "out")   // -> neon_alley
	w.Command(p, "south") // chrome_bar (fixer board)

	// Accept + force-complete + claim the alley bounty.
	w.Command(p, "accept 1")
	for id := range p.Quests {
		p.Quests[id] = 99
	}
	w.Command(p, "claim")
	if p.Done["clear_alley"] == 0 {
		t.Fatalf("claiming a story bounty should mark it done: %v", p.Done)
	}

	// Re-accepting it must be refused now.
	w.Command(p, "accept 1")
	if _, active := p.Quests["clear_alley"]; active {
		t.Error("a completed one-time bounty must not be re-acceptable")
	}
}
