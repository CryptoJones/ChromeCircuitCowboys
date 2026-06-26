package tests

import (
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestAcceptAll checks "accept all" takes every eligible bounty on the board
// and skips the ones gated above the player's level (#17).
func TestAcceptAll(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, _ := sink()
	p := w.Connect("Case", out) // level 1

	w.Command(p, "out")   // -> neon_alley
	w.Command(p, "south") // chrome_bar (street broker / fixer board)
	w.Command(p, "accept all")

	if len(p.Quests) < 1 {
		t.Fatalf("accept all took no bounties: %v", p.Quests)
	}
	// "Ghost in the Machine" is MinLevel 5 — a level-1 runner must not get it.
	if _, got := p.Quests["ghost_machine"]; got {
		t.Error("accept all should skip the level-5 bounty for a level-1 runner")
	}
	// A level-1 eligible bounty should be on.
	if _, got := p.Quests["clear_alley"]; !got {
		t.Error("accept all should have taken the level-1 alley bounty")
	}
}

// TestAcceptMultiple checks "accept 1 2" takes exactly those two.
func TestAcceptMultiple(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, _ := sink()
	p := w.Connect("Case", out)
	p.Level = 5 // eligible for the first two board bounties

	w.Command(p, "out")
	w.Command(p, "south")
	w.Command(p, "accept 1 2")
	if len(p.Quests) != 2 {
		t.Fatalf("accept 1 2 should take two bounties, got %d: %v", len(p.Quests), p.Quests)
	}
}
