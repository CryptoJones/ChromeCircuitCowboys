package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestQuestBoardStateTags checks the giver's board tags an accepted bounty as
// accepted (greyed) and a completed one as READY (#10).
func TestQuestBoardStateTags(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out)

	w.Command(p, "out")   // -> neon_alley
	w.Command(p, "south") // chrome_bar (fixer board)
	w.Command(p, "accept 1")

	buf.Reset()
	w.Command(p, "quests")
	if !strings.Contains(buf.String(), "[accepted") {
		t.Errorf("an in-progress bounty should be tagged accepted; got:\n%s", buf.String())
	}

	// Complete it and re-check: now READY.
	for id := range p.Quests {
		p.Quests[id] = 99 // force complete
	}
	buf.Reset()
	w.Command(p, "quests")
	if !strings.Contains(buf.String(), "READY — turn in") {
		t.Errorf("a completed bounty should be tagged READY on the board; got:\n%s", buf.String())
	}
}
