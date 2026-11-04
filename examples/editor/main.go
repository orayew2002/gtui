// A file tree beside code with a gutter and an active line.
//
// The screen itself lives in demos/editor.go, so the gallery can open the same
// app on its projects page.
package main

import (
	"github.com/orayew2002/gtui"
	"github.com/orayew2002/gtui/demos"
)

func main() {
	app := gtui.Init()

	app.View(func() gtui.Component { return demos.Editor(app, 0) })
	app.QuitOn("q", "ctrl+c")
	app.Render()
}
