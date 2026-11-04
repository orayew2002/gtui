package main

import (
	"strings"

	"github.com/orayew2002/gtui"
)

func pageInputs() gtui.Component {
	// OnInput and OnChange attach by selector after the page is on screen,
	// which is the alternative to wiring a callback at build time
	defer app.Do(wireWatchers)

	return gtui.Column("grow gap-1",
		cards(
			card("a", "text fields",
				note("type, then try ← → home end backspace"),
				gtui.Field("gap-0", "name",
					gtui.Input("w-full px-1 border-b border-line focus:text-accent #in-text", "text", "plain text").
						On("input", func(e *gtui.Event) { say("in-out", "text is %q", e.Value) }),
				),
				gtui.Row("gap-x-2 gap-y-1 flex-wrap",
					gtui.Number("w-10 px-1 border-b border-line focus:text-accent #in-num", "number", "digits"),
					gtui.Password("w-12 px-1 border-b border-line focus:text-accent #in-pass", "password", "masked"),
				),
				out("in-out", "nothing typed"),
			),
			card("b", "text area and values",
				note("enter makes a line, ↑↓ move between lines"),
				gtui.TextArea("w-full rounded border-line px-1 focus:border-accent #in-area", "area", 3).
					On("input", func(e *gtui.Event) {
						say("area-out", "%d lines, %d runes", strings.Count(e.Value, "\n")+1, len([]rune(e.Value)))
					}),
				gtui.Row("gap-x-2 flex-wrap",
					gtui.Button("pill border-accent text-accent #read-values", "read values", func(e *gtui.Event) {
						values := app.Values()
						say("area-out", "text=%q area=%q", values["text"], values["area"])
					}),
					gtui.Button("pill border-line text-muted", "fill", func(e *gtui.Event) {
						app.SetValue("#in-text", "written from code")
						app.SetValue("#in-area", "line one\nline two")
					}),
					gtui.Button("pill border-line text-muted", "clear", func(e *gtui.Event) {
						app.SetValue("#in-text", "")
						app.SetValue("#in-area", "")
					}),
				),
				out("area-out", "1 line, 0 runes"),
			),
		),
		cards(
			card("d", "toggles",
				note("enter or space flips the focused one"),
				gtui.Switch("#in-switch", "switch", "notifications", true).
					On("change", func(e *gtui.Event) { say("tog-out", "switch=%s", e.Value) }),
				gtui.Checkbox("#in-check", "checkbox", "remember me", false).
					On("change", func(e *gtui.Event) { say("tog-out", "checkbox=%s", e.Value) }),
				out("tog-out", "—"),
			),
			card("f", "search in a group",
				note("a prefix, the field and a key hint in one border"),
				gtui.InputGroup("w-full",
					gtui.Icon("", "search"),
					gtui.Search("w-full #in-find", "find", "search the gallery").
						On("input", func(e *gtui.Event) { say("find-out", "looking for %q", e.Value) }),
					gtui.Kbd("", "/"),
				),
				out("find-out", "nothing typed"),
			),
		),
		gtui.Separator("", "one field per type"),
		gtui.Panel("", "the filtered types", nil,
			note("each one refuses the runes it does not want"),
			gtui.Row("gap-x-2 gap-y-1 flex-wrap",
				gtui.Field("gap-0", "email",
					gtui.Email("w-18 px-1 border-b border-line focus:text-accent", "email", "name@site.tm"),
				),
				gtui.Field("gap-0", "url",
					gtui.URL("w-18 px-1 border-b border-line focus:text-accent", "url", "https://gtui.dev"),
				),
				gtui.Field("gap-0", "phone",
					gtui.Phone("w-16 px-1 border-b border-line focus:text-accent", "phone", "+993 6x xx"),
				),
				gtui.Field("gap-0", "date",
					gtui.Date("w-12 px-1 border-b border-line focus:text-accent", "date"),
				),
			),
		),
		card("j", "watching a field from the app",
			note("app.OnInput fires on every keystroke, app.OnChange when it settles"),
			gtui.Field("gap-0", "watched",
				gtui.Input("w-full px-1 border-b border-line focus:text-accent #watched", "watched", "type here"),
			),
			out("watch-out", "no keystroke seen yet"),
			gtui.Row("gap-x-2 flex-wrap mt-1",
				gtui.Button("pill border-accent text-accent", "set from code", func(e *gtui.Event) {
					app.SetValue("#watched", "written by SetValue")
				}),
				gtui.Button("pill border-line text-muted", "toggle a class", func(e *gtui.Event) {
					bold = !bold
					app.Class("#watch-out", "bold", bold)
					say("watch-out", "app.Class turned bold %v", bold)
				}),
			),
		),
	)
}

var (
	bold  = false
	wired = false
)

// wireWatchers attaches the selector based handlers once, after the first build
// has put the field on screen.
func wireWatchers() {
	if wired || app.Query("#watched") == nil {
		return
	}
	wired = true

	app.OnInput("#watched", func(e *gtui.Event) { say("watch-out", "OnInput saw %q", e.Value) })
	app.OnChange("#watched", func(e *gtui.Event) { say("watch-out", "OnChange settled on %q", e.Value) })
}
