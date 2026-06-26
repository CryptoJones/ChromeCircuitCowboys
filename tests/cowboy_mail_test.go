package tests

import (
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestTerminalMailAndWire checks SEND/MAIL and WIRE between players at a data
// terminal, including an offline recipient (#24).
func TestTerminalMailAndWire(t *testing.T) {
	st := cowboy.NewMemStore()
	w := cowboy.NewWorld(st)

	// Seed Jett with 100 scrip, then take them offline (persists).
	oj, _ := sink()
	jett := w.Connect("Jett", oj)
	jett.Eddies = 100
	w.Disconnect(jett)

	// Rook, at the Night Market terminal, mails + wires the offline Jett.
	or, rbuf := sink()
	rook := w.Connect("Rook", or)
	rook.RoomID = "market" // vendor → terminal
	rook.Eddies = 500
	w.Command(rook, "send Jett meet me at the Night Market")
	w.Command(rook, "wire Jett 150")
	if rook.Eddies != 350 {
		t.Fatalf("wire should debit the sender 150: have %d", rook.Eddies)
	}

	// Jett logs back in: scrip credited, message delivered on login.
	o4, buf4 := sink()
	jett2 := w.Connect("Jett", o4)
	if jett2.Eddies != 250 {
		t.Errorf("offline wire should credit Jett: have %d want 250", jett2.Eddies)
	}
	if !strings.Contains(buf4.String(), "Night Market") {
		t.Errorf("Jett should receive the queued message on login; got:\n%s", buf4.String())
	}

	// A terminal is required: off one, SEND is refused.
	rook.RoomID = "back_alley"
	rbuf.Reset()
	w.Command(rook, "send Jett hi")
	if !strings.Contains(rbuf.String(), "need a data terminal") {
		t.Errorf("SEND off a terminal should be refused; got:\n%s", rbuf.String())
	}
}
