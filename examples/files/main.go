// A tree, a thousand-row list built a row at a time, split panes.
//
// The screen itself lives in demos/files.go, so the gallery can open the same
// app on its projects page.
package main

import (
	"github.com/orayew2002/gtui"
	"github.com/orayew2002/gtui/demos"
)

func main() {
	app := gtui.Init()

	app.View(func() gtui.Component { return demos.Files(app, 0) })
	app.QuitOn("q", "ctrl+c")
	app.Render()
}
