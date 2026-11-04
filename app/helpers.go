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

func listPage(key string, cols []gtui.Col, rows [][]gtui.Component, info string) gtui.Component {
	return gtui.DataTable(gtui.TableOptions{
		Class:      "grow key-hide key-" + key,
		Cols:       cols,
		Rows:       rows,
		Zebra:      true,
		Compact:    true,
		Selectable: true,
		OnSelect:   func(index int) { app.Toast("", gtui.T("col.name")+" "+info) },
		Page:       info,
	})
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
