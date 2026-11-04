package dom

import (
	"strings"

	"github.com/orayew2002/gtui/src/colors"
)

type NodeType uint8

const (
	ElementNode NodeType = iota
	TextNode
	CommentNode
)

type State struct {
	Hover    bool
	Focus    bool
	Active   bool
	Current  bool
	Disabled bool
	Checked  bool
}

type Node struct {
	Type NodeType
	Tag  string
	Data string

	ID      string
	Classes []string
	Attr    map[string]string

	Parent   *Node
	Children []*Node

	State State

	ScrollY int

	Canvas *Canvas

	Handlers map[string][]func(*Event)
}

type Event struct {
	Type   string
	Target *Node
	Key    string
	Rune   rune
	X, Y   int
	Value  string
	stop   bool
}

func (e *Event) StopPropagation() { e.stop = true }
func (e *Event) Stopped() bool    { return e.stop }

func NewElement(tag string) *Node {
	return &Node{Type: ElementNode, Tag: strings.ToLower(tag), Attr: map[string]string{}}
}

func NewText(text string) *Node {
	return &Node{Type: TextNode, Data: text}
}

func (n *Node) IsElement() bool { return n.Type == ElementNode }
func (n *Node) IsText() bool    { return n.Type == TextNode }

func (n *Node) Append(child *Node) *Node {
	if child == nil {
		return n
	}
	child.Detach()
	child.Parent = n
	n.Children = append(n.Children, child)
	return n
}

func (n *Node) Insert(i int, child *Node) *Node {
	if child == nil {
		return n
	}
	child.Detach()
	if i < 0 {
		i = 0
	}
	if i > len(n.Children) {
		i = len(n.Children)
	}
	child.Parent = n
	n.Children = append(n.Children, nil)
	copy(n.Children[i+1:], n.Children[i:])
	n.Children[i] = child
	return n
}

func (n *Node) Detach() {
	p := n.Parent
	if p == nil {
		return
	}
	for i, c := range p.Children {
		if c == n {
			p.Children = append(p.Children[:i], p.Children[i+1:]...)
			break
		}
	}
	n.Parent = nil
}

func (n *Node) Clear() {
	for _, c := range n.Children {
		c.Parent = nil
	}
	n.Children = n.Children[:0]
}

func (n *Node) Get(key string) string {
	if n.Attr == nil {
		return ""
	}
	return n.Attr[strings.ToLower(key)]
}

func (n *Node) Has(key string) bool {
	if n.Attr == nil {
		return false
	}
	_, ok := n.Attr[strings.ToLower(key)]
	return ok
}

func (n *Node) Set(key, value string) *Node {
	key = strings.ToLower(key)
	if n.Attr == nil {
		n.Attr = map[string]string{}
	}
	n.Attr[key] = value

	switch key {
	case "id":
		n.ID = value
	case "class":
		n.Classes = strings.Fields(value)
	}
	return n
}

func (n *Node) HasClass(name string) bool {
	for _, c := range n.Classes {
		if c == name {
			return true
		}
	}
	return false
}

func (n *Node) AddClass(name string) *Node {
	if name == "" || n.HasClass(name) {
		return n
	}
	n.Classes = append(n.Classes, name)
	n.syncClassAttr()
	return n
}

func (n *Node) RemoveClass(name string) *Node {
	for i, c := range n.Classes {
		if c == name {
			n.Classes = append(n.Classes[:i], n.Classes[i+1:]...)
			n.syncClassAttr()
			break
		}
	}
	return n
}

func (n *Node) ToggleClass(name string) bool {
	if n.HasClass(name) {
		n.RemoveClass(name)
		return false
	}
	n.AddClass(name)
	return true
}

func (n *Node) syncClassAttr() {
	if n.Attr == nil {
		n.Attr = map[string]string{}
	}
	n.Attr["class"] = strings.Join(n.Classes, " ")
}

func (n *Node) Text() string {
	var sb strings.Builder
	n.writeText(&sb)
	return sb.String()
}

func (n *Node) writeText(sb *strings.Builder) {
	if n.Type == TextNode {
		sb.WriteString(n.Data)
		return
	}
	for _, c := range n.Children {
		c.writeText(sb)
	}
}

func (n *Node) SetText(text string) *Node {
	if n.Type == TextNode {
		n.Data = text
		return n
	}
	n.Clear()
	n.Append(NewText(text))
	return n
}

func (n *Node) On(event string, fn func(*Event)) *Node {
	if fn == nil {
		return n
	}
	if n.Handlers == nil {
		n.Handlers = map[string][]func(*Event){}
	}
	event = strings.ToLower(event)
	n.Handlers[event] = append(n.Handlers[event], fn)
	return n
}

func (n *Node) Dispatch(e *Event) {
	if e.Target == nil {
		e.Target = n
	}
	for node := n; node != nil; node = node.Parent {
		for _, fn := range node.Handlers[e.Type] {
			fn(e)
			if e.stop {
				return
			}
		}
	}
}

func (n *Node) Walk(fn func(*Node) bool) {
	if n == nil || !fn(n) {
		return
	}
	for _, c := range n.Children {
		c.Walk(fn)
	}
}

func (n *Node) ElementChildren() []*Node {
	out := make([]*Node, 0, len(n.Children))
	for _, c := range n.Children {
		if c.Type == ElementNode {
			out = append(out, c)
		}
	}
	return out
}

func (n *Node) Depth() int {
	d := 0
	for p := n.Parent; p != nil; p = p.Parent {
		d++
	}
	return d
}

func (n *Node) Clone() *Node {
	if n == nil {
		return nil
	}
	cp := &Node{
		Type:    n.Type,
		Tag:     n.Tag,
		Data:    n.Data,
		ID:      n.ID,
		Classes: append([]string(nil), n.Classes...),
		State:   n.State,
	}
	if n.Attr != nil {
		cp.Attr = make(map[string]string, len(n.Attr))
		for k, v := range n.Attr {
			cp.Attr[k] = v
		}
	}
	for _, c := range n.Children {
		cp.Append(c.Clone())
	}
	return cp
}

const RootTag = "root"

type Pixel struct {
	Ch     rune
	Fg, Bg colors.Color
	HasFg  bool
	HasBg  bool
	Bold   bool
}

type Canvas struct {
	W, H   int
	Pixels []Pixel
}

func NewCanvas(w, h int) *Canvas {
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	return &Canvas{W: w, H: h, Pixels: make([]Pixel, w*h)}
}

func (c *Canvas) Set(x, y int, p Pixel) {
	if c == nil || x < 0 || y < 0 || x >= c.W || y >= c.H {
		return
	}
	c.Pixels[y*c.W+x] = p
}

func (c *Canvas) At(x, y int) Pixel {
	if c == nil || x < 0 || y < 0 || x >= c.W || y >= c.H {
		return Pixel{}
	}
	return c.Pixels[y*c.W+x]
}
