package main

import "github.com/orayew2002/gtui"

var (
	app      *gtui.App
	position = 65
	playing  = true
	volume   = 70
)

var queue = []string{"Garagum gijeleri", "Ýaz ýagşy", "Altyn asyr"}

func main() {
	app = gtui.Init()
	app.View(view)

	app.Animate(2, func(tick int) {
		if playing && position < 300 {
			position++
			app.Refresh()
		}
	})

	app.Bind("space", func() { playing = !playing })
	app.Bind("left", func() { position = max(0, position-10) })
	app.Bind("right", func() { position = min(300, position+10) })

	app.Render()
}

func view() gtui.Component {
	return gtui.Center("full bg-canvas",
		gtui.Column("w-60 gap-1",
			gtui.Player(gtui.PlayerOptions{
				Class:    "key-p",
				Title:    queue[0],
				Subtitle: "tk · 4k · 2.1 GB",
				Position: position,
				Duration: 300,
				Playing:  playing,
				Volume:   volume,
				Width:    44,
				Queue:    queue,
			}),
			gtui.Row("gap-2 px-1",
				gtui.Hint("", "space", "play"),
				gtui.Hint("", "← →", "seek"),
				gtui.Spacer(),
				gtui.Hint("", "q", "quit"),
			),
		),
	)
}
