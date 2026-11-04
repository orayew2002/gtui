package main

import (
	"math/rand"

	"github.com/orayew2002/gtui"
)

func main() {
	app := gtui.Init()

	traffic := make([]float64, 32)
	for i := range traffic {
		traffic[i] = 10 + rand.Float64()*20
	}

	app.HTML(
		gtui.Column("full bg-canvas p-1 gap-1",
			gtui.Row("gap-1",
				gtui.StatCard("grow", "movies", "1 284", "412 GB"),
				gtui.StatCard("grow", "songs", "8 942", "96 GB"),
				gtui.StatCard("grow", "books", "317", "48 GB"),
			),
			gtui.Row("gap-1 grow",
				gtui.Card("grow rounded border-line bg-panel px-2", "traffic",
					gtui.Sparkline("text-success", traffic),
					gtui.LineChart("", traffic, 40, 6),
				),
				gtui.DataTable(gtui.TableOptions{
					Class:   "w-1/3",
					Title:   "queue",
					Caption: "3 jobs",
					Cols: []gtui.Col{
						{Label: "job", Width: "grow"},
						{Label: "state", Width: "w-10", Right: true},
					},
					Rows: [][]gtui.Component{
						{gtui.Text("", "garagum"), gtui.Pill("text-success", "done")},
						{gtui.Text("", "ýaz ýagşy"), gtui.Pill("text-warning", "encoding")},
						{gtui.Text("", "altyn asyr"), gtui.Pill("text-danger", "failed")},
					},
					Zebra:   true,
					Compact: true,
				}),
			),
		),
	)

	app.Render()
}
