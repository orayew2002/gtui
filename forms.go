package gtui

import (
	"strconv"
	"strings"

	"github.com/orayew2002/gtui/src/dom"
	"github.com/orayew2002/gtui/src/layout"
	"github.com/orayew2002/gtui/src/render"
)

const (
	InputText     = "text"
	InputPassword = "password"
	InputNumber   = "number"
	InputArea     = "area"
)

const controlSync = "sync"

func Input(class, name, placeholder string) Component {
	return field(class, name, placeholder, InputText, 1)
}

func Password(class, name, placeholder string) Component {
	return field(class, name, placeholder, InputPassword, 1)
}

func Number(class, name, placeholder string) Component {
	return field(class, name, placeholder, InputNumber, 1)
}

func TextArea(class, name string, rows int) Component {
	if rows < 1 {
		rows = 3
	}
	return field(class, name, "", InputArea, rows)
}

func field(class, name, placeholder, kind string, rows int) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)
			base = join("focusable nowrap clip focus:bold", base)

			n := dom.NewElement("div")
			n.Set("class", base)
			if id != "" {
				n.Set("id", id)
			}

			n.Set("name", name)
			n.Set("type", kind)
			n.Set("control", "field")
			n.Set("value", "")
			n.Set("cursor", "0")
			n.Set("placeholder", placeholder)
			if kind == InputArea {
				n.Set("rows", itoa(rows))
			}

			render := func() {
				measureField(a, n)
				paintField(n)
			}
			render()

			n.On("focus", func(e *Event) { render() })
			n.On("blur", func(e *Event) { render() })

			n.On(controlSync, func(e *Event) {
				if e.Target != n {
					return
				}
				e.StopPropagation()

				setCaret(n, len([]rune(n.Get("value"))))
				render()
			})

			n.On("keydown", func(e *Event) {
				consumed, changed := editField(n, e, kind)
				if !consumed {
					return
				}

				e.StopPropagation()
				render()

				if changed {
					n.Dispatch(&dom.Event{
						Type:   "input",
						Target: n,
						Value:  n.Get("value"),
					})
				}

				if a != nil {
					a.Redraw()
				}
			})

			return n
		},
	}
}

func caret(n *Node) int {
	at, err := strconv.Atoi(n.Get("cursor"))
	if err != nil {
		at = len([]rune(n.Get("value")))
	}
	return clamp(at, 0, len([]rune(n.Get("value"))))
}

func setCaret(n *Node, at int) {
	n.Set("cursor", strconv.Itoa(clamp(at, 0, len([]rune(n.Get("value"))))))
}

// editField reports whether the field owns the key and whether the value moved:
// a focused field swallows every printable rune and backspace/delete even when
// it rejects them, so a stray letter cannot reach the app shortcuts, while the
// "input" event stays reserved for real edits.
func editField(n *Node, e *Event, kind string) (consumed, changed bool) {
	runes := []rune(n.Get("value"))
	at := caret(n)

	write := func(next []rune, cursor int) (bool, bool) {
		n.Set("value", string(next))
		setCaret(n, cursor)
		return true, true
	}

	switch {
	case e.Key == "left":
		if at == 0 {
			return false, false
		}
		setCaret(n, at-1)
		return true, false

	case e.Key == "right":
		if at >= len(runes) {
			return false, false
		}
		setCaret(n, at+1)
		return true, false

	case e.Key == "up" && kind == InputArea:
		setCaret(n, moveLine(runes, at, -1))
		return true, false

	case e.Key == "down" && kind == InputArea:
		setCaret(n, moveLine(runes, at, 1))
		return true, false

	case e.Key == "home":
		if kind != InputArea {
			setCaret(n, 0)
			return true, false
		}
		start, _ := lineBounds(runes, at)
		setCaret(n, start)
		return true, false

	case e.Key == "end":
		if kind != InputArea {
			setCaret(n, len(runes))
			return true, false
		}
		_, end := lineBounds(runes, at)
		setCaret(n, end)
		return true, false

	case e.Key == "backspace":
		if at == 0 {
			return true, false
		}
		return write(append(append([]rune{}, runes[:at-1]...), runes[at:]...), at-1)

	case e.Key == "delete":
		if at >= len(runes) {
			return true, false
		}
		return write(append(append([]rune{}, runes[:at]...), runes[at+1:]...), at)

	case e.Key == "enter" && kind == InputArea:
		return write(insertRune(runes, at, '\n'), at+1)

	case printable(e):
		if !allowed(kind, e.Rune) {
			return true, false
		}
		return write(insertRune(runes, at, e.Rune), at+1)
	}

	return false, false
}

