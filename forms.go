package gtui

import (
	"strings"

	"github.com/orayew2002/gtui/src/dom"
)

const (
	InputText     = "text"
	InputPassword = "password"
	InputNumber   = "number"
	InputArea     = "area"
)

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
			base = withClass(base, "focusable")
			base = withClass(base, "nowrap")
			if kind == InputArea {
				base = withClass(base, "h-"+itoa(rows))
			}

			n := dom.NewElement("div")
			if base != "" {
				n.Set("class", base)
			}
			if id != "" {
				n.Set("id", id)
			}

			n.Set("name", name)
			n.Set("type", kind)
			n.Set("value", "")
			n.Set("placeholder", placeholder)

			render := func() { paintField(n) }
			render()

			n.On("focus", func(e *Event) { render() })
			n.On("blur", func(e *Event) { render() })

			n.On("keydown", func(e *Event) {
				if !editField(n, e, kind) {
					return
				}

				e.StopPropagation()
				render()

				n.Dispatch(&dom.Event{
					Type:   "input",
					Target: n,
					Value:  n.Get("value"),
				})

				if a != nil {
					a.Redraw()
				}
			})

			return n
		},
	}
}

func editField(n *Node, e *Event, kind string) bool {
	value := n.Get("value")

	switch {
	case e.Key == "backspace":
		if value == "" {
			return false
		}
		runes := []rune(value)
		n.Set("value", string(runes[:len(runes)-1]))
		return true

	case e.Key == "enter" && kind == InputArea:
		n.Set("value", value+"\n")
		return true

	case e.Rune >= ' ' && e.Rune != 0:
		if !allowed(kind, e.Rune) {
			return false
		}
		n.Set("value", value+string(e.Rune))
		return true
	}

	return false
}

func paintField(n *Node) {
	value := n.Get("value")
	shown := value

	if n.Get("type") == InputPassword {
		shown = strings.Repeat("•", len([]rune(value)))
	}

	ink := ""
	if shown == "" {
		shown = n.Get("placeholder")
		ink = "dim"
	}

	if n.State.Focus {
		shown += "▌"
	}

	n.Clear()

	text := dom.NewElement("div")
	if ink != "" {
		text.Set("class", ink)
	}
	if n.Get("type") == InputArea {
		text.Set("class", withClass(ink, "pre"))
	}
	text.Append(dom.NewText(shown))

	n.Append(text)
}

func Checkbox(class, name, label string, checked bool) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)
			base = withClass(base, "row")
			base = withClass(base, "gap-1")
			base = withClass(base, "focusable")

			n := dom.NewElement("div")
			if base != "" {
				n.Set("class", base)
			}
			if id != "" {
				n.Set("id", id)
			}
			n.Set("name", name)
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
			base = withClass(base, "col")

			group := dom.NewElement("div")
			if base != "" {
				group.Set("class", base)
			}
			if id != "" {
				group.Set("id", id)
			}
			group.Set("name", name)
			group.Set("value", "")

			items := make([]*Node, 0, len(options))

			for _, option := range options {
				label := option

				item := dom.NewElement("div")
				item.Set("class", "row gap-1 focusable")
				item.Set("value", label)

				render := func() {
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
					caption.Append(dom.NewText(label))

					item.Append(dot)
					item.Append(caption)
				}
				render()

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

			button := dom.NewElement("div")
			button.Set("class", "row gap-1 focusable nowrap")

			list := dom.NewElement("div")
			list.Set("class", "col none rounded border")

			label := func() string {
				if value := wrap.Get("value"); value != "" {
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

			for _, option := range options {
				choice := option

				item := dom.NewElement("div")
				item.Set("class", "px-1 nowrap focusable hover:bold focus:bold")
				item.Append(dom.NewText(choice))

				pick := func(e *Event) {
					wrap.Set("value", choice)
					setOpen(false)

					wrap.Dispatch(&dom.Event{Type: "change", Target: wrap, Value: choice})
					if a != nil {
						a.Redraw()
					}
				}

				item.On("click", pick)
				item.On("keydown", func(e *Event) {
					if e.Key == "enter" || e.Rune == ' ' {
						e.StopPropagation()
						pick(e)
					}
				})

				list.Append(item)
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
				case e.Key == "esc":
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
	if node.Has("type") {
		paintField(node)
	}

	return a.Redraw()
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
