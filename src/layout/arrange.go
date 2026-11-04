package layout

import (
	"sort"

	"github.com/orayew2002/gtui/src/render"
	"github.com/orayew2002/gtui/src/style"
)

func arrange(b *Box, x, y, availW, availH int, clip render.Rect) {
	st := b.Style

	b.Frame = render.Rect{X: x, Y: y, W: b.OuterW, H: b.OuterH}
	b.Content = b.Frame.Inset(b.Ins)
	b.Clip = clip.Intersect(b.Frame)

	if b.Node == nil || len(b.Children) == 0 {
		return
	}

	childClip := clip
	if st.Overflow != style.OverflowVisible {
		childClip = clip.Intersect(b.Content)
	}

	scroll := scrollOffset(b)
	shift := scrollOffsetX(b)

	content := b.Content
	if shift != 0 {
		b.Content.X -= shift
	}

	switch st.Display {
	case style.DisplayFlex:
		switch {
		case st.Direction == style.Row && st.FlexWrap:
			arrangeWrapped(b, childClip, scroll)
		case st.Direction == style.Row:
			arrangeRow(b, childClip, scroll)
		default:
			arrangeColumn(b, childClip, scroll)
		}
	default:
		arrangeBlock(b, childClip, scroll)
	}

	b.Content = content

	arrangeAbsolute(b, childClip)
	sortByZ(b)
}

func scrollOffset(b *Box) int {
	flowH := 0
	for _, c := range flowChildren(b) {
		flowH += c.OuterH + c.Style.Margin.Vertical()
	}
	if b.Style.Display == style.DisplayFlex && b.Style.Direction == style.Row {
		flowH = 0
		for _, c := range flowChildren(b) {
			if h := c.OuterH + c.Style.Margin.Vertical(); h > flowH {
				flowH = h
			}
		}
	}
	if n := len(flowChildren(b)); n > 1 && b.Style.Direction == style.Column {
		flowH += b.Style.RowGap() * (n - 1)
	}

	b.ScrollMax = flowH - b.Content.H
	if b.ScrollMax < 0 {
		b.ScrollMax = 0
	}

	if b.Node == nil || b.Style.Overflow == style.OverflowVisible {
		return 0
	}

	if b.Node.ScrollY > b.ScrollMax {
		b.Node.ScrollY = b.ScrollMax
	}
	if b.Node.ScrollY < 0 {
		b.Node.ScrollY = 0
	}
	return b.Node.ScrollY
}

func scrollOffsetX(b *Box) int {
	if b.Node == nil || !b.Style.ScrollX {
		return 0
	}

	if b.Node.ScrollX > b.ScrollMaxX {
		b.Node.ScrollX = b.ScrollMaxX
	}
	if b.Node.ScrollX < 0 {
		b.Node.ScrollX = 0
	}

	return b.Node.ScrollX
}

func arrangeBlock(b *Box, clip render.Rect, scroll int) {
	y := b.Content.Y - scroll

	for _, c := range flowChildren(b) {
		m := c.Style.Margin
		x := b.Content.X + m.Left

		switch crossAlign(b, c) {
		case style.AlignCenter:
			x = b.Content.X + (b.Content.W-c.OuterW)/2
		case style.AlignEnd:
			x = b.Content.Right() - c.OuterW - m.Right
		}
		if x < b.Content.X && crossAlign(b, c) == style.AlignCenter {
			x = b.Content.X
		}

		y += m.Top
		arrange(c, x, y, b.Content.W, b.Content.H, clip)
		y += c.OuterH + m.Bottom
	}
}

func arrangeRow(b *Box, clip render.Rect, scroll int) {
	flow := flowChildren(b)
	if len(flow) == 0 {
		return
	}

	used := 0
	for _, c := range flow {
		used += c.OuterW + c.Style.Margin.Horizontal()
	}
	used += b.Style.ColumnGap() * (len(flow) - 1)

	start, between := mainOffsets(b.Style.Justify, b.Content.W-used, len(flow))
	x := b.Content.X + start

	for i, c := range flow {
		if i > 0 {
			x += b.Style.ColumnGap() + between
		}
		m := c.Style.Margin
		x += m.Left

		y := crossStart(b.Content.Y, b.Content.H, c.OuterH, m.Top, m.Bottom, crossAlign(b, c)) - scroll
		arrange(c, x, y, b.Content.W, b.Content.H, clip)

		x += c.OuterW + m.Right
	}
}

