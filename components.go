package gtui

import (
	"strings"
	"sync/atomic"
	"unicode/utf8"

	"github.com/orayew2002/gtui/src/dom"
)

type Col struct {
	Label string
	Width string
	Class string
	Right bool
}

func Clickable(class string, onClick func(*Event), child ...Component) Component {
	return Component{
		class:   join("row nowrap focusable focus:bold focus:underline", class),
		child:   child,
		onClick: onClick,
	}
}

func Button(class, label string, onClick func(*Event)) Component {
	button := Clickable(join("px-2", class), onClick)
	button.text = label
	return button
}

func IconButton(class, icon, label string, onClick func(*Event)) Component {
	parts := make([]Component, 0, 2)
	parts = append(parts, Icon("", icon))

	if label != "" {
		parts = append(parts, Text("nowrap", label))
	}

	return Clickable(join("px-2 gap-1", class), onClick, parts...)
}

func Badge(class, label string) Component {
	return Text(join("px-1 nowrap", class), label)
}

func Pill(class, label string) Component {
	return Text(join("nowrap", class), "• "+label)
}

func Alert(class, kind, title, body string) Component {
	tone, mark := "info", "i"

	switch kind {
	case "success":
		tone, mark = "success", "✓"
	case "warning":
		tone, mark = "warning", "!"
	case "danger":
		tone, mark = "danger", "✗"
	}

	return Div(join("row rounded px-2 gap-1 border-"+tone, class),
		Text("text-"+tone+" bold nowrap", mark),
		Div("col grow",
			Text("text-"+tone+" bold nowrap", title),
			Text("", body),
		),
	)
}

func Card(class, title string, body ...Component) Component {
	head := []Component{}
	if title != "" {
		head = append(head, Text("bold nowrap", title))
	}

	return Div(join("col rounded px-2 gap-0", class), append(head, body...)...)
}

func Divider(class string) Component {
	return Div(join("h-1 border-t w-full", class))
}

func Avatar(class, name string) Component {
	initials := make([]rune, 0, 2)

	for _, word := range strings.Fields(name) {
		first, _ := utf8.DecodeRuneInString(word)
		if first == utf8.RuneError {
			continue
		}

		initials = append(initials, first)
		if len(initials) == 2 {
			break
		}
	}

	if len(initials) == 0 {
		return Text(join("px-1 nowrap bold", class), "?")
	}

	return Text(join("px-1 nowrap bold", class), strings.ToUpper(string(initials)))
}

func List(class string, items ...Component) Component {
	return Div(join("col gap-0", class), items...)
}

func ListRow(class, label string, value ...Component) Component {
	cells := make([]Component, 0, len(value)+1)
	cells = append(cells, Text("grow nowrap", label))
	cells = append(cells, value...)

	return Div(join("row gap-1", class), cells...)
}

func ListItem(class, label, value string) Component {
	return ListRow(class, label, Text("nowrap", value))
}

func Breadcrumb(class string, parts ...string) Component {
	cells := make([]Component, 0, len(parts)*2)

	for i, part := range parts {
		if i > 0 {
			cells = append(cells, Text("dim nowrap", "›"))
		}

		ink := "dim nowrap"
		if i == len(parts)-1 {
			ink = "bold nowrap"
		}
		cells = append(cells, Text(ink, part))
	}

	return Div(join("row gap-1", class), cells...)
}

type Tab struct {
	Key    string
	Icon   string
	Label  string
	Page   []Component
	Active bool
}

func Tabs(class, target string, tabs ...Tab) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)

			bar := dom.NewElement("div")
			bar.Set("class", join("row gap-1 border-b", base))
			if id != "" {
				bar.Set("id", id)
			}

			active := activeTab(tabs)

			for i, tab := range tabs {
				node := tabHead(tab, target).node(a)
				if i == active {
					node.State.Current = true
				}
				bar.Append(node)
			}

			if a != nil && target != "" && active < len(tabs) {
				if panel := a.Query(target); panel != nil && len(panel.Children) == 0 {
					a.Open(target, tabs[active].Page...)
				}
			}

			return bar
		},
	}
}

