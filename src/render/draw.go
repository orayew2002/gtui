package render

import "io"

func Draw(w io.Writer, buf *Buffer, originX, originY int) error {
	if buf == nil {
		return nil
	}

	p := newPen(w)

	for y := range buf.H {
		row := buf.Cells[y*buf.W : (y+1)*buf.W]

		for x := range buf.W {
			cell := row[x]
			if cell.Cont {
				continue
			}

			screenX, screenY := originX+x, originY+y
			if screenX < 0 || screenY < 0 {
				continue
			}

			p.moveTo(screenX, screenY)
			p.writeCell(cell)
		}

		p.forgetPosition()
	}

	p.writeString("\x1b[0m")
	return p.flush()
}
