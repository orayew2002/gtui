package gtui

import (
	"strings"
	"time"

	"github.com/orayew2002/gtui/src/render"
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

				// the handle has to win even while the task queue is full
				select {
				case a.tasks <- func() { fn(at) }:
				case <-stop:
					return
				case <-a.stopCh:
					return
				}
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

// After runs fn on the UI goroutine once delay has elapsed. A zero or negative
// delay runs fn on the next loop iteration, the way setTimeout(fn, 0) does.
func (a *App) After(delay time.Duration, fn func()) {
	if fn == nil {
		return
	}
	if delay < 0 {
		delay = 0
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
	if render.StringWidth(text) <= width {
		return Text(join("nowrap", class), text)
	}

	runes := []rune(text + "   ")
	count := len(runes)
	offset := ((tick % count) + count) % count

	var sb strings.Builder
	cells := 0

	for i := 0; cells < width && i < count+width; i++ {
		r := runes[(offset+i)%count]
		w := render.RuneWidth(r)

		if w == 0 && cells == 0 {
			continue
		}
		if cells+w > width {
			// a wide rune would straddle the edge, so pad the cell instead
			sb.WriteByte(' ')
			break
		}

		sb.WriteRune(r)
		cells += w
	}

	return Text(join("nowrap", class), sb.String())
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
	if ((tick%2)+2)%2 == 1 {
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
