package gtui

type Action struct {
	Key   string
	Label string
	Run   func()
}

func (a *App) Palette(title string, actions ...Action) *App {
	if len(actions) == 0 {
		return a
	}

	labels := make([]string, 0, len(actions))
	index := map[string]func(){}

	for _, action := range actions {
		label := action.Label
		if action.Key != "" {
			label += "   " + action.Key
		}

		labels = append(labels, label)
		index[label] = action.Run
	}

	picker := Selector("w-full #gtui-palette", "palette", "type to filter", labels...)

	a.Modal(title, picker)

	if node := a.Query("#gtui-palette"); node != nil {
		node.On("change", func(e *Event) {
			run := index[e.Value]
			a.CloseModal()

			if run != nil {
				run()
			}
			a.Redraw()
		})
	}

	if first := firstFocusableNode(a.Query("#" + modalID)); first != nil {
		a.Focus(first)
	}

	return a
}
