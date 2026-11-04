package main

import "github.com/orayew2002/gtui"

func col(key, width string, right bool) gtui.Col {
	return gtui.Col{Label: gtui.T("col." + key), Width: width, Right: right}
}

func statusPill(status string) gtui.Component {
	switch status {
	case "ready", "done", "active", "paid":
		return gtui.Pill("text-success", status)
	case "encoding", "pending", "queued":
		return gtui.Pill("text-warning", status)
	case "failed", "error", "expired":
		return gtui.Pill("text-danger", status)
	}
	return gtui.Pill("text-muted", status)
}

func cells(values ...string) []gtui.Component {
	out := make([]gtui.Component, 0, len(values))

	for _, value := range values {
		switch value {
		case "on":
			out = append(out, gtui.Toggle("", true))
		case "off":
			out = append(out, gtui.Toggle("", false))
		case "ready", "done", "active", "encoding", "pending", "queued", "failed", "expired":
			out = append(out, statusPill(value))
		default:
			out = append(out, gtui.Text("text-ink", value))
		}
	}

	return out
}

func listPage(key string, cols []gtui.Col, rows [][]gtui.Component, info string, names []string) gtui.Component {
	label := func(index int) string {
		if index >= 0 && index < len(names) {
			return names[index]
		}
		return gtui.T("col.name")
	}

	return gtui.DataTable(gtui.TableOptions{
		Class:      "grow key-" + key,
		Cols:       cols,
		Rows:       rows,
		Zebra:      true,
		Compact:    true,
		Selectable: true,
		OnSelect:   func(index int) { app.Toast("", label(index)) },
		RowActions: []gtui.RowAction{
			{Label: "edit", Class: "text-accent", Run: func(index int) { editRow(label(index)) }},
			{Label: "del", Class: "text-danger", Run: func(index int) { deleteRow(label(index)) }},
		},
		Page: info,
	})
}

func editRow(name string) {
	open("editor")

	app.SetValue("#f-title", name)
	app.Toast("", gtui.T("action.edit")+": "+name)
}

func deleteRow(name string) {
	app.Modal(gtui.T("form.deleteTitle"),
		gtui.Text("text-muted", name),
		gtui.Text("text-muted dim", gtui.T("form.deleteBody")),
		gtui.Row("gap-1",
			gtui.Button("rounded border-danger text-danger focus:bold #row-yes", gtui.T("action.delete"), func(e *gtui.Event) {
				app.CloseModal()
				app.Toast("danger", gtui.T("form.deleted")+": "+name)
			}),
			gtui.Button("rounded border-line text-muted focus:bold #row-no", gtui.T("action.cancel"), func(e *gtui.Event) {
				app.CloseModal()
			}),
		),
	)
}

func names(values [][]string) []string {
	out := make([]string, 0, len(values))
	for _, row := range values {
		if len(row) > 0 {
			out = append(out, row[0])
		}
	}
	return out
}

func field(label string, input gtui.Component) gtui.Component {
	return gtui.Column("",
		gtui.Text("text-muted dim nowrap", label),
		input,
	)
}

func rows(values [][]string) [][]gtui.Component {
	out := make([][]gtui.Component, 0, len(values))
	for _, row := range values {
		out = append(out, cells(row...))
	}
	return out
}
