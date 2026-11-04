package gtui

import (
	"strings"
	"unicode/utf8"

	"github.com/orayew2002/gtui/src/dom"
)

func Navbar(class string, brand Component, items ...Component) Component {
	cells := make([]Component, 0, len(items)+2)

	if !blank(brand) {
		cells = append(cells, brand)
	}
	cells = append(cells, Spacer())
	cells = append(cells, items...)

	return Div(join("row gap-2 items-center px-2 bg-panel border-b border-line", class), cells...)
}

func ButtonGroup(class string, buttons ...Component) Component {
	cells := make([]Component, 0, len(buttons)*2)

	for i, button := range buttons {
		if i > 0 {
			cells = append(cells, Text("dim nowrap", "│"))
		}
		cells = append(cells, button)
	}

	return Div(join("row gap-0 items-center rounded border-line", class), cells...)
}

func InputGroup(class string, prefix, input, suffix Component) Component {
	cells := make([]Component, 0, 3)

	if !blank(prefix) {
		cells = append(cells, Div("px-1 shrink-0 bg-raised text-muted", prefix))
	}
	cells = append(cells, Div("grow", input))
	if !blank(suffix) {
		cells = append(cells, Div("px-1 shrink-0 bg-raised text-muted", suffix))
	}

	return Div(join("row gap-0 items-center rounded border-line", class), cells...)
}

func Banner(class, kind, title, body string, actions ...Component) Component {
	tone := toneOr(kind, "accent")
	rows := make([]Component, 0, 3)

	rows = append(rows, Text("text-"+tone+" bold nowrap", title))
	if body != "" {
		rows = append(rows, Text("text-muted", body))
	}
	if len(actions) > 0 {
		rows = append(rows, Div("row gap-1 pt-1", actions...))
	}

	return Div(join("col w-full rounded px-2 gap-0 bg-panel border-"+tone, class), rows...)
}

func EmptyState(class, icon, title, body string, actions ...Component) Component {
	rows := make([]Component, 0, 4)

	if glyph, ok := IconGlyph(icon); ok {
		rows = append(rows, Text("text-muted nowrap", glyph))
	}
	rows = append(rows, Text("bold nowrap", title))
	if body != "" {
		rows = append(rows, Text("text-muted text-center", body))
	}
	if len(actions) > 0 {
		rows = append(rows, Div("row gap-1 pt-1", actions...))
	}

	return Div(join("col w-full items-center justify-center gap-0 py-1", class), rows...)
}

func Panel(class, title string, actions []Component, body ...Component) Component {
	head := make([]Component, 0, len(actions)+1)
	head = append(head, Text("bold grow nowrap", title))
	head = append(head, actions...)

	rows := make([]Component, 0, len(body)+1)
	rows = append(rows, Div("row gap-1 items-center border-b border-line", head...))
	rows = append(rows, body...)

	return Div(join("col rounded border-line bg-panel px-2 gap-0", class), rows...)
}

type TimelineItem struct {
	Icon  string
	Title string
	Time  string
	Body  string
	Kind  string
}

func Timeline(class string, items ...TimelineItem) Component {
	rows := make([]Component, 0, len(items)*2)

	for i, item := range items {
		mark := "●"
		if glyph, ok := IconGlyph(item.Icon); ok {
			mark = glyph
		}

		rows = append(rows, Div("row gap-1",
			Text("w-1 nowrap text-"+toneOr(item.Kind, "accent"), mark),
			Text("grow bold nowrap", item.Title),
			Text("text-muted dim nowrap", item.Time),
		))

		last := i == len(items)-1
		if last && item.Body == "" {
			continue
		}

		// the rail only continues while another entry is still coming
		rail := " "
		if !last {
			rail = "│"
		}

		rows = append(rows, Div("row gap-1",
			Text("w-1 nowrap dim", rail),
			Text("grow text-muted", item.Body),
		))
	}

	return Div(join("col gap-0", class), rows...)
}

type StepItem struct {
	Title string
	Note  string
}

func Steps(class string, step int, items ...StepItem) Component {
	cells := make([]Component, 0, len(items)*2)

	for i, item := range items {
		if i > 0 {
			cells = append(cells, Text("dim nowrap self-start", "──"))
		}

		mark, ink := "○", "text-muted"
		switch {
		case i < step:
			mark, ink = "✓", "text-success"
		case i == step:
			mark, ink = "◉", "text-accent"
		}

		head := Div("row gap-1",
			Text(ink+" nowrap", mark),
			Text(ink+" bold nowrap", item.Title),
		)

		if item.Note == "" {
			cells = append(cells, head)
			continue
		}

		cells = append(cells, Column("gap-0",
			head,
			Text("pl-2 text-muted dim nowrap", item.Note),
		))
	}

	return Div(join("row gap-1 items-start", class), cells...)
}

