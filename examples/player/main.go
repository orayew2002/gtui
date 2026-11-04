// A music player with a queue.
//
// The screen itself lives in demos/player.go, so the gallery can open the same
// app on its projects page.
package main

import (
	"github.com/orayew2002/gtui"
	"github.com/orayew2002/gtui/demos"
)

func main() {
	app := gtui.Init()

	// the project draws itself from the tick it is handed, so the timer that
	// moves it belongs to whoever is hosting it
	tick := 0

	app.View(func() gtui.Component { return demos.Player(app, tick) })
	app.Animate(8, func(at int) {
		tick = at
		app.Refresh()
	})

	app.QuitOn("q", "ctrl+c")
	app.Render()
}
