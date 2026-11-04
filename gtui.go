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

// paintedCaret is the glyph paintField draws at the caret.
const paintedCaret = '▌'

type Event = dom.Event

type Node = dom.Node

type App struct {
	root *dom.Node
	base style.Style

	terminal *term.Terminal
	screen   *render.Screen
	tree     *layout.Box
	stale    bool
	diags    []style.Diagnostic

	focus *dom.Node
	hover *dom.Node

	quitKeys map[string]bool
	binds    map[string]func()
	tasks    chan func()
	redrawCh chan struct{}
	stopCh   chan struct{}

	onResize    func(int, int)
	viewW       int
	viewH       int
	mouse       bool
	cursor      bool
	modalOpen   bool
	beforeModal *dom.Node
	toastSeq    int
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
	a.resetOverlays()
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
	a.resetOverlays()
	return a.Redraw()
}

func (a *App) Set(sel string, components ...Component) *App {
	n := a.Query(sel)
	if n == nil {
		return a
	}

	keep := a.focusID()

	n.Clear()
	for _, child := range a.nodes(components) {
		n.Append(child)
	}

	a.restoreFocus(keep)
	return a.Redraw()
}

func (a *App) Add(sel string, components ...Component) *App {
	n := a.Query(sel)
	if n == nil {
		return a
	}

	keep := a.focusID()

	for _, child := range a.nodes(components) {
		n.Append(child)
	}

	a.restoreFocus(keep)
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

// mark clears the current flag on the links that share this one's target, so
// every nav group (sidebar, tabs) keeps its own active item.
func (a *App) mark(node *dom.Node) *App {
	target := node.Get("href")

	for _, link := range a.QueryAll("div") {
		if link.Has("href") && link.Get("href") == target {
			link.State.Current = false
		}
	}

	node.State.Current = true
	return a.Redraw()
}

func (a *App) attached(n *dom.Node) bool { return descends(n, a.root) }

func (a *App) byID(id string) *dom.Node {
	if id == "" {
		return nil
	}
	return a.Query("#" + id)
}

func (a *App) focusID() string {
	if a.focus == nil {
		return ""
	}
	return a.focus.ID
}

// restoreFocus hands focus to the node that replaced the focused one. Rebuilds
// make fresh nodes, so the id is the only thing that survives; focus is dropped
// only when nothing carries it any more.
func (a *App) restoreFocus(id string) {
	if a.focus == nil || a.attached(a.focus) {
		return
	}

	if next := a.byID(id); next != nil {
		a.Focus(next)
		return
	}

	a.focus.State.Focus = false
	a.focus = nil
}

func (a *App) resetOverlays() {
	a.modalOpen = false
	a.beforeModal = nil
}

func (a *App) Render() error {
	defer a.Stop()

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
	a.resize(w, h)

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
			a.drainTasks()
			a.frame()

		case <-t.Resized():
			w, h := t.Size()
			a.screen.Resize(w, h)
			a.resize(w, h)

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
	a.drainTasks()

	a.screen = render.NewScreen(w, width, height)
	a.resize(width, height)
	a.frame()

	return a.err
}

// resize tells the document how much room it has, then rebuilds a registered
// view so pages that lay themselves out from app.Size() see the new size.
// Refresh carries the state of every node with an id across that rebuild.
func (a *App) resize(width, height int) {
	if a.viewW == width && a.viewH == height {
		return
	}

	a.viewW, a.viewH = width, height

	if a.view != nil {
		a.Refresh()
	}
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
	return a.onEach(sel, "click", fn)
}

func (a *App) OnChange(sel string, fn func(*dom.Event)) *App {
	return a.onEach(sel, "change", fn)
}

func (a *App) OnInput(sel string, fn func(*dom.Event)) *App {
	return a.onEach(sel, "input", fn)
}

func (a *App) onEach(sel, event string, fn func(*dom.Event)) *App {
	for _, n := range a.QueryAll(sel) {
		n.On(event, fn)
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
	if a.viewW > 0 && a.viewH > 0 {
		return a.viewW, a.viewH
	}
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

func (a *App) Press(key string) *App {
	a.handleKey(pressKey(key))
	return a
}

func (a *App) Type(text string) *App {
	for _, r := range text {
		a.Press(string(r))
	}
	return a
}

// Click presses and releases the mouse on one cell, the way Press sends one
// key: the cell is hit tested, focus moves to the control under it and the
// click is dispatched from there. It is how a test drives the pointer.
func (a *App) Click(x, y int) *App {
	a.handleMouse(term.Mouse{Kind: term.MousePress, X: x, Y: y})
	a.handleMouse(term.Mouse{Kind: term.MouseRelease, X: x, Y: y})
	return a
}

// Hover moves the pointer without pressing, so hover: classes light up.
func (a *App) Hover(x, y int) *App {
	a.handleMouse(term.Mouse{Kind: term.MouseMove, X: x, Y: y})
	return a
}

// Wheel scrolls at a cell. A negative count scrolls up, a positive one down.
func (a *App) Wheel(x, y, lines int) *App {
	kind := term.MouseWheelDown
	if lines < 0 {
		kind, lines = term.MouseWheelUp, -lines
	}

	for range lines {
		a.handleMouse(term.Mouse{Kind: kind, X: x, Y: y})
	}
	return a
}

// pressKey names control characters the way the terminal decoder does, so
// Type("a\nb") inserts a newline instead of typing an unknown key.
func pressKey(key string) term.Key {
	switch key {
	case "\r", "\n":
		return term.Key{Name: "enter", Rune: '\r'}
	case "\t":
		return term.Key{Name: "tab", Rune: '\t'}
	case "\x7f":
		return term.Key{Name: "backspace"}
	case "\x1b":
		return term.Key{Name: "esc"}
	}

	rune_ := rune(0)
	if runes := []rune(key); len(runes) == 1 {
		rune_ = runes[0]
	}

	return term.Key{Name: strings.ToLower(key), Rune: rune_}
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
	// whatever asked for a redraw has changed the document, so the boxes the
	// last layout produced no longer describe it
	a.stale = true

	select {
	case a.redrawCh <- struct{}{}:
	default:
	}
	return a
}

// boxes is the box tree, laid out again first if the document moved on since
// the last frame. Input arrives in bursts and the whole burst is handled before
// the next frame, so without this a second click in one burst would be matched
// against boxes the first click already threw away.
func (a *App) boxes() *layout.Box {
	if !a.stale {
		return a.tree
	}

	width, height := a.viewW, a.viewH
	if a.screen != nil {
		width, height = a.screen.Size()
	}
	if width <= 0 || height <= 0 {
		return a.tree
	}

	return a.layout(width, height)
}

// Do queues fn for the next Render loop pass, Snapshot, RenderTo or Flush. It
// never blocks: with nothing draining the queue the stalest task is dropped to
// make room, so a headless or stopped app cannot wedge its callers.
func (a *App) Do(fn func()) {
	if fn == nil {
		return
	}

	select {
	case <-a.stopCh:
		return
	default:
	}

	for range cap(a.tasks) + 1 {
		select {
		case a.tasks <- fn:
			return
		default:
		}

		select {
		case <-a.tasks:
		default:
		}
	}
}

// Flush runs every task queued with Do, After, Every and Animate on the calling
// goroutine. Render, Snapshot and RenderTo do this for you.
func (a *App) Flush() *App {
	a.drainTasks()
	return a
}

func (a *App) Stop() {
	a.stopOnce.Do(func() { close(a.stopCh) })
}

func (a *App) Focus(n *dom.Node) {
	if a.focus == n {
		return
	}

	prev := a.focus
	if prev != nil {
		prev.State.Focus = false
	}

	a.focus = n
	if n != nil {
		n.State.Focus = true
	}

	// blur reaches the old node once focus has already moved, so a handler can
	// ask where it went: a dropdown keeps its list open while focus is still
	// inside it and closes when focus leaves for good
	if prev != nil {
		prev.Dispatch(&dom.Event{Type: "blur", Target: prev})
	}
	if n != nil {
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
	a.drainTasks()
	a.resize(width, height)

	tree := a.layout(width, height)
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
	a.viewW, a.viewH = width, height

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
	a.stale = false

	return tree
}

// layout lays the document out and, like the live loop, scrolls the focused
// node into view, so a headless Snapshot shows what the screen shows.
func (a *App) layout(width, height int) *layout.Box {
	tree := a.build(width, height)
	if tree == nil {
		return nil
	}

	// a field cannot window its value around the caret until it knows how wide
	// it ended up, and only layout knows that
	if a.remeasureFields() {
		tree = a.build(width, height)
	}

	if a.focusVisible() {
		tree = a.build(width, height)
	}

	return tree
}

func (a *App) frame() {
	if a.screen == nil {
		return
	}

	w, h := a.screen.Size()
	tree := a.layout(w, h)
	if tree == nil {
		return
	}

	blank := render.Cell{Ch: ' '}
	if a.base.Background.Set {
		blank.Bg, blank.HasBg = a.base.Background.Color, true
	}
	a.screen.Clear(blank)

	paint.Paint(a.screen.Back(), tree)

	if err := a.screen.Flush(); err != nil {
		a.fail(err)
		return
	}

	if !a.cursor {
		return
	}

	// the caret escape has to reach the terminal after the cells it sits on,
	// or the hardware cursor ends up wherever the diff stopped writing
	a.placeCursor()
	if err := a.screen.Flush(); err != nil {
		a.fail(err)
	}
}

func (a *App) fail(err error) {
	if err == nil || a.err != nil {
		return
	}

	a.err = err
	a.Stop()
}

func (a *App) placeCursor() {
	box := a.findBox(a.focus)
	if box == nil || a.focus == nil || !a.focus.Has("type") {
		a.screen.HideCursor()
		return
	}

	line, col, ok := caretCell(a.focus)
	if !ok {
		a.screen.HideCursor()
		return
	}

	// the terminal caret belongs on a cell the field really painted, and a
	// clipping ancestor can cut that cell away even when it sits inside the
	// field's own box
	x, y := box.Content.X+col, box.Content.Y+line
	if !box.Content.Intersect(box.Clip).Contains(x, y) {
		a.screen.HideCursor()
		return
	}

	a.screen.ShowCursorAt(x, y)
}

// caretCell finds the caret glyph the field painted and reports the row and
// column it sits on, which is the one place that already accounts for text area
// lines, wide runes, password bullets and the window a long value scrolls in.
func caretCell(n *dom.Node) (int, int, bool) {
	if len(n.Children) == 0 {
		return 0, 0, false
	}

	text := n.Children[0]

	shown := text.Text()
	at := strings.IndexRune(shown, paintedCaret)
	if at < 0 {
		return 0, 0, false
	}

	before := shown[:at]
	line := strings.Count(before, "\n") - text.ScrollY

	return line, render.StringWidth(before[strings.LastIndexByte(before, '\n')+1:]), true
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
