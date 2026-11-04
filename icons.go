package gtui

import (
	"slices"
	"strings"
	"sync"

	"github.com/orayew2002/gtui/src/render"
)

// icons keeps one glyph per name so apps can use a stable vocabulary instead
// of hunting for code points. Every shipped glyph is exactly one cell wide so
// rows of icons line up.
var icons = map[string]string{
	"home": "⌂", "search": "⌕", "settings": "⚙", "user": "◍", "users": "⚇",
	"file": "▤", "folder": "▰", "folderOpen": "▱", "image": "▣", "video": "▶",
	"music": "♪", "sound": "♫", "mute": "⊘", "play": "▶", "pause": "⏸",
	"stop": "■", "next": "⏭", "prev": "⏮", "record": "●", "disk": "⊙",
	"check": "✓", "cross": "✗", "plus": "+", "minus": "−", "edit": "✎",
	"trash": "⌫", "save": "⤓", "upload": "⤒", "download": "⤓", "refresh": "⟳",
	"clock": "◷", "calendar": "▦", "bell": "✦", "flag": "⚑", "star": "★",
	"heart": "♥", "info": "i", "warning": "!", "error": "✗", "ok": "✓",
	"left": "‹", "right": "›", "up": "⌃", "down": "⌄", "back": "←",
	"forward": "→", "menu": "≡", "more": "⋯", "grid": "⌗", "list": "☰",
	"lock": "⚿", "unlock": "⚷", "eye": "◉", "hidden": "◌", "link": "⇗",
	"mail": "✉", "phone": "☎", "chat": "❝", "send": "➤", "pin": "⌖",
	"filter": "⚟", "sort": "⇅", "chart": "◫", "terminal": "▌", "code": "⌘",
	"wifi": "≋", "battery": "▰", "power": "⏻", "cloud": "☁", "key": "⚿",

	"arrowUp": "↑", "arrowDown": "↓", "arrowLeft": "←", "arrowRight": "→",
	"arrowUpRight": "↗", "arrowDownRight": "↘", "arrowUpLeft": "↖", "arrowDownLeft": "↙",
	"sortAsc": "▲", "sortDesc": "▼", "expand": "⤢", "collapse": "⤡",
	"pullRequest": "⇵", "shuffle": "⇄", "sync": "↻", "undo": "↶", "redo": "↷",
	"enter": "⏎", "eject": "⏏", "keyboard": "⌨", "delete": "⌦", "cancel": "⌧",

	"close": "✕", "question": "?", "alert": "⚠", "copy": "⧉", "paste": "⎘",
	"clipboard": "⌷", "pencil": "✐", "pen": "✒", "eraser": "⎚", "scissors": "✂",
	"attach": "⌇", "printer": "⎙", "camera": "⌾", "mic": "⚲", "flash": "⌁",
	"zoomIn": "⊕", "zoomOut": "⊖", "zoomReset": "⊝", "filterOff": "⚞",

	"circle": "○", "square": "◻", "squareFilled": "◼", "triangle": "△",
	"diamond": "◆", "diamondOutline": "◇", "hexagon": "⬡", "pentagon": "⬠",
	"bullet": "▪", "bulletOutline": "▫", "circleHalf": "◐", "target": "◎",
	"radioOn": "⦿", "radioOff": "◯",

	"book": "▥", "map": "▧", "files": "▨", "layers": "▩", "panel": "◧",
	"split": "◨", "card": "▢",

	"desktop": "▬", "laptop": "▭", "mobile": "▯", "server": "⌸", "database": "⛁",
	"cpu": "⊞", "memory": "⊟", "usb": "⚯", "bluetooth": "ᛒ", "network": "⍟",

	"volumeLow": "▁", "volumeMid": "▄", "volumeHigh": "█",
	"batteryEmpty": "░", "batteryMid": "▒", "batteryFull": "▓", "batteryOutline": "▱",
	"wifiOff": "≁", "wifiLow": "∼", "wifiMid": "≈", "wifiHigh": "≋",

	"rewind": "«", "fastForward": "»", "speaker": "⊲", "headphones": "◠",

	"userOnline": "☻", "userOffline": "☺", "bellRing": "⍾", "bellOff": "✧",
	"award": "✪", "crown": "♔", "starOutline": "☆", "heartOutline": "♡",

	"sun": "☀", "moon": "☾", "snow": "❄", "rain": "⛆", "umbrella": "☂",
	"fire": "♨", "drop": "⬩", "leaf": "❧", "flower": "⚘", "tree": "♣",

	"globe": "⊜", "compass": "✥", "shield": "⛨", "beaker": "⚗", "atom": "⚛",
	"bug": "⚉", "rocket": "⟰", "puzzle": "⬔", "gift": "✤", "cart": "⊓",
	"coin": "⊚", "tag": "⬟", "bookmark": "⚐", "graduation": "⏢", "truck": "⛟",

	"chartBar": "▊", "chartLine": "∿", "chartPie": "◔", "chartArea": "◣",
	"branch": "⑂", "merge": "⑃", "commit": "◈", "diff": "±",
}

var iconsMu sync.RWMutex

// RegisterIcon adds or replaces a named glyph. Glyphs that are zero width or
// wider than two cells are ignored, because every caller sizes an icon as one
// cell of text.
func RegisterIcon(name, glyph string) {
	if name == "" || glyph == "" {
		return
	}
	if w := render.StringWidth(glyph); w < 1 || w > 2 {
		return
	}

	iconsMu.Lock()
	icons[name] = glyph
	iconsMu.Unlock()
}

// IconGlyph returns the glyph registered under name.
func IconGlyph(name string) (string, bool) {
	iconsMu.RLock()
	glyph, ok := icons[name]
	iconsMu.RUnlock()

	return glyph, ok
}

// Icon renders one named glyph. Unknown names fall back to a dot so a typo is
// visible instead of invisible.
func Icon(class, name string) Component {
	glyph, ok := IconGlyph(name)
	if !ok {
		glyph = "·"
	}

	return Text(join("nowrap", class), glyph)
}

// IconText pairs a glyph with a label, the way most UI rows want it.
func IconText(class, name, label string) Component {
	return Row(join("gap-1 nowrap", class),
		Icon("", name),
		Text("nowrap", label),
	)
}

// IconNames lists every icon the framework ships plus the registered ones,
// sorted.
func IconNames() []string {
	iconsMu.RLock()
	out := make([]string, 0, len(icons))
	for name := range icons {
		out = append(out, name)
	}
	iconsMu.RUnlock()

	slices.Sort(out)

	return out
}

// Equalizer draws a music-style level meter that moves with tick.
func Equalizer(class string, bars, tick int) Component {
	if bars < 1 {
		bars = 12
	}

	levels := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	period := len(levels) * 2

	var sb strings.Builder
	for i := range bars {
		step := (((i*3 + tick*2) % period) + period) % period
		if step >= len(levels) {
			step = period - step - 1
		}
		sb.WriteRune(levels[step])
	}

	return Text(join("nowrap", class), sb.String())
}
