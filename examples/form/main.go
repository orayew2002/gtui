package main

import "github.com/orayew2002/gtui"

func main() {
	app := gtui.Init()

	app.HTML(
		gtui.Center("full bg-canvas",
			gtui.Column("w-50 rounded border-line bg-panel px-2 gap-0 key-f",
				gtui.Text("text-ink bold", "new movie"),
				gtui.Text("text-muted dim", "title"),
				gtui.Input("rounded border-line px-1 #title", "title", "Garagum gijeleri"),
				gtui.Text("text-muted dim", "storage"),
				gtui.Selector("#disk", "disk", "pick a disk", "disk-01", "disk-02", "disk-03"),
				gtui.Text("text-muted dim", "quality"),
				gtui.Range("#quality", "quality", 0, 100, 10, 70),
				gtui.Switch("#publish", "publish", "publish right away", true),
				gtui.Row("gap-1",
					gtui.Button("rounded border-success text-success focus:bold", "save", func(e *gtui.Event) {
						app.Toast("success", "saved "+app.Value("#title"))
					}),
					gtui.Button("rounded border-line text-muted focus:bold", "clear", func(e *gtui.Event) {
						app.SetValue("#title", "")
					}),
				),
			),
		),
	)

	app.Bind("ctrl+s", func() { app.Toast("", "ctrl+s saved") })
	app.Render()
}
