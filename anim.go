package gtui

import (
	"strings"
	"time"
)

func (a *App) Every(interval time.Duration, fn func(tick int)) func() {
	if interval <= 0 || fn == nil {
		return func() {}
	}

	stop := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		tick := 0
		for {
			select {
			case <-stop:
				return
			case <-a.stopCh:
				return
			case <-ticker.C:
				tick++
				at := tick
				a.Do(func() { fn(at) })
			}
		}
	}()

	return func() {
		select {
		case <-stop:
		default:
			close(stop)
		}
	}
}

func (a *App) Animate(fps int, fn func(tick int)) func() {
	if fps < 1 {
		fps = 10
	}
	return a.Every(time.Second/time.Duration(fps), fn)
}

func (a *App) After(delay time.Duration, fn func()) {
	if delay <= 0 || fn == nil {
		return
	}

	go func() {
		select {
		case <-time.After(delay):
			a.Do(fn)
		case <-a.stopCh:
		}
	}()
}

func Marquee(class, text string, width, tick int) Component {
	if width < 1 || text == "" {
		return Text(class, text)
	}

	runes := []rune(text + "   ")
	if len(runes) <= width {
		return Text(join("nowrap", class), text)
	}

	offset := ((tick % len(runes)) + len(runes)) % len(runes)
	window := make([]rune, 0, width)

	for i := range width {
		window = append(window, runes[(offset+i)%len(runes)])
	}

	return Text(join("nowrap", class), string(window))
}

func Typewriter(class, text string, tick int) Component {
	runes := []rune(text)
	shown := clampInt(tick, 0, len(runes))

	caret := ""
	if shown < len(runes) {
		caret = "▌"
	}

	return Text(join("nowrap", class), string(runes[:shown])+caret)
}

func Blink(class, text string, tick int) Component {
	if tick%2 == 1 {
		return Text(join("nowrap dim", class), text)
	}
	return Text(join("nowrap", class), text)
}

func Pulse(class string, width, tick int) Component {
	if width < 3 {
		width = 12
	}

	head := ((tick % width) + width) % width
	row := make([]rune, width)

	for i := range row {
		switch {
		case i == head:
			row[i] = '█'
		case i == (head+width-1)%width, i == (head+1)%width:
			row[i] = '▓'
		default:
			row[i] = '░'
		}
	}

	return Text(join("nowrap", class), string(row))
}

func Wave(class string, width, tick int) Component {
	glyphs := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█', '▇', '▆', '▅', '▄', '▃', '▂'}
	if width < 1 {
		width = 16
	}

	var sb strings.Builder
	for i := range width {
		sb.WriteRune(glyphs[((i+tick)%len(glyphs)+len(glyphs))%len(glyphs)])
	}

	return Text(join("nowrap", class), sb.String())
}

func Dots(class string, tick int) Component {
	return Text(join("nowrap", class), strings.Repeat(".", ((tick%3)+3)%3+1))
}