func printable(e *Event) bool {
	if e.Rune < ' ' {
		return false
	}
	return !hotkey(e.Key)
}

func hotkey(key string) bool {
	return strings.HasPrefix(key, "ctrl+") || strings.HasPrefix(key, "alt+") ||
		strings.HasPrefix(key, "meta+") || strings.HasPrefix(key, "super+")
}

func lineBounds(runes []rune, at int) (int, int) {
	start, end := clamp(at, 0, len(runes)), clamp(at, 0, len(runes))

	for start > 0 && runes[start-1] != '\n' {
		start--
	}
	for end < len(runes) && runes[end] != '\n' {
		end++
	}

	return start, end
}

func moveLine(runes []rune, at, delta int) int {
	start, end := lineBounds(runes, at)
	col := at - start

	if delta < 0 {
		if start == 0 {
			return at
		}
		prev, prevEnd := lineBounds(runes, start-1)
		return prev + min(col, prevEnd-prev)
	}

	if end >= len(runes) {
		return at
	}

	next, nextEnd := lineBounds(runes, end+1)
	return next + min(col, nextEnd-next)
}

func insertRune(runes []rune, at int, r rune) []rune {
	out := make([]rune, 0, len(runes)+1)
	out = append(out, runes[:at]...)
	out = append(out, r)
	return append(out, runes[at:]...)
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func paintField(n *Node) {
	kind := n.Get("type")
	value := n.Get("value")

	shown := value
	if kind == InputPassword {
		shown = strings.Repeat("•", len([]rune(value)))
	}

	ink := ""
	if shown == "" {
		shown, ink = n.Get("placeholder"), "dim"
	}

	runes := []rune(shown)
	at := 0
	if value != "" {
		at = clamp(caret(n), 0, len(runes))
	}

	scroll := 0
	if kind == InputArea {
		scroll = areaScroll(runes, at, fieldRows(n))
	} else {
		runes, at = fieldWindow(n, runes, at)
	}

	if n.State.Focus {
		runes = insertRune(runes, at, '▌')
	}

	n.Clear()

	text := dom.NewElement("div")
	if class := fieldTextClass(n, kind, ink); class != "" {
		text.Set("class", class)
	}
	text.ScrollY = scroll
	text.Append(dom.NewText(string(runes)))

	n.Append(text)
}

func fieldTextClass(n *Node, kind, ink string) string {
	if kind != InputArea {
		return ink
	}
	if rows := n.Get("rows"); rows != "" {
		return join(ink, "pre clip h-"+rows)
	}
	return join(ink, "pre")
}

func fieldRows(n *Node) int {
	rows, err := strconv.Atoi(n.Get("rows"))
	if err != nil {
		return 0
	}
	return rows
}

// areaScroll keeps the caret line inside the rows a text area shows.
func areaScroll(runes []rune, at, rows int) int {
	if rows < 1 {
		return 0
	}

	line := 0
	for _, r := range runes[:at] {
		if r == '\n' {
			line++
		}
	}

	if line < rows {
		return 0
	}
	return line - rows + 1
}

// fieldWindow slides a window over a value that is wider than its box so the
// caret stays on screen, and remembers where the window starts in ScrollX.
func fieldWindow(n *Node, runes []rune, at int) ([]rune, int) {
	width := fieldView(n)
	if width <= 0 {
		return runes, at
	}

	pad := 0
	if n.State.Focus {
		pad = 1
	}

	start := clamp(n.ScrollX, 0, at)
	head, tail := runeCells(runes[start:at]), runeCells(runes[at:])

	for start > 0 && head+tail+render.RuneWidth(runes[start-1])+pad <= width {
		start--
		head += render.RuneWidth(runes[start])
	}
	for start < at && head+pad > width {
		head -= render.RuneWidth(runes[start])
		start++
	}
	n.ScrollX = start

	end, shown := len(runes), head+tail
	for end > at && shown+pad > width {
		end--
		shown -= render.RuneWidth(runes[end])
	}

	return runes[start:end], at - start
}

func fieldView(n *Node) int {
	width, err := strconv.Atoi(n.Get("view"))
	if err != nil {
		return 0
	}
	return width
}

// remeasureFields hands every field the width layout just gave it and repaints
// the ones whose width changed. A field windows its value around the caret, and
// that window needs a width, so without this pass a value longer than the field
// is painted whole on the first frame: the overflow is clipped and the caret
// goes with it, until the next keystroke happens to re-measure.
func (a *App) remeasureFields() bool {
	if a.tree == nil {
		return false
	}

	changed := false

	a.tree.Walk(func(b *layout.Box) bool {
		n := b.Node
		if n == nil || b.Content.W <= 0 || n.Get("control") != "field" {
			return true
		}

		width := itoa(visibleWidth(b))
		if n.Get("view") == width {
			return true
		}

		n.Set("view", width)
		paintField(n)
		changed = true
		return true
	})

	return changed
}

// visibleWidth is how much of a box the viewer can actually see. A field inside
// a narrower clipping pane has to window its value into the visible part, or
// the caret is painted in cells the pane cuts away. A pane scrolled past the
// field's left edge is left alone, because scrolling the focused node back into
// view already fixes that case.
func visibleWidth(b *layout.Box) int {
	area := b.Content.Intersect(b.Clip)
	if area.Empty() || area.X > b.Content.X || area.W >= b.Content.W {
		return b.Content.W
	}
	return area.W
}

func measureField(a *App, n *Node) {
	if a == nil {
		return
	}
	if box := a.findBox(n); box != nil && box.Content.W > 0 {
		n.Set("view", itoa(box.Content.W))
	}
}

func runeCells(runes []rune) int {
	width := 0
	for _, r := range runes {
		width += render.RuneWidth(r)
	}
	return width
}

func Checkbox(class, name, label string, checked bool) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)
			base = join("row gap-1 focusable focus:bold focus:text-accent", base)

			n := dom.NewElement("div")
			n.Set("class", base)
			if id != "" {
				n.Set("id", id)
			}
			n.Set("name", name)
			n.Set("control", "checkbox")
			n.State.Checked = checked
			n.Set("value", boolValue(checked))

			render := func() {
				n.Clear()

				mark := "☐"
				if n.State.Checked {
					mark = "☑"
				}

				box := dom.NewElement("div")
				box.Set("class", "nowrap")
				box.Append(dom.NewText(mark))

				caption := dom.NewElement("div")
				caption.Set("class", "grow nowrap")
				caption.Append(dom.NewText(label))

				n.Append(box)
				n.Append(caption)
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

			toggle := func(e *Event) {
				n.State.Checked = !n.State.Checked
				n.Set("value", boolValue(n.State.Checked))
				render()

				n.Dispatch(&dom.Event{Type: "change", Target: n, Value: n.Get("value")})
				if a != nil {
					a.Redraw()
				}
			}

			n.On("click", toggle)
			n.On("keydown", func(e *Event) {
				if e.Key == "enter" || e.Rune == ' ' {
					e.StopPropagation()
					toggle(e)
				}
			})

			return n
		},
	}
}

