// The self check. It drives the gallery through its own components and reports
// pass or fail for each one, so "does this framework behave" is one command
// away: go run ./app check
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/orayew2002/gtui"
	"github.com/orayew2002/gtui/src/render"
)

// Interactions run at the wide size, where every demo is on screen. Both sizes
// are checked for overflow.
const (
	wideW, wideH     = 120, 40
	narrowW, narrowH = 80, 24
)

// keyNames maps the control characters a test author types to the names the
// framework gives those keys.
var keyNames = map[string]string{
	"\r":   "enter",
	"\n":   "enter",
	"\t":   "tab",
	"\x1b": "esc",
	"\x7f": "backspace",
}

// check drives the gallery through its own components and reports what worked,
// so "does this framework actually behave" is one command away:
//
//	go run ./app check
//
// Every step has to prove the component reacted: the screen changes, and the
// result is the one only that component produces.
func check() bool {
	if len(os.Args) < 2 || os.Args[1] != "check" {
		return false
	}

	fmt.Println("gtui gallery self check")

	run("jump badges", func() {
		start("inputs")

		// a card wears its key in its border; the sidebar and the page pane
		// carry theirs with key-hide, because a single rule has nowhere to
		// print a badge
		expect("┤a├")
		expect("┤b├")

		press("a")
		changed()

		if app.Focused() == nil {
			fail("the a badge moved focus nowhere")
		}

		for _, key := range []string{"n", "c"} {
			press(key)
			if app.Focused() == nil {
				fail("the hidden %s key moved focus nowhere", key)
			}
		}
	})

	run("sidebar marks one page", func() {
		start("buttons")
		press("5")
		changed()

		marked := make([]string, 0, 1)
		for _, item := range pages() {
			if node := app.Query("#nav" + item.Key); node != nil && node.State.Current {
				marked = append(marked, item.Title)
			}
		}

		if len(marked) != 1 || marked[0] != "lists" {
			fail("the sidebar marks %v", marked)
		}
	})

	run("button click", func() {
		start("buttons")
		press("a", "\r")
		changed()
		expectOut("btn-out", "primary clicked")
	})

	run("icon button", func() {
		start("buttons")
		press("b", "\r")
		changed()
		expectOut("icon-out", "edit")
	})

	run("dropdown picks", func() {
		start("buttons")
		focus("#dd")
		press("\r")
		changed()
		expect("duplicate")

		press("down", "\r")
		changed()
		expectOut("dd-out", "duplicate")
		missing("rename")
	})

	run("tabs swap", func() {
		start("buttons")
		expect("first page")

		press("p")
		changed()
		expect("second page")
		missing("first page")
	})

	run("prompt asks and answers", func() {
		start("buttons")
		hit("prompt")
		changed()
		expect("new name")

		typed("disk-09")
		press("\r")
		changed()
		expectOut("btn-out", `renamed to "disk-09"`)
	})

	run("typing and reading", func() {
		start("inputs")
		focus("#in-text")
		typed("hello")
		changed()
		expectValue("#in-text", "hello")
		expectOut("in-out", "hello")

		hit("read values")
		expectOut("area-out", `text="hello"`)
	})

	run("text area newline", func() {
		start("inputs")
		focus("#in-area")
		typed("one")
		press("\r")
		typed("two")
		changed()
		expectValue("#in-area", "one\ntwo")
	})

	run("switch and checkbox", func() {
		start("inputs")
		focus("#in-switch")
		press("\r")
		changed()
		expectValue("#in-switch", "false")
		expectOut("tog-out", "switch=false")

		focus("#in-check")
		press("\r")
		changed()
		expectValue("#in-check", "true")
		expectOut("tog-out", "checkbox=true")
	})

	run("search types", func() {
		start("inputs")
		focus("#in-find")
		typed("disk")
		changed()
		expectValue("#in-find", "disk")
		expectOut("find-out", "disk")
	})

	run("set values from code", func() {
		start("inputs")
		hit("fill")
		changed()
		expectValue("#in-text", "written from code")
		expectValue("#in-area", "line one\nline two")

		hit("clear")
		changed()
		expectValue("#in-text", "")
	})

	run("selector picks", func() {
		start("pickers")
		focus("#pick-disk")
		press("\r")
		changed()
		expect("disk-03")

		press("down", "\r")
		changed()
		expectValue("#pick-disk", "disk-02")
		expectOut("pick-out", "disk-02")
	})

	run("autocomplete suggests", func() {
		start("pickers")
		focus("#pick-find")
		typed("nas")
		changed()
		expect("nas-01")

		press("down", "\r")
		expectValue("#pick-find", "nas-01")
	})

	run("range moves", func() {
		start("pickers")
		focus("#pick-range")
		press("right", "right")
		changed()
		expectValue("#pick-range", "50")
		expectOut("range-out", "50")
	})

	run("radios pick", func() {
		start("pickers")
		focus("#pick-range")
		press("\t", "\r")
		changed()

		if quality := app.Values()["quality"]; quality != "draft" {
			fail("quality is %q", quality)
		}
	})

	run("select picks", func() {
		start("pickers")
		focus("#pick-sort")
		press("\r", "down", "down", "\r")
		changed()
		expectValue("#pick-sort", "oldest")
	})

	run("calendar moves", func() {
		start("pickers")
		focus("#pick-day")
		press("right", "\r")
		changed()

		if day := app.Value("#pick-day"); day == "" {
			fail("no day selected")
		}
	})

	run("table selects a row", func() {
		start("tables")
		press("a", "\t", "\r")
		changed()
		expectOut("table-out", "selected Ada Lovelace")
	})

	run("table row action", func() {
		start("tables")
		press("a", "\t", "\t", "\r")
		changed()
		expectOut("table-out", "edit Ada Lovelace")
	})

	run("table asks before deleting", func() {
		start("tables")
		press("a", "\t", "\t", "\t", "\r")
		changed()
		expect("delete Ada Lovelace?")

		press("\r")
		changed()
		expectOut("table-out", "deleted Ada Lovelace")
	})

	run("tree folds", func() {
		start("lists")
		press("a")
		expect("internal")

		press("\r")
		changed()
		missing("internal")

		press("\r")
		changed()
		expect("internal")
	})

	run("accordion opens", func() {
		start("lists")
		press("b", "down", "\r")
		changed()
		expect("hidden until you open it")
	})

	run("scroll list scrolls", func() {
		start("lists")
		focus("#big-list")
		expect("row number 0")

		press("pagedown")
		changed()
		missing("row number 0")
		expectOut("list-out", "first visible row")
	})

	run("tag removes itself", func() {
		start("lists")
		hit("×")
		changed()
		expectOut("tag-out", "removed")
	})

	run("toast appears", func() {
		start("feedback")
		press("a", "\r")
		changed()
		expect("saved")
	})

	run("dialog opens and closes", func() {
		start("feedback")
		hit("alert")
		changed()
		expect("modals take focus")

		press("\x1b")
		changed()
		missing("modals take focus")
	})

	run("confirm runs its action", func() {
		start("feedback")
		hit("confirm")
		changed()
		expect("delete the pool?")

		press("\r")
		changed()
		expectOut("busy-out", "pool deleted")
	})

	run("charts draw", func() {
		start("charts")
		expect("█")
		expect("▁")
	})

	run("picture viewer", func() {
		start("media")
		focus("#mv-viewer")
		press("right")
		changed()
		expectOut("mv-out", "picture 2 of 3")
	})

	run("music player", func() {
		start("media")
		was := nodeText("#mp-meter")

		click("#mp-next")
		changed()

		click("#mp-play")
		changed()
		expect("▶")

		if now := nodeText("#mp-meter"); now == was {
			fail("the transport stayed %q", was)
		}
	})

	run("video player seeks", func() {
		start("media")
		was := nodeText("#mv-frame")

		// the player moves the model and the next tick draws it, which is how
		// one timer drives every moving box on the page
		focus("#mv-player")
		press("right")
		animate(1)
		changed()

		if now := nodeText("#mv-frame"); now == was {
			fail("the frame stayed %q", was)
		}
	})

	run("motion moves on its own", func() {
		start("motion")
		animate(3)
		changed()
		animate(4)
		changed()
	})

	run("the tick leaves still pages alone", func() {
		start("layout")
		animate(5)
		still()
	})

	run("layout scrolls both ways", func() {
		start("layout")
		expect("tile 6")

		focus("#both")
		press("down", "right")
		changed()
	})

	run("icons page and filter", func() {
		start("icons")

		names := gtui.IconNames()
		if len(names) < 100 {
			fail("only %d icons", len(names))
		}
		expect(names[0])

		const filter = "battery"
		focus("#ic-find")
		typed(filter)
		changed()
		expect("battery")
		missing(names[0])
		expectOut("ic-info", fmt.Sprintf("of %d", len(names)))

		for range len(filter) {
			press("\x7f")
		}
		changed()
		expect(names[0])
	})

	run("theme switches", func() {
		start("themes")
		expect("active: dark")

		click("#th-light")
		changed()
		expect("active: light")

		if name := app.ThemeName(); name != "light" {
			fail("the theme is %q", name)
		}

		click("#th-new")
		changed()
		expect("active: ocean")
	})

	run("theme button toggles", func() {
		start("buttons")

		// a theme is colour only, so the snapshot cannot show it: the proof
		// that the header button reacted is the active theme itself
		was := app.ThemeName()
		click("#mode")

		if now := app.ThemeName(); now == was {
			fail("the theme stayed %q", now)
		}

		click("#mode")
		if back := app.ThemeName(); back != was {
			fail("the theme is %q, want %q", back, was)
		}
	})

	run("language switches", func() {
		start("i18n")
		expect("gtui gallery")

		click("#lang-tk")
		changed()
		expect("gtui galereýasy")
		expect("hereket")

		click("#lang-en")
		changed()
		expect("gtui gallery")
	})

	run("style tabs and states", func() {
		start("style")
		press("z")
		changed()
		expect("tab here, then enter")

		focus("#st-try")
		press("\r")
		changed()
		expectOut("st-out", "clicked the focus demo")
	})

	run("extras tabs swap", func() {
		start("extras")
		expect("Centered inside it")

		press("w")
		changed()
		expect("18 402")
		missing("Centered inside it")
	})

	run("palette opens and goes", func() {
		start("buttons")
		press("ctrl+p")
		changed()
		expect("type to filter")

		typed("charts")
		press("\r")
		changed()
		expectOut("head", "charts")
	})

	for _, item := range pages() {
		run("page "+item.Title, func() {
			from := "buttons"
			if item.Title == from {
				from = "layout"
			}

			start(from)
			press(item.Key)
			changed()

			expectOut("head", phrase("nav."+item.Title, item.Title))
			if body := nodeText("#main"); body == "" {
				fail("#main is empty")
			}
		})
	}

	for _, item := range pages() {
		run("fits "+item.Title, func() {
			fits(item, narrowW, narrowH)
			fits(item, wideW, wideH)
		})
	}

	fmt.Printf("\n%d checks, %d failed, %d pages at %dx%d and %dx%d\n",
		checks, failed, len(pages()), narrowW, narrowH, wideW, wideH)

	return true
}

