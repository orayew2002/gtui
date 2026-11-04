// Languages. The gallery loads app/locales at startup, and this page shows the
// keys, the arguments, the fallback, and a language added from code.
package main

import "github.com/orayew2002/gtui"

var shownKeys = []string{
	"app.name", "nav.media", "nav.themes",
	"i18n.keys", "i18n.fallback", "msg.saved",
}

// pageI18n is built from T() calls, so switching the language rebuilds this
// page in the new one with no other work.
func pageI18n() gtui.Component {
	rows := []gtui.Component{
		card("a", phrase("i18n.title", "languages"),
			gtui.Row("gap-1 items-center", languageButtons()...),
			note("app.Language(code) sets the language and rebuilds the view"),
		),
		gtui.Row("gap-1",
			card("d", phrase("i18n.keys", "keys and translations"), keyTable()),
			card("f", phrase("i18n.format", "T with arguments"), phrases()...),
		),
	}

	if wide() {
		rows = append(rows, card("g", phrase("i18n.live", "switching rebuilds the view"),
			gtui.CodeView(gtui.CodeOptions{
				Lines: []string{
					`gtui.LoadLanguages("app/locales")`,
					`app.View(shell)`,
					`app.Language("tk") // SetLanguage then Refresh`,
				},
				First: 1,
			}),
		))
	}

	rows = append(rows, card("h", phrase("i18n.add", "a language from code"),
		note("AddLanguage takes a flat map, so a language need not be a file"),
		gtui.Row("gap-x-2 items-center flex-wrap",
			gtui.Button("pill border-accent text-accent #i18n-add", "add pirate", func(e *gtui.Event) {
				gtui.AddLanguage("pt", map[string]string{
					"i18n.title": "tongues",
					"nav.i18n":   "tongues",
				})
				app.Language("pt")
			}),
			gtui.Text("dim truncate", `AddLanguage("pt", map[string]string{…})`),
		),
	))

	return gtui.Column("grow gap-1", rows...)
}

func languageButtons() []gtui.Component {
	codes := gtui.Languages()
	buttons := make([]gtui.Component, 0, len(codes)+1)

	for _, code := range codes {
		class := "px-1 text-muted focus:bold focus:text-accent #lang-" + code
		if code == gtui.Language() {
			class = "px-1 text-accent bold focus:bold #lang-" + code
		}

		buttons = append(buttons, gtui.Button(class, code, func(e *gtui.Event) {
			app.Language(code)
			app.Select("#nav" + currentKey)
		}))
	}

	return append(buttons,
		gtui.Spacer(),
		gtui.Text("dim truncate", "%s · fallback en", gtui.T("i18n.pick")),
	)
}

func keyTable() gtui.Component {
	keys := shownKeys[:3]
	if wide() {
		keys = shownKeys
	}

	rows := make([][]gtui.Component, 0, len(keys))
	for _, key := range keys {
		rows = append(rows, gtui.Cells(key, gtui.T(key)))
	}

	return gtui.Table("",
		[]gtui.Col{{Label: "key", Width: "w-14"}, {Label: gtui.Language(), Class: "text-accent"}},
		rows...,
	)
}

// phrases shows what T does with arguments, with a key no translator has
// covered, and with a key that is in no file at all.
func phrases() []gtui.Component {
	rows := []gtui.Component{
		gtui.Text("dim truncate", `T("msg.greeting", "Ada", 3)`),
		gtui.Text("text-accent truncate", gtui.T("msg.greeting", "Ada", 3)),
		gtui.Text("dim truncate", `T("msg.onlyEnglish")`),
		gtui.Text("text-accent truncate", gtui.T("msg.onlyEnglish")),
		gtui.Text("dim truncate", `T("msg.nowhere")`),
		gtui.Text("text-warning truncate", gtui.T("msg.nowhere")),
	}

	if wide() {
		rows = append(rows, note("en.json is the fallback, and a key no file has shows itself"))
	}

	return rows
}
