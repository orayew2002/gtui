// A settings screen: every kind of field, the pickers, and a settingsSave that reads
// the whole form back in one call.
package demos

import (
	"strconv"
	"time"

	"github.com/orayew2002/gtui"
)

func Settings(a *gtui.App, _ int) gtui.Component {
	host = a

	return gtui.Column("full",
		gtui.Row("px-1 items-center gap-2 border-b border-line",
			gtui.IconText("text-accent bold nowrap", "settings", "settings"),
			gtui.Text("text-muted dim truncate", "tab moves, enter opens a picker, ctrl+c quits"),
			gtui.Spacer(),
			gtui.Kbd("", "ctrl+c"),
		),
		gtui.Div("grow clip scroll px-1 py-1",
			gtui.TabView("grow",
				gtui.Tab{Key: "a", Label: "account", Active: true, Page: settingsAccount()},
				gtui.Tab{Key: "s", Label: "server", Page: settingsServer()},
				gtui.Tab{Key: "n", Label: "notices", Page: settingsNotices()},
			),
		),
		gtui.Row("border-t border-line px-1 gap-2 items-center",
			gtui.Text("text-muted dim grow truncate #saved", "nothing saved yet"),
			gtui.ButtonGroup("shrink-0",
				gtui.Button("text-muted", "reset", func(e *gtui.Event) {
					host.Refresh()
					host.Text("#saved", "back to the defaults")
				}),
				gtui.Button("text-accent bold", "save", settingsSave),
			),
		),
	)
}

func settingsAccount() []gtui.Component {
	return []gtui.Component{
		gtui.Row("gap-2 flex-wrap",
			gtui.Field("w-28", "name", gtui.Input("w-full rounded border-line #name", "name", "full name")),
			gtui.Field("w-28", "email", gtui.Email("w-full rounded border-line #email", "email", "you@example.com")),
		),
		gtui.Row("gap-2 flex-wrap",
			gtui.Field("w-28", "password", gtui.Password("w-full rounded border-line #pass", "pass", "secret")),
			gtui.Field("w-28", "phone", gtui.Phone("w-full rounded border-line #phone", "phone", "+993 6x xx xx xx")),
		),
		gtui.Field("w-58", "about", gtui.TextArea("w-full rounded border-line #about", "about", 3)),
		gtui.Row("gap-2 items-center flex-wrap pt-1",
			gtui.Checkbox("", "news", "send me news", true),
			gtui.Switch("", "beta", "join the beta", false),
			gtui.Toggle("", true),
			gtui.Rating("text-warning", 4, 5),
		),
	}
}

func settingsServer() []gtui.Component {
	return []gtui.Component{
		gtui.Row("gap-2 flex-wrap",
			gtui.Field("w-28", "address",
				gtui.InputGroup("",
					gtui.Text("px-1 text-muted nowrap", "https://"),
					gtui.URL("grow rounded border-line #url", "url", "example.com"),
					gtui.Button("rounded border-accent text-accent", "test", func(e *gtui.Event) {
						done := host.Busy("reaching the server")
						host.After(time.Second, func() {
							done()
							host.Toast("success", "answered in 41ms")
						})
					}),
				),
			),
			gtui.Field("w-16", "port", gtui.Number("w-full rounded border-line #port", "port", "8080")),
		),
		gtui.Row("gap-2 flex-wrap",
			gtui.Field("w-28", "region", gtui.Select("w-full rounded border-line #region", "region", "choose", "eu", "us", "asia")),
			gtui.Field("w-28", "pool", gtui.Selector("w-full rounded border-line #pool", "pool", "filter", "fast", "archive", "cold")),
		),
		gtui.Row("gap-2 flex-wrap",
			gtui.Field("w-28", "image", gtui.Autocomplete("w-full rounded border-line #image", "image", "type a name", "alpine", "almalinux", "arch")),
			gtui.Field("w-28", "search", gtui.Search("w-full rounded border-line #find", "find", "anything")),
		),
		gtui.Field("w-28", "workers", gtui.Range("w-full", "workers", 1, 16, 1, 4)),
		gtui.Row("gap-2 items-center flex-wrap pt-1",
			gtui.Text("text-muted nowrap", "plan"),
			gtui.Radios("", "plan", "free", "team", "business"),
		),
	}
}

func settingsNotices() []gtui.Component {
	return []gtui.Component{
		gtui.Row("gap-2 flex-wrap",
			gtui.Field("w-20", "from", gtui.Date("w-full rounded border-line #from", "from")),
			gtui.Column("gap-0",
				gtui.Text("text-muted nowrap", "when"),
				gtui.Calendar("rounded border-line #day", "day", time.Now()),
			),
		),
		gtui.Separator("", "channels"),
		gtui.Row("gap-2 flex-wrap items-center",
			gtui.Checkbox("", "mail", "mail", true),
			gtui.Checkbox("", "push", "push", false),
			gtui.Button("pill border-accent text-accent", "pick one", func(e *gtui.Event) {
				host.Choose("send through", []string{"mail", "push", "both"}, func(v string) {
					host.Text("#saved", "notices go by "+v)
				})
			}),
		),
	}
}

// settingsSave reads every control at once, which is what Values is for.
func settingsSave(e *gtui.Event) {
	values := host.Values()

	host.Text("#saved", "saved "+strconv.Itoa(len(values))+" fields, name is "+settingsQuote(values["name"]))
	host.Toast("success", "settings saved")
}

func settingsQuote(s string) string {
	if s == "" {
		return "empty"
	}
	return `"` + s + `"`
}
