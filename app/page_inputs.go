package main

import "github.com/orayew2002/gtui"

func pageInputs() gtui.Component {
	return gtui.Column("grow gap-1",
		gtui.Row("gap-1",
			card("a", "text fields",
				note("type, then use ← → home end delete"),
				gtui.Input("rounded border-line px-1 #in-text", "text", "plain text"),
				gtui.Number("rounded border-line px-1 w-12 #in-num", "number", "digits only"),
				gtui.Password("rounded border-line px-1 #in-pass", "password", "masked"),
				gtui.Email("rounded border-line px-1 #in-mail", "email", "name@site.tm"),
			),
			card("b", "text area",
				note("enter makes a new line"),
				gtui.TextArea("rounded border-line px-1 #in-area", "area", 4),
				gtui.Button("rounded border-accent text-accent focus:bold #read-values", "read values", func(e *gtui.Event) {
					values := app.Values()
					say("in-out", "text=%q area=%q", values["text"], values["area"])
				}),
				out("in-out", "click read values"),
			),
		),
		gtui.Row("gap-1",
			card("d", "toggles",
				note("enter or space flips them"),
				gtui.Switch("#in-switch", "switch", "switch", true),
				gtui.Checkbox("#in-check", "checkbox", "checkbox", false),
				gtui.Button("rounded border-line text-muted focus:bold", "read", func(e *gtui.Event) {
					values := app.Values()
					say("tog-out", "switch=%s checkbox=%s", values["switch"], values["checkbox"])
				}),
				out("tog-out", "—"),
			),
			card("e", "search",
				note("typing fires the input event"),
				gtui.Search("rounded border-line px-1 #in-find", "find", "search"),
				out("find-out", "nothing typed"),
			),
			card("f", "set values from code",
				note("buttons write into the fields above"),
				gtui.Row("gap-1",
					gtui.Button("rounded border-line focus:bold", "fill", func(e *gtui.Event) {
						app.SetValue("#in-text", "written from code")
						app.SetValue("#in-area", "line one\nline two")
					}),
					gtui.Button("rounded border-line focus:bold", "clear", func(e *gtui.Event) {
						app.SetValue("#in-text", "")
						app.SetValue("#in-area", "")
					}),
				),
			),
		),
	)
}
