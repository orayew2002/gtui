package demos

import (
	"math"

	"github.com/orayew2002/gtui"
)

// traffic is a stand-in for a real series.
var traffic = series(32)

func series(n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = 20 + 12*math.Sin(float64(i)/4)
	}
	return out
}

// Dashboard is the shape most admin screens take: numbers on top, a chart and a
// table under them.
func Dashboard(a *gtui.App, tick int) gtui.Component {
	host = a

	return gtui.Column("full p-1 gap-1",
		gtui.Row("gap-1",
			gtui.StatCard("grow", "movies", "1 284", "412 GB"),
			gtui.StatCard("grow", "songs", "8 942", "96 GB"),
			gtui.StatCard("grow", "books", "317", "48 GB"),
		),
		gtui.Row("gap-1 grow clip",
			gtui.Card("grow clip rounded border-line bg-panel px-2", "traffic",
				gtui.Sparkline("text-success", traffic),
				gtui.LineChart("", traffic, 40, 6),
			),
			gtui.DataTable(gtui.TableOptions{
				Class:   "w-1/3 shrink-0",
				Title:   "queue",
				Caption: "3 jobs",
				Cols: []gtui.Col{
					{Label: "job", Width: "grow"},
					{Label: "state", Width: "w-10", Right: true},
				},
				Rows: [][]gtui.Component{
					{gtui.Text("truncate", "garagum"), gtui.Pill("text-success", "done")},
					{gtui.Text("truncate", "yaz yagsy"), gtui.Pill("text-warning", "encoding")},
					{gtui.Text("truncate", "altyn asyr"), gtui.Pill("text-danger", "failed")},
				},
				Zebra:   true,
				Compact: true,
			}),
		),
	)
}
