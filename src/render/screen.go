package render

import "io"

type Screen struct {
	pen *pen

	back  *Buffer
	front *Buffer

	full bool
}

func NewScreen(w io.Writer, width, height int) *Screen {
	s := &Screen{
		pen:   newPen(w),
		back:  NewBuffer(width, height),
		front: NewBuffer(width, height),
		full:  true,
	}
	s.back.Reset(Blank)
	s.front.Reset(Blank)
	return s
}

func (s *Screen) Back() *Buffer { return s.back }

func (s *Screen) Size() (int, int) { return s.back.W, s.back.H }

func (s *Screen) Resize(w, h int) {
	if w == s.back.W && h == s.back.H {
		return
	}
	s.back.Resize(w, h)
	s.front.Resize(w, h)
	s.front.Reset(Blank)
	s.Invalidate()
}

func (s *Screen) Invalidate() {
	s.full = true
	s.pen.reset()
}

func (s *Screen) Clear(c Cell) { s.back.Reset(c) }

func (s *Screen) Flush() error {
	if s.full {
		s.pen.clearScreen()
	}

	w, h := s.back.W, s.back.H
	for y := range h {
		rowOff := y * w
		for x := range w {
			i := rowOff + x
			cell := s.back.Cells[i]

			if cell.Cont {
				s.front.Cells[i] = cell
				continue
			}
			if !s.full && cell.eq(s.front.Cells[i]) {
				continue
			}

			s.pen.moveTo(x, y)
			s.pen.writeCell(cell)
			s.front.Cells[i] = cell

			if cx, _ := s.pen.at(); cx >= w {
				s.pen.forgetPosition()
			}
		}
	}

	s.full = false
	return s.pen.flush()
}

func (s *Screen) Write(p []byte) (int, error) { return s.pen.write(p) }

func (s *Screen) ShowCursorAt(x, y int) {
	s.pen.moveTo(x, y)
	s.pen.writeString("\x1b[?25h")
}

func (s *Screen) HideCursor() {
	s.pen.writeString("\x1b[?25l")
}
