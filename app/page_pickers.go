// The controls that open something to choose from: a filterable selector, an
// autocomplete, a plain select, a range, radios and a calendar.
package main

import (
	"time"

	"github.com/orayew2002/gtui"
)

var disks = []string{"disk-01", "disk-02", "disk-03", "usb pool", "nas-01"}

// pagePickers covers the controls that open something: each one reports what it
// picked into a live line, so a click or a keystroke visibly does something.
func pagePickers() gtui.Component {
	return gtui.Column("grow gap-1",
		cards(
			card("a", "selector",
				note("enter opens, type to filter, enter picks"),
				gtui.Selector("#pick-disk", "disk", "pick a disk", disks...).
					On("change", func(e *gtui.Event) { say("pick-out", "picked %s", e.Value) }),
				out("pick-out", "nothing picked"),
			),
			card("b", "autocomplete",
				note("type, arrows move, enter accepts"),
				gtui.Autocomplete("#pick-find", "suggest", "type a disk name", disks...).
					On("input", func(e *gtui.Event) { say("find2-out", "filtering %q", e.Value) }).
					On("change", func(e *gtui.Event) { say("find2-out", "accepted %s", e.Value) }),
				out("find2-out", "—"),
			),
		),
		cards(
			card("d", "range and radios",
				note("arrows move the slider, enter picks a radio"),
				gtui.Range("#pick-range", "range", 0, 100, 5, 40).
					On("change", func(e *gtui.Event) { say("range-out", "range=%s", e.Value) }),
				gtui.Radios("#pick-quality", "quality", "draft", "standard", "high").
					On("change", func(e *gtui.Event) { say("range-out", "quality=%s", e.Value) }),
				out("range-out", "—"),
			),
			card("f", "calendar",
				note("←→ a day, ↑↓ a week, ⇞⇟ a month"),
				gtui.Calendar("#pick-day", "day", time.Now()).
					On("change", func(e *gtui.Event) { say("day-out", "picked %s", e.Value) }),
				out("day-out", "enter selects a day"),
			),
		),
		card("g", "select and choose",
			note("Select drops a list down, Choose opens it as a dialog"),
			gtui.Row("gap-x-2 gap-y-1 items-center flex-wrap",
				gtui.Select("#pick-sort", "sort", "sort by", "newest", "oldest", "name").
					On("change", func(e *gtui.Event) { say("sort-out", "sort by %s", e.Value) }),
				gtui.Button("pill border-accent text-accent", "choose a disk", func(e *gtui.Event) {
					app.Choose("which disk?", disks, func(picked string) {
						say("sort-out", "chose %s", picked)
					})
				}),
			),
			out("sort-out", "—"),
		),
	)
}
