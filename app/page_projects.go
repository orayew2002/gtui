package main

// The projects page is the odd one out: every other page demonstrates
// components, this one lists the example apps in demos/ and hands the whole
// page over to the one you click. The same function draws that app when you run
// it on its own with go run ./examples/<name>.

import (
	"github.com/orayew2002/gtui"
	"github.com/orayew2002/gtui/demos"
)

// openProject is the project filling the page, empty while the list is showing.
var openProject = ""

func pageProjects() gtui.Component {
	if project, ok := projectByKey(openProject); ok {
		return gtui.Column("grow clip gap-1",
			// one bar to get back out, so the project below owns the rest
			gtui.Row("gap-2 items-center shrink-0",
				gtui.Button("pill border-accent text-accent nowrap shrink-0 #pj-back", "‹ projects", func(e *gtui.Event) {
					openProject = ""
					show(current())
				}),
				gtui.Text("text-ink bold nowrap shrink-0", project.Name),
				// only the description gives way on a narrow screen: the way
				// out and the name of the thing you are in must not
				gtui.Text("text-muted dim truncate", project.About),
			),
			// the project draws itself from the gallery's tick, the same way it
			// draws itself from its own command's timer
			gtui.Div("grow clip rounded border-line #pj-open", project.Build(app, frame)),
		)
	}

	return gtui.Column("grow clip gap-1",
		note("click one to open it here, or run it on its own: go run ./examples/<name>"),
		gtui.Grid("gap-1", projectColumns(), projectCards()...),
	)
}

func projectCards() []gtui.Component {
	all := demos.All()
	cards := make([]gtui.Component, 0, len(all))

	for _, project := range all {
		pick := project

		cards = append(cards, gtui.Clickable(
			"col clip rounded border-line bg-panel px-1 hover:border-accent focus:border-accent focus:bold #pj-"+pick.Key,
			func(e *gtui.Event) {
				openProject = pick.Key
				show(current())
			},
			gtui.Row("gap-1 items-center",
				gtui.Text("text-accent bold nowrap", pick.Key),
				gtui.Text("text-ink bold truncate", pick.Name),
			),
			gtui.Text("text-muted dim truncate", pick.About),
		))
	}

	return cards
}

// projectColumns fits as many cards across as the page is wide.
func projectColumns() int {
	width, _ := room()

	switch {
	case width < 56:
		return 1
	case width < 96:
		return 2
	default:
		return 3
	}
}

func projectByKey(key string) (demos.Project, bool) {
	if key == "" {
		return demos.Project{}, false
	}

	for _, project := range demos.All() {
		if project.Key == key {
			return project, true
		}
	}
	return demos.Project{}, false
}
