package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/orayew2002/gtui"
	"github.com/orayew2002/gtui/src/dom"
	"github.com/orayew2002/gtui/src/render"
)

// audit drives every page with the mouse and the keyboard and checks, after
// every single event, that the screen still makes sense:
//
//	go run ./app audit
//
// Where check proves each component does its job, audit proves none of them
// breaks the page while being poked at.
func audit() bool {
	if len(os.Args) < 2 || os.Args[1] != "audit" {
		return false
	}

	fmt.Println("gtui gallery interaction audit")

	for _, item := range pages() {
		run("clicking every cell of "+item.Title, func() { auditClicks(item, wideW, wideH) })
		run("hovering all over "+item.Title, func() { auditHover(item) })
		run("the wheel on "+item.Title, func() { auditWheel(item) })
		run("the keyboard on "+item.Title, func() { auditKeys(item) })
		run("focus stays visible on "+item.Title, func() { auditFocusVisible(item) })
		run(item.Title+" at eighty columns", func() { auditClicks(item, narrowW, narrowH) })
	}

	run("every page comes back the way it left", auditRevisit)

	fmt.Printf("\n%d audits, %d failed, %d pages poked at %dx%d and %dx%d\n",
		checks, failed, len(pages()), narrowW, narrowH, wideW, wideH)

	if failed > 0 {
		os.Exit(1)
	}
	return true
}

// begin puts one page on screen at a given size, with nothing left over from
// the page before.
func begin(item page, width, height int) string {
	app.ThemeByName("dark")
	gtui.SetLanguage("en")
	goTo(item.Title)

	app.HTML(shell())
	app.Select("#nav" + currentKey)

	return app.Snapshot(width, height)
}

// sane is every rule that has to hold after any interaction at all.
func sane(where string, shot string, width int) {
	if diags := app.Diagnostics(); len(diags) > 0 {
		fail("%s: style complained about %v", where, diags[0])
		return
	}

	if focused := app.Focused(); focused != nil {
		if !holds(app.Root(), focused) {
			fail("%s: focus points at a node that is no longer on the page: id=%q class=%q text=%q",
				where, focused.Get("id"), focused.Get("class"), longest(focused.Text()))
			return
		}
		if !jumpable(focused) {
			fail("%s: focus sits on a box that cannot hold it: %q", where, focused.Get("class"))
			return
		}
	}

	for i, line := range strings.Split(strings.TrimRight(shot, "\n"), "\n") {
		if got := render.StringWidth(line); got > width {
			fail("%s: line %d is %d cells wide on a %d column screen", where, i, got, width)
			return
		}

		for _, r := range line {
			if r != '\t' && unicode.IsControl(r) {
				fail("%s: line %d holds the control character %q", where, i, r)
				return
			}
		}
	}

	bad := ""
	app.Root().Walk(func(n *dom.Node) bool {
		if n.ScrollY < 0 || n.ScrollX < 0 {
			bad = fmt.Sprintf("%s: a box scrolled past its start, y=%d x=%d", where, n.ScrollY, n.ScrollX)
			return false
		}
		return true
	})
	if bad != "" {
		fail("%s", bad)
	}
}

func holds(root, want *dom.Node) bool {
	found := false
	root.Walk(func(n *dom.Node) bool {
		if n == want {
			found = true
			return false
		}
		return true
	})
	return found
}

// jumpable mirrors src/style/utility.go: focusable, scroll, scroll-x and a
// key-x jump target each make a box able to hold the keyboard.
func jumpable(n *dom.Node) bool {
	for _, class := range strings.Fields(n.Get("class")) {
		switch {
		case class == "focusable", class == "scroll", class == "scroll-x":
			return true
		case strings.HasPrefix(class, "key-") && class != "key-hide" && class != "key-show":
			return true
		}
	}
	return false
}

// pane is the first column of the content, past the sidebar, which collapses
// to keys only on a narrow screen.
func pane(width int) int {
	if width < 90 {
		return 8
	}
	return 21
}

