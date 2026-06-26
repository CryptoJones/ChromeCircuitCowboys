package cowboy

import (
	"strings"
	"testing"
)

func botsIn(w *World) []*Player {
	var out []*Player
	for _, p := range w.players {
		if p.IsBot {
			out = append(out, p)
		}
	}
	return out
}

func TestEnableBotsSeedsRunners(t *testing.T) {
	w := NewWorld(NewMemStore())
	w.EnableBots(4)

	bots := botsIn(w)
	if len(bots) != 4 {
		t.Fatalf("expected 4 bots, got %d", len(bots))
	}
	for _, b := range bots {
		if w.byName[strings.ToLower(b.Name)] != b {
			t.Errorf("bot %q not registered in byName", b.Name)
		}
		if realm, _ := areaInfo(b.RoomID); realm != "city" {
			t.Errorf("bot %q spawned outside the surface, in realm %q (%s)", b.Name, realm, b.RoomID)
		}
		if r := w.room(b.RoomID); r == nil || r.Private {
			t.Errorf("bot %q spawned in a missing/private room %s", b.Name, b.RoomID)
		}
	}
	// n<=0 is a no-op; cap at roster size.
	w2 := NewWorld(NewMemStore())
	w2.EnableBots(0)
	if len(botsIn(w2)) != 0 {
		t.Error("EnableBots(0) should seed no bots")
	}
	w3 := NewWorld(NewMemStore())
	w3.EnableBots(999)
	if len(botsIn(w3)) != len(botRoster) {
		t.Errorf("EnableBots(999) should cap at roster size %d, got %d", len(botRoster), len(botsIn(w3)))
	}
}

func TestBotsAreNeverSaved(t *testing.T) {
	store := NewMemStore()
	w := NewWorld(store)
	w.EnableBots(3)
	w.SaveAll() // must skip bots

	for _, b := range botsIn(w) {
		if _, ok, _ := store.Load(b.Name); ok {
			t.Errorf("bot %q was persisted to the store; bots must be ephemeral", b.Name)
		}
	}
}

func TestBotChatterOnlyNearRealRunner(t *testing.T) {
	w := NewWorld(NewMemStore())
	w.SetRoll(func(n int) int { return 0 })
	w.EnableBots(1)
	b := botsIn(w)[0]

	// Alone: chatter produces nothing.
	w.botChatter(b)
	if w.realRunnerNear(b) {
		t.Fatal("no real runner placed yet, but realRunnerNear is true")
	}

	// Drop a real runner into the bot's room: now chatter reaches them.
	p, drain := newTestPlayer(w, "Human", b.RoomID)
	_ = p
	if !w.realRunnerNear(b) {
		t.Fatal("real runner shares the room but realRunnerNear is false")
	}
	w.botChatter(b)
	if out := drain(); out == "" {
		t.Error("bot should chatter when a real runner is present")
	}
}

func TestBotWanderStaysOnSurface(t *testing.T) {
	w := NewWorld(NewMemStore())
	w.SetRoll(func(n int) int { return 0 }) // first exit each time
	w.EnableBots(3)

	// Drive several wander steps; bots must never leave the city realm.
	for step := 0; step < 20; step++ {
		for _, b := range botsIn(w) {
			w.botWander(b)
			if realm, _ := areaInfo(b.RoomID); realm != "city" {
				t.Fatalf("bot %q wandered off the surface into %q (%s)", b.Name, realm, b.RoomID)
			}
		}
	}
}
