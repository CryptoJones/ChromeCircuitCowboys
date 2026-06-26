package tests

import (
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestPartyVsPartySpar checks a crew-vs-crew gym match: knocked-out fighters
// stay down until one crew is wiped, then everyone is restored (#45).
func TestPartyVsPartySpar(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(alwaysHit)
	mk := func(name string) *cowboy.Player {
		o, _ := sink()
		p := w.Connect(name, o)
		p.RoomID = "sb_gym"
		p.HP, p.MaxHP = 400, 400
		return p
	}
	rook := mk("Rook")
	jett := mk("Jett")
	bo := mk("Bo")
	cal := mk("Cal")

	// Two crews.
	w.Command(rook, "invite Jett")
	w.Command(jett, "accept")
	w.Command(bo, "invite Cal")
	w.Command(cal, "accept")

	knock := func(attacker, target *cowboy.Player) {
		target.HP = 1
		w.Command(attacker, "attack "+target.Name)
		for i := 0; i < 5 && !target.Downed(); i++ {
			w.Tick()
		}
	}

	// Drop Bo: team match → Bo stays down, match continues.
	knock(rook, bo)
	if !bo.Downed() {
		t.Fatalf("Bo should be downed in a team spar")
	}
	if cal.Downed() {
		t.Fatalf("match should still be live (Cal up)")
	}

	// Drop Cal: crew B wiped → match over, everyone restored.
	knock(jett, cal)
	if rook.Downed() || jett.Downed() || bo.Downed() || cal.Downed() {
		t.Errorf("match over should clear all downed flags")
	}
	if bo.HP != bo.MaxHP || cal.HP != cal.MaxHP {
		t.Errorf("match over should restore HP: bo=%d/%d cal=%d/%d", bo.HP, bo.MaxHP, cal.HP, cal.MaxHP)
	}
}
