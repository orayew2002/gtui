package dom

import "strings"

func Query(root *Node, sel string) *Node {
	var found *Node
	queryAll(root, sel, func(n *Node) bool {
		found = n
		return false
	})
	return found
}

func QueryAll(root *Node, sel string) []*Node {
	var out []*Node
	queryAll(root, sel, func(n *Node) bool {
		out = append(out, n)
		return true
	})
	return out
}

func queryAll(root *Node, sel string, visit func(*Node) bool) {
	groups := parseSelector(sel)
	if len(groups) == 0 || root == nil {
		return
	}

	stopped := false

	root.Walk(func(n *Node) bool {
		if stopped {
			return false
		}
		if n.Type != ElementNode || n == root {
			return true
		}

		for _, g := range groups {
			if g.matches(n) {
				if !visit(n) {
					stopped = true
					return false
				}
				break
			}
		}
		return true
	})
}

func Matches(n *Node, sel string) bool {
	for _, g := range parseSelector(sel) {
		if g.matches(n) {
			return true
		}
	}
	return false
}

type compound struct {
	tag     string
	id      string
	classes []string
}

type chain []compound

func parseSelector(sel string) []chain {
	var out []chain

	for _, alt := range strings.Split(sel, ",") {
		parts := strings.Fields(alt)
		if len(parts) == 0 {
			continue
		}

		c := make(chain, 0, len(parts))
		for i := len(parts) - 1; i >= 0; i-- {
			c = append(c, parseCompound(parts[i]))
		}
		out = append(out, c)
	}

	return out
}

func parseCompound(s string) compound {
	var c compound

	for s != "" {
		i := strings.IndexAny(s[1:], ".#")
		var part string
		if i < 0 {
			part, s = s, ""
		} else {
			part, s = s[:i+1], s[i+1:]
		}

		switch part[0] {
		case '.':
			c.classes = append(c.classes, part[1:])
		case '#':
			c.id = part[1:]
		default:
			if part != "*" {
				c.tag = strings.ToLower(part)
			}
		}
	}

	return c
}

func (c compound) matches(n *Node) bool {
	if n.Type != ElementNode {
		return false
	}
	if c.tag != "" && c.tag != n.Tag {
		return false
	}
	if c.id != "" && c.id != n.ID {
		return false
	}
	for _, class := range c.classes {
		if !n.HasClass(class) {
			return false
		}
	}
	return true
}

func (ch chain) matches(n *Node) bool {
	if len(ch) == 0 || !ch[0].matches(n) {
		return false
	}

	node := n.Parent
	for _, c := range ch[1:] {
		for node != nil && !c.matches(node) {
			node = node.Parent
		}
		if node == nil {
			return false
		}
		node = node.Parent
	}

	return true
}
