package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/orayew2002/gtui"
)

var (
	shot    = 0
	zoom    = 1
	playing = true
	beat    = 0
	track   = 0
	tuneAt  = 96

	gallery = []image.Image{picture(0), picture(1), picture(2)}
	reel    = clip()
	queue   = []string{"Nocturne in E flat", "Gymnopédie no 1", "Clair de lune"}
)

// picture paints a sample image so the gallery ships no binary files. It is
// twice as wide as it is tall, the shape a terminal draws best.
func picture(seed int) image.Image {
	const size = 96
	img := image.NewRGBA(image.Rect(0, 0, size*2, size))

	for y := range size {
		for x := range size * 2 {
			wave := math.Sin(float64(x+seed*20)/9) * math.Cos(float64(y)/11)
			img.Set(x, y, color.RGBA{
				R: uint8(90 + 70*wave),
				G: uint8(70 + 120*float64(y)/size),
				B: uint8(150 + 80*math.Sin(float64(y+seed*25)/13)),
				A: 255,
			})
		}
	}

	return img
}

// clip builds a short animation in memory. OpenClip("file.gif") does the same
// from a real file.
func clip() *gtui.Clip {
	frames := make([]image.Image, 0, 8)

	for step := range 8 {
		img := image.NewRGBA(image.Rect(0, 0, 64, 32))
		for y := range 32 {
			for x := range 64 {
				v := math.Sin(float64(x)/6 + float64(step)/1.2)
				img.Set(x, y, color.RGBA{
					R: uint8(40 + 60*v),
					G: uint8(120 + 90*math.Sin(float64(y)/5+float64(step)/2)),
					B: uint8(200 - 40*v),
					A: 255,
				})
			}
		}
		frames = append(frames, img)
	}

	return &gtui.Clip{Name: "sample.gif", Frames: frames, Delays: []int{8}}
}

var samples struct {
	once  sync.Once
	paths []string
	reel  string
}

// samplePaths writes the sample pictures once, because Viewer, Backdrop and Fit
// read files: pointing them at a photo folder is the real use.
func samplePaths() []string {
	samples.once.Do(func() {
		dir, err := os.MkdirTemp("", "gtui-gallery")
		if err != nil {
			return
		}

		for i, img := range gallery {
			path := filepath.Join(dir, "picture-"+strconv.Itoa(i+1)+".png")

			file, err := os.Create(path)
			if err != nil {
				return
			}

			err = png.Encode(file, img)
			file.Close()

			if err != nil {
				return
			}
			samples.paths = append(samples.paths, path)
		}

		samples.reel = writeReel(dir)
	})

	return samples.paths
}

