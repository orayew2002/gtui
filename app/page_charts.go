// Drawings made of cells: bars, a sparkline, a line chart, gauges and donuts.
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

// pageCharts sizes its drawings from the room the page has, so a chart never
// pushes its card past the edge of the terminal.
func pageCharts() gtui.Component {
	bars, deep := chartSize()

	return gtui.Column("grow gap-1",
		gtui.Row("gap-1",
			card("a", "bar chart",
				gtui.BarChart("", []gtui.Bar{
					{Label: "go", Value: 1284, Color: "text-accent"},
					{Label: "rust", Value: 894, Color: "text-success"},
					{Label: "zig", Value: 317, Color: "text-warning"},
					{Label: "c", Value: 642, Color: "text-danger"},
				}, bars),
			),
			card("b", "sparkline and line",
				gtui.Sparkline("text-accent", series),
				gtui.LineChart("", series, bars*2, deep),
			),
			card("d", "gauges",
				gtui.Gauge("text-success", 8, 10, bars),
				gtui.Gauge("text-warning", 5, 10, bars),
				gtui.Gauge("text-danger", 2, 10, bars),
				gtui.Row("gap-1 items-center",
					gtui.Donut("text-accent", 1, 4),
					gtui.Donut("text-accent", 2, 4),
					gtui.Donut("text-accent", 3, 4),
					gtui.Donut("text-accent", 4, 4),
				),
			),
		),
		gtui.Row("gap-1",
			gtui.StatCard("grow", "requests", "18 402", "last hour"),
			gtui.StatCard("grow", "errors", "12", "0.07 percent"),
			gtui.StatCard("grow", "p95", "41 ms", "steady"),
		),
	)
}

// chartSize keeps the drawings inside their cards: three cards share the width,
// and the line chart takes the rows the stat cards leave over.
func chartSize() (bars, deep int) {
	width, height := room()

	bars = min(max(width/3-10, 8), 26)
	deep = min(max(height-12, 3), 10)

	return bars, deep
}
