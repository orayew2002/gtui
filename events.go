package gtui

import (
	"github.com/orayew2002/gtui/src/dom"
	"github.com/orayew2002/gtui/src/layout"
	"github.com/orayew2002/gtui/src/style"
	"github.com/orayew2002/gtui/src/term"
)

func (a *App) handle(ev term.Event) {
	switch e := ev.(type) {
	case term.Key:
		a.handleKey(e)
	case term.Mouse:
		a.handleMouse(e)
	case term.Paste:
		a.handlePaste(e)
	}
}

func (a *App) handleKey(k term.Key) {
	// a modified key never types a character, so alt+s stays a shortcut
	// instead of landing an "s" in whatever field holds focus
	typed := k.Rune
	if k.Alt || k.Ctrl {
		typed = 0
	}

	layer := a.modalLayer()

	target := a.focus
	if layer != nil && !descends(target, layer) {
		target = layer
	}
	if target == nil {
		target = a.root
	}

	event := &dom.Event{Type: "keydown", Target: target, Key: k.Name, Rune: typed}
	target.Dispatch(event)
	if event.Stopped() {
		return
	}

	if layer != nil {
		a.modalKey(k, typed, layer)
		return
	}

	if fn := a.binds[k.Name]; fn != nil {
		fn()
		a.Redraw()
		return
	}

	if a.jumpTo(k.Name) {
		return
	}

	switch k.Name {
	case "tab":
		a.cycleFocus(1)
		return
	case "shift+tab":
		a.cycleFocus(-1)
		return
	case "enter":
		if a.activateFocus() {
			return
		}
	}

	if typed == ' ' && a.activateFocus() {
		return
	}

	if a.scrollSideways(k.Name) {
		return
	}

	if a.moveFocus(k.Name) {
		return
	}

	if a.scrollFocused(k.Name) {
		return
	}

	if a.quitKeys[k.Name] {
		a.Stop()
	}
}

// modalLayer returns the layer of the modal that owns the keyboard. A rebuild
// can take the layer away without CloseModal running, so a missing or hidden
// layer also clears the flag instead of trapping every later key.
func (a *App) modalLayer() *dom.Node {
	if !a.modalOpen {
		return nil
	}

	layer := a.Query("#" + modalID)
	if layer == nil || layer.HasClass("none") {
		a.resetOverlays()
		return nil
	}

	return layer
}

func (a *App) modalKey(k term.Key, typed rune, layer *dom.Node) {
	switch {
	case k.Name == "esc":
		a.CloseModal()
	case k.Name == "tab", k.Name == "down", k.Name == "right":
		a.focusWithin(a.modalFocusables(), 1)
	case k.Name == "shift+tab", k.Name == "up", k.Name == "left":
		a.focusWithin(a.modalFocusables(), -1)
	case k.Name == "enter", typed == ' ':
		if descends(a.focus, layer) {
			a.activateFocus()
		}
	case a.quitKeys[k.Name]:
		a.Stop()
	}
}

func descends(n, root *dom.Node) bool {
	if root == nil {
		return false
	}

	for ; n != nil; n = n.Parent {
		if n == root {
			return true
		}
	}
	return false
}

func (a *App) jumpTo(key string) bool {
	tree := a.boxes()
	if key == "" || tree == nil {
		return false
	}

	var target *layout.Box
	tree.Walk(func(b *layout.Box) bool {
		if target != nil {
			return false
		}
		if b.Node != nil && b.Style.Shortcut == key {
			target = b
			return false
		}
		return true
	})

	if target == nil {
		return false
	}

	if target.Node.Has("href") {
		a.Focus(target.Node)
		a.activateFocus()
		return true
	}

	if target.Style.Overflow == style.OverflowScroll {
		a.Focus(target.Node)
		return true
	}

	if inner := firstFocusableIn(target); inner != nil {
		a.Focus(inner)
		return true
	}

	a.Focus(target.Node)
	return true
}

