// A file browser: a filesTree on the left, a long list in the filesMiddle that scrolls
// without building a filesRow it does not show, and what you filesPicked on the right.
package demos

import (
	"strconv"

	"github.com/orayew2002/gtui"
)

var (
	filesPicked = "README.md"
)

func Files(a *gtui.App, _ int) gtui.Component {
	host = a

	return gtui.Column("full",
		gtui.Row("px-1 items-center gap-2 border-b border-line",
			gtui.IconText("text-accent bold nowrap", "folder", "files"),
			gtui.Breadcrumb("text-muted truncate", "project", "internal", filesPicked),
			gtui.Spacer(),
			gtui.Chips("", "go", "1.24"),
		),
		gtui.Split("grow clip px-1 py-1", "w-20", filesTree(), filesMiddle()),
		gtui.StatusBar("border-t border-line px-1",
			[]gtui.Component{gtui.Text("nowrap text-muted", "1 000 files")},
			[]gtui.Component{gtui.Text("nowrap text-muted dim", "arrows move · q quits")},
		),
	)
}

func filesTree() gtui.Component {
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
			filesPicked = item.Label
			host.Refresh()
		},
	})
}

func filesMiddle() gtui.Component {
	return gtui.Row("gap-1 clip",
		gtui.Column("grow clip gap-1",
			gtui.Card("grow clip rounded border-line bg-panel px-1", "every file",
				gtui.Text("text-muted dim", "a thousand rows, only the visible ones are built"),
				gtui.ScrollList("grow #list", 10, 1000, filesRow),
			),
			filesRecent(),
		),
		filesDetails(),
	)
}

// filesRow builds one line of the long list, and is only ever called for the lines
// on screen.
func filesRow(at int) gtui.Component {
	return gtui.ListRow("px-1 focusable hover:bg-raised focus:bg-raised",
		"file-"+strconv.Itoa(at)+".go",
		gtui.Text("text-muted nowrap", "%d B", 400+at*7),
	)
}

func filesRecent() gtui.Component {
	return gtui.Card("shrink-0 rounded border-line bg-panel px-1", "recent",
		gtui.List("",
			gtui.ListItem("", "api.go", "8.4k"),
			gtui.ListItem("", "store.go", "3.1k"),
			gtui.ListItem("", "main.go", "1.2k"),
		),
	)
}

func filesDetails() gtui.Component {
	return gtui.Column("w-24 shrink-0 clip gap-1",
		gtui.Window("clip", filesPicked,
			gtui.KeyValue("",
				[2]string{"size", "8.4k"},
				[2]string{"kind", "go source"},
				[2]string{"read", "yes"},
			),
			gtui.Row("gap-1 pt-1",
				gtui.Tag("", "tracked", nil),
				gtui.Tag("", "staged", func(e *gtui.Event) { host.Toast("info", "unstaged") }),
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
