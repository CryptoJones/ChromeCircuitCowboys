package tests

import (
	"testing"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// TestResleeveBuff checks a fresh sleeve wakes with a +15% HP buffer (overheal
// above max) after a flatline (#26).
func TestResleeveBuff(t *testing.T) {
	w := cowboy.NewWorld(cowboy.NewMemStore())
	w.SetRoll(alwaysHit)
	out, _ := sink()
	p := w.Connect("Case", out)

	// Pick a fight we'll lose: drop to 1 HP in the back alley and let the ganger
	// flatline us, triggering the re-sleeve.
	p.RoomID = "back_alley"
	p.HP = 1
	w.Command(p, "attack ganger")
	for i := 0; i < 40 && p.HP <= p.MaxHP; i++ {
		w.Tick()
		if p.HP == 1 { // still alive and untouched? force the loss by staying at 1
			p.HP = 1
		}
	}

	if p.HP <= p.MaxHP {
		t.Fatalf("after re-sleeve, HP should exceed max (the +15%% buffer): HP=%d MaxHP=%d", p.HP, p.MaxHP)
	}
	want := p.MaxHP + p.MaxHP*15/100
	if p.HP != want {
		t.Errorf("re-sleeve HP = %d, want %d (max + 15%%)", p.HP, want)
	}
}
