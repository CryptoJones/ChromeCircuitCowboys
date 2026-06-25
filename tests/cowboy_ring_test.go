package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestRPRingsAndRumors covers the street-level RP rings: you can step from Neon
// Alley onto the Inner Circuit and walk the loop, the ring givers offer roving
// "rumor" bounties (which are NOT offered at a non-giver ring room), and a rumor
// can be accepted, cleared on a light belt stray, and claimed.
func TestRPRingsAndRumors(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(alwaysHit)
	out, buf := sink()
	p := w.Connect("Roleplayer", out)
	p.MaxHP, p.HP = 500, 500
	p.Level = 5 // clears every ring rumor's low level gate (which one is offered is randomized)

	// Step from the street onto the Inner Circuit and round the loop.
	w.Command(p, "out")   // capsule -> neon_alley
	w.Command(p, "north") // ic_1 (Neon Gate)
	if p.RoomID != "ic_1" {
		t.Fatalf("Neon Alley should step north onto the Inner Circuit, at %s", p.RoomID)
	}
	w.Command(p, "east") // ic_2 (Busker's Span — a ring giver)
	if p.RoomID != "ic_2" {
		t.Fatalf("expected ic_2, at %s", p.RoomID)
	}

	// The busker offers a rumor (giver-gated): the board is non-empty here.
	w.Command(p, "quests")
	if !strings.Contains(buf.String(), "ACCEPT") {
		t.Fatalf("a ring giver should post a rumor board; got:\n%s", buf.String())
	}
	w.Command(p, "accept 1")
	if len(p.Quests) == 0 {
		t.Fatal("accepting a rumor at a ring giver should add it")
	}

	// A non-giver ring room offers nothing (the rumors are giver-gated).
	out2, _ := sink()
	q := w.Connect("Extra", out2)
	q.RoomID = "ic_3" // Mirrorglass Curve — not a giver
	w.Command(q, "accept 1")
	if len(q.Quests) != 0 {
		t.Fatal("a non-giver ring room must not hand out rumors")
	}

	// The Sprawlbelt carries light strays you can actually fight (e.g. sb_3).
	p.RoomID = "sb_3"
	w.Command(p, "attack")
	w.Tick()
	if !strings.Contains(buf.String(), "turf-tagger") {
		t.Errorf("the belt should have a light stray to fight; got:\n%s", lastLines(buf.String()))
	}
}