func TabView(class string, tabs ...Tab) Component {
	base, id := splitClass(class)
	target := panelTarget(id)

	return Component{
		class: class,
		build: func(a *App) *Node {
			wrap := dom.NewElement("div")
			wrap.Set("class", join("col", base))
			if id != "" {
				wrap.Set("id", id)
			}

			panel := dom.NewElement("div")
			panel.Set("class", "col grow")
			panel.Set("id", strings.TrimPrefix(target, "#"))

			if active := activeTab(tabs); active < len(tabs) {
				for _, page := range tabs[active].Page {
					panel.Append(page.node(a))
				}
			}

			wrap.Append(Tabs("", target, tabs...).node(a))
			wrap.Append(panel)

			return wrap
		},
	}
}

func tabHead(tab Tab, target string) Component {
	class := "px-2 nowrap focus:bold current:bold current:underline"
	if tab.Key != "" {
		class += " key-" + tab.Key + " key-hide"
	}

	label := tab.Label
	if glyph, ok := IconGlyph(tab.Icon); ok {
		label = glyph + " " + label
	}

	return Href(class, label, target, tab.Page...)
}

func activeTab(tabs []Tab) int {
	for i, tab := range tabs {
		if tab.Active {
			return i
		}
	}
	return 0
}

func Pagination(class, info string, prev, next func(*Event)) Component {
	items := []Component{
		Text("nowrap", info),
		Div("grow"),
	}

	if prev != nil {
		items = append(items, Button("dim", "‹ prev", prev))
	}
	if next != nil {
		items = append(items, Button("dim", "next ›", next))
	}

	return Div(join("row gap-2", class), items...)
}

func Table(class string, cols []Col, rows ...[]Component) Component {
	body := make([]Component, 0, len(rows)+1)
	body = append(body, TableHead("", cols))

	for _, row := range rows {
		body = append(body, TableRow("", cols, row...))
	}

	return Div(join("col", class), body...)
}

func TableHead(class string, cols []Col) Component {
	cells := make([]Component, 0, len(cols))

	for _, col := range cols {
		cell := col.Width + " shrink-0 dim nowrap"
		if col.Right {
			cell += " text-right"
		}
		cells = append(cells, Text(join(cell, col.Class), strings.ToUpper(col.Label)))
	}

	return Div(join("row gap-1 px-1 border-b", class), cells...)
}

func TableRow(class string, cols []Col, cells ...Component) Component {
	out := make([]Component, 0, len(cols))

	for i, col := range cols {
		cell := Div("")
		if i < len(cells) {
			cell = cells[i]
		}

		own := col.Width + " shrink-0 nowrap"
		if col.Right {
			own += " text-right"
		}

		cell.class = join(join(own, col.Class), cell.class)
		out = append(out, cell)
	}

	return Div(join("row gap-1 px-1", class), out...)
}

func Cells(values ...string) []Component {
	out := make([]Component, 0, len(values))
	for _, value := range values {
		out = append(out, Text("", value))
	}
	return out
}

func Toggle(class string, on bool) Component {
	if on {
		return Text(join("nowrap text-success", class), "●──")
	}
	return Text(join("nowrap dim", class), "──○")
}

func Spinner(class string, frame int) Component {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return Text(join("nowrap", class), frames[((frame%len(frames))+len(frames))%len(frames)])
}

func Grid(class string, columns int, items ...Component) Component {
	if columns < 1 {
		columns = 1
	}

	rows := make([]Component, 0, (len(items)+columns-1)/columns)

	for start := 0; start < len(items); start += columns {
		cells := make([]Component, 0, columns)

		for i := start; i < start+columns; i++ {
			if i < len(items) {
				cells = append(cells, Div("grow basis-0", items[i]))
				continue
			}
			cells = append(cells, Div("grow basis-0"))
		}

		rows = append(rows, Div("row gap-1", cells...))
	}

	return Div(join("col gap-1", class), rows...)
}

var panelSeq atomic.Uint64

func panelTarget(id string) string {
	if id != "" {
		return "#" + id + "-panel"
	}
	return "#gtui-panel-" + itoa(int(panelSeq.Add(1)))
}

func join(base, extra string) string {
	switch {
	case base == "":
		return extra
	case extra == "":
		return base
	}
	return base + " " + extra
}

func inkOf(class string) string {
	ink := ""

	for _, token := range strings.Fields(class) {
		switch token {
		case "text-left", "text-center", "text-right":
			continue
		}
		if strings.HasPrefix(token, "text-") {
			ink = token
		}
	}

	return ink
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}

	digits := []byte{}
	for v > 0 {
		digits = append([]byte{byte('0' + v%10)}, digits...)
		v /= 10
	}
	return string(digits)
}
