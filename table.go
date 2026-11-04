package gtui

import (
	"strings"

	"github.com/orayew2002/gtui/src/render"
)

type TableOptions struct {
	Class   string
	Title   string
	Caption string
	Actions []Component

	Cols []Col
	Rows [][]Component

	Footer []Component
	Empty  string

	Zebra      bool
	Compact    bool
	Selectable bool
	OnSelect   func(index int)

	RowActions []RowAction
	RowWidth   string

	Page string
	Prev func(*Event)
	Next func(*Event)
}

type RowAction struct {
	Label string
	Class string
	Run   func(index int)
}

func rowActions(index int, opts TableOptions) Component {
	buttons := make([]Component, 0, len(opts.RowActions))

	for _, action := range opts.RowActions {
		at, run := index, action.Run

		class := action.Class
		if class == "" {
			class = "text-muted"
		}

		buttons = append(buttons, Button(class, action.Label,
			func(e *Event) {
				e.StopPropagation()
				if run != nil {
					run(at)
				}
			},
		))
	}

	return Div(opts.RowWidth+" shrink-0 row gap-1 justify-end", buttons...)
}

func rowActionsWidth(opts TableOptions) string {
	if opts.RowWidth != "" {
		return opts.RowWidth
	}

	width := 0
	for i, action := range opts.RowActions {
		if i > 0 {
			width++
		}
		width += render.StringWidth(action.Label) + 4
	}

	return "w-" + itoa(width)
}

func DataTable(opts TableOptions) Component {
	if opts.OnSelect != nil {
		opts.Selectable = true
	}

	pad := "px-2"
	if opts.Compact {
		pad = "px-1"
	}

	body := []Component{}

	if opts.Title != "" || len(opts.Actions) > 0 || opts.Caption != "" {
		body = append(body, tableToolbar(pad, opts))
	}

	cols := opts.Cols
	if len(opts.RowActions) > 0 {
		opts.RowWidth = rowActionsWidth(opts)
		cols = append(append([]Col{}, cols...), Col{Width: opts.RowWidth, Right: true})
		opts.Cols = cols
	}

	body = append(body, tableHeader(pad, cols))

	switch {
	case len(opts.Rows) == 0:
		body = append(body, tableEmpty(pad, opts))
	default:
		for i, row := range opts.Rows {
			if len(opts.RowActions) > 0 {
				row = append(padCells(row, len(cols)-1), rowActions(i, opts))
			}
			body = append(body, tableBodyRow(pad, i, row, opts))
		}
	}

	if len(opts.Footer) > 0 {
		body = append(body, tableFooter(pad, opts))
	}

	if opts.Page != "" || opts.Prev != nil || opts.Next != nil {
		body = append(body, Pagination(pad+" text-muted border-t border-line", opts.Page, opts.Prev, opts.Next))
	}

	return Div(join("col rounded border-line bg-panel", opts.Class), body...)
}

func padCells(row []Component, n int) []Component {
	out := make([]Component, 0, n+1)
	out = append(out, row...)

	for len(out) < n {
		out = append(out, Div(""))
	}
	return out
}

func tableToolbar(pad string, opts TableOptions) Component {
	ink := inkOf(opts.Class)
	left := []Component{}

	if opts.Title != "" {
		left = append(left, Text("bold nowrap", opts.Title))
	}
	if opts.Caption != "" {
		left = append(left, Text(join("text-muted dim nowrap", ink), opts.Caption))
	}

	head := []Component{Div("col grow", left...)}
	head = append(head, opts.Actions...)

	return Div("row gap-2 items-center border-b border-line "+pad, head...)
}

func tableHeader(pad string, cols []Col) Component {
	cells := make([]Component, 0, len(cols))

	for _, col := range cols {
		class := col.Width + " shrink-0 text-muted dim nowrap"
		if col.Right {
			class += " text-right"
		}
		cells = append(cells, Text(join(class, col.Class), strings.ToUpper(col.Label)))
	}

	return Div("row gap-1 border-b border-line "+pad, cells...)
}

func tableBodyRow(pad string, index int, row []Component, opts TableOptions) Component {
	cells := tableCells(opts.Cols, row, "")

	class := "row gap-1 " + pad + " hover:bg-raised"

	if opts.Zebra && index%2 == 1 {
		class += " bg-raised"
	}

	if opts.Selectable {
		class += " focusable focus:bg-raised focus:bold current:bg-raised current:bold"
	}

	component := Div(class, cells...)

	if opts.Selectable && opts.OnSelect != nil {
		at := index
		component.onClick = func(e *Event) { opts.OnSelect(at) }
	}

	return component
}

func tableFooter(pad string, opts TableOptions) Component {
	cells := tableCells(opts.Cols, opts.Footer, "bold")

	return Div("row gap-1 border-t border-line "+pad, cells...)
}

func tableCells(cols []Col, row []Component, extra string) []Component {
	cells := make([]Component, 0, len(cols))

	for i, col := range cols {
		cell := Div("")
		if i < len(row) {
			cell = row[i]
		}

		class := join(col.Width+" shrink-0 nowrap", extra)
		if col.Right {
			class += " text-right"
		}

		cell.class = join(join(class, col.Class), cell.class)
		cells = append(cells, cell)
	}

	return cells
}

func tableEmpty(pad string, opts TableOptions) Component {
	message := opts.Empty
	if message == "" {
		message = "nothing here yet"
	}

	return Div("row center "+pad,
		Text("text-muted dim nowrap", message),
	)
}
