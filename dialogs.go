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

func (a *App) Choose(title string, options []string, onPick func(string)) *App {
	if len(options) == 0 {
		return a
	}

	a.Modal(title, Selector("w-full #"+chooseID, "choose", "filter", options...))

	if node := a.Query("#" + chooseID); node != nil {
		// a chooser that arrives collapsed hides the very thing it was opened
		// to show, so unfold it the way a click would
		node.Dispatch(&Event{Type: "click", Target: node})

		node.On("change", func(e *Event) {
			picked := e.Value

			a.CloseModal()
			if onPick != nil {
				onPick(picked)
			}
			a.Redraw()
		})
	}

	return a
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
