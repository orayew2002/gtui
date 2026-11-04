package main

import (
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

// The gallery ships no binary files, so the pictures it shows are painted here
// and the ones that have to come off disk are written to a temp directory the
// first time a page asks for them.

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
