package gtui

import (
	"strings"

	"github.com/orayew2002/gtui/src/dom"
)

type Section struct {
	Key   string
	Title string
	Body  []Component
	Open  bool
}

func Accordion(class string, sections ...Section) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)

			wrap := dom.NewElement("div")
			wrap.Set("class", join("col", base))
			if id != "" {
				wrap.Set("id", id)
			}

			for _, section := range sections {
				panel := dom.NewElement("div")
				panel.Set("class", "col")

				body := dom.NewElement("div")
				body.Set("class", "col px-2")
				if !section.Open {
					body.AddClass("none")
				}
				for _, item := range section.Body {
					body.Append(item.node(a))
				}

				head := dom.NewElement("div")
				head.Set("class", "row gap-1 focusable focus:bold border-b border-line")
				if section.Key != "" {
					head.AddClass("key-" + section.Key)
				}

				title := section.Title
				render := func() {
					head.Clear()

					mark := dom.NewElement("div")
					mark.Set("class", "nowrap dim")
					if body.HasClass("none") {
						mark.Append(dom.NewText("▸"))
					} else {
						mark.Append(dom.NewText("▾"))
					}

					text := dom.NewElement("div")
					text.Set("class", "grow nowrap")
					text.Append(dom.NewText(title))

					head.Append(mark)
					head.Append(text)
				}
				render()

				toggle := func() {
					if body.HasClass("none") {
						body.RemoveClass("none")
					} else {
						body.AddClass("none")
					}
					render()
					if a != nil {
						a.Redraw()
					}
				}

				head.On("click", func(e *Event) { toggle() })
				head.On("keydown", func(e *Event) {
					if e.Key == "enter" || e.Rune == ' ' {
						e.StopPropagation()
						toggle()
					}
				})

				panel.Append(head)
				panel.Append(body)
				wrap.Append(panel)
			}

			return wrap
		},
	}
}

type TreeNode struct {
	Label    string
	Value    string
	Icon     string
	Open     bool
	Children []TreeNode
}

type TreeOptions struct {
	Class    string
	Roots    []TreeNode
	OnSelect func(TreeNode)
}

func Tree(class string, roots ...TreeNode) Component {
	return TreeView(TreeOptions{Class: class, Roots: roots})
}

func TreeView(opts TreeOptions) Component {
	return Component{
		class: opts.Class,
		build: func(a *App) *Node {
			base, id := splitClass(opts.Class)

			wrap := dom.NewElement("div")
			wrap.Set("class", join("col", base))
			if id != "" {
				wrap.Set("id", id)
				wrap.Set("name", id)
			}

			var (
				rows  []*Node
				mount func(parent *Node, items []TreeNode, depth int)
			)

			choose := func(row *Node, item TreeNode, value string) {
				for _, other := range rows {
					other.State.Current = other == row
				}

				wrap.Set("value", value)
				row.Dispatch(&dom.Event{Type: "select", Target: row, Value: value})

				if opts.OnSelect != nil {
					opts.OnSelect(item)
				}
			}

			mount = func(parent *Node, items []TreeNode, depth int) {
				for _, item := range items {
					branch := item

					value := branch.Value
					if value == "" {
						value = branch.Label
					}

					row := dom.NewElement("div")
					row.Set("class", "row gap-1 focusable focus:bold current:bold")
					if depth > 0 {
						row.AddClass("pl-" + itoa(depth*2))
					}
					row.Set("value", value)
					if id != "" && branch.Value != "" {
						row.Set("id", id+"-"+branch.Value)
					}
					rows = append(rows, row)

					children := dom.NewElement("div")
					children.Set("class", "col")
					if !branch.Open {
						children.AddClass("none")
					}

					glyph, hasIcon := IconGlyph(branch.Icon)
					leaf := len(branch.Children) == 0

					mark := dom.NewElement("div")
					mark.Set("class", "nowrap dim")

					label := dom.NewElement("div")
					label.Set("class", "grow nowrap")
					label.Append(dom.NewText(branch.Label))

					render := func() {
						mark.Clear()
						switch {
						case leaf && hasIcon:
							mark.Append(dom.NewText(glyph))
						case leaf:
							mark.Append(dom.NewText("·"))
						case children.HasClass("none"):
							mark.Append(dom.NewText("▸"))
						default:
							mark.Append(dom.NewText("▾"))
						}
					}
					render()

					activate := func() {
						switch {
						case leaf:
							choose(row, branch, value)
						case children.HasClass("none"):
							children.RemoveClass("none")
							render()
						default:
							children.AddClass("none")
							render()
						}
						if a != nil {
							a.Redraw()
						}
					}

					row.On("click", func(e *Event) { activate() })
					row.On("keydown", func(e *Event) {
						if e.Key == "enter" || e.Rune == ' ' {
							e.StopPropagation()
							activate()
						}
					})

					row.Append(mark)
					if hasIcon && !leaf {
						icon := dom.NewElement("div")
						icon.Set("class", "nowrap")
						icon.Append(dom.NewText(glyph))
						row.Append(icon)
					}
					row.Append(label)

					parent.Append(row)
					parent.Append(children)

					mount(children, branch.Children, depth+1)
				}
			}

			mount(wrap, opts.Roots, 0)
			return wrap
		},
	}
}

