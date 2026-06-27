package cowboy

import "testing"

func TestDeepestZoneHasForwardGuidance(t *testing.T) {
	w := NewWorld(NewMemStore())

	// The deepest band (z10, L91-99) has no harder area...
	if _, ok := w.onwardStep("z10_01", true); ok {
		t.Skip("z10 unexpectedly has a harder area; layout changed")
	}
	// ...so MAP must fall back to pointing at the zone objective instead of going silent.
	dir, ok := w.zoneObjectiveStep("z10_01")
	if !ok {
		t.Fatal("deepest zone should still point toward its final objective")
	}
	// dir is either a real move or "" (already at the objective room).
	if dir != "" {
		if _, valid := dirShort[dir]; !valid {
			t.Errorf("objective direction %q is not a known move", dir)
		}
	}

	// A non-deepest zone keeps using PROCEED (objective fallback is only for the end).
	if _, ok := w.onwardStep("z1_01", true); !ok {
		t.Error("a shallow zone should still have a harder area to PROCEED to")
	}
}
