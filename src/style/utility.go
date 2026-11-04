package style

import (
	"strconv"
	"strings"

	"github.com/orayew2002/gtui/src/colors"
)

const Prefix = "box"

func Apply(s *Style, class string) bool {
	if class == "" || class == Prefix {
		return true
	}

	rest := trimPrefix(class)

	key, arb, hasArb := splitArbitrary(rest)
	if hasArb {
		return applyArbitrary(s, key, arb)
	}

	if applyKeyword(s, rest) {
		return true
	}
	return applyValued(s, rest)
}

// Classes apply left to right, so a later class overrides an earlier one and a
// caller's class always wins over a component's defaults.
func ApplyAll(s *Style, classes []string) []string {
	var unknown []string
	for _, c := range classes {
		if c == "" {
			continue
		}
		if !Apply(s, c) {
			unknown = append(unknown, c)
		}
	}
	return unknown
}

func trimPrefix(class string) string {
	if strings.HasPrefix(class, Prefix+"-") {
		return class[len(Prefix)+1:]
	}
	return class
}

func splitArbitrary(rest string) (key, value string, ok bool) {
	open := strings.IndexByte(rest, '[')
	if open <= 0 || !strings.HasSuffix(rest, "]") {
		return rest, "", false
	}
	return strings.TrimSuffix(rest[:open], "-"), rest[open+1 : len(rest)-1], true
}

