// A code editor you can type in. The buffer belongs to this file, CodeEdit
// draws it and reports what the keys did, and a language server — gopls here —
// can be attached to it for real diagnostics, completion and hover.
package demos

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/orayew2002/gtui"
	"github.com/orayew2002/gtui/lsp"
)

var editor = struct {
	lines []string
	line  int
	col   int

	marks    map[int]string
	problems []lsp.Diagnostic

	// everything the language server side needs
	client  *lsp.Client
	dir     string
	path    string
	status  string
	suggest []lsp.Completion
	picking int
	hover   string

	once sync.Once
}{
	// spaces, not tabs: a tab's width depends on where the run it sits in
	// starts, and the caret splits a line into runs, so a buffer of spaces is
	// the one that always lines up
	lines: strings.Split(`package main

import "fmt"

func main() {
    count := 3
    fmt.Println("hello", count)
}`, "\n"),
	line:   6,
	col:    2,
	marks:  map[int]string{},
	status: "not connected",
}

// Editor is the project: a file tree, the editor, and what the server says.
func Editor(a *gtui.App, _ int) gtui.Component {
	host = a

	return gtui.Column("full gap-0",
		editorBar(),
		gtui.Row("grow clip",
			editorTree(),
			gtui.Column("grow clip gap-0",
				// the editor itself: the buffer goes in, the keys come back out
				gtui.CodeEdit(gtui.EditorOptions{
					Class: "grow clip px-1 #ed",
					Lines: editor.lines,
					Line:  editor.line,
					Col:   editor.col,
					Marks: editor.marks,
					Rows:  editorRows(),
					// the editor redraws its own lines, but the bars around it
					// are built from this state, so the view is rebuilt too or
					// the line count and the caret position go stale
					OnChange: func(lines []string, line, col int) {
						editor.lines, editor.line, editor.col = lines, line, col
						editor.suggest = nil
						editorSend()
						host.Refresh()
					},
					OnMove: func(line, col int) {
						editor.line, editor.col = line, col
						editor.suggest = nil
						host.Refresh()
					},
					OnKey: editorHotkey,
				}),
				editorPanel(),
			),
		),
		editorStatus(),
	)
}

func editorBar() gtui.Component {
	return gtui.Row("px-1 gap-2 items-center shrink-0 border-b border-line",
		gtui.IconText("text-accent bold nowrap", "code", "main.go"),
		gtui.Text("text-muted dim nowrap", "go · utf-8"),
		gtui.Spacer(),
		gtui.Kbd("", "ctrl+space"),
		gtui.Text("text-muted dim nowrap", "complete"),
		gtui.Button("text-accent focus:bold nowrap #ed-lsp", editorLabel(), func(e *gtui.Event) {
			editorConnect()
		}),
	)
}

func editorLabel() string {
	if editor.client != nil {
		return "⚡ gopls"
	}
	return "connect gopls"
}

func editorTree() gtui.Component {
	return gtui.Tree("w-18 shrink-0 px-1 border-r border-line",
		gtui.TreeNode{Label: "probe", Open: true, Children: []gtui.TreeNode{
			{Label: "main.go", Icon: "code", Value: "open"},
			{Label: "go.mod", Icon: "file"},
		}},
	)
}

// editorPanel is whichever of the three things there is to show: the
// suggestions, the problems, or nothing at all.
func editorPanel() gtui.Component {
	switch {
	case len(editor.suggest) > 0:
		return editorSuggestions()
	case len(editor.problems) > 0:
		return editorProblems()
	case editor.hover != "":
		return gtui.Div("shrink-0 clip border-t border-line px-1",
			gtui.Text("text-muted truncate", editor.hover),
		)
	}
	return gtui.Div("")
}

func editorSuggestions() gtui.Component {
	rows := make([]gtui.Component, 0, 6)

	for at, item := range editor.suggest {
		if at >= 6 {
			break
		}

		class := "row gap-1 px-1 nowrap"
		if at == editor.picking {
			class += " bg-raised bold"
		}

		rows = append(rows, gtui.Div(class,
			gtui.Text("w-16 truncate text-ink", item.Label),
			gtui.Text("w-8 truncate text-muted dim", item.Kind),
			gtui.Text("grow truncate text-muted dim", item.Detail),
		))
	}

	return gtui.Column("shrink-0 clip border-t border-line",
		gtui.Text("px-1 text-muted dim nowrap", "↑↓ choose · enter insert · esc close"),
		gtui.Column("gap-0", rows...),
	)
}

func editorProblems() gtui.Component {
	rows := make([]gtui.Component, 0, 4)

	for at, item := range editor.problems {
		if at >= 4 {
			break
		}

		ink := "text-danger"
		if item.Severity != "error" {
			ink = "text-warning"
		}

		rows = append(rows, gtui.Div("row gap-1 px-1 nowrap",
			gtui.Text("w-8 shrink-0 "+ink, "%d:%d", item.Line, item.Col),
			gtui.Text("grow truncate "+ink, item.Message),
		))
	}

	return gtui.Column("shrink-0 clip border-t border-line",
		gtui.Text("px-1 text-muted dim nowrap", "%d problem(s) from gopls", len(editor.problems)),
		gtui.Column("gap-0", rows...),
	)
}

