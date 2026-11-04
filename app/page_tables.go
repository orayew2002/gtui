package main

import "github.com/orayew2002/gtui"

var people = [][]string{
	{"Ada Lovelace", "analyst", "ready"},
	{"Alan Turing", "engineer", "pending"},
	{"Grace Hopper", "captain", "ready"},
	{"Edsger Dijkstra", "professor", "failed"},
}

func pageTables() gtui.Component {
	rows := make([][]gtui.Component, 0, len(people))

	for _, person := range people {
		rows = append(rows, []gtui.Component{
			gtui.Text("truncate", person[0]),
			gtui.Text("text-muted truncate", person[1]),
			gtui.Row("gap-1 items-center",
				gtui.StatusDot("", person[2]),
				gtui.Text("truncate", person[2]),
			),
		})
	}

	return gtui.Column("grow gap-1",
		gtui.DataTable(gtui.TableOptions{
			Class:   "key-a",
			Title:   "people",
			Caption: "zebra rows, selectable rows, row actions and a pager",
			Actions: []gtui.Component{
				gtui.IconButton("pill border-accent text-accent", "plus", "new", func(e *gtui.Event) {
					say("table-out", "new row")
				}),
			},
			Cols: []gtui.Col{
				{Label: "name", Width: "grow"},
				{Label: "role", Width: "w-9"},
				{Label: "status", Width: "w-9"},
			},
			Rows:       rows,
			Zebra:      true,
			Compact:    true,
			Selectable: true,
			OnSelect:   func(index int) { say("table-out", "selected %s", people[index][0]) },
			RowActions: []gtui.RowAction{
				{Label: "edit", Class: "text-accent", Run: func(index int) { say("table-out", "edit %s", people[index][0]) }},
				{Label: "del", Class: "text-danger", Run: func(index int) { confirmDelete(people[index][0]) }},
			},
			Page: "page 1 of 1",
			Prev: func(e *gtui.Event) { say("table-out", "previous page") },
			Next: func(e *gtui.Event) { say("table-out", "next page") },
		}),
		out("table-out", "click a row, or edit and del on the right"),
		gtui.Columns("", []string{"w-14", "w-14", "w-16"},
			gtui.Metric("", "people", "4", "+1"),
			gtui.Metric("", "failed", "1", "-2"),
			gtui.Metric("", "last sync", "09:40", ""),
		),
		cards(
			gtui.Window("key-b grow basis-0 clip", "plain table",
				gtui.Table("",
					[]gtui.Col{{Label: "key", Width: "w-10"}, {Label: "value", Width: "grow", Right: true}},
					gtui.Cells("rows", "4"),
					gtui.Cells("columns", "3"),
				),
			),
			gtui.DataTable(gtui.TableOptions{
				Class:   "key-d grow basis-0 clip",
				Title:   "backups",
				Cols:    []gtui.Col{{Label: "name", Width: "grow"}, {Label: "size", Width: "w-8", Right: true}},
				Footer:  gtui.Cells("total", "0 B"),
				Empty:   "no backups yet",
				Compact: true,
			}),
		),
		card("j", "a table built by hand",
			note("Table, or TableHead and TableRow, when DataTable is too much"),
			gtui.Table("", handCols,
				gtui.Cells("fast", "2", "61%"),
				gtui.Cells("archive", "10", "12%"),
			),
			gtui.Text("text-muted dim mt-1", "the same two rows, written row by row"),
			gtui.TableHead("", handCols),
			gtui.TableRow("", handCols, gtui.Cells("fast", "2", "61%")...),
			gtui.TableRow("bg-raised", handCols, gtui.Cells("archive", "10", "12%")...),
		),
	)
}

var handCols = []gtui.Col{
	{Label: "pool", Width: "w-12"},
	{Label: "disks", Width: "w-7", Right: true},
	{Label: "used", Width: "w-7", Right: true},
}

func confirmDelete(name string) {
	app.Confirm("delete "+name+"?", "the row goes away for good", func() {
		app.Toast("danger", "deleted "+name)
		say("table-out", "deleted %s", name)
	})
}
