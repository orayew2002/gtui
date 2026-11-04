package main

import (
	"math"

	"github.com/orayew2002/gtui"
)

var series = sine(36)

func sine(n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = 20 + 12*math.Sin(float64(i)/4)
	}
	return out
}

func pageCharts() gtui.Component {
	return gtui.Column("grow gap-1",
		gtui.Row("gap-1 grow",
			card("a", "bar chart",
				gtui.BarChart("", []gtui.Bar{
					{Label: "go", Value: 1284, Color: "text-accent"},
					{Label: "rust", Value: 894, Color: "text-success"},
					{Label: "zig", Value: 317, Color: "text-warning"},
					{Label: "c", Value: 642, Color: "text-danger"},
				}, 16),
			),
			card("b", "sparkline and line",
				note("same data, two densities"),
				gtui.Sparkline("text-accent", series),
				gtui.LineChart("", series, 26, 5),
			),
			card("d", "gauges",
				gtui.Gauge("text-success", 8, 10, 16),
				gtui.Gauge("text-warning", 5, 10, 16),
				gtui.Gauge("text-danger", 2, 10, 16),
				gtui.Row("gap-1 items-center",
					gtui.Donut("text-accent", 1, 4),
					gtui.Donut("text-accent", 2, 4),
					gtui.Donut("text-accent", 3, 4),
					gtui.Donut("text-accent", 4, 4),
				),
			),
		),
		gtui.Row("gap-1",
			card("e", "stat cards",
				gtui.Row("gap-1",
					gtui.StatCard("grow", "requests", "18 402", "last hour"),
					gtui.StatCard("grow", "errors", "12", "0.07 percent"),
					gtui.StatCard("grow", "p95", "41 ms", "steady"),
				),
			),
		),
	)
}
