package main

import (
	"time"

	"github.com/orayew2002/gtui"
)

func pageWidgets() gtui.Component {
	return gtui.Column("grow gap-1",
		gtui.Row("gap-1",
			gtui.Card("w-1/3 rounded border-line bg-panel px-2 key-a", "inputs",
				field("text", gtui.Input("rounded border-line px-1 #w-text", "wText", "type here")),
				field("number", gtui.Number("rounded border-line px-1 w-10 #w-num", "wNum", "2026")),
				field("password", gtui.Password("rounded border-line px-1 #w-pass", "wPass", "secret")),
				field("email", gtui.Email("rounded border-line px-1 #w-mail", "wMail", "name@site.tm")),
				field("search", gtui.Search("rounded border-line px-1 #w-find", "wFind", "find")),
			),
			gtui.Card("w-1/3 rounded border-line bg-panel px-2 key-b", "pickers",
				field("selector", gtui.Selector("#w-disk", "wDisk", "pick a disk", storages...)),
				field("range", gtui.Range("#w-range", "wRange", 0, 100, 5, 60)),
				gtui.Switch("#w-switch", "wSwitch", "switch", true),
				gtui.Checkbox("#w-check", "wCheck", "checkbox", false),
				gtui.Radios("", "wRadio", "one", "two"),
			),
			gtui.Card("grow rounded border-line bg-panel px-2 key-d", "calendar",
				gtui.Calendar("#w-day", "wDay", time.Now()),
			),
		),
		gtui.Row("gap-1 grow",
			gtui.Card("w-1/3 rounded border-line bg-panel px-2 key-k", "code",
				gtui.CodeView(gtui.CodeOptions{
					Lines: []string{
						"func main() {",
						"\tgtui.Init().HTML(",
						"\t\tgtui.Text(\"bold\", \"hi\"),",
						"\t).Render()",
						"}",
					},
					First:  1,
					Active: 3,
					Marks:  map[int]string{2: "●"},
				}),
			),
			gtui.Card("w-1/3 rounded border-line bg-panel px-2 key-x", "chat",
				gtui.Chat("",
					gtui.Message{Role: "user", Who: "you", Text: "how do I center a box?"},
					gtui.Message{Role: "assistant", Text: "Center(\"full\", child) — center is a class."},
				),
			),
			gtui.Card("grow rounded border-line bg-panel px-2 key-w", "lists",
				gtui.Tabs("", "#w-tab",
					gtui.Tab{Key: "o", Label: "scroll", Page: []gtui.Component{
						gtui.ScrollList("#w-list", 5, 200, func(i int) gtui.Component {
							return gtui.Text("nowrap", "row %d of 200", i)
						}),
					}},
					gtui.Tab{Key: "p", Label: "accordion", Page: []gtui.Component{
						gtui.Accordion("",
							gtui.Section{Title: "general", Open: true, Body: []gtui.Component{
								gtui.KeyValue("", [2]string{"language", gtui.Language()}),
							}},
							gtui.Section{Title: "advanced", Body: []gtui.Component{
								gtui.Text("text-muted", "hidden until opened"),
							}},
						),
					}},
				),
				gtui.Div("#w-tab",
					gtui.ScrollList("#w-list", 5, 200, func(i int) gtui.Component {
						return gtui.Text("nowrap", "row %d of 200", i)
					}),
				),
				gtui.Breadcrumb("text-muted", "admin", "widgets", "lists"),
				gtui.Pagination("text-muted", "page 1 of 9", nil, nil),
			),
		),
	)
}
