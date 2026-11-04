package main

import "github.com/orayew2002/gtui"

func main() {
	gtui.Init().HTML(
		gtui.Center("full bg-canvas",
			gtui.Column("rounded border-accent bg-panel px-3 py-1 gap-0",
				gtui.Text("text-accent bold", "hello from gtui"),
				gtui.Text("text-muted", "press q to quit"),
			),
		),
	).Render()
}
