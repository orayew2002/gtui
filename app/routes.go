package main

import "github.com/orayew2002/gtui"

type page struct {
	Key   string
	Title string
	About string
	Body  func() gtui.Component
}

// pages is the whole gallery. Add a component demo by adding one entry here
// and one function in its own file.
func pages() []page {
	return []page{
		{Key: "1", Title: "buttons", About: "buttons, links, dropdowns and the command palette", Body: pageButtons},
		{Key: "2", Title: "inputs", About: "text, number, password, area, checkbox, switch", Body: pageInputs},
		{Key: "3", Title: "pickers", About: "selector, autocomplete, range, radios, calendar", Body: pagePickers},
		{Key: "4", Title: "tables", About: "data tables with zebra rows, actions and a pager", Body: pageTables},
		{Key: "5", Title: "lists", About: "trees, accordions, scroll lists and key/value rows", Body: pageLists},
		{Key: "6", Title: "feedback", About: "alerts, toasts, modals, progress and spinners", Body: pageFeedback},
		{Key: "7", Title: "charts", About: "bars, sparklines, line charts, gauges and donuts", Body: pageCharts},
		{Key: "8", Title: "media", About: "pictures, backdrops, video and music players", Body: pageMedia},
		{Key: "9", Title: "motion", About: "animated text, pulses, waves and tickers", Body: pageMotion},
		{Key: "0", Title: "layout", About: "rows, columns, grids, splits and scrolling", Body: pageLayout},
		{Key: "i", Title: "icons", About: "every icon the framework ships", Body: pageIcons},
		{Key: "t", Title: "themes", About: "theme tokens and the light and dark palettes", Body: pageThemes},
	}
}

var currentKey = "1"

func current() page {
	for _, item := range pages() {
		if item.Key == currentKey {
			return item
		}
	}
	return pages()[0]
}

// goTo switches the page by title, which is what the palette and the snapshot
// mode use.
func goTo(title string) {
	for _, item := range pages() {
		if item.Title == title {
			currentKey = item.Key
			return
		}
	}
}

func onNavigate(e *gtui.Event) {
	for _, item := range pages() {
		if e.Target.Get("id") == "nav"+item.Key {
			currentKey = item.Key
			app.Open("#head", header(item)...)
		}
	}
}

func openPalette() {
	actions := make([]gtui.Action, 0, len(pages()))

	for _, item := range pages() {
		title := item.Title
		actions = append(actions, gtui.Action{Key: item.Key, Label: title, Run: func() {
			goTo(title)
			show(current())
		}})
	}

	app.Palette("go to", actions...)
}

// show swaps the page, the header and the selected nav item together.
func show(item page) {
	app.Open("#main", item.Body())
	app.Open("#head", header(item)...)
	app.Select("#nav" + item.Key)
}
