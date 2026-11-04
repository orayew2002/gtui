package main

import "github.com/orayew2002/gtui"

// card wraps one component demo: a title, the live component, and a line that
// says what to try. Every page is a handful of these.
func card(key, title string, body ...gtui.Component) gtui.Component {
	head := []gtui.Component{gtui.Text("text-ink bold nowrap", title)}

	class := "col grow rounded border-line bg-panel px-2 gap-0"
	if key != "" {
		class += " key-" + key
	}

	return gtui.Div(class, append(head, body...)...)
}

// note is the small grey line under a demo.
func note(text string, args ...any) gtui.Component {
	return gtui.Text("text-muted dim nowrap", text, args...)
}

// out shows a live result, so clicking or typing visibly does something.
func out(id, text string) gtui.Component {
	return gtui.Text("text-accent nowrap #"+id, text)
}

// say writes into an out() box.
func say(id, text string, args ...any) {
	app.Open("#"+id, gtui.Text("text-accent nowrap", text, args...))
}
