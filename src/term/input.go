package term

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

type Event interface{ event() }

type Key struct {
	Name string
	Rune rune
	Alt  bool
	Ctrl bool
}

func (Key) event() {}

type MouseKind uint8

const (
	MousePress MouseKind = iota
	MouseRelease
	MouseMove
	MouseWheelUp
	MouseWheelDown
)

type Mouse struct {
	Kind   MouseKind
	Button int
	X, Y   int
}

func (Mouse) event() {}

func decode(buf []byte) (Event, int) {
	if len(buf) == 0 {
		return nil, 0
	}

	c := buf[0]

	if c == 0x1b {
		if len(buf) == 1 {
			return Key{Name: "esc"}, 1
		}

		switch buf[1] {
		case '[':
			return decodeCSI(buf)
		case 'O':
			return decodeSS3(buf)
		case 0x1b:
			return Key{Name: "esc"}, 1
		}

		ev, used := decode(buf[1:])
		if used == 0 {
			return nil, 0
		}
		if key, ok := ev.(Key); ok {
			key.Alt = true
			key.Name = "alt+" + key.Name
			return key, used + 1
		}
		return ev, used + 1
	}

	switch c {
	case '\r', '\n':
		return Key{Name: "enter", Rune: '\r'}, 1
	case '\t':
		return Key{Name: "tab", Rune: '\t'}, 1
	case 0x7f:
		return Key{Name: "backspace"}, 1
	case 0:
		return Key{Name: "ctrl+space", Ctrl: true}, 1
	}

	if c < 0x20 {
		return Key{Name: "ctrl+" + string(rune('a'+c-1)), Ctrl: true, Rune: rune(c)}, 1
	}

	r, size := utf8.DecodeRune(buf)
	if r == utf8.RuneError && size <= 1 {
		if !utf8.FullRune(buf) {
			return nil, 0
		}
		return nil, 1
	}

	return Key{Name: strings.ToLower(string(r)), Rune: r}, size
}

func decodeCSI(buf []byte) (Event, int) {

	end := -1
	for i := 2; i < len(buf); i++ {
		if b := buf[i]; b >= 0x40 && b <= 0x7e {
			end = i
			break
		}
	}
	if end < 0 {
		return nil, 0
	}

	params := string(buf[2:end])
	final := buf[end]
	used := end + 1

	if final == 'M' || final == 'm' {
		if strings.HasPrefix(params, "<") {
			return decodeSGRMouse(params[1:], final == 'm'), used
		}
	}

	switch final {
	case 'A':
		return modified("up", params), used
	case 'B':
		return modified("down", params), used
	case 'C':
		return modified("right", params), used
	case 'D':
		return modified("left", params), used
	case 'H':
		return Key{Name: "home"}, used
	case 'F':
		return Key{Name: "end"}, used
	case 'Z':
		return Key{Name: "shift+tab"}, used
	case '~':
		return tildeKey(params), used
	}

	return nil, used
}

func decodeSS3(buf []byte) (Event, int) {
	if len(buf) < 3 {
		return nil, 0
	}

	switch buf[2] {
	case 'A':
		return Key{Name: "up"}, 3
	case 'B':
		return Key{Name: "down"}, 3
	case 'C':
		return Key{Name: "right"}, 3
	case 'D':
		return Key{Name: "left"}, 3
	case 'H':
		return Key{Name: "home"}, 3
	case 'F':
		return Key{Name: "end"}, 3
	case 'P':
		return Key{Name: "f1"}, 3
	case 'Q':
		return Key{Name: "f2"}, 3
	case 'R':
		return Key{Name: "f3"}, 3
	case 'S':
		return Key{Name: "f4"}, 3
	}
	return nil, 3
}

func modified(name, params string) Key {
	key := Key{Name: name}

	parts := strings.Split(params, ";")
	if len(parts) < 2 {
		return key
	}

	mod, err := strconv.Atoi(parts[1])
	if err != nil || mod < 2 {
		return key
	}

	bits := mod - 1
	if bits&4 != 0 {
		key.Name = "ctrl+" + key.Name
		key.Ctrl = true
	}
	if bits&2 != 0 {
		key.Name = "alt+" + key.Name
		key.Alt = true
	}
	if bits&1 != 0 {
		key.Name = "shift+" + key.Name
	}
	return key
}

func tildeKey(params string) Event {
	code := params
	if i := strings.IndexByte(params, ';'); i >= 0 {
		code = params[:i]
	}

	names := map[string]string{
		"1": "home", "2": "insert", "3": "delete", "4": "end",
		"5": "pageup", "6": "pagedown", "7": "home", "8": "end",
		"11": "f1", "12": "f2", "13": "f3", "14": "f4", "15": "f5",
		"17": "f6", "18": "f7", "19": "f8", "20": "f9", "21": "f10",
		"23": "f11", "24": "f12",
	}

	name, ok := names[code]
	if !ok {
		return nil
	}
	return modified(name, params)
}

func decodeSGRMouse(params string, release bool) Event {
	fields := strings.Split(params, ";")
	if len(fields) != 3 {
		return nil
	}

	btn, err1 := strconv.Atoi(fields[0])
	col, err2 := strconv.Atoi(fields[1])
	row, err3 := strconv.Atoi(fields[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return nil
	}

	m := Mouse{X: col - 1, Y: row - 1, Button: btn & 3}

	switch {
	case btn&64 != 0:
		if btn&1 == 0 {
			m.Kind = MouseWheelUp
		} else {
			m.Kind = MouseWheelDown
		}
	case btn&32 != 0:
		m.Kind = MouseMove
	case release:
		m.Kind = MouseRelease
	default:
		m.Kind = MousePress
	}

	return m
}
