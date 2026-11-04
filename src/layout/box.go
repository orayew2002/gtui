package layout

import (
	"github.com/orayew2002/gtui/src/dom"
	"github.com/orayew2002/gtui/src/render"
	"github.com/orayew2002/gtui/src/style"
)

type Run struct {
	Text  string
	X     int
	W     int
	Style *style.Style
	Node  *dom.Node
}

type Box struct {
	Node     *dom.Node
	Style    style.Style
	Children []*Box
	Parent   *Box

	Ins style.Edges

	Flow []*Box
	Abs  []*Box

	FlexLines [][]*Box

	OuterW, OuterH     int
	ContentW, ContentH int

	Frame   render.Rect
	Content render.Rect

	Clip render.Rect

	Lines  []Line
	Source []Run

	ScrollMax  int
	ScrollMaxX int

	intrinsic    int
	hasIntrinsic bool
	minContent   int
	hasMinWidth  bool

	wrapWidth int
	wrapped   bool

	runBuf  []Run
	lineBuf []Line
	spans   []span
}

func (b *Box) IsInline() bool { return b.Node == nil || b.Style.Display == style.DisplayInline }

func (b *Box) Walk(fn func(*Box) bool) {
	if b == nil || !fn(b) {
		return
	}
	for _, c := range b.Children {
		c.Walk(fn)
	}
}

func (b *Box) HitTest(x, y int) *Box {
	if b == nil || !b.Frame.Contains(x, y) || !b.Clip.Contains(x, y) {
		return nil
	}

	for i := len(b.Children) - 1; i >= 0; i-- {
		if hit := b.Children[i].HitTest(x, y); hit != nil {
			return hit
		}
	}

	if b.Node == nil {
		return nil
	}
	return b
}
