package gtui

import (
	"strconv"
	"strings"

	"github.com/orayew2002/gtui/src/dom"
)

type CodeOptions struct {
	Class    string
	Lines    []string
	First    int
	Active   int
	Marks    map[int]string
	NoGutter bool
	Wrap     bool
}

func CodeView(opts CodeOptions) Component {
	width := len(strconv.Itoa(opts.First + len(opts.Lines)))
	rows := make([]Component, 0, len(opts.Lines))

	for i, line := range opts.Lines {
		number := opts.First + i

		cells := make([]Component, 0, 3)

		if !opts.NoGutter {
			ink := "text-muted dim nowrap"
			if number == opts.Active {
				ink = "text-accent bold nowrap"
			}
			cells = append(cells, Text("w-"+itoa(width)+" text-right "+ink, strconv.Itoa(number)))

			mark := opts.Marks[number]
			if mark == "" {
				mark = " "
			}
			cells = append(cells, Text("w-1 nowrap text-warning", mark))
		}

		body := "nowrap"
		if opts.Wrap {
			body = "wrap"
		}
		if number == opts.Active {
			body += " bg-raised"
		}

		cells = append(cells, Text("grow "+body, line))
		rows = append(rows, Div("row gap-1", cells...))
	}

	return Div(join("col gap-0", opts.Class), rows...)
}

type Message struct {
	Role string
	Who  string
	Text string
}

func Chat(class string, messages ...Message) Component {
	rows := make([]Component, 0, len(messages)*2)

	for _, message := range messages {
		ink, mark := "text-ink", "▸"

		switch message.Role {
		case "user":
			ink, mark = "text-accent", "▸"
		case "assistant":
			ink, mark = "text-success", "◂"
		case "system":
			ink, mark = "text-muted dim", "·"
		case "error":
			ink, mark = "text-danger", "!"
		}

		who := message.Who
		if who == "" {
			who = message.Role
		}

		rows = append(rows,
			Div("row gap-1",
				Text(ink+" bold nowrap", mark+" "+who),
			),
			Div("row gap-1 px-2",
				Text("grow", message.Text),
			),
		)
	}

	return Div(join("col gap-0", class), rows...)
}

func Logs(class string, lines []string, stamps []string) Component {
	rows := make([]Component, 0, len(lines))

	for i, line := range lines {
		cells := make([]Component, 0, 2)

		if i < len(stamps) {
			cells = append(cells, Text("w-9 text-muted dim nowrap", stamps[i]))
		}

		ink := "grow nowrap"
		switch {
		case strings.Contains(line, "error"), strings.Contains(line, "failed"):
			ink += " text-danger"
		case strings.Contains(line, "warn"):
			ink += " text-warning"
		}

		cells = append(cells, Text(ink, line))
		rows = append(rows, Div("row gap-1", cells...))
	}

	return Div(join("col gap-0", class), rows...)
}

func ScrollList(class string, height int, total int, row func(index int) Component) Component {
	if height < 1 {
		height = 10
	}

	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)
			base = withClass(withClass(base, "col"), "focusable")
			base = withClass(base, "h-"+itoa(height))
			base = withClass(base, "clip")

			n := dom.NewElement("div")
			if base != "" {
				n.Set("class", base)
			}
			if id != "" {
				n.Set("id", id)
			}

			n.Set("name", id)
			n.Set("value", "0")

			render := func() {
				offset, _ := strconv.Atoi(n.Get("value"))
				if offset > total-height {
					offset = total - height
				}
				if offset < 0 {
					offset = 0
				}
				n.Set("value", strconv.Itoa(offset))

				n.Clear()
				for i := offset; i < offset+height && i < total; i++ {
					n.Append(row(i).node(a))
				}
			}
			render()

			move := func(delta int) {
				offset, _ := strconv.Atoi(n.Get("value"))
				n.Set("value", strconv.Itoa(offset+delta))
				render()

				n.Dispatch(&dom.Event{Type: "scroll", Target: n, Value: n.Get("value")})
				if a != nil {
					a.Redraw()
				}
			}

			n.On("keydown", func(e *Event) {
				switch e.Key {
				case "down":
					e.StopPropagation()
					move(1)
				case "up":
					e.StopPropagation()
					move(-1)
				case "pagedown":
					e.StopPropagation()
					move(height)
				case "pageup":
					e.StopPropagation()
					move(-height)
				case "home":
					e.StopPropagation()
					move(-total)
				case "end":
					e.StopPropagation()
					move(total)
				}
			})

			n.On("wheelup", func(e *Event) { move(-3) })
			n.On("wheeldown", func(e *Event) { move(3) })

			return n
		},
	}
}

type PlayerOptions struct {
	Class    string
	Title    string
	Subtitle string
	Position int
	Duration int
	Playing  bool
	Volume   int
	Width    int
	Queue    []string
}

func Player(opts PlayerOptions) Component {
	if opts.Width < 8 {
		opts.Width = 28
	}

	mark := "▶"
	if opts.Playing {
		mark = "⏸"
	}

	head := Div("row gap-1",
		Text("text-accent bold nowrap", mark),
		Column("grow",
			Text("text-ink bold nowrap", opts.Title),
			Text("text-muted dim nowrap", opts.Subtitle),
		),
		Text("text-muted nowrap", clock(opts.Position)+" / "+clock(opts.Duration)),
	)

	bars := Div("row gap-1 items-center",
		SeekBar("grow", opts.Position, opts.Duration, opts.Width),
		Text("text-muted nowrap", "♪"),
		Progress("text-accent nowrap", float64(opts.Volume), 100, 6),
	)

	body := []Component{head, bars}

	if len(opts.Queue) > 0 {
		items := make([]Component, 0, len(opts.Queue))
		for i, track := range opts.Queue {
			ink := "text-muted nowrap"
			if i == 0 {
				ink = "text-ink nowrap"
			}
			items = append(items, Text(ink, track))
		}
		body = append(body, Divider("border-line"), Column("gap-0", items...))
	}

	return Div(join("col rounded border-line bg-panel px-2 gap-0", opts.Class), body...)
}

func SeekBar(class string, position, duration, width int) Component {
	if width < 4 {
		width = 20
	}

	at := scale(float64(position), float64(duration), width)
	if at >= width {
		at = width - 1
	}

	return Text(join("nowrap text-accent", class),
		strings.Repeat("━", at)+"◉"+strings.Repeat("─", width-at-1))
}

func AddressBar(class, name, value string, onGo func(string)) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)
			base = withClass(withClass(base, "row"), "gap-1")

			wrap := dom.NewElement("div")
			if base != "" {
				wrap.Set("class", base)
			}
			if id != "" {
				wrap.Set("id", id)
			}

			back := Text("nowrap text-muted focusable focus:text-ink", "‹").node(a)
			forward := Text("nowrap text-muted focusable focus:text-ink", "›").node(a)
			reload := Text("nowrap text-muted focusable focus:text-ink", "⟳").node(a)

			input := Input("grow rounded border-line px-1", name, "https://").node(a)
			input.Set("value", value)
			paintField(input)

			input.On("keydown", func(e *Event) {
				if e.Key == "enter" && onGo != nil {
					e.StopPropagation()
					onGo(input.Get("value"))
					if a != nil {
						a.Redraw()
					}
				}
			})

			wrap.Append(back)
			wrap.Append(forward)
			wrap.Append(reload)
			wrap.Append(input)

			return wrap
		},
	}
}

func Toolbar(class string, items ...Component) Component {
	return Div(join("row gap-1 items-center border-b border-line", class), items...)
}
