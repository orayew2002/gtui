package main

import "github.com/orayew2002/gtui"

func pageLists() gtui.Component {
	return gtui.Column("grow gap-1",
		gtui.Row("gap-1 grow",
			card("a", "tree",
				note("enter or click folds a branch"),
				gtui.Tree("",
					gtui.Node2{Label: "project", Open: true, Children: []gtui.Node2{
						{Label: "cmd", Children: []gtui.Node2{{Label: "main.go"}}},
						{Label: "internal", Children: []gtui.Node2{
							{Label: "store", Children: []gtui.Node2{{Label: "db.go"}}},
							{Label: "api.go"},
						}},
						{Label: "README.md"},
					}},
				),
			),
			card("b", "accordion",
				note("enter opens one section"),
				gtui.Accordion("",
					gtui.Section{Title: "general", Open: true, Body: []gtui.Component{
						gtui.KeyValue("", [2]string{"theme", "dark"}, [2]string{"width", "120"}),
					}},
					gtui.Section{Title: "advanced", Body: []gtui.Component{
						gtui.Text("text-muted", "hidden until you open it"),
					}},
					gtui.Section{Title: "about", Body: []gtui.Component{
						gtui.Text("text-muted", "gtui gallery"),
					}},
				),
			),
			card("d", "scroll list",
				note("1000 rows, only the visible ones are built"),
				gtui.ScrollList("#big-list", 8, 1000, func(i int) gtui.Component {
					return gtui.Row("gap-1",
						gtui.Text("w-6 dim nowrap", "%04d", i),
						gtui.Text("grow nowrap", "row number %d", i),
					)
				}),
				gtui.Hint("", "↑↓", "scroll after jumping here"),
			),
		),
		gtui.Row("gap-1",
			card("e", "key and value",
				gtui.KeyValue("",
					[2]string{"framework", "gtui"},
					[2]string{"language", "go"},
					[2]string{"license", "MIT"},
				),
			),
			card("f", "chips and badges",
				gtui.Chips("", "go", "tui", "flexbox", "ansi"),
				gtui.Row("gap-1 items-center",
					gtui.Badge("bg-accent text-canvas", " new "),
					gtui.Badge("bg-success text-canvas", " 12 "),
					gtui.Avatar("bg-accent text-canvas", "Ada Lovelace"),
					gtui.Rating("", 4, 5),
				),
			),
			card("g", "list rows",
				gtui.List("",
					gtui.ListItem("", "first", "1"),
					gtui.ListItem("", "second", "2"),
					gtui.ListItem("", "third", "3"),
				),
				gtui.Breadcrumb("text-muted", "gallery", "lists", "rows"),
			),
		),
	)
}
