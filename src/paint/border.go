package paint

import "github.com/orayew2002/gtui/src/style"

type glyphs struct {
	tl, tr, bl, br rune
	h, v           rune

	top, bottom rune
	left, right rune
}

func borderGlyphs(s style.BorderStyle) glyphs {
	switch s {
	case style.BorderRounded:
		return glyphs{tl: '╭', tr: '╮', bl: '╰', br: '╯', h: '─', v: '│'}
	case style.BorderThick:
		return glyphs{tl: '┏', tr: '┓', bl: '┗', br: '┛', h: '━', v: '┃'}
	case style.BorderDouble:
		return glyphs{tl: '╔', tr: '╗', bl: '╚', br: '╝', h: '═', v: '║'}
	case style.BorderDashed:
		return glyphs{tl: '┌', tr: '┐', bl: '└', br: '┘', h: '┄', v: '┆'}
	case style.BorderSoft:
		return glyphs{
			tl: ' ', tr: ' ', bl: ' ', br: ' ',
			h: '▄', v: '▌',
			top: '▄', bottom: '▀', left: '▐', right: '▌',
		}
	case style.BorderTight:
		return glyphs{
			tl: ' ', tr: ' ', bl: ' ', br: ' ',
			h: '▁', v: '▏',
			top: '▁', bottom: '▔', left: '▕', right: '▏',
		}
	default:
		return glyphs{tl: '┌', tr: '┐', bl: '└', br: '┘', h: '─', v: '│'}
	}
}

func (g glyphs) topRune() rune {
	if g.top != 0 {
		return g.top
	}
	return g.h
}

func (g glyphs) bottomRune() rune {
	if g.bottom != 0 {
		return g.bottom
	}
	return g.h
}

func (g glyphs) leftRune() rune {
	if g.left != 0 {
		return g.left
	}
	return g.v
}

func (g glyphs) rightRune() rune {
	if g.right != 0 {
		return g.right
	}
	return g.v
}
