package layout

import (
	"strings"

	"github.com/orayew2002/gtui/src/render"
)

type span struct {
	start, end, w int
}

func (b *Box) wrapText(width int, nowrap bool) {
	b.runBuf = b.runBuf[:0]
	b.spans = b.spans[:0]
	b.lineBuf = b.lineBuf[:0]

	if len(b.Source) == 0 {
		b.Lines = nil
		return
	}

	cur, lineW, pending := span{}, 0, 0

	flush := func() {
		cur.end = len(b.runBuf)
		cur.w = lineW
		b.spans = append(b.spans, cur)
		cur = span{start: len(b.runBuf)}
		lineW = 0
	}

	for i := range b.Source {
		src := &b.Source[i]
		st, node, text := src.Style, src.Node, src.Text
		if st == nil {
			st = &b.Style
		}

		segStart, segEnd, segW, segX := -1, 0, 0, 0

		closeSeg := func() {
			if segStart < 0 {
				return
			}
			b.runBuf = append(b.runBuf, Run{
				Text:  text[segStart:segEnd],
				X:     segX,
				W:     segW,
				Style: st,
				Node:  node,
			})
			segStart, segEnd, segW = -1, 0, 0
		}

		breakLine := func() {
			closeSeg()
			pending = 0
			flush()
		}

		eachToken(text, st.Pre, func(tok token) {
			switch {
			case tok.br:
				breakLine()

			case tok.space:
				if lineW == 0 {
					return
				}
				pending = tok.w

			default:
				if !nowrap && width > 0 && lineW+pending+tok.w > width && lineW > 0 {
					breakLine()
				}

				if !nowrap && width > 0 && tok.w > width {
					closeSeg()
					rest, offset := tok.text, tok.start

					for render.StringWidth(rest) > width-lineW {
						head, headW := render.TruncateWidth(rest, width-lineW)
						if headW == 0 {
							if lineW == 0 {
								break
							}
							flush()
							continue
						}

						b.runBuf = append(b.runBuf, Run{
							Text: text[offset : offset+len(head)], X: lineW,
							W: headW, Style: st, Node: node,
						})
						lineW += headW
						rest = rest[len(head):]
						offset += len(head)
						flush()
					}

					if rest != "" {
						segStart, segEnd = offset, offset+len(rest)
						segX, segW = lineW, render.StringWidth(rest)
						lineW += segW
					}
					pending = 0
					return
				}

				if segStart < 0 {
					segStart, segX = tok.start, lineW+pending
				} else {
					segEnd = tok.start
					segW += pending
				}

				segEnd = tok.end
				segW += tok.w
				lineW += pending + tok.w
				pending = 0
			}
		})

		closeSeg()
	}

	flush()

	for _, sp := range b.spans {
		b.lineBuf = append(b.lineBuf, Line{Runs: b.runBuf[sp.start:sp.end], W: sp.w})
	}
	b.Lines = b.lineBuf
}

type Line struct {
	Runs []Run
	W    int
}

type token struct {
	text  string
	start int
	end   int
	w     int
	space bool
	br    bool
}

func eachToken(text string, pre bool, fn func(token)) {
	if pre {
		offset := 0
		for {
			i := strings.IndexByte(text[offset:], '\n')
			if i < 0 {
				break
			}
			if i > 0 {
				line := text[offset : offset+i]
				fn(token{text: line, start: offset, end: offset + i, w: render.StringWidth(line)})
			}
			fn(token{br: true})
			offset += i + 1
		}
		if offset < len(text) {
			rest := text[offset:]
			fn(token{text: rest, start: offset, end: len(text), w: render.StringWidth(rest)})
		}
		return
	}

	start := -1
	emitWord := func(end int) {
		if start < 0 {
			return
		}
		word := text[start:end]
		fn(token{text: word, start: start, end: end, w: render.StringWidth(word)})
		start = -1
	}

	for i := 0; i < len(text); i++ {
		switch text[i] {
		case ' ':
			emitWord(i)
			fn(token{text: " ", start: i, end: i + 1, w: 1, space: true})
		case '\n':
			emitWord(i)
			fn(token{br: true})
		default:
			if start < 0 {
				start = i
			}
		}
	}
	emitWord(len(text))
}

func maxLineWidth(src []Run) int {
	widest, cur := 0, 0

	for i := range src {
		text := src[i].Text
		for {
			j := strings.IndexByte(text, '\n')
			if j < 0 {
				cur += render.StringWidth(text)
				break
			}
			cur += render.StringWidth(text[:j])
			if cur > widest {
				widest = cur
			}
			cur = 0
			text = text[j+1:]
		}
	}

	if cur > widest {
		widest = cur
	}
	return widest
}
