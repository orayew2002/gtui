// Tabs and links that swap one box, which is the routing model.
//
// The screen itself lives in demos/browser.go, so the gallery can open the same
// app on its projects page.
package main

import (
	"github.com/orayew2002/gtui"
	"github.com/orayew2002/gtui/demos"
)

func main() {
	app := gtui.Init()

	app.View(func() gtui.Component { return demos.Browser(app, 0) })
	app.QuitOn("q", "ctrl+c")
	app.Render()
}
