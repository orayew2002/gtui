// A storefront in three languages, with a hand built table, a checkout that
// steps forward, and pictures read off disk.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/orayew2002/gtui"
)

var (
	app   = gtui.Init()
	step  = 0
	cart  []string
	tick  = 0
	files = write()
)

func main() {
	languages()

	// one glyph of our own, next to the 198 the framework ships
	gtui.RegisterIcon("basket", "⌱")

	app.View(view)
	app.QuitOn("q", "ctrl+c")
	app.Animate(4, func(at int) { tick = at; app.Refresh() })
	app.Render()
}

// languages adds the three the shop speaks. A real app keeps them in JSON files
// and calls gtui.LoadLanguages("locales") instead.
func languages() {
	gtui.AddLanguage("en", map[string]string{
		"shop.name": "corner shop", "shop.cart": "basket", "shop.pay": "pay",
		"shop.empty": "your basket is empty", "shop.add": "add something to it",
	})
	gtui.AddLanguage("tk", map[string]string{
		"shop.name": "burçdaky dükan", "shop.cart": "sebet", "shop.pay": "töle",
		"shop.empty": "sebediňiz boş", "shop.add": "bir zat goşuň",
	})
	gtui.AddLanguage("ru", map[string]string{
		"shop.name": "лавка", "shop.cart": "корзина", "shop.pay": "оплатить",
		"shop.empty": "корзина пуста", "shop.add": "положите что-нибудь",
	})

	gtui.SetFallback("en")
	gtui.SetLanguage("en")
}

func view() gtui.Component {
	return gtui.Column("full bg-canvas",
		header(),
		gtui.Row("grow clip px-1 py-1 gap-1",
			gtui.Column("grow clip gap-1", goods(), checkout()),
			basket(),
		),
		gtui.StatusBar("border-t border-line px-1",
			[]gtui.Component{gtui.Text("nowrap text-muted", "language %s of %v", gtui.Language(), gtui.Languages())},
			[]gtui.Component{gtui.Text("nowrap text-muted dim", "q quits")},
		),
	)
}

func header() gtui.Component {
	buttons := make([]gtui.Component, 0, len(gtui.Languages()))
	for _, code := range gtui.Languages() {
		pick := code

		class := "px-1 text-muted focus:text-ink"
		if code == gtui.Language() {
			class = "px-1 text-accent bold"
		}

		buttons = append(buttons, gtui.Button(class, pick, func(e *gtui.Event) {
			app.Language(pick)
		}))
	}

	return gtui.Row("px-1 items-center gap-2 border-b border-line",
		gtui.IconText("text-accent bold nowrap", "basket", gtui.T("shop.name")),
		gtui.Typewriter("text-muted nowrap", "fresh every morning", tick),
		gtui.Spacer(),
		gtui.Row("gap-1", buttons...),
	)
}

// goods is a table written row by row, which is what Table is for when the data
// is not a slice of records.
func goods() gtui.Component {
	cols := []gtui.Col{
		{Label: "item", Width: "grow"},
		{Label: "price", Width: "w-8", Right: true},
		{Label: "", Width: "w-8", Right: true},
	}

	rows := make([]gtui.Component, 0, len(stock))
	for _, item := range stock {
		name, price := item[0], item[1]

		rows = append(rows, gtui.TableRow("hover:bg-raised", cols,
			gtui.Text("truncate", name),
			gtui.Text("nowrap", price),
			gtui.Clickable("text-accent focus:bold nowrap", func(e *gtui.Event) {
				cart = append(cart, name)
				app.Refresh()
			}, gtui.Text("nowrap", "add")),
		))
	}

	body := append([]gtui.Component{gtui.TableHead("", cols)}, rows...)

	return gtui.Card("grow clip rounded border-line bg-panel px-1", "today",
		gtui.Scrollbox("", 7, body...),
		gtui.Pagination("pt-1 text-muted", "page 1 of 3",
			func(e *gtui.Event) { app.Toast("info", "first page already") },
			func(e *gtui.Event) { app.Toast("info", "nothing newer") },
		),
	)
}

func checkout() gtui.Component {
	return gtui.Card("shrink-0 rounded border-line bg-panel px-1", gtui.T("shop.pay"),
		gtui.Stepper("", step, "basket", "address", "card"),
		gtui.Row("gap-2 items-center pt-1",
			gtui.Button("pill border-accent text-accent", "next", func(e *gtui.Event) {
				step = (step + 1) % 3
				if step == 0 {
					app.Alert(gtui.T("shop.pay"), "thank you, it is on its way")
				}
				app.Refresh()
			}),
			gtui.Spinner("text-accent", tick),
			gtui.Blink("text-warning nowrap", "card reader ready", tick),
		),
	)
}

