package gtui

import (
	"strconv"
	"strings"

	"github.com/orayew2002/gtui/src/dom"
	"github.com/orayew2002/gtui/src/render"
)

// EditorOptions describes an editable view of some text. The lines belong to
// the app, not to the component: every keystroke reports the new text through
// OnChange, so the app stays the single place the buffer lives and undo,
// saving or sending the text to a language server are all the app's own code.
type EditorOptions struct {
	Class string

	Lines []string // the buffer, one entry per line
	Line  int      // the caret's line, counting from one
	Col   int      // the caret's column in runes, counting from one

	Marks map[int]string // a glyph in the gutter, by line, for diagnostics
	Rows  int            // visible lines; zero means fill the box

	// Tab is how many spaces the tab key inserts, zero meaning four. The editor
	// draws the lines exactly as they are given, and a tab's width depends on
	// where the run holding it starts, so a buffer that uses spaces is the one
	// whose columns always line up. Expand tabs before handing text in.
	Tab      int
	NoGutter bool // leave out the line numbers
	ReadOnly bool

	OnChange func(lines []string, line, col int)
	OnMove   func(line, col int)

	// OnKey sees every key before the editor does. Return true to say the key
	// was handled, which is how an app adds its own bindings, a completion
	// popup for instance, without the editor also acting on them.
	OnKey func(e *Event) bool
}

const editorCaret = '▌'

// CodeEdit is a text editor: it draws the lines with a gutter and a caret, and
// reports what the keys did. It is the editable twin of CodeView.
func CodeEdit(opts EditorOptions) Component {
	if len(opts.Lines) == 0 {
		opts.Lines = []string{""}
	}
	if opts.Tab < 1 {
		opts.Tab = 4
	}

	opts.Line, opts.Col = clampCaret(opts.Lines, opts.Line, opts.Col)

	return Component{
		class: opts.Class,
		build: func(a *App) *Node {
			base, id := splitClass(opts.Class)
			base = join("col focusable clip", base)

			n := dom.NewElement("div")
			n.Set("class", base)
			if id != "" {
				n.Set("id", id)
			}
			n.Set("control", "editor")

			draw := func() {
				n.Clear()
				for _, row := range editorRows(opts) {
					n.Append(row.node(a))
				}
			}
			draw()

			n.On("keydown", func(e *Event) {
				if opts.OnKey != nil && opts.OnKey(e) {
					e.StopPropagation()
					return
				}
				editorKey(a, &opts, e, draw)
			})

			n.On("paste", func(e *Event) {
				if opts.ReadOnly || e.Value == "" {
					return
				}
				e.StopPropagation()

				lines, line, col := pasteInto(opts.Lines, opts.Line, opts.Col, e.Value)
				editorApply(a, &opts, lines, line, col, draw)
			})

			return n
		},
	}
}

// editorRows draws the window of lines the caret is in.
func editorRows(opts EditorOptions) []Component {
	rows := opts.Rows
	if rows < 1 {
		rows = len(opts.Lines)
	}

	first := opts.Line - rows/2
	first = clampInt(first, 1, max(1, len(opts.Lines)-rows+1))

	width := len(strconv.Itoa(len(opts.Lines)))
	out := make([]Component, 0, rows)

	for at := first; at < first+rows && at <= len(opts.Lines); at++ {
		out = append(out, editorRow(opts, at, width))
	}

	return out
}

func editorRow(opts EditorOptions, at, width int) Component {
	cells := make([]Component, 0, 4)

	if !opts.NoGutter {
		ink := "text-muted dim nowrap shrink-0"
		if at == opts.Line {
			ink = "text-accent bold nowrap shrink-0"
		}
		cells = append(cells, Text("w-"+itoa(width)+" text-right "+ink, strconv.Itoa(at)))

		mark := opts.Marks[at]
		if mark == "" {
			mark = " "
		}
		cells = append(cells, Text("w-1 nowrap shrink-0 text-danger", mark))
	}

	// the text goes in a row of its own with no gap: the gap that separates it
	// from the gutter must not also open up between the pieces the caret
	// splits a line into, or the code stops lining up
	cells = append(cells, Div("row gap-0 grow clip", editorText(opts, at)...))

	class := "row gap-1"
	if at == opts.Line {
		class += " bg-raised"
	}

	return Div(class, cells...)
}

