package render

import (
	"github.com/orayew2002/gtui/src/colors"
	"github.com/orayew2002/gtui/src/style"
)

type Rect struct {
	X, Y, W, H int
}

func (r Rect) Right() int  { return r.X + r.W }
func (r Rect) Bottom() int { return r.Y + r.H }
func (r Rect) Empty() bool { return r.W <= 0 || r.H <= 0 }

func (r Rect) Intersect(o Rect) Rect {
	x := max(r.X, o.X)
	y := max(r.Y, o.Y)
	right := min(r.Right(), o.Right())
	bottom := min(r.Bottom(), o.Bottom())
	if right <= x || bottom <= y {
		return Rect{}
	}
	return Rect{x, y, right - x, bottom - y}
}

func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x < r.Right() && y >= r.Y && y < r.Bottom()
}

func (r Rect) Inset(e style.Edges) Rect {
	out := Rect{
		X: r.X + e.Left,
		Y: r.Y + e.Top,
		W: r.W - e.Horizontal(),
		H: r.H - e.Vertical(),
	}
	if out.W < 0 {
		out.W = 0
	}
	if out.H < 0 {
		out.H = 0
	}
	return out
}

type Cell struct {
	Ch    rune
	Fg    colors.Color
	Bg    colors.Color
	Attr  style.Attrs
	HasFg bool
	HasBg bool

	Cont bool
}

var Blank = Cell{Ch: ' '}

func (c Cell) eq(o Cell) bool {
	return c.Ch == o.Ch &&
		c.Attr == o.Attr &&
		c.HasFg == o.HasFg && c.HasBg == o.HasBg &&
		c.Cont == o.Cont &&
		(!c.HasFg || c.Fg == o.Fg) &&
		(!c.HasBg || c.Bg == o.Bg)
}

type Buffer struct {
	W, H  int
	Cells []Cell
}

func NewBuffer(w, h int) *Buffer {
	b := &Buffer{}
	b.Resize(w, h)
	return b
}

func (b *Buffer) Resize(w, h int) {
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	b.W, b.H = w, h
	need := w * h
	if cap(b.Cells) < need {
		b.Cells = make([]Cell, need)
	}
	b.Cells = b.Cells[:need]
}

func (b *Buffer) Reset(c Cell) {
	for i := range b.Cells {
		b.Cells[i] = c
	}
}

func (b *Buffer) At(x, y int) *Cell {
	if x < 0 || y < 0 || x >= b.W || y >= b.H {
		return nil
	}
	return &b.Cells[y*b.W+x]
}

func (b *Buffer) Set(x, y int, c Cell) {
	if p := b.At(x, y); p != nil {
		*p = c
	}
}

func (b *Buffer) SetGlyph(x, y int, c Cell) {
	p := b.At(x, y)
	if p == nil {
		return
	}

	if !c.HasBg {
		c.Bg, c.HasBg = p.Bg, p.HasBg
	}
	*p = c
}

func (b *Buffer) Fill(r Rect, c Cell) {
	r = r.Intersect(Rect{0, 0, b.W, b.H})
	for y := r.Y; y < r.Bottom(); y++ {
		row := b.Cells[y*b.W : (y+1)*b.W]
		for x := r.X; x < r.Right(); x++ {
			row[x] = c
		}
	}
}

// Blend mixes bg into the cells of r, keeping whatever is already painted there
// as the layer underneath: alpha 0 leaves them untouched and 255 covers them.
// A cell with no background of its own is mixed against the theme's canvas,
// which is the colour the app would have painted there.
func (b *Buffer) Blend(r Rect, bg colors.Color, alpha uint8) {
	if alpha == 0 {
		return
	}

	r = r.Intersect(Rect{0, 0, b.W, b.H})
	canvas := colors.Active().Canvas

	for y := r.Y; y < r.Bottom(); y++ {
		row := b.Cells[y*b.W : (y+1)*b.W]
		for x := r.X; x < r.Right(); x++ {
			cell := &row[x]

			under := canvas
			if cell.HasBg {
				under = cell.Bg
			}

			cell.Bg, cell.HasBg = Mix(under, bg, alpha), true

			// a cell already carrying a glyph has part of the picture in its
			// foreground, so the scrim has to reach that too
			if cell.HasFg {
				cell.Fg = Mix(cell.Fg, bg, alpha)
			}
		}
	}
}

// Mix lays over on top of under at the given alpha.
func Mix(under, over colors.Color, alpha uint8) colors.Color {
	if alpha >= 255 {
		return over
	}

	a, inv := int(alpha), 255-int(alpha)

	return colors.Color{
		R: uint8((int(over.R)*a + int(under.R)*inv) / 255),
		G: uint8((int(over.G)*a + int(under.G)*inv) / 255),
		B: uint8((int(over.B)*a + int(under.B)*inv) / 255),
	}
}

func (b *Buffer) SetString(x, y int, s string, proto Cell, clip Rect) int {
	if y < clip.Y || y >= clip.Bottom() || y < 0 || y >= b.H {
		return 0
	}

	left, right := max(clip.X, 0), min(clip.Right(), b.W)
	if left >= right {
		return 0
	}

	row := b.Cells[y*b.W : (y+1)*b.W]
	drawn := 0

	for _, r := range s {
		w := RuneWidth(r)
		if w == 0 {
			continue
		}
		if x >= right {
			break
		}

		if x >= left {
			cell := proto
			cell.Ch = r
			if !cell.HasBg {
				cell.Bg, cell.HasBg = row[x].Bg, row[x].HasBg
			}
			row[x] = cell

			if w == 2 && x+1 >= left && x+1 < right {
				cont := proto
				cont.Ch = 0
				cont.Cont = true
				if !cont.HasBg {
					cont.Bg, cont.HasBg = row[x+1].Bg, row[x+1].HasBg
				}
				row[x+1] = cont
			}
		}

		x += w
		drawn += w
	}

	return drawn
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
