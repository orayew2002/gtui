package html

import (
	"strings"

	"github.com/orayew2002/gtui/src/dom"
)

const FragmentTag = "fragment"

var voidElements = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"source": true, "track": true, "wbr": true,

	"spacer": true, "gauge": true, "sparkline": true,
}

var rawTextElements = map[string]bool{
	"style": true, "script": true, "textarea": true,
}

var closedBy = map[string]map[string]bool{
	"p":  {"p": true, "div": true, "ul": true, "ol": true, "li": true, "h1": true, "h2": true, "h3": true, "h4": true, "hr": true, "pre": true, "table": true, "section": true, "header": true, "footer": true},
	"li": {"li": true},
	"dt": {"dt": true, "dd": true},
	"dd": {"dt": true, "dd": true},
	"tr": {"tr": true},
	"td": {"td": true, "th": true, "tr": true},
	"th": {"td": true, "th": true, "tr": true},
}

type Document struct {
	Root   *dom.Node
	Styles []string
}

func Parse(src string) *Document {
	p := &parser{src: src}
	root := dom.NewElement(FragmentTag)
	p.stack = []*dom.Node{root}
	p.run()

	return &Document{Root: root, Styles: p.styles}
}

type parser struct {
	src    string
	pos    int
	stack  []*dom.Node
	styles []string
}

func (p *parser) top() *dom.Node { return p.stack[len(p.stack)-1] }

func (p *parser) push(n *dom.Node) { p.stack = append(p.stack, n) }

func (p *parser) pop() {
	if len(p.stack) > 1 {
		p.stack = p.stack[:len(p.stack)-1]
	}
}

func (p *parser) run() {
	for p.pos < len(p.src) {
		if p.src[p.pos] != '<' {
			p.text()
			continue
		}

		switch {
		case p.lookahead("<!--"):
			p.comment()
		case p.lookahead("<!"), p.lookahead("<?"):
			p.skipTo(">")
		case p.lookahead("</"):
			p.endTag()
		case p.pos+1 < len(p.src) && isNameStart(p.src[p.pos+1]):
			p.startTag()
		default:

			p.top().Append(dom.NewText("<"))
			p.pos++
		}
	}
}

func (p *parser) lookahead(s string) bool {
	return strings.HasPrefix(p.src[p.pos:], s)
}

func (p *parser) skipTo(delim string) {
	if i := strings.Index(p.src[p.pos:], delim); i >= 0 {
		p.pos += i + len(delim)
		return
	}
	p.pos = len(p.src)
}

func (p *parser) text() {
	start := p.pos
	if i := strings.IndexByte(p.src[p.pos:], '<'); i >= 0 {
		p.pos += i
	} else {
		p.pos = len(p.src)
	}

	raw := p.src[start:p.pos]
	if raw == "" {
		return
	}
	p.top().Append(dom.NewText(Unescape(raw)))
}

func (p *parser) comment() {
	p.pos += len("<!--")
	start := p.pos
	if i := strings.Index(p.src[p.pos:], "-->"); i >= 0 {
		p.pos += i
		p.top().Append(&dom.Node{Type: dom.CommentNode, Data: p.src[start:p.pos]})
		p.pos += len("-->")
		return
	}
	p.pos = len(p.src)
}

func (p *parser) endTag() {
	p.pos += len("</")
	name := p.name()
	p.skipTo(">")

	if name == "" {
		return
	}

	for i := len(p.stack) - 1; i > 0; i-- {
		if p.stack[i].Tag == name {
			p.stack = p.stack[:i]
			return
		}
	}
}

func (p *parser) startTag() {
	p.pos++
	name := p.name()
	if name == "" {
		return
	}

	el := dom.NewElement(name)
	selfClosing := p.attributes(el)

	if set := closedBy[p.top().Tag]; set != nil && set[name] {
		p.pop()
	}

	p.top().Append(el)

	if selfClosing || voidElements[name] {
		return
	}

	if rawTextElements[name] {
		p.rawText(el)
		return
	}

	p.push(el)
}

func (p *parser) rawText(el *dom.Node) {
	closing := "</" + el.Tag
	start := p.pos

	if i := strings.Index(strings.ToLower(p.src[p.pos:]), closing); i >= 0 {
		p.pos += i
		el.Append(dom.NewText(p.src[start:p.pos]))
		p.skipTo(">")
	} else {
		el.Append(dom.NewText(p.src[start:]))
		p.pos = len(p.src)
	}

	if el.Tag == "style" {
		p.styles = append(p.styles, el.Text())
		el.Detach()
	}
}

func (p *parser) attributes(el *dom.Node) bool {
	for p.pos < len(p.src) {
		p.skipSpace()
		if p.pos >= len(p.src) {
			return false
		}

		switch p.src[p.pos] {
		case '>':
			p.pos++
			return false
		case '/':
			p.pos++
			p.skipSpace()
			if p.pos < len(p.src) && p.src[p.pos] == '>' {
				p.pos++
			}
			return true
		}

		key := p.attrName()
		if key == "" {
			p.pos++
			continue
		}

		value := ""
		p.skipSpace()
		if p.pos < len(p.src) && p.src[p.pos] == '=' {
			p.pos++
			p.skipSpace()
			value = p.attrValue()
		}

		el.Set(key, Unescape(value))
	}

	return false
}

func (p *parser) attrValue() string {
	if p.pos >= len(p.src) {
		return ""
	}

	switch q := p.src[p.pos]; q {
	case '"', '\'':
		p.pos++
		start := p.pos
		if i := strings.IndexByte(p.src[p.pos:], q); i >= 0 {
			p.pos += i
			v := p.src[start:p.pos]
			p.pos++
			return v
		}
		p.pos = len(p.src)
		return p.src[start:]
	}

	start := p.pos
	for p.pos < len(p.src) && !isSpace(p.src[p.pos]) && p.src[p.pos] != '>' {
		p.pos++
	}
	return p.src[start:p.pos]
}

func (p *parser) name() string {
	start := p.pos
	for p.pos < len(p.src) && isNameByte(p.src[p.pos]) {
		p.pos++
	}
	return strings.ToLower(p.src[start:p.pos])
}

func (p *parser) attrName() string {
	start := p.pos
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		if isSpace(c) || c == '=' || c == '>' || c == '/' {
			break
		}
		p.pos++
	}
	return strings.ToLower(p.src[start:p.pos])
}

func (p *parser) skipSpace() {
	for p.pos < len(p.src) && isSpace(p.src[p.pos]) {
		p.pos++
	}
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f'
}

func isNameStart(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isNameByte(c byte) bool {
	return isNameStart(c) || c >= '0' && c <= '9' || c == '-' || c == '_' || c == ':'
}
