package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestPointsAlwaysShown checks the sheet shows 0 when there are no character
// points, and a bold spend note when there are some (#31).
func TestPointsAlwaysShown(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out)

	// No points → "Character points: 0", no spend note.
	p.StatPoints = 0
	buf.Reset()
	w.Command(p, "score")
	got := buf.String()
	if !strings.Contains(got, "Character points: 0") {
		t.Errorf("sheet should show 0 points; got:\n%s", got)
	}
	if strings.Contains(got, "to spend") {
		t.Errorf("no spend note when 0 points; got:\n%s", got)
	}

	// Some points → bold spend note.
	p.StatPoints = 3
	buf.Reset()
	w.Command(p, "score")
	got = buf.String()
	if !strings.Contains(got, "Character points: 3") {
		t.Errorf("sheet should show the integer; got:\n%s", got)
	}
	if !strings.Contains(got, "You have character points to spend.") {
		t.Errorf("expected the bold spend note; got:\n%s", got)
	}
}
