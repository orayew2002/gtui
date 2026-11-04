package layout

import (
	"github.com/orayew2002/gtui/src/render"
	"github.com/orayew2002/gtui/src/style"
)

func Layout(root *Box, viewport render.Rect) {
	if root == nil {
		return
	}

	measure(root, viewport.W, viewport.H)
	arrange(root, viewport.X, viewport.Y, viewport.W, viewport.H, viewport)
}

func measure(b *Box, availW, availH int) {
	if b.Node == nil {
		measureText(b, availW)
		return
	}

	st := b.Style
	ins := b.Ins

	if c := b.Node.Canvas; c != nil {
		w, ok := resolveMain(st.Width, availW, st.Margin.Horizontal(), ins.Horizontal())
		if !ok {
			w = c.W
		}
		finishCanvas(b, w, c.H, availH)
		return
	}

	contentW, definite := resolveMain(st.Width, availW, st.Margin.Horizontal(), ins.Horizontal())
	if !definite {
		switch {
		case st.Display == style.DisplayInline || st.Position == style.PositionAbsolute:
			contentW = intrinsicWidth(b, availW)
		case availW >= 0:
			contentW = availW - st.Margin.Horizontal() - ins.Horizontal()
		default:
			contentW = intrinsicWidth(b, -1)
		}
	}

	finish(b, contentW, availW, availH)
}

func measureOuterW(b *Box, outerW, availW, availH int) {
	if b.Node == nil {
		measureText(b, outerW)
		return
	}

	if c := b.Node.Canvas; c != nil {
		finishCanvas(b, outerW-b.Ins.Horizontal(), c.H, availH)
		return
	}

	ins := b.Ins
	finish(b, outerW-ins.Horizontal(), availW, availH)
}

func measureOuterH(b *Box, outerH, availW int) {
	if outerH < 0 {
		outerH = 0
	}

	if b.Node == nil {
		b.ContentH, b.OuterH = outerH, outerH
		return
	}

	ins := b.Ins
	contentH := outerH - ins.Vertical()
	if contentH < 0 {
		contentH = 0
	}

	switch b.Style.Display {
	case style.DisplayFlex:
		measureFlex(b, b.ContentW, contentH, true)
	default:
		measureBlock(b, b.ContentW, contentH)
	}

	b.ContentH = contentH
	b.OuterH = contentH + ins.Vertical()
}

func finishCanvas(b *Box, contentW, canvasH, availH int) {
	st := b.Style
	ins := b.Ins

	if contentW < 0 {
		contentW = 0
	}

	contentH := canvasH
	if h, ok := resolveMain(st.Height, availH, st.Margin.Vertical(), ins.Vertical()); ok {
		contentH = h
	}

	b.ContentW, b.ContentH = contentW, contentH
	b.OuterW = contentW + ins.Horizontal()
	b.OuterH = contentH + ins.Vertical()

	// A canvas can also be a backdrop: lay its children out inside it so an
	// image can sit behind ordinary content.
	if len(b.Children) == 0 {
		return
	}

	switch st.Display {
	case style.DisplayFlex:
		measureFlex(b, contentW, contentH, true)
	default:
		measureBlock(b, contentW, contentH)
	}

	b.ContentW, b.ContentH = contentW, contentH
	b.OuterW = contentW + ins.Horizontal()
	b.OuterH = contentH + ins.Vertical()
}

func finish(b *Box, contentW, availW, availH int) {
	st := b.Style
	ins := b.Ins

	contentW = clampDim(contentW, st.MinWidth, st.MaxWidth, availW, ins.Horizontal())
	b.ContentW = contentW

	contentH, definiteH := resolveMain(st.Height, availH, st.Margin.Vertical(), ins.Vertical())
	childAvailH := contentH
	if !definiteH {
		childAvailH = -1
	}

	childW := contentW
	if st.ScrollX {
		if natural := naturalWidth(b); natural > childW {
			childW = natural
		}
	}

	switch st.Display {
	case style.DisplayFlex:
		measureFlex(b, childW, childAvailH, definiteH)
	default:
		measureBlock(b, childW, childAvailH)
	}

	b.ScrollMaxX = childW - contentW
	if b.ScrollMaxX < 0 {
		b.ScrollMaxX = 0
	}

	if !definiteH {
		contentH = b.ContentH
	}
	contentH = clampDim(contentH, st.MinHeight, st.MaxHeight, availH, ins.Vertical())

	b.ContentH = contentH
	b.OuterW = contentW + ins.Horizontal()
	b.OuterH = contentH + ins.Vertical()
}

