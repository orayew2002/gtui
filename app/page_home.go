package main

import "github.com/orayew2002/gtui"

func pageHome() gtui.Component {
	return gtui.Column("grow gap-1",
		gtui.Row("gap-1 shrink-0",
			gtui.StatCard("grow", gtui.T("movies"), "1 284", "412 GB"),
			gtui.StatCard("grow", gtui.T("serials"), "96", "1.1 TB"),
			gtui.StatCard("grow", gtui.T("books"), "317", "48 GB"),
			gtui.StatCard("grow", gtui.T("songs"), "8 942", "96 GB"),
		),
		gtui.Row("gap-1 grow",
			gtui.Card("grow rounded border-line bg-panel px-2 key-d", gtui.T("panel.releases"),
				gtui.KeyValue("",
					[2]string{"usb-014 · ashgabat north", "done"},
					[2]string{"usb-013 · turkmenabat", "done"},
					[2]string{"usb-012 · mary depot", "encoding"},
					[2]string{"usb-011 · balkanabat", "failed"},
				),
			),
			gtui.Card("w-1/3 rounded border-line bg-panel px-2 key-b", gtui.T("panel.usage"),
				gtui.Gauge("text-success", 18, 20, 18),
				gtui.Gauge("text-success", 12, 20, 18),
				gtui.Gauge("text-warning", 7, 10, 18),
				gtui.Divider("border-line"),
				gtui.Sparkline("text-accent", []float64{3, 7, 4, 9, 12, 8, 14, 11, 16}),
			),
		),
	)
}
