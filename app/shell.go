package main

import (
	"github.com/orayew2002/gtui"
	"github.com/orayew2002/gtui/src/colors"
)

func shell() gtui.Component {
	return gtui.Row("full bg-canvas p-1 gap-1",
		sidebar(),
		gtui.Column("grow gap-1",
			gtui.Div("row rounded border-line bg-panel px-2 gap-2 #head", header(active)...),
			gtui.Div("grow scroll key-c key-hide #main", active.Page()),
		),
	)
}

func sidebar() gtui.Component {
	items := []gtui.Component{
		gtui.Row("gap-1 px-1",
			gtui.Text("text-ink bold nowrap", "▣ "+gtui.T("app")),
			gtui.Text("grow text-right text-muted nowrap", gtui.T("admin")),
		),
		gtui.Row("rounded border-line px-1",
			gtui.Text("text-muted nowrap", "⌕ "+gtui.T("search")),
		),
	}

	for _, item := range routes {
		if item.Section != "" {
			items = append(items, gtui.Text("text-muted dim nowrap", gtui.T(item.Section)))
		}
		items = append(items, navLink(item))
	}

	items = append(items,
		gtui.Spacer(),
		gtui.Row("border-t border-line gap-1",
			gtui.Text("text-ink nowrap", "◍ "+gtui.T("admin")),
			gtui.Text("grow text-right text-muted nowrap", "ctrl+p"),
		),
	)

	return gtui.Column("w-20 shrink-0 rounded border-line bg-panel px-1 key-n key-hide #nav", items...)
}

func navLink(item route) gtui.Component {
	class := "row px-1 nowrap text-muted" +
		" hover:text-ink" +
		" focus:bg-line focus:text-ink focus:bold" +
		" current:bg-line current:text-ink current:bold" +
		" key-" + item.Key + " key-hide" +
		" #nav" + item.Key

	return gtui.Href(class, item.Key+" "+gtui.T(item.Label), "#main", item.Page())
}

func header(item route) []gtui.Component {
	body := []gtui.Component{
		gtui.Column("grow",
			gtui.Text("text-ink bold nowrap", gtui.T(item.Label)),
			gtui.Text("text-muted nowrap", description(item.Label)),
		),
		languageSwitch(),
		gtui.Button("shrink-0 text-muted focus:text-ink focus:bold #mode", "◐ "+gtui.T("theme"),
			func(e *gtui.Event) { app.ToggleTheme() },
		),
	}

	if item.Action != "" {
		label := gtui.T("action." + item.Action)
		if item.Action == "create" {
			label = gtui.T("action.create", gtui.T(item.Label))
		}

		body = append(body, gtui.Button("shrink-0 bg-ink text-canvas bold px-2 focus:underline #action", label,
			func(e *gtui.Event) { app.Toast("", label) },
		))
	}

	return body
}

func description(label string) string {
	key := "desc." + label
	if text := gtui.T(key); text != key {
		return text
	}
	return ""
}

func languageSwitch() gtui.Component {
	codes := gtui.Languages()
	items := make([]gtui.Component, 0, len(codes))

	for _, code := range codes {
		pick := code
		class := "px-1 nowrap text-muted focus:text-ink focus:bold #lang-" + code

		if code == gtui.Language() {
			class = "px-1 nowrap text-ink bold focus:underline #lang-" + code
		}

		items = append(items, gtui.Button(class, pick, func(e *gtui.Event) { app.Language(pick) }))
	}

	return gtui.Row("shrink-0 gap-1", items...)
}

func themeSwitch() gtui.Component {
	items := make([]gtui.Component, 0, 2)

	for _, name := range app.Themes() {
		pick := name
		mark := "○ "
		class := "row gap-1 text-muted focus:text-ink focus:bold #theme-" + name

		if name == app.ThemeName() {
			mark = "◉ "
			class = "row gap-1 text-ink bold focus:underline #theme-" + name
		}

		items = append(items, gtui.Button(class, mark+name, func(e *gtui.Event) {
			app.ThemeByName(pick)
			app.Refresh()
		}))
	}

	_ = colors.Dark
	return gtui.Column("", items...)
}