func measureText(b *Box, availW int) {
	width := availW
	if b.Style.NoWrap {
		width = 0
	}

	if !b.wrapped || b.wrapWidth != width {
		b.wrapText(width, b.Style.NoWrap)
		b.wrapWidth, b.wrapped = width, true
	}
	b.ContentH = len(b.Lines)

	b.ContentW = 0
	for _, l := range b.Lines {
		if l.W > b.ContentW {
			b.ContentW = l.W
		}
	}
	if availW >= 0 {
		b.ContentW = availW
	}

	b.OuterW, b.OuterH = b.ContentW, b.ContentH
}

func measureBlock(b *Box, contentW, availH int) {
	height := 0

	for _, c := range b.Abs {
		measure(c, contentW, availH)
	}

	for _, c := range b.Flow {
		measure(c, contentW, availH)
		height += c.OuterH + c.Style.Margin.Vertical()
	}

	b.ContentH = height
}

func measureFlex(b *Box, contentW, availH int, definiteH bool) {
	st := b.Style
	flow := flowChildren(b)

	gap := st.RowGap()
	if st.Direction == Row() {
		gap = st.ColumnGap()
	}

	gaps := 0
	if n := len(flow); n > 1 {
		gaps = gap * (n - 1)
	}

	for _, c := range b.Abs {
		measure(c, contentW, availH)
	}

	if st.Direction == Row() {
		mainAvail := contentW - gaps

		for _, c := range flow {
			if w, ok := basisWidth(c, contentW); ok {
				measureOuterW(c, w, contentW, availH)
				continue
			}
			measureOuterW(c, intrinsicWidth(c, mainAvail)+childInsets(c), contentW, availH)
		}

		if st.FlexWrap {
			measureWrapped(b, flow, contentW, availH)
			return
		}

		b.FlexLines = [][]*Box{flow}
		distribute(flow, mainAvail, true, contentW, availH)

		tallest := 0
		for _, c := range flow {
			if h := c.OuterH + c.Style.Margin.Vertical(); h > tallest {
				tallest = h
			}
		}
		b.ContentH = tallest

		if st.Align == style.AlignStretch {
			target := b.ContentH
			if definiteH {
				target = availH
			}
			for _, c := range flow {
				if !c.Style.Height.IsAuto() || c.Node == nil {
					continue
				}
				if c.Style.HasAlignSelf && c.Style.AlignSelf != style.AlignStretch {
					continue
				}

				h := target - c.Style.Margin.Vertical()

				if h == c.OuterH || (!definiteH && h < c.OuterH) {
					continue
				}
				if h < minMain(c, false) {
					h = minMain(c, false)
				}
				measureOuterH(c, h, contentW)
			}
		}
		return
	}

	for _, c := range flow {
		if st.Align != style.AlignStretch && c.Style.Width.IsAuto() {
			measureOuterW(c, intrinsicWidth(c, contentW)+childInsets(c), contentW, -1)
			continue
		}
		measure(c, contentW, -1)
	}

	b.FlexLines = [][]*Box{flow}

	if definiteH && availH >= 0 {
		distribute(flow, availH-gaps, false, contentW, availH)
		b.ContentH = availH
		return
	}

	total := gaps
	for _, c := range flow {
		total += c.OuterH + c.Style.Margin.Vertical()
	}
	b.ContentH = total
}

func measureWrapped(b *Box, flow []*Box, contentW, availH int) {
	gapX, gapY := b.Style.ColumnGap(), b.Style.RowGap()

	lines := [][]*Box{}
	line := []*Box{}
	used := 0

	for _, c := range flow {
		width := c.OuterW + c.Style.Margin.Horizontal()

		if len(line) > 0 && used+gapX+width > contentW {
			lines = append(lines, line)
			line, used = []*Box{}, 0
		}

		if len(line) > 0 {
			used += gapX
		}
		used += width
		line = append(line, c)
	}

	if len(line) > 0 {
		lines = append(lines, line)
	}

	total := 0
	for i, items := range lines {
		if i > 0 {
			total += gapY
		}

		avail := contentW
		if n := len(items); n > 1 {
			avail -= gapX * (n - 1)
		}
		distribute(items, avail, true, contentW, availH)

		tallest := 0
		for _, c := range items {
			if h := c.OuterH + c.Style.Margin.Vertical(); h > tallest {
				tallest = h
			}
		}
		total += tallest
	}

	b.FlexLines = lines
	b.ContentH = total
}

