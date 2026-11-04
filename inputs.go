package gtui

import (
	"strconv"
	"strings"
	"time"

	"github.com/orayew2002/gtui/src/dom"
)

const (
	InputEmail  = "email"
	InputURL    = "url"
	InputSearch = "search"
	InputPhone  = "phone"
	InputDate   = "date"
)

func Email(class, name, placeholder string) Component {
	return field(class, name, placeholder, InputEmail, 1)
}

func URL(class, name, placeholder string) Component {
	return field(class, name, placeholder, InputURL, 1)
}

func Search(class, name, placeholder string) Component {
	return field(class, name, placeholder, InputSearch, 1)
}

func Phone(class, name, placeholder string) Component {
	return field(class, name, placeholder, InputPhone, 1)
}

func Date(class, name string) Component {
	return field(class, name, "dd.mm.yyyy", InputDate, 1)
}

func allowed(kind string, r rune) bool {
	switch kind {
	case InputNumber:
		return r >= '0' && r <= '9'
	case InputPhone:
		return r >= '0' && r <= '9' || r == '+' || r == ' ' || r == '-'
	case InputDate:
		return r >= '0' && r <= '9' || r == '.'
	case InputEmail, InputURL:
		return r > ' ' && r != ','
	}
	return true
}

func Range(class, name string, min, max, step, value int) Component {
	if step < 1 {
		step = 1
	}
	if max <= min {
		max = min + 1
	}

	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)
			base = withClass(withClass(base, "row"), "gap-1")
			base = withClass(base, "focusable")

			n := dom.NewElement("div")
			if base != "" {
				n.Set("class", base)
			}
			if id != "" {
				n.Set("id", id)
			}

			n.Set("name", name)
			n.Set("value", strconv.Itoa(clampInt(value, min, max)))

			render := func() {
				at, _ := strconv.Atoi(n.Get("value"))
				filled := scale(float64(at-min), float64(max-min), 12)

				n.Clear()

				track := dom.NewElement("div")
				track.Set("class", "nowrap text-accent")
				track.Append(dom.NewText(strings.Repeat("━", filled) + "●" + strings.Repeat("─", 12-filled)))

				label := dom.NewElement("div")
				label.Set("class", "nowrap")
				label.Append(dom.NewText(strconv.Itoa(at)))

				n.Append(track)
				n.Append(label)
			}
			render()

			move := func(delta int) {
				at, _ := strconv.Atoi(n.Get("value"))
				at = clampInt(at+delta*step, min, max)

				n.Set("value", strconv.Itoa(at))
				render()

				n.Dispatch(&dom.Event{Type: "change", Target: n, Value: n.Get("value")})
				if a != nil {
					a.Redraw()
				}
			}

			n.On("keydown", func(e *Event) {
				switch e.Key {
				case "left", "down":
					e.StopPropagation()
					move(-1)
				case "right", "up":
					e.StopPropagation()
					move(1)
				case "home":
					e.StopPropagation()
					move(-(max - min) / step)
				case "end":
					e.StopPropagation()
					move((max - min) / step)
				}
			})

			return n
		},
	}
}

const dropdownRows = 6

func Selector(class, name, placeholder string, options ...string) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)
			base = withClass(base, "col")

			wrap := dom.NewElement("div")
			if base != "" {
				wrap.Set("class", base)
			}
			if id != "" {
				wrap.Set("id", id)
			}
			wrap.Set("name", name)
			wrap.Set("value", "")
			wrap.Set("query", "")

			head := dom.NewElement("div")
			head.Set("class", "row gap-1 rounded border-line px-1 focusable focus:border-accent")

			panel := dom.NewElement("div")
			panel.Set("class", "col none rounded border-line bg-panel")

			cursor := 0

			matches := func() []string {
				query := strings.ToLower(wrap.Get("query"))
				if query == "" {
					return options
				}

				out := make([]string, 0, len(options))
				for _, option := range options {
					if strings.Contains(strings.ToLower(option), query) {
						out = append(out, option)
					}
				}
				return out
			}

			render := func() {
				open := !panel.HasClass("none")
				found := matches()

				cursor = clamp(cursor, 0, max(0, len(found)-1))

				head.Clear()
				head.Append(text("grow nowrap", headLabel(wrap, placeholder, open)))
				head.Append(text("nowrap dim", caretGlyph(open)))

				panel.Clear()
				if !open {
					return
				}

				panel.Append(text("px-1 nowrap dim", "⌕ "+queryLabel(wrap)))

				if len(found) == 0 {
					panel.Append(text("px-1 nowrap dim", "no matches"))
					return
				}

				first := clamp(cursor-dropdownRows+1, 0, max(0, len(found)-dropdownRows))

				for i := first; i < len(found) && i < first+dropdownRows; i++ {
					panel.Append(optionRow(found[i], i == cursor))
				}

				if len(found) > dropdownRows {
					panel.Append(text("px-1 nowrap dim", itoa(cursor+1)+" / "+itoa(len(found))))
				}
			}

			setOpen := func(open bool) {
				if open {
					panel.RemoveClass("none")
				} else {
					panel.AddClass("none")
					wrap.Set("query", "")
					cursor = 0
				}
				render()
			}

			pick := func() {
				found := matches()
				if cursor < len(found) {
					wrap.Set("value", found[cursor])
				}

				setOpen(false)
				wrap.Dispatch(&dom.Event{Type: "change", Target: wrap, Value: wrap.Get("value")})
			}

			head.On("click", func(e *Event) {
				setOpen(panel.HasClass("none"))
				if a != nil {
					a.Redraw()
				}
			})

			head.On("keydown", func(e *Event) {
				open := !panel.HasClass("none")

				switch {
				case e.Key == "enter":
					e.StopPropagation()
					if open {
						pick()
					} else {
						setOpen(true)
					}

				case e.Key == "esc":
					e.StopPropagation()
					setOpen(false)

				case e.Key == "down" && open:
					e.StopPropagation()
					cursor++
					render()

				case e.Key == "up" && open:
					e.StopPropagation()
					cursor--
					render()

				case e.Key == "backspace" && open:
					e.StopPropagation()
					trimQuery(wrap)
					cursor = 0
					render()

				case e.Rune >= ' ' && e.Rune != 0:
					e.StopPropagation()
					if !open {
						setOpen(true)
					}
					wrap.Set("query", wrap.Get("query")+string(e.Rune))
					cursor = 0
					render()

				default:
					return
				}

				if a != nil {
					a.Redraw()
				}
			})

			render()
			wrap.Append(head)
			wrap.Append(panel)

			return wrap
		},
	}
}

