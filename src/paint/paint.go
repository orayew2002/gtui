package paint

import (
	"unicode"

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
	if st.Shadow {
		Shadow(buf, b.Frame, outerClip(buf, b))
	}
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
		if st.Background.Translucent() {
			buf.Blend(area, st.Background.Color, st.Background.Alpha)
		} else {
			buf.Fill(area, render.Cell{Ch: ' ', Bg: st.Background.Color, HasBg: true})
		}
		cut.restore(buf)
	}

	if !st.HasBorder() {
		return
	}

	proto := render.Cell{}
	if ink := st.BorderColor.Or(st.Foreground); ink.Set {
		proto.Fg, proto.HasFg = ink.Color, true
	}
	if st.Background.Set && !st.Background.Translucent() {
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

	caps := g.capL != 0 && (sides.Left > 0 || sides.Right > 0)
	if caps && frame.H == 1 {
		if sides.Left > 0 {
			set(frame.X, frame.Y, g.capL)
		}
		if sides.Right > 0 && frame.W > 1 {
			set(frame.Right()-1, frame.Y, g.capR)
		}
		return
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

	if (sides.Top > 0 || caps) && sides.Left > 0 {
		set(frame.X, frame.Y, g.tl)
	}
	if (sides.Top > 0 || caps) && sides.Right > 0 && frame.W > 1 {
		set(frame.Right()-1, frame.Y, g.tr)
	}
	if (sides.Bottom > 0 || caps) && sides.Left > 0 && frame.H > 1 {
		set(frame.X, frame.Bottom()-1, g.bl)
	}
	if (sides.Bottom > 0 || caps) && sides.Right > 0 && frame.W > 1 && frame.H > 1 {
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

	label := "┤" + st.Shortcut + "├"
	if !st.HasBorder() {
		label = st.Shortcut
	}

	width := render.StringWidth(label)
	if frame.W < width {
		return
	}

	inset := st.BorderWidths()

	x := frame.Right() - width - inset.Right
	y := frame.Y + inset.Top

	if st.HasBorder() && inset.Top > 0 {
		// Sit in the top border, like a legend, so the badge never covers
		// content or a scrollbar.
		x = frame.Right() - width - 1
		y = frame.Y
	} else if !st.HasBorder() {
		x = frame.Right() - width
	}

	cell := render.Cell{Attr: style.AttrBold}
	cell.Fg, cell.HasFg = colors.Active().Accent, true
	if st.Background.Set {
		cell.Bg, cell.HasBg = st.Background.Color, true
	}

	line := render.Rect{X: frame.X, Y: y, W: frame.W, H: 1}.Intersect(clip)
	buf.SetString(x, y, label, cell, line)
}

func Shadow(buf *render.Buffer, frame, clip render.Rect) {
	if frame.Empty() {
		return
	}

	right := render.Rect{X: frame.Right(), Y: frame.Y + 1, W: 1, H: frame.H}
	bottom := render.Rect{X: frame.X + 1, Y: frame.Bottom(), W: frame.W, H: 1}

	shade(buf, right.Intersect(clip))
	shade(buf, bottom.Intersect(clip))
}

func shade(buf *render.Buffer, area render.Rect) {
	if area.Empty() {
		return
	}

	fallback := colors.Active().Line

	for y := area.Y; y < area.Bottom(); y++ {
		for x := area.X; x < area.Right(); x++ {
			p := buf.At(x, y)
			if p == nil {
				continue
			}

			cell := render.Cell{Ch: '░', Attr: style.AttrDim}
			cell.Fg, cell.HasFg = fallback, true
			if p.HasBg {
				cell.Bg, cell.HasBg = p.Bg, true
				cell.Fg = darken(p.Bg)
			}

			*p = cell
		}
	}
}

func darken(c colors.Color) colors.Color {
	return colors.RGB(c.R/2, c.G/2, c.B/2)
}

func outerClip(buf *render.Buffer, b *layout.Box) render.Rect {
	p := b.Parent
	if p == nil {
		return render.Rect{W: buf.W, H: buf.H}
	}
	if p.Style.Overflow != style.OverflowVisible {
		return p.Clip.Intersect(p.Content)
	}
	return p.Clip
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
	case style.BorderPill:
		sides = style.Edges{Top: 1, Right: sides.Right, Bottom: 1, Left: sides.Left}
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

			x := b.Content.X + offset + run.X
			if ink.Transform == style.TransformNone {
				buf.SetString(x, y, run.Text, cell, clip)
				continue
			}
			setTransformed(buf, x, y, run.Text, cell, clip, ink.Transform)
		}
	}
}

func setTransformed(buf *render.Buffer, x, y int, text string, proto render.Cell, clip render.Rect, t style.Transform) {
	start := true

	for _, r := range text {
		g := transformRune(r, t, start)
		start = !isWordRune(r)

		w := render.RuneWidth(g)
		if w == 0 {
			continue
		}
		if x >= clip.Right() {
			return
		}

		if clip.Contains(x, y) {
			cell := proto
			cell.Ch = g
			buf.SetGlyph(x, y, cell)

			if w == 2 && clip.Contains(x+1, y) {
				cont := proto
				cont.Ch, cont.Cont = 0, true
				buf.SetGlyph(x+1, y, cont)
			}
		}

		x += w
	}
}

func transformRune(r rune, t style.Transform, start bool) rune {
	switch t {
	case style.TransformUpper:
		return unicode.ToUpper(r)
	case style.TransformLower:
		return unicode.ToLower(r)
	case style.TransformTitle:
		if start {
			return unicode.ToUpper(r)
		}
	}
	return r
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '\''
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
