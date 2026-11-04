package main

import "github.com/orayew2002/gtui"

func pageButtons() gtui.Component {
	return gtui.Column("grow gap-1",
		demo("a", "button variants", "one look per button, written only with classes",
			gtui.Row("gap-x-2 gap-y-1 items-center flex-wrap",
				gtui.Button("pill bg-accent text-canvas border-accent", "primary", func(e *gtui.Event) {
					say("btn-out", "primary clicked")
				}),
				gtui.Button("pill border-line text-muted", "secondary", func(e *gtui.Event) {
					say("btn-out", "secondary clicked")
				}),
				gtui.Button("pill border-danger text-danger", "danger", func(e *gtui.Event) {
					say("btn-out", "danger clicked")
				}),
				gtui.Button("rounded shadow bg-panel text-ink border-line uppercase", "shadow", func(e *gtui.Event) {
					say("btn-out", "shadow clicked")
				}),
			),
			out("btn-out", "nothing clicked yet"),
		),
		cards(
			card("b", "groups and clickables",
				note("ButtonGroup joins buttons, Clickable wraps a row"),
				gtui.ButtonGroup("",
					gtui.IconButton("", "edit", "edit", func(e *gtui.Event) { say("icon-out", "edit") }),
					gtui.IconButton("", "save", "save", func(e *gtui.Event) { say("icon-out", "save") }),
					gtui.IconButton("text-danger", "trash", "delete", func(e *gtui.Event) { say("icon-out", "delete") }),
				),
				gtui.Clickable("pill px-2 gap-1 border-accent text-accent", func(e *gtui.Event) {
					say("icon-out", "the whole row is the button")
				},
					gtui.Icon("", "star"),
					gtui.Text("truncate", "rate this gallery"),
				),
				out("icon-out", "—"),
			),
			card("d", "dropdown",
				note("enter opens, arrows move, enter picks"),
				gtui.Dropdown("#dd", "actions",
					gtui.MenuItem("text-ink", "rename", func(e *gtui.Event) { say("dd-out", "rename") }),
					gtui.MenuItem("text-ink", "duplicate", func(e *gtui.Event) { say("dd-out", "duplicate") }),
					gtui.MenuItem("text-danger", "delete", func(e *gtui.Event) { say("dd-out", "delete") }),
				),
				out("dd-out", "—"),
			),
		),
		cards(
			card("f", "tab view",
				note("click a tab, or press o and p"),
				gtui.TabView("#btn-tabs",
					gtui.Tab{Key: "o", Icon: "file", Label: "first", Page: []gtui.Component{
						gtui.Text("text-accent", "first page"),
					}},
					gtui.Tab{Key: "p", Icon: "folder", Label: "second", Page: []gtui.Component{
						gtui.Text("text-accent", "second page"),
					}},
				),
			),
			card("g", "dialogs",
				note("each one takes focus and gives it back"),
				gtui.Row("gap-x-2 gap-y-1 flex-wrap",
					gtui.Button("pill border-accent text-accent", "prompt", func(e *gtui.Event) {
						app.Prompt("rename", "new name", "disk-04", func(value string) {
							say("btn-out", "renamed to %q", value)
						})
					}),
					gtui.Button("pill border-danger text-danger", "confirm", func(e *gtui.Event) {
						app.Confirm("delete the pool?", "every snapshot goes with it", func() {
							say("btn-out", "pool deleted")
						})
					}),
					gtui.Button("pill border-line text-muted", "palette", func(e *gtui.Event) {
						openPalette()
					}),
				),
				gtui.Hint("", "ctrl+p", "command palette"),
			),
		),
	)
}
