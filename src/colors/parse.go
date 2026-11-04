package colors

import (
	"strconv"
	"strings"
)

var Named = map[string]Color{
	"black":     Black,
	"white":     White,
	"gray":      Gray,
	"grey":      Gray,
	"red":       Red,
	"green":     Green,
	"blue":      Blue,
	"yellow":    Yellow,
	"orange":    Orange,
	"purple":    Purple,
	"pink":      Pink,
	"cyan":      Cyan,
	"teal":      Teal,
	"indigo":    Indigo,
	"primary":   Primary,
	"secondary": Secondary,
	"success":   Success,
	"danger":    Danger,
	"warning":   Warning,
	"info":      Info,
	"light":     RGB(248, 249, 250),
	"dark":      RGB(33, 37, 41),
	"base":      Base,
	"surface":   Surface,
	"overlay":   Overlay,
	"text":      Text,
	"lavender":  Lavender,
	"sky":       Sky,
	"mauve":     Mauve,
	"peach":     Peach,
	"rosewater": Rosewater,
}

func Parse(s string) (Color, bool) {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return Color{}, false
	}

	if c, ok := Token(s); ok {
		return c, true
	}

	if c, ok := Named[s]; ok {
		return c, true
	}

	if strings.HasPrefix(s, "#") {
		return parseHex(s[1:])
	}

	if strings.HasPrefix(s, "rgb(") && strings.HasSuffix(s, ")") {
		return parseFunc(s[4 : len(s)-1])
	}

	return Color{}, false
}

func parseHex(h string) (Color, bool) {
	switch len(h) {
	case 3:
		r, ok1 := hexNibble(h[0])
		g, ok2 := hexNibble(h[1])
		b, ok3 := hexNibble(h[2])
		if !ok1 || !ok2 || !ok3 {
			return Color{}, false
		}
		return Color{r * 17, g * 17, b * 17}, true

	case 6:
		v, err := strconv.ParseUint(h, 16, 32)
		if err != nil {
			return Color{}, false
		}
		return Color{uint8(v >> 16), uint8(v >> 8), uint8(v)}, true
	}

	return Color{}, false
}

func parseFunc(body string) (Color, bool) {
	parts := strings.Split(body, ",")
	if len(parts) != 3 {
		return Color{}, false
	}

	var out [3]uint8
	for i, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return Color{}, false
		}
		if n < 0 {
			n = 0
		}
		if n > 255 {
			n = 255
		}
		out[i] = uint8(n)
	}

	return Color{out[0], out[1], out[2]}, true
}

func hexNibble(c byte) (uint8, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

func Mix(a, b Color, t float64) Color {
	if t <= 0 {
		return a
	}
	if t >= 1 {
		return b
	}
	lerp := func(x, y uint8) uint8 {
		return uint8(float64(x) + (float64(y)-float64(x))*t + 0.5)
	}
	return Color{lerp(a.R, b.R), lerp(a.G, b.G), lerp(a.B, b.B)}
}

func (c Color) Hex() string {
	const digits = "0123456789abcdef"
	buf := []byte{'#', 0, 0, 0, 0, 0, 0}
	v := []uint8{c.R, c.G, c.B}
	for i, x := range v {
		buf[1+i*2] = digits[x>>4]
		buf[2+i*2] = digits[x&0x0f]
	}
	return string(buf)
}