func arrangeWrapped(b *Box, clip render.Rect, scroll int) {
	y := b.Content.Y - scroll

	for i, line := range b.FlexLines {
		if i > 0 {
			y += b.Style.RowGap()
		}

		used := 0
		for _, c := range line {
			used += c.OuterW + c.Style.Margin.Horizontal()
		}
		if n := len(line); n > 1 {
			used += b.Style.ColumnGap() * (n - 1)
		}

		start, between := mainOffsets(b.Style.Justify, b.Content.W-used, len(line))
		x := b.Content.X + start

		tallest := 0
		for j, c := range line {
			if j > 0 {
				x += b.Style.ColumnGap() + between
			}

			m := c.Style.Margin
			x += m.Left

			cross := crossAlign(b, c)
			cy := crossStart(y, lineHeight(line), c.OuterH, m.Top, m.Bottom, cross)
			arrange(c, x, cy, b.Content.W, b.Content.H, clip)

			x += c.OuterW + m.Right
			if h := c.OuterH + m.Vertical(); h > tallest {
				tallest = h
			}
		}

		y += tallest
	}
}

func lineHeight(line []*Box) int {
	tallest := 0
	for _, c := range line {
		if h := c.OuterH + c.Style.Margin.Vertical(); h > tallest {
			tallest = h
		}
	}
	return tallest
}

func crossAlign(parent, child *Box) style.Align {
	if child.Style.HasAlignSelf {
		return child.Style.AlignSelf
	}
	return parent.Style.Align
}

func arrangeColumn(b *Box, clip render.Rect, scroll int) {
	flow := flowChildren(b)
	if len(flow) == 0 {
		return
	}

	used := 0
	for _, c := range flow {
		used += c.OuterH + c.Style.Margin.Vertical()
	}
	used += b.Style.RowGap() * (len(flow) - 1)

	start, between := mainOffsets(b.Style.Justify, b.Content.H-used, len(flow))
	y := b.Content.Y + start - scroll

	for i, c := range flow {
		if i > 0 {
			y += b.Style.RowGap() + between
		}
		m := c.Style.Margin
		y += m.Top

		x := crossStart(b.Content.X, b.Content.W, c.OuterW, m.Left, m.Right, crossAlign(b, c))
		arrange(c, x, y, b.Content.W, b.Content.H, clip)

		y += c.OuterH + m.Bottom
	}
}

func mainOffsets(j style.Justify, free, count int) (start, between int) {
	if free < 0 {
		free = 0
	}

	switch j {
	case style.JustifyCenter:
		return free / 2, 0
	case style.JustifyEnd:
		return free, 0
	case style.JustifyBetween:
		if count < 2 {
			return 0, 0
		}
		return 0, free / (count - 1)
	case style.JustifyAround:
		if count == 0 {
			return 0, 0
		}
		unit := free / (count * 2)
		return unit, unit * 2
	case style.JustifyEvenly:
		unit := free / (count + 1)
		return unit, unit
	}
	return 0, 0
}

// crossStart places a box across the line it sits on. A box larger than the
// line is left at the start rather than centred or pushed past it: centring it
// would put half of it before the container, where it would paint over
// whatever sits above or to the left of it.
func crossStart(origin, extent, size, before, after int, a style.Align) int {
	free := extent - size
	if free < 0 {
		free = 0
	}

	switch a {
	case style.AlignCenter:
		return origin + free/2
	case style.AlignEnd:
		return origin + free - after + min(after, 0)
	default:
		return origin + before
	}
}

func arrangeAbsolute(b *Box, clip render.Rect) {
	for _, c := range b.Abs {
		st := c.Style
		x := b.Content.X
		y := b.Content.Y

		if v, ok := st.Left.Resolve(b.Content.W); ok {
			x = b.Content.X + v
		} else if v, ok := st.Right.Resolve(b.Content.W); ok {
			x = b.Content.Right() - v - c.OuterW
		}

		if v, ok := st.Top.Resolve(b.Content.H); ok {
			y = b.Content.Y + v
		} else if v, ok := st.Bottom.Resolve(b.Content.H); ok {
			y = b.Content.Bottom() - v - c.OuterH
		}

		arrange(c, x, y, b.Content.W, b.Content.H, clip)
	}
}

func sortByZ(b *Box) {
	needsSort := false
	for _, c := range b.Children {
		if c.Style.ZIndex != 0 {
			needsSort = true
			break
		}
	}
	if !needsSort {
		return
	}

	sort.SliceStable(b.Children, func(i, j int) bool {
		return b.Children[i].Style.ZIndex < b.Children[j].Style.ZIndex
	})
}
