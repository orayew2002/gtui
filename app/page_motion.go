package main

import "github.com/orayew2002/gtui"

// frame counts animation ticks for every page that wants movement.
var frame = 0

type liveBox struct {
	id   string
	rows func() []gtui.Component
}

// liveBoxes lists the small boxes the tick redraws on the page that is open.
// Redrawing one box instead of the whole page keeps the tick cheap and leaves
// the node that holds keyboard focus where it was.
func liveBoxes() []liveBox {
	switch current().Title {
	case "motion":
		return []liveBox{{"#mo-text", motionText}, {"#mo-meters", motionMeters}}
	case "media":
		return []liveBox{{"#mv-frame", videoFrame}, {"#mp-meter", musicMeter}}
	case "feedback":
		return []liveBox{{"#fb-spin", spinnerRow}}
	}
	return nil
}

func onTick(tick int) {
	frame = tick

	if playing && tick%2 == 0 {
		beat = (beat + 1) % reel.Count()
	}

	for _, box := range liveBoxes() {
		if app.Query(box.id) == nil {
			continue
		}
		app.Open(box.id, box.rows()...)
	}

	liveThemeName()
}

func pageMotion() gtui.Component {
	rows := []gtui.Component{
		gtui.Row("gap-1",
			card("a", "text in motion",
				note("Marquee, Typewriter, Blink"),
				gtui.Column("gap-0 #mo-text", motionText()...),
			),
			card("b", "meters",
				note("Pulse, Wave, Dots, Spinner"),
				gtui.Column("gap-0 #mo-meters", motionMeters()...),
			),
		),
		gtui.Row("gap-1", card("d", "how it works", code()...)),
	}

	return gtui.Column("grow gap-1", rows...)
}

func code() []gtui.Component {
	lines := []string{
		`app.Animate(8, func(tick int) {`,
		`    app.Open("#meters", meters(tick))`,
	}

	body := []gtui.Component{}
	if wideScreen() {
		lines = append(lines, `})`)
		body = append(body, note("the tick redraws one box, not the page, so focus stays put"))
	}

	body = append(body, gtui.CodeView(gtui.CodeOptions{Lines: lines, First: 1}))
	if wideScreen() {
		body = append(body, gtui.Hint("", "app.Every", "slower timers, app.After for one shot"))
	}

	return body
}

func motionText() []gtui.Component {
	width, _ := room()
	band := width/2 - 8

	return []gtui.Component{
		gtui.Marquee("text-accent", "gtui writes only the cells that changed, so motion is cheap", band, frame),
		gtui.Typewriter("text-success", "typed one rune at a time", frame/2),
		gtui.Blink("text-danger", "● recording", frame/4),
	}
}

func motionMeters() []gtui.Component {
	width, _ := room()
	bars := width/2 - 10

	return []gtui.Component{
		gtui.Pulse("text-warning", bars, frame),
		gtui.Wave("text-accent", bars, frame),
		gtui.Equalizer("text-success", bars, frame),
		gtui.Row("gap-2 items-center",
			gtui.Spinner("text-accent", frame),
			gtui.Dots("text-muted", frame),
			gtui.Text("text-muted nowrap", "tick %d", frame),
		),
	}
}

// spinnerRow feeds the feedback page's #fb-spin box, so that page can move
// without the tick rebuilding anything that holds focus.
func spinnerRow() []gtui.Component {
	return []gtui.Component{
		gtui.Row("gap-1 items-center",
			gtui.Spinner("text-accent", frame),
			gtui.Dots("text-muted", frame),
			gtui.Text("text-muted nowrap", "working"),
		),
	}
}

// room reports the cells a page has inside #main: the shell spends the rest on
// the sidebar, the header, the gaps and the borders, and the sidebar itself is
// narrower on a small terminal.
func room() (int, int) {
	width, height := app.Size()

	cols := width - 23
	if compact() {
		cols = width - 10
	}
	if cols < 24 {
		cols = 24
	}

	rows := height - 9
	if rows < 8 {
		rows = 8
	}

	return cols, rows
}

func wideScreen() bool {
	width, _ := room()
	return width >= 80
}

// demo is card() with a line of explanation that only a terminal with rows to
// spare has to give up.
func demo(key, title, about string, body ...gtui.Component) gtui.Component {
	if about != "" && wideScreen() {
		body = append([]gtui.Component{note(about)}, body...)
	}
	return card(key, title, body...)
}