func Autocomplete(class, name, placeholder string, options ...string) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)
			base = withClass(base, "col")

			wrap := dom.NewElement("div")
			if base != "" {
				wrap.Set("class", base)
			}
			if id != "" {
				wrap.Set("id", id)
			}
			wrap.Set("name", name)

			input := Search("rounded border-line px-1 focus:border-accent", name, placeholder).node(a)
			panel := dom.NewElement("div")
			panel.Set("class", "col none rounded border-line bg-panel")

			cursor := 0

			matches := func() []string {
				query := strings.ToLower(input.Get("value"))
				if query == "" {
					return nil
				}

				out := make([]string, 0, len(options))
				for _, option := range options {
					if strings.Contains(strings.ToLower(option), query) {
						out = append(out, option)
					}
				}
				return out
			}

			render := func() {
				found := matches()
				cursor = clamp(cursor, 0, max(0, len(found)-1))

				panel.Clear()

				if input.Get("value") == "" || !input.State.Focus {
					panel.AddClass("none")
					return
				}

				panel.RemoveClass("none")

				if len(found) == 0 {
					panel.Append(text("px-1 nowrap dim", "no matches"))
					return
				}

				for i := 0; i < len(found) && i < dropdownRows; i++ {
					panel.Append(optionRow(found[i], i == cursor))
				}
			}

			accept := func() {
				found := matches()
				if cursor >= len(found) {
					return
				}

				input.Set("value", found[cursor])
				setCaret(input, len([]rune(found[cursor])))
				paintField(input)

				panel.Clear()
				panel.AddClass("none")

				wrap.Set("value", found[cursor])
				wrap.Dispatch(&dom.Event{Type: "change", Target: wrap, Value: found[cursor]})
			}

			input.On("input", func(e *Event) {
				wrap.Set("value", input.Get("value"))
				cursor = 0
				render()
			})

			input.On("focus", func(e *Event) { render() })
			input.On("blur", func(e *Event) {
				panel.Clear()
				panel.AddClass("none")
			})

			input.On("keydown", func(e *Event) {
				switch {
				case e.Key == "down":
					e.StopPropagation()
					cursor++
					render()
				case e.Key == "up":
					e.StopPropagation()
					cursor--
					render()
				case e.Key == "enter":
					e.StopPropagation()
					accept()
				case e.Key == "esc":
					e.StopPropagation()
					panel.Clear()
					panel.AddClass("none")
				default:
					return
				}

				if a != nil {
					a.Redraw()
				}
			})

			wrap.Append(input)
			wrap.Append(panel)
			render()

			return wrap
		},
	}
}

func optionRow(option string, active bool) *Node {
	class, mark := "row gap-1 px-1", " "
	if active {
		class, mark = "row gap-1 px-1 bg-accent text-canvas bold", "›"
	}

	row := dom.NewElement("div")
	row.Set("class", class)
	row.Append(text("w-1 nowrap", mark))
	row.Append(text("grow nowrap", option))

	return row
}

func text(class, content string) *Node {
	n := dom.NewElement("div")
	if class != "" {
		n.Set("class", class)
	}
	n.Append(dom.NewText(content))
	return n
}

func headLabel(wrap *Node, placeholder string, open bool) string {
	if open && wrap.Get("query") != "" {
		return wrap.Get("query")
	}
	if value := wrap.Get("value"); value != "" {
		return value
	}
	return placeholder
}

func queryLabel(wrap *Node) string {
	if query := wrap.Get("query"); query != "" {
		return query
	}
	return "type to filter"
}

