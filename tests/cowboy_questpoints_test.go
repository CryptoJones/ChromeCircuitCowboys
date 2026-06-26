package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestQuestAwardsPoint checks an arc-climax bounty grants a character point on
// claim (#32).
func TestQuestAwardsPoint(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	out, buf := sink()
	p := w.Connect("Case", out)

	// ug1_heist (Jax @ z1_11) awards a point. Mark it complete and claim at the
	// giver. Use a high level so the XP reward won't trigger a level-up (which
	// would grant its own points and confound the count).
	p.Level = 50
	p.Quests["ug1_heist"] = 1
	p.RoomID = "z1_11"
	before := p.StatPoints
	buf.Reset()
	w.Command(p, "claim")
	if p.StatPoints != before+1 {
		t.Errorf("claiming a Points bounty should grant a character point: %d -> %d", before, p.StatPoints)
	}
	if !strings.Contains(buf.String(), "character point") {
		t.Errorf("claim should mention the character point reward; got:\n%s", buf.String())
	}
}
