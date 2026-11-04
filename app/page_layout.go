// The layout primitives on their own: rows and columns, a grid that wraps, a
// split, scrolling in both directions, centring, cards and target-swapping tabs.
package main

import (
	"fmt"

	"github.com/orayew2002/gtui"
)

// pageLayout gives the scroll and centre boxes whatever rows the two rows of
// cards leave over, so the page fills the screen without overflowing it.
func pageLayout() gtui.Component {
	_, height := room()

	// the scroll and centre boxes take the rows the two rows of cards leave
	deep := height - 9
	if wide() {
		deep = height - 16
	}
	deep = min(max(deep, 3), 12)

	return gtui.Column("grow gap-1",
		gtui.Row("gap-1",
			demo("a", "row and column", `Row() and Column(), the last box has grow`,
				gtui.Row("gap-1",
					gtui.Div("w-6 h-2 rounded border-accent"),
					gtui.Div("w-6 h-2 rounded border-success"),
					gtui.Div("grow h-2 rounded border-warning"),
				),
			),
			demo("b", "grid that wraps", `Grid("", 3, tiles…) wraps into rows`,
				gtui.Grid("gap-0", 3, tiles()...),
			),
		),
		gtui.Row("gap-1",
			demo("d", "split", `Split("", "w-1/3", left, right)`,
				gtui.Split("", "w-1/3",
					gtui.Div("rounded border-line px-1", gtui.Text("nowrap", "side")),
					gtui.Div("rounded border-line px-1", gtui.Text("nowrap", "main area")),
				),
			),
			demo("f", "scroll both ways", "f jumps here, arrows scroll",
				gtui.Div(fmt.Sprintf("h-%d scroll scroll-x rounded border-line #both", deep),
					gtui.Column("gap-0", lines()...),
				),
			),
			demo("g", "centering", `Center("full", child)`,
				gtui.Center(fmt.Sprintf("h-%d rounded border-line", deep),
					gtui.Text("text-accent bold nowrap", "centered"),
				),
			),
		),
		cards(
			demo("h", "card", `Card("", title, body…) frames a titled box`,
				gtui.Card("", "disk-01",
					gtui.KeyValue("", [2]string{"pool", "fast"}, [2]string{"used", "61%"}),
				),
			),
			demo("j", "tabs that swap a box", "Tabs writes into the box it targets",
				gtui.Tabs("", "#tabbox",
					gtui.Tab{Key: "1", Label: "summary", Active: true, Page: []gtui.Component{
						gtui.Text("nowrap", "12 disks, 2 pools"),
					}},
					gtui.Tab{Key: "2", Label: "errors", Page: []gtui.Component{
						gtui.Text("nowrap text-danger", "disk-07 checksum failed"),
					}},
				),
				gtui.Div("rounded border-line px-1 mt-1 #tabbox",
					gtui.Text("nowrap", "12 disks, 2 pools"),
				),
			),
		),
	)
}

func tiles() []gtui.Component {
	class := "bg-raised px-1"
	if wide() {
		class = "rounded border-line bg-panel px-1"
	}

	cells := make([]gtui.Component, 0, 6)
	for i := range 6 {
		cells = append(cells, gtui.Div(class, gtui.Text("nowrap", "tile %d", i+1)))
	}

	return cells
}

func lines() []gtui.Component {
	rows := []string{
		"a wide line that keeps going past the edge of this box",
		"another wide line that also keeps going past the edge",
		"third", "fourth", "fifth", "sixth", "seventh", "eighth",
	}

	cells := make([]gtui.Component, 0, len(rows))
	for _, text := range rows {
		cells = append(cells, gtui.Text("nowrap", text))
	}

	return cells
}