func distribute(flow []*Box, mainAvail int, horizontal bool, contentW, availH int) {
	if mainAvail < 0 || len(flow) == 0 {
		return
	}

	type item struct {
		box    *Box
		size   int
		min    int
		max    int
		weight int
		frozen bool
	}

	items := make([]item, len(flow))
	used := 0
	for i, c := range flow {
		size := mainOuter(c, horizontal)
		used += size
		items[i] = item{
			box:  c,
			size: size,
			min:  minMain(c, horizontal) + marginMain(c, horizontal),
			max:  maxMain(c, horizontal, contentW, availH),
		}
	}

	free := mainAvail - used
	if free == 0 {
		return
	}

	for i := range items {
		it := &items[i]
		switch {
		case free < 0 && it.size <= it.min:
			it.frozen = true
		case free > 0 && it.max >= 0 && it.size >= it.max:
			it.frozen = true
		}
	}

	const maxPasses = 8

	for pass := 0; pass < maxPasses && free != 0; pass++ {
		total := 0
		for i := range items {
			it := &items[i]
			if it.frozen {
				it.weight = 0
				continue
			}

			if free > 0 {
				it.weight = it.box.Style.Grow
			} else {

				it.weight = it.box.Style.Shrink * it.size
			}
			total += it.weight
		}
		if total == 0 {
			break
		}

		share := free
		assigned, cumulative := 0, 0

		for i := range items {
			it := &items[i]
			if it.weight == 0 {
				continue
			}

			cumulative += it.weight
			target := share * cumulative / total
			delta := target - assigned
			assigned = target

			want := it.size + delta
			switch {
			case want < it.min:
				want = it.min
				it.frozen = true
			case it.max >= 0 && want > it.max:
				want = it.max
				it.frozen = true
			}

			free -= want - it.size
			it.size = want
		}
	}

	for i := range items {
		it := items[i]
		if it.size != mainOuter(it.box, horizontal) {
			resizeMain(it.box, it.size, horizontal, contentW, availH)
		}
	}
}

func marginMain(c *Box, horizontal bool) int {
	if horizontal {
		return c.Style.Margin.Horizontal()
	}
	return c.Style.Margin.Vertical()
}

func minMain(c *Box, horizontal bool) int {
	if c.Style.Overflow != style.OverflowVisible {
		return 0
	}

	if horizontal {
		if w, ok := c.Style.Width.Resolve(-1); ok {
			return w
		}
		return minContentWidth(c) + childInsets(c)
	}

	return minContentHeight(c)
}

func minContentHeight(b *Box) int {
	if b.Node == nil {
		return b.ContentH
	}
	if b.Node.Canvas != nil {
		return b.ContentH
	}
	if b.Style.Overflow != style.OverflowVisible {
		return 0
	}
	if h, ok := b.Style.Height.Resolve(-1); ok {
		return h
	}

	flow := flowChildren(b)
	inner := 0
	stacked := b.Style.Display != style.DisplayFlex || b.Style.Direction == style.Column

	for _, c := range flow {
		ch := minContentHeight(c) + c.Style.Margin.Vertical()
		if stacked {
			inner += ch
		} else if ch > inner {
			inner = ch
		}
	}
	if stacked && len(flow) > 1 {
		inner += b.Style.RowGap() * (len(flow) - 1)
	}

	return inner + b.Ins.Vertical()
}

func maxMain(c *Box, horizontal bool, contentW, availH int) int {
	d := c.Style.MaxHeight
	base := availH
	if horizontal {
		d, base = c.Style.MaxWidth, contentW
	}
	if v, ok := d.Resolve(base); ok {
		return v + marginMain(c, horizontal)
	}
	return -1
}

func minContentWidth(b *Box) int {
	if b.hasMinWidth {
		return b.minContent
	}

	w := minContentWidthUncached(b)
	b.minContent, b.hasMinWidth = w, true
	return w
}

func minContentWidthUncached(b *Box) int {
	if b.Node != nil && b.Node.Canvas != nil {
		return b.Node.Canvas.W
	}

	if b.Node == nil {
		w := 0
		for i := range b.Source {
			run := &b.Source[i]
			eachToken(run.Text, run.Style != nil && run.Style.Pre, func(tok token) {
				if tok.w > w {
					w = tok.w
				}
			})
		}
		return w
	}

	if w, ok := b.Style.Width.Resolve(-1); ok {
		return w - b.Ins.Horizontal()
	}

	inner := 0
	row := b.Style.Display == style.DisplayFlex && b.Style.Direction == style.Row

	for _, c := range flowChildren(b) {
		cw := minContentWidth(c) + c.Style.Margin.Horizontal() + childInsets(c)
		if row {
			inner += cw
		} else if cw > inner {
			inner = cw
		}
	}

	return inner
}

