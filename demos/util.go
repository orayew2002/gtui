package demos

import "github.com/orayew2002/gtui"

// clamp keeps v inside lo and hi, which several projects need when a key or a
// click moves a position.
func clamp(v, lo, hi int) int { return min(max(v, lo), hi) }

// narrow reports whether the screen is too tight to put panels side by side.
// The gallery's projects page hands a project less room than the whole
// terminal, so the threshold is deliberately generous.
func narrow(a *gtui.App) bool {
	width, _ := a.Size()
	return width < 100
}
