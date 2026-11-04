package main

import (
	"fmt"

	"github.com/orayew2002/gtui"
)

// Everything every page is built from lives here: one way to frame a demo, one
// way to explain it, one way to show a live result, and the two numbers a page
// needs to lay itself out. A new page needs nothing but these.

// narrow is the width under which the sidebar shrinks to its shortcut keys, so
// a page still has room on an eighty column terminal.
const narrow = 90

// card frames one demo: a title, the live component, and whatever else the page
// wants to put inside.
//
// key is the letter of the ┤x├ jump badge. c and n belong to the content pane
// and the sidebar, so cards use a, b, d, e, f, g, h, j.
func card(key, title string, body ...gtui.Component) gtui.Component {
	head := []gtui.Component{gtui.Text("text-ink bold truncate", title)}

	// clip is what keeps a wide child inside the card instead of pushing the
	// row past the right edge of the terminal
	class := "col grow clip rounded border-line bg-panel px-2 gap-0"
	if key != "" {
		class += " key-" + key
	}

	return gtui.Div(class, append(head, body...)...)
}

// demo is card with a line of explanation, which a terminal short of rows gives
// up before it gives up the component.
func demo(key, title, about string, body ...gtui.Component) gtui.Component {
	if about != "" && wide() {
		body = append([]gtui.Component{note(about)}, body...)
	}
	return card(key, title, body...)
}

// cards is a row of demos that stacks instead of squeezing on a narrow screen.
func cards(items ...gtui.Component) gtui.Component {
	if compact() {
		return gtui.Column("gap-1", items...)
	}
	return gtui.Row("gap-1", items...)
}

// note is the small grey line under a demo. It wraps, so a long instruction is
// never cut mid-word.
func note(text string, args ...any) gtui.Component {
	return gtui.Text("text-muted dim", text, args...)
}

// out is a live result, so clicking or typing visibly does something.
func out(id, text string) gtui.Component {
	return gtui.Text("text-accent #"+id, text)
}

// say writes into an out box.
func say(id, text string, args ...any) {
	if len(args) > 0 {
		text = fmt.Sprintf(text, args...)
	}
	app.Text("#"+id, text)
}

// compact reports whether the shell is in its narrow layout.
func compact() bool {
	width, _ := app.Size()
	return width < narrow
}

// room is the cells a page has inside #main: the shell spends the rest on the
// sidebar, the header, the gaps and the borders, and the sidebar itself is
// narrower on a small terminal.
func room() (int, int) {
	width, height := app.Size()

	cols := width - 23
	if compact() {
		cols = width - 10
	}
	cols = max(cols, 24)

	return cols, max(height-9, 8)
}

// wide reports whether a page has room for the extra explanation a demo can
// carry.
func wide() bool {
	cols, _ := room()
	return cols >= 80
}
