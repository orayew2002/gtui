package main

import "github.com/orayew2002/gtui"

func pageButtons() gtui.Component {
	return gtui.Column("grow gap-1",
		gtui.Row("gap-1",
			card("a", "buttons",
				note("click or press enter"),
				gtui.Row("gap-1 flex-wrap",
					gtui.Button("rounded border-accent text-accent focus:bold", "primary", func(e *gtui.Event) {
						say("btn-out", "primary clicked")
					}),
					gtui.Button("rounded border-line text-muted focus:bold", "secondary", func(e *gtui.Event) {
						say("btn-out", "secondary clicked")
					}),
					gtui.Button("rounded border-danger text-danger focus:bold", "danger", func(e *gtui.Event) {
						say("btn-out", "danger clicked")
					}),
				),
				out("btn-out", "nothing clicked yet"),
			),
			card("b", "icon buttons",
				note("the same button with an icon label"),
				gtui.Row("gap-1 flex-wrap",
					gtui.Button("rounded border-line focus:bold", "✎ edit", func(e *gtui.Event) {
						say("icon-out", "edit")
					}),
					gtui.Button("rounded border-line focus:bold", "⌫ delete", func(e *gtui.Event) {
						say("icon-out", "delete")
					}),
					gtui.Button("rounded border-line focus:bold", "⤓ save", func(e *gtui.Event) {
						say("icon-out", "save")
					}),
				),
				out("icon-out", "—"),
			),
		),
		gtui.Row("gap-1",
			card("d", "dropdown",
				note("enter opens, arrows move, enter picks"),
				gtui.Dropdown("#dd", "actions",
					gtui.MenuItem("text-ink", "rename", func(e *gtui.Event) { say("dd-out", "rename") }),
					gtui.MenuItem("text-ink", "duplicate", func(e *gtui.Event) { say("dd-out", "duplicate") }),
					gtui.MenuItem("text-danger", "delete", func(e *gtui.Event) { say("dd-out", "delete") }),
				),
				out("dd-out", "—"),
			),
			card("e", "links and tabs",
				note("links swap the box below"),
				gtui.Tabs("", "#tab-out",
					gtui.Tab{Key: "o", Label: "first", Page: []gtui.Component{gtui.Text("text-accent", "first page")}},
					gtui.Tab{Key: "p", Label: "second", Page: []gtui.Component{gtui.Text("text-accent", "second page")}},
				),
				gtui.Div("#tab-out", gtui.Text("text-accent", "first page")),
			),
			card("f", "palette",
				note("ctrl+p anywhere, or click"),
				gtui.Button("rounded border-accent text-accent focus:bold", "open palette", func(e *gtui.Event) {
					openPalette()
				}),
				gtui.Hint("", "ctrl+p", "command palette"),
			),
		),
	)
}