func applyKeyword(s *Style, rest string) bool {
	switch rest {

	case "block":
		s.Display = DisplayBlock
	case "inline":
		s.Display = DisplayInline
	case "none":
		s.Display = DisplayNone
	case "flex", "row":
		s.Display, s.Direction = DisplayFlex, Row
	case "col", "column":
		s.Display, s.Direction = DisplayFlex, Column

	case "border":
		s.Border, s.BorderSides = BorderSharp, allSides()
	case "rounded":
		s.Border, s.BorderSides = BorderRounded, allSides()
	case "thick":
		s.Border, s.BorderSides = BorderThick, allSides()
	case "double":
		s.Border, s.BorderSides = BorderDouble, allSides()
	case "dashed":
		s.Border, s.BorderSides = BorderDashed, allSides()
	case "soft":
		s.Border, s.BorderSides = BorderSoft, allSides()
	case "tight":
		s.Border, s.BorderSides = BorderTight, allSides()
	case "pill", "rounded-pill":
		// A pill only brackets its content, so it claims no vertical chrome
		// and a one-line pill stays one row tall.
		s.Border, s.BorderSides = BorderPill, Edges{Left: 1, Right: 1}
	case "borderless", "border-none", "rounded-none":
		s.Border, s.BorderSides = BorderNone, Edges{}
	case "shadow":
		s.Shadow = true
	case "shadow-none":
		s.Shadow = false

	case "border-t":
		s.BorderSides = Edges{Top: 1}
		s.defaultBorder()
	case "border-r":
		s.BorderSides = Edges{Right: 1}
		s.defaultBorder()
	case "border-b":
		s.BorderSides = Edges{Bottom: 1}
		s.defaultBorder()
	case "border-l":
		s.BorderSides = Edges{Left: 1}
		s.defaultBorder()
	case "border-x":
		s.BorderSides = Edges{Left: 1, Right: 1}
		s.defaultBorder()
	case "border-y":
		s.BorderSides = Edges{Top: 1, Bottom: 1}
		s.defaultBorder()

	case "clip":
		s.Overflow = OverflowHidden
	case "visible":
		s.Overflow = OverflowVisible
	case "scroll":
		s.Overflow, s.Scroll, s.Focusable = OverflowScroll, true, true
	case "scroll-x":
		s.ScrollX, s.Focusable = true, true
		if s.Overflow == OverflowVisible {
			s.Overflow = OverflowHidden
		}

	case "justify-start":
		s.Justify = JustifyStart
	case "justify-center", "center-x":
		s.Justify = JustifyCenter
	case "justify-end":
		s.Justify = JustifyEnd
	case "justify-between", "between":
		s.Justify = JustifyBetween
	case "justify-around", "around":
		s.Justify = JustifyAround
	case "justify-evenly", "evenly":
		s.Justify = JustifyEvenly

	case "items-start":
		s.Align = AlignStart
	case "items-center", "center-y":
		s.Align = AlignCenter
	case "items-end":
		s.Align = AlignEnd
	case "items-stretch":
		s.Align = AlignStretch
	case "self-start":
		s.AlignSelf, s.HasAlignSelf = AlignStart, true
	case "self-center":
		s.AlignSelf, s.HasAlignSelf = AlignCenter, true
	case "self-end":
		s.AlignSelf, s.HasAlignSelf = AlignEnd, true
	case "self-stretch":
		s.AlignSelf, s.HasAlignSelf = AlignStretch, true
	case "flex-wrap":
		s.FlexWrap = true
	case "flex-nowrap":
		s.FlexWrap = false
	case "center":
		s.Display, s.Justify, s.Align = DisplayFlex, JustifyCenter, AlignCenter

	case "grow":
		s.Grow = 1
	case "grow-0":
		s.Grow = 0
	case "shrink":
		s.Shrink = 1
	case "shrink-0":
		s.Shrink = 0

	case "w-full":
		s.Width = Full()
	case "h-full":
		s.Height = Full()
	case "w-auto":
		s.Width = Auto()
	case "h-auto":
		s.Height = Auto()
	case "full":
		s.Width, s.Height = Full(), Full()

	case "text-left":
		s.TextAlign = AlignStart
	case "text-center":
		s.TextAlign = AlignCenter
	case "text-right":
		s.TextAlign = AlignEnd
	case "nowrap":
		s.NoWrap = true
	case "wrap":
		s.NoWrap = false
	case "pre":
		s.Pre, s.NoWrap = true, true
	case "truncate", "text-truncate":
		s.NoWrap, s.Overflow = true, OverflowHidden
	case "uppercase", "text-uppercase":
		s.Transform = TransformUpper
	case "lowercase", "text-lowercase":
		s.Transform = TransformLower
	case "capitalize", "text-capitalize":
		s.Transform = TransformTitle
	case "normal-case":
		s.Transform = TransformNone
	case "bold":
		s.Attr |= AttrBold
	case "dim":
		s.Attr |= AttrDim
	case "italic":
		s.Attr |= AttrItalic
	case "underline":
		s.Attr |= AttrUnderline
	case "reverse":
		s.Attr |= AttrReverse
	case "strike":
		s.Attr |= AttrStrike
	case "blink":
		s.Attr |= AttrBlink
	case "opacity-50":
		s.Attr |= AttrDim
	case "opacity-100":
		s.Attr &^= AttrDim
	case "normal":
		s.Attr = 0

	case "relative", "static":
		s.Position = PositionStatic
	case "absolute", "fixed":
		s.Position = PositionAbsolute

	case "focusable":
		s.Focusable = true
	case "key-hide":
		s.HideShortcut = true
	case "key-show":
		s.HideShortcut = false

	case "title-left":
		s.TitleAlign = AlignStart
	case "title-center":
		s.TitleAlign = AlignCenter
	case "title-right":
		s.TitleAlign = AlignEnd

	default:
		return false
	}

	return true
}

