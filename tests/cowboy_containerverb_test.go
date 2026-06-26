package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestContainerVerb checks attacking/opening a cache uses a container verb, not
// "lunge at" (#46).
func TestContainerVerb(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(func(n int) int { return 0 }) // first container verb deterministically
	out, buf := sink()
	p := w.Connect("Case", out)

	p.RoomID = "z1_02_cache"
	buf.Reset()
	w.Command(p, "open")
	got := buf.String()
	if strings.Contains(got, "lunge at") {
		t.Errorf("a cache must not be 'lunged at'; got:\n%s", got)
	}
	if !strings.Contains(got, "pry at") { // roll=0 → first verb "You pry at "
		t.Errorf("opening a cache should use a container verb; got:\n%s", got)
	}
}
