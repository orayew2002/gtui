package gtui

import "strings"

type Col struct {
	Label string
	Width string
	Right bool
}

func Button(class, label string, onClick func(*Event)) Component {
	return Component{
		class:   join("row px-2 nowrap focusable", class),
		text:    label,
		onClick: onClick,
	}
}

func Badge(class, label string) Component {
	return Text(join("px-1 nowrap", class), label)
}

func Pill(class, label string) Component {
	return Text(join("nowrap", class), "• "+label)
}

func Alert(class, kind, title, body string) Component {
	ink, mark := "text-sky", "i"

	switch kind {
	case "success":
		ink, mark = "text-green", "✓"
	case "warning":
		ink, mark = "text-yellow", "!"
	case "danger":
		ink, mark = "text-red", "✗"
	}

	return Div(join("row rounded px-2 gap-1 border-"+strings.TrimPrefix(ink, "text-"), class),
		Text(ink+" bold nowrap", mark),
		Div("col grow",
			Text(ink+" bold nowrap", title),
			Text("nowrap", body),
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
	initials := ""
	for _, word := range strings.Fields(name) {
		initials += strings.ToUpper(word[:1])
		if len(initials) == 2 {
			break
		}
	}
	if initials == "" {
		initials = "?"
	}

	return Text(join("px-1 nowrap bold", class), initials)
}

func List(class string, items ...Component) Component {
	return Div(join("col gap-0", class), items...)
}

func ListItem(class, label, value string) Component {
	return Div(join("row gap-1", class),
		Text("grow nowrap", label),
		Text("nowrap", value),
	)
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
	Key   string
	Label string
	Page  []Component
}

func Tabs(class, target string, tabs ...Tab) Component {
	items := make([]Component, 0, len(tabs))

	for _, tab := range tabs {
		item := "px-2 nowrap focus:bold current:bold current:underline"
		if tab.Key != "" {
			item += " key-" + tab.Key + " key-hide"
		}
		items = append(items, Href(item, tab.Label, target, tab.Page...))
	}

	return Div(join("row gap-1 border-b", class), items...)
}

func Pagination(class, info string, prev, next func(*Event)) Component {
	items := []Component{
		Text("nowrap", info),
		Div("grow"),
	}

	if prev != nil {
		items = append(items, Button("dim focus:bold", "‹ prev", prev))
	}
	if next != nil {
		items = append(items, Button("dim focus:bold", "next ›", next))
	}

	return Div(join("row gap-2", class), items...)
}

func Table(class string, cols []Col, rows ...[]Component) Component {
	body := make([]Component, 0, len(rows)+1)
	body = append(body, TableHead(cols))

	for _, row := range rows {
		body = append(body, TableRow(cols, row...))
	}

	return Div(join("col", class), body...)
}

func TableHead(cols []Col) Component {
	cells := make([]Component, 0, len(cols))

	for _, col := range cols {
		class := col.Width + " shrink-0 dim nowrap"
		if col.Right {
			class += " text-right"
		}
		cells = append(cells, Text(class, strings.ToUpper(col.Label)))
	}

	return Div("row gap-1 px-1 border-b", cells...)
}

func TableRow(cols []Col, cells ...Component) Component {
	out := make([]Component, 0, len(cols))

	for i, col := range cols {
		if i >= len(cells) {
			break
		}

		class := col.Width + " shrink-0 nowrap"
		if col.Right {
			class += " text-right"
		}

		cell := cells[i]
		cell.class = join(class, cell.class)
		out = append(out, cell)
	}

	return Div("row gap-1 px-1", out...)
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
		return Text(join("nowrap text-green", class), "●──")
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

	width := "w-" + itoa(100/columns) + "%"
	cells := make([]Component, 0, len(items))

	for _, item := range items {
		cells = append(cells, Div(width+" shrink-0", item))
	}

	return Div(join("row flex-wrap gap-1", class), cells...)
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
