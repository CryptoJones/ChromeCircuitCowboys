package cowboy

import (
	"strings"
	"testing"
)

func TestStepToRooms(t *testing.T) {
	w := NewWorld(NewMemStore())
	// Already at a goal → ok with empty dir ("here").
	if dir, ok := w.stepToRooms("z1_01", map[string]bool{"z1_01": true}); !ok || dir != "" {
		t.Errorf("start==goal should be ok with empty dir, got (%q,%v)", dir, ok)
	}
	// A reachable goal one step away returns a real direction.
	start := "z1_01"
	r := w.room(start)
	if r == nil || len(r.Exits) == 0 {
		t.Skip("no exits from test room")
	}
	var nbr, wantDir string
	for _, d := range mapDirs {
		if dest, ok := r.Exits[d]; ok {
			nbr, wantDir = dest, d
			break
		}
	}
	if dir, ok := w.stepToRooms(start, map[string]bool{nbr: true}); !ok || dir != wantDir {
		t.Errorf("step to adjacent %s should be %q, got (%q,%v)", nbr, wantDir, dir, ok)
	}
	// Unreachable / empty goals → not ok.
	if _, ok := w.stepToRooms(start, map[string]bool{"no_such_room": true}); ok {
		t.Error("unknown goal room should be unreachable")
	}
	if _, ok := w.stepToRooms(start, map[string]bool{}); ok {
		t.Error("empty goal set should be not-ok")
	}
}

func TestQuestDirectionPointsToGiverWhenReady(t *testing.T) {
	w := NewWorld(NewMemStore())
	q, ok := questByID("ug1_snatch") // Giver: z1_01
	if !ok {
		t.Skip("quest fixture missing")
	}
	// Stand at the giver room → READY should say turn in HERE.
	p, _ := newTestPlayer(w, "Boss", q.Giver)
	if got := w.questDirection(p, q, true); !strings.Contains(got, "HERE") {
		t.Errorf("at the giver room, READY hint should say HERE, got %q", got)
	}
	// Stand somewhere else → READY should point a direction to CLAIM.
	p.RoomID = "z1_05"
	if got := w.questDirection(p, q, true); !strings.Contains(got, "CLAIM") {
		t.Errorf("away from the giver, READY hint should point to CLAIM, got %q", got)
	}
	// Not ready → hint should point toward the target (no CLAIM wording).
	if got := w.questDirection(p, q, false); strings.Contains(got, "CLAIM") {
		t.Errorf("in-progress hint should point at the target, not CLAIM, got %q", got)
	}
}
