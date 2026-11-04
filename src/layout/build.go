package layout

import (
	"strings"

	"github.com/orayew2002/gtui/src/dom"
	"github.com/orayew2002/gtui/src/render"
	"github.com/orayew2002/gtui/src/style"
)

func Build(node *dom.Node, parent style.Style, diags *[]style.Diagnostic) *Box {
	st := style.Compute(node, parent, diags)
	if st.Display == style.DisplayNone {
		return nil
	}

	b := &Box{Node: node, Style: st, Ins: st.Insets()}
	buildChildren(b, node, st, diags)
	splitChildren(b)
	return b
}

func buildChildren(b *Box, node *dom.Node, parentStyle style.Style, diags *[]style.Diagnostic) {
	var pending []Run

	flush := func() {

		for len(pending) > 0 && isSpaceOnly(pending[len(pending)-1].Text) {
			pending = pending[:len(pending)-1]
		}
		if len(pending) == 0 {
			return
		}
		inline := &Box{
			Style:  style.Inherit(parentStyle),
			Source: pending,
			Parent: b,
		}
		inline.Style.Display = style.DisplayInline
		b.Children = append(b.Children, inline)
		pending = nil
	}

	for _, child := range node.Children {
		switch child.Type {
		case dom.CommentNode:
			continue

		case dom.TextNode:

			if !parentStyle.Pre && isSpaceOnly(child.Data) {
				if len(pending) > 0 {
					pending = append(pending, spaceRun(child, &parentStyle))
				}
				continue
			}

			if run, ok := textRun(child, &parentStyle); ok {
				pending = append(pending, run)
			}

		case dom.ElementNode:
			st := style.Compute(child, parentStyle, diags)
			if st.Display == style.DisplayNone {
				continue
			}

			if st.Display == style.DisplayInline {
				pending = append(pending, inlineRuns(child, st, diags)...)
				continue
			}

			flush()
			box := &Box{Node: child, Style: st, Parent: b, Ins: st.Insets()}
			buildChildren(box, child, st, diags)
			splitChildren(box)
			b.Children = append(b.Children, box)
		}
	}

	flush()
}

func splitChildren(b *Box) {
	if len(b.Children) == 0 {
		return
	}

	flow := make([]*Box, 0, len(b.Children))
	for _, c := range b.Children {
		if c.Style.Position == style.PositionAbsolute {
			b.Abs = append(b.Abs, c)
			continue
		}
		flow = append(flow, c)
	}
	b.Flow = flow
}

func inlineRuns(node *dom.Node, st style.Style, diags *[]style.Diagnostic) []Run {
	var out []Run

	shared := &st

	for _, child := range node.Children {
		switch child.Type {
		case dom.TextNode:
			if !st.Pre && isSpaceOnly(child.Data) {
				if len(out) > 0 {
					out = append(out, spaceRun(child, shared))
				}
				continue
			}

			if run, ok := textRun(child, shared); ok {
				run.Node = node
				out = append(out, run)
			}

		case dom.ElementNode:
			childStyle := style.Compute(child, st, diags)
			if childStyle.Display == style.DisplayNone {
				continue
			}
			out = append(out, inlineRuns(child, childStyle, diags)...)
		}
	}

	return out
}

func spaceRun(node *dom.Node, st *style.Style) Run {
	return Run{Text: " ", W: 1, Style: st, Node: node}
}

func plain(s string) bool {
	if s[0] == ' ' || s[len(s)-1] == ' ' {
		return false
	}

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\t', '\n', '\r', '\f':
			return false
		case ' ':
			if i+1 < len(s) && s[i+1] == ' ' {
				return false
			}
		}
	}

	return true
}

func isSpaceOnly(text string) bool {
	return strings.TrimLeft(text, " \t\n\r\f") == ""
}

func textRun(node *dom.Node, st *style.Style) (Run, bool) {
	text := node.Data
	if st.Pre {

		text = strings.ReplaceAll(text, "\t", "    ")
	}
	if !st.Pre {
		text = collapse(text)
		if text == "" {
			return Run{}, false
		}
	}
	if text == "" {
		return Run{}, false
	}

	return Run{Text: text, W: render.StringWidth(text), Style: st, Node: node}, true
}

func collapse(s string) string {
	if s == "" {
		return ""
	}

	if plain(s) {
		return s
	}

	var sb strings.Builder
	sb.Grow(len(s))

	space := false
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f' {
			space = true
			continue
		}
		if space {
			sb.WriteByte(' ')
		}
		space = false
		sb.WriteRune(r)
	}

	out := sb.String()
	if space && out != "" {
		out += " "
	}
	if out == " " {
		return ""
	}
	return out
}
