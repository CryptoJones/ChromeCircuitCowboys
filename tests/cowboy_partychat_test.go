package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestPartyChatShortcut checks the ";" prefix routes to crew chat (#43).
func TestPartyChatShortcut(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out)

	// Solo: ";hi" routes to group chat, which reports no crew (proves wiring).
	buf.Reset()
	w.Command(p, ";hi crew")
	if !strings.Contains(buf.String(), "no crew") {
		t.Errorf("; prefix should route to crew chat; got:\n%s", buf.String())
	}
	if strings.Contains(buf.String(), "Unknown command") {
		t.Errorf("; prefix must not be an unknown command; got:\n%s", buf.String())
	}
}
