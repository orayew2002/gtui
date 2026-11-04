package main

import "github.com/orayew2002/gtui"

func pageStyle() gtui.Component {
	return card("a", "every class a component takes",
		gtui.TabView("grow #st",
			gtui.Tab{Key: "u", Label: "border", Active: true, Page: borderClasses()},
			gtui.Tab{Key: "v", Label: "colour", Page: colourClasses()},
			gtui.Tab{Key: "w", Label: "text", Page: textClasses()},
			gtui.Tab{Key: "x", Label: "space", Page: spaceClasses()},
			gtui.Tab{Key: "y", Label: "align", Page: alignClasses()},
			gtui.Tab{Key: "z", Label: "state", Page: stateClasses()},
		),
	)
}

func borderClasses() []gtui.Component {
	boxed := []string{"border", "rounded", "thick", "double", "dashed", "soft", "tight"}

	boxes := make([]gtui.Component, 0, len(boxed))
	for _, class := range boxed {
		boxes = append(boxes, gtui.Div(class+" px-1 border-line", gtui.Text("nowrap", class)))
	}

	return []gtui.Component{
		gtui.Div("row flex-wrap gap-1", boxes...),
		gtui.Row("gap-2 items-center",
			gtui.Text("pill border-line nowrap", "pill"),
			gtui.Text("nowrap bg-raised px-1", "rounded-none"),
			gtui.Text("dim truncate", "border-t|r|b|l|x|y and border-<token> pick the sides and the colour"),
			gtui.Div("rounded border-line shadow px-1 shrink-0", gtui.Text("nowrap", "shadow")),
		),
	}
}

func colourClasses() []gtui.Component {
	fills := make([]gtui.Component, 0, len(themeTokens))
	inks := make([]gtui.Component, 0, len(themeTokens))

	for _, name := range themeTokens {
		fills = append(fills, gtui.Text("px-1 nowrap text-canvas bg-"+name, name))
		inks = append(inks, gtui.Text("px-1 nowrap text-"+name, name))
	}

	// the same colour at every strength, over a light strip, so the mixing is
	// plain to see
	steps := make([]gtui.Component, 0, 5)
	for _, pct := range []string{"0", "25", "50", "75", "100"} {
		steps = append(steps, gtui.Text("px-1 nowrap text-ink bg-[#000000]/"+pct, pct+"%"))
	}

	return []gtui.Component{
		note("bg-<token>"),
		gtui.Div("row flex-wrap gap-1", fills...),
		note("text-<token>"),
		gtui.Div("row flex-wrap gap-1", inks...),
		gtui.Row("gap-1 items-center",
			gtui.Text("px-1 nowrap bg-[#7f5af0] text-canvas", "bg-[#7f5af0]"),
			gtui.Text("pill nowrap border-accent text-accent", "border-accent"),
			gtui.Text("pill nowrap border-danger text-danger", "border-danger"),
		),
		note("bg-<colour>/<percent> mixes with whatever is painted underneath"),
		gtui.Div("row bg-[#e8e8e8]", steps...),
		gtui.Row("gap-2 items-center",
			gtui.Text("px-1 nowrap bg-transparent", "bg-transparent"),
			gtui.Text("px-1 nowrap bg-opacity-0", "bg-opacity-0"),
			gtui.Text("px-1 nowrap bg-accent/40 text-ink", "bg-accent/40"),
			gtui.Text("dim truncate", "both hand the cells back to whatever is behind"),
		),
	}
}

func textClasses() []gtui.Component {
	return []gtui.Component{
		gtui.Row("gap-1",
			gtui.Text("w-16 nowrap uppercase", "uppercase"),
			gtui.Text("w-16 nowrap lowercase", "LOWERCASE"),
			gtui.Text("grow nowrap capitalize", "capitalize each word"),
		),
		gtui.Row("gap-1",
			gtui.Text("w-16 nowrap bold", "bold"),
			gtui.Text("w-16 nowrap dim", "dim"),
			gtui.Text("w-10 nowrap italic", "italic"),
			gtui.Text("grow nowrap underline", "underline"),
		),
		gtui.Row("gap-1",
			gtui.Text("w-16 nowrap strike", "strike"),
			gtui.Text("w-16 nowrap opacity-50", "opacity-50"),
			gtui.Text("grow nowrap text-right", "text-right"),
		),
		gtui.Row("gap-1 items-center",
			gtui.Text("w-16 truncate bg-raised", "truncate cuts this line at the box edge"),
			gtui.Text("dim nowrap", "truncate · nowrap · pre · wrap"),
		),
	}
}

func spaceClasses() []gtui.Component {
	pads := make([]gtui.Component, 0, 4)
	for _, class := range []string{"p-0", "p-1", "px-2", "py-1"} {
		pads = append(pads, gtui.Div("rounded border-line "+class, gtui.Text("nowrap", class)))
	}

	return []gtui.Component{
		note("padding, then the gap between children"),
		gtui.Div("row flex-wrap gap-1 items-start", pads...),
		gtui.Row("gap-0 items-center",
			gtui.Text("bg-accent text-canvas nowrap", " gap-0 "),
			gtui.Text("bg-success text-canvas nowrap", " no "),
			gtui.Text("bg-warning text-canvas nowrap", " space "),
		),
		gtui.Row("gap-2 items-center",
			gtui.Text("bg-accent text-canvas nowrap", " gap-2 "),
			gtui.Text("bg-success text-canvas nowrap", " two "),
			gtui.Text("bg-warning text-canvas nowrap", " cells "),
		),
		note("sizing: w-N h-N w-1/3 w-50% grow grow-N shrink-0 basis-N"),
		gtui.Row("gap-1",
			gtui.Text("w-8 bg-raised nowrap", "w-8"),
			gtui.Text("w-1/3 bg-line nowrap", "w-1/3"),
			gtui.Text("grow bg-raised nowrap", "grow"),
			gtui.Text("shrink-0 bg-line nowrap", "shrink-0"),
		),
	}
}

func alignClasses() []gtui.Component {
	rows := make([]gtui.Component, 0, 4)

	for _, class := range []string{"justify-start", "justify-center", "justify-end", "justify-between", "justify-evenly"} {
		rows = append(rows, gtui.Row("gap-1 bg-raised "+class,
			gtui.Text("bg-accent text-canvas nowrap", " a "),
			gtui.Text("bg-success text-canvas nowrap", " b "),
			gtui.Text("dim nowrap", class),
		))
	}

	return []gtui.Component{
		gtui.Column("gap-0", rows...),
		gtui.Row("gap-1 h-3 items-center rounded border-line",
			gtui.Text("bg-raised nowrap self-start", "self-start"),
			gtui.Text("bg-raised nowrap self-center", "items-center"),
			gtui.Text("bg-raised nowrap self-end", "self-end"),
		),
	}
}

func stateClasses() []gtui.Component {
	return []gtui.Component{
		note("prefix any class: hover: focus: active: current:"),
		gtui.Row("gap-1 flex-wrap",
			gtui.Clickable("px-1 rounded border-line hover:bg-raised focus:bold focus:border-accent active:bg-accent active:text-canvas #st-try",
				func(e *gtui.Event) { say("st-out", "clicked the focus demo") },
				gtui.Text("nowrap", "tab here, then enter"),
			),
			gtui.Href("px-1 rounded border-line current:bold current:text-accent #st-link", "a link keeps current:", "#st-out",
				gtui.Text("text-accent nowrap", "the link marked itself current"),
			),
		),
		out("st-out", "nothing yet"),
		gtui.Hint("", "focusable", "add it to any div to make it tabbable"),
	}
}
