package gtui

import (
	"strings"

	"github.com/orayew2002/gtui/src/colors"
	"github.com/orayew2002/gtui/src/dom"
)

type Bar struct {
	Label string
	Value float64
	Color string
}

func Progress(class string, value, max float64, width int) Component {
	if width < 1 {
		width = 10
	}

	filled := scale(value, max, width)

	return Text(join("nowrap", class),
		strings.Repeat("█", filled)+strings.Repeat("░", width-filled))
}

func Gauge(class string, value, max float64, width int) Component {
	if width < 4 {
		width = 10
	}

	percent := 0
	if max > 0 {
		percent = int(value / max * 100)
	}

	label := itoa(clampInt(percent, 0, 100)) + "%"
	track := width - len(label) - 1
	if track < 1 {
		track = 1
	}

	return Div(join("row gap-1", class),
		Progress("", value, max, track),
		Text("nowrap", label),
	)
}

func BarChart(class string, bars []Bar, width int) Component {
	if width < 1 {
		width = 20
	}

	max := 0.0
	label := 0
	for _, bar := range bars {
		if bar.Value > max {
			max = bar.Value
		}
		if len(bar.Label) > label {
			label = len(bar.Label)
		}
	}

	rows := make([]Component, 0, len(bars))
	for _, bar := range bars {
		ink := bar.Color
		if ink == "" {
			ink = "text-sky"
		}

		rows = append(rows, Div("row gap-1",
			Text("w-"+itoa(label)+" nowrap dim", bar.Label),
			Text("nowrap "+ink, strings.Repeat("█", scale(bar.Value, max, width))),
		))
	}

	return Div(join("col gap-0", class), rows...)
}

func Sparkline(class string, values []float64) Component {
	glyphs := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	min, max := bounds(values)
	span := max - min

	var sb strings.Builder
	for _, value := range values {
		level := 0
		if span > 0 {
			level = int((value - min) / span * float64(len(glyphs)-1))
		}
		sb.WriteRune(glyphs[clampInt(level, 0, len(glyphs)-1)])
	}

	return Text(join("nowrap", class), sb.String())
}

func LineChart(class string, values []float64, width, height int) Component {
	if width < 1 || height < 1 || len(values) == 0 {
		return Div(class)
	}

	canvas := dom.NewCanvas(width, height)
	min, max := bounds(values)
	span := max - min

	for x := range width {
		index := x * len(values) / width
		value := values[index]

		level := 0
		if span > 0 {
			level = int((value - min) / span * float64(height*2-1))
		}
		level = clampInt(level, 0, height*2-1)

		row := height - 1 - level/2
		glyph := '▄'
		if level%2 == 1 {
			glyph = '▀'
		}

		canvas.Set(x, row, dom.Pixel{Ch: glyph, Fg: colors.Sky, HasFg: true})
	}

	return canvasComponent(class, canvas)
}

func Donut(class string, value, max float64) Component {
	glyphs := []string{"○", "◔", "◑", "◕", "●"}
	step := scale(value, max, len(glyphs)-1)

	return Text(join("nowrap", class), glyphs[clampInt(step, 0, len(glyphs)-1)])
}

func canvasComponent(class string, canvas *dom.Canvas, child ...Component) Component {
	return Component{
		class: class,
		child: child,
		build: func(a *App) *Node {
			n := dom.NewElement("div")

			base, id := splitClass(class)
			if base != "" {
				n.Set("class", base)
			}
			if id != "" {
				n.Set("id", id)
			}

			n.Canvas = canvas

			for _, item := range child {
				n.Append(item.node(a))
			}

			return n
		},
	}
}

func scale(value, max float64, width int) int {
	if max <= 0 || value <= 0 {
		return 0
	}

	filled := int(value / max * float64(width))
	if filled == 0 {
		filled = 1
	}
	return clampInt(filled, 0, width)
}

func bounds(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}

	min, max := values[0], values[0]
	for _, value := range values {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
