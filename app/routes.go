package main

import "github.com/orayew2002/gtui"

type route struct {
	Section string
	Key     string
	Label   string
	Page    func() gtui.Component
	Action  string
}

var routes = []route{
	{Section: "sections.main", Key: "1", Label: "home", Page: pageHome},
	{Key: "y", Label: "dynamic", Page: pageDynamic},

	{Section: "sections.content", Key: "2", Label: "movies", Page: pageMovies, Action: "create"},
	{Key: "3", Label: "serials", Page: pageSerials, Action: "create"},
	{Key: "4", Label: "books", Page: pageBooks, Action: "create"},
	{Key: "5", Label: "songs", Page: pageSongs, Action: "create"},
	{Key: "6", Label: "editor", Page: pageForm, Action: "save"},

	{Section: "sections.system", Key: "u", Label: "uikit", Page: pageUIKit},
	{Key: "s", Label: "storages", Page: pageStorages, Action: "create"},
	{Key: "l", Label: "logs", Page: pageLogs, Action: "export"},
	{Key: "g", Label: "settings", Page: pageSettings, Action: "save"},
}

var active = routes[0]

func show(label string) {
	for _, item := range routes {
		if item.Label == label {
			active = item
			return
		}
	}
}

func onNavigate(e *gtui.Event) {
	for _, item := range routes {
		if e.Target.Get("id") != "nav"+item.Key {
			continue
		}

		active = item
		app.Open("#head", header(item)...)
	}
}

func openPalette() {
	actions := make([]gtui.Action, 0, len(routes))

	for _, item := range routes {
		page := item
		actions = append(actions, gtui.Action{
			Key:   page.Key,
			Label: gtui.T(page.Label),
			Run: func() {
				active = page
				app.Open("#main", page.Page())
				app.Open("#head", header(page)...)
				app.Select("#nav" + page.Key)
			},
		})
	}

	app.Palette(gtui.T("search"), actions...)
}
