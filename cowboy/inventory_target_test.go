package cowboy

import (
	"strings"
	"testing"
)

func TestResolveInvItem(t *testing.T) {
	w := NewWorld(NewMemStore())
	p, _ := newTestPlayer(w, "Tester", startRoom)
	p.Inv = map[string]int{"medkit": 2, "ammo": 1, "ram-chip": 1}
	// sortedInv is alphabetical: ammo(1), medkit(2), ram-chip(3).

	cases := map[string]string{
		"1":        "ammo",
		"2":        "medkit",
		"3":        "ram-chip",
		"medkit":   "medkit", // a name passes through untouched
		"MEDKIT":   "medkit", // lowercased
		"9":        "9",      // out of range falls through as its digits
		"0":        "0",      // 0 is not a valid 1-based index
		"medkit 2": "medkit 2",
	}
	for token, want := range cases {
		if got := resolveInvItem(p, token); got != want {
			t.Errorf("resolveInvItem(%q) = %q, want %q", token, got, want)
		}
	}
}

func TestDropByNumber(t *testing.T) {
	w := NewWorld(NewMemStore())
	// Use a non-private room so floor piles persist (capsule is private/isolated).
	p, drain := newTestPlayer(w, "Tester", "ic_1")
	p.Inv = map[string]int{"alpha": 1, "beta": 1} // sorted: alpha(1), beta(2)

	w.drop(p, "2") // should drop beta
	out := drain()
	if !strings.Contains(out, "beta") {
		t.Errorf("DROP 2 should drop beta, output: %q", out)
	}
	if p.Inv["beta"] != 0 {
		t.Errorf("beta should be gone from the pack, still have %d", p.Inv["beta"])
	}
	if p.Inv["alpha"] != 1 {
		t.Errorf("alpha should be untouched, have %d", p.Inv["alpha"])
	}
	if w.floor["ic_1"]["beta"] != 1 {
		t.Errorf("beta should be on the floor, floor has %d", w.floor["ic_1"]["beta"])
	}
}

func TestGiveByNumber(t *testing.T) {
	w := NewWorld(NewMemStore())
	giver, gdrain := newTestPlayer(w, "Giver", "ic_1")
	receiver, _ := newTestPlayer(w, "Receiver", "ic_1")
	giver.Inv = map[string]int{"alpha": 1, "beta": 1} // sorted: alpha(1), beta(2)

	w.give(giver, "1 Receiver") // give alpha to Receiver
	out := gdrain()
	if !strings.Contains(out, "alpha") {
		t.Errorf("GIVE 1 Receiver should hand over alpha, output: %q", out)
	}
	if giver.Inv["alpha"] != 0 {
		t.Errorf("giver should no longer have alpha, have %d", giver.Inv["alpha"])
	}
	if receiver.Inv["alpha"] != 1 {
		t.Errorf("receiver should have alpha, have %d", receiver.Inv["alpha"])
	}
}
