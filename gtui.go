package gtui

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/orayew2002/gtui/src/colors"
	"github.com/orayew2002/gtui/src/dom"
	"github.com/orayew2002/gtui/src/html"
	"github.com/orayew2002/gtui/src/layout"
	"github.com/orayew2002/gtui/src/paint"
	"github.com/orayew2002/gtui/src/render"
	"github.com/orayew2002/gtui/src/style"
	"github.com/orayew2002/gtui/src/term"
)

const scrollStep = 3

type Event = dom.Event

type Node = dom.Node

type App struct {
	root *dom.Node
	base style.Style

	terminal *term.Terminal
	screen   *render.Screen
	tree     *layout.Box
	diags    []style.Diagnostic

	focus *dom.Node
	hover *dom.Node

	quitKeys map[string]bool
	binds    map[string]func()
	tasks    chan func()
	redrawCh chan struct{}
	stopCh   chan struct{}

	onResize    func(int, int)
	mouse       bool
	cursor      bool
	modalOpen   bool
	beforeModal *dom.Node
	view        func() Component

	stopOnce sync.Once
	mu       sync.Mutex
	err      error
}

func Init() *App {
	colors.Use(colors.Dark)

	base := style.New()
	base.Foreground = style.P(colors.Active().Ink)
	base.Background = style.P(colors.Active().Canvas)

	return &App{
		root:     dom.NewElement(dom.RootTag),
		base:     base,
		mouse:    true,
		quitKeys: map[string]bool{"q": true, "ctrl+c": true},
		binds:    map[string]func(){},
		tasks:    make(chan func(), 64),
		redrawCh: make(chan struct{}, 1),
		stopCh:   make(chan struct{}),
	}
}

func New() *App { return Init() }

func (a *App) HTML(components ...Component) *App {
	a.root.Clear()
	for _, n := range a.nodes(components) {
		a.root.Append(n)
	}

	a.focus = nil
	a.hover = nil
	return a.Redraw()
}

func (a *App) Markup(src string) *App {
	doc := html.Parse(src)

	a.root.Clear()
	for _, child := range append([]*dom.Node(nil), doc.Root.Children...) {
		a.root.Append(child)
	}

	a.focus = nil
	a.hover = nil
	return a.Redraw()
}

func (a *App) Set(sel string, components ...Component) *App {
	n := a.Query(sel)
	if n == nil {
		return a
	}

	n.Clear()
	for _, child := range a.nodes(components) {
		n.Append(child)
	}

	a.dropDetachedFocus()
	return a.Redraw()
}

func (a *App) Add(sel string, components ...Component) *App {
	n := a.Query(sel)
	if n == nil {
		return a
	}

	for _, child := range a.nodes(components) {
		n.Append(child)
	}
	return a.Redraw()
}

func (a *App) Open(target string, page ...Component) *App {
	if target == "" {
		return a
	}
	return a.Set(target, page...)
}

func (a *App) Select(sel string) *App {
	node := a.Query(sel)
	if node == nil {
		return a
	}
	return a.mark(node)
}

func (a *App) mark(node *dom.Node) *App {
	for _, link := range a.QueryAll("div") {
		if link.Has("href") {
			link.State.Current = false
		}
	}

	node.State.Current = true
	return a.Redraw()
}

func (a *App) dropDetachedFocus() {
	if a.focus == nil {
		return
	}

	for n := a.focus; n != nil; n = n.Parent {
		if n == a.root {
			return
		}
	}

	a.focus.State.Focus = false
	a.focus = nil
}

