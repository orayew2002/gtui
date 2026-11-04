package html

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

var entities = map[string]rune{
	"amp":    '&',
	"lt":     '<',
	"gt":     '>',
	"quot":   '"',
	"apos":   '\'',
	"nbsp":   ' ',
	"copy":   '©',
	"reg":    '®',
	"trade":  '™',
	"hellip": '…',
	"mdash":  '—',
	"ndash":  '–',
	"laquo":  '«',
	"raquo":  '»',
	"larr":   '←',
	"rarr":   '→',
	"uarr":   '↑',
	"darr":   '↓',
	"middot": '·',
	"bull":   '•',
	"check":  '✓',
	"cross":  '✗',
	"star":   '★',
	"deg":    '°',
	"euro":   '€',
	"pound":  '£',
	"yen":    '¥',
	"sect":   '§',
	"para":   '¶',
	"dagger": '†',
}

func Unescape(s string) string {
	if !strings.ContainsRune(s, '&') {
		return s
	}

	var sb strings.Builder
	sb.Grow(len(s))

	for i := 0; i < len(s); {
		if s[i] != '&' {
			sb.WriteByte(s[i])
			i++
			continue
		}

		end := strings.IndexByte(s[i:], ';')
		if end < 0 || end > 32 {
			sb.WriteByte('&')
			i++
			continue
		}

		name := s[i+1 : i+end]
		if r, ok := decodeRef(name); ok {
			sb.WriteRune(r)
			i += end + 1
			continue
		}

		sb.WriteByte('&')
		i++
	}

	return sb.String()
}

func decodeRef(name string) (rune, bool) {
	if name == "" {
		return 0, false
	}

	if name[0] == '#' {
		digits, base := name[1:], 10
		if len(digits) > 1 && (digits[0] == 'x' || digits[0] == 'X') {
			digits, base = digits[1:], 16
		}
		v, err := strconv.ParseInt(digits, base, 32)
		if err != nil || v <= 0 || v > utf8.MaxRune {
			return 0, false
		}
		return rune(v), true
	}

	r, ok := entities[strings.ToLower(name)]
	return r, ok
}
