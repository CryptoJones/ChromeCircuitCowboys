package tests

import (
	"bufio"
	"strings"
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

func TestCowboyCharacterCreation(t *testing.T) {
	// Class 2 = Enforcer (B13 R11 I6); spend 3 into BODY, 1 into REFLEXES, rest (2) to INT.
	in := bufio.NewReader(strings.NewReader("2\r\n3\r\n1\r\n"))
	var out strings.Builder
	spec, err := cowboy.RunCharacterCreation(in, func(s string) { out.WriteString(s) })
	if err != nil {
		t.Fatal(err)
	}
	if spec.ClassID != "enforcer" {
		t.Fatalf("class = %q, want enforcer", spec.ClassID)
	}
	if spec.Body != 16 || spec.Reflexes != 12 || spec.Intelligence != 8 {
		t.Fatalf("stats = B%d R%d I%d, want B16 R12 I8", spec.Body, spec.Reflexes, spec.Intelligence)
	}
	s := out.String()
	if !strings.Contains(s, "CHARACTER CREATION") || !strings.Contains(s, "Hacker") || !strings.Contains(s, "Enforcer") {
		t.Error("creation screen should list the classes")
	}
}

func TestCowboyCreationRejectsGarbageThenAccepts(t *testing.T) {
	// Garbage class + non-numeric point entries must RE-PROMPT (never silently
	// pick a class or eat points), then accept the valid follow-up.
	// class: 99 (out of range) -> z (invalid) -> 2 (Enforcer)
	// BODY:  h (invalid) -> 3        REFLEXES: y (invalid) -> 1
	in := bufio.NewReader(strings.NewReader("99\r\nz\r\n2\r\nh\r\n3\r\ny\r\n1\r\n"))
	var out strings.Builder
	spec, err := cowboy.RunCharacterCreation(in, func(s string) { out.WriteString(s) })
	if err != nil {
		t.Fatal(err)
	}
	if spec.ClassID != "enforcer" {
		t.Fatalf("class = %q, want enforcer after re-prompts", spec.ClassID)
	}
	// Enforcer base B13 R11 I6 + (3 BODY, 1 REFLEXES, 2 leftover INT).
	if spec.Body != 16 || spec.Reflexes != 12 || spec.Intelligence != 8 {
		t.Fatalf("stats = B%d R%d I%d, want B16 R12 I8", spec.Body, spec.Reflexes, spec.Intelligence)
	}
	if !strings.Contains(out.String(), "Enter a number") {
		t.Error("invalid input should have re-prompted")
	}
}

func TestCowboyCreationQuit(t *testing.T) {
	// Typing Q at any creation prompt jacks out (ErrQuit), which the caller turns
	// into a clean disconnect.
	for _, seq := range []string{"q\r\n", "2\r\nq\r\n", "2\r\n3\r\nquit\r\n"} {
		in := bufio.NewReader(strings.NewReader(seq))
		if _, err := cowboy.RunCharacterCreation(in, func(string) {}); err != cowboy.ErrQuit {
			t.Fatalf("seq %q: want ErrQuit, got %v", seq, err)
		}
	}
}

func TestCowboyCreateCharacterAppliesClassAndPersists(t *testing.T) {
	store := cowboy.NewMemStore()
	w := cowboy.NewWorld(store)
	out, _ := sink()
	spec := cowboy.CharSpec{ClassID: "enforcer", Body: 16, Reflexes: 12, Intelligence: 8}

	if w.HasCharacter("Rogue") {
		t.Fatal("brand-new name should not exist yet")
	}
	p := w.CreateCharacter("Rogue", spec, out)
	if p.Class != "Enforcer" || p.Body != 16 || p.Reflexes != 12 || p.Intelligence != 8 {
		t.Fatalf("class/stats not applied: %+v", p)
	}
	if p.HP <= 40 { // higher Body -> higher MaxHP than the default 40
		t.Fatalf("MaxHP should reflect higher Body, got HP=%d", p.HP)
	}
	// Persisted immediately, and a later login reloads the class.
	if !w.HasCharacter("Rogue") {
		t.Fatal("created character should be saved")
	}
	w.Disconnect(p)
	w2 := cowboy.NewWorld(store)
	out2, _ := sink()
	p2 := w2.Connect("Rogue", out2)
	if p2.Class != "Enforcer" || p2.Body != 16 {
		t.Fatalf("returning login lost class/stats: %+v", p2)
	}
}