// editorText is the line itself, with the caret drawn as a reversed cell so the
// code keeps its columns instead of being pushed sideways by a caret glyph.
func editorText(opts EditorOptions, at int) []Component {
	text := []rune(opts.Lines[at-1])

	if at != opts.Line {
		return []Component{Text("grow nowrap pre", string(text))}
	}

	col := clampInt(opts.Col, 1, len(text)+1)
	before, after := string(text[:col-1]), ""

	under := " "
	if col-1 < len(text) {
		under = string(text[col-1])
		after = string(text[col:])
	}

	return []Component{
		Text("nowrap pre shrink-0", before),
		Text("nowrap pre shrink-0 reverse", under),
		Text("grow nowrap pre", after),
	}
}

// editorKey is the whole keyboard: the keys that move, the keys that edit, and
// everything else typed straight in.
func editorKey(a *App, opts *EditorOptions, e *Event, draw func()) {
	lines, line, col := opts.Lines, opts.Line, opts.Col
	moved, changed := false, false

	switch e.Key {
	case "left":
		line, col, moved = stepLeft(lines, line, col)
	case "right":
		line, col, moved = stepRight(lines, line, col)
	case "up":
		if line > 1 {
			line, moved = line-1, true
		}
	case "down":
		if line < len(lines) {
			line, moved = line+1, true
		}
	case "home":
		col, moved = 1, true
	case "end":
		col, moved = len([]rune(lines[line-1]))+1, true
	case "pageup":
		line, moved = max(1, line-10), true
	case "pagedown":
		line, moved = min(len(lines), line+10), true

	case "enter":
		if !opts.ReadOnly {
			lines, line, col = splitLine(lines, line, col)
			changed = true
		}
	case "backspace":
		if !opts.ReadOnly {
			lines, line, col, changed = deleteBefore(lines, line, col)
		}
	case "delete":
		if !opts.ReadOnly {
			lines, changed = deleteAt(lines, line, col)
		}
	case "tab":
		if !opts.ReadOnly {
			lines, col = insertRunes(lines, line, col, []rune(strings.Repeat(" ", opts.Tab)))
			changed = true
		}

	default:
		if opts.ReadOnly || e.Rune == 0 || e.Rune < ' ' {
			return
		}
		lines, col = insertRunes(lines, line, col, []rune{e.Rune})
		changed = true
	}

	if !moved && !changed {
		return
	}
	e.StopPropagation()

	if changed {
		editorApply(a, opts, lines, line, col, draw)
		return
	}

	opts.Line, opts.Col = clampCaret(lines, line, col)
	draw()

	if opts.OnMove != nil {
		opts.OnMove(opts.Line, opts.Col)
	}
	if a != nil {
		a.Redraw()
	}
}

// editorApply takes the new text, redraws, and tells the app.
func editorApply(a *App, opts *EditorOptions, lines []string, line, col int, draw func()) {
	opts.Lines = lines
	opts.Line, opts.Col = clampCaret(lines, line, col)
	draw()

	if opts.OnChange != nil {
		opts.OnChange(lines, opts.Line, opts.Col)
	}
	if a != nil {
		a.Redraw()
	}
}

func clampCaret(lines []string, line, col int) (int, int) {
	line = clampInt(line, 1, max(1, len(lines)))
	return line, clampInt(col, 1, len([]rune(lines[line-1]))+1)
}

func stepLeft(lines []string, line, col int) (int, int, bool) {
	if col > 1 {
		return line, col - 1, true
	}
	if line > 1 {
		return line - 1, len([]rune(lines[line-2])) + 1, true
	}
	return line, col, false
}

