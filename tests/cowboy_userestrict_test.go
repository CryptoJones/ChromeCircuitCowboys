package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestUseReasons checks USE explains why instead of "you can't use that", and
// class-restricted gear is blocked with a reason (#49).
func TestUseReasons(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out)

	// Trying to USE cyberware tells you to INSTALL it, not a bare refusal.
	p.Inv["ice-breaker"] = 1
	buf.Reset()
	w.Command(p, "use ice-breaker")
	got := buf.String()
	if strings.Contains(got, "You can't use that.") {
		t.Errorf("USE should explain why, not bare-refuse; got:\n%s", got)
	}
	if !strings.Contains(got, "INSTALL") {
		t.Errorf("USE of cyberware should point to INSTALL; got:\n%s", got)
	}

	// Class-restricted gear names the requirement.
	p.Class = "hacker"
	p.Inv["berserker-core"] = 1
	buf.Reset()
	w.Command(p, "use berserker-core")
	if !strings.Contains(buf.String(), "Enforcer gear") {
		t.Errorf("class-restricted USE should name the reason; got:\n%s", buf.String())
	}
}