func mainOuter(c *Box, horizontal bool) int {
	if horizontal {
		return c.OuterW + c.Style.Margin.Horizontal()
	}
	return c.OuterH + c.Style.Margin.Vertical()
}

func resizeMain(c *Box, outerWithMargin int, horizontal bool, contentW, availH int) {
	st := c.Style

	if horizontal {
		w := outerWithMargin - st.Margin.Horizontal()
		if w < 0 {
			w = 0
		}
		measureOuterW(c, w, contentW, availH)
		return
	}

	outerH := outerWithMargin - st.Margin.Vertical()
	ins := c.Ins
	contentH := clampDim(outerH-ins.Vertical(), st.MinHeight, st.MaxHeight, availH, ins.Vertical())
	measureOuterH(c, contentH+ins.Vertical(), contentW)
}

func flowChildren(b *Box) []*Box { return b.Flow }

func intrinsicWidth(b *Box, limit int) int {
	if b.hasIntrinsic {
		return capTo(b.intrinsic, limit)
	}

	st := b.Style
	w := 0

	if b.Node != nil && b.Node.Canvas != nil {
		b.intrinsic, b.hasIntrinsic = b.Node.Canvas.W, true
		return capTo(b.intrinsic, limit)
	}

	if b.Node == nil {
		w = maxLineWidth(b.Source)
	} else if fixed, ok := st.Width.Resolve(limit); ok {
		w = fixed - b.Ins.Horizontal()
	} else {
		inner := 0
		if st.Display == style.DisplayFlex && st.Direction == Row() {
			for i, c := range flowChildren(b) {
				if i > 0 {
					inner += st.ColumnGap()
				}
				inner += intrinsicWidth(c, -1) + c.Style.Margin.Horizontal() + childInsets(c)
			}
		} else {
			for _, c := range b.Flow {
				cw := intrinsicWidth(c, -1) + c.Style.Margin.Horizontal() + childInsets(c)
				if cw > inner {
					inner = cw
				}
			}
		}
		w = inner
	}

	if t := st.Title; t != "" {
		if tw := render.StringWidth(t) + 2; tw > w {
			w = tw
		}
	}

	if w < 0 {
		w = 0
	}

	b.intrinsic, b.hasIntrinsic = w, true
	return capTo(w, limit)
}

func naturalWidth(b *Box) int {
	if b.Node != nil && b.Node.Canvas != nil {
		return b.Node.Canvas.W
	}
	if b.Node == nil {
		return maxLineWidth(b.Source)
	}

	flow := flowChildren(b)
	row := b.Style.Display == style.DisplayFlex && b.Style.Direction == style.Row

	widest := 0
	for i, c := range flow {
		own := intrinsicWidth(c, -1) + c.Style.Margin.Horizontal() + childInsets(c)

		if row {
			if i > 0 {
				widest += b.Style.ColumnGap()
			}
			widest += own
			continue
		}

		if own > widest {
			widest = own
		}
	}

	return widest
}

func childInsets(c *Box) int {
	if c.Node == nil {
		return 0
	}
	return c.Ins.Horizontal()
}

func capTo(w, limit int) int {
	if limit >= 0 && w > limit {
		return limit
	}
	return w
}

func resolveMain(d style.Dim, avail, margin, insets int) (int, bool) {
	v, ok := d.Resolve(avail)
	if !ok {
		return 0, false
	}
	if d.Kind == style.DimPercent || d.Kind == style.DimFull {
		v -= margin
	}
	v -= insets
	if v < 0 {
		v = 0
	}
	return v, true
}

func clampDim(v int, minD, maxD style.Dim, avail, insets int) int {
	if m, ok := minD.Resolve(avail); ok {
		if lo := m - insets; v < lo {
			v = lo
		}
	}
	if m, ok := maxD.Resolve(avail); ok {
		if hi := m - insets; v > hi {
			v = hi
		}
	}
	if v < 0 {
		return 0
	}
	return v
}

func basisWidth(c *Box, parentContentW int) (int, bool) {
	if w, ok := c.Style.Basis.Resolve(parentContentW); ok {
		return w, true
	}
	if w, ok := c.Style.Width.Resolve(parentContentW); ok {
		return w, true
	}
	return 0, false
}

func Row() style.Direction { return style.Row }
