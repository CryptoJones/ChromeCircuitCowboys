package cowboy

import "strings"

// Color themes (#38). The game emits its default neon palette everywhere; for
// players who pick a colorblind-friendly scheme we remap those exact SGR codes
// to a distinguishable set at send-time — so no call site changes. Modeled on
// Claude Code's colorblind palette: success → blue, danger → orange (instead of
// the green/red that red-green colorblind players can't tell apart), in dark and
// light variants.

// recolor rewrites the default palette codes in s to the player's theme. The
// empty/"default" theme is a no-op.
func recolor(theme, s string) string {
	r := themeReplacers[theme]
	if r == nil {
		return s
	}
	return r.Replace(s)
}

var themeReplacers = map[string]*strings.Replacer{
	// Colorblind, dark background. The WHOLE palette is remapped to an
	// Okabe-Ito-inspired, colorblind-safe set — so the entire UI (the cyan-
	// bordered MAP included) visibly shifts, not just green/red text. Hues are
	// chosen so no two semantic colors collapse for red-green colorblindness:
	// system=azure, success=blue, reward=amber, alert=pink, danger=orange.
	"cbdark": strings.NewReplacer(
		neon, "\x1b[38;5;39m", // system/headers → azure (was cyan)
		hot, "\x1b[38;5;213m", // alerts → pink
		gold, "\x1b[38;5;220m", // reward → amber
		green, "\x1b[38;5;27m", // success → strong blue
		red, "\x1b[38;5;208m", // danger → orange
		dim, "\x1b[38;5;250m", // ambience → lighter grey
	),
	// Colorblind, light background. Dark, non-bold hues that read on a light
	// terminal AND stay colorblind-safe — a deliberately different (darker,
	// muted) look from both default and cbdark so the three are unmistakable.
	"cblight": strings.NewReplacer(
		neon, "\x1b[38;5;26m", // system/headers → medium blue
		hot, "\x1b[38;5;90m", // alerts → dark magenta
		gold, "\x1b[38;5;94m", // reward → dark amber
		green, "\x1b[38;5;19m", // success → deep navy
		red, "\x1b[38;5;166m", // danger → dark orange
		dim, "\x1b[38;5;240m", // ambience → dark grey for light bg
	),
}

// themeName maps a user-typed scheme to its key (and validates it).
func themeName(arg string) (key, label string, ok bool) {
	switch strings.ToLower(strings.TrimSpace(arg)) {
	case "default", "neon", "":
		return "", "default (neon)", true
	case "cbdark", "colorblind-dark", "cb-dark":
		return "cbdark", "colorblind (dark)", true
	case "cblight", "colorblind-light", "cb-light":
		return "cblight", "colorblind (light)", true
	}
	return "", "", false
}

// theme handles the THEME command: switch and persist the color scheme.
func (w *World) theme(p *Player, arg string) {
	if strings.TrimSpace(arg) == "" {
		_, label, _ := themeName(p.Theme)
		p.send(style(neon, "Color scheme: ") + label + crlf)
		p.send(style(dim, "THEME default | cbdark | cblight  (colorblind-friendly: success=blue, danger=orange)") + crlf)
		return
	}
	key, label, ok := themeName(arg)
	if !ok {
		p.send(style(dim, "Unknown scheme. THEME default | cbdark | cblight") + crlf)
		return
	}
	p.Theme = key
	w.save(p)
	p.send(style(green, "Color scheme set to "+label+".") + crlf)
	p.send(style(red, "  danger ") + style(green, "success ") + style(gold, "reward ") + style(neon, "system ") + style(dim, "(sample)") + crlf)
}
