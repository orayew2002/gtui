package gtui

import (
	"time"

	"github.com/orayew2002/gtui/src/dom"
)

const (
	modalID    = "gtui-modal"
	toastID    = "gtui-toast"
	toastLife  = 2 * time.Second
	toastStack = 3
)

func (a *App) Modal(title string, body ...Component) *App {
	layer := a.layer(modalID, "absolute top-0 left-0 full z-9 center")

	panel := Div("col rounded w-1/2 bg-panel px-2 gap-1 #gtui-modal-panel",
		append([]Component{
			Div("row gap-1 border-b",
				Text("bold grow nowrap", title),
				Text("dim nowrap", "esc"),
			),
		}, body...)...,
	)

	node := panel.node(a)

	layer.Clear()
	layer.Append(node)
	layer.RemoveClass("none")

	// a modal opened over a modal still returns to the page, not to the panel
	// its own Clear just detached
	if !a.modalOpen {
		a.beforeModal = a.focus
	}
	a.modalOpen = true

	// whatever the pointer was over is now behind the panel
	a.setHover(nil)

	first := firstFocusableNode(layer)
	if first == nil {
		node.AddClass("focusable")
		first = node
	}
	a.Focus(first)

	return a.Redraw()
}

func (a *App) CloseModal() *App {
	if !a.modalOpen {
		return a
	}
	a.modalOpen = false

	if layer := a.Query("#" + modalID); layer != nil {
		layer.Clear()
		layer.AddClass("none")
	}

	restore := a.beforeModal
	a.beforeModal = nil

	if restore != nil && !a.attached(restore) {
		restore = a.byID(restore.ID)
	}
	a.Focus(restore)

	return a.Redraw()
}

func (a *App) ModalOpen() bool { return a.modalOpen }

func (a *App) Toast(kind, message string) *App {
	layer := a.layer(toastID, "absolute top-0 right-0 z-9 px-1 col")

	ink := "text-sky"
	switch kind {
	case "success":
		ink = "text-green"
	case "warning":
		ink = "text-yellow"
	case "danger":
		ink = "text-red"
	}

	a.toastSeq++

	toast := Div("row rounded px-2 bg-panel "+ink, Text("nowrap", message)).node(a)
	toast.Set("id", toastID+"-"+itoa(a.toastSeq))

	for len(layer.Children) >= toastStack {
		layer.Children[0].Detach()
	}

	layer.Append(toast)
	layer.RemoveClass("none")

	// each toast times out on its own node, so a later one cannot cut the
	// earlier ones short
	a.After(toastLife, func() {
		toast.Detach()
		if len(layer.Children) == 0 {
			layer.AddClass("none")
		}
		a.Redraw()
	})

	return a.Redraw()
}

func firstFocusableNode(root *Node) *Node {
	var found *Node

	root.Walk(func(n *Node) bool {
		if found != nil {
			return false
		}
		if n == root || n.Type != dom.ElementNode {
			return true
		}
		if n.HasClass("focusable") || n.Has("href") {
			found = n
			return false
		}
		return true
	})

	return found
}

func (a *App) layer(id, class string) *Node {
	if node := a.Query("#" + id); node != nil {
		return node
	}

	node := dom.NewElement("div")
	node.Set("class", class)
	node.Set("id", id)
	a.root.Append(node)

	return node
}

func Dropdown(class, label string, items ...Component) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)

			// sizing belongs to the wrapper, everything else to the button, so
			// colours, padding and a border radius land where they can be seen
			outer, inner := controlClass(base, "row gap-1 focusable nowrap focus:bold")

			wrap := dom.NewElement("div")
			wrap.Set("class", outer)

			button := dom.NewElement("div")
			button.Set("class", inner)
			if id != "" {
				button.Set("id", id)
			}
			button.Append(dom.NewText(label + " ▾"))

			menu := dom.NewElement("div")
			menu.Set("class", "col none rounded border bg-panel")

			dismiss := func() {
				menu.AddClass("none")
				if a != nil {
					a.Focus(button)
					a.Redraw()
				}
			}

			// a menu left open after focus moved on would float over whatever
			// the user went to look at
			closeIfAway := func(e *Event) {
				if a == nil || menu.HasClass("none") {
					return
				}
				if wrap.Contains(a.Focused()) {
					return
				}
				menu.AddClass("none")
				a.Redraw()
			}

			button.On("blur", closeIfAway)

			for _, item := range items {
				// closing before the handler runs keeps a modal it opens
				// anchored to the button rather than to a hidden item
				pick := item.onClick
				item.onClick = func(e *Event) {
					dismiss()
					if pick != nil {
						pick(e)
					}
				}

				node := item.node(a)
				node.On("blur", closeIfAway)
				menu.Append(node)
			}

			menu.On("click", func(e *Event) { dismiss() })

			toggle := func(e *Event) {
				if !menu.HasClass("none") {
					dismiss()
					return
				}

				menu.RemoveClass("none")
				if a == nil {
					return
				}

				// an open menu owns the arrows, so it owns the focus too
				if first := firstFocusableNode(menu); first != nil {
					a.Focus(first)
				}
				a.Redraw()
			}

			button.On("click", toggle)
			button.On("keydown", func(e *Event) {
				if e.Key == "enter" || e.Rune == ' ' {
					e.StopPropagation()
					toggle(e)
				}
			})

			// esc reaches wrap from the button and from every item inside
			wrap.On("keydown", func(e *Event) {
				if e.Key == "esc" && !menu.HasClass("none") {
					e.StopPropagation()
					dismiss()
				}
			})

			wrap.Append(button)
			wrap.Append(menu)

			return wrap
		},
	}
}

func MenuItem(class, label string, onClick func(*Event)) Component {
	return Component{
		class:   join("px-1 nowrap focusable focus:bold hover:bold", class),
		text:    label,
		onClick: onClick,
	}
}
