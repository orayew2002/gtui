package styles

import (
	"io"
	"strings"

	"github.com/orayew2002/gtui/src/colors"
	"github.com/orayew2002/gtui/src/paint"
	"github.com/orayew2002/gtui/src/render"
	"github.com/orayew2002/gtui/src/style"
)

type Box struct {
	x, y int
	w, h int
	st   style.Style
}

func NewBox(w, h int) *Box {
	st := style.New()
	st.Border = style.BorderRounded
	st.BorderSides = style.Edges{Top: 1, Right: 1, Bottom: 1, Left: 1}

	return &Box{w: w, h: h, st: st}
}

func (b *Box) At(x, y int) *Box {
	b.x, b.y = x, y
	return b
}

func (b *Box) Size(w, h int) *Box {
	b.w, b.h = w, h
	return b
}

func (b *Box) Border(s style.BorderStyle) *Box {
	b.st.Border = s
	if s == style.BorderNone {
		b.st.BorderSides = style.Edges{}
		return b
	}
	b.st.BorderSides = style.Edges{Top: 1, Right: 1, Bottom: 1, Left: 1}
	return b
}

func (b *Box) BorderSides(e style.Edges) *Box {
	b.st.BorderSides = e
	return b
}

func (b *Box) BorderColor(c colors.Color) *Box {
	b.st.BorderColor = style.P(c)
	return b
}

func (b *Box) Background(c colors.Color) *Box {
	b.st.Background = style.P(c)
	return b
}

func (b *Box) Foreground(c colors.Color) *Box {
	b.st.Foreground = style.P(c)
	return b
}

func (b *Box) Title(text string) *Box {
	b.st.Title = text
	return b
}

func (b *Box) TitleAlign(a style.Align) *Box {
	b.st.TitleAlign = a
	return b
}

func (b *Box) Style() style.Style { return b.st }

func (b *Box) Rect() render.Rect {
	return render.Rect{X: b.x, Y: b.y, W: b.w, H: b.h}
}

func (b *Box) Buffer() *render.Buffer {
	buf := render.NewBuffer(b.w, b.h)
	if b.w <= 0 || b.h <= 0 {
		return buf
	}

	buf.Reset(render.Cell{Ch: ' '})

	local := render.Rect{W: b.w, H: b.h}
	paint.Frame(buf, local, local, b.st)

	return buf
}

func (b *Box) Render(w io.Writer) error {
	if b.w <= 0 || b.h <= 0 {
		return nil
	}
	return render.Draw(w, b.Buffer(), b.x, b.y)
}

func (b *Box) String() string {
	buf := b.Buffer()

	var sb strings.Builder
	for y := range buf.H {
		row := make([]rune, 0, buf.W)

		for x := range buf.W {
			cell := buf.At(x, y)
			switch {
			case cell == nil, cell.Cont:
				continue
			case cell.Ch == 0:
				row = append(row, ' ')
			default:
				row = append(row, cell.Ch)
			}
		}

		sb.WriteString(strings.TrimRight(string(row), " "))
		if y < buf.H-1 {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}
