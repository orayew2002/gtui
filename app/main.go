// The gallery's entry point: build the app, load the languages, register the
// view, and either run it in the terminal or answer one of the two modes the
// scripts use, check and audit.
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
	languages()

	app.View(shell)

	app.On("click", onNavigate)
	app.Bind("ctrl+p", openPalette)

	app.Select("#nav" + current().Key)

	if check() {
		return
	}
	if audit() {
		return
	}
	if snapshot() {
		return
	}

	// the self check drives onTick itself, so the timer starts only for the
	// terminal, where it keeps the moving pages moving
	app.Animate(8, onTick)
	app.Render()
}

// languages loads app/locales when it is there, from the repo root or from the
// app directory, so the i18n page has real languages to switch between.
func languages() {
	for _, dir := range []string{"app/locales", "locales"} {
		if _, err := os.Stat(dir); err != nil {
			continue
		}
		if err := gtui.LoadLanguages(dir); err != nil {
			fmt.Fprintln(os.Stderr, "locales:", err)
		}
		gtui.SetFallback("en")
		return
	}
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