func Stepper(class string, step int, labels ...string) Component {
	cells := make([]Component, 0, len(labels)*2)

	for i, label := range labels {
		if i > 0 {
			cells = append(cells, Text("dim nowrap", "──"))
		}

		mark, ink := "○", "dim nowrap"
		switch {
		case i < step:
			mark, ink = "●", "text-success nowrap"
		case i == step:
			mark, ink = "◉", "text-accent bold nowrap"
		}

		cells = append(cells, Text(ink, mark+" "+label))
	}

	return Div(join("row gap-1 items-center", class), cells...)
}

func Rating(class string, value, max int) Component {
	if max < 1 {
		max = 5
	}

	return Text(join("nowrap text-warning", class),
		strings.Repeat("★", clampInt(value, 0, max))+strings.Repeat("☆", max-clampInt(value, 0, max)))
}

func Skeleton(class string, rows, width int) Component {
	if rows < 1 {
		rows = 3
	}
	if width < 1 {
		width = 20
	}

	lines := make([]Component, 0, rows)
	for i := range rows {
		length := width
		if i == rows-1 {
			length = width * 2 / 3
		}
		lines = append(lines, Text("dim nowrap", strings.Repeat("▒", length)))
	}

	return Div(join("col gap-0", class), lines...)
}

func Chips(class string, labels ...string) Component {
	items := make([]Component, 0, len(labels))
	for _, label := range labels {
		items = append(items, Text("px-1 nowrap bg-raised", label))
	}

	return Div(join("row flex-wrap gap-1", class), items...)
}

func StatCard(class, label, value, note string) Component {
	ink := inkOf(class)

	return Div(join("col rounded border-line bg-panel px-2", class),
		Text(join("text-muted dim nowrap", ink), strings.ToUpper(label)),
		Text("bold nowrap", value),
		Text(join("text-muted nowrap", ink), note),
	)
}

func KeyValue(class string, pairs ...[2]string) Component {
	ink := inkOf(class)
	rows := make([]Component, 0, len(pairs))

	for _, pair := range pairs {
		rows = append(rows, Div("row gap-1",
			Text(join("grow text-muted nowrap", ink), pair[0]),
			Text("nowrap", pair[1]),
		))
	}

	return Div(join("col gap-0", class), rows...)
}

func Split(class, leftWidth string, left, right Component) Component {
	if leftWidth == "" {
		leftWidth = "w-1/3"
	}

	return Div(join("row gap-1", class),
		Div(leftWidth+" shrink-0", left),
		Div("grow", right),
	)
}

func StatusBar(class string, left, right []Component) Component {
	items := make([]Component, 0, len(left)+len(right)+1)
	items = append(items, left...)
	items = append(items, Spacer())
	items = append(items, right...)

	return Div(join("row gap-2 px-1 text-muted", class), items...)
}

func Hint(class, key, label string) Component {
	ink := inkOf(class)

	return Div(join("row gap-1 nowrap", class),
		Text(join("bold text-accent nowrap", ink), key),
		Text(join("text-muted nowrap", ink), label),
	)
}