func Radios(class, name string, options ...string) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)
			base = join("col", base)

			group := dom.NewElement("div")
			group.Set("class", base)
			if id != "" {
				group.Set("id", id)
			}
			group.Set("name", name)
			group.Set("control", "radios")
			group.Set("value", "")

			items := make([]*Node, 0, len(options))

			for _, option := range options {
				item := dom.NewElement("div")
				item.Set("class", "row gap-1 focusable focus:bold focus:text-accent")
				item.Set("value", option)

				items = append(items, item)
				group.Append(item)
			}

			renderAll := func() {
				for _, item := range items {
					item.Clear()

					mark := "○"
					if item.State.Checked {
						mark = "◉"
					}

					dot := dom.NewElement("div")
					dot.Set("class", "nowrap")
					dot.Append(dom.NewText(mark))

					caption := dom.NewElement("div")
					caption.Set("class", "grow nowrap")
					caption.Append(dom.NewText(item.Get("value")))

					item.Append(dot)
					item.Append(caption)
				}
			}
			renderAll()

			group.On(controlSync, func(e *Event) {
				if e.Target != group {
					return
				}
				e.StopPropagation()

				for _, item := range items {
					item.State.Checked = item.Get("value") == group.Get("value")
				}
				renderAll()
			})

			for _, item := range items {
				choice := item

				pick := func(e *Event) {
					for _, other := range items {
						other.State.Checked = other == choice
					}
					group.Set("value", choice.Get("value"))
					renderAll()

					group.Dispatch(&dom.Event{Type: "change", Target: group, Value: group.Get("value")})
					if a != nil {
						a.Redraw()
					}
				}

				choice.On("click", pick)
				choice.On("keydown", func(e *Event) {
					if e.Key == "enter" || e.Rune == ' ' {
						e.StopPropagation()
						pick(e)
					}
				})
			}

			return group
		},
	}
}

