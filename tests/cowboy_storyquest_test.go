package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestStoryQuestGiverGatingAndClaim covers the giver-gated story bounties: a
// quest is only offered in its quest-giver's room, and the accept -> kill the
// named boss -> claim loop pays out.
func TestStoryQuestGiverGatingAndClaim(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(alwaysHit)
	out, buf := sink()
	p := w.Connect("Case", out)
	p.Intelligence, p.MaxHP, p.RAM = 80, 50000, 50000
	p.HP = p.MaxHP

	// The L1-10 Net bounty is offered by Fixer-7 at the first Net shell only.
	p.RoomID = "nz1_1_top"
	w.Command(p, "quests")
	if !strings.Contains(buf.String(), "Fixer-7") || !strings.Contains(buf.String(), "GigaMesh Ledger") {
		t.Fatalf("Fixer-7 should offer the GigaMesh Ledger here; got:\n%s", buf.String())
	}
	w.Command(p, "accept 1")
	if _, ok := p.Quests["net1_trace"]; !ok {
		t.Fatal("accepting at the giver should add the bounty")
	}

	// It is NOT offered elsewhere (giver-gated): a fresh runner can't accept it
	// from an unrelated room.
	out2, _ := sink()
	q := w.Connect("Molly", out2)
	q.RoomID = "nz1_2_mid"
	w.Command(q, "accept 1")
	if _, ok := q.Quests["net1_trace"]; ok {
		t.Fatal("the bounty must not be acceptable away from its giver")
	}

	// Kill the named boss (Tracewright lives in nz1_5_bot) to satisfy the bounty.
	p.RoomID = "nz1_5_bot"
	for i := 0; i < 60 && p.Quests["net1_trace"] < 1; i++ {
		w.Command(p, "attack tracewright")
		w.Tick()
	}
	if p.Quests["net1_trace"] < 1 {
		t.Fatalf("never landed the boss kill; progress=%d", p.Quests["net1_trace"])
	}

	// Claim back at a broker (the Net shell is a vendor) and get paid.
	p.RoomID = "nz1_1_top"
	xp0, scrip0 := p.XP, p.Eddies
	w.Command(p, "claim")
	if _, still := p.Quests["net1_trace"]; still {
		t.Fatal("claim should clear the completed bounty")
	}
	if p.XP <= xp0 || p.Eddies <= scrip0 {
		t.Fatalf("claim paid nothing: dXP=%d dScrip=%d", p.XP-xp0, p.Eddies-scrip0)
	}
}
