package cowboy

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestGenQuestDoc writes docs/quests.md (the quest-giver + bounty reference) from
// the live quest table when GEN_ROOM_MAP=1, so the doc can't drift from the code.
func TestGenQuestDoc(t *testing.T) {
	if os.Getenv("GEN_ROOM_MAP") == "" {
		t.Skip("set GEN_ROOM_MAP=1 to regenerate docs/quests.md")
	}
	tmpls := buildMobTemplates()
	rooms := buildRooms()
	target := func(q Quest) string {
		if mt := tmpls[q.Target]; mt != nil {
			return mt.Name
		}
		return q.Target
	}
	giver := func(q Quest) string {
		if q.Giver == "" {
			return "street brokers (Chrome Rose / Night Market)"
		}
		name := q.GiverName
		if name == "" {
			name = "a fixer"
		}
		room := q.Giver
		if r := rooms[q.Giver]; r != nil {
			room = r.Name
		}
		return fmt.Sprintf("%s — `%s` (%s)", name, q.Giver, room)
	}

	var b strings.Builder
	b.WriteString("# Chrome Circuit Cowboys — Quest Givers & Bounties\n\n")
	b.WriteString("_Generated from `cowboy/quests.go`. Accept a bounty in its giver's room (`ACCEPT <#>`), ")
	b.WriteString("kill the target, then `CLAIM` at any broker (vendor room). Bounties are repeatable._\n\n")

	section := func(title string, pred func(Quest) bool) {
		fmt.Fprintf(&b, "## %s\n\n", title)
		for _, q := range quests {
			if !pred(q) {
				continue
			}
			fmt.Fprintf(&b, "- **%s** _(L%d+)_ — %s\n", q.Name, q.MinLevel, q.Desc)
			fmt.Fprintf(&b, "    - giver: %s\n", giver(q))
			fmt.Fprintf(&b, "    - target: %s ×%d · reward: +%d XP, €$%d\n", target(q), q.Count, q.XP, q.Eddies)
		}
		b.WriteString("\n")
	}
	section("Street bounties", func(q Quest) bool { return q.Giver == "" })
	section("Meatspace — the underground descent", func(q Quest) bool { return strings.HasPrefix(q.ID, "ug") })
	section("Netspace — the Net ascent", func(q Quest) bool { return strings.HasPrefix(q.ID, "net") })
	b.WriteString("*Proudly Made in Nebraska. Go Big Red! 🌽 <https://xkcd.com/2347/>*\n")

	if err := os.WriteFile("../docs/quests.md", []byte(b.String()), 0644); err != nil {
		t.Fatalf("write quests doc: %v", err)
	}
	t.Logf("wrote ../docs/quests.md")
}

// TestQuestsWellFormed guards every bounty: its Target must be a real mob
// template (so the kill can be credited) and its Giver, if set, must be a real
// room. This catches dead quests — e.g. a target mob that was renamed or removed.
func TestQuestsWellFormed(t *testing.T) {
	tmpls := buildMobTemplates()
	rooms := buildRooms()
	seen := map[string]bool{}
	for _, q := range quests {
		if seen[q.ID] {
			t.Errorf("duplicate quest id %q", q.ID)
		}
		seen[q.ID] = true
		if tmpls[q.Target] == nil {
			t.Errorf("quest %q targets missing mob %q", q.ID, q.Target)
		}
		if q.Count < 1 {
			t.Errorf("quest %q has Count %d", q.ID, q.Count)
		}
		if q.Giver != "" && rooms[q.Giver] == nil {
			t.Errorf("quest %q giver room %q does not exist", q.ID, q.Giver)
		}
	}
}