func editorStatus() gtui.Component {
	return gtui.StatusBar("border-t border-line px-1 shrink-0",
		[]gtui.Component{
			gtui.Text("nowrap", "%d:%d", editor.line, editor.col),
			gtui.Text("nowrap text-muted dim", "%d lines", len(editor.lines)),
		},
		[]gtui.Component{
			gtui.Row("gap-1", gtui.StatusDot("", editorHealth()), gtui.Text("nowrap text-muted", editor.status)),
		},
	)
}

func editorHealth() string {
	switch {
	case editor.client == nil:
		return "idle"
	case len(editor.problems) > 0:
		return "down"
	}
	return "ok"
}

// editorRows leaves room for the bars and whichever panel is open.
func editorRows() int {
	_, height := host.Size()

	rows := height - 4
	if len(editor.suggest) > 0 || len(editor.problems) > 0 {
		rows -= 6
	}

	return max(rows, 3)
}

// editorHotkey is the editor's own bindings: completion, and walking the list
// when one is open. Returning true keeps the key away from the editing keys.
func editorHotkey(e *gtui.Event) bool {
	if e.Key == "ctrl+space" {
		editorComplete()
		return true
	}

	if len(editor.suggest) == 0 {
		return false
	}

	switch e.Key {
	case "down":
		editor.picking = (editor.picking + 1) % len(editor.suggest)
	case "up":
		editor.picking = (editor.picking + len(editor.suggest) - 1) % len(editor.suggest)
	case "enter":
		editorInsert(editor.suggest[editor.picking])
	case "esc":
		editor.suggest = nil
	default:
		return false
	}

	host.Refresh()
	return true
}

// editorInsert puts a suggestion in, replacing the word being typed.
func editorInsert(item lsp.Completion) {
	line := []rune(editor.lines[editor.line-1])
	start := editor.col - 1

	for start > 0 && isWord(line[start-1]) {
		start--
	}

	grown := make([]rune, 0, len(line)+len(item.Insert))
	grown = append(grown, line[:start]...)
	grown = append(grown, []rune(item.Insert)...)
	grown = append(grown, line[editor.col-1:]...)

	editor.lines = append([]string(nil), editor.lines...)
	editor.lines[editor.line-1] = string(grown)
	editor.col = start + len([]rune(item.Insert)) + 1
	editor.suggest = nil

	editorSend()
}

func isWord(r rune) bool {
	return r == '_' || r == '.' ||
		(r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

// editorConnect starts gopls on a scratch module holding the buffer, which is
// all a language server needs to answer about it.
func editorConnect() {
	if editor.client != nil {
		editor.status = "already connected"
		host.Refresh()
		return
	}

	server, err := exec.LookPath("gopls")
	if err != nil {
		editor.status = "gopls is not installed"
		host.Refresh()
		return
	}

	dir, err := os.MkdirTemp("", "gtui-editor")
	if err != nil {
		editor.status = "no temp directory"
		host.Refresh()
		return
	}

	editor.dir = dir
	editor.path = filepath.Join(dir, "main.go")

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module probe\n\ngo 1.21\n"), 0o644); err != nil {
		editor.status = "cannot write the module"
		host.Refresh()
		return
	}
	if err := os.WriteFile(editor.path, []byte(editorText()), 0o644); err != nil {
		editor.status = "cannot write the file"
		host.Refresh()
		return
	}

	client, err := lsp.Start(dir, server)
	if err != nil {
		editor.status = "gopls did not start"
		host.Refresh()
		return
	}

	// the server answers on its own goroutine, so hand the work back to the
	// app's loop before touching anything the view reads
	client.OnDiagnostics(func(path string, items []lsp.Diagnostic) {
		host.Do(func() {
			editor.problems = items
			editor.marks = map[int]string{}

			for _, item := range items {
				editor.marks[item.Line] = "●"
			}

			editor.status = "gopls is watching"
			host.Refresh()
		})
	})

	editor.client = client
	editor.status = "gopls connected"

	_ = client.Open(editor.path, "go", editorText())
	host.Refresh()
}

// editorSend keeps the server's copy of the file up to date.
func editorSend() {
	if editor.client == nil {
		return
	}
	_ = editor.client.Change(editor.path, editorText())
}

func editorComplete() {
	if editor.client == nil {
		editor.status = "connect gopls first"
		host.Refresh()
		return
	}

	items, err := editor.client.Complete(editor.path, editor.line, editor.col, editorText())
	if err != nil {
		editor.status = "completion: " + err.Error()
		host.Refresh()
		return
	}

	editor.suggest, editor.picking = items, 0
	editor.status = "gopls offered " + itoa(len(items))
	host.Refresh()
}

func editorText() string { return strings.Join(editor.lines, "\n") }

func itoa(v int) string {
	if v == 0 {
		return "0"
	}

	digits := ""
	for v > 0 {
		digits = string(rune('0'+v%10)) + digits
		v /= 10
	}
	return digits
}