func Select(class, name, placeholder string, options ...string) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)
			outer, inner := controlClass(base, "row gap-1 nowrap focusable focus:bold focus:text-accent")

			wrap := dom.NewElement("div")
			wrap.Set("class", outer)

			button := dom.NewElement("div")
			button.Set("class", inner)
			if id != "" {
				button.Set("id", id)
			}
			button.Set("name", name)
			button.Set("control", "select")
			button.Set("value", "")

			list := dom.NewElement("div")
			list.Set("class", "col none rounded border")

			label := func() string {
				if value := button.Get("value"); value != "" {
					return value
				}
				return placeholder
			}

			renderButton := func() {
				button.Clear()

				caption := dom.NewElement("div")
				caption.Set("class", "grow nowrap")
				caption.Append(dom.NewText(label()))

				arrow := dom.NewElement("div")
				arrow.Set("class", "nowrap dim")
				if list.HasClass("none") {
					arrow.Append(dom.NewText("▾"))
				} else {
					arrow.Append(dom.NewText("▴"))
				}

				button.Append(caption)
				button.Append(arrow)
			}

			setOpen := func(open bool) {
				if open {
					list.RemoveClass("none")
				} else {
					list.AddClass("none")
				}
				renderButton()
			}

			dismiss := func() {
				setOpen(false)
				if a != nil {
					a.Focus(button)
				}
			}

			var (
				items []*Node
				step  func(delta int)
			)

			// a list left open after focus moved on would float over whatever
			// the user went to look at
			closeIfAway := func(e *Event) {
				if a == nil || list.HasClass("none") {
					return
				}
				if wrap.Contains(a.Focused()) {
					return
				}
				setOpen(false)
				a.Redraw()
			}

			button.On("blur", closeIfAway)

			button.On(controlSync, func(e *Event) {
				if e.Target != button {
					return
				}
				e.StopPropagation()
				renderButton()
			})

			for _, option := range options {
				choice := option

				item := dom.NewElement("div")
				item.Set("class", "px-1 nowrap focusable hover:bold focus:bold")
				item.Append(dom.NewText(choice))

				pick := func(e *Event) {
					button.Set("value", choice)
					dismiss()

					button.Dispatch(&dom.Event{Type: "change", Target: button, Value: choice})
					if a != nil {
						a.Redraw()
					}
				}

				item.On("blur", closeIfAway)
				item.On("click", pick)
				item.On("keydown", func(e *Event) {
					switch {
					case e.Key == "enter", e.Rune == ' ':
						e.StopPropagation()
						pick(e)
					case e.Key == "down":
						e.StopPropagation()
						step(1)
					case e.Key == "up":
						e.StopPropagation()
						step(-1)
					case e.Key == "esc":
						e.StopPropagation()
						dismiss()
					}
				})

				items = append(items, item)
				list.Append(item)
			}

			// opening a list the arrows cannot reach would leave the control
			// usable by mouse only
			step = func(delta int) {
				if a == nil || len(items) == 0 {
					return
				}
				setOpen(true)

				at := -1
				if delta < 0 {
					at = len(items)
				}
				for i, item := range items {
					if item == a.Focused() {
						at = i
						break
					}
				}

				a.Focus(items[clampInt(at+delta, 0, len(items)-1)])
				a.Redraw()
			}

			toggle := func(e *Event) {
				setOpen(list.HasClass("none"))
				if a != nil {
					a.Redraw()
				}
			}

			button.On("click", toggle)
			button.On("keydown", func(e *Event) {
				switch {
				case e.Key == "enter", e.Rune == ' ':
					e.StopPropagation()
					toggle(e)
				case e.Key == "down":
					e.StopPropagation()
					step(1)
				case e.Key == "up":
					e.StopPropagation()
					step(-1)
				case e.Key == "esc":
					if list.HasClass("none") {
						return
					}
					e.StopPropagation()
					setOpen(false)
				}
			})

			renderButton()
			wrap.Append(button)
			wrap.Append(list)

			return wrap
		},
	}
}

func (a *App) Value(sel string) string {
	if node := a.Query(sel); node != nil {
		return node.Get("value")
	}
	return ""
}

func (a *App) SetValue(sel, value string) *App {
	node := a.Query(sel)
	if node == nil {
		return a
	}

	node.Set("value", value)
	syncControl(node)

	return a.Redraw()
}

// syncControl redraws a control from the value it carries, so setting a value
// from code moves the glyph, the caret and the selection, not just the attribute.
func syncControl(n *Node) {
	if n == nil || !n.Has("control") {
		return
	}
	n.Dispatch(&dom.Event{Type: controlSync, Target: n, Value: n.Get("value")})
}

func (a *App) Values() map[string]string {
	out := map[string]string{}

	for _, node := range a.QueryAll("div") {
		if name := node.Get("name"); name != "" {
			out[name] = node.Get("value")
		}
	}

	return out
}

func boolValue(on bool) string {
	if on {
		return "true"
	}
	return "false"
}

func truthy(value string) bool {
	switch strings.ToLower(value) {
	case "true", "1", "on", "yes":
		return true
	}
	return false
}
