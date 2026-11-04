// An admin panel: a table of records with row actions, dialogs that ask before
// they act, and the numbers above it.
package main

import "github.com/orayew2002/gtui"

var app = gtui.Init()

type member struct {
	name, role, state string
}

var members = []member{
	{"Ada Lovelace", "analyst", "ready"},
	{"Alan Turing", "engineer", "pending"},
	{"Grace Hopper", "captain", "ready"},
	{"Edsger Dijkstra", "professor", "failed"},
}

func main() {
	// View registers the builder, so Refresh redraws the whole screen from the
	// data after anything changes it
	app.View(view)
	app.QuitOn("q", "ctrl+c")
	app.Render()
}

func view() gtui.Component {
	return gtui.Column("full bg-canvas",
		gtui.Navbar("px-1 border-b border-line",
			gtui.IconText("text-accent bold nowrap", "terminal", "admin"),
			gtui.IconButton("text-muted hover:text-ink", "search", "", func(e *gtui.Event) {
				app.Prompt("search", "name", "who are you looking for", func(v string) {
					app.Toast("info", "searched "+v)
				})
			}),
			gtui.IconButton("text-muted hover:text-ink", "bell", "", func(e *gtui.Event) {
				app.Toast("info", "nothing new")
			}),
			gtui.Avatar("bg-accent text-canvas", "Ada Lovelace"),
			gtui.Kbd("", "q"),
		),
		gtui.Row("px-1 items-center gap-2",
			gtui.Breadcrumb("text-muted", "team", "members"),
			gtui.Spacer(),
			gtui.Text("text-muted dim nowrap", "4 of 4"),
		),
		gtui.Row("px-1 gap-1",
			gtui.Metric("grow", "members", "4", "+1"),
			gtui.Metric("grow", "ready", "2", "flat"),
			gtui.Metric("grow", "failed", "1", "-2"),
		),
		gtui.Row("grow clip px-1 py-1 gap-1",
			gtui.Div("grow clip col", table()),
			side(),
		),
		gtui.StatusBar("border-t border-line px-1",
			[]gtui.Component{gtui.Row("gap-1", gtui.StatusDot("", "ok"), gtui.Text("nowrap", "connected"))},
			[]gtui.Component{gtui.Text("nowrap text-muted dim", "enter acts · q quits")},
		),
	)
}

func table() gtui.Component {
	rows := make([][]gtui.Component, 0, len(members))
	for _, who := range members {
		rows = append(rows, []gtui.Component{
			gtui.Text("truncate", who.name),
			gtui.Text("text-muted truncate", who.role),
			gtui.Row("gap-1 items-center",
				gtui.StatusDot("", state(who.state)),
				gtui.Text("nowrap", who.state),
			),
		})
	}

	return gtui.DataTable(gtui.TableOptions{
		Title:   "members",
		Caption: "edit renames, del asks first",
		Actions: []gtui.Component{
			gtui.Dropdown("text-muted", "bulk",
				gtui.MenuItem("", "mark ready", func(e *gtui.Event) { app.Toast("success", "all ready") }),
				gtui.MenuItem("", "export", func(e *gtui.Event) { app.Toast("info", "exported") }),
			),
			gtui.Button("pill border-accent text-accent", "+ new", func(e *gtui.Event) {
				app.Prompt("new member", "name", "full name", func(v string) {
					members = append(members, member{v, "guest", "pending"})
					app.Refresh()
				})
			}),
		},
		Cols: []gtui.Col{
			{Label: "name", Width: "grow"},
			{Label: "role", Width: "w-10"},
			{Label: "state", Width: "w-12"},
		},
		Rows:   rows,
		Footer: gtui.Cells("4 members", "", "1 failed"),
		Zebra:  true,
		RowActions: []gtui.RowAction{
			{Label: "edit", Class: "text-accent", Run: rename},
			{Label: "del", Class: "text-danger", Run: remove},
		},
		OnSelect: func(at int) { app.Toast("info", members[at].name) },
	})
}

func side() gtui.Component {
	return gtui.Panel("w-26 shrink-0", "picked",
		[]gtui.Component{gtui.Badge("bg-accent text-canvas", "live")},
		gtui.KeyValue("",
			[2]string{"name", members[0].name},
			[2]string{"role", members[0].role},
			[2]string{"since", "2021"},
		),
		gtui.Separator("", "tags"),
		gtui.Row("gap-1 flex-wrap",
			gtui.Tag("", "owner", func(e *gtui.Event) { app.Toast("info", "tag gone") }),
			gtui.Tag("", "on call", nil),
		),
	)
}

func rename(at int) {
	app.Prompt("rename", "name", members[at].name, func(v string) {
		members[at].name = v
		app.Refresh()
		app.Toast("success", "renamed")
	})
}

func remove(at int) {
	app.Confirm("remove "+members[at].name+"?", "the row goes for good", func() {
		members = append(members[:at], members[at+1:]...)
		app.Refresh()
		app.Toast("danger", "removed")
	})
}

func state(name string) string {
	switch name {
	case "ready":
		return "ok"
	case "failed":
		return "down"
	}
	return "busy"
}
