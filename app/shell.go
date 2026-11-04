package main

import "github.com/orayew2002/gtui"

// narrow is the width under which the sidebar shrinks to its shortcut keys, so
// a page still has room on an 80 column terminal.
const narrow = 90

// compact reports whether the shell is in its narrow layout.
func compact() bool {
	width, _ := app.Size()
	return width < narrow
}

var seeThrough = false

func shell() gtui.Component {
	// seeThrough is the themes page switch: with it on the gallery paints no
	// background of its own and the terminal's shows through
	fill := "bg-canvas"
	if seeThrough {
		fill = "bg-transparent"
	}

	return gtui.Row("full "+fill+" p-1 gap-1",
		sidebar(),
		// clip is the min-width:0 of this layout: without it the content pane
		// is as wide as its widest child and pushes past the right edge.
		gtui.Column("grow clip gap-1",
			gtui.Div("row shrink-0 clip rounded border-line bg-panel px-2 gap-2 #head", header(current())...),
			gtui.Div("grow clip scroll rounded border-line key-c #main", current().Body()),
		),
	)
}

func sidebar() gtui.Component {
	wide := !compact()

	items := []gtui.Component{brand(wide), gtui.Divider("border-line")}

	for _, item := range pages() {
		items = append(items, navLink(item, wide))
	}

	items = append(items, gtui.Spacer())
	if wide {
		items = append(items, gtui.Divider("border-line"), gtui.IconText("text-muted", "search", "ctrl+p"))
	} else {
		items = append(items, gtui.Icon("text-muted", "search"))
	}

	width := "w-5"
	if wide {
		width = "w-18"
	}

	return gtui.Column(width+" shrink-0 clip rounded border-line bg-panel px-1 key-n #nav", items...)
}

func brand(wide bool) gtui.Component {
	if !wide {
		return gtui.Icon("text-accent bold", "terminal")
	}

	return gtui.Row("gap-1",
		gtui.Icon("text-accent bold", "terminal"),
		gtui.Text("grow text-ink bold truncate", phrase("app.name", "gtui gallery")),
	)
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
func navLink(item page, wide bool) gtui.Component {
	title, box := item.Key, "row text-center truncate text-muted"
	if wide {
		title, box = item.Key+" "+phrase("nav."+item.Title, item.Title), "row px-1 truncate text-muted"
	}

	class := box +
		" hover:text-ink" +
		" focus:bg-line focus:text-ink focus:bold" +
		" current:bg-line current:text-ink current:bold" +
		" key-" + item.Key + " key-hide" +
		" #nav" + item.Key

	return gtui.Href(class, title, "")
}

func header(item page) []gtui.Component {
	return []gtui.Component{
		gtui.Column("grow clip",
			gtui.Text("text-ink bold truncate", item.Title),
			gtui.Text("text-muted truncate", item.About),
		),
		gtui.Button("shrink-0 text-muted focus:text-ink focus:bold #mode", "◐ theme",
			func(e *gtui.Event) { app.ToggleTheme() },
		),
	}
}
