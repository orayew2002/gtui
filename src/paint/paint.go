package paint

import (
	"github.com/orayew2002/gtui/src/colors"
	"github.com/orayew2002/gtui/src/layout"
	"github.com/orayew2002/gtui/src/render"
	"github.com/orayew2002/gtui/src/style"
)

func Paint(buf *render.Buffer, root *layout.Box) {
	if root == nil {
		return
	}
	paintBox(buf, root)
}

func paintBox(buf *render.Buffer, b *layout.Box) {
	if b.Clip.Empty() {
		return
	}

	if b.Node == nil {
		paintText(buf, b)
		return
	}

	st := b.Style
	Frame(buf, b.Frame, b.Clip, st)

	if b.Node.Canvas != nil {
		paintCanvas(buf, b)
	}

	for _, c := range b.Children {
		paintBox(buf, c)
	}

	Hint(buf, b.Frame, b.Clip, st)

	if st.Overflow == style.OverflowScroll && b.ScrollMax > 0 {
		paintScrollbar(buf, b)
	}
}

func Frame(buf *render.Buffer, frame, clip render.Rect, st style.Style) {
	area := clip.Intersect(frame)
	if area.Empty() {
		return
	}

	g := borderGlyphs(st.Border)
	sides := st.BorderWidths()
	cut := cutCorners(frame, sides, st)

	if st.Background.Set {
		cut.snapshot(buf)
		buf.Fill(area, render.Cell{Ch: ' ', Bg: st.Background.Color, HasBg: true})
		cut.restore(buf)
	}

	if !st.HasBorder() {
		return
	}

	proto := render.Cell{}
	if ink := st.BorderColor.Or(st.Foreground); ink.Set {
		proto.Fg, proto.HasFg = ink.Color, true
	}
	if st.Background.Set {
		proto.Bg, proto.HasBg = st.Background.Color, true
	}

	set := func(x, y int, r rune) {
		if !clip.Contains(x, y) {
			return
		}

		c := proto
		c.Ch = r

		if cut.has(x, y) {
			c.Bg, c.HasBg = render.Cell{}.Bg, false
		}

		buf.SetGlyph(x, y, c)
	}

	if sides.Top > 0 {
		for x := frame.X; x < frame.Right(); x++ {
			set(x, frame.Y, g.topRune())
		}
	}
	if sides.Bottom > 0 && (frame.H > 1 || sides.Top == 0) {
		for x := frame.X; x < frame.Right(); x++ {
			set(x, frame.Bottom()-1, g.bottomRune())
		}
	}
	if sides.Left > 0 {
		for y := frame.Y; y < frame.Bottom(); y++ {
			set(frame.X, y, g.leftRune())
		}
	}
	if sides.Right > 0 && (frame.W > 1 || sides.Left == 0) {
		for y := frame.Y; y < frame.Bottom(); y++ {
			set(frame.Right()-1, y, g.rightRune())
		}
	}

	if sides.Top > 0 && sides.Left > 0 {
		set(frame.X, frame.Y, g.tl)
	}
	if sides.Top > 0 && sides.Right > 0 && frame.W > 1 {
		set(frame.Right()-1, frame.Y, g.tr)
	}
	if sides.Bottom > 0 && sides.Left > 0 && frame.H > 1 {
		set(frame.X, frame.Bottom()-1, g.bl)
	}
	if sides.Bottom > 0 && sides.Right > 0 && frame.W > 1 && frame.H > 1 {
		set(frame.Right()-1, frame.Bottom()-1, g.br)
	}

	if st.Title != "" && sides.Top > 0 {
		paintTitle(buf, frame, clip, st, proto)
	}
}

func Hint(buf *render.Buffer, frame, clip render.Rect, st style.Style) {
	if st.Shortcut == "" || st.HideShortcut {
		return
	}

	label := "[" + st.Shortcut + "]"
	width := render.StringWidth(label)
	if frame.W < width {
		return
	}

	inset := st.BorderWidths()
	x := frame.Right() - width - inset.Right
	y := frame.Y + inset.Top
	if !st.HasBorder() {
		x = frame.Right() - width
	}

	cell := render.Cell{}
	if ink := st.BorderColor.Or(st.Foreground); ink.Set {
		cell.Fg, cell.HasFg = ink.Color, true
	}
	if st.Background.Set {
		cell.Bg, cell.HasBg = st.Background.Color, true
	}

	line := render.Rect{X: frame.X, Y: y, W: frame.W, H: 1}.Intersect(clip)
	buf.SetString(x, y, label, cell, line)
}

type corners struct {
	n     int
	x, y  [4]int
	saved [4]render.Cell
}

func cutCorners(frame render.Rect, sides style.Edges, st style.Style) corners {
	var c corners

	if frame.W < 2 || frame.H < 2 {
		return c
	}

	switch st.Border {
	case style.BorderRounded, style.BorderSoft, style.BorderTight:
	default:
		return c
	}

	if sides.Top > 0 && sides.Left > 0 {
		c.add(frame.X, frame.Y)
	}
	if sides.Top > 0 && sides.Right > 0 {
		c.add(frame.Right()-1, frame.Y)
	}
	if sides.Bottom > 0 && sides.Left > 0 {
		c.add(frame.X, frame.Bottom()-1)
	}
	if sides.Bottom > 0 && sides.Right > 0 {
		c.add(frame.Right()-1, frame.Bottom()-1)
	}

	return c
}

