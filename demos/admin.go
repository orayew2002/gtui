// An admin panel: a adminTable of records with row actions, dialogs that ask before
// they act, and the numbers above it.
package demos

import "github.com/orayew2002/gtui"

type adminMember struct {
	name, role, adminState string
}

var adminMembers = []adminMember{
	{"Ada Lovelace", "analyst", "ready"},
	{"Alan Turing", "engineer", "pending"},
	{"Grace Hopper", "captain", "ready"},
	{"Edsger Dijkstra", "professor", "failed"},
}

func Admin(a *gtui.App, _ int) gtui.Component {
	host = a

	return gtui.Column("full",
		gtui.Navbar("px-1 border-b border-line",
			gtui.IconText("text-accent bold nowrap", "terminal", "admin"),
			gtui.IconButton("text-muted hover:text-ink", "search", "", func(e *gtui.Event) {
				host.Prompt("search", "name", "who are you looking for", func(v string) {
					host.Toast("info", "searched "+v)
				})
			}),
			gtui.IconButton("text-muted hover:text-ink", "bell", "", func(e *gtui.Event) {
				host.Toast("info", "nothing new")
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
		adminBody(a),
		gtui.StatusBar("border-t border-line px-1",
			[]gtui.Component{gtui.Row("gap-1", gtui.StatusDot("", "ok"), gtui.Text("nowrap", "connected"))},
			[]gtui.Component{gtui.Text("nowrap text-muted dim", "enter acts · q quits")},
		),
	)
}

// adminBody puts the table and the panel side by side, or stacks them when the
// screen is too tight for both.
func adminBody(a *gtui.App) gtui.Component {
	// stacked, neither panel may be squeezed: the column scrolls instead
	if narrow(a) {
		return gtui.Column("grow clip scroll px-1 py-1 gap-1",
			gtui.Div("shrink-0 col", adminTable()),
			adminSide("shrink-0"),
		)
	}

	return gtui.Row("grow clip px-1 py-1 gap-1",
		gtui.Div("grow clip col", adminTable()),
		adminSide("w-26 shrink-0"),
	)
}

func adminTable() gtui.Component {
	rows := make([][]gtui.Component, 0, len(adminMembers))
	for _, who := range adminMembers {
		rows = append(rows, []gtui.Component{
			gtui.Text("truncate", who.name),
			gtui.Text("text-muted truncate", who.role),
			gtui.Row("gap-1 items-center",
				gtui.StatusDot("", adminState(who.adminState)),
				gtui.Text("nowrap", who.adminState),
			),
		})
	}

	return gtui.DataTable(gtui.TableOptions{
		Title:   "members",
		Caption: "edit renames, del asks first",
		Actions: []gtui.Component{
			gtui.Dropdown("text-muted", "bulk",
				gtui.MenuItem("", "mark ready", func(e *gtui.Event) { host.Toast("success", "all ready") }),
				gtui.MenuItem("", "export", func(e *gtui.Event) { host.Toast("info", "exported") }),
			),
			gtui.Button("pill border-accent text-accent", "+ new", func(e *gtui.Event) {
				host.Prompt("new member", "name", "full name", func(v string) {
					adminMembers = append(adminMembers, adminMember{v, "guest", "pending"})
					host.Refresh()
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
			{Label: "edit", Class: "text-accent", Run: adminRename},
			{Label: "del", Class: "text-danger", Run: adminRemove},
		},
		OnSelect: func(at int) { host.Toast("info", adminMembers[at].name) },
	})
}

func adminSide(class string) gtui.Component {
	return gtui.Panel(class, "picked",
		[]gtui.Component{gtui.Badge("bg-accent text-canvas", "live")},
		gtui.KeyValue("",
			[2]string{"name", adminMembers[0].name},
			[2]string{"role", adminMembers[0].role},
			[2]string{"since", "2021"},
		),
		gtui.Separator("", "tags"),
		gtui.Row("gap-1 flex-wrap",
			gtui.Tag("", "owner", func(e *gtui.Event) { host.Toast("info", "tag gone") }),
			gtui.Tag("", "on call", nil),
		),
	)
}

func adminRename(at int) {
	host.Prompt("rename", "name", adminMembers[at].name, func(v string) {
		adminMembers[at].name = v
		host.Refresh()
		host.Toast("success", "renamed")
	})
}

func adminRemove(at int) {
	host.Confirm("remove "+adminMembers[at].name+"?", "the row goes for good", func() {
		adminMembers = append(adminMembers[:at], adminMembers[at+1:]...)
		host.Refresh()
		host.Toast("danger", "removed")
	})
}

func adminState(name string) string {
	switch name {
	case "ready":
		return "ok"
	case "failed":
		return "down"
	}
	return "busy"
}
