package demos

import "github.com/orayew2002/gtui"

// Hello is the smallest app: one centred card.
func Hello(a *gtui.App, tick int) gtui.Component {
	host = a

	return gtui.Center("full",
		gtui.Column("rounded border-accent bg-panel px-3 py-1 gap-0",
			gtui.Text("text-accent bold", "hello from gtui"),
			gtui.Text("text-muted", "press q to quit"),
		),
	)
}