func (a *App) Render() error {
	options := term.DefaultOptions()
	options.Mouse = a.mouse
	options.HideCursor = !a.cursor

	t, err := term.Open(options)
	if err != nil {
		return fmt.Errorf("gtui: open terminal: %w", err)
	}
	defer t.Close()

	a.terminal = t
	w, h := t.Size()
	a.screen = render.NewScreen(t.Out, w, h)

	a.frame()

	for {
		select {
		case <-a.stopCh:
			return a.err

		case fn := <-a.tasks:
			fn()
			a.drainTasks()
			a.frame()

		case <-a.redrawCh:
			a.frame()

		case <-t.Resized():
			w, h := t.Size()
			a.screen.Resize(w, h)

			if a.onResize != nil {
				a.onResize(w, h)
			}
			a.frame()

		case ev, ok := <-t.Events():
			if !ok {
				return a.err
			}
			a.handle(ev)
			a.drainEvents(t)
			a.frame()
		}
	}
}

func (a *App) RenderTo(w io.Writer, width, height int) error {
	a.screen = render.NewScreen(w, width, height)
	a.frame()
	return a.err
}

func (a *App) Root() *dom.Node { return a.root }

func (a *App) Query(sel string) *dom.Node { return dom.Query(a.root, sel) }

func (a *App) QueryAll(sel string) []*dom.Node { return dom.QueryAll(a.root, sel) }

func (a *App) Text(sel, text string) *App {
	if n := a.Query(sel); n != nil {
		n.SetText(text)
		a.Redraw()
	}
	return a
}

func (a *App) Class(sel, class string, on bool) *App {
	for _, n := range a.QueryAll(sel) {
		if on {
			n.AddClass(class)
		} else {
			n.RemoveClass(class)
		}
	}
	return a.Redraw()
}

func (a *App) On(event string, fn func(*dom.Event)) *App {
	a.root.On(event, fn)
	return a
}

func (a *App) OnClick(sel string, fn func(*dom.Event)) *App {
	for _, n := range a.QueryAll(sel) {
		n.On("click", fn)
	}
	return a
}

func (a *App) OnKey(fn func(*dom.Event)) *App {
	return a.On("keydown", fn)
}

func (a *App) Theme(theme colors.Theme) *App {
	colors.Use(theme)

	a.base.Foreground = style.P(theme.Ink)
	a.base.Background = style.P(theme.Canvas)

	return a.Redraw()
}

func (a *App) Colors(fg, bg colors.Color) *App {
	a.base.Foreground = style.P(fg)
	a.base.Background = style.P(bg)
	return a.Redraw()
}

func (a *App) Tokens(edit func(*colors.Theme)) *App {
	if edit == nil {
		return a
	}

	theme := colors.Active()
	edit(&theme)

	colors.Register(theme)
	return a.Theme(theme)
}

func (a *App) ThemeByName(name string) *App {
	if theme, ok := colors.ByName(name); ok {
		return a.Theme(theme)
	}
	return a
}

func (a *App) Themes() []string { return colors.Names() }

func (a *App) ToggleTheme() *App {
	if colors.Active().Name == "light" {
		return a.Theme(colors.Dark)
	}
	return a.Theme(colors.Light)
}

func (a *App) ThemeName() string { return colors.Active().Name }

func (a *App) UseMouse(on bool) *App {
	a.mouse = on

	if a.terminal != nil {
		a.terminal.Mouse(on)
	}
	return a
}

func (a *App) UseCursor(on bool) *App {
	a.cursor = on
	return a.Redraw()
}

func (a *App) Size() (int, int) {
	if a.screen != nil {
		return a.screen.Size()
	}
	if a.terminal != nil {
		return a.terminal.Size()
	}
	return 80, 24
}

func (a *App) OnResize(fn func(width, height int)) *App {
	a.onResize = fn
	return a
}

