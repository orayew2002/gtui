package main

import "github.com/orayew2002/gtui"

var (
	app      *gtui.App
	messages = []gtui.Message{
		{Role: "system", Who: "gtui", Text: "ask me about the framework"},
		{Role: "user", Who: "you", Text: "how do I center a box?"},
		{Role: "assistant", Who: "assistant", Text: `gtui.Center("full", gtui.Text("", "hi")) — center is a class, so any div can use it.`},
	}
)

func main() {
	app = gtui.Init()
	app.View(view)
	app.Render()
}

func view() gtui.Component {
	return gtui.Column("full bg-canvas p-1 gap-1",
		gtui.Div("grow scroll rounded border-line bg-panel px-2 key-c #thread",
			gtui.Chat("", messages...),
		),
		gtui.Row("gap-1",
			gtui.Input("grow rounded border-line px-1 #prompt", "prompt", "send a message"),
			gtui.Button("rounded border-accent text-accent focus:bold", "send", send),
		),
	)
}

func send(e *gtui.Event) {
	text := app.Value("#prompt")
	if text == "" {
		return
	}

	messages = append(messages,
		gtui.Message{Role: "user", Who: "you", Text: text},
		gtui.Message{Role: "assistant", Who: "assistant", Text: "typing..."},
	)

	app.SetValue("#prompt", "")
	app.Refresh()

	app.After(600_000_000, func() {
		messages[len(messages)-1].Text = "that is " + text + " — handled by your own code"
		app.Refresh()
	})
}
