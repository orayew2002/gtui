package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/orayew2002/gtui"
	"github.com/orayew2002/gtui/src/colors"
)

// app is the single application handle every page uses.
var app *gtui.App

func main() {
	app = gtui.Init().Theme(colors.Dark)
	app.View(shell)

	app.On("click", onNavigate)
	app.Bind("ctrl+p", openPalette)
	app.Animate(8, onTick)

	app.Select("#nav" + current().Key)

	if check() {
		return
	}
	if snapshot() {
		return
	}

	app.Render()
}

// snapshot prints one frame as text instead of starting the terminal, which
// makes the gallery easy to check from a script: go run ./app snapshot 110 30
func snapshot() bool {
	if len(os.Args) != 4 || os.Args[1] != "snapshot" {
		return false
	}

	width, _ := strconv.Atoi(os.Args[2])
	height, _ := strconv.Atoi(os.Args[3])

	if os.Getenv("MODE") == "light" {
		app.Theme(colors.Light)
	}
	if page := os.Getenv("PAGE"); page != "" {
		goTo(page)
	}

	app.Refresh()
	app.Select("#nav" + current().Key)

	fmt.Println(app.Snapshot(width, height))
	return true
}