var (
	checks      int
	failed      int
	lastFailure string

	screen   string // the gallery as of the last draw
	baseline string // the gallery as of the last checkpoint
)

func run(name string, steps func()) {
	checks++
	lastFailure = ""

	steps()

	if lastFailure == "" {
		fmt.Printf("  ok    %s\n", name)
		return
	}

	fmt.Printf("  FAIL  %s: %s\n", name, lastFailure)
	failed++
}

func fail(format string, args ...any) {
	if lastFailure == "" {
		lastFailure = fmt.Sprintf(format, args...)
	}
}

// start rebuilds the whole shell on one page, so no check inherits typed text,
// an open menu or a stale theme from the check before it.
func start(title string) {
	app.ThemeByName("dark")
	gtui.SetLanguage("en")
	goTo(title)

	app.HTML(shell())
	app.Select("#nav" + currentKey)

	snap()
	baseline = screen
}

// snap renders the gallery the way the terminal would.
func snap() string {
	screen = app.Snapshot(wideW, wideH)
	return screen
}

// press sends keys the way a person would, redrawing in between.
func press(names ...string) {
	for _, name := range names {
		if real, ok := keyNames[name]; ok {
			name = real
		}

		app.Press(name)
		snap()
	}
}

// typed types into whatever holds focus.
func typed(text string) {
	app.Type(text)
	snap()
}

