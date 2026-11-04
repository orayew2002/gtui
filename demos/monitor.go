// A monitoring wall: numbers, charts and a log that all move on their own.
package demos

import (
	"math"

	"github.com/orayew2002/gtui"
)

var (
	monitorTick = 0
)

func Monitor(a *gtui.App, tick int) gtui.Component {
	host = a
	monitorTick = tick

	return gtui.Column("full",
		gtui.Row("px-1 items-center gap-2 border-b border-line",
			gtui.IconText("text-accent bold nowrap", "chart", "monitor"),
			gtui.Marquee("text-muted truncate", "four regions · 18 402 requests a minute · all green", 40, monitorTick),
			gtui.Spacer(),
			gtui.Loader("text-accent nowrap", "live", monitorTick),
		),
		gtui.Div("grow clip scroll px-1 py-1 gap-1 col",
			gtui.Banner("", "warning", "disk 91% full", "the archive pool needs attention",
				gtui.Button("pill border-warning text-warning", "clean up", func(e *gtui.Event) {
					host.Toast("success", "4 GB freed")
				}),
			),
			gtui.Grid("gap-1", metricColumns(a),
				gtui.Metric("", "requests", "18.4k", "+8%"),
				gtui.Metric("", "errors", "12", "-3%"),
				gtui.Metric("", "p95", "41ms", "flat"),
				gtui.Metric("", "uptime", "99.98%", "+0.01"),
			),
			gtui.Row("gap-1",
				gtui.Card("grow rounded border-line bg-panel px-1", "traffic",
					gtui.Sparkline("text-success", monitorWave(32, 0)),
					gtui.LineChart("", monitorWave(32, monitorTick), 36, 6),
				),
				gtui.Card("w-30 shrink-0 rounded border-line bg-panel px-1", "by language",
					gtui.BarChart("", []gtui.Bar{
						{Label: "go", Value: 1284, Color: "text-accent"},
						{Label: "rust", Value: 894, Color: "text-success"},
						{Label: "zig", Value: 317, Color: "text-warning"},
					}, 14),
				),
			),
			gtui.Row("gap-1",
				gtui.Card("grow rounded border-line bg-panel px-1", "pools",
					gtui.Gauge("text-success", 8, 10, 18),
					gtui.Gauge("text-warning", 9.1, 10, 18),
					gtui.Progress("text-accent", float64(monitorTick%20), 20, 18),
					gtui.Row("gap-2 items-center pt-1",
						gtui.Donut("text-accent", 3, 4),
						gtui.Equalizer("text-success", 10, monitorTick),
						gtui.Dots("text-muted", monitorTick),
						gtui.Pulse("text-danger", 6, monitorTick),
						gtui.Wave("text-accent", 8, monitorTick),
					),
				),
				gtui.Card("w-34 shrink-0 rounded border-line bg-panel px-1", "rollout",
					gtui.Steps("", 1,
						gtui.StepItem{Title: "build", Note: "done"},
						gtui.StepItem{Title: "canary", Note: "running"},
						gtui.StepItem{Title: "all"},
					),
					gtui.Separator("", "history"),
					gtui.Timeline("",
						gtui.TimelineItem{Icon: "check", Kind: "ok", Title: "1.8.2 live", Time: "09:12"},
						gtui.TimelineItem{Icon: "refresh", Kind: "busy", Title: "1.8.3 canary", Time: "09:40"},
					),
				),
			),
			gtui.Card("rounded border-line bg-panel px-1", "log",
				gtui.Logs("", monitorLines(), monitorStamps()),
			),
		),
		gtui.StatusBar("border-t border-line px-1",
			[]gtui.Component{gtui.Row("gap-1", gtui.StatusDot("", "ok"), gtui.Text("nowrap", "all regions up"))},
			[]gtui.Component{gtui.Text("nowrap text-muted dim", "q quits")},
		),
	)
}

// metricColumns fits the four numbers across the screen, two at a time when it
// is narrow.
func metricColumns(a *gtui.App) int {
	if narrow(a) {
		return 2
	}
	return 4
}

// monitorWave is a stand-in for a real series.
func monitorWave(n, at int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = 20 + 12*math.Sin(float64(i+at)/4)
	}
	return out
}

func monitorLines() []string {
	return []string{
		"GET /disks 200 1.4ms",
		"warn: disk-07 answered slowly",
		"error: checksum failed on disk-07",
		"GET /disks 200 1.1ms",
	}
}

func monitorStamps() []string {
	return []string{"09:12:04", "09:12:09", "09:12:11", "09:12:18"}
}