func (a *App) Copy(text string) *App {
	if a.screen == nil || text == "" {
		return a
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	fmt.Fprintf(a.screen, "\x1b]52;c;%s\x07", encoded)

	return a
}

func (a *App) Bind(key string, fn func()) *App {
	if fn == nil {
		delete(a.binds, strings.ToLower(key))
		return a
	}

	a.binds[strings.ToLower(key)] = fn
	return a
}

func (a *App) QuitOn(keys ...string) *App {
	a.quitKeys = make(map[string]bool, len(keys))
	for _, k := range keys {
		a.quitKeys[strings.ToLower(k)] = true
	}
	return a
}

func (a *App) Redraw() *App {
	select {
	case a.redrawCh <- struct{}{}:
	default:
	}
	return a
}

func (a *App) Do(fn func()) {
	if fn == nil {
		return
	}
	select {
	case a.tasks <- fn:
	case <-a.stopCh:
	}
}

func (a *App) Stop() {
	a.stopOnce.Do(func() { close(a.stopCh) })
}

func (a *App) Focus(n *dom.Node) {
	if a.focus == n {
		return
	}

	if a.focus != nil {
		a.focus.State.Focus = false
		a.focus.Dispatch(&dom.Event{Type: "blur", Target: a.focus})
	}

	a.focus = n
	if n != nil {
		n.State.Focus = true
		n.Dispatch(&dom.Event{Type: "focus", Target: n})
	}
	a.Redraw()
}

func (a *App) Focused() *dom.Node { return a.focus }

func (a *App) Diagnostics() []style.Diagnostic {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]style.Diagnostic(nil), a.diags...)
}

func (a *App) Snapshot(width, height int) string {
	tree := a.build(width, height)
	if tree == nil {
		return ""
	}

	buf := render.NewBuffer(width, height)
	buf.Reset(render.Cell{Ch: ' '})
	paint.Paint(buf, tree)

	var sb strings.Builder
	for y := range height {
		row := make([]rune, 0, width)

		for x := range width {
			cell := buf.At(x, y)
			switch {
			case cell == nil, cell.Cont:
				continue
			case cell.Ch == 0:
				row = append(row, ' ')
			default:
				row = append(row, cell.Ch)
			}
		}

		sb.WriteString(strings.TrimRight(string(row), " "))
		if y < height-1 {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}

func (a *App) build(width, height int) *layout.Box {
	a.mu.Lock()
	a.diags = a.diags[:0]
	tree := layout.Build(a.root, a.base, &a.diags)
	a.mu.Unlock()

	if tree == nil {
		return nil
	}

	if tree.Style.Width.IsAuto() {
		tree.Style.Width = style.Full()
	}
	if tree.Style.Height.IsAuto() {
		tree.Style.Height = style.Full()
	}

	layout.Layout(tree, render.Rect{W: width, H: height})
	a.tree = tree

	return tree
}

func (a *App) frame() {
	if a.screen == nil {
		return
	}

	w, h := a.screen.Size()
	tree := a.build(w, h)
	if tree == nil {
		return
	}

	blank := render.Cell{Ch: ' '}
	if a.base.Background.Set {
		blank.Bg, blank.HasBg = a.base.Background.Color, true
	}
	a.screen.Clear(blank)

	if a.focusVisible() {
		tree = a.build(w, h)
		if tree == nil {
			return
		}
		a.screen.Clear(blank)
	}

	paint.Paint(a.screen.Back(), tree)

	if a.cursor {
		a.placeCursor()
	}

	if err := a.screen.Flush(); err != nil && a.err == nil {
		a.err = err
		a.Stop()
	}
}

func (a *App) placeCursor() {
	box := a.findBox(a.focus)
	if box == nil || a.focus == nil || !a.focus.Has("type") {
		a.screen.HideCursor()
		return
	}

	at := caret(a.focus)
	a.screen.ShowCursorAt(box.Content.X+at, box.Content.Y)
}

func (a *App) drainTasks() {
	for {
		select {
		case fn := <-a.tasks:
			fn()
		default:
			return
		}
	}
}

func (a *App) drainEvents(t *term.Terminal) {
	for {
		select {
		case ev, ok := <-t.Events():
			if !ok {
				return
			}
			a.handle(ev)
		default:
			return
		}
	}
}
