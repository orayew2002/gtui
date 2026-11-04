package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/orayew2002/gtui"
)

var (
	frame   int
	traffic = make([]float64, 36)
	queue   = []job{
		{"garagum gijeleri", 12},
		{"ýaz ýagşy", 48},
		{"altyn asyr", 86},
	}
)

type job struct {
	Name    string
	Percent int
}

func startTicker() {
	for i := range traffic {
		traffic[i] = 20 + 12*math.Sin(float64(i)/4)
	}

	app.Animate(8, func(tick int) {
		frame = tick

		if tick%4 == 0 {
			for i := range queue {
				queue[i].Percent = (queue[i].Percent + 1 + rand.Intn(3)) % 101
			}

			copy(traffic, traffic[1:])
			traffic[len(traffic)-1] = 10 + rand.Float64()*30
		}

		if active().Label == "dynamic" {
			app.Open("#main", pageDynamic())
		}
	})
}

func pageDynamic() gtui.Component {
	jobs := make([]gtui.Component, 0, len(queue))

	for _, item := range queue {
		jobs = append(jobs, gtui.Row("gap-1 items-center",
			gtui.Text("w-18 nowrap text-ink", item.Name),
			gtui.Progress("text-accent", float64(item.Percent), 100, 18),
			gtui.Text("w-5 text-right nowrap text-muted", "%d%%", item.Percent),
		))
	}

	return gtui.Column("grow gap-1",
		gtui.Row("gap-1",
			gtui.Column("grow rounded border-line bg-panel px-2 gap-0 key-d",
				gtui.Row("gap-1",
					gtui.Text("text-ink bold nowrap", gtui.T("live.queue")),
					gtui.Spacer(),
					gtui.Spinner("text-accent", frame),
					gtui.Text("text-muted nowrap", time.Now().Format("15:04:05")),
				),
				gtui.Column("gap-0", jobs...),
			),
			gtui.Column("w-1/3 rounded border-line bg-panel px-2 gap-0 key-b",
				gtui.Text("text-ink bold nowrap", gtui.T("live.traffic")),
				gtui.Sparkline("text-success", traffic),
				gtui.LineChart("", traffic, 22, 4),
				gtui.Wave("text-accent", 20, frame),
			),
		),
		gtui.Row("gap-1",
			gtui.Column("grow rounded border-line bg-panel px-2 gap-0 key-w",
				gtui.Text("text-ink bold nowrap", gtui.T("live.motion")),
				gtui.Marquee("text-accent", gtui.T("live.note"), 42, frame),
				gtui.Typewriter("text-success", gtui.T("live.typed"), frame/2),
				gtui.Row("gap-2",
					gtui.Pulse("text-warning", 14, frame),
					gtui.Dots("text-muted", frame),
					gtui.Blink("text-danger", "● rec", frame/4),
				),
			),
			gtui.Column("w-1/3 rounded border-line bg-panel px-2 gap-0 key-x",
				gtui.Text("text-ink bold nowrap", gtui.T("live.controls")),
				gtui.Row("gap-1",
					gtui.Button("rounded border-accent text-accent focus:bold #toast", "toast", func(e *gtui.Event) {
						app.Toast("success", gtui.T("live.done"))
					}),
					gtui.Button("rounded border-line text-muted focus:bold #shuffle", "shuffle", func(e *gtui.Event) {
						for i := range queue {
							queue[i].Percent = rand.Intn(101)
						}
						app.Open("#main", pageDynamic())
					}),
				),
				gtui.Range("#speed", "speed", 1, 16, 1, 8),
				gtui.Switch("#live", "live", gtui.T("live.live"), true),
			),
		),
	)
}
