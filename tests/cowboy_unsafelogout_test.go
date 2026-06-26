package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestUnsafeLogoutPenalty checks logging out in an unsafe room costs 5% HP on
// return, while a safe room is free (#48).
func TestUnsafeLogoutPenalty(t *testing.T) {
	st := cowboy.NewMemStore()
	w := cowboy.NewWorld(st)

	// Save a character logged off in an UNSAFE room (the back alley).
	o1, _ := sink()
	jett := w.Connect("Jett", o1)
	jett.RoomID = "back_alley" // not Safe
	jett.HP, jett.MaxHP = 100, 100
	w.Disconnect(jett)

	o2, buf := sink()
	jett2 := w.Connect("Jett", o2)
	if jett2.HP != 95 { // 100 - 5%
		t.Errorf("unsafe logout should cost 5%% HP: have %d want 95", jett2.HP)
	}
	if !strings.Contains(buf.String(), "while you were logged off") {
		t.Errorf("expected the offline-jump message; got:\n%s", buf.String())
	}

	// A SAFE logout (the clone pod) is free.
	o3, _ := sink()
	rook := w.Connect("Rook", o3)
	rook.RoomID = "capsule" // Safe (the pod)
	rook.HP, rook.MaxHP = 100, 100
	w.Disconnect(rook)
	o4, _ := sink()
	rook2 := w.Connect("Rook", o4)
	if rook2.HP != 100 {
		t.Errorf("safe logout should not be penalized: have %d", rook2.HP)
	}
}
