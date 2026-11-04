// Package demos holds the example apps. Each one is a single function that
// builds a screen, so the same app can be run on its own from examples/ and
// opened inside the gallery's projects page.
//
// A project never binds a global key and never starts its own timer: the host
// passes a tick and the project draws itself from it. That is what lets the
// gallery open one, close it again, and leave nothing behind.
package demos

import "github.com/orayew2002/gtui"

// Project is one example app.
type Project struct {
	Key   string // the gallery's one-key shortcut
	Name  string
	About string
	Build func(a *gtui.App, tick int) gtui.Component
}

// All is every example, in the order a newcomer should meet them.
func All() []Project {
	return []Project{
		{"1", "hello", "the smallest app there is", Hello},
		{"2", "dashboard", "stat cards, charts and a compact table", Dashboard},
		{"3", "admin", "records, row actions and dialogs that ask first", Admin},
		{"4", "files", "a tree, a thousand rows, split panes", Files},
		{"5", "monitor", "charts and gauges that move on their own", Monitor},
		{"6", "settings", "every field and picker, read back in one call", Settings},
		{"7", "studio", "a picture behind the content and both players", Studio},
		{"8", "shop", "three languages and pictures off disk", Shop},
		{"9", "form", "one panel of inputs with a save", Form},
		{"0", "browser", "tabs and links that swap one box", Browser},
		{"e", "editor", "type code, with gopls for problems", Editor},
		{"r", "player", "a music player with a queue", Player},
		{"t", "chat", "an assistant thread with a prompt box", Chat},
	}
}

// host is the app the project is drawing into. Build sets it, so the handlers a
// project hangs off its buttons can reach the app without threading it through
// every function.
var host *gtui.App
