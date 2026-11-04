package demos

import "github.com/orayew2002/gtui"

// Browser shows the routing model: the tab strip is Tabs writing into the page
// box, and an Href inside the page swaps the same box.
func Browser(a *gtui.App, tick int) gtui.Component {
	host = a

	return gtui.Column("full gap-0",
		gtui.Row("px-1 gap-2 shrink-0",
			gtui.Tabs("", "#page",
				gtui.Tab{Key: "1", Label: "github", Active: true, Page: []gtui.Component{repo()}},
				gtui.Tab{Key: "2", Label: "go.dev", Page: []gtui.Component{docs()}},
				gtui.Tab{Key: "3", Label: "news", Page: []gtui.Component{news()}},
			),
		),
		gtui.Row("px-1 items-center gap-1 shrink-0",
			gtui.Icon("text-success shrink-0", "lock"),
			gtui.AddressBar("grow py-0", "url", "https://github.com/orayew2002/gtui", func(url string) {
				host.Toast("", "loading "+url)
				host.Open("#page", repo())
			}),
		),
		gtui.Div("grow clip scroll px-2 #page", repo()),
		gtui.StatusBar("border-t border-line px-1 shrink-0",
			[]gtui.Component{gtui.Text("nowrap", "done")},
			[]gtui.Component{gtui.Hint("", "tab", "focus")},
		),
	)
}

func repo() gtui.Component {
	return gtui.Column("gap-0",
		gtui.Text("text-ink bold", "gtui"),
		gtui.Text("text-muted", "github.com/orayew2002/gtui"),
		gtui.Divider("border-line"),
		gtui.Text("", "Terminal user interfaces for Go, written like HTML. A div is a box, classes are utilities, and links swap a target div."),
		gtui.Text("text-accent", "★ 1.2k   forks 84   MIT"),
		gtui.Divider("border-line"),
		gtui.Text("text-muted dim", "a link inside the page swaps it too:"),
		gtui.Row("gap-2",
			gtui.Href("text-accent underline nowrap", "read the docs", "#page", docs()),
			gtui.Href("text-accent underline nowrap", "what is new", "#page", news()),
		),
	)
}

func docs() gtui.Component {
	cols := []gtui.Col{
		{Label: "package", Width: "w-14"},
		{Label: "what it holds", Width: "grow"},
	}

	return gtui.Column("gap-0",
		gtui.Text("text-ink bold", "go.dev/github.com/orayew2002/gtui"),
		gtui.Divider("border-line"),
		gtui.Table("", cols,
			gtui.Cells("gtui", "the components and the app"),
			gtui.Cells("src/layout", "block and flex layout"),
			gtui.Cells("src/paint", "borders, text and canvases"),
			gtui.Cells("src/render", "the diffed cell buffer"),
		),
		back(),
	)
}

func news() gtui.Component {
	return gtui.Column("gap-0",
		gtui.Text("text-ink bold", "what is new"),
		gtui.Divider("border-line"),
		gtui.List("",
			gtui.ListItem("", "background pictures fill their box", "today"),
			gtui.ListItem("", "bg-black/40 mixes with what is under it", "today"),
			gtui.ListItem("", "the pointer stops at a dialog", "yesterday"),
		),
		back(),
	)
}

// back swaps the page from a handler. Href builds the page it carries straight
// away, so a page cannot hold an Href to one that holds an Href back to it.
func back() gtui.Component {
	return gtui.Button("pt-1 text-accent underline nowrap", "back to the repo", func(e *gtui.Event) {
		host.Open("#page", repo())
	})
}
