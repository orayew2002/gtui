package main

import (
	"github.com/orayew2002/gtui"
	"github.com/orayew2002/gtui/src/colors"
)

func pageThemes() gtui.Component {
	tokens := []struct{ name, class string }{
		{"canvas", "bg-canvas"},
		{"panel", "bg-panel"},
		{"raised", "bg-raised"},
		{"line", "bg-line"},
		{"muted", "bg-muted"},
		{"ink", "bg-ink"},
		{"accent", "bg-accent"},
		{"success", "bg-success"},
		{"warning", "bg-warning"},
		{"danger", "bg-danger"},
		{"info", "bg-info"},
	}

	swatches := make([]gtui.Component, 0, len(tokens))
	for _, token := range tokens {
		swatches = append(swatches, gtui.Row("gap-1 items-center",
			gtui.Div("w-8 h-1 "+token.class),
			gtui.Text("nowrap text-muted", token.name),
		))
	}

	buttons := make([]gtui.Component, 0, 3)
	for _, name := range app.Themes() {
		pick := name
		buttons = append(buttons, gtui.Button("rounded border-line focus:bold", name, func(e *gtui.Event) {
			app.ThemeByName(pick)
			app.Refresh()
		}))
	}

	buttons = append(buttons, gtui.Button("rounded border-accent text-accent focus:bold", "custom", func(e *gtui.Event) {
		app.Tokens(func(t *colors.Theme) {
			t.Name = "ocean"
			t.Canvas = colors.RGB(8, 18, 28)
			t.Panel = colors.RGB(14, 28, 42)
			t.Raised = colors.RGB(20, 38, 56)
			t.Line = colors.RGB(30, 54, 76)
			t.Accent = colors.RGB(120, 200, 255)
		})
		app.Refresh()
	}))

	return gtui.Column("grow gap-1",
		gtui.Row("gap-1 grow",
			card("a", "tokens",
				note("every colour class reads from the active theme"),
				gtui.Column("gap-0", swatches...),
			),
			card("b", "switch",
				note("switching costs one frame, no rebuild"),
				gtui.Row("gap-1 flex-wrap", buttons...),
				gtui.Text("text-muted nowrap", "active: %s", app.ThemeName()),
				gtui.Hint("", "◐ theme", "the header button toggles dark and light"),
			),
			card("d", "on a panel",
				gtui.Text("text-ink bold nowrap", "ink on panel"),
				gtui.Text("text-muted nowrap", "muted on panel"),
				gtui.Text("text-accent nowrap", "accent on panel"),
				gtui.Text("text-success nowrap", "success"),
				gtui.Text("text-warning nowrap", "warning"),
				gtui.Text("text-danger nowrap", "danger"),
			),
		),
	)
}
