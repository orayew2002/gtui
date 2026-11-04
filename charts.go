package gtui

import (
	"math"
	"strings"

	"github.com/orayew2002/gtui/src/colors"
	"github.com/orayew2002/gtui/src/dom"
	"github.com/orayew2002/gtui/src/render"
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
	// one track cell, a gap and the widest label ("100%") is the true minimum
	if width < 6 {
		width = 10
	}

	percent := int(ratio(value, max) * 100)
	label := itoa(clampInt(percent, 0, 100)) + "%"

	return Div(join("row gap-1", class),
		Progress("", value, max, width-5),
		Text("w-4 text-right nowrap", label),
	)
}

func BarChart(class string, bars []Bar, width int) Component {
	if width < 1 {
		width = 20
	}

	max := 0.0
	label := 0
	for _, bar := range bars {
		if bar.Value > max && !math.IsInf(bar.Value, 1) {
			max = bar.Value
		}
		if w := render.StringWidth(bar.Label); w > label {
			label = w
		}
	}

	rows := make([]Component, 0, len(bars))
	for _, bar := range bars {
		ink := bar.Color
		if ink == "" {
			ink = "text-accent"
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
		level := int(norm(value, min, span) * float64(len(glyphs)-1))
		sb.WriteRune(glyphs[clampInt(level, 0, len(glyphs)-1)])
	}

	return Text(join("nowrap", class), sb.String())
}

func LineChart(class string, values []float64, width, height int) Component {
	if width < 1 || height < 1 || len(values) == 0 {
		return Div(class)
	}

	ink, ok := inkFromClass(class)
	if !ok {
		ink = colors.Active().Accent
	}

	canvas := dom.NewCanvas(width, height)
	min, max := bounds(values)
	span := max - min

	for x := range width {
		index := x * len(values) / width
		value := values[index]

		level := clampInt(int(norm(value, min, span)*float64(height*2-1)), 0, height*2-1)

		row := height - 1 - level/2
		glyph := '▄'
		if level%2 == 1 {
			glyph = '▀'
		}

		canvas.Set(x, row, dom.Pixel{Ch: glyph, Fg: ink, HasFg: true})
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

// inkFromClass reads the last text-* colour out of a class string, so a canvas
// can paint itself in the colour the caller asked for.
func inkFromClass(class string) (colors.Color, bool) {
	var (
		out   colors.Color
		found bool
	)

	for _, token := range strings.Fields(class) {
		value, ok := strings.CutPrefix(strings.TrimPrefix(token, "box-"), "text-")
		if !ok {
			continue
		}
		if inner, ok := strings.CutPrefix(value, "["); ok {
			value = strings.TrimSuffix(inner, "]")
		}
		if c, ok := colors.Parse(value); ok {
			out, found = c, true
		}
	}

	return out, found
}

func scale(value, max float64, width int) int {
	if width < 1 {
		return 0
	}

	r := ratio(value, max)
	if r <= 0 {
		return 0
	}

	filled := int(r * float64(width))
	if filled < 1 {
		filled = 1
	}
	return clampInt(filled, 0, width)
}

// ratio clamps value/max into 0..1 before any int conversion, so NaN, infinity
// and huge ratios cannot overflow into a negative number of cells.
func ratio(value, max float64) float64 {
	if !(max > 0) || !(value > 0) {
		return 0
	}
	if value >= max {
		return 1
	}
	return value / max
}

// norm places value inside min..min+span as a 0..1 ratio, clamped in float
// space for the same reason as ratio.
func norm(value, min, span float64) float64 {
	if !(span > 0) {
		return 0
	}

	r := (value - min) / span
	if !(r > 0) {
		return 0
	}
	if r > 1 {
		return 1
	}
	return r
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
