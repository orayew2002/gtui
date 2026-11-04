package style

import (
	"strings"

	"github.com/orayew2002/gtui/src/dom"
)

type Diagnostic struct {
	Node  *dom.Node
	Class string
}

func Compute(n *dom.Node, parent Style, diags *[]Diagnostic) Style {
	s := Inherit(parent)

	if n.Type != dom.ElementNode {
		return s
	}

	applyTagDefaults(&s, n.Tag)

	for _, class := range n.Classes {
		if strings.ContainsRune(class, ':') {
			continue
		}
		if Apply(&s, class) {
			continue
		}
		if diags != nil {
			*diags = append(*diags, Diagnostic{Node: n, Class: class})
		}
	}

	applyAttrs(&s, n)
	applyState(&s, n)

	return s
}

func applyTagDefaults(s *Style, tag string) {
	switch tag {
	case "div", dom.RootTag:
		s.Display = DisplayBlock
	default:
		s.Display = DisplayBlock
	}
}

func applyAttrs(s *Style, n *dom.Node) {
	if t := n.Get("title"); t != "" {
		s.Title = t
		s.defaultBorder()
	}

	if inline := n.Get("style"); inline != "" {
		for _, class := range strings.Fields(inline) {
			Apply(s, class)
		}
	}

	if n.Has("hidden") {
		s.Display = DisplayNone
	}
	if n.Has("tabindex") {
		s.Focusable = true
	}
}

func applyState(s *Style, n *dom.Node) {
	for _, class := range n.Classes {
		variant, base, ok := strings.Cut(class, ":")
		if !ok {
			continue
		}
		switch variant {
		case "hover":
			if n.State.Hover {
				Apply(s, base)
			}
		case "focus":
			if n.State.Focus {
				Apply(s, base)
			}
		case "active":
			if n.State.Active {
				Apply(s, base)
			}
		case "current":
			if n.State.Current {
				Apply(s, base)
			}
		}
	}
}
