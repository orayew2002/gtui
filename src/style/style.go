package style

import "github.com/orayew2002/gtui/src/colors"

type Display uint8

const (
	DisplayBlock Display = iota
	DisplayFlex
	DisplayInline
	DisplayNone
)

type Direction uint8

const (
	Row Direction = iota
	Column
)

type Justify uint8

const (
	JustifyStart Justify = iota
	JustifyCenter
	JustifyEnd
	JustifyBetween
	JustifyAround
	JustifyEvenly
)

type Align uint8

const (
	AlignStart Align = iota
	AlignCenter
	AlignEnd
	AlignStretch
)

type Overflow uint8

const (
	OverflowVisible Overflow = iota
	OverflowHidden
	OverflowScroll
)

type Position uint8

const (
	PositionStatic Position = iota
	PositionAbsolute
)

type BorderStyle uint8

const (
	BorderNone BorderStyle = iota
	BorderSharp
	BorderRounded
	BorderThick
	BorderDouble
	BorderSoft
	BorderTight
	BorderDashed
)

type DimKind uint8

const (
	DimAuto DimKind = iota
	DimCells
	DimPercent
	DimFull
)

type Dim struct {
	Kind DimKind
	Val  int
}

func Auto() Dim            { return Dim{Kind: DimAuto} }
func Cells(n int) Dim      { return Dim{Kind: DimCells, Val: n} }
func Percent(n int) Dim    { return Dim{Kind: DimPercent, Val: n} }
func Full() Dim            { return Dim{Kind: DimFull} }
func (d Dim) IsAuto() bool { return d.Kind == DimAuto }

func (d Dim) Resolve(parent int) (int, bool) {
	switch d.Kind {
	case DimCells:
		return d.Val, true
	case DimPercent:
		if parent < 0 {
			return 0, false
		}
		return parent * d.Val / 100, true
	case DimFull:
		if parent < 0 {
			return 0, false
		}
		return parent, true
	}
	return 0, false
}

type Edges struct {
	Top, Right, Bottom, Left int
}

func (e Edges) Horizontal() int { return e.Left + e.Right }
func (e Edges) Vertical() int   { return e.Top + e.Bottom }

type Paint struct {
	Color colors.Color
	Set   bool
}

func P(c colors.Color) Paint { return Paint{Color: c, Set: true} }

func (p Paint) Or(fallback Paint) Paint {
	if p.Set {
		return p
	}
	return fallback
}

type Attrs uint8

const (
	AttrBold Attrs = 1 << iota
	AttrDim
	AttrItalic
	AttrUnderline
	AttrReverse
	AttrStrike
	AttrBlink
)

func (a Attrs) Has(f Attrs) bool { return a&f != 0 }

type Style struct {
	Display   Display
	Position  Position
	Direction Direction
	Justify   Justify
	Align     Align
	Overflow  Overflow

	Width, Height       Dim
	MinWidth, MinHeight Dim
	MaxWidth, MaxHeight Dim
	Top, Right          Dim
	Bottom, Left        Dim

	Padding Edges
	Margin  Edges

	Border      BorderStyle
	BorderSides Edges
	BorderColor Paint
	Background  Paint
	Foreground  Paint

	Attr      Attrs
	TextAlign Align
	NoWrap    bool
	Pre       bool

	Gap    int
	Grow   int
	Shrink int
	Basis  Dim

	AlignSelf    Align
	HasAlignSelf bool
	FlexWrap     bool

	ZIndex    int
	Scroll    bool
	ScrollX   bool
	Focusable bool

	Title      string
	TitleAlign Align

	Shortcut     string
	HideShortcut bool
}

func New() Style {
	return Style{
		Width:  Auto(),
		Height: Auto(),
		Basis:  Auto(),
		Shrink: 1,
		Align:  AlignStretch,
	}
}

func (s Style) inherit() Style {
	out := New()
	out.Foreground = s.Foreground
	out.Attr = s.Attr
	out.TextAlign = s.TextAlign
	out.NoWrap = s.NoWrap
	out.Pre = s.Pre
	return out
}

func Inherit(parent Style) Style { return parent.inherit() }

func (s Style) HasBorder() bool {
	return s.Border != BorderNone &&
		(s.BorderSides.Top|s.BorderSides.Right|s.BorderSides.Bottom|s.BorderSides.Left) != 0
}

func (s Style) BorderWidths() Edges {
	if !s.HasBorder() {
		return Edges{}
	}
	clamp := func(v int) int {
		if v > 0 {
			return 1
		}
		return 0
	}
	return Edges{
		Top:    clamp(s.BorderSides.Top),
		Right:  clamp(s.BorderSides.Right),
		Bottom: clamp(s.BorderSides.Bottom),
		Left:   clamp(s.BorderSides.Left),
	}
}

func (s Style) Gutter() int {
	if s.Overflow == OverflowScroll {
		return 1
	}
	return 0
}

func (s Style) Insets() Edges {
	b := s.BorderWidths()
	return Edges{
		Top:    s.Padding.Top + b.Top,
		Right:  s.Padding.Right + b.Right + s.Gutter(),
		Bottom: s.Padding.Bottom + b.Bottom,
		Left:   s.Padding.Left + b.Left,
	}
}
