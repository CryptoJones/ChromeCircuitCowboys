package tests

import "strings"

// lastLines returns the final few lines of output (handy in failure messages).
func lastLines(s string) string {
	parts := strings.Split(strings.TrimRight(s, "\r\n"), "\n")
	if len(parts) > 6 {
		parts = parts[len(parts)-6:]
	}
	return strings.Join(parts, "\n")
}