// focus puts the keyboard on one control, so a check says which component it is
// about instead of counting tab stops.
func focus(sel string) {
	node := app.Query(sel)
	if node == nil {
		fail("no node %s", sel)
		return
	}

	app.Focus(node)
	snap()
}

// click activates a node by id.
func click(sel string) {
	node := app.Query(sel)
	if node == nil {
		fail("no node %s", sel)
		return
	}

	node.Dispatch(&gtui.Event{Type: "click", Target: node})
	snap()
}

// hit clicks the control labelled text, which is how a check reaches a button
// that carries no id of its own.
func hit(text string) {
	for _, node := range app.QueryAll(".focusable") {
		if strings.TrimSpace(node.Text()) == text {
			node.Dispatch(&gtui.Event{Type: "click", Target: node})
			snap()
			return
		}
	}

	fail("no control labelled %q", text)
}

// animate runs one frame of the gallery's timer, the way app.Animate does.
func animate(at int) {
	onTick(at)
	snap()
}

// changed proves the component reacted: the screen is not the one from the last
// checkpoint. It then moves the checkpoint on, so a check can prove every step
// instead of only the first.
func changed() {
	if screen == baseline {
		fail("the screen did not change")
	}
	baseline = screen
}

// still is the opposite, for the pages the timer must not touch.
func still() {
	if screen != baseline {
		fail("the screen changed")
	}
	baseline = screen
}

