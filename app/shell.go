package main

import "github.com/orayew2002/gtui"

// The shell is one bar across the top, a column of links down the left and the
// page itself. Each of the three is separated by a single rule rather than a
// box of its own, so the only borders on screen belong to the demos.

var seeThrough = false

func shell() gtui.Component {
	// seeThrough is the themes page switch: with it on the gallery paints no
	// background of its own and the terminal's shows through
	fill := "bg-canvas"
	if seeThrough {
		fill = "bg-transparent"
	}

	return gtui.Column("full "+fill,
		// the bar is swapped on every navigation, so the box that holds it
		// carries the id and the bar itself stays a plain component
		gtui.Div("col shrink-0 #head", topbar()),
		// clip is the min-width:0 of this layout: without it the content pane
		// is as wide as its widest child and pushes past the right edge
		gtui.Row("grow clip",
			sidebar(),
			gtui.Div("grow clip scroll px-1 py-1 gap-1 key-c key-hide #main", current().Body()),
		),
	)
}

// topbar carries the brand, where you are, and the two switches that belong to
// the whole app.
func topbar() gtui.Component {
	item := current()

	left := []gtui.Component{
		gtui.IconText("text-accent bold nowrap shrink-0", "terminal", phrase("app.name", "gtui")),
		gtui.Text("text-line nowrap shrink-0", "│"),
		gtui.Text("text-ink bold nowrap shrink-0", phrase("nav."+item.Title, item.Title)),
	}

	if !compact() {
		left = append(left, gtui.Text("text-muted dim truncate", item.About))
	}

	right := []gtui.Component{
		gtui.Spacer(),
		gtui.Button("shrink-0 text-muted focus:text-ink focus:bold #mode", "◐ "+app.ThemeName(),
			func(e *gtui.Event) { app.ToggleTheme() },
		),
		gtui.Text("text-muted dim nowrap shrink-0", "⌕ ctrl+p"),
	}

	return gtui.Column("shrink-0",
		gtui.Row("gap-2 items-center px-1", append(left, right...)...),
		gtui.Divider("border-line"),
	)
}

// sidebar is the page list. It keeps its own column of keys on a narrow screen
// and spells the titles out when there is room.
func sidebar() gtui.Component {
	roomy := !compact()

	items := make([]gtui.Component, 0, len(pages())+2)
	for _, item := range pages() {
		items = append(items, navLink(item, roomy))
	}

	width := "w-8"
	if roomy {
		width = "w-17"
	}

	// the list scrolls rather than squeezes: on a short terminal there are more
	// pages than rows, and a squeezed link is a link that disappears
	return gtui.Column(width+" shrink-0 clip scroll border-r border-line py-1 key-n key-hide #nav", items...)
}

// phrase is T with something to fall back on: T answers with the key itself
// when no translation is loaded, and a key is not something to show a user.
func phrase(key, fallback string) string {
	if text := gtui.T(key); text != key {
		return text
	}
	return fallback
}

// navLink is a link with an empty target: the click handler in routes.go does
// the swapping, so the sidebar costs nothing to build.
func navLink(item page, roomy bool) gtui.Component {
	// the icon says what the page is at a glance, the key says how to reach it,
	// and the title only fits while the sidebar is open. The glyph goes in the
	// label because Href's trailing arguments are the page it swaps in, not the
	// contents of the link.
	glyph, _ := gtui.IconGlyph(item.Icon)

	label := glyph + " " + item.Key
	if roomy {
		label += " " + phrase("nav."+item.Title, item.Title)
	}

	// the selected item is marked with a fill, not a border: a border would
	// take a column and shove the label sideways
	class := "row gap-1 px-1 shrink-0 truncate text-muted" +
		" hover:text-ink" +
		" focus:text-accent focus:bold" +
		" current:bg-raised current:text-ink current:bold" +
		" key-" + item.Key + " key-hide" +
		" #nav" + item.Key

	return gtui.Href(class, label, "")
}
