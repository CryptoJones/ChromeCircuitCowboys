package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestColorblindTheme checks switching to the colorblind theme remaps the
// success/danger colors at send-time (#38).
func TestColorblindTheme(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out)

	// Default: red danger uses the standard red SGR (1;31m).
	buf.Reset()
	w.Command(p, "score") // emits some default-colored text
	if !strings.Contains(buf.String(), "\x1b[1;3") {
		t.Fatalf("expected ANSI color in output")
	}

	// Switch to colorblind-dark. The WHOLE palette is remapped (not just
	// green/red) so the entire UI — the cyan-bordered MAP included — visibly
	// shifts: danger→orange (208), success→blue (27), system→azure (39).
	buf.Reset()
	w.Command(p, "theme cbdark")
	got := buf.String()
	if p.Theme != "cbdark" {
		t.Errorf("theme should be set to cbdark, got %q", p.Theme)
	}
	if strings.Contains(got, "\x1b[1;31m") {
		t.Errorf("colorblind theme should remap the default red away; got raw red in:\n%q", got)
	}
	if strings.Contains(got, "\x1b[1;36m") {
		t.Errorf("colorblind theme should remap the default cyan (system/map borders) away; got raw cyan in:\n%q", got)
	}
	if !strings.Contains(got, "\x1b[38;5;208m") {
		t.Errorf("colorblind-dark should render danger as orange (208); got:\n%q", got)
	}
	if !strings.Contains(got, "\x1b[38;5;39m") {
		t.Errorf("colorblind-dark should remap system/map borders to azure (39); got:\n%q", got)
	}
}
