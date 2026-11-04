package main

import "github.com/orayew2002/gtui"

func pageFeedback() gtui.Component {
	return gtui.Column("grow gap-1",
		gtui.Row("gap-1",
			gtui.Alert("grow border-info", "info", "info alert", "this is what an alert looks like"),
			gtui.Alert("grow border-success", "success", "success alert", "something finished"),
		),
		gtui.Row("gap-1",
			gtui.Alert("grow border-warning", "warning", "warning alert", "something needs attention"),
			gtui.Alert("grow border-danger", "danger", "danger alert", "something failed"),
		),
		gtui.Row("gap-1 grow",
			card("a", "toasts",
				note("they disappear on their own"),
				gtui.Row("gap-1 flex-wrap",
					gtui.Button("rounded border-success text-success focus:bold", "success", func(e *gtui.Event) {
						app.Toast("success", "saved")
					}),
					gtui.Button("rounded border-warning text-warning focus:bold", "warning", func(e *gtui.Event) {
						app.Toast("warning", "disk almost full")
					}),
					gtui.Button("rounded border-danger text-danger focus:bold", "danger", func(e *gtui.Event) {
						app.Toast("danger", "upload failed")
					}),
				),
			),
			card("b", "modal",
				note("esc closes it, tab moves inside it"),
				gtui.Button("rounded border-accent text-accent focus:bold", "open modal", func(e *gtui.Event) {
					app.Modal("confirm",
						gtui.Text("text-muted", "modals take focus and give it back"),
						gtui.Row("gap-1",
							gtui.Button("rounded border-success text-success focus:bold #ok", "ok", func(e *gtui.Event) {
								app.CloseModal()
								app.Toast("success", "confirmed")
							}),
							gtui.Button("rounded border-line text-muted focus:bold #cancel", "cancel", func(e *gtui.Event) {
								app.CloseModal()
							}),
						),
					)
				}),
			),
			card("d", "progress",
				note("progress, gauge, spinner and skeleton"),
				gtui.Progress("text-accent", 7, 10, 18),
				gtui.Gauge("text-success", 3, 4, 18),
				gtui.Row("gap-1 items-center",
					gtui.Spinner("text-accent", frame),
					gtui.Dots("text-muted", frame),
					gtui.Text("text-muted nowrap", "working"),
				),
				gtui.Skeleton("", 2, 16),
			),
		),
	)
}
