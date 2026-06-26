package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestStatImplant checks a per-band stat implant is stocked, buyable, and that
// installing it grants the permanent stat boost (#15).
func TestStatImplant(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out)

	// Vance's Back-Door Clinic (z1_09) is a tier-1 zone vendor AND a Emergency
	// Medic, so it stocks the tier-1 implants and can install them.
	p.RoomID = "z1_09"
	p.Eddies = 1000
	body0, hp0 := p.Body, p.MaxHP

	w.Command(p, "buy reflex-booster")
	if p.Inv["reflex-booster"] <= 0 {
		t.Fatalf("should have bought the reflex-booster: inv=%v", p.Inv)
	}
	ref0 := p.Reflexes
	buf.Reset()
	w.Command(p, "install reflex-booster")
	if p.Reflexes != ref0+2 {
		t.Errorf("installing reflex-booster should +2 Reflexes: %d -> %d", ref0, p.Reflexes)
	}
	if !strings.Contains(buf.String(), "Reflexes +2") {
		t.Errorf("install should report the boost; got:\n%s", buf.String())
	}

	// A Body implant should also lift MaxHP.
	w.Command(p, "buy subdermal-plating")
	w.Command(p, "install subdermal-plating")
	if p.Body != body0+2 {
		t.Errorf("subdermal-plating should +2 Body: %d -> %d", body0, p.Body)
	}
	if p.MaxHP <= hp0 {
		t.Errorf("a Body implant should raise MaxHP: %d -> %d", hp0, p.MaxHP)
	}
}

// TestBandStocksImplants checks a band vendor stocks stat implants (#15).
func TestBandStocksImplants(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out)

	// z1_03 (Gutter Bazaar) is a tier-1 zone vendor.
	p.RoomID = "z1_03"
	w.Command(p, "list")
	got := buf.String()
	if !strings.Contains(got, "booster") && !strings.Contains(got, "plating") && !strings.Contains(got, "coprocessor") {
		t.Errorf("a tier-1 vendor should stock stat implants; got:\n%s", got)
	}
}
