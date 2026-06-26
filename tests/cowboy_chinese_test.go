package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestChineseNPC checks a Chinese-speaking local trash-talks you, with a gloss (#22).
func TestChineseNPC(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(func(n int) int { return 0 })
	out, buf := sink()
	p := w.Connect("Case", out)

	p.RoomID = "sb_3" // Sprawlbelt :: Tagged Underpass — the sneering tagger
	buf.Reset()
	w.Command(p, "talk")
	got := buf.String()
	if !strings.Contains(got, "克隆") { // contains Chinese (克隆 = clone)
		t.Errorf("the tagger should trash-talk in Chinese; got:\n%s", got)
	}
	if !strings.Contains(got, "vat") { // dim English gloss present
		t.Errorf("expected an English gloss; got:\n%s", got)
	}
}