func applyValued(s *Style, rest string) bool {
	prop, value, ok := cut(rest)
	if !ok {
		return false
	}

	switch prop {

	case "p", "px", "py", "pt", "pr", "pb", "pl":
		n, ok := atoi(value)
		if !ok {
			return false
		}
		setEdge(&s.Padding, prop[1:], n)
		return true

	case "m", "mx", "my", "mt", "mr", "mb", "ml":
		n, ok := atoi(value)
		if !ok {
			return false
		}
		setEdge(&s.Margin, prop[1:], n)
		return true

	case "gap":
		n, ok := atoi(value)
		if !ok {
			return false
		}
		s.Gap, s.HasGapX, s.HasGapY = n, false, false
		return true

	case "gap-x":
		n, ok := atoi(value)
		if !ok {
			return false
		}
		s.GapX, s.HasGapX = n, true
		return true

	case "gap-y":
		n, ok := atoi(value)
		if !ok {
			return false
		}
		s.GapY, s.HasGapY = n, true
		return true

	case "w":
		return setDim(&s.Width, value)
	case "h":
		return setDim(&s.Height, value)
	case "min-w":
		return setDim(&s.MinWidth, value)
	case "min-h":
		return setDim(&s.MinHeight, value)
	case "max-w":
		return setDim(&s.MaxWidth, value)
	case "max-h":
		return setDim(&s.MaxHeight, value)
	case "basis":
		return setDim(&s.Basis, value)

	case "top":
		return setDim(&s.Top, value)
	case "right":
		return setDim(&s.Right, value)
	case "bottom":
		return setDim(&s.Bottom, value)
	case "left":
		return setDim(&s.Left, value)

	case "grow":
		n, ok := atoi(value)
		if !ok {
			return false
		}
		s.Grow = n
		return true
	case "shrink":
		n, ok := atoi(value)
		if !ok {
			return false
		}
		s.Shrink = n
		return true
	case "z":
		n, ok := atoi(value)
		if !ok {
			return false
		}
		s.ZIndex = n
		return true

	case "key":
		if value == "" {
			return false
		}
		s.Shortcut = value
		s.Focusable = true
		return true

	case "bg":
		if value == "transparent" || value == "none" {
			s.Background = Paint{}
			return true
		}
		c, ok := colors.Parse(value)
		if !ok {
			return false
		}
		s.Background = P(c)
		return true

	case "text":
		c, ok := colors.Parse(value)
		if !ok {
			return false
		}
		s.Foreground = P(c)
		return true

	case "border":
		c, ok := colors.Parse(value)
		if !ok {
			return false
		}
		s.BorderColor = P(c)
		s.defaultBorder()
		return true
	}

	return false
}

func applyArbitrary(s *Style, key, value string) bool {
	switch key {
	case "bg", "text", "border":
		return applyValued(s, key+"-"+value)
	case "w", "h", "min-w", "min-h", "max-w", "max-h", "p", "px", "py", "pt",
		"pr", "pb", "pl", "m", "mx", "my", "mt", "mr", "mb", "ml",
		"gap", "gap-x", "gap-y",
		"top", "right", "bottom", "left", "z", "grow", "shrink", "basis":
		return applyValued(s, key+"-"+value)
	}
	return false
}

func (s *Style) defaultBorder() {
	if s.Border == BorderNone {
		s.Border = BorderSharp
	}
	if s.BorderSides == (Edges{}) {
		s.BorderSides = allSides()
	}
}

func allSides() Edges { return Edges{Top: 1, Right: 1, Bottom: 1, Left: 1} }

func cut(rest string) (prop, value string, ok bool) {
	i := strings.LastIndexByte(rest, '-')
	if i <= 0 || i == len(rest)-1 {
		return "", "", false
	}
	return rest[:i], rest[i+1:], true
}

func setEdge(e *Edges, side string, n int) {
	switch side {
	case "":
		*e = Edges{n, n, n, n}
	case "x":
		e.Left, e.Right = n, n
	case "y":
		e.Top, e.Bottom = n, n
	case "t":
		e.Top = n
	case "r":
		e.Right = n
	case "b":
		e.Bottom = n
	case "l":
		e.Left = n
	}
}

func setDim(d *Dim, value string) bool {
	switch value {
	case "auto":
		*d = Auto()
		return true
	case "full":
		*d = Full()
		return true
	}

	if i := strings.IndexByte(value, '/'); i > 0 {
		num, ok1 := atoi(value[:i])
		den, ok2 := atoi(value[i+1:])
		if !ok1 || !ok2 || den == 0 {
			return false
		}
		*d = Percent(num * 100 / den)
		return true
	}

	if strings.HasSuffix(value, "%") {
		n, ok := atoi(strings.TrimSuffix(value, "%"))
		if !ok {
			return false
		}
		*d = Percent(n)
		return true
	}

	n, ok := atoi(value)
	if !ok {
		return false
	}
	*d = Cells(n)
	return true
}

func atoi(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}
