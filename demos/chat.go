package demos

import (
	"time"

	"github.com/orayew2002/gtui"
)

var thread = []gtui.Message{
	{Role: "system", Who: "gtui", Text: "ask me about the framework"},
	{Role: "user", Who: "you", Text: "how do I center a box?"},
	{Role: "assistant", Who: "assistant", Text: `gtui.Center("full", gtui.Text("", "hi")) — center is a class, so any div can use it.`},
}

// Chat is a thread that grows, which is a scroll pane and an input.
func Chat(a *gtui.App, tick int) gtui.Component {
	host = a

	return gtui.Column("full p-1 gap-1",
		gtui.Div("grow clip scroll rounded border-line bg-panel px-2 #thread",
			gtui.Chat("", thread...),
		),
		gtui.Row("gap-1 shrink-0",
			gtui.Input("grow rounded border-line px-1 #prompt", "prompt", "send a message"),
			gtui.Button("rounded border-accent text-accent focus:bold", "send", send),
		),
	)
}

func send(e *gtui.Event) {
	text := host.Value("#prompt")
	if text == "" {
		return
	}

	thread = append(thread,
		gtui.Message{Role: "user", Who: "you", Text: text},
		gtui.Message{Role: "assistant", Who: "assistant", Text: "typing..."},
	)

	host.SetValue("#prompt", "")
	host.Refresh()

	// the answer arrives later, the way a real one would
	host.After(600*time.Millisecond, func() {
		thread[len(thread)-1].Text = "that is " + text + " — handled by your own code"
		host.Refresh()
	})
}
