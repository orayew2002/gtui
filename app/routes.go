package main

import (
	"strings"

	"github.com/orayew2002/gtui"
)

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
		{Key: "l", Title: "i18n", About: "languages, translation keys and switching at runtime", Body: pageI18n},
		{Key: "s", Title: "style", About: "the utility classes every component takes", Body: pageStyle},
		{Key: "e", Title: "extras", About: "navbars, panels, windows, timelines and tags", Body: pageExtras},
		{Key: "v", Title: "views", About: "chat, logs, code, address and status bars", Body: pageViews},
	}
}

var currentKey = "1"

func current() page {
	if item, ok := byKey(currentKey); ok {
		return item
	}
	return pages()[0]
}

func byKey(key string) (page, bool) {
	for _, item := range pages() {
		if item.Key == key {
			return item, true
		}
	}
	return page{}, false
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

// show swaps the page, the bar above it and the selected nav item together.
// Every way of navigating ends here.
func show(item page) {
	currentKey = item.Key

	app.Open("#main", item.Body())
	app.Open("#head", topbar())
	app.Select("#nav" + item.Key)
}

// onNavigate handles a click on a sidebar link, which is also what the one key
// shortcuts and the palette trigger.
func onNavigate(e *gtui.Event) {
	id := e.Target.Get("id")
	if !strings.HasPrefix(id, "nav") {
		return
	}

	if item, ok := byKey(strings.TrimPrefix(id, "nav")); ok {
		show(item)
	}
}

func openPalette() {
	actions := make([]gtui.Action, 0, len(pages()))

	for _, item := range pages() {
		actions = append(actions, gtui.Action{Key: item.Key, Label: item.Title, Run: func() {
			show(item)
		}})
	}

	app.Palette("go to", actions...)
}
