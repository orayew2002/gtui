package main

import "github.com/orayew2002/gtui"

type route struct {
	Section string
	Key     string
	Label   string
	Page    func() gtui.Component
	Action  string
}

// routes is a function, not a variable, so pages are free to navigate back
// through it without creating an initialisation cycle.
func routes() []route {
	return []route{
		{Section: "sections.main", Key: "1", Label: "home", Page: pageHome},
		{Key: "y", Label: "dynamic", Page: pageDynamic},

		{Section: "sections.content", Key: "2", Label: "movies", Page: pageMovies, Action: "create"},
		{Key: "3", Label: "serials", Page: pageSerials, Action: "create"},
		{Key: "4", Label: "books", Page: pageBooks, Action: "create"},
		{Key: "5", Label: "songs", Page: pageSongs, Action: "create"},
		{Key: "6", Label: "editor", Page: pageForm, Action: "save"},

		{Section: "sections.system", Key: "u", Label: "uikit", Page: pageUIKit},
		{Key: "i", Label: "widgets", Page: pageWidgets},
		{Key: "s", Label: "storages", Page: pageStorages, Action: "create"},
		{Key: "l", Label: "logs", Page: pageLogs, Action: "export"},
		{Key: "g", Label: "settings", Page: pageSettings, Action: "save"},
	}
}

var activeLabel = "home"

func active() route {
	for _, item := range routes() {
		if item.Label == activeLabel {
			return item
		}
	}
	return routes()[0]
}

func show(label string) {
	for _, item := range routes() {
		if item.Label == label {
			activeLabel = label
			return
		}
	}
}

// open switches the page, the header and the selected nav item together.
func open(label string) {
	show(label)
	item := active()

	app.Open("#main", item.Page())
	app.Open("#head", header(item)...)
	app.Select("#nav" + item.Key)
}

func onNavigate(e *gtui.Event) {
	for _, item := range routes() {
		if e.Target.Get("id") == "nav"+item.Key {
			show(item.Label)
			app.Open("#head", header(item)...)
		}
	}
}

func openPalette() {
	actions := make([]gtui.Action, 0, 12)

	for _, item := range routes() {
		label := item.Label
		actions = append(actions, gtui.Action{
			Key:   item.Key,
			Label: gtui.T(label),
			Run:   func() { open(label) },
		})
	}

	app.Palette(gtui.T("search"), actions...)
}
