package main

import (
	"strings"

	"github.com/orayew2002/gtui"
)

var (
	app  *gtui.App
	line = 3
)

var source = strings.Split(`package main

import "github.com/orayew2002/gtui"

func main() {
	gtui.Init().HTML(
		gtui.Center("full bg-canvas",
			gtui.Text("text-accent bold", "hello"),
		),
	).Render()
}`, "\n")

func main() {
	app = gtui.Init()
	app.View(view)

	app.Bind("down", func() { line = min(len(source), line+1) })
	app.Bind("up", func() { line = max(1, line-1) })
	app.Bind("ctrl+s", func() { app.Toast("success", "saved main.go") })

	app.Render()
}

func view() gtui.Component {
	return gtui.Column("full bg-canvas gap-0",
		gtui.Toolbar("px-1 gap-2",
			gtui.Text("text-ink bold", "main.go"),
			gtui.Text("text-muted dim", "go · utf-8"),
			gtui.Spacer(),
			gtui.Hint("", "ctrl+s", "save"),
		),
		gtui.Row("grow gap-1",
			gtui.Tree("w-20 shrink-0 px-1 border-r border-line",
				gtui.Node2{Label: "gtui", Open: true, Children: []gtui.Node2{
					{Label: "main.go"}, {Label: "div.go"}, {Label: "table.go"},
				}},
			),
			gtui.Div("grow scroll px-1 key-c",
				gtui.CodeView(gtui.CodeOptions{
					Lines:  source,
					First:  1,
					Active: line,
					Marks:  map[int]string{6: "●"},
				}),
			),
		),
		gtui.StatusBar("border-t border-line",
			[]gtui.Component{gtui.Text("", "line %d of %d", line, len(source))},
			[]gtui.Component{gtui.Text("", "spaces: 4"), gtui.Hint("", "q", "quit")},
		),
	)
}
