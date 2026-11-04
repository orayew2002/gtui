package gtui

import "strings"

// icons keeps one glyph per name so apps can use a stable vocabulary instead
// of hunting for code points.
var icons = map[string]string{
	"home": "⌂", "search": "⌕", "settings": "⚙", "user": "◍", "users": "⚇",
	"file": "▤", "folder": "▸", "folderOpen": "▾", "image": "▣", "video": "▶",
	"music": "♪", "sound": "♫", "mute": "⊘", "play": "▶", "pause": "⏸",
	"stop": "■", "next": "⏭", "prev": "⏮", "record": "●", "disk": "◉",
	"check": "✓", "cross": "✗", "plus": "＋", "minus": "−", "edit": "✎",
	"trash": "⌫", "save": "⤓", "upload": "⤒", "download": "⤓", "refresh": "⟳",
	"clock": "◷", "calendar": "▦", "bell": "✦", "flag": "⚑", "star": "★",
	"heart": "♥", "info": "i", "warning": "!", "error": "✗", "ok": "✓",
	"left": "‹", "right": "›", "up": "⌃", "down": "⌄", "back": "←",
	"forward": "→", "menu": "≡", "more": "⋯", "grid": "▦", "list": "☰",
	"lock": "⚿", "unlock": "⚿", "eye": "◉", "hidden": "◌", "link": "⇗",
	"mail": "✉", "phone": "☎", "chat": "❝", "send": "➤", "pin": "⌖",
	"filter": "⚟", "sort": "⇅", "chart": "◫", "terminal": "▌", "code": "⌘",
	"wifi": "≋", "battery": "▰", "power": "⏻", "cloud": "☁", "key": "⚿",
}

// Icon renders one named glyph. Unknown names fall back to a dot so a typo is
// visible instead of invisible.
func Icon(class, name string) Component {
	glyph, ok := icons[name]
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

// IconNames lists every icon the framework ships, sorted.
func IconNames() []string {
	out := make([]string, 0, len(icons))
	for name := range icons {
		out = append(out, name)
	}

	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}

	return out
}

// Equalizer draws a music-style level meter that moves with tick.
func Equalizer(class string, bars, tick int) Component {
	if bars < 1 {
		bars = 12
	}

	levels := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	var sb strings.Builder

	for i := range bars {
		step := (i*3 + tick*2) % (len(levels) * 2)
		if step >= len(levels) {
			step = len(levels)*2 - step - 1
		}
		sb.WriteRune(levels[step])
	}

	return Text(join("nowrap", class), sb.String())
}