func StatusDot(class, kind string) Component {
	return Text(join("nowrap text-"+toneOr(kind, "muted"), class), "●")
}

func Metric(class, label, value, delta string) Component {
	ink := inkOf(class)
	cells := make([]Component, 0, 2)

	cells = append(cells, Text("bold nowrap", value))

	if delta != "" {
		tone, mark := deltaTone(delta)
		if mark != "" {
			delta = mark + " " + delta
		}
		cells = append(cells, Text("text-"+tone+" nowrap", delta))
	}

	return Div(join("col rounded border-line bg-panel px-2 gap-0", class),
		Text(join("text-muted dim nowrap", ink), strings.ToUpper(label)),
		Div("row gap-1 items-center", cells...),
	)
}

func Tag(class, label string, onRemove func(*Event)) Component {
	if onRemove == nil {
		return Text(join("px-1 nowrap bg-raised", class), label)
	}

	return Div(join("row gap-1 px-1 items-center bg-raised", class),
		Text("nowrap", label),
		Button("px-0 text-muted hover:text-danger focus:text-danger", "×", onRemove),
	)
}

func Kbd(class, keys string) Component {
	return Text(join("pill nowrap bold border-line bg-raised", class), keys)
}

func Loader(class, label string, frame int) Component {
	ink := inkOf(class)

	return Div(join("row gap-1 items-center", class),
		Spinner(join("text-accent", ink), frame),
		Text("grow nowrap", label),
	)
}

func Separator(class, label string) Component {
	if label == "" {
		return Divider(class)
	}

	return Div(join("row gap-1 items-center w-full", class),
		Div("grow h-1 border-t border-line"),
		Text("text-muted dim nowrap", label),
		Div("grow h-1 border-t border-line"),
	)
}

func Collapse(class, title string, open bool, body ...Component) Component {
	return Accordion(class, Section{Title: title, Body: body, Open: open})
}

func Scrollbox(class string, height int, child ...Component) Component {
	if height < 1 {
		height = 10
	}

	return Div(join("col scroll gap-0 h-"+itoa(height), class), child...)
}

func Centered(class string, child ...Component) Component {
	return Div(join("center full", class), child...)
}

func Window(class, title string, body ...Component) Component {
	return Component{
		class: class,
		build: func(a *App) *Node {
			base, id := splitClass(class)

			n := dom.NewElement("div")
			n.Set("class", join("col rounded border-line bg-panel px-1 gap-0", base))
			if id != "" {
				n.Set("id", id)
			}
			// the title rides the top border, which the paint pass draws
			if title != "" {
				n.Set("title", title)
			}

			for _, item := range body {
				n.Append(item.node(a))
			}

			return n
		},
	}
}

func Columns(class string, widths []string, cells ...Component) Component {
	out := make([]Component, 0, len(cells))

	for i, cell := range cells {
		width := "grow basis-0"
		if i < len(widths) && widths[i] != "" {
			width = join(widths[i], "shrink-0")
		}
		out = append(out, Div(width, cell))
	}

	return Div(join("row gap-1 items-start", class), out...)
}

// blank reports the zero Component, which callers pass for a slot they have no
// content for. Dropping it keeps a row's gap off an empty cell.
func blank(c Component) bool {
	return c.class == "" && c.text == "" && c.build == nil && !c.link && len(c.child) == 0
}

// toneOr maps a component kind onto a theme token, so "ok", "down" and the
// token names themselves all pick a colour.
func toneOr(kind, fallback string) string {
	switch kind {
	case "accent", "success", "warning", "danger", "info", "muted":
		return kind
	case "ok", "on", "up", "done", "live":
		return "success"
	case "warn", "busy", "pending", "slow":
		return "warning"
	case "error", "fail", "down", "offline":
		return "danger"
	case "off", "idle", "none", "draft":
		return "muted"
	}
	return fallback
}

func deltaTone(delta string) (tone, mark string) {
	switch first, _ := utf8.DecodeRuneInString(delta); first {
	case '+':
		return "success", "↑"
	case '-', '−':
		return "danger", "↓"
	}
	return "muted", ""
}
