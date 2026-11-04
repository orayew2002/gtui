// A media studio: a picture behind the content, a viewer, and the two players.
// The pictures are painted in code, so this example ships no binary files.
package main

import (
	"image"
	"image/color"
	"math"

	"github.com/orayew2002/gtui"
)

var (
	app = gtui.Init()

	shots   = []image.Image{picture(0), picture(1), picture(2)}
	shot    = 0
	zoom    = 1
	at      = 96
	playing = true
	tick    = 0
)

func main() {
	app.View(view)
	app.QuitOn("q", "ctrl+c")

	app.Bind("space", func() { playing = !playing; app.Refresh() })
	app.Animate(8, func(step int) {
		tick = step
		if playing {
			at = (at + 1) % 268
		}
		app.Refresh()
	})

	app.Render()
}

func view() gtui.Component {
	return gtui.Column("full bg-canvas",
		// a picture fills the bar, and a scrim over it keeps the text readable
		gtui.BackdropFrom("h-3 shrink-0 border-b border-line", shots[shot], 0, 0,
			gtui.Row("px-1 py-1 items-center gap-2",
				gtui.IconText("bg-[#000000]/55 text-canvas bold nowrap", "image", " studio "),
				gtui.Text("bg-[#000000]/55 text-canvas nowrap", " space plays · q quits "),
			),
		),
		gtui.Row("grow clip px-1 py-1 gap-1",
			gtui.Column("grow clip gap-1",
				viewer(),
				video(),
			),
			gtui.Column("w-34 shrink-0 clip gap-1",
				music(),
				sizes(),
			),
		),
	)
}

func viewer() gtui.Component {
	return gtui.Card("grow clip rounded border-line bg-panel px-1", "viewer",
		gtui.Viewer(gtui.ViewerOptions{
			Class:  "#view",
			Images: shots,
			Index:  shot,
			Zoom:   zoom,
			Width:  30,
			Height: 8,
		}),
		gtui.Row("gap-2 items-center",
			gtui.FitFrom("rounded border-line", shots[shot], 12, 4),
			gtui.Column("gap-0",
				gtui.Text("text-muted dim nowrap", "arrows walk the pictures"),
				gtui.Row("gap-1",
					gtui.Button("px-1 text-muted focus:text-ink", "‹", func(e *gtui.Event) {
						shot = (shot + len(shots) - 1) % len(shots)
						app.Refresh()
					}),
					gtui.Button("px-1 text-muted focus:text-ink", "›", func(e *gtui.Event) {
						shot = (shot + 1) % len(shots)
						app.Refresh()
					}),
					gtui.Button("px-1 text-muted focus:text-ink", "+", func(e *gtui.Event) {
						zoom = min(zoom+1, 4)
						app.Refresh()
					}),
					gtui.Button("px-1 text-muted focus:text-ink", "−", func(e *gtui.Event) {
						zoom = max(zoom-1, 1)
						app.Refresh()
					}),
				),
			),
		),
	)
}

func video() gtui.Component {
	return gtui.Card("shrink-0 rounded border-line bg-panel px-1", "video",
		gtui.VideoPlayer(gtui.VideoOptions{
			Title:    "the long way round",
			Position: at,
			Duration: 268,
			Playing:  playing,
			Step:     5,
			OnToggle: func() { playing = !playing; app.Refresh() },
			OnSeek: func(delta int) {
				at = clamp(at+delta, 0, 268)
				app.Refresh()
			},
			OnFullscreen: func() { app.Toast("info", "as big as it gets") },
		}),
	)
}

func music() gtui.Component {
	return gtui.Card("shrink-0 rounded border-line bg-panel px-1", "music",
		gtui.Player(gtui.PlayerOptions{
			Title:    "Nocturne in E flat",
			Subtitle: "piano · 320 kbps",
			Position: at,
			Duration: 268,
			Playing:  playing,
			Volume:   55,
			Width:    12,
			Queue:    []string{"Gymnopédie no 1", "Clair de lune"},
		}),
		gtui.Row("gap-1 items-center pt-1",
			gtui.Text("text-muted nowrap", "level"),
			gtui.Equalizer("text-success", 10, tick),
		),
		gtui.SeekBar("pt-1", at, 268, 20),
	)
}

// sizes shows the two ways a picture can meet a box: stretched to it, or fitted
// inside it with its shape kept.
func sizes() gtui.Component {
	return gtui.Card("grow clip rounded border-line bg-panel px-1", "stretch and fit",
		gtui.Row("gap-1 items-center",
			gtui.ImageFrom("rounded border-line", shots[shot], 9, 4),
			gtui.FitFrom("rounded border-line", shots[shot], 9, 4),
		),
		gtui.Text("text-muted dim", "ImageFrom fills the box, FitFrom keeps the shape"),
		gtui.Video("pt-1", "trailer.mp4", at, 268, playing),
	)
}

// picture paints a sample image, twice as wide as it is tall, which is the
// shape a terminal draws best.
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

func clamp(v, lo, hi int) int { return min(max(v, lo), hi) }