// writeReel saves the in-memory clip as a real gif, so the gallery can show
// OpenClip reading one back without shipping a binary file.
func writeReel(dir string) string {
	path := filepath.Join(dir, "sample.gif")

	file, err := os.Create(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	frames := make([]*image.Paletted, 0, len(reel.Frames))
	delays := make([]int, 0, len(reel.Frames))

	for _, frame := range reel.Frames {
		bounds := frame.Bounds()
		shot := image.NewPaletted(bounds, palette.Plan9)
		draw.Draw(shot, bounds, frame, bounds.Min, draw.Src)

		frames = append(frames, shot)
		delays = append(delays, 8)
	}

	if err := gif.EncodeAll(file, &gif.GIF{Image: frames, Delay: delays}); err != nil {
		return ""
	}

	return path
}

// openedReel loads the gif from disk once, the way an app with real media does.
func openedReel() *gtui.Clip {
	loaded.once.Do(func() {
		samplePaths()
		if samples.reel == "" {
			return
		}
		loaded.clip, _ = gtui.OpenClip(samples.reel)
	})

	return loaded.clip
}

var loaded struct {
	once sync.Once
	clip *gtui.Clip
}

func pageMedia() gtui.Component {
	// on a small terminal the three players share one row, which leaves each
	// of them enough rows for a picture worth looking at. The thirds are
	// explicit, so a long title cannot take a picture's cells away.
	if !wideScreen() {
		third := fmt.Sprintf("w-%d col clip shrink-0", mediaThird())

		return gtui.Row("gap-1 grow",
			gtui.Div(third, pictureCard()),
			gtui.Div(third, videoCard()),
			gtui.Div(third, musicCard()),
		)
	}

	return gtui.Column("grow gap-1",
		gtui.Row("gap-1", pictureCard(), videoCard()),
		gtui.Row("gap-1", musicCard()),
		barsCard(),
	)
}

// barsCard shows the two plain pieces the players are built from, so they can
// be dropped into a layout of your own.
func barsCard() gtui.Component {
	width, _ := room()

	return cards(
		card("j", "a picture as the background",
			note("width and height of zero means fill the box, like a wallpaper"),
			gtui.BackdropFrom("h-6 rounded border-line", gallery[shot], 0, 0,
				gtui.Column("gap-0 px-1",
					gtui.Text("bg-[#000000]/55 text-canvas bold nowrap", " night sky "),
					gtui.Text("bg-[#000000]/55 text-canvas nowrap", " a scrim keeps it readable "),
				),
			),
		),
		card("g", "plain transport pieces",
			note("Video is the one line view, SeekBar the bar inside a player"),
			gtui.Video("", queue[track], tuneAt, 268, playing),
			gtui.SeekBar("mt-1", tuneAt, 268, min(max(width-24, 10), 40)),
		),
		twinsCard(),
	)
}

// twinsCard puts the path and the image.Image versions of the same component
// side by side: Image stretches to the box, Fit keeps the aspect ratio.
func twinsCard() gtui.Component {
	paths := samplePaths()
	if len(paths) == 0 {
		return card("h", "pictures from memory",
			note("ImageFrom stretches, FitFrom keeps the shape"),
			gtui.Row("gap-1",
				gtui.ImageFrom("rounded border-line", gallery[shot], 10, 4),
				gtui.FitFrom("rounded border-line", gallery[shot], 10, 4),
			),
		)
	}

	at := shot % len(paths)

	body := []gtui.Component{
		note("Image and Fit read a file, ImageFrom and FitFrom take an image.Image"),
		gtui.Row("gap-1 items-center",
			gtui.Image("rounded border-line", paths[at], 9, 4),
			gtui.Fit("rounded border-line", paths[at], 9, 4),
			gtui.ImageFrom("rounded border-line", gallery[shot], 9, 4),
		),
	}

	tail := []gtui.Component{
		gtui.Backdrop("rounded border-line", paths[at], 9, 4,
			gtui.Text("px-1 bold text-canvas nowrap", "behind"),
		),
	}

	if disk := openedReel(); disk != nil {
		tail = append(tail, disk.Frame("rounded border-line", beat, 9, 4))
		body = append(body, note("Backdrop and a gif OpenClip read back, %d frames", disk.Count()))
	} else {
		body = append(body, note("Backdrop puts a picture behind ordinary content"))
	}

	body = append(body, gtui.Row("gap-1 items-center", tail...))

	return card("h", "from a path, or from memory", body...)
}

func pictureCard() gtui.Component {
	width, height := pictureSize()
	paths := samplePaths()

	if len(paths) == 0 {
		return card("a", "picture viewer",
			note("no temp directory, so this one comes from memory"),
			gtui.FitFrom("rounded border-line", gallery[shot], width, height),
		)
	}

	body := []gtui.Component{
		gtui.Viewer(gtui.ViewerOptions{
			Class:  "#mv-viewer",
			Paths:  paths,
			Index:  shot,
			Zoom:   zoom,
			Width:  width,
			Height: height,
			OnChange: func(index, size int) {
				shot, zoom = index, size
				say("mv-out", "picture %d of %d at x%d", index+1, len(paths), size)
			},
		}),
	}

	if wideScreen() {
		body = append(body, out("mv-out", "Viewer draws every picture with Fit, so the ratio survives"))
	}

	return card("a", "picture viewer", body...)
}

func videoCard() gtui.Component {
	body := []gtui.Component{gtui.Column("gap-0 #mv-frame", videoFrame()...)}
	if wideScreen() {
		body = append(body, note("key b puts the keyboard on the player"))
	}

	return card("b", "video player", body...)
}

// videoFrame is the moving half of the video card: the tick redraws this box
// alone, and the player keeps the keyboard because it carries an id.
func videoFrame() []gtui.Component {
	width, height := pictureSize()

	frame := "rounded border-line"
	chrome := "#mv-player"
	if !wideScreen() {
		frame, chrome = "", "rounded-none #mv-player"
	}

	return []gtui.Component{
		gtui.BackdropFrom(frame, reel.Frames[beat%len(reel.Frames)], width, height,
			gtui.Column("px-1",
				gtui.Text("bold text-canvas nowrap", reel.Name),
				gtui.Text("text-canvas nowrap", "frame %d of %d", beat+1, reel.Count()),
			),
		),
		gtui.VideoPlayer(gtui.VideoOptions{
			Class:    chrome,
			Title:    "Backdrop and VideoPlayer",
			Position: beat,
			Duration: reel.Count(),
			Playing:  playing,
			Step:     1,
			OnToggle: func() { playing = !playing },
			OnSeek: func(delta int) {
				beat = (beat + reel.Count() + delta) % reel.Count()
			},
		}),
	}
}

func musicCard() gtui.Component {
	width, _ := room()

	if !wideScreen() {
		return card("d", "music player",
			gtui.Row("gap-1 items-center",
				gtui.Icon("text-accent", "music"),
				gtui.Text("grow truncate text-ink", queue[track]),
				gtui.Text("shrink-0 text-muted nowrap", "%d:%02d", tuneAt/60, tuneAt%60),
			),
			gtui.Row("gap-1 items-center #mp-meter", musicMeter()...),
		)
	}

	return card("d", "music player",
		gtui.Row("gap-1",
			gtui.BackdropFrom("rounded border-line shrink-0", gallery[track%len(gallery)], 12, coverRows(),
				gtui.Text("px-1 bold text-canvas nowrap", "cover"),
			),
			gtui.Column("grow gap-0",
				gtui.Player(gtui.PlayerOptions{
					Title:    queue[track],
					Subtitle: "piano · 320 kbps",
					Position: tuneAt,
					Duration: 268,
					Playing:  playing,
					Volume:   55,
					Width:    width - 40,
					Queue:    tracks(),
				}),
				gtui.Row("gap-1 items-center #mp-meter", musicMeter()...),
			),
		),
	)
}

// musicMeter is the transport and the level meter: the tick redraws it, so the
// meter moves while the buttons keep working.
func musicMeter() []gtui.Component {
	width, _ := room()

	return []gtui.Component{
		gtui.Button("px-1 focus:bold focus:text-accent #mp-prev", "⏮", func(e *gtui.Event) {
			track = (track + len(queue) - 1) % len(queue)
			tuneAt = 0
			show(current())
		}),
		gtui.Button("px-1 focus:bold focus:text-accent #mp-play", playMark(), func(e *gtui.Event) {
			playing = !playing
			show(current())
		}),
		gtui.Button("px-1 focus:bold focus:text-accent #mp-next", "⏭", func(e *gtui.Event) {
			track = (track + 1) % len(queue)
			tuneAt = 0
			show(current())
		}),
		gtui.Button("px-1 focus:bold focus:text-accent #mp-seek", "»", func(e *gtui.Event) {
			tuneAt = (tuneAt + 15) % 268
			show(current())
		}),
		gtui.Equalizer("text-accent", max(width/5, 8), beat),
	}
}

// tracks rotates the queue so the one playing is first, and leaves it out
// when the card has no rows for it.
func tracks() []string {
	if !musicQueue() {
		return nil
	}

	list := make([]string, 0, len(queue))
	for i := range queue {
		list = append(list, queue[(track+i)%len(queue)])
	}
	return list
}

// pictureSize keeps the pictures inside the cards, whatever the terminal is.
func pictureSize() (int, int) {
	width, height := room()

	if !wideScreen() {
		// the card spends six cells on its frame, the viewer four more
		return mediaThird() - 10, min(max((height-7)/2, 3), 6)
	}

	// the picture takes the rows the video chrome, its note and the music card
	// leave over
	spare := height - 22
	if musicQueue() {
		spare = height - 26
	}

	return min(max(width/2-10, 12), 36), min(max(spare, 2), 6)
}

// musicQueue reports whether there are rows to spare for the track list.
func musicQueue() bool {
	_, height := room()
	return height >= 30
}

// coverRows keeps the album art the same height as the player beside it.
func coverRows() int {
	if musicQueue() {
		return 8
	}
	return 4
}

func mediaThird() int {
	width, _ := room()
	return max((width-2)/3, 18)
}

func playMark() string {
	if playing {
		return "⏸"
	}
	return "▶"
}
