package main

import "github.com/orayew2002/gtui"

func pageLayout() gtui.Component {
	tiles := make([]gtui.Component, 0, 6)
	for i := range 6 {
		tiles = append(tiles, gtui.Div("rounded border-line bg-panel px-2",
			gtui.Text("nowrap", "tile %d", i+1),
		))
	}

	return gtui.Column("grow gap-1",
		gtui.Row("gap-1",
			card("a", "row and column",
				note(`Row("gap-1", …) and Column("gap-1", …)`),
				gtui.Row("gap-1",
					gtui.Div("w-6 h-2 rounded border-accent"),
					gtui.Div("w-6 h-2 rounded border-success"),
					gtui.Div("grow h-2 rounded border-warning"),
				),
				note("the last box has grow, so it takes the rest"),
			),
			card("b", "grid that wraps",
				note(`Grid("", 3, tiles…) with flex-wrap`),
				gtui.Grid("", 3, tiles...),
			),
		),
		gtui.Row("gap-1 grow",
			card("d", "split",
				note(`Split("", "w-1/3", left, right)`),
				gtui.Split("", "w-1/3",
					gtui.Div("rounded border-line px-1", gtui.Text("nowrap", "side")),
					gtui.Div("rounded border-line px-1", gtui.Text("nowrap", "main area")),
				),
			),
			card("e", "scroll both ways",
				note("jump here, then arrows scroll it"),
				gtui.Div("h-6 scroll scroll-x rounded border-line #both",
					gtui.Column("gap-0",
						gtui.Text("nowrap", "a wide line that keeps going past the edge of this box"),
						gtui.Text("nowrap", "another wide line that also keeps going past the edge"),
						gtui.Text("nowrap", "third"),
						gtui.Text("nowrap", "fourth"),
						gtui.Text("nowrap", "fifth"),
						gtui.Text("nowrap", "sixth"),
						gtui.Text("nowrap", "seventh"),
						gtui.Text("nowrap", "eighth"),
					),
				),
			),
			card("f", "centering",
				note(`Center("full", child)`),
				gtui.Center("h-6 rounded border-line",
					gtui.Text("text-accent bold nowrap", "centered"),
				),
			),
		),
	)
}
