package main

import "github.com/orayew2002/gtui"

// frame counts animation ticks for every page that wants movement.
var frame = 0

// onTick runs eight times a second and only repaints pages that move.
func onTick(tick int) {
	frame = tick

	if playing && tick%2 == 0 {
		beat = (beat + 1) % reel.Count()
	}

	switch current().Title {
	case "motion", "media", "feedback":
		app.Open("#main", current().Body())
	}
}

func pageMotion() gtui.Component {
	return gtui.Column("grow gap-1",
		gtui.Row("gap-1",
			card("a", "text in motion",
				note("Marquee, Typewriter and Blink take a tick number"),
				gtui.Marquee("text-accent", "gtui writes only the cells that changed, so motion is cheap", 40, frame),
				gtui.Typewriter("text-success", "typed one rune at a time", frame/2),
				gtui.Blink("text-danger", "● recording", frame/4),
			),
			card("b", "meters",
				note("Pulse, Wave, Dots, Spinner and Equalizer"),
				gtui.Pulse("text-warning", 18, frame),
				gtui.Wave("text-accent", 18, frame),
				gtui.Equalizer("text-success", 18, frame),
				gtui.Row("gap-2 items-center",
					gtui.Spinner("text-accent", frame),
					gtui.Dots("text-muted", frame),
					gtui.Text("text-muted nowrap", "tick %d", frame),
				),
			),
		),
		gtui.Row("gap-1",
			card("d", "how it works",
				note("the whole page above is three lines of code"),
				gtui.CodeView(gtui.CodeOptions{
					Lines: []string{
						`app.Animate(8, func(tick int) {`,
						`    app.Open("#main", page(tick))`,
						`})`,
					},
					First: 1,
				}),
				gtui.Hint("", "app.Every", "for slower timers, app.After for one shot"),
			),
		),
	)
}