// auditClicks presses the mouse on the content of a page, cell by cell. It
// stays out of the sidebar, or the first click would navigate away and leave
// the rest of the page unpoked.
func auditClicks(item page, width, height int) {
	begin(item, width, height)

	for y := range height {
		for x := pane(width); x < width; x += 2 {
			app.Click(x, y)

			if currentKey != item.Key {
				fail("a click at %d,%d navigated to %q", x, y, currentKey)
				return
			}

			sane(fmt.Sprintf("after clicking %d,%d", x, y), app.Snapshot(width, height), width)
			if lastFailure != "" {
				return
			}

			// whatever a click opened, esc has to take away again
			if app.ModalOpen() {
				app.Press("esc")
				if app.ModalOpen() {
					app.CloseModal()
					fail("esc did not close what the click at %d,%d opened", x, y)
					return
				}
			}
		}
	}
}

func auditHover(item page) {
	begin(item, wideW, wideH)

	for y := range wideH {
		for x := 0; x < wideW; x += 2 {
			app.Hover(x, y)
		}
	}

	sane("after hovering everywhere", app.Snapshot(wideW, wideH), wideW)
}

func auditWheel(item page) {
	begin(item, wideW, wideH)

	spots := [][2]int{{2, 2}, {wideW / 2, wideH / 2}, {wideW - 3, wideH - 3}, {pane(wideW), wideH - 6}}

	for _, spot := range spots {
		for _, lines := range []int{-6, 6, 40, -40} {
			app.Wheel(spot[0], spot[1], lines)

			sane(fmt.Sprintf("after the wheel at %d,%d", spot[0], spot[1]),
				app.Snapshot(wideW, wideH), wideW)
			if lastFailure != "" {
				return
			}
		}
	}
}

func auditKeys(item page) {
	begin(item, wideW, wideH)

	keys := []string{"enter", "down", "up", "right", "left", "pagedown", "pageup", "home", "end", "esc"}
	stops := 0

	app.Root().Walk(func(n *dom.Node) bool {
		if n.IsElement() && jumpable(n) {
			stops++
		}
		return true
	})

	seen := map[*dom.Node]bool{}

	for step := 0; step < stops+4; step++ {
		app.Press("tab")

		if focused := app.Focused(); focused != nil {
			seen[focused] = true
		}

		for _, key := range keys {
			app.Press(key)

			sane("after "+key, app.Snapshot(wideW, wideH), wideW)
			if lastFailure != "" {
				return
			}
		}

		if app.ModalOpen() {
			app.Press("esc")
			if app.ModalOpen() {
				app.CloseModal()
				fail("a dialog survived esc at tab stop %d", step)
				return
			}
		}
	}

	if stops > 1 && len(seen) < 2 {
		fail("tab reached %d of %d stops, the keyboard is stuck", len(seen), stops)
	}
}

// auditFocusVisible tabs through a page and asks that whatever holds the
// keyboard is painted, so focus never hides below the fold.
func auditFocusVisible(item page) {
	begin(item, wideW, wideH)

	for step := 0; step < 24; step++ {
		app.Press("tab")
		shot := app.Snapshot(wideW, wideH)

		focused := app.Focused()
		if focused == nil || len(focused.ElementChildren()) > 0 {
			continue
		}

		// match on the longest plain word: a node's text runs its glyphs
		// together while the screen lays them out with gaps
		word := longest(focused.Text())
		if len(word) < 4 {
			continue
		}

		if !strings.Contains(shot, word) {
			fail("focus sits on %q at tab %d but it is not on screen", word, step)
			return
		}
	}
}

func longest(text string) string {
	best, run := "", strings.Builder{}

	flush := func() {
		if run.Len() > len(best) {
			best = run.String()
		}
		run.Reset()
	}

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' {
			run.WriteRune(r)
			continue
		}
		flush()
	}
	flush()

	return best
}

// auditRevisit leaves every page and comes back to it, because state left
// behind shows up as a different screen.
func auditRevisit() {
	// these pages move on their own, so two snapshots of them differ by design
	moving := map[string]bool{"motion": true, "media": true, "feedback": true}

	for _, item := range pages() {
		if moving[item.Title] {
			continue
		}

		first := begin(item, wideW, wideH)

		goTo("charts")
		app.HTML(shell())
		app.Snapshot(wideW, wideH)

		if again := begin(item, wideW, wideH); again != first {
			fail("%s does not come back the way it left", item.Title)
			return
		}
	}
}