func (c *corners) add(x, y int) {
	c.x[c.n], c.y[c.n] = x, y
	c.n++
}

func (c *corners) has(x, y int) bool {
	for i := range c.n {
		if c.x[i] == x && c.y[i] == y {
			return true
		}
	}
	return false
}

func (c *corners) snapshot(buf *render.Buffer) {
	for i := range c.n {
		if p := buf.At(c.x[i], c.y[i]); p != nil {
			c.saved[i] = *p
		}
	}
}

func (c *corners) restore(buf *render.Buffer) {
	for i := range c.n {
		buf.Set(c.x[i], c.y[i], c.saved[i])
	}
}

func paintTitle(buf *render.Buffer, frame, clip render.Rect, st style.Style, proto render.Cell) {
	avail := frame.W - 2
	if avail <= 2 {
		return
	}

	text, w := render.TruncateWidth(st.Title, avail-3)
	if w == 0 {
		return
	}

	x := frame.X + 2
	switch st.TitleAlign {
	case style.AlignCenter:
		x = frame.X + (frame.W-w-2)/2
	case style.AlignEnd:
		x = frame.Right() - w - 3
	}
	if x < frame.X+1 {
		x = frame.X + 1
	}

	cell := proto
	if st.Foreground.Set {
		cell.Fg, cell.HasFg = st.Foreground.Color, true
	}
	cell.Attr = st.Attr

	line := render.Rect{X: frame.X, Y: frame.Y, W: frame.W, H: 1}.Intersect(clip)
	buf.SetString(x, frame.Y, " "+text+" ", cell, line)
}

func paintCanvas(buf *render.Buffer, b *layout.Box) {
	canvas := b.Node.Canvas
	area := b.Content.Intersect(b.Clip)

	for y := area.Y; y < area.Bottom(); y++ {
		for x := area.X; x < area.Right(); x++ {
			pixel := canvas.At(x-b.Content.X, y-b.Content.Y)
			if pixel.Ch == 0 && !pixel.HasBg {
				continue
			}

			cell := render.Cell{Ch: pixel.Ch, Fg: pixel.Fg, Bg: pixel.Bg, HasFg: pixel.HasFg, HasBg: pixel.HasBg}
			if cell.Ch == 0 {
				cell.Ch = ' '
			}
			if pixel.Bold {
				cell.Attr |= style.AttrBold
			}

			buf.SetGlyph(x, y, cell)
		}
	}
}

func paintText(buf *render.Buffer, b *layout.Box) {
	if len(b.Lines) == 0 {
		return
	}

	clip := b.Clip
	align := b.Style.TextAlign

	for i, line := range b.Lines {
		y := b.Content.Y + i
		if y < clip.Y || y >= clip.Bottom() {
			continue
		}

		offset := 0
		switch align {
		case style.AlignCenter:
			offset = (b.Content.W - line.W) / 2
		case style.AlignEnd:
			offset = b.Content.W - line.W
		}
		if offset < 0 {
			offset = 0
		}

		for _, run := range line.Runs {
			ink := run.Style
			if ink == nil {
				ink = &b.Style
			}

			cell := render.Cell{Attr: ink.Attr}
			if ink.Foreground.Set {
				cell.Fg, cell.HasFg = ink.Foreground.Color, true
			}
			if ink.Background.Set {
				cell.Bg, cell.HasBg = ink.Background.Color, true
			}
			buf.SetString(b.Content.X+offset+run.X, y, run.Text, cell, clip)
		}
	}
}

func paintScrollbar(buf *render.Buffer, b *layout.Box) {
	track := render.Rect{X: b.Content.Right(), Y: b.Content.Y, W: 1, H: b.Content.H}
	track = track.Intersect(b.Clip)
	if track.Empty() {
		return
	}

	total := b.Content.H + b.ScrollMax
	thumb := track.H * b.Content.H / total
	if thumb < 1 {
		thumb = 1
	}

	pos := 0
	if b.ScrollMax > 0 {
		pos = (track.H - thumb) * b.Node.ScrollY / b.ScrollMax
	}

	fg := colors.Overlay
	if b.Style.BorderColor.Set {
		fg = b.Style.BorderColor.Color
	}

	for y := 0; y < track.H; y++ {
		ch := '│'
		cell := render.Cell{Ch: ch, Fg: colors.Overlay, HasFg: true}
		if y >= pos && y < pos+thumb {
			cell.Ch = '█'
			cell.Fg = fg
		}
		if b.Style.Background.Set {
			cell.Bg, cell.HasBg = b.Style.Background.Color, true
		}
		buf.SetGlyph(track.X, track.Y+y, cell)
	}
}
