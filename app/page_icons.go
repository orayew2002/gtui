package main

import "github.com/orayew2002/gtui"

func pageIcons() gtui.Component {
	names := gtui.IconNames()
	cells := make([]gtui.Component, 0, len(names))

	for _, name := range names {
		cells = append(cells, gtui.IconText("w-16 text-muted", name, name))
	}

	return gtui.Column("grow gap-1",
		card("a", "icons",
			note("Icon(class, name) and IconText(class, name, label) — %d of them", len(names)),
			gtui.Div("row flex-wrap gap-1", cells...),
		),
		gtui.Row("gap-1",
			card("b", "in buttons and rows",
				gtui.Row("gap-1 flex-wrap",
					gtui.Button("rounded border-line focus:bold", "⤓ save", nil),
					gtui.Button("rounded border-line focus:bold", "✎ edit", nil),
					gtui.Button("rounded border-danger text-danger focus:bold", "⌫ delete", nil),
				),
				gtui.List("",
					gtui.ListItem("", "⚙ settings", "ctrl+,"),
					gtui.ListItem("", "⌕ search", "ctrl+p"),
					gtui.ListItem("", "⏻ quit", "q"),
				),
			),
		),
	)
}
