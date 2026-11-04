package main

import "github.com/orayew2002/gtui"

var tabs = []string{"gtui — github", "go.dev", "news"}

func main() {
	app := gtui.Init()

	app.HTML(
		gtui.Column("full bg-canvas gap-0",
			gtui.Toolbar("px-1 gap-2", tabBar()),
			gtui.AddressBar("px-1 py-0", "url", "https://github.com/orayew2002/gtui", func(url string) {
				app.Toast("", "loading "+url)
				app.Open("#page", page(url))
			}),
			gtui.Div("grow scroll px-2 key-c #page", page("https://github.com/orayew2002/gtui")),
			gtui.StatusBar("border-t border-line",
				[]gtui.Component{gtui.Text("", "done")},
				[]gtui.Component{gtui.Hint("", "tab", "focus"), gtui.Hint("", "q", "quit")},
			),
		),
	)

	app.Render()
}

func tabBar() gtui.Component {
	items := make([]gtui.Component, 0, len(tabs))

	for i, title := range tabs {
		class := "px-2 nowrap text-muted focusable focus:bold"
		if i == 0 {
			class = "px-2 nowrap text-ink bold bg-raised focusable"
		}
		items = append(items, gtui.Text(class, title))
	}

	return gtui.Row("gap-1", items...)
}

func page(url string) gtui.Component {
	return gtui.Column("gap-0",
		gtui.Text("text-ink bold", "gtui"),
		gtui.Text("text-muted", url),
		gtui.Divider("border-line"),
		gtui.Text("", "Terminal user interfaces for Go, written like HTML. A div is a box, classes are utilities, and links swap a target div."),
		gtui.Text("text-accent", "★ 1.2k   forks 84   MIT"),
	)
}
