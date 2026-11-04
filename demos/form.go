package demos

import "github.com/orayew2002/gtui"

// Form is one panel of inputs with a save, which is most of what a small tool
// ever needs.
func Form(a *gtui.App, tick int) gtui.Component {
	host = a

	return gtui.Center("full",
		gtui.Column("w-50 rounded border-line bg-panel px-2 gap-0",
			gtui.Text("text-ink bold", "new movie"),
			gtui.Text("text-muted dim", "title"),
			gtui.Input("rounded border-line px-1 #title", "title", "Garagum gijeleri"),
			gtui.Text("text-muted dim", "storage"),
			gtui.Selector("#disk", "disk", "pick a disk", "disk-01", "disk-02", "disk-03"),
			gtui.Text("text-muted dim", "quality"),
			gtui.Range("#quality", "quality", 0, 100, 10, 70),
			gtui.Switch("#publish", "publish", "publish right away", true),
			gtui.Row("gap-1 pt-1",
				gtui.Button("rounded border-success text-success focus:bold", "save", func(e *gtui.Event) {
					host.Toast("success", "saved "+host.Value("#title"))
				}),
				gtui.Button("rounded border-line text-muted focus:bold", "clear", func(e *gtui.Event) {
					host.SetValue("#title", "")
				}),
			),
		),
	)
}
