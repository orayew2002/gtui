package main

import (
	"fmt"

	"github.com/orayew2002/gtui"
)

// card wraps one component demo: a title, the live component, and a line that
// says what to try. Every page is a handful of these.
//
// key is the letter of the ┤x├ jump badge. c and n are taken by the content
// pane and the sidebar, so cards use a, b, d, e, f, g.
func card(key, title string, body ...gtui.Component) gtui.Component {
	head := []gtui.Component{gtui.Text("text-ink bold truncate", title)}

	// clip is what keeps a wide child inside the card instead of pushing the
	// row past the right edge of the terminal.
	class := "col grow clip rounded border-line bg-panel px-2 gap-0"
	if key != "" {
		class += " key-" + key
	}

	return gtui.Div(class, append(head, body...)...)
}

// cards is a row of card()s that stacks instead of squeezing on a narrow
// terminal. Pages use it for every row of demos.
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

// out shows a live result, so clicking or typing visibly does something.
func out(id, text string) gtui.Component {
	return gtui.Text("text-accent #"+id, text)
}

// say writes into an out() box.
func say(id, text string, args ...any) {
	if len(args) > 0 {
		text = fmt.Sprintf(text, args...)
	}
	app.Text("#"+id, text)
}
