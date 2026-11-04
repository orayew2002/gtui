package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/orayew2002/gtui"
)

// check drives the gallery through its own components and reports what worked,
// so "does this framework actually behave" is one command away:
//
//	go run ./app check
func check() bool {
	if len(os.Args) < 2 || os.Args[1] != "check" {
		return false
	}

	failed := 0

	run := func(name string, steps func()) {
		goTo("buttons")
		app.Refresh()
		app.Snapshot(120, 40)

		steps()

		if lastFailure == "" {
			fmt.Printf("  ok    %s\n", name)
			return
		}

		fmt.Printf("  FAIL  %s: %s\n", name, lastFailure)
		lastFailure = ""
		failed++
	}

	fmt.Println("gtui gallery self check")

	run("button click", func() {
		press("1", "a", "\r")
		expect("primary clicked")
	})

	run("icon button", func() {
		press("1", "b", "\r")
		expect("edit")
	})

	run("dropdown picks", func() {
		press("1", "d", "\r", "down", "\r")
		expect("duplicate")
	})

	run("tabs swap", func() {
		press("1", "p")
		expect("second page")
	})

	run("typing and reading", func() {
		press("2", "a")
		app.Type("hello")
		click("#read-values")
		expect(`text="hello"`)
	})

	run("text area newline", func() {
		press("2", "b")
		app.Type("one")
		press("\r")
		app.Type("two")
		expect("one")
		expect("two")
	})

	run("switch and checkbox", func() {
		press("2", "d", "\r", "down", "\r")
		press("down", "\r")
		expect("switch=")
	})

	run("selector picks", func() {
		press("3", "a", "\r", "down", "\r")
		expect("disk-02")
	})

	run("autocomplete suggests", func() {
		press("3", "b")
		app.Type("nas")
		expect("nas-01")
	})

	run("range moves", func() {
		press("3", "d", "right", "right")
		expect("50")
	})

	run("calendar moves", func() {
		press("3", "e", "right")
		expect("←→ ↑↓ ⇞⇟")
	})

	run("table row action", func() {
		press("4", "a", "\t", "\r")
		expect("Ada Lovelace")
	})

	run("tree expands", func() {
		press("5", "a")
		expect("internal")

		press("\r")
		missing("internal")

		press("\r")
		expect("internal")
	})

	run("accordion opens", func() {
		press("5", "b", "down", "down", "\r")
		expect("gtui gallery")
	})

	run("scroll list scrolls", func() {
		press("5", "d", "pagedown")
		expect("row number 8")
	})

	run("toast appears", func() {
		press("6", "a", "\r")
		expect("saved")
	})

	run("modal opens and closes", func() {
		press("6", "b", "\r")
		expect("modals take focus")
		press("esc")
		missing("modals take focus")
	})

	run("charts draw", func() {
		press("7")
		expect("bar chart")
		expect("gauges")
	})

	run("picture viewer", func() {
		press("8", "a", "\r")
		expect("picture")
	})

	run("music player", func() {
		press("8")
		expect("Nocturne in E flat")
	})

	run("motion frames", func() {
		press("9")
		expect("recording")
	})

	run("layout wraps", func() {
		press("0")
		expect("tile 6")
	})

	run("icons listed", func() {
		press("i")
		expect("terminal")
	})

	run("theme switches", func() {
		press("t", "b", "\r")
		if app.ThemeName() == "" {
			fail("no theme name")
		}
		expect("active:")
	})

	run("palette opens", func() {
		openPalette()
		expect("type to filter")
		press("esc")
	})

	if failed > 0 {
		fmt.Printf("\n%d checks failed\n", failed)
		return true
	}

	fmt.Println("\nall checks passed")
	return true
}

var (
	lastFailure string
	screen      string
)

// press sends keys the way a person would, redrawing in between.
func press(keys ...string) {
	for _, key := range keys {
		switch key {
		case "\r":
			app.Press("enter")
		default:
			app.Press(key)
		}
		screen = app.Snapshot(120, 40)
	}
}

func expect(text string) {
	screen = app.Snapshot(120, 40)
	if !strings.Contains(screen, text) {
		fail("missing %q", text)
	}
}

func missing(text string) {
	screen = app.Snapshot(120, 40)
	if strings.Contains(screen, text) {
		fail("still shows %q", text)
	}
}

func fail(format string, args ...any) {
	if lastFailure == "" {
		lastFailure = fmt.Sprintf(format, args...)
	}
}

// click presses a button by id, for places where tabbing there would be
// ambiguous.
func click(sel string) {
	node := app.Query(sel)
	if node == nil {
		fail("no node %s", sel)
		return
	}

	node.Dispatch(&gtui.Event{Type: "click", Target: node})
	screen = app.Snapshot(120, 40)
}
