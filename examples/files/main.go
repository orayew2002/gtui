// A file browser: a tree on the left, a long list in the middle that scrolls
// without building a row it does not show, and what you picked on the right.
package main

import "github.com/orayew2002/gtui"

var (
	app    = gtui.Init()
	picked = "README.md"
)

func main() {
	app.View(view)
	app.QuitOn("q", "ctrl+c")
	app.Render()
}

func view() gtui.Component {
	return gtui.Column("full bg-canvas",
		gtui.Row("px-1 items-center gap-2 border-b border-line",
			gtui.IconText("text-accent bold nowrap", "folder", "files"),
			gtui.Breadcrumb("text-muted truncate", "project", "internal", picked),
			gtui.Spacer(),
			gtui.Chips("", "go", "1.24"),
		),
		gtui.Split("grow clip px-1 py-1", "w-20", tree(), middle()),
		gtui.StatusBar("border-t border-line px-1",
			[]gtui.Component{gtui.Text("nowrap text-muted", "1 000 files")},
			[]gtui.Component{gtui.Text("nowrap text-muted dim", "arrows move · q quits")},
		),
	)
}

func tree() gtui.Component {
	return gtui.TreeView(gtui.TreeOptions{
		Class: "clip",
		Roots: []gtui.TreeNode{{Label: "project", Open: true, Children: []gtui.TreeNode{
			{Label: "cmd", Children: []gtui.TreeNode{
				{Label: "main.go", Icon: "code", Value: "1.2k"},
			}},
			{Label: "internal", Open: true, Children: []gtui.TreeNode{
				{Label: "api.go", Icon: "code", Value: "8.4k"},
				{Label: "store.go", Icon: "code", Value: "3.1k"},
			}},
			{Label: "README.md", Icon: "file", Value: "900"},
		}}},
		OnSelect: func(item gtui.TreeNode) {
			picked = item.Label
			app.Refresh()
		},
	})
}

func middle() gtui.Component {
	return gtui.Row("gap-1 clip",
		gtui.Column("grow clip gap-1",
			gtui.Card("grow clip rounded border-line bg-panel px-1", "every file",
				gtui.Text("text-muted dim", "a thousand rows, only the visible ones are built"),
				gtui.ScrollList("grow #list", 10, 1000, row),
			),
			recent(),
		),
		details(),
	)
}

// row builds one line of the long list, and is only ever called for the lines
// on screen.
func row(at int) gtui.Component {
	return gtui.ListRow("px-1 focusable hover:bg-raised focus:bg-raised",
		"file-"+itoa(at)+".go",
		gtui.Text("text-muted nowrap", "%d B", 400+at*7),
	)
}

func recent() gtui.Component {
	return gtui.Card("shrink-0 rounded border-line bg-panel px-1", "recent",
		gtui.List("",
			gtui.ListItem("", "api.go", "8.4k"),
			gtui.ListItem("", "store.go", "3.1k"),
			gtui.ListItem("", "main.go", "1.2k"),
		),
	)
}

func details() gtui.Component {
	return gtui.Column("w-24 shrink-0 clip gap-1",
		gtui.Window("clip", picked,
			gtui.KeyValue("",
				[2]string{"size", "8.4k"},
				[2]string{"kind", "go source"},
				[2]string{"read", "yes"},
			),
			gtui.Row("gap-1 pt-1",
				gtui.Tag("", "tracked", nil),
				gtui.Tag("", "staged", func(e *gtui.Event) { app.Toast("info", "unstaged") }),
			),
		),
		gtui.Accordion("rounded border-line bg-panel px-1",
			gtui.Section{Title: "history", Open: true, Body: []gtui.Component{
				gtui.Text("text-muted dim nowrap", "3 commits"),
			}},
			gtui.Section{Title: "blame", Body: []gtui.Component{
				gtui.Text("text-muted dim nowrap", "nobody yet"),
			}},
		),
		gtui.Collapse("rounded border-line bg-panel px-1", "trash", false,
			gtui.Centered("h-3", gtui.Text("text-muted dim nowrap", "empty")),
		),
	)
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}

	digits := ""
	for v > 0 {
		digits = string(rune('0'+v%10)) + digits
		v /= 10
	}
	return digits
}
