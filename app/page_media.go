package main

import (
	"image"
	"image/color"
	"math"

	"github.com/orayew2002/gtui"
)

var (
	shot    = 0
	zoom    = 1
	playing = true
	beat    = 0

	gallery = []image.Image{picture(0), picture(1), picture(2)}
	reel    = clip()
)

// picture paints a sample image so the gallery needs no binary files.
func picture(seed int) image.Image {
	const size = 96
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	for y := range size {
		for x := range size {
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

func pageMedia() gtui.Component {
	width, rows := pictureSize()

	return gtui.Column("grow gap-1",
		gtui.Row("gap-1",
			card("a", "picture viewer",
				note("‹ › change picture, zoom doubles the size"),
				gtui.BackdropFrom("rounded border-line", gallery[shot], width*zoom, rows*zoom,
					gtui.Column("px-1",
						gtui.Text("bold text-canvas nowrap", "picture %d", shot+1),
						gtui.Text("text-canvas nowrap", "%d / %d · zoom x%d", shot+1, len(gallery), zoom),
					),
				),
				gtui.Row("gap-1",
					gtui.Button("rounded border-line focus:bold #prev", "‹", func(e *gtui.Event) {
						shot = (shot + len(gallery) - 1) % len(gallery)
						show(current())
					}),
					gtui.Button("rounded border-line focus:bold #next", "›", func(e *gtui.Event) {
						shot = (shot + 1) % len(gallery)
						show(current())
					}),
					gtui.Button("rounded border-accent text-accent focus:bold #zoom", "zoom", func(e *gtui.Event) {
						zoom = zoom%2 + 1
						show(current())
					}),
				),
			),
			card("b", "video player",
				note("space plays, the frame animates while this page is open"),
				reel.Frame("rounded border-line", beat, width, rows),
				gtui.Player(gtui.PlayerOptions{
					Class:    "key-hide",
					Title:    reel.Name,
					Subtitle: "generated in memory",
					Position: beat,
					Duration: reel.Count(),
					Playing:  playing,
					Volume:   70,
					Width:    width - 8,
				}),
			),
		),
		gtui.Row("gap-1",
			card("d", "music player",
				note("same Player with a queue and a level meter"),
				gtui.Player(gtui.PlayerOptions{
					Class:    "key-hide",
					Title:    "Nocturne in E flat",
					Subtitle: "piano · 320 kbps",
					Position: 96,
					Duration: 268,
					Playing:  playing,
					Volume:   55,
					Width:    width - 8,
					Queue:    []string{"Nocturne in E flat", "Gymnopédie no 1", "Clair de lune"},
				}),
				gtui.Equalizer("text-accent", 18, beat),
			),
			card("e", "controls",
				gtui.Row("gap-1",
					gtui.Button("rounded border-accent text-accent focus:bold #play", playMark(), func(e *gtui.Event) {
						playing = !playing
						show(current())
					}),
					gtui.Button("rounded border-line focus:bold #step", "step", func(e *gtui.Event) {
						beat = (beat + 1) % reel.Count()
						show(current())
					}),
				),
				gtui.IconText("text-muted", "music", "Icon(\"music\") and friends"),
				gtui.IconText("text-muted", "video", "Icon(\"video\")"),
				gtui.Hint("", "space", "play or pause"),
			),
		),
	)
}

// pictureSize keeps the pictures inside the terminal, whatever its width.
func pictureSize() (int, int) {
	width, height := app.Size()

	picture := (width - 46) / 2
	if picture < 12 {
		picture = 12
	}
	if picture > 40 {
		picture = 40
	}

	rows := (height - 16) / 2
	if rows < 4 {
		rows = 4
	}
	if rows > 9 {
		rows = 9
	}

	return picture, rows
}

func playMark() string {
	if playing {
		return "⏸"
	}
	return "▶"
}
