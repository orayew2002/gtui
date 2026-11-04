package main

import "github.com/orayew2002/gtui"

func pageForm() gtui.Component {
	return gtui.Row("gap-1 grow",
		gtui.Column("grow rounded border-line bg-panel px-2 gap-0 key-f",
			gtui.Text("text-ink bold nowrap", gtui.T("form.title")),
			gtui.Div("#form-alert"),

			field(gtui.T("col.name"), gtui.Input("rounded border-line px-1 #f-title", "title", "Garagum gijeleri")),
			gtui.Row("gap-1",
				field(gtui.T("col.year"), gtui.Number("rounded border-line px-1 w-10 #f-year", "year", "2026")),
				field(gtui.T("col.duration"), gtui.Input("rounded border-line px-1 w-12 #f-time", "duration", "1:52:10")),
			),
			field(gtui.T("col.storage"), gtui.Selector("#f-disk", "storage", gtui.T("form.pick"), storages...)),
			field(gtui.T("form.quality"), gtui.Range("#f-quality", "quality", 0, 100, 10, 70)),
			field(gtui.T("form.about"), gtui.TextArea("rounded border-line px-1 #f-about", "about", 3)),

			gtui.Divider("border-line"),
			formButtons(),
		),

		gtui.Column("w-1/3 rounded border-line bg-panel px-2 gap-0 key-k",
			gtui.Text("text-ink bold nowrap", gtui.T("form.options")),
			gtui.Text("text-muted dim nowrap", gtui.T("col.tracks")),
			gtui.Radios("", "lang", "turkmen", "russian", "english"),
			gtui.Divider("border-line"),
			gtui.Switch("#f-publish", "publish", gtui.T("form.publish"), true),
			gtui.Switch("#f-notify", "notify", gtui.T("form.notify"), false),
			gtui.Divider("border-line"),
			gtui.Text("text-muted dim nowrap", gtui.T("form.tags")),
			gtui.Chips("", "drama", "tk", "4k"),
		),
	)
}

func formButtons() gtui.Component {
	return gtui.Row("gap-1",
		gtui.Button("rounded border-success text-success focus:bold #f-save", gtui.T("action.save"), saveForm),
		gtui.Button("rounded border-line text-muted focus:bold #f-reset", gtui.T("action.reset"), resetForm),
		gtui.Spacer(),
		gtui.Button("rounded border-danger text-danger focus:bold #f-delete", gtui.T("action.delete"), confirmDelete),
	)
}

func saveForm(e *gtui.Event) {
	values := app.Values()

	if values["title"] == "" {
		app.Open("#form-alert", gtui.Alert("border-danger", "danger",
			gtui.T("form.invalid"), gtui.T("form.needTitle")))

		if node := app.Query("#f-title"); node != nil {
			app.Focus(node)
		}
		return
	}

	app.Open("#form-alert", gtui.Alert("border-success", "success",
		gtui.T("form.saved"), values["title"]+" · "+values["storage"]))

	app.Toast("success", gtui.T("form.saved"))
}

func resetForm(e *gtui.Event) {
	for _, id := range []string{"#f-title", "#f-year", "#f-time", "#f-about"} {
		app.SetValue(id, "")
	}

	app.Open("#form-alert")
	app.Toast("", gtui.T("action.reset"))
}

func confirmDelete(e *gtui.Event) {
	app.Modal(gtui.T("form.deleteTitle"),
		gtui.Text("text-muted", gtui.T("form.deleteBody")),
		gtui.Row("gap-1",
			gtui.Button("rounded border-danger text-danger focus:bold #yes", gtui.T("action.delete"), func(e *gtui.Event) {
				app.CloseModal()
				app.Toast("danger", gtui.T("form.deleted"))
			}),
			gtui.Button("rounded border-line text-muted focus:bold #no", gtui.T("action.cancel"), func(e *gtui.Event) {
				app.CloseModal()
			}),
		),
	)
}
