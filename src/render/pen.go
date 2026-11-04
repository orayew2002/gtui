package render

import (
	"bufio"
	"io"
	"strconv"
	"unicode/utf8"

	"github.com/orayew2002/gtui/src/colors"
	"github.com/orayew2002/gtui/src/style"
)

const writeBuffer = 1 << 16

type pen struct {
	out *bufio.Writer

	fg, bg       colors.Color
	hasFg, hasBg bool
	attr         style.Attrs

	cx, cy int

	scratch []byte
}

func newPen(w io.Writer) *pen {
	return &pen{
		out:     bufio.NewWriterSize(w, writeBuffer),
		cx:      -1,
		cy:      -1,
		scratch: make([]byte, 0, 64),
	}
}

func (p *pen) reset() {
	p.hasFg, p.hasBg, p.attr = false, false, 0
	p.cx, p.cy = -1, -1
}

func (p *pen) forgetPosition() { p.cx, p.cy = -1, -1 }

func (p *pen) at() (int, int) { return p.cx, p.cy }

func (p *pen) flush() error { return p.out.Flush() }

func (p *pen) writeString(s string) { p.out.WriteString(s) }

func (p *pen) write(b []byte) (int, error) { return p.out.Write(b) }

func (p *pen) clearScreen() {
	p.writeString("\x1b[0m\x1b[2J")
	p.reset()
}

func (p *pen) moveTo(x, y int) {
	if p.cx == x && p.cy == y {
		return
	}

	if p.cy == y && x == p.cx+1 {
		p.writeString("\x1b[C")
		p.cx, p.cy = x, y
		return
	}

	if p.cy == y && x == 0 {
		p.writeString("\r")
		p.cx, p.cy = 0, y
		return
	}

	p.scratch = p.scratch[:0]
	p.scratch = append(p.scratch, 0x1b, '[')
	p.scratch = strconv.AppendInt(p.scratch, int64(y+1), 10)
	p.scratch = append(p.scratch, ';')
	p.scratch = strconv.AppendInt(p.scratch, int64(x+1), 10)
	p.scratch = append(p.scratch, 'H')
	p.out.Write(p.scratch)

	p.cx, p.cy = x, y
}

func (p *pen) setPen(c Cell) {
	needReset := p.attr&^c.Attr != 0 ||
		(p.hasFg && !c.HasFg) ||
		(p.hasBg && !c.HasBg)

	if needReset {
		p.writeString("\x1b[0m")
		p.hasFg, p.hasBg, p.attr = false, false, 0
	}

	if add := c.Attr &^ p.attr; add != 0 {
		p.writeAttrs(add)
		p.attr = c.Attr
	}

	if c.HasFg && (!p.hasFg || p.fg != c.Fg) {
		p.writeColor(38, c.Fg)
		p.fg, p.hasFg = c.Fg, true
	}

	if c.HasBg && (!p.hasBg || p.bg != c.Bg) {
		p.writeColor(48, c.Bg)
		p.bg, p.hasBg = c.Bg, true
	}
}

func (p *pen) writeAttrs(a style.Attrs) {
	for _, code := range attrCodes {
		if a.Has(code.flag) {
			p.writeString(code.seq)
		}
	}
}

var attrCodes = [...]struct {
	flag style.Attrs
	seq  string
}{
	{style.AttrBold, "\x1b[1m"},
	{style.AttrDim, "\x1b[2m"},
	{style.AttrItalic, "\x1b[3m"},
	{style.AttrUnderline, "\x1b[4m"},
	{style.AttrBlink, "\x1b[5m"},
	{style.AttrReverse, "\x1b[7m"},
	{style.AttrStrike, "\x1b[9m"},
}

func (p *pen) writeColor(base int, c colors.Color) {
	p.scratch = p.scratch[:0]
	p.scratch = append(p.scratch, 0x1b, '[')
	p.scratch = strconv.AppendInt(p.scratch, int64(base), 10)
	p.scratch = append(p.scratch, ';', '2', ';')
	p.scratch = strconv.AppendInt(p.scratch, int64(c.R), 10)
	p.scratch = append(p.scratch, ';')
	p.scratch = strconv.AppendInt(p.scratch, int64(c.G), 10)
	p.scratch = append(p.scratch, ';')
	p.scratch = strconv.AppendInt(p.scratch, int64(c.B), 10)
	p.scratch = append(p.scratch, 'm')
	p.out.Write(p.scratch)
}

func (p *pen) writeCell(c Cell) {
	p.setPen(c)
	p.writeRune(c.Ch)

	advance := RuneWidth(c.Ch)
	if advance < 1 {
		advance = 1
	}
	p.cx += advance
}

func (p *pen) writeRune(r rune) {
	if r == 0 || r < 0x20 {
		p.out.WriteByte(' ')
		return
	}
	if r < utf8.RuneSelf {
		p.out.WriteByte(byte(r))
		return
	}
	p.out.WriteRune(r)
}
