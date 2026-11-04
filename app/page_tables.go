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

	for _, row := range people {
		rows = append(rows, []gtui.Component{
			gtui.Text("", row[0]),
			gtui.Text("text-muted", row[1]),
			statusPill(row[2]),
			gtui.Toggle("", row[2] == "ready"),
		})
	}

	return gtui.Column("grow gap-1",
		gtui.DataTable(gtui.TableOptions{
			Class:   "key-a",
			Title:   "data table",
			Caption: "zebra rows, selectable rows, row actions and a pager",
			Actions: []gtui.Component{
				gtui.Button("text-accent focus:bold", "＋ new", func(e *gtui.Event) { say("table-out", "new row") }),
			},
			Cols: []gtui.Col{
				{Label: "name", Width: "w-1/4"},
				{Label: "role", Width: "w-12"},
				{Label: "status", Width: "w-12"},
				{Label: "active", Width: "grow", Right: true},
			},
			Rows:       rows,
			Footer:     gtui.Cells("4 people", "", "", ""),
			Zebra:      true,
			Compact:    true,
			Selectable: true,
			OnSelect:   func(index int) { say("table-out", "selected %s", people[index][0]) },
			RowActions: []gtui.RowAction{
				{Label: "edit", Class: "text-accent", Run: func(index int) { say("table-out", "edit %s", people[index][0]) }},
				{Label: "del", Class: "text-danger", Run: func(index int) { confirmDelete(people[index][0]) }},
			},
			Page: "page 1 of 1",
		}),
		gtui.Row("gap-1",
			card("b", "result", out("table-out", "click a row, edit or del")),
			card("d", "simple table",
				note("Table is the plain version without the chrome"),
				gtui.Table("rounded border-line",
					[]gtui.Col{{Label: "key", Width: "w-10"}, {Label: "value", Width: "grow", Right: true}},
					gtui.Cells("rows", "4"),
					gtui.Cells("columns", "4"),
				),
			),
		),
	)
}

func statusPill(status string) gtui.Component {
	switch status {
	case "ready":
		return gtui.Pill("text-success", status)
	case "pending":
		return gtui.Pill("text-warning", status)
	}
	return gtui.Pill("text-danger", status)
}

func confirmDelete(name string) {
	app.Modal("delete row",
		gtui.Text("text-muted", name),
		gtui.Row("gap-1",
			gtui.Button("rounded border-danger text-danger focus:bold #yes", "delete", func(e *gtui.Event) {
				app.CloseModal()
				app.Toast("danger", "deleted "+name)
			}),
			gtui.Button("rounded border-line text-muted focus:bold #no", "cancel", func(e *gtui.Event) {
				app.CloseModal()
			}),
		),
	)
}