func firstFocusableIn(section *layout.Box) *dom.Node {
	var found *dom.Node

	section.Walk(func(b *layout.Box) bool {
		if found != nil {
			return false
		}
		if b == section || b.Node == nil {
			return true
		}
		if b.Style.Focusable && b.Style.Overflow != style.OverflowScroll {
			found = b.Node
			return false
		}
		return true
	})

	return found
}

func (a *App) handlePaste(p term.Paste) {
	target := a.focus
	if target == nil {
		target = a.root
	}

	event := &dom.Event{Type: "paste", Target: target, Value: p.Text}
	target.Dispatch(event)

	if event.Stopped() || a.focus == nil || !a.focus.Has("type") {
		a.Redraw()
		return
	}

	field := a.focus
	kind := field.Get("type")
	runes := []rune(p.Text)
	edited := false

	for i := 0; i < len(runes); i++ {
		key := &dom.Event{Rune: runes[i]}

		// terminals end a pasted line with LF, CR or CRLF, and a field only
		// takes a line break from the enter key
		if runes[i] == '\n' || runes[i] == '\r' {
			if runes[i] == '\r' && i+1 < len(runes) && runes[i+1] == '\n' {
				i++
			}
			if kind != InputArea {
				continue
			}
			key = &dom.Event{Key: "enter"}
		}

		if _, changed := editField(field, key, kind); changed {
			edited = true
		}
	}

	paintField(field)

	if edited {
		field.Dispatch(&dom.Event{Type: "input", Target: field, Value: field.Get("value")})
	}

	a.Redraw()
}

func (a *App) activateFocus() bool {
	if a.focus == nil {
		return false
	}

	node := a.focus
	node.State.Active = true
	node.Dispatch(&dom.Event{Type: "click", Target: node})
	node.State.Active = false

	a.Redraw()
	return true
}

func (a *App) moveFocus(key string) bool {
	delta := 0
	switch key {
	case "down", "right":
		delta = 1
	case "up", "left":
		delta = -1
	default:
		return false
	}

	targets := a.scope()
	if len(targets) == 0 {
		return false
	}

	if a.focus == nil {
		if delta > 0 {
			a.Focus(targets[0])
		} else {
			a.Focus(targets[len(targets)-1])
		}
		return true
	}

	if box := a.findBox(a.focus); box != nil && box.Style.Overflow == style.OverflowScroll {
		return false
	}

	if len(targets) < 2 {
		return false
	}

	a.focusWithin(targets, delta)
	return true
}

func (a *App) scrollSideways(key string) bool {
	box := a.scrollableXAncestor(a.findBox(a.focus))
	if box == nil || box.ScrollMaxX == 0 {
		return false
	}

	node := box.Node

	switch key {
	case "left":
		node.ScrollX -= 4
	case "right":
		node.ScrollX += 4
	case "home":
		node.ScrollX = 0
	case "end":
		node.ScrollX = box.ScrollMaxX
	default:
		return false
	}

	if node.ScrollX < 0 {
		node.ScrollX = 0
	}
	if node.ScrollX > box.ScrollMaxX {
		node.ScrollX = box.ScrollMaxX
	}

	a.Redraw()
	return true
}

func (a *App) scrollableXAncestor(b *layout.Box) *layout.Box {
	for ; b != nil; b = b.Parent {
		if b.Node != nil && b.Style.ScrollX && b.ScrollMaxX > 0 {
			return b
		}
	}
	return nil
}

func (a *App) scrollFocused(key string) bool {

	box := a.scrollableAncestor(a.findBox(a.focus))
	if box == nil || box.ScrollMax == 0 {
		return false
	}

	node := box.Node
	switch key {
	case "up":
		node.ScrollY -= 1
	case "down":
		node.ScrollY += 1
	case "pageup":
		node.ScrollY -= box.Content.H
	case "pagedown":
		node.ScrollY += box.Content.H
	case "home":
		node.ScrollY = 0
	case "end":
		node.ScrollY = box.ScrollMax
	default:
		return false
	}

	a.clampScroll(box)
	a.Redraw()
	return true
}

