package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/orayew2002/gtui"
	"github.com/orayew2002/gtui/src/colors"
)

var app *gtui.App

func main() {
	loadLanguages()

	app = gtui.Init().Theme(colors.Dark)
	app.View(shell)

	app.On("click", onNavigate)
	app.Bind("ctrl+p", openPalette)

	app.Select("#nav" + active().Key)
	startTicker()

	if snapshot() {
		return
	}

	app.Render()
}

func loadLanguages() {
	for _, dir := range []string{"app/locales", "locales"} {
		if err := gtui.LoadLanguages(dir); err == nil {
			gtui.SetFallback("en")
			gtui.SetLanguage("en")
			return
		}
	}

	fmt.Fprintln(os.Stderr, "gtui: locales not found")
	os.Exit(1)
}

func snapshot() bool {
	if len(os.Args) != 4 || os.Args[1] != "snapshot" {
		return false
	}

	width, _ := strconv.Atoi(os.Args[2])
	height, _ := strconv.Atoi(os.Args[3])

	if code := os.Getenv("LANG_CODE"); code != "" {
		gtui.SetLanguage(code)
	}
	if os.Getenv("MODE") == "light" {
		app.Theme(colors.Light)
	}
	if label := os.Getenv("PAGE"); label != "" {
		show(label)
	}

	app.Refresh()
	app.Select("#nav" + active().Key)

	fmt.Println(app.Snapshot(width, height))
	return true
}
