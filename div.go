package gtui

import (
	"fmt"
	"strings"

	"github.com/orayew2002/gtui/src/dom"
)

type Component struct {
	class    string
	text     string
	target   string
	page     []Component
	child    []Component
	link     bool
	onClick  func(*Event)
	handlers []listener
	build    func(*App) *Node
}

type listener struct {
	event string
	fn    func(*Event)
}

// On attaches a handler to the component itself, so it travels with the
// description and comes back with every rebuild. It returns a copy: the
// component it is called on is left alone.
func (c Component) On(event string, fn func(*Event)) Component {
	if event == "" || fn == nil {
		return c
	}

	c.handlers = append(c.handlers[:len(c.handlers):len(c.handlers)], listener{strings.ToLower(event), fn})
	return c
}

func (c Component) OnClick(fn func(*Event)) Component { return c.On("click", fn) }

func Div(class string, child ...Component) Component {
	return Component{class: class, child: child}
}

func Text(class, content string, args ...any) Component {
	if len(args) > 0 {
		content = fmt.Sprintf(content, args...)
	}
	return Component{class: class, text: content}
}

func Row(class string, child ...Component) Component {
	return Div(join("row", class), child...)
}

func Column(class string, child ...Component) Component {
	return Div(join("col", class), child...)
}

func Center(class string, child ...Component) Component {
	return Div(join("center", class), child...)
}

func Spacer() Component {
	return Div("grow")
}

func Href(class, title, target string, page ...Component) Component {
	return Component{
		class:  class,
		text:   title,
		target: target,
		page:   page,
		link:   true,
	}
}

func (c Component) Node() *Node {
	return c.node(nil)
}

func (c Component) node(a *App) *Node {
	if c.build != nil {
		return c.attach(c.build(a))
	}

	n := dom.NewElement("div")

	class, id := splitClass(c.class)
	if c.link || c.clickable() {
		class = withClass(class, "focusable")
	}
	if class != "" {
		n.Set("class", class)
	}
	if id != "" {
		n.Set("id", id)
	}

	if c.text != "" {
		n.Append(dom.NewText(c.text))
	}
	for _, child := range c.child {
		n.Append(child.node(a))
	}

	c.attach(n)

	if c.link {
		n.Set("href", c.target)
		if a != nil {
			target, page := c.target, c.page
			link := n
			n.On("click", func(e *Event) {
				a.mark(link)
				a.Open(target, page...)
			})
		}
	}

	return n
}

// attach wires the component's own handlers onto the node it just built, which
// is what lets a plain div, a link and a build-backed control all take a click.
func (c Component) attach(n *Node) *Node {
	if n == nil {
		return nil
	}

	if c.onClick != nil {
		n.On("click", c.onClick)
	}
	for _, h := range c.handlers {
		n.On(h.event, h.fn)
	}

	return n
}

func (c Component) clickable() bool {
	if c.onClick != nil {
		return true
	}
	for _, h := range c.handlers {
		if h.event == "click" {
			return true
		}
	}
	return false
}

func splitClass(class string) (string, string) {
	if !strings.ContainsRune(class, '#') {
		return class, ""
	}

	var (
		keep []string
		id   string
	)

	for _, token := range strings.Fields(class) {
		if strings.HasPrefix(token, "#") {
			id = token[1:]
			continue
		}
		keep = append(keep, token)
	}

	return strings.Join(keep, " "), id
}

func withClass(class, extra string) string {
	for _, token := range strings.Fields(class) {
		if token == extra {
			return class
		}
	}
	if class == "" {
		return extra
	}
	return class + " " + extra
}

func (a *App) nodes(components []Component) []*Node {
	out := make([]*Node, 0, len(components))
	for _, c := range components {
		out = append(out, c.node(a))
	}
	return out
}
