package cowboy

import (
	"strings"
	"testing"
)

// newTestPlayer wires a player to a capturing sink so command output can be
// asserted. Returns the player and a func that drains everything sent so far.
func newTestPlayer(w *World, name, room string) (*Player, func() string) {
	var sb strings.Builder
	p := w.newPlayer(name, func(s string) { sb.WriteString(s) })
	p.RoomID = room
	w.players[p.ID] = p
	w.byName[strings.ToLower(name)] = p
	return p, func() string { out := sb.String(); sb.Reset(); return out }
}

func TestClassifyTalk(t *testing.T) {
	cases := map[string]string{
		"hi":                     "greet",
		"hey choom":              "greet",
		"hola":                   "greet",
		"later":                  "bye",
		"adios amigo":            "bye",
		"thanks":                 "thanks",
		"thank you":              "thanks",
		"you suck":               "insult",
		"fuck off":               "insult",
		"where is the fixer?":    "ask",
		"how do I get deeper":    "ask",
		"who are you":            "ask",
		"nice weather we having": "smalltalk",
		"":                       "smalltalk",
	}
	for input, want := range cases {
		if got := classifyTalk(input); got != want {
			t.Errorf("classifyTalk(%q) = %q, want %q", input, got, want)
		}
	}
	// Quoted input must classify the same as bare input.
	if got := classifyTalk(`"hi"`); got != "greet" {
		t.Errorf("classifyTalk(quoted hi) = %q, want greet", got)
	}
}

func TestTalkWithArgReplies(t *testing.T) {
	w := NewWorld(NewMemStore())
	w.SetRoll(func(n int) int { return 0 }) // always the first line, for determinism

	// In a room with a named flavor NPC (Rosa at ic_1), the reply is fronted by her.
	p, drain := newTestPlayer(w, "Tester", "ic_1")
	w.talk(p, "hello there")
	out := drain()
	if !strings.Contains(out, "Rosa") {
		t.Errorf("expected the room's named NPC (Rosa) to answer, got: %q", out)
	}
	if !strings.Contains(out, talkReplies["greet"][0]) {
		t.Errorf("expected the first greet reply, got: %q", out)
	}

	// An insult routes to the insult pool.
	w.talk(p, "you suck, clone")
	if out := drain(); !strings.Contains(out, talkReplies["insult"][0]) {
		t.Errorf("expected an insult reply, got: %q", out)
	}
}

func TestTalkBareStillGivesLore(t *testing.T) {
	w := NewWorld(NewMemStore())
	w.SetRoll(func(n int) int { return 0 })

	// Bare TALK at a flavor NPC still gives the NPC's authored patter, not a reply.
	p, drain := newTestPlayer(w, "Tester", "ic_1")
	w.talk(p, "")
	if out := drain(); !strings.Contains(out, roomNPC["ic_1"].lines[0]) {
		t.Errorf("bare TALK should give the NPC's first authored line, got: %q", out)
	}
}
