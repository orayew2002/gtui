package gtui

import (
	"time"

	"github.com/orayew2002/gtui/src/dom"
)

const (
	modalID = "gtui-modal"
	toastID = "gtui-toast"
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

	layer.Clear()
	layer.Append(panel.node(a))
	layer.RemoveClass("none")

	a.modalOpen = true
	a.beforeModal = a.focus

	if first := firstFocusableNode(layer); first != nil {
		a.Focus(first)
	}

	return a.Redraw()
}

func (a *App) CloseModal() *App {
	layer := a.Query("#" + modalID)
	if layer == nil {
		return a
	}

	layer.Clear()
	layer.AddClass("none")
	a.modalOpen = false

	a.Focus(a.beforeModal)
	a.beforeModal = nil

	return a.Redraw()
}

func (a *App) ModalOpen() bool { return a.modalOpen }

func (a *App) Toast(kind, message string) *App {
	layer := a.layer(toastID, "absolute top-0 right-0 z-9 px-1")

	ink := "text-sky"
	switch kind {
	case "success":
		ink = "text-green"
	case "warning":
		ink = "text-yellow"
	case "danger":
		ink = "text-red"
	}

	layer.Clear()
	layer.Append(Div("row rounded px-2 bg-panel "+ink, Text("nowrap", message)).node(a))
	layer.RemoveClass("none")

	go func() {
		time.Sleep(2 * time.Second)
		a.Do(func() {
			if node := a.Query("#" + toastID); node != nil {
				node.Clear()
				node.AddClass("none")
			}
		})
	}()

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
			base = withClass(base, "col")

			wrap := dom.NewElement("div")
			if base != "" {
				wrap.Set("class", base)
			}
			if id != "" {
				wrap.Set("id", id)
			}

			button := dom.NewElement("div")
			button.Set("class", "row gap-1 focusable nowrap")
			button.Append(dom.NewText(label + " ▾"))

			menu := dom.NewElement("div")
			menu.Set("class", "col none rounded border bg-panel")
			for _, item := range items {
				menu.Append(item.node(a))
			}

			toggle := func(e *Event) {
				if menu.HasClass("none") {
					menu.RemoveClass("none")
				} else {
					menu.AddClass("none")
				}
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
					menu.AddClass("none")
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