func (a *App) clampScroll(box *layout.Box) {
	if box.Node.ScrollY < 0 {
		box.Node.ScrollY = 0
	}
	if box.Node.ScrollY > box.ScrollMax {
		box.Node.ScrollY = box.ScrollMax
	}
}

func (a *App) handleMouse(m term.Mouse) {
	hit := a.hitBox(m.X, m.Y)

	// a dialog covers the page, so the pointer stops at it the way the
	// keyboard does: a click beside an open dialog must not reach the button
	// it hides, and the wheel must not scroll the page behind it
	if layer := a.modalLayer(); layer != nil && !descends(nodeOf(hit), layer) {
		a.setHover(nil)
		return
	}

	switch m.Kind {
	case term.MouseMove:
		a.setHover(nodeOf(hit))

	case term.MousePress:
		node := nodeOf(hit)
		a.setHover(node)

		if focusable := a.focusableAncestor(hit); focusable != nil {
			a.Focus(focusable.Node)
		}

		if node != nil {
			node.State.Active = true
			node.Dispatch(&dom.Event{Type: "click", Target: node, X: m.X, Y: m.Y})
			node.State.Active = false
			a.Redraw()
		}

	case term.MouseRelease:
		if node := nodeOf(hit); node != nil {
			node.Dispatch(&dom.Event{Type: "mouseup", Target: node, X: m.X, Y: m.Y})
		}

	case term.MouseWheelUp, term.MouseWheelDown:
		box := a.scrollableAncestor(hit)
		if box == nil {
			return
		}
		if m.Kind == term.MouseWheelUp {
			box.Node.ScrollY -= scrollStep
		} else {
			box.Node.ScrollY += scrollStep
		}
		a.clampScroll(box)
		a.Redraw()

	case term.MouseWheelLeft, term.MouseWheelRight:
		box := a.scrollableXAncestor(hit)
		if box == nil {
			return
		}
		if m.Kind == term.MouseWheelLeft {
			box.Node.ScrollX -= scrollStep
		} else {
			box.Node.ScrollX += scrollStep
		}
		if box.Node.ScrollX < 0 {
			box.Node.ScrollX = 0
		}
		if box.Node.ScrollX > box.ScrollMaxX {
			box.Node.ScrollX = box.ScrollMaxX
		}
		a.Redraw()
	}
}

func (a *App) setHover(node *dom.Node) {
	if a.hover == node {
		return
	}

	if a.hover != nil {
		a.hover.State.Hover = false
	}
	a.hover = node
	if node != nil {
		node.State.Hover = true
	}
	a.Redraw()
}

func nodeOf(b *layout.Box) *dom.Node {
	if b == nil {
		return nil
	}
	return b.Node
}

func (a *App) hitBox(x, y int) *layout.Box {
	tree := a.boxes()
	if tree == nil {
		return nil
	}
	return tree.HitTest(x, y)
}

func (a *App) findBox(node *dom.Node) *layout.Box {
	tree := a.boxes()
	if node == nil || tree == nil {
		return nil
	}

	var found *layout.Box
	tree.Walk(func(b *layout.Box) bool {
		if found != nil {
			return false
		}
		if b.Node == node {
			found = b
			return false
		}
		return true
	})
	return found
}

func (a *App) focusableAncestor(b *layout.Box) *layout.Box {
	for ; b != nil; b = b.Parent {
		if b.Node != nil && b.Style.Focusable {
			return b
		}
	}
	return nil
}

// scrollingAncestor finds the scroll pane a box sits inside, or nil when it
// sits in none.
func scrollingAncestor(b *layout.Box) *layout.Box {
	for ; b != nil; b = b.Parent {
		if b.Node != nil && b.Style.Overflow == style.OverflowScroll && b.ScrollMax > 0 {
			return b
		}
	}
	return nil
}

