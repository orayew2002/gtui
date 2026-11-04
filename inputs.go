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
			head.Set("class", "row gap-1 focusable nowrap")

			list := dom.NewElement("div")
			list.Set("class", "col none")

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

			var render func()

			render = func() {
				open := !list.HasClass("none")

				head.Clear()
				caption := wrap.Get("value")
				if caption == "" {
					caption = placeholder
				}
				if open && wrap.Get("query") != "" {
					caption = wrap.Get("query")
				}

				text := dom.NewElement("div")
				text.Set("class", "grow nowrap")
				text.Append(dom.NewText(caption))

				mark := dom.NewElement("div")
				mark.Set("class", "nowrap dim")
				if open {
					mark.Append(dom.NewText("⌕"))
				} else {
					mark.Append(dom.NewText("▾"))
				}

				head.Append(text)
				head.Append(mark)

				list.Clear()
				found := matches()

				if len(found) == 0 {
					empty := dom.NewElement("div")
					empty.Set("class", "px-1 nowrap dim")
					empty.Append(dom.NewText("no matches"))
					list.Append(empty)
					return
				}

				if cursor >= len(found) {
					cursor = len(found) - 1
				}
				if cursor < 0 {
					cursor = 0
				}

				for i, option := range found {
					row := dom.NewElement("div")
					class := "px-1 nowrap"
					if i == cursor {
						class += " bg-raised bold"
					}
					row.Set("class", class)
					row.Append(dom.NewText(option))
					list.Append(row)
				}
			}

			setOpen := func(open bool) {
				if open {
					list.RemoveClass("none")
				} else {
					list.AddClass("none")
					wrap.Set("query", "")
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
				setOpen(list.HasClass("none"))
				if a != nil {
					a.Redraw()
				}
			})

			head.On("keydown", func(e *Event) {
				open := !list.HasClass("none")

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
					query := wrap.Get("query")
					if query != "" {
						runes := []rune(query)
						wrap.Set("query", string(runes[:len(runes)-1]))
					}
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
			wrap.Append(list)

			return wrap
		},
	}
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
				head.Set("class", "row gap-1")

				prev := dom.NewElement("div")
				prev.Set("class", "nowrap dim")
				prev.Append(dom.NewText("‹"))

				title := dom.NewElement("div")
				title.Set("class", "grow text-center nowrap bold")
				title.Append(dom.NewText(day.Format("January 2006")))

				next := dom.NewElement("div")
				next.Set("class", "nowrap dim")
				next.Append(dom.NewText("›"))

				head.Append(prev)
				head.Append(title)
				head.Append(next)
				n.Append(head)

				week := dom.NewElement("div")
				week.Set("class", "row gap-1")
				for _, label := range []string{"mo", "tu", "we", "th", "fr", "sa", "su"} {
					cell := dom.NewElement("div")
					cell.Set("class", "w-2 nowrap dim")
					cell.Append(dom.NewText(label))
					week.Append(cell)
				}
				n.Append(week)

				first := time.Date(day.Year(), day.Month(), 1, 0, 0, 0, 0, day.Location())
				offset := (int(first.Weekday()) + 6) % 7
				cursor := first.AddDate(0, 0, -offset)

				for range 6 {
					row := dom.NewElement("div")
					row.Set("class", "row gap-1")

					for range 7 {
						cell := dom.NewElement("div")
						class := "w-2 nowrap"

						switch {
						case cursor.Month() != day.Month():
							class += " dim"
						case cursor.Day() == day.Day():
							class += " bg-accent text-canvas bold"
						}

						cell.Set("class", class)
						cell.Append(dom.NewText(strconv.Itoa(cursor.Day())))
						row.Append(cell)

						cursor = cursor.AddDate(0, 0, 1)
					}

					n.Append(row)
				}
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
