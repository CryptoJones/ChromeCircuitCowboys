package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestBoothOnboarding checks TALK in the Re-Clone Bay gives the new-player
// command primer (#20).
func TestBoothOnboarding(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out) // a fresh clone starts in the Re-Clone Bay

	buf.Reset()
	w.Command(p, "talk")
	got := buf.String()
	if !strings.Contains(got, "Doc Splice") {
		t.Errorf("TALK in the booth should be the onboarding tech; got:\n%s", got)
	}
	for _, kw := range []string{"MAP", "QUESTS", "ATTACK"} {
		if !strings.Contains(got, kw) {
			t.Errorf("onboarding should mention %s; got:\n%s", kw, got)
		}
	}
}
