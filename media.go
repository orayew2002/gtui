package gtui

import (
	"image"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
	"time"

	"github.com/orayew2002/gtui/src/colors"
	"github.com/orayew2002/gtui/src/dom"
)

func Image(class, path string, width, height int) Component {
	return imageComponent(class, path, width, height)
}

func ImageFrom(class string, src image.Image, width, height int) Component {
	return imageFromComponent(class, src, width, height)
}

func imageComponent(class, path string, width, height int, child ...Component) Component {
	file, err := os.Open(path)
	if err != nil {
		return imageFallback(class, path, width, height)
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return imageFallback(class, path, width, height)
	}

	return imageFromComponent(class, src, width, height, child...)
}

func imageFromComponent(class string, src image.Image, width, height int, child ...Component) Component {
	if src == nil || width < 1 || height < 1 {
		return Div(class, child...)
	}

	canvas := dom.NewCanvas(width, height)
	bounds := src.Bounds()

	for y := range height {
		for x := range width {
			top := sample(src, bounds, x, y*2, width, height*2)
			bottom := sample(src, bounds, x, y*2+1, width, height*2)

			canvas.Set(x, y, dom.Pixel{
				Ch:    '▀',
				Fg:    top,
				Bg:    bottom,
				HasFg: true,
				HasBg: true,
			})
		}
	}

	return canvasComponent(class, canvas, child...)
}

func imageFallback(class, path string, width, height int) Component {
	if width < 1 {
		width = 12
	}
	if height < 1 {
		height = 4
	}

	name := path
	if i := strings.LastIndexByte(name, '/'); i >= 0 {
		name = name[i+1:]
	}

	return Div(join("col border center w-"+itoa(width)+" h-"+itoa(height), class),
		Text("dim nowrap", "▣ "+name),
	)
}

func Video(class, title string, position, duration int, playing bool) Component {
	mark := "▶"
	if playing {
		mark = "⏸"
	}

	return Div(join("col rounded px-2 gap-0", class),
		Div("row gap-1",
			Text("bold nowrap", mark),
			Text("grow nowrap", title),
			Text("dim nowrap", clock(position)+" / "+clock(duration)),
		),
		Progress("text-sky nowrap", float64(position), float64(duration), 24),
		Div("row gap-2 dim",
			Text("nowrap", "space play"),
			Text("nowrap", "← → seek"),
			Text("grow text-right nowrap", "f fullscreen"),
		),
	)
}

func sample(src image.Image, bounds image.Rectangle, x, y, width, height int) colors.Color {
	sx := bounds.Min.X + x*bounds.Dx()/width
	sy := bounds.Min.Y + y*bounds.Dy()/height

	r, g, b, _ := src.At(sx, sy).RGBA()
	return colors.RGB(uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

func clock(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}

	minutes := seconds / 60
	rest := seconds % 60

	pad := ":"
	if rest < 10 {
		pad = ":0"
	}

	return itoa(minutes) + pad + itoa(rest)
}

type Clip struct {
	Frames []image.Image
	Delays []int
	Name   string
}

func OpenClip(path string) (*Clip, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	name := path
	if i := strings.LastIndexByte(name, '/'); i >= 0 {
		name = name[i+1:]
	}

	if strings.HasSuffix(strings.ToLower(path), ".gif") {
		decoded, err := gif.DecodeAll(file)
		if err != nil {
			return nil, err
		}

		return &Clip{Name: name, Delays: decoded.Delay, Frames: composite(decoded)}, nil
	}

	still, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return &Clip{Name: name, Frames: []image.Image{still}, Delays: []int{0}}, nil
}

// composite flattens a gif into whole frames: the format stores later frames
// as patches, so each one has to be drawn over the canvas that came before.
func composite(src *gif.GIF) []image.Image {
	if len(src.Image) == 0 {
		return nil
	}

	bounds := src.Image[0].Bounds()
	for _, frame := range src.Image {
		bounds = bounds.Union(frame.Bounds())
	}

	canvas := image.NewRGBA(bounds)
	out := make([]image.Image, 0, len(src.Image))

	for i, frame := range src.Image {
		previous := image.NewRGBA(bounds)
		copy(previous.Pix, canvas.Pix)

		draw.Draw(canvas, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)

		whole := image.NewRGBA(bounds)
		copy(whole.Pix, canvas.Pix)
		out = append(out, whole)

		if i < len(src.Disposal) {
			switch src.Disposal[i] {
			case gif.DisposalBackground:
				draw.Draw(canvas, frame.Bounds(), image.Transparent, image.Point{}, draw.Src)
			case gif.DisposalPrevious:
				copy(canvas.Pix, previous.Pix)
			}
		}
	}

	return out
}

func (c *Clip) Count() int {
	if c == nil {
		return 0
	}
	return len(c.Frames)
}

func (c *Clip) Delay(index int) time.Duration {
	if c == nil || len(c.Delays) == 0 {
		return 100 * time.Millisecond
	}

	hundredths := c.Delays[((index%len(c.Delays))+len(c.Delays))%len(c.Delays)]
	if hundredths <= 0 {
		return 100 * time.Millisecond
	}

	return time.Duration(hundredths) * 10 * time.Millisecond
}

func (c *Clip) Frame(class string, index, width, height int) Component {
	if c.Count() == 0 {
		return Div(class)
	}

	at := ((index % c.Count()) + c.Count()) % c.Count()
	return ImageFrom(class, c.Frames[at], width, height)
}

func Backdrop(class, path string, width, height int, child ...Component) Component {
	return imageComponent(class, path, width, height, child...)
}

func BackdropFrom(class string, src image.Image, width, height int, child ...Component) Component {
	return imageFromComponent(class, src, width, height, child...)
}

type ViewerOptions struct {
	Class  string
	Paths  []string
	Index  int
	Width  int
	Height int
	Zoom   int
}

func Viewer(opts ViewerOptions) Component {
	if len(opts.Paths) == 0 {
		return Div(opts.Class, Text("dim", "no images"))
	}

	if opts.Width < 4 {
		opts.Width = 24
	}
	if opts.Height < 2 {
		opts.Height = 10
	}
	if opts.Zoom < 1 {
		opts.Zoom = 1
	}

	at := ((opts.Index % len(opts.Paths)) + len(opts.Paths)) % len(opts.Paths)
	path := opts.Paths[at]

	name := path
	if i := strings.LastIndexByte(name, '/'); i >= 0 {
		name = name[i+1:]
	}

	return Div(join("col rounded border-line bg-panel px-1 gap-0", opts.Class),
		Div("row gap-1",
			Text("grow nowrap text-ink bold", name),
			Text("nowrap dim", "%d / %d", at+1, len(opts.Paths)),
		),
		Center("",
			Image("", path, opts.Width*opts.Zoom, opts.Height*opts.Zoom),
		),
		Div("row gap-2 dim",
			Text("nowrap", "← → image"),
			Text("nowrap", "+ − zoom"),
			Text("grow text-right nowrap", "x%d", opts.Zoom),
		),
	)
}