func stepRight(lines []string, line, col int) (int, int, bool) {
	if col <= len([]rune(lines[line-1])) {
		return line, col + 1, true
	}
	if line < len(lines) {
		return line + 1, 1, true
	}
	return line, col, false
}

// insertRunes puts runes in at the caret and reports where the caret lands.
func insertRunes(lines []string, line, col int, add []rune) ([]string, int) {
	out := append([]string(nil), lines...)
	text := []rune(out[line-1])
	col = clampInt(col, 1, len(text)+1)

	grown := make([]rune, 0, len(text)+len(add))
	grown = append(grown, text[:col-1]...)
	grown = append(grown, add...)
	grown = append(grown, text[col-1:]...)

	out[line-1] = string(grown)
	return out, col + len(add)
}

// splitLine is what enter does: the tail of the line becomes the next one.
func splitLine(lines []string, line, col int) ([]string, int, int) {
	text := []rune(lines[line-1])
	col = clampInt(col, 1, len(text)+1)

	out := make([]string, 0, len(lines)+1)
	out = append(out, lines[:line-1]...)
	out = append(out, string(text[:col-1]), string(text[col-1:]))
	out = append(out, lines[line:]...)

	return out, line + 1, 1
}

// deleteBefore is backspace: inside a line it removes a rune, at the start of
// one it joins this line to the one above.
func deleteBefore(lines []string, line, col int) ([]string, int, int, bool) {
	if col > 1 {
		out := append([]string(nil), lines...)
		text := []rune(out[line-1])
		out[line-1] = string(append(text[:col-2:col-2], text[col-1:]...))
		return out, line, col - 1, true
	}

	if line == 1 {
		return lines, line, col, false
	}

	above := []rune(lines[line-2])

	out := make([]string, 0, len(lines)-1)
	out = append(out, lines[:line-2]...)
	out = append(out, string(above)+lines[line-1])
	out = append(out, lines[line:]...)

	return out, line - 1, len(above) + 1, true
}

// deleteAt is the delete key: it takes the rune under the caret, or pulls the
// next line up when the caret sits at the end of this one.
func deleteAt(lines []string, line, col int) ([]string, bool) {
	text := []rune(lines[line-1])

	if col <= len(text) {
		out := append([]string(nil), lines...)
		out[line-1] = string(append(text[:col-1:col-1], text[col:]...))
		return out, true
	}

	if line >= len(lines) {
		return lines, false
	}

	out := make([]string, 0, len(lines)-1)
	out = append(out, lines[:line-1]...)
	out = append(out, lines[line-1]+lines[line])
	out = append(out, lines[line+1:]...)

	return out, true
}

// pasteInto drops a whole clipboard in at the caret, however many lines it has.
func pasteInto(lines []string, line, col int, text string) ([]string, int, int) {
	parts := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")

	if len(parts) == 1 {
		out, at := insertRunes(lines, line, col, []rune(parts[0]))
		return out, line, at
	}

	head := []rune(lines[line-1])
	col = clampInt(col, 1, len(head)+1)

	out := make([]string, 0, len(lines)+len(parts))
	out = append(out, lines[:line-1]...)
	out = append(out, string(head[:col-1])+parts[0])
	out = append(out, parts[1:len(parts)-1]...)
	out = append(out, parts[len(parts)-1]+string(head[col-1:]))
	out = append(out, lines[line:]...)

	return out, line + len(parts) - 1, len([]rune(parts[len(parts)-1])) + 1
}

// EditorWidth is how wide the gutter and the longest line make an editor, which
// an app needs when it puts one in a box of its own.
func EditorWidth(lines []string) int {
	widest := 0
	for _, line := range lines {
		if w := render.StringWidth(line); w > widest {
			widest = w
		}
	}
	return widest + len(strconv.Itoa(len(lines))) + 2
}
