package cowboy

import "testing"

func TestParqueCentralEightSpokeHub(t *testing.T) {
	w := NewWorld(NewMemStore())
	park := w.room("parque_central")
	if park == nil {
		t.Fatal("parque_central room should exist")
	}
	// The park fans OUT the eight compass ways to its eight gate rooms...
	spokes := map[string]string{
		"north": "ic_1", "east": "ic_3", "south": "ic_4", "west": "ic_6",
		"northeast": "sb_2", "southeast": "sb_5", "southwest": "sb_8", "northwest": "sb_10",
	}
	if len(park.Exits) != len(spokes) {
		t.Errorf("park should have exactly %d exits, has %d: %v", len(spokes), len(park.Exits), park.Exits)
	}
	for dir, gate := range spokes {
		if park.Exits[dir] != gate {
			t.Errorf("park %s should lead to %s, got %q", dir, gate, park.Exits[dir])
		}
		// ...and each gate opens IN toward the centre.
		if g := w.room(gate); g == nil || g.Exits["in"] != "parque_central" {
			t.Errorf("gate %s should have IN -> parque_central", gate)
		}
	}
	// A non-gate ring room does NOT open in (you walk the ring to a gate).
	if r := w.room("sb_7"); r != nil && r.Exits["in"] != "" {
		t.Errorf("non-gate sb_7 should not open IN, got %q", r.Exits["in"])
	}
}

func TestDiagonalDirections(t *testing.T) {
	for sh, full := range map[string]string{"ne": "northeast", "se": "southeast", "sw": "southwest", "nw": "northwest"} {
		if dirAliases[sh] != full {
			t.Errorf("alias %q should map to %q, got %q", sh, full, dirAliases[sh])
		}
	}
	if opposite("northeast") != "southwest" || opposite("southeast") != "northwest" {
		t.Error("diagonal opposites are wrong")
	}
	// You can actually walk a diagonal spoke and come back IN.
	w := NewWorld(NewMemStore())
	p, _ := newTestPlayer(w, "Boss", "parque_central")
	w.move(p, "northeast")
	if p.RoomID != "sb_2" {
		t.Fatalf("NE from the park should reach sb_2, in %s", p.RoomID)
	}
	w.move(p, "in")
	if p.RoomID != "parque_central" {
		t.Errorf("IN from a gate should return to the park, in %s", p.RoomID)
	}
}

func TestWarmPulseOneWayGate(t *testing.T) {
	w := NewWorld(NewMemStore())
	p, _ := newTestPlayer(w, "Boss", "z10_14")

	// Going UP doesn't move you — it arms a confirmation.
	w.move(p, "up")
	if p.RoomID != "z10_14" {
		t.Fatalf("UP should not move yet (awaiting confirm), in %s", p.RoomID)
	}
	if p.confirmExit != "z10_14" {
		t.Fatalf("UP should arm the departure confirm, got %q", p.confirmExit)
	}

	// NO cancels and keeps you put.
	w.cancelDepart(p)
	if p.confirmExit != "" || p.RoomID != "z10_14" {
		t.Fatalf("NO should cancel and keep you in z10_14")
	}

	// Re-arm, then YES makes the one-way trip to the park.
	w.move(p, "up")
	w.confirmDepart(p)
	if p.RoomID != "parque_central" {
		t.Errorf("YES should surface you in parque_central, in %s", p.RoomID)
	}
	if p.confirmExit != "" {
		t.Error("confirm should be cleared after departure")
	}
	// One-way: the park has no exit back to The Warm Pulse.
	for _, dest := range w.room("parque_central").Exits {
		if dest == "z10_14" {
			t.Error("parque_central must not lead back to the Undercity")
		}
	}
}

func TestOrdinaryMoveCancelsPendingDepart(t *testing.T) {
	w := NewWorld(NewMemStore())
	p, _ := newTestPlayer(w, "Boss", "z10_14")
	w.move(p, "up") // arm
	if p.confirmExit == "" {
		t.Fatal("expected armed confirm")
	}
	// Move back toward the boss room (the real cardinal exit) — should clear the confirm.
	var back string
	for d, dest := range w.room("z10_14").Exits {
		if dest == "z10_13" {
			back = d
		}
	}
	if back == "" {
		t.Skip("no back-exit found")
	}
	w.move(p, back)
	if p.confirmExit != "" {
		t.Error("an ordinary move should cancel the pending departure")
	}
}
