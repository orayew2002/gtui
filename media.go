package gtui

import (
	"bufio"
	"image"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"strings"
	"sync"
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

// Fit draws the picture inside a width x height box without stretching it: the
// ratio stays, the spare cells stay empty.
func Fit(class, path string, width, height int) Component {
	src, err := loadImage(path)
	if err != nil {
		return imageFallback(class, path, width, height)
	}

	return FitFrom(class, src, width, height)
}

func FitFrom(class string, src image.Image, width, height int) Component {
	if src == nil || width < 1 || height < 1 {
		return Div(class)
	}

	w, h := containImage(src.Bounds(), width, height)

	return Div(join("center w-"+itoa(width)+" h-"+itoa(height), class),
		imageFromComponent("", src, w, h),
	)
}

func imageComponent(class, path string, width, height int, child ...Component) Component {
	src, err := loadImage(path)
	if err != nil {
		return imageFallback(class, path, width, height, child...)
	}

	return imageFromComponent(class, src, width, height, child...)
}

func imageFromComponent(class string, src image.Image, width, height int, child ...Component) Component {
	if src == nil {
		return Div(class, child...)
	}

	width, height = fitImage(src.Bounds(), width, height)
	if width < 1 || height < 1 {
		return Div(class, child...)
	}

	canvas := dom.NewCanvas(width, height)
	pick := newSampler(src, width, height*2)

	for y := range height {
		for x := range width {
			top, hasTop := pick.at(x, y*2)
			bottom, hasBottom := pick.at(x, y*2+1)

			switch {
			case hasTop && hasBottom:
				canvas.Set(x, y, dom.Pixel{Ch: '▀', Fg: top, Bg: bottom, HasFg: true, HasBg: true})
			case hasTop:
				canvas.Set(x, y, dom.Pixel{Ch: '▀', Fg: top, HasFg: true})
			case hasBottom:
				canvas.Set(x, y, dom.Pixel{Ch: '▄', Fg: bottom, HasFg: true})
			}
		}
	}

	return canvasComponent(class, canvas, child...)
}

func imageFallback(class, path string, width, height int, child ...Component) Component {
	if width < 1 {
		width = 12
	}
	if height < 1 {
		height = 4
	}

	size := "w-" + itoa(width) + " h-" + itoa(height)

	// a Backdrop still has to show its content when the picture is gone
	if len(child) > 0 {
		return Div(join("col bg-panel "+size, class), child...)
	}

	name := path
	if i := strings.LastIndexByte(name, '/'); i >= 0 {
		name = name[i+1:]
	}

	return Div(join("col border center "+size, class),
		Text("dim nowrap", "▣ "+name),
	)
}

// fitImage fills in a missing dimension from the picture's ratio, counting two
// source rows per cell.
func fitImage(b image.Rectangle, width, height int) (int, int) {
	dx, dy := b.Dx(), b.Dy()
	if dx < 1 || dy < 1 {
		return width, height
	}

	switch {
	case width >= 1 && height < 1:
		height = max(1, (width*dy+dx)/(2*dx))
	case height >= 1 && width < 1:
		width = max(1, (2*height*dx+dy/2)/dy)
	}

	return width, height
}

// containImage is fitImage for both dimensions at once: the result never spills
// out of width x height.
func containImage(b image.Rectangle, width, height int) (int, int) {
	dx, dy := b.Dx(), b.Dy()
	if dx < 1 || dy < 1 {
		return width, height
	}

	w, h := width, max(1, (width*dy+dx)/(2*dx))
	if h > height {
		h = height
		w = min(width, max(1, (2*height*dx+dy/2)/dy))
	}

	return w, h
}

type sampler struct {
	src    image.Image
	fast   image.RGBA64Image
	bounds image.Rectangle
	width  int
	height int
}

func newSampler(src image.Image, width, height int) sampler {
	fast, _ := src.(image.RGBA64Image)
	return sampler{src: src, fast: fast, bounds: src.Bounds(), width: width, height: height}
}

// at reads one source pixel as a straight colour, and reports false for a
// mostly transparent one so the cell keeps the background behind it.
func (s sampler) at(x, y int) (colors.Color, bool) {
	sx := s.bounds.Min.X + x*s.bounds.Dx()/s.width
	sy := s.bounds.Min.Y + y*s.bounds.Dy()/s.height

	var r, g, b, a uint32
	if s.fast != nil {
		c := s.fast.RGBA64At(sx, sy)
		r, g, b, a = uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
	} else {
		c := s.src.At(sx, sy)
		if c == nil {
			return colors.Color{}, false
		}
		r, g, b, a = c.RGBA()
	}

	if a < 0x8000 {
		return colors.Color{}, false
	}

	return colors.RGB(straight(r, a), straight(g, a), straight(b, a)), true
}

func straight(v, a uint32) uint8 {
	if a == 0 {
		return 0
	}
	if a < 0xffff {
		v = v * 0xffff / a
	}
	if v > 0xffff {
		v = 0xffff
	}

	return uint8(v >> 8)
}

const (
	imageCacheMax   = 16
	imageCacheBytes = 64 << 20
)

type imageEntry struct {
	src  image.Image
	size int64
	mod  time.Time
}

var (
	imageMu    sync.Mutex
	imageCache = map[string]imageEntry{}
	imageOrder []string
)

// loadImage decodes path, or hands back the decode a previous call made: a view
// rebuild and an animation tick ask for the same picture many times a second.
func loadImage(path string) (image.Image, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if src := cachedImage(path, info); src != nil {
		return src, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	storeImage(path, info, src)

	return src, nil
}

// ClearImageCache forgets every decoded picture, for an app that rewrites files
// without changing their size or timestamp.
func ClearImageCache() {
	imageMu.Lock()
	imageCache = map[string]imageEntry{}
	imageOrder = imageOrder[:0]
	imageMu.Unlock()
}

func cachedImage(path string, info os.FileInfo) image.Image {
	imageMu.Lock()
	defer imageMu.Unlock()

	entry, ok := imageCache[path]
	if !ok || entry.size != info.Size() || !entry.mod.Equal(info.ModTime()) {
		return nil
	}

	touchImage(path)

	return entry.src
}

func storeImage(path string, info os.FileInfo, src image.Image) {
	imageMu.Lock()
	defer imageMu.Unlock()

	imageCache[path] = imageEntry{src: src, size: info.Size(), mod: info.ModTime()}
	touchImage(path)
	trimImages()
}

func touchImage(path string) {
	for i, p := range imageOrder {
		if p == path {
			imageOrder = append(imageOrder[:i], imageOrder[i+1:]...)
			break
		}
	}

	imageOrder = append(imageOrder, path)
}

func trimImages() {
	total := 0
	for _, p := range imageOrder {
		total += imageBytes(imageCache[p].src)
	}

	for len(imageOrder) > imageCacheMax || (total > imageCacheBytes && len(imageOrder) > 1) {
		oldest := imageOrder[0]
		total -= imageBytes(imageCache[oldest].src)

		delete(imageCache, oldest)
		imageOrder = append(imageOrder[:0], imageOrder[1:]...)
	}
}

func imageBytes(src image.Image) int {
	if src == nil {
		return 0
	}

	b := src.Bounds()

	return b.Dx() * b.Dy() * 4
}

type VideoOptions struct {
	Class        string
	Title        string
	Position     int
	Duration     int
	Playing      bool
	Step         int
	OnToggle     func()
	OnSeek       func(delta int)
	OnFullscreen func()
}

func Video(class, title string, position, duration int, playing bool) Component {
	return VideoPlayer(VideoOptions{
		Class:    class,
		Title:    title,
		Position: position,
		Duration: duration,
		Playing:  playing,
	})
}

// VideoPlayer is Video with keys: it only prints the hints it can honour, and
// takes space, the arrows and f when the matching handler is there.
func VideoPlayer(opts VideoOptions) Component {
	if opts.Step < 1 {
		opts.Step = 5
	}

	mark := "▶"
	if opts.Playing {
		mark = "⏸"
	}

	chrome := []Component{
		Div("row gap-1",
			Text("bold nowrap", mark),
			Text("grow nowrap", opts.Title),
			Text("dim nowrap", clock(opts.Position)+" / "+clock(opts.Duration)),
		),
		Progress("text-accent nowrap", float64(opts.Position), float64(opts.Duration), 24),
	}

	if hints := videoHints(opts); len(hints) > 0 {
		chrome = append(chrome, Div("row gap-2 dim", hints...))
	}

	base := join("col rounded px-2 gap-0", opts.Class)
	if opts.OnToggle == nil && opts.OnSeek == nil && opts.OnFullscreen == nil {
		return Div(base, chrome...)
	}

	return Component{
		class: opts.Class,
		build: func(a *App) *Node {
			n := mediaNode(withClass(base, "focusable"))
			for _, item := range chrome {
				n.Append(item.node(a))
			}

			n.On("keydown", func(e *Event) {
				switch {
				case e.Key == " " && opts.OnToggle != nil:
					e.StopPropagation()
					opts.OnToggle()
				case e.Key == "left" && opts.OnSeek != nil:
					e.StopPropagation()
					opts.OnSeek(-opts.Step)
				case e.Key == "right" && opts.OnSeek != nil:
					e.StopPropagation()
					opts.OnSeek(opts.Step)
				case e.Key == "f" && opts.OnFullscreen != nil:
					e.StopPropagation()
					opts.OnFullscreen()
				}
			})

			return n
		},
	}
}

func videoHints(opts VideoOptions) []Component {
	hints := make([]Component, 0, 3)

	if opts.OnToggle != nil {
		hints = append(hints, Text("nowrap", "space play"))
	}
	if opts.OnSeek != nil {
		hints = append(hints, Text("nowrap", "← → seek"))
	}
	if opts.OnFullscreen != nil {
		hints = append(hints, Text("grow text-right nowrap", "f fullscreen"))
	}

	return hints
}

func mediaNode(class string) *Node {
	n := dom.NewElement("div")

	base, id := splitClass(class)
	if base != "" {
		n.Set("class", base)
	}
	if id != "" {
		n.Set("id", id)
	}

	return n
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
	return openClip(path, 0, 0)
}

// OpenClipSized decodes the clip straight into cols x rows cells, so a long gif
// keeps kilobytes per frame instead of a whole picture.
func OpenClipSized(path string, cols, rows int) (*Clip, error) {
	return openClip(path, cols, rows)
}

func openClip(path string, cols, rows int) (*Clip, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	name := path
	if i := strings.LastIndexByte(name, '/'); i >= 0 {
		name = name[i+1:]
	}

	reader := bufio.NewReader(file)

	// the name says nothing about the content: a gif keeps its magic up front
	if head, _ := reader.Peek(6); string(head) == "GIF87a" || string(head) == "GIF89a" {
		decoded, err := gif.DecodeAll(reader)
		if err != nil {
			return nil, err
		}

		return &Clip{Name: name, Delays: decoded.Delay, Frames: composite(decoded, cols, rows)}, nil
	}

	still, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return &Clip{Name: name, Frames: []image.Image{still}, Delays: []int{0}}, nil
}

// composite flattens a gif into whole frames: the format stores later frames
// as patches, so each one has to be drawn over the canvas that came before.
func composite(src *gif.GIF, cols, rows int) []image.Image {
	if len(src.Image) == 0 {
		return nil
	}

	bounds := src.Image[0].Bounds()
	for _, frame := range src.Image {
		bounds = bounds.Union(frame.Bounds())
	}

	canvas := image.NewRGBA(bounds)
	out := make([]image.Image, 0, len(src.Image))

	var previous *image.RGBA

	for i, frame := range src.Image {
		if i < len(src.Disposal) && src.Disposal[i] == gif.DisposalPrevious {
			if previous == nil {
				previous = image.NewRGBA(bounds)
			}
			copy(previous.Pix, canvas.Pix)
		}

		draw.Draw(canvas, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)
		out = append(out, clipFrame(canvas, cols, rows))

		if i < len(src.Disposal) {
			switch src.Disposal[i] {
			case gif.DisposalBackground:
				draw.Draw(canvas, frame.Bounds(), image.Transparent, image.Point{}, draw.Src)
			case gif.DisposalPrevious:
				if previous != nil {
					copy(canvas.Pix, previous.Pix)
				}
			}
		}
	}

	return out
}

func clipFrame(canvas *image.RGBA, cols, rows int) image.Image {
	bounds := canvas.Bounds()

	if cols < 1 || rows < 1 || bounds.Dx() < 1 || bounds.Dy() < 1 {
		whole := image.NewRGBA(bounds)
		copy(whole.Pix, canvas.Pix)

		return whole
	}

	small := image.NewRGBA(image.Rect(0, 0, cols, rows*2))
	for y := range rows * 2 {
		sy := bounds.Min.Y + y*bounds.Dy()/(rows*2)
		for x := range cols {
			small.SetRGBA64(x, y, canvas.RGBA64At(bounds.Min.X+x*bounds.Dx()/cols, sy))
		}
	}

	return small
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
	Class    string
	Paths    []string
	Index    int
	Width    int
	Height   int
	Zoom     int
	OnChange func(index, zoom int)
}

const viewerMaxZoom = 8

// Viewer owns its keys: the arrows walk the list, + and − zoom, and OnChange
// reports both back to the app.
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
	opts.Zoom = clampInt(opts.Zoom, 1, viewerMaxZoom)

	start := wrapIndex(opts.Index, len(opts.Paths))

	return Component{
		class: opts.Class,
		build: func(a *App) *Node {
			n := mediaNode(withClass(join("col rounded border-line bg-panel px-1 gap-0", opts.Class), "focusable"))
			n.Set("index", itoa(start))
			n.Set("zoom", itoa(opts.Zoom))

			draw := func() {
				index, _ := strconv.Atoi(n.Get("index"))
				zoom, _ := strconv.Atoi(n.Get("zoom"))

				n.Clear()
				for _, item := range viewerChrome(opts, index, zoom) {
					n.Append(item.node(a))
				}
			}
			draw()

			move := func(index, zoom int) {
				index = wrapIndex(index, len(opts.Paths))
				zoom = clampInt(zoom, 1, viewerMaxZoom)

				n.Set("index", itoa(index))
				n.Set("zoom", itoa(zoom))
				draw()

				n.Dispatch(&dom.Event{Type: "change", Target: n, Value: itoa(index)})
				if opts.OnChange != nil {
					opts.OnChange(index, zoom)
				}
				if a != nil {
					a.Redraw()
				}
			}

			n.On("keydown", func(e *Event) {
				index, _ := strconv.Atoi(n.Get("index"))
				zoom, _ := strconv.Atoi(n.Get("zoom"))

				switch e.Key {
				case "left":
					e.StopPropagation()
					move(index-1, zoom)
				case "right":
					e.StopPropagation()
					move(index+1, zoom)
				case "+", "=":
					e.StopPropagation()
					move(index, zoom+1)
				case "-", "_":
					e.StopPropagation()
					move(index, zoom-1)
				}
			})

			return n
		},
	}
}

func viewerChrome(opts ViewerOptions, index, zoom int) []Component {
	path := opts.Paths[wrapIndex(index, len(opts.Paths))]

	name := path
	if i := strings.LastIndexByte(name, '/'); i >= 0 {
		name = name[i+1:]
	}

	return []Component{
		Div("row gap-1",
			Text("grow nowrap text-ink bold", name),
			Text("nowrap dim", "%d / %d", wrapIndex(index, len(opts.Paths))+1, len(opts.Paths)),
		),
		Fit("", path, opts.Width*zoom, opts.Height*zoom),
		Div("row gap-2 dim",
			Text("nowrap", "← → image"),
			Text("nowrap", "+ − zoom"),
			Text("grow text-right nowrap", "x%d", zoom),
		),
	}
}

func wrapIndex(index, count int) int {
	if count < 1 {
		return 0
	}
	return ((index % count) + count) % count
}
