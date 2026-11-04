// The views a whole app is made of: a chat thread, coloured logs, code with a
// gutter, an address bar and a status bar.
package main

import "github.com/orayew2002/gtui"

var chatLog = []gtui.Message{
	{Role: "system", Who: "session", Text: "restore started for disk-01"},
	{Role: "user", Who: "you", Text: "how long will the verify step take?"},
	{Role: "assistant", Who: "gtui", Text: "about four minutes for 12 disks. I will post each one as it lands."},
	{Role: "error", Who: "disk-07", Text: "checksum mismatch, retrying once"},
}

var logTimes = []string{"09:12:04", "09:12:09", "09:12:11", "09:12:18", "09:12:25"}

var logRows = []string{
	"listening on :8080",
	"GET /disks 200 1.4ms",
	"warn: disk-07 answered slowly",
	"error: checksum failed on disk-07",
	"GET /disks 200 1.1ms",
}

var codeLines = []string{
	"func pageViews() gtui.Component {",
	"\treturn gtui.Column(\"grow gap-1\",",
	"\t\tcard(\"a\", \"chat\", gtui.Chat(\"\", chatLog...)),",
	"\t)",
	"}",
}

// pageViews shows the big views a whole app is made of, each one a single call.
func pageViews() gtui.Component {
	return gtui.Column("grow gap-1",
		cards(
			card("a", "chat",
				note("one block per message, coloured by role"),
				gtui.Chat("", chatLog...),
			),
			card("b", "logs",
				note("a line saying error or warn picks its own colour"),
				gtui.Logs("", logRows, logTimes),
			),
		),
		cards(
			card("d", "code view",
				note("line numbers, a mark column and an active line"),
				gtui.CodeView(gtui.CodeOptions{
					Lines:  codeLines,
					First:  18,
					Active: 20,
					Marks:  map[int]string{19: "●"},
				}),
			),
			card("e", "address bar and toolbar",
				note("type a url and press enter"),
				gtui.AddressBar("#url", "url", "https://pkg.go.dev", func(value string) {
					say("opened", "went to %s", value)
				}),
				out("opened", "nothing opened yet"),
				gtui.Toolbar("mt-1",
					gtui.IconButton("", "save", "save", func(e *gtui.Event) { say("opened", "saved") }),
					gtui.IconButton("", "refresh", "reload", func(e *gtui.Event) { say("opened", "reloaded") }),
					gtui.Spacer(),
					gtui.Kbd("", "ctrl+s"),
				),
			),
		),
		card("f", "trails and pagers",
			note("Breadcrumb, Stepper and Pagination, all click to move"),
			gtui.Breadcrumb("", "disks", "disk-01", "snapshots"),
			gtui.Stepper("mt-1", 1, "pick", "verify", "restore"),
			gtui.Pagination("mt-1 #pager", "page 2 of 9",
				func(e *gtui.Event) { say("opened", "previous page") },
				func(e *gtui.Event) { say("opened", "next page") },
			),
		),
		gtui.StatusBar("shrink-0 rounded border-line bg-panel px-2",
			[]gtui.Component{
				gtui.Row("gap-1", gtui.StatusDot("", "ok"), gtui.Text("nowrap", "12 disks")),
			},
			[]gtui.Component{
				gtui.Text("nowrap text-muted dim", "StatusBar sits at the bottom of a real app"),
			},
		),
	)
}
