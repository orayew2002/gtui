package main

import (
	"time"

	"github.com/orayew2002/gtui"
)

var disks = []string{"disk-01", "disk-02", "disk-03", "usb pool", "nas-01"}

func pagePickers() gtui.Component {
	return gtui.Column("grow gap-1",
		gtui.Row("gap-1",
			card("a", "selector",
				note("enter opens, type to filter, enter picks"),
				gtui.Selector("#pick-disk", "disk", "pick a disk", disks...),
				out("pick-out", "nothing picked"),
			),
			card("b", "autocomplete",
				note("type to see suggestions, arrows and enter accept"),
				gtui.Autocomplete("#pick-find", "suggest", "type a disk name", disks...),
				out("find2-out", "—"),
			),
		),
		gtui.Row("gap-1",
			card("d", "range and radios",
				note("arrows move the slider; enter picks a radio"),
				gtui.Range("#pick-range", "range", 0, 100, 5, 40),
				gtui.Radios("", "quality", "draft", "standard", "high"),
				gtui.Button("rounded border-line focus:bold", "read", func(e *gtui.Event) {
					values := app.Values()
					say("range-out", "range=%s quality=%s", values["range"], values["quality"])
				}),
				out("range-out", "—"),
			),
			card("e", "calendar",
				note("arrows move days, pgup/pgdn months, enter selects"),
				gtui.Calendar("#pick-day", "day", time.Now()),
			),
		),
	)
}
