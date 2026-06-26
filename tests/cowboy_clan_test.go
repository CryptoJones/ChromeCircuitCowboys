package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestClanAndRewardBonus checks clan create/join + that clanmates partying
// together earn the boosted reward multiplier (#44).
func TestClanAndRewardBonus(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	o1, b1 := sink()
	o2, _ := sink()
	a := w.Connect("Rook", o1)
	b := w.Connect("Jett", o2)

	// Both join the same clan.
	w.Command(a, "clan create Chrome Reapers")
	w.Command(b, "clan join Chrome Reapers")
	if a.Clan != "Chrome Reapers" || b.Clan != "Chrome Reapers" {
		t.Fatalf("clan join failed: a=%q b=%q", a.Clan, b.Clan)
	}

	// Solo (party nil) → no bonus.
	if pct := w.RewardPct(a); pct != 100 {
		t.Errorf("solo reward pct should be 100, got %d", pct)
	}

	// Party up in the same room → clan bonus applies (115 * 180% = 207).
	a.RoomID, b.RoomID = "back_alley", "back_alley"
	w.Command(a, "invite Jett")
	w.Command(b, "accept")
	if pct := w.RewardPct(a); pct != 207 {
		t.Errorf("clanmate party reward pct should be 207 (115*1.8), got %d", pct)
	}

	// CLAN with no arg shows the clan.
	b1.Reset()
	w.Command(a, "clan")
	if !strings.Contains(b1.String(), "Chrome Reapers") {
		t.Errorf("CLAN should show your clan; got:\n%s", b1.String())
	}
}
