package gtui

import (
	"strconv"
	"strings"
	"time"

	"github.com/orayew2002/gtui/src/colors"
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

func Field(class, label string, input Component) Component {
	return Column(class, Text("text-muted nowrap", label), input)
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
			ink, width := inkTokens(base), trackWidth(base, max)
			base = join("row gap-1 focusable focus:bold focus:text-accent", base)

			n := dom.NewElement("div")
			n.Set("class", base)
			if id != "" {
				n.Set("id", id)
			}

			n.Set("name", name)
			n.Set("control", "range")
			n.Set("value", strconv.Itoa(clampInt(value, min, max)))

			render := func() {
				at, _ := strconv.Atoi(n.Get("value"))
				filled := scale(float64(at-min), float64(max-min), width)

				n.Clear()

				track := dom.NewElement("div")
				track.Set("class", join("nowrap text-accent", ink))
				track.Append(dom.NewText(strings.Repeat("━", filled) + "●" + strings.Repeat("─", width-filled)))

				label := dom.NewElement("div")
				label.Set("class", "nowrap")
				label.Append(dom.NewText(strconv.Itoa(at)))

				n.Append(track)
				n.Append(label)
			}
			render()

			n.On(controlSync, func(e *Event) {
				if e.Target != n {
					return
				}
				e.StopPropagation()

				at, _ := strconv.Atoi(n.Get("value"))
				n.Set("value", strconv.Itoa(clampInt(at, min, max)))
				render()
			})

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
			outer, inner := controlClass(base, "row gap-1 rounded border-line px-1 focusable focus:bold focus:border-accent")

			wrap := dom.NewElement("div")
			wrap.Set("class", outer)

			head := dom.NewElement("div")
			head.Set("class", inner)
			if id != "" {
				head.Set("id", id)
			}
			head.Set("name", name)
			head.Set("control", "selector")
			head.Set("value", "")
			head.Set("query", "")

			panel := dom.NewElement("div")
			panel.Set("class", "col none rounded border-line bg-panel")

			cursor := 0

			matches := func() []string {
				query := strings.ToLower(head.Get("query"))
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

			var choose func(string)

			render := func() {
				open := !panel.HasClass("none")
				found := matches()

				cursor = clamp(cursor, 0, max(0, len(found)-1))

				head.Clear()
				head.Append(text("grow nowrap", headLabel(head, placeholder, open)))
				head.Append(text("nowrap dim", caretGlyph(open)))

				panel.Clear()
				if !open {
					return
				}

				panel.Append(text("px-1 nowrap dim", "⌕ "+queryLabel(head)))

				if len(found) == 0 {
					panel.Append(text("px-1 nowrap dim", "no matches"))
					return
				}

				first := clamp(cursor-dropdownRows+1, 0, max(0, len(found)-dropdownRows))

				for i := first; i < len(found) && i < first+dropdownRows; i++ {
					row, option := optionRow(found[i], i == cursor), found[i]
					row.On("click", func(e *Event) {
						e.StopPropagation()
						choose(option)
						if a != nil {
							a.Redraw()
						}
					})
					panel.Append(row)
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
					head.Set("query", "")
					cursor = 0
				}
				render()
			}

			choose = func(option string) {
				head.Set("value", option)
				setOpen(false)

				if a != nil {
					a.Focus(head)
				}
				head.Dispatch(&dom.Event{Type: "change", Target: head, Value: option})
			}

			pick := func() {
				found := matches()
				if cursor >= len(found) {
					setOpen(false)
					return
				}
				choose(found[cursor])
			}

			head.On("click", func(e *Event) {
				setOpen(panel.HasClass("none"))
				if a != nil {
					a.Redraw()
				}
			})

			head.On("blur", func(e *Event) { setOpen(false) })

			head.On(controlSync, func(e *Event) {
				if e.Target != head {
					return
				}
				e.StopPropagation()
				render()
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
					if !open {
						return
					}
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
					trimQuery(head)
					cursor = 0
					render()

				case e.Rune == ' ' && (!open || head.Get("query") == ""):
					e.StopPropagation()
					if !open {
						setOpen(true)
					}

				case printable(e):
					e.StopPropagation()
					if !open {
						setOpen(true)
					}
					head.Set("query", head.Get("query")+string(e.Rune))
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
			outer, inner := controlClass(base, "rounded border-line px-1 focus:border-accent")

			wrap := dom.NewElement("div")
			wrap.Set("class", outer)

			input := Search(inner, name, placeholder).node(a)
			if id != "" {
				input.Set("id", id)
			}

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

			var fill func(string)

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
					row, option := optionRow(found[i], i == cursor), found[i]
					row.On("click", func(e *Event) {
						e.StopPropagation()
						fill(option)
						if a != nil {
							a.Redraw()
						}
					})
					panel.Append(row)
				}
			}

			fill = func(option string) {
				if a != nil {
					a.Focus(input)
				}

				input.Set("value", option)
				setCaret(input, len([]rune(option)))
				paintField(input)

				panel.Clear()
				panel.AddClass("none")

				input.Dispatch(&dom.Event{Type: "change", Target: input, Value: option})
			}

			accept := func() {
				found := matches()
				if cursor >= len(found) {
					return
				}
				fill(found[cursor])
			}

			input.On("input", func(e *Event) {
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
					if panel.HasClass("none") {
						return
					}
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

func headLabel(head *Node, placeholder string, open bool) string {
	if open && head.Get("query") != "" {
		return head.Get("query")
	}
	if value := head.Get("value"); value != "" {
		return value
	}
	return placeholder
}

func queryLabel(head *Node) string {
	if query := head.Get("query"); query != "" {
		return query
	}
	return "type to filter"
}

func trimQuery(head *Node) {
	query := head.Get("query")
	if query == "" {
		return
	}

	runes := []rune(query)
	head.Set("query", string(runes[:len(runes)-1]))
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

// controlClass keeps the sizing tokens on the wrapper and hands everything else
// to the control the user actually sees, so a border or a background class
// styles the box with the text in it instead of an invisible column.
func controlClass(base, defaults string) (string, string) {
	var outer, inner []string
	fill := false

	for _, token := range strings.Fields(base) {
		if !sizeToken(token) {
			inner = append(inner, token)
			continue
		}

		outer = append(outer, token)
		fill = fill || token == "grow" || strings.HasPrefix(token, "w-") ||
			strings.HasPrefix(token, "basis-")
	}

	if fill {
		inner = append([]string{"w-full"}, inner...)
	}

	return join("col", strings.Join(outer, " ")), join(defaults, strings.Join(inner, " "))
}

var sizePrefixes = []string{"w-", "min-w-", "max-w-", "basis-", "self-"}

func sizeToken(token string) bool {
	switch token {
	case "grow", "grow-0", "shrink", "shrink-0":
		return true
	}

	for _, prefix := range sizePrefixes {
		if strings.HasPrefix(token, prefix) {
			return true
		}
	}
	return false
}

func inkTokens(base string) string {
	var out []string

	for _, token := range strings.Fields(base) {
		value, ok := strings.CutPrefix(token, "text-")
		if !ok {
			continue
		}
		if _, ok := colors.Parse(value); ok {
			out = append(out, token)
		}
	}

	return strings.Join(out, " ")
}

func classWidth(base string) int {
	for _, token := range strings.Fields(base) {
		value, ok := strings.CutPrefix(token, "w-")
		if !ok {
			continue
		}
		if width, err := strconv.Atoi(value); err == nil {
			return width
		}
	}
	return 0
}

func trackWidth(base string, max int) int {
	box := classWidth(base)
	if box <= 0 {
		return 12
	}
	return clamp(box-len(strconv.Itoa(max))-2, 2, box-1)
}

func Calendar(class, name string, start time.Time) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)
			base = join("col focusable focus:bold", base)

			n := dom.NewElement("div")
			n.Set("class", base)
			if id != "" {
				n.Set("id", id)
			}
			n.Set("name", name)
			n.Set("control", "calendar")

			if start.IsZero() {
				start = time.Now()
			}
			day := start

			var (
				move func(days, months int)
				pick func(on time.Time)
			)

			render := func() {
				n.Set("value", day.Format("02.01.2006"))
				n.Clear()

				prev, next := text("nowrap dim", "‹"), text("nowrap dim", "›")
				prev.On("click", func(e *Event) { e.StopPropagation(); move(0, -1) })
				next.On("click", func(e *Event) { e.StopPropagation(); move(0, 1) })

				head := dom.NewElement("div")
				head.Set("class", "row gap-1 border-b border-line")
				head.Append(prev)
				head.Append(text("grow text-center nowrap bold", day.Format("January 2006")))
				head.Append(next)
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
						at := cursor
						cell := text(dayClass(at, day, today), strconv.Itoa(at.Day()))
						cell.On("click", func(e *Event) {
							e.StopPropagation()
							pick(at)
						})

						row.Append(cell)
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

			n.On(controlSync, func(e *Event) {
				if e.Target != n {
					return
				}
				e.StopPropagation()

				if picked, err := time.ParseInLocation("02.01.2006", n.Get("value"), day.Location()); err == nil {
					day = picked
				}
				render()
			})

			pick = func(on time.Time) {
				day = on
				render()

				n.Dispatch(&dom.Event{Type: "change", Target: n, Value: n.Get("value")})
				if a != nil {
					a.Redraw()
				}
			}

			move = func(days, months int) {
				pick(addMonths(day, months).AddDate(0, 0, days))
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

// addMonths lands on the same day of another month, clamped to its length, so
// paging from the 31st does not skip a short month the way AddDate does.
func addMonths(day time.Time, months int) time.Time {
	if months == 0 {
		return day
	}

	first := time.Date(day.Year(), day.Month()+time.Month(months), 1, 0, 0, 0, 0, day.Location())
	last := first.AddDate(0, 1, -1).Day()

	return first.AddDate(0, 0, min(day.Day(), last)-1)
}

func dayClass(cell, selected, today time.Time) string {
	class := "w-2 text-right nowrap hover:bg-raised hover:bold"

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
			ink := inkTokens(base)
			base = join("row gap-1 focusable focus:bold focus:text-accent", base)

			n := dom.NewElement("div")
			n.Set("class", base)
			if id != "" {
				n.Set("id", id)
			}

			n.Set("name", name)
			n.Set("control", "switch")
			n.State.Checked = on
			n.Set("value", boolValue(on))

			render := func() {
				n.Clear()

				caption := dom.NewElement("div")
				caption.Set("class", "grow nowrap")
				caption.Append(dom.NewText(label))

				knob := dom.NewElement("div")
				if n.State.Checked {
					knob.Set("class", join("nowrap text-success", ink))
					knob.Append(dom.NewText("●──"))
				} else {
					knob.Set("class", join("nowrap dim", ink))
					knob.Append(dom.NewText("──○"))
				}

				n.Append(caption)
				n.Append(knob)
			}
			render()

			n.On(controlSync, func(e *Event) {
				if e.Target != n {
					return
				}
				e.StopPropagation()

				n.State.Checked = truthy(n.Get("value"))
				n.Set("value", boolValue(n.State.Checked))
				render()
			})

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
