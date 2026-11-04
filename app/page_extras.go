// The pieces an app is framed with: navbars, panels, windows, banners, empty
// states, timelines, steps, metrics and the small marks that go beside them.
package main

import "github.com/orayew2002/gtui"

// pageExtras is the catalogue page: it shows the framing pieces an app is built
// from rather than the controls a user touches.
func pageExtras() gtui.Component {
	return card("a", "the rest of the kit",
		gtui.TabView("grow #ex",
			gtui.Tab{Key: "u", Label: "shell", Active: true, Page: shellExtras()},
			gtui.Tab{Key: "v", Label: "panels", Page: panelExtras()},
			gtui.Tab{Key: "w", Label: "status", Page: statusExtras()},
			gtui.Tab{Key: "x", Label: "flow", Page: flowExtras()},
			gtui.Tab{Key: "y", Label: "bits", Page: bitExtras()},
			gtui.Tab{Key: "z", Label: "notices", Page: noticeExtras()},
		),
	)
}

func shellExtras() []gtui.Component {
	return []gtui.Component{
		gtui.Navbar("",
			gtui.IconText("text-accent bold", "terminal", "gtui"),
			gtui.Text("text-muted nowrap", "docs"),
			gtui.Text("text-ink nowrap", "download"),
			gtui.StatusDot("", "live"),
		),
		gtui.Window("", "Window(class, title, body…)",
			gtui.Text("text-muted nowrap", "the title rides the top border"),
			gtui.Centered("h-2", gtui.Text("text-accent nowrap", "Centered inside it")),
		),
	}
}

func panelExtras() []gtui.Component {
	return []gtui.Component{
		gtui.Panel("", "Panel with actions",
			[]gtui.Component{gtui.Text("dim nowrap", "⟳"), gtui.Text("dim nowrap", "⋯")},
			gtui.Text("text-muted nowrap", "a titled box with a row of actions"),
			gtui.Separator("", "Separator"),
			gtui.Text("text-muted nowrap", "and a rule with a label in it"),
		),
		gtui.Scrollbox("rounded border-line px-1", scrollRows(), logLines()...),
	}
}

func statusExtras() []gtui.Component {
	return []gtui.Component{
		gtui.Row("gap-1",
			gtui.Metric("grow", "requests", "18 402", "+12%"),
			gtui.Metric("grow", "errors", "12", "-3%"),
			gtui.Metric("grow", "p95", "41 ms", ""),
		),
		gtui.Row("gap-2 items-center",
			gtui.StatusDot("", "ok"),
			gtui.Text("text-muted nowrap", "ok"),
			gtui.StatusDot("", "warn"),
			gtui.Text("text-muted nowrap", "warn"),
			gtui.StatusDot("", "down"),
			gtui.Text("text-muted nowrap", "down"),
			gtui.Loader("grow", "Loader takes the animation tick", frame),
		),
	}
}

func flowExtras() []gtui.Component {
	return []gtui.Component{
		gtui.Steps("", 1,
			gtui.StepItem{Title: "plan"},
			gtui.StepItem{Title: "build"},
			gtui.StepItem{Title: "ship"},
		),
		gtui.Timeline("",
			gtui.TimelineItem{Icon: "check", Kind: "ok", Title: "deployed", Time: "09:12"},
			gtui.TimelineItem{Icon: "alert", Kind: "warn", Title: "retried once", Time: "09:14",
				Body: "the rail stops at the last entry"},
		),
		gtui.Columns("", []string{"w-12", ""},
			gtui.Text("text-muted nowrap", "Columns"),
			gtui.Text("text-ink truncate", "a fixed first cell, the rest shares what is left"),
		),
	}
}

func bitExtras() []gtui.Component {
	return []gtui.Component{
		gtui.Row("gap-1 items-center",
			gtui.Kbd("", "ctrl+p"),
			gtui.Kbd("", "esc"),
			gtui.Tag("", "go", nil),
			gtui.Tag("", "tui", func(e *gtui.Event) { say("bit-out", "removed tui") }),
			gtui.Text("dim truncate", "Kbd and Tag"),
		),
		gtui.Row("gap-1 items-start",
			gtui.ButtonGroup("",
				gtui.Button("focus:bold focus:text-accent", "cut", func(e *gtui.Event) { say("bit-out", "cut") }),
				gtui.Button("focus:bold focus:text-accent", "copy", func(e *gtui.Event) { say("bit-out", "copy") }),
				gtui.Button("focus:bold focus:text-accent", "paste", func(e *gtui.Event) { say("bit-out", "paste") }),
			),
			gtui.InputGroup("grow",
				gtui.Text("nowrap", "$"),
				gtui.Input("px-1 #ex-amount", "amount", "0.00"),
				gtui.Text("nowrap", "USD"),
			),
		),
		out("bit-out", "—"),
	}
}

func noticeExtras() []gtui.Component {
	banner := gtui.Banner("", "warning", "Banner", "a wide notice with its own buttons",
		gtui.Button("px-1 focus:bold focus:text-accent", "update", func(e *gtui.Event) {
			app.Toast("warning", "updating")
		}),
	)

	if !wide() {
		return []gtui.Component{banner, gtui.EmptyState("", "files", "EmptyState", "")}
	}

	return []gtui.Component{
		banner,
		gtui.EmptyState("rounded border-line", "files", "EmptyState", "nothing here yet",
			gtui.Button("px-1 focus:bold focus:text-accent", "add one", func(e *gtui.Event) {
				app.Toast("success", "added")
			}),
		),
	}
}

func scrollRows() int {
	_, height := room()
	return min(max(height-12, 3), 8)
}

func logLines() []gtui.Component {
	rows := []string{
		"Scrollbox keeps a tall list inside a short box",
		"jump here with v and the arrows scroll it",
		"03:14  worker started",
		"03:15  job 41 done",
		"03:16  job 42 done",
		"03:17  job 43 queued",
		"03:18  worker idle",
	}

	cells := make([]gtui.Component, 0, len(rows))
	for _, text := range rows {
		cells = append(cells, gtui.Text("nowrap", text))
	}

	return cells
}
