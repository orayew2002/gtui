package styles

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"strconv"
	"sync/atomic"

	"golang.org/x/sys/unix"
)

const (
	chunkSize      = 4096
	defaultCellW   = 10
	defaultCellH   = 20
	maxImagePixels = 16 << 20
)

var (
	ErrNoColor  = errors.New("styles: pixel box needs a background or border color")
	ErrTooLarge = errors.New("styles: pixel box exceeds the image size limit")
	ErrEmpty    = errors.New("styles: pixel box has no area")
)

var nextImageID atomic.Uint32

type Pixel struct {
	box         *Box
	radius      float64
	borderWidth float64
	id          uint32
	cellW       int
	cellH       int
}

func (b *Box) Pixel() *Pixel {
	return &Pixel{
		box:         b,
		radius:      4,
		borderWidth: 1,
		id:          nextImageID.Add(1),
	}
}

func (p *Pixel) Radius(r float64) *Pixel {
	p.radius = math.Max(0, r)
	return p
}

func (p *Pixel) BorderWidth(w float64) *Pixel {
	p.borderWidth = math.Max(0, w)
	return p
}

func (p *Pixel) Cell(w, h int) *Pixel {
	p.cellW, p.cellH = w, h
	return p
}

func (p *Pixel) ID() uint32 { return p.id }

func (p *Pixel) cellSize() (int, int) {
	if p.cellW > 0 && p.cellH > 0 {
		return p.cellW, p.cellH
	}

	w, h, err := CellSize()
	if err != nil {
		return defaultCellW, defaultCellH
	}
	return w, h
}

func CellSize() (int, int, error) {
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return defaultCellW, defaultCellH, err
	}
	if ws.Xpixel == 0 || ws.Ypixel == 0 || ws.Col == 0 || ws.Row == 0 {
		return defaultCellW, defaultCellH, errors.New("styles: terminal reports no pixel size")
	}

	return int(ws.Xpixel) / int(ws.Col), int(ws.Ypixel) / int(ws.Row), nil
}

func (p *Pixel) Image() (*image.NRGBA, error) {
	cellW, cellH := p.cellSize()
	w, h := p.box.w*cellW, p.box.h*cellH

	if w <= 0 || h <= 0 {
		return nil, ErrEmpty
	}
	if w*h > maxImagePixels {
		return nil, ErrTooLarge
	}

	st := p.box.Style()
	if !st.Background.Set && !st.BorderColor.Set {
		return nil, ErrNoColor
	}

	fill := color.NRGBA{}
	if st.Background.Set {
		c := st.Background.Color
		fill = color.NRGBA{R: c.R, G: c.G, B: c.B, A: 255}
	}

	border := fill
	if st.BorderColor.Set {
		c := st.BorderColor.Color
		border = color.NRGBA{R: c.R, G: c.G, B: c.B, A: 255}
	}

	if !st.Background.Set {
		fill = color.NRGBA{R: border.R, G: border.G, B: border.B}
	}

	radius := math.Min(p.radius, math.Min(float64(w), float64(h))/2)
	stroke := math.Min(p.borderWidth, math.Min(float64(w), float64(h))/2)

	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	fw, fh := float64(w), float64(h)

	for y := range h {
		for x := range w {
			d := roundedRectDistance(float64(x)+0.5, float64(y)+0.5, fw, fh, radius)

			coverage := clamp01(0.5 - d)
			if coverage == 0 {
				continue
			}

			inside := clamp01(0.5 - (d + stroke))
			c := mix(border, fill, inside)
			c.A = scaleAlpha(c.A, coverage)
			img.SetNRGBA(x, y, c)
		}
	}

	return img, nil
}

func (p *Pixel) Render(w io.Writer) error {
	img, err := p.Image()
	if err != nil {
		return err
	}

	var png bytes.Buffer
	if err := encodePNG(&png, img); err != nil {
		return err
	}

	out := bufio.NewWriter(w)
	moveTo(out, p.box.x, p.box.y)
	p.transmit(out, png.Bytes())

	return out.Flush()
}

func (p *Pixel) Delete(w io.Writer) error {
	out := bufio.NewWriter(w)
	out.WriteString("\x1b_Ga=d,d=i,i=")
	out.WriteString(strconv.FormatUint(uint64(p.id), 10))
	out.WriteString("\x1b\\")
	return out.Flush()
}

func (p *Pixel) transmit(out *bufio.Writer, data []byte) {
	encoded := base64.StdEncoding.EncodeToString(data)
	first := true

	for len(encoded) > 0 {
		n := min(chunkSize, len(encoded))
		chunk := encoded[:n]
		encoded = encoded[n:]

		more := 0
		if len(encoded) > 0 {
			more = 1
		}

		out.WriteString("\x1b_G")
		if first {
			out.WriteString("a=T,f=100,q=2,z=-1,C=1,i=")
			out.WriteString(strconv.FormatUint(uint64(p.id), 10))
			out.WriteString(",c=")
			out.WriteString(strconv.Itoa(p.box.w))
			out.WriteString(",r=")
			out.WriteString(strconv.Itoa(p.box.h))
			out.WriteString(",m=")
			first = false
		} else {
			out.WriteString("m=")
		}
		out.WriteString(strconv.Itoa(more))
		out.WriteByte(';')
		out.WriteString(chunk)
		out.WriteString("\x1b\\")
	}
}

func encodePNG(buf *bytes.Buffer, img *image.NRGBA) error {
	enc := png.Encoder{CompressionLevel: png.BestSpeed}
	return enc.Encode(buf, img)
}

func moveTo(out *bufio.Writer, x, y int) {
	out.WriteString("\x1b[")
	out.WriteString(strconv.Itoa(y + 1))
	out.WriteByte(';')
	out.WriteString(strconv.Itoa(x + 1))
	out.WriteByte('H')
}

func roundedRectDistance(px, py, w, h, r float64) float64 {
	qx := math.Abs(px-w/2) - (w/2 - r)
	qy := math.Abs(py-h/2) - (h/2 - r)

	outside := math.Hypot(math.Max(qx, 0), math.Max(qy, 0))
	inside := math.Min(math.Max(qx, qy), 0)

	return outside + inside - r
}

func mix(a, b color.NRGBA, t float64) color.NRGBA {
	if t <= 0 {
		return a
	}
	if t >= 1 {
		return b
	}

	lerp := func(x, y uint8) uint8 {
		return uint8(float64(x) + (float64(y)-float64(x))*t + 0.5)
	}

	return color.NRGBA{
		R: lerp(a.R, b.R),
		G: lerp(a.G, b.G),
		B: lerp(a.B, b.B),
		A: lerp(a.A, b.A),
	}
}

func scaleAlpha(a uint8, coverage float64) uint8 {
	return uint8(float64(a)*coverage + 0.5)
}

func clamp01(v float64) float64 {
	return math.Max(0, math.Min(1, v))
}
