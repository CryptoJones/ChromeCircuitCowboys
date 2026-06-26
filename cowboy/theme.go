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
	// Colorblind, dark background: green→blue, red→orange; keep cyan/yellow/magenta.
	"cbdark": strings.NewReplacer(
		green, "\x1b[1;34m", // success → bright blue
		red, "\x1b[38;5;208m", // danger → orange
	),
	// Colorblind, light background: darker, non-bold hues that read on a light terminal.
	"cblight": strings.NewReplacer(
		neon, "\x1b[36m",
		hot, "\x1b[35m",
		gold, "\x1b[38;5;94m", // dark amber
		green, "\x1b[34m", // success → dark blue
		red, "\x1b[38;5;166m", // danger → dark orange
		dim, "\x1b[38;5;240m", // darker grey for light bg
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