func basket() gtui.Component {
	if len(cart) == 0 {
		return gtui.Column("w-28 shrink-0 clip gap-1",
			gtui.Card("rounded border-line bg-panel px-1", gtui.T("shop.cart"),
				gtui.EmptyState("", "basket", gtui.T("shop.empty"), gtui.T("shop.add")),
			),
			gtui.Card("rounded border-line bg-panel px-1", "loading",
				gtui.Skeleton("", 3, 22),
			),
			picture(),
		)
	}

	lines := make([]gtui.Component, 0, len(cart))
	for _, name := range cart {
		lines = append(lines, gtui.Text("truncate", name))
	}

	return gtui.Column("w-28 shrink-0 clip gap-1",
		gtui.Card("rounded border-line bg-panel px-1", gtui.T("shop.cart"),
			gtui.Columns("", []string{"grow", "w-6"},
				gtui.Column("gap-0", lines...),
				gtui.Text("text-right nowrap", "%d", len(cart)),
			),
			gtui.Alert("mt-1 border-success", "success", "ready", gtui.T("shop.pay")),
		),
		picture(),
	)
}

// picture shows the two ways a file meets a box, and a gif read back as a clip.
func picture() gtui.Component {
	if files == "" {
		return gtui.Card("rounded border-line bg-panel px-1", "picture",
			gtui.Text("text-muted dim nowrap", "no temp directory"),
		)
	}

	body := []gtui.Component{
		gtui.Row("gap-1",
			gtui.Image("rounded border-line", filepath.Join(files, "shelf.png"), 11, 4),
			gtui.Fit("rounded border-line", filepath.Join(files, "shelf.png"), 11, 4),
		),
	}

	if clip, err := gtui.OpenClipSized(filepath.Join(files, "shelf.gif"), 22, 4); err == nil {
		body = append(body, clip.Frame("rounded border-line", tick, 22, 4))
	}

	return gtui.Card("rounded border-line bg-panel px-1", "picture",
		gtui.Backdrop("h-3 rounded", filepath.Join(files, "shelf.png"), 0, 0,
			gtui.Text("bg-[#000000]/55 text-canvas nowrap px-1", " on the shelf "),
		),
		gtui.Column("gap-1 pt-1", body...),
	)
}

var stock = [][2]string{
	{"bread", "4.00"}, {"milk", "2.50"}, {"apples", "6.20"},
	{"cheese", "9.80"}, {"coffee", "14.00"}, {"honey", "11.40"},
}

// write paints a picture and an animation, and saves them where the components
// that read files can find them. A real app ships its own.
func write() string {
	dir, err := os.MkdirTemp("", "gtui-shop")
	if err != nil {
		return ""
	}

	still, err := os.Create(filepath.Join(dir, "shelf.png"))
	if err != nil {
		return ""
	}
	defer still.Close()

	if png.Encode(still, shelf(0)) != nil {
		return ""
	}

	moving, err := os.Create(filepath.Join(dir, "shelf.gif"))
	if err != nil {
		return ""
	}
	defer moving.Close()

	frames := make([]*image.Paletted, 0, 6)
	delays := make([]int, 0, 6)

	for step := range 6 {
		shot := image.NewPaletted(image.Rect(0, 0, 64, 32), palette())
		src := shelf(step)

		for y := range 32 {
			for x := range 64 {
				shot.Set(x, y, src.At(x, y))
			}
		}

		frames = append(frames, shot)
		delays = append(delays, 12)
	}

	if gif.EncodeAll(moving, &gif.GIF{Image: frames, Delay: delays}) != nil {
		return ""
	}

	return dir
}

func shelf(step int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 64, 32))

	for y := range 32 {
		for x := range 64 {
			wave := math.Sin(float64(x+step*4)/7 + float64(y)/9)
			img.Set(x, y, color.RGBA{
				R: uint8(120 + 60*wave),
				G: uint8(90 + 40*math.Cos(float64(y)/6)),
				B: uint8(70 + 50*wave),
				A: 255,
			})
		}
	}

	return img
}

func palette() color.Palette {
	out := make(color.Palette, 0, 64)
	for i := range 64 {
		v := uint8(i * 4)
		out = append(out, color.RGBA{v, v / 2, v / 3, 255})
	}
	return out
}