func trimQuery(wrap *Node) {
	query := wrap.Get("query")
	if query == "" {
		return
	}

	runes := []rune(query)
	wrap.Set("query", string(runes[:len(runes)-1]))
}

func caretGlyph(open bool) string {
	if open {
		return "⌃"
	}
	return "⌄"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Calendar(class, name string, start time.Time) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)
			base = withClass(withClass(base, "col"), "focusable")

			n := dom.NewElement("div")
			if base != "" {
				n.Set("class", base)
			}
			if id != "" {
				n.Set("id", id)
			}
			n.Set("name", name)

			if start.IsZero() {
				start = time.Now()
			}
			day := start

			render := func() {
				n.Set("value", day.Format("02.01.2006"))
				n.Clear()

				head := dom.NewElement("div")
				head.Set("class", "row gap-1 border-b border-line")
				head.Append(text("nowrap dim", "‹"))
				head.Append(text("grow text-center nowrap bold", day.Format("January 2006")))
				head.Append(text("nowrap dim", "›"))
				n.Append(head)

				week := dom.NewElement("div")
				week.Set("class", "row gap-1")
				for i, label := range []string{"mo", "tu", "we", "th", "fr", "sa", "su"} {
					class := "w-2 text-right nowrap dim"
					if i >= 5 {
						class = "w-2 text-right nowrap dim text-danger"
					}
					week.Append(text(class, label))
				}
				n.Append(week)

				today := time.Now()
				first := time.Date(day.Year(), day.Month(), 1, 0, 0, 0, 0, day.Location())
				cursor := first.AddDate(0, 0, -((int(first.Weekday()) + 6) % 7))

				for range 6 {
					row := dom.NewElement("div")
					row.Set("class", "row gap-1")

					for range 7 {
						row.Append(text(dayClass(cursor, day, today), strconv.Itoa(cursor.Day())))
						cursor = cursor.AddDate(0, 0, 1)
					}

					n.Append(row)
				}

				foot := dom.NewElement("div")
				foot.Set("class", "row gap-1 border-t border-line")
				foot.Append(text("grow nowrap text-accent bold", day.Format("Mon 02.01.2006")))
				foot.Append(text("nowrap dim", "←→ ↑↓ ⇞⇟"))
				n.Append(foot)
			}

			render()

			move := func(days, months int) {
				day = day.AddDate(0, months, days)
				render()

				n.Dispatch(&dom.Event{Type: "change", Target: n, Value: n.Get("value")})
				if a != nil {
					a.Redraw()
				}
			}

			n.On("keydown", func(e *Event) {
				switch e.Key {
				case "left":
					e.StopPropagation()
					move(-1, 0)
				case "right":
					e.StopPropagation()
					move(1, 0)
				case "up":
					e.StopPropagation()
					move(-7, 0)
				case "down":
					e.StopPropagation()
					move(7, 0)
				case "pageup":
					e.StopPropagation()
					move(0, -1)
				case "pagedown":
					e.StopPropagation()
					move(0, 1)
				case "enter":
					e.StopPropagation()
					n.Dispatch(&dom.Event{Type: "select", Target: n, Value: n.Get("value")})
				}
			})

			return n
		},
	}
}

func dayClass(cell, selected, today time.Time) string {
	class := "w-2 text-right nowrap"

	switch {
	case sameDay(cell, selected):
		return class + " bg-accent text-canvas bold"
	case cell.Month() != selected.Month():
		return class + " dim"
	case sameDay(cell, today):
		return class + " text-accent bold underline"
	case cell.Weekday() == time.Saturday, cell.Weekday() == time.Sunday:
		return class + " text-danger"
	}

	return class
}

func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

func Switch(class, name, label string, on bool) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)
			base = withClass(withClass(base, "row"), "gap-1")
			base = withClass(base, "focusable")

			n := dom.NewElement("div")
			if base != "" {
				n.Set("class", base)
			}
			if id != "" {
				n.Set("id", id)
			}

			n.Set("name", name)
			n.State.Checked = on
			n.Set("value", boolValue(on))

			render := func() {
				n.Clear()

				caption := dom.NewElement("div")
				caption.Set("class", "grow nowrap")
				caption.Append(dom.NewText(label))

				knob := dom.NewElement("div")
				if n.State.Checked {
					knob.Set("class", "nowrap text-success")
					knob.Append(dom.NewText("●──"))
				} else {
					knob.Set("class", "nowrap dim")
					knob.Append(dom.NewText("──○"))
				}

				n.Append(caption)
				n.Append(knob)
			}
			render()

			flip := func() {
				n.State.Checked = !n.State.Checked
				n.Set("value", boolValue(n.State.Checked))
				render()

				n.Dispatch(&dom.Event{Type: "change", Target: n, Value: n.Get("value")})
				if a != nil {
					a.Redraw()
				}
			}

			n.On("click", func(e *Event) { flip() })
			n.On("keydown", func(e *Event) {
				if e.Key == "enter" || e.Rune == ' ' {
					e.StopPropagation()
					flip()
				}
			})

			return n
		},
	}
}
