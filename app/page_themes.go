package main

import (
	"github.com/orayew2002/gtui"
	"github.com/orayew2002/gtui/src/colors"
)

var themeTokens = []string{
	"canvas", "panel", "raised", "line", "muted", "ink",
	"accent", "success", "warning", "danger", "info",
}

// shownTheme is the theme the page was last built with. The header button and
// the palette only redraw, so the tick compares the two and rebuilds the page
// when something else swapped the palette underneath it.
var shownTheme = ""

var (
	caretOn = true
	mouseOn = true
)

func pageThemes() gtui.Component {
	shownTheme = app.ThemeName()

	return gtui.Row("gap-1 grow",
		card("a", "tokens", swatches()...),
		gtui.Column("grow gap-1",
			card("b", "switch",
				gtui.Row("gap-1 flex-wrap", themeButtons()...),
				gtui.Column("gap-0 #theme-now", themeNow()...),
			),
			card("d", "the same card, either theme",
				gtui.Row("gap-1 items-center",
					gtui.Icon("text-accent", "server"),
					gtui.Text("text-ink bold nowrap", "api-01"),
					gtui.Spacer(),
					gtui.Text("pill border-success text-success nowrap", "live"),
				),
				gtui.Progress("text-accent", 7, 10, 12),
				note("one set of classes, no palette in sight"),
				gtui.Button("px-1 focus:bold focus:text-accent #th-flip", "flip", func(e *gtui.Event) {
					app.ToggleTheme()
					show(current())
				}),
			),
			card("f", "terminal switches",
				note("the caret, the mouse and the background can go back to the terminal"),
				gtui.Row("gap-x-2 flex-wrap",
					gtui.Button("pill border-accent text-accent #th-bare", "see through", func(e *gtui.Event) {
						seeThrough = !seeThrough
						app.Transparent(seeThrough)
						show(current())
						say("th-switches", "Transparent(%v)", seeThrough)
					}),
					gtui.Button("pill border-accent text-accent #th-caret", "caret", func(e *gtui.Event) {
						caretOn = !caretOn
						app.UseCursor(caretOn)
						say("th-switches", "UseCursor(%v)", caretOn)
					}),
					gtui.Button("pill border-accent text-accent #th-mouse", "mouse", func(e *gtui.Event) {
						mouseOn = !mouseOn
						app.UseMouse(mouseOn)
						say("th-switches", "UseMouse(%v)", mouseOn)
					}),
				),
				out("th-switches", "caret on, mouse on"),
			),
		),
	)
}

func swatches() []gtui.Component {
	rows := make([]gtui.Component, 0, len(themeTokens))

	for _, name := range themeTokens {
		cells := []gtui.Component{
			gtui.Text("shrink-0 nowrap text-"+name, "██"),
			gtui.Text("shrink-0 nowrap bg-"+name, "  "),
			gtui.Text("w-8 truncate text-muted", name),
		}
		if wide() {
			if colour, ok := colors.Token(name); ok {
				cells = append(cells, gtui.Text("dim nowrap", colour.Hex()))
			}
		}

		rows = append(rows, gtui.Row("gap-1 items-center", cells...))
	}

	return rows
}

func themeButtons() []gtui.Component {
	names := app.Themes()
	buttons := make([]gtui.Component, 0, len(names)+1)

	for _, pick := range names {
		buttons = append(buttons, gtui.Button("px-1 focus:bold focus:text-accent #th-"+pick, pick, func(e *gtui.Event) {
			app.ThemeByName(pick)
			show(current())
		}))
	}

	return append(buttons, gtui.Button("px-1 text-accent focus:bold #th-new", "ocean", func(e *gtui.Event) {
		app.Tokens(func(t *colors.Theme) {
			t.Name = "ocean"
			t.Canvas = colors.RGB(8, 18, 28)
			t.Panel = colors.RGB(14, 28, 42)
			t.Raised = colors.RGB(20, 38, 56)
			t.Line = colors.RGB(30, 54, 76)
			t.Accent = colors.RGB(120, 200, 255)
		})
		show(current())
	}))
}

func themeNow() []gtui.Component {
	return []gtui.Component{
		gtui.Text("text-muted nowrap", "active: %s", app.ThemeName()),
		gtui.Hint("", "◐ theme", "the header toggles too"),
	}
}

// liveThemeName rebuilds the page when the palette changed without it, which is
// what the header button does: Theme only repaints.
func liveThemeName() {
	if app.ThemeName() == shownTheme || app.Query("#theme-now") == nil {
		return
	}
	show(current())
}
