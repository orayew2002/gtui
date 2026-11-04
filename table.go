package gtui

import "strings"

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

		buttons = append(buttons, Button(class+" focus:bold focus:underline", action.Label,
			func(e *Event) {
				e.StopPropagation()
				if run != nil {
					run(at)
				}
			},
		))
	}

	width := opts.RowWidth
	if width == "" {
		width = "w-14"
	}

	return Div(width+" shrink-0 row gap-1 justify-end", buttons...)
}

func DataTable(opts TableOptions) Component {
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
		width := opts.RowWidth
		if width == "" {
			width = "w-14"
		}
		cols = append(append([]Col{}, cols...), Col{Label: "", Width: width, Right: true})
		opts.Cols = cols
	}

	body = append(body, tableHeader(pad, cols))

	switch {
	case len(opts.Rows) == 0:
		body = append(body, tableEmpty(pad, opts))
	default:
		for i, row := range opts.Rows {
			if len(opts.RowActions) > 0 {
				row = append(append([]Component{}, row...), rowActions(i, opts))
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

func tableToolbar(pad string, opts TableOptions) Component {
	left := []Component{}

	if opts.Title != "" {
		left = append(left, Text("text-ink bold nowrap", opts.Title))
	}
	if opts.Caption != "" {
		left = append(left, Text("text-muted dim nowrap", opts.Caption))
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
		cells = append(cells, Text(class, strings.ToUpper(col.Label)))
	}

	return Div("row gap-1 border-b border-line "+pad, cells...)
}

func tableBodyRow(pad string, index int, row []Component, opts TableOptions) Component {
	cells := make([]Component, 0, len(opts.Cols))

	for i, col := range opts.Cols {
		if i >= len(row) {
			break
		}

		class := col.Width + " shrink-0 nowrap"
		if col.Right {
			class += " text-right"
		}

		cell := row[i]
		cell.class = join(class, cell.class)
		cells = append(cells, cell)
	}

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
	cells := make([]Component, 0, len(opts.Cols))

	for i, col := range opts.Cols {
		if i >= len(opts.Footer) {
			break
		}

		class := col.Width + " shrink-0 nowrap bold"
		if col.Right {
			class += " text-right"
		}

		cell := opts.Footer[i]
		cell.class = join(class, cell.class)
		cells = append(cells, cell)
	}

	return Div("row gap-1 border-t border-line "+pad, cells...)
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
