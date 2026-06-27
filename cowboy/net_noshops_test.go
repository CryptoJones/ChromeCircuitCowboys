package cowboy

import "testing"

// The Net should have no shops or medics — you gear up in meatspace.
func TestNoShopsOrMedicsInTheNet(t *testing.T) {
	w := NewWorld(NewMemStore())
	for id, r := range w.rooms {
		realm, _ := areaInfo(id)
		if realm != "net" {
			continue
		}
		if r.Vendor {
			t.Errorf("Net room %s (%s) is a vendor — no shops allowed in the Net", id, r.Name)
		}
		if r.Medic {
			t.Errorf("Net room %s (%s) is a medic — no medics allowed in the Net", id, r.Name)
		}
	}
}
