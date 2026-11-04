package gtui

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"github.com/orayew2002/gtui/src/colors"
	"github.com/orayew2002/gtui/src/dom"
)

func Image(class, path string, width, height int) Component {
	file, err := os.Open(path)
	if err != nil {
		return imageFallback(class, path, width, height)
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return imageFallback(class, path, width, height)
	}

	return ImageFrom(class, src, width, height)
}

func ImageFrom(class string, src image.Image, width, height int) Component {
	if src == nil || width < 1 || height < 1 {
		return Div(class)
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

	return canvasComponent(class, canvas)
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
