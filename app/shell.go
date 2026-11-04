package main

import "github.com/orayew2002/gtui"

func shell() gtui.Component {
	return gtui.Row("full bg-canvas p-1 gap-1",
		sidebar(),
		gtui.Column("grow gap-1",
			gtui.Div("row rounded border-line bg-panel px-2 gap-2 #head", header(current())...),
			gtui.Div("grow scroll rounded border-line key-c #main", current().Body()),
		),
	)
}

func sidebar() gtui.Component {
	items := []gtui.Component{
		gtui.Row("gap-1 px-1",
			gtui.Icon("text-accent bold", "terminal"),
			gtui.Text("text-ink bold nowrap", "gtui"),
			gtui.Text("grow text-right text-muted nowrap", "gallery"),
		),
		gtui.Divider("border-line"),
	}

	for _, item := range pages() {
		items = append(items, navLink(item))
	}

	items = append(items,
		gtui.Spacer(),
		gtui.Divider("border-line"),
		gtui.IconText("text-muted", "search", "ctrl+p"),
	)

	return gtui.Column("w-18 shrink-0 rounded border-line bg-panel px-1 key-n #nav", items...)
}

func navLink(item page) gtui.Component {
	class := "row px-1 nowrap text-muted" +
		" hover:text-ink" +
		" focus:bg-line focus:text-ink focus:bold" +
		" current:bg-line current:text-ink current:bold" +
		" key-" + item.Key + " key-hide" +
		" #nav" + item.Key

	return gtui.Href(class, item.Key+" "+item.Title, "#main", item.Body())
}

func header(item page) []gtui.Component {
	return []gtui.Component{
		gtui.Column("grow",
			gtui.Text("text-ink bold nowrap", item.Title),
			gtui.Text("text-muted nowrap", item.About),
		),
		gtui.Button("shrink-0 text-muted focus:text-ink focus:bold #mode", "◐ theme",
			func(e *gtui.Event) { app.ToggleTheme() },
		),
	}
}
