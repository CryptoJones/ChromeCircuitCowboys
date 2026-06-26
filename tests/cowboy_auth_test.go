package tests

import (
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestPasswordAuth checks set/verify and the passwordless-legacy state that
// drives the first-login migration (#55/#56).
func TestPasswordAuth(t *testing.T) {
	st := cowboy.NewMemStore()
	w := cowboy.NewWorld(st)

	// Seed a character with no password (legacy), then take it offline.
	o, _ := sink()
	jett := w.Connect("Jett", o)
	w.Disconnect(jett)

	a := w.AuthInfo("Jett")
	if !a.Exists || a.HasPassword {
		t.Fatalf("legacy character should exist with no password: %+v", a)
	}
	if w.AuthInfo("Nobody").Exists {
		t.Error("unknown character should not exist")
	}

	// Migration: set a password.
	if err := w.SetPassword("Jett", "correct horse"); err != nil {
		t.Fatalf("SetPassword: %v", err)
	}
	if a := w.AuthInfo("Jett"); !a.HasPassword {
		t.Error("after SetPassword the character should have a password")
	}

	// Verify.
	if !w.CheckPassword("Jett", "correct horse") {
		t.Error("the correct password should verify")
	}
	if w.CheckPassword("Jett", "battery staple") {
		t.Error("a wrong password must not verify")
	}
}
