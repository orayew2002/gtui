// A media studio: a studioPicture behind the content, a studioViewer, and the two players.
// The pictures are painted in code, so this example ships no binary files.
package demos

import (
	"image"
	"image/color"
	"math"

	"github.com/orayew2002/gtui"
)

var (
	studioShots   = []image.Image{studioPicture(0), studioPicture(1), studioPicture(2)}
	studioShot    = 0
	studioZoom    = 1
	studioAt      = 96
	studioPlaying = true
	studioTick    = 0
)

func Studio(a *gtui.App, tick int) gtui.Component {
	host = a
	studioTick = tick

	return gtui.Column("full",
		// a studioPicture fills the bar, and a scrim over it keeps the text readable
		gtui.BackdropFrom("h-3 shrink-0 border-b border-line", studioShots[studioShot], 0, 0,
			// one scrim across the whole bar, so the picture darkens evenly
			// instead of only behind the words
			gtui.Row("full px-1 items-center gap-2 bg-[#000000]/55",
				gtui.IconText("text-canvas bold nowrap", "image", "studio"),
				gtui.Text("text-canvas nowrap truncate", "the buttons play, seek and zoom"),
			),
		),
		gtui.Row("grow clip px-1 py-1 gap-1",
			gtui.Column("grow clip gap-1",
				studioViewer(),
				studioVideo(),
			),
			gtui.Column("w-34 shrink-0 clip gap-1",
				studioMusic(),
				studioSizes(),
			),
		),
	)
}

func studioViewer() gtui.Component {
	return gtui.Card("grow clip rounded border-line bg-panel px-1", "viewer",
		gtui.Viewer(gtui.ViewerOptions{
			Class:  "#view",
			Images: studioShots,
			Index:  studioShot,
			Zoom:   studioZoom,
			Width:  30,
			Height: 8,
		}),
		gtui.Row("gap-2 items-center",
			gtui.FitFrom("rounded border-line", studioShots[studioShot], 12, 4),
			gtui.Column("gap-0",
				gtui.Text("text-muted dim nowrap", "arrows walk the pictures"),
				gtui.Row("gap-1",
					gtui.Button("px-1 text-muted focus:text-ink", "‹", func(e *gtui.Event) {
						studioShot = (studioShot + len(studioShots) - 1) % len(studioShots)
						host.Refresh()
					}),
					gtui.Button("px-1 text-muted focus:text-ink", "›", func(e *gtui.Event) {
						studioShot = (studioShot + 1) % len(studioShots)
						host.Refresh()
					}),
					gtui.Button("px-1 text-muted focus:text-ink", "+", func(e *gtui.Event) {
						studioZoom = min(studioZoom+1, 4)
						host.Refresh()
					}),
					gtui.Button("px-1 text-muted focus:text-ink", "−", func(e *gtui.Event) {
						studioZoom = max(studioZoom-1, 1)
						host.Refresh()
					}),
				),
			),
		),
	)
}

func studioVideo() gtui.Component {
	return gtui.Card("shrink-0 rounded border-line bg-panel px-1", "video",
		gtui.VideoPlayer(gtui.VideoOptions{
			Title:    "the long way round",
			Position: studioAt,
			Duration: 268,
			Playing:  studioPlaying,
			Step:     5,
			OnToggle: func() { studioPlaying = !studioPlaying; host.Refresh() },
			OnSeek: func(delta int) {
				studioAt = clamp(studioAt+delta, 0, 268)
				host.Refresh()
			},
			OnFullscreen: func() { host.Toast("info", "as big as it gets") },
		}),
	)
}

func studioMusic() gtui.Component {
	return gtui.Card("shrink-0 rounded border-line bg-panel px-1", "music",
		gtui.Player(gtui.PlayerOptions{
			Title:    "Nocturne in E flat",
			Subtitle: "piano · 320 kbps",
			Position: studioAt,
			Duration: 268,
			Playing:  studioPlaying,
			Volume:   55,
			Width:    12,
			Queue:    []string{"Gymnopédie no 1", "Clair de lune"},
		}),
		gtui.Row("gap-1 items-center pt-1",
			gtui.Text("text-muted nowrap", "level"),
			gtui.Equalizer("text-success", 10, studioTick),
		),
		gtui.SeekBar("pt-1", studioAt, 268, 20),
	)
}

// studioSizes shows the two ways a studioPicture can meet a box: stretched to it, or fitted
// inside it with its shape kept.
func studioSizes() gtui.Component {
	return gtui.Card("grow clip rounded border-line bg-panel px-1", "stretch and fit",
		gtui.Row("gap-1 items-center",
			gtui.ImageFrom("rounded border-line", studioShots[studioShot], 9, 4),
			gtui.FitFrom("rounded border-line", studioShots[studioShot], 9, 4),
		),
		gtui.Text("text-muted dim", "ImageFrom fills the box, FitFrom keeps the shape"),
		gtui.Video("pt-1", "trailer.mp4", studioAt, 268, studioPlaying),
	)
}

// studioPicture paints a sample image, twice as wide as it is tall, which is the
// shape a terminal draws best.
func studioPicture(seed int) image.Image {
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
