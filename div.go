package gtui

import (
	"fmt"
	"strings"

	"github.com/orayew2002/gtui/src/dom"
)

type Component struct {
	class   string
	text    string
	target  string
	page    []Component
	child   []Component
	link    bool
	onClick func(*Event)
	build   func(*App) *Node
}

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
		return c.build(a)
	}

	n := dom.NewElement("div")

	class, id := splitClass(c.class)
	if c.link {
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

	if c.onClick != nil {
		n.On("click", c.onClick)
	}

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
