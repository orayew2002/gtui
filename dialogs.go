package gtui

import "sync"

const (
	dialogPrimary   = "rounded border-accent text-accent focus:bold"
	dialogSecondary = "rounded border-line text-muted focus:bold"
)

const (
	dialogYesID    = "gtui-dialog-yes"
	dialogOKID     = "gtui-dialog-ok"
	promptInputID  = "gtui-prompt-input"
	chooseID       = "gtui-choose"
	chooseFirstID  = "gtui-choose-first"
	busyGlyphID    = "gtui-busy-glyph"
	busyGlyphSel   = "#" + busyGlyphID
	busyFramesPerS = 10
)

// Confirm asks a yes/cancel question. Yes takes the focus the modal hands out,
// so enter answers it and esc cancels, the way every other modal closes.
func (a *App) Confirm(title, body string, onYes func()) *App {
	return a.Modal(title,
		Text("", body),
		Div("row gap-1 justify-end pt-1",
			Button(dialogPrimary+" #"+dialogYesID, "yes", func(e *Event) {
				a.CloseModal()
				if onYes != nil {
					onYes()
				}
			}),
			Button(dialogSecondary, "cancel", func(e *Event) { a.CloseModal() }),
		),
	)
}

func (a *App) Alert(title, body string) *App {
	return a.Modal(title,
		Text("", body),
		Div("row gap-1 justify-end pt-1",
			Button(dialogPrimary+" #"+dialogOKID, "ok", func(e *Event) { a.CloseModal() }),
		),
	)
}

func (a *App) Prompt(title, label, placeholder string, onSubmit func(value string)) *App {
	submit := func() {
		value := a.Value("#" + promptInputID)

		a.CloseModal()
		if onSubmit != nil {
			onSubmit(value)
		}
	}

	a.Modal(title,
		Field("gap-0", label,
			Input("w-full rounded border-line focus:border-accent #"+promptInputID, "prompt", placeholder),
		),
		Div("row gap-1 justify-end pt-1",
			Button(dialogPrimary, "ok", func(e *Event) { submit() }),
			Button(dialogSecondary, "cancel", func(e *Event) { a.CloseModal() }),
		),
	)

	// a one-line field leaves enter alone, so the dialog claims it as submit
	if node := a.Query("#" + promptInputID); node != nil {
		node.On("keydown", func(e *Event) {
			if e.Key == "enter" {
				e.StopPropagation()
				submit()
			}
		})
	}

	return a
}

// Choose asks for one of a handful of options. The options are rows rather than
// a dropdown: a dropdown inside a dialog would take the first esc to fold
// itself away, leaving a dialog with nothing in it, and the person pressing esc
// means cancel.
func (a *App) Choose(title string, options []string, onPick func(string)) *App {
	if len(options) == 0 {
		return a
	}

	rows := make([]Component, 0, len(options))
	for i, option := range options {
		pick := option

		rows = append(rows, Clickable(chooseRowClass(i), func(e *Event) {
			e.StopPropagation()

			a.CloseModal()
			if onPick != nil {
				onPick(pick)
			}
			a.Redraw()
		}, Text("grow nowrap", option)))
	}

	a.Modal(title, Div("col gap-0 #"+chooseID, rows...))

	a.walkChooseRows(func(index int, row *Node) {
		row.On("keydown", func(e *Event) {
			switch e.Key {
			case "down":
				e.StopPropagation()
				a.stepChoose(index + 1)
			case "up":
				e.StopPropagation()
				a.stepChoose(index - 1)
			}
		})
	})

	return a
}

func chooseRowClass(index int) string {
	class := "row gap-1 px-1 focusable hover:bg-raised focus:bg-raised focus:bold"
	if index == 0 {
		class += " #" + chooseFirstID
	}
	return class
}

// walkChooseRows hands every option row to fn, in the order they were added.
func (a *App) walkChooseRows(fn func(index int, row *Node)) {
	list := a.Query("#" + chooseID)
	if list == nil {
		return
	}

	index := 0
	for _, row := range list.ElementChildren() {
		fn(index, row)
		index++
	}
}

// stepChoose moves the keyboard between option rows, so the arrows work the way
// they do in a list.
func (a *App) stepChoose(to int) {
	var rows []*Node
	a.walkChooseRows(func(_ int, row *Node) { rows = append(rows, row) })

	if len(rows) == 0 {
		return
	}

	a.Focus(rows[clampInt(to, 0, len(rows)-1)])
	a.Redraw()
}

// Busy shows a spinning modal and returns the func that stops it. The closer is
// safe to call twice and from another goroutine, and it leaves a modal that has
// since replaced this one alone.
func (a *App) Busy(label string) func() {
	a.Modal(label, Center("w-full py-1", Spinner("text-accent #"+busyGlyphID, 0)))

	stop := a.Animate(busyFramesPerS, func(tick int) {
		if a.Query(busyGlyphSel) == nil {
			return
		}
		// Spinner owns the frame table, so ask it for the glyph
		a.Text(busyGlyphSel, Spinner("", tick).text)
	})

	var once sync.Once

	return func() {
		once.Do(func() {
			stop()

			if a.Query(busyGlyphSel) != nil {
				a.CloseModal()
			}
		})
	}
}
