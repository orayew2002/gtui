package demos

import (
	"strings"

	"github.com/orayew2002/gtui"
)

var (
	editorLine   = 3
	editorSource = strings.Split(`package main

import "github.com/orayew2002/gtui"

func main() {
	gtui.Init().HTML(
		gtui.Center("full",
			gtui.Text("text-accent bold", "hello"),
		),
	).Render()
}`, "\n")
)

// Editor is a file tree beside code with a gutter, an active line and a mark.
func Editor(a *gtui.App, tick int) gtui.Component {
	host = a

	return gtui.Column("full gap-0",
		gtui.Row("px-1 gap-2 items-center shrink-0 border-b border-line",
			gtui.Text("text-ink bold nowrap", "main.go"),
			gtui.Text("text-muted dim nowrap", "go · utf-8"),
			gtui.Spacer(),
			gtui.Button("text-muted focus:text-ink nowrap", "▲", func(e *gtui.Event) {
				editorLine = max(1, editorLine-1)
				host.Refresh()
			}),
			gtui.Button("text-muted focus:text-ink nowrap", "▼", func(e *gtui.Event) {
				editorLine = min(len(editorSource), editorLine+1)
				host.Refresh()
			}),
			gtui.Button("text-accent focus:bold nowrap", "save", func(e *gtui.Event) {
				host.Toast("success", "saved main.go")
			}),
		),
		gtui.Row("grow clip",
			gtui.Tree("w-20 shrink-0 px-1 border-r border-line",
				gtui.TreeNode{Label: "gtui", Open: true, Children: []gtui.TreeNode{
					{Label: "main.go", Icon: "code"},
					{Label: "div.go", Icon: "code"},
					{Label: "table.go", Icon: "code"},
				}},
			),
			gtui.Div("grow clip scroll px-1",
				gtui.CodeView(gtui.CodeOptions{
					Lines:  editorSource,
					First:  1,
					Active: editorLine,
					Marks:  map[int]string{6: "●"},
				}),
			),
		),
		gtui.StatusBar("border-t border-line px-1 shrink-0",
			[]gtui.Component{gtui.Text("nowrap", "line %d of %d", editorLine, len(editorSource))},
			[]gtui.Component{gtui.Text("nowrap text-muted dim", "spaces: 4")},
		),
	)
}