// scrollableAncestor is scrollingAncestor for the wheel and the arrows, where
// falling back to the one pane on screen is friendlier than doing nothing.
func (a *App) scrollableAncestor(b *layout.Box) *layout.Box {
	if found := scrollingAncestor(b); found != nil {
		return found
	}
	return a.firstScrollable()
}

func (a *App) firstScrollable() *layout.Box {
	tree := a.boxes()
	if tree == nil {
		return nil
	}

	var found *layout.Box
	tree.Walk(func(b *layout.Box) bool {
		if found != nil {
			return false
		}
		if b.Node != nil && b.Style.Overflow == style.OverflowScroll && b.ScrollMax > 0 {
			found = b
			return false
		}
		return true
	})
	return found
}

func (a *App) focusVisible() bool {
	box := a.findBox(a.focus)
	if box == nil {
		return false
	}

	// only a real ancestor can scroll the focused box into view: the fallback
	// pane would be some unrelated box, and a pane cannot reveal itself, which
	// used to walk a focused scroll pane back to the top one row per redraw
	pane := scrollingAncestor(box.Parent)
	if pane == nil || pane.Node == nil {
		return false
	}

	before := pane.Node.ScrollY

	if top := box.Frame.Y - pane.Content.Y; top < 0 {
		pane.Node.ScrollY += top
	} else if bottom := box.Frame.Bottom() - pane.Content.Bottom(); bottom > 0 {
		pane.Node.ScrollY += bottom
	}

	if pane.Node.ScrollY < 0 {
		pane.Node.ScrollY = 0
	}
	if pane.Node.ScrollY > pane.ScrollMax {
		pane.Node.ScrollY = pane.ScrollMax
	}

	return pane.Node.ScrollY != before
}

func (a *App) focusables() []*dom.Node {
	tree := a.boxes()
	if tree == nil {
		return nil
	}

	var out []*dom.Node
	tree.Walk(func(b *layout.Box) bool {
		if b.Node != nil && b.Style.Focusable {
			out = append(out, b.Node)
		}
		return true
	})
	return out
}

func (a *App) cycleFocus(delta int) {
	a.focusWithin(a.focusables(), delta)
}

func (a *App) focusWithin(nodes []*dom.Node, delta int) {
	if len(nodes) == 0 {
		return
	}

	index := -1
	for i, n := range nodes {
		if n == a.focus {
			index = i
			break
		}
	}

	next := 0
	if index >= 0 {
		next = (index + delta + len(nodes)) % len(nodes)
	} else if delta < 0 {
		next = len(nodes) - 1
	}

	a.Focus(nodes[next])
}

func (a *App) modalFocusables() []*dom.Node {
	box := a.findBox(a.Query("#" + modalID))
	if box == nil {
		return nil
	}
	return focusablesIn(box)
}

func (a *App) scope() []*dom.Node {
	for box := a.findBox(a.focus); box != nil; box = box.Parent {
		if box.Node == nil || box.Style.Shortcut == "" {
			continue
		}
		if inside := focusablesIn(box); len(inside) > 1 {
			return inside
		}
		break
	}
	return a.links()
}

func focusablesIn(section *layout.Box) []*dom.Node {
	var out []*dom.Node

	section.Walk(func(b *layout.Box) bool {
		if b == section || b.Node == nil {
			return true
		}
		if b.Style.Focusable && b.Style.Overflow != style.OverflowScroll {
			out = append(out, b.Node)
		}
		return true
	})

	return out
}

func (a *App) links() []*dom.Node {
	tree := a.boxes()
	if tree == nil {
		return nil
	}

	var out []*dom.Node
	tree.Walk(func(b *layout.Box) bool {
		if b.Node != nil && b.Style.Focusable && b.Style.Overflow != style.OverflowScroll {
			out = append(out, b.Node)
		}
		return true
	})
	return out
}
