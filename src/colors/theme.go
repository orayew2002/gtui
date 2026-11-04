package colors

import "sync/atomic"

type Theme struct {
	Name string

	Canvas Color
	Panel  Color
	Raised Color
	Line   Color
	Muted  Color
	Ink    Color
	Accent Color

	Success Color
	Warning Color
	Danger  Color
	Info    Color
}

var Dark = Theme{
	Name:    "dark",
	Canvas:  RGB(11, 11, 11),
	Panel:   RGB(28, 28, 28),
	Raised:  RGB(36, 36, 36),
	Line:    RGB(46, 46, 46),
	Muted:   RGB(138, 139, 139),
	Ink:     RGB(232, 232, 232),
	Accent:  RGB(137, 220, 235),
	Success: RGB(126, 203, 145),
	Warning: RGB(233, 196, 106),
	Danger:  RGB(231, 87, 99),
	Info:    RGB(137, 180, 250),
}

var Light = Theme{
	Name:    "light",
	Canvas:  RGB(245, 247, 248),
	Panel:   RGB(255, 255, 255),
	Raised:  RGB(243, 244, 246),
	Line:    RGB(221, 224, 228),
	Muted:   RGB(109, 110, 110),
	Ink:     RGB(31, 32, 36),
	Accent:  RGB(23, 105, 170),
	Success: RGB(25, 135, 84),
	Warning: RGB(176, 128, 9),
	Danger:  RGB(200, 36, 51),
	Info:    RGB(13, 110, 253),
}

var (
	active atomic.Pointer[Theme]

	registry = map[string]Theme{
		"dark":  Dark,
		"light": Light,
	}
)

func Register(t Theme) {
	if t.Name == "" {
		return
	}
	registry[t.Name] = t
}

func ByName(name string) (Theme, bool) {
	t, ok := registry[name]
	return t, ok
}

func Names() []string {
	out := make([]string, 0, len(registry))
	for name := range registry {
		out = append(out, name)
	}

	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}

	return out
}

func Use(t Theme) {
	active.Store(&t)
}

func Active() Theme {
	if t := active.Load(); t != nil {
		return *t
	}
	return Dark
}

func Token(name string) (Color, bool) {
	t := Active()

	switch name {
	case "canvas":
		return t.Canvas, true
	case "panel":
		return t.Panel, true
	case "raised":
		return t.Raised, true
	case "line":
		return t.Line, true
	case "muted":
		return t.Muted, true
	case "ink":
		return t.Ink, true
	case "accent":
		return t.Accent, true
	case "success":
		return t.Success, true
	case "warning":
		return t.Warning, true
	case "danger":
		return t.Danger, true
	case "info":
		return t.Info, true
	}

	return Color{}, false
}
