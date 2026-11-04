package main

import (
	"fmt"
	"strings"

	"github.com/orayew2002/gtui"
)

var (
	iconFilter = ""
	iconPage   = 0
)

func pageIcons() gtui.Component {
	return card("a", "icons",
		gtui.Row("gap-1 items-center",
			gtui.Search("px-1 bg-raised w-12 #ic-find", "icon", "filter").On("input", filterIcons),
			gtui.Text("dim nowrap #ic-info", iconInfo()),
			gtui.IconButton("px-1 focus:bold focus:text-accent #ic-prev", "left", "", func(e *gtui.Event) {
				turnIcons(-1)
			}),
			gtui.IconButton("px-1 focus:bold focus:text-accent #ic-next", "right", "", func(e *gtui.Event) {
				turnIcons(1)
			}),
			gtui.Text("grow truncate text-accent #ic-out", "click one to copy its name"),
		),
		gtui.Div("col gap-0 #ic-grid", iconGrid()),
		gtui.Row("gap-1 items-center mt-1",
			note("RegisterIcon adds your own glyph to the same set:"),
			gtui.IconText("text-accent nowrap", "gopher", "gopher"),
			gtui.Text("dim nowrap", `RegisterIcon("gopher", "ᵹ")`),
		),
	)
}

// the gallery registers one glyph of its own, so the icon set is clearly open
func init() { gtui.RegisterIcon("gopher", "ᵹ") }

// iconGrid lays the icons out for the room there is: as many columns as fit, as
// many rows as the card has, and a pager when the names need more.
func iconGrid() gtui.Component {
	names := iconMatches()
	if len(names) == 0 {
		return gtui.Text("text-muted nowrap", "no icon matches %q", iconFilter)
	}

	columns, rows, cell := iconShape()
	size := columns * rows
	iconPage = clampPage(iconPage, len(names), size)

	page := names[iconPage*size : min(iconPage*size+size, len(names))]
	grid := make([]gtui.Component, 0, rows)

	for start := 0; start < len(page); start += columns {
		cells := make([]gtui.Component, 0, columns)
		for _, name := range page[start:min(start+columns, len(page))] {
			cells = append(cells, iconCell(name, cell))
		}
		grid = append(grid, gtui.Row("gap-0", cells...))
	}

	return gtui.Column("gap-0", grid...)
}

// iconCell is one glyph and its name in a box of a known width, so a long name
// cannot widen the grid and push the rest of the gallery around.
func iconCell(name string, cell int) gtui.Component {
	label := name
	if len(label) > cell-2 {
		label = label[:cell-3] + "…"
	}

	return gtui.Div(fmt.Sprintf("w-%d clip shrink-0", cell),
		gtui.IconText("text-muted hover:text-ink", name, label).OnClick(func(e *gtui.Event) {
			app.Copy(name)
			app.Text("#ic-out", "copied "+name)
		}),
	)
}

// iconShape splits the card body into columns that fill it: a cell wide enough
// for a glyph, a space and most of a name.
func iconShape() (columns, rows, cell int) {
	width, height := room()

	rows = height - 4
	if rows < 3 {
		rows = 3
	}

	// the card spends two cells on its border and four on its padding
	body := width - 6

	columns = max(body/12, 1)
	cell = max(body/columns, 6)

	return columns, rows, cell
}

func iconMatches() []string {
	names := gtui.IconNames()
	if iconFilter == "" {
		return names
	}

	found := make([]string, 0, len(names))
	for _, name := range names {
		if strings.Contains(strings.ToLower(name), iconFilter) {
			found = append(found, name)
		}
	}

	return found
}

func iconInfo() string {
	names := iconMatches()
	columns, rows, _ := iconShape()

	size := columns * rows
	pages := (len(names) + size - 1) / size

	if pages < 2 {
		return fmt.Sprintf("%d of %d", len(names), len(gtui.IconNames()))
	}

	return fmt.Sprintf("%d of %d · %d/%d", len(names), len(gtui.IconNames()),
		clampPage(iconPage, len(names), size)+1, pages)
}

func filterIcons(e *gtui.Event) {
	iconFilter = strings.ToLower(strings.TrimSpace(e.Value))
	iconPage = 0

	redrawIcons()
}

func turnIcons(step int) {
	columns, rows, _ := iconShape()

	iconPage = clampPage(iconPage+step, len(iconMatches()), columns*rows)
	redrawIcons()
}

// redrawIcons repaints the grid and the counter without rebuilding the page, so
// the filter field keeps both its text and the keyboard.
func redrawIcons() {
	app.Open("#ic-grid", iconGrid())
	app.Text("#ic-info", iconInfo())
}

func clampPage(page, count, size int) int {
	if size < 1 {
		return 0
	}

	last := (count - 1) / size
	if page > last {
		page = last
	}
	if page < 0 {
		page = 0
	}

	return page
}