func expect(text string) {
	if !strings.Contains(screen, text) {
		fail("missing %q", text)
	}
}

func missing(text string) {
	if strings.Contains(screen, text) {
		fail("still shows %q", text)
	}
}

// expectOut reads a box out of the document, so a narrow card cannot hide the
// result a check is about.
func expectOut(id, want string) {
	if got := nodeText("#" + id); !strings.Contains(got, want) {
		fail("#%s says %q, want %q", id, got, want)
	}
}

func expectValue(sel, want string) {
	if got := app.Value(sel); got != want {
		fail("%s is %q, want %q", sel, got, want)
	}
}

func nodeText(sel string) string {
	node := app.Query(sel)
	if node == nil {
		fail("no node %s", sel)
		return ""
	}
	return strings.TrimSpace(node.Text())
}

// fits renders one page at a terminal size and fails if anything escaped the
// shell: a row wider than the terminal, or a right border that is not where the
// frame belongs. The shell pads by one cell, so the last column stays empty and
// the frame sits one column in.
func fits(item page, width, height int) {
	goTo(item.Title)
	app.HTML(shell())

	lines := strings.Split(app.Snapshot(width, height), "\n")
	if len(lines) != height {
		fail("%dx%d: %d rows", width, height, len(lines))
		return
	}

	for y, line := range lines {
		if w := render.StringWidth(line); w > width {
			fail("%dx%d: row %d is %d cells wide: %q", width, height, y, w, line)
			return
		}
	}

	// the rule beside the sidebar holds its column on every row of the body,
	// which it only manages while nothing is spilling into it
	rule := sidebarRule(lines)
	if rule < 0 {
		fail("%dx%d: no rule between the sidebar and the page", width, height)
		return
	}

	for y := 2; y < height; y++ {
		if r := cellAt(lines[y], rule); r != '│' && r != ' ' {
			fail("%dx%d: row %d puts %q where the sidebar rule belongs", width, height, y, r)
			return
		}
	}
}

// sidebarRule is the column the sidebar's right rule sits in. The sidebar is
// the leftmost box, so the first rule on a body row is its own.
func sidebarRule(lines []string) int {
	for y := 2; y < len(lines); y++ {
		for x := range 20 {
			if cellAt(lines[y], x) == '│' {
				return x
			}
		}
	}
	return -1
}

const borderRunes = "╭╮╰╯│─┤├"

// cellAt returns the rune in one display column of a snapshot row.
func cellAt(line string, col int) rune {
	x := 0
	for _, r := range line {
		if x+render.RuneWidth(r) > col {
			return r
		}
		x += render.RuneWidth(r)
	}
	return ' '
}
