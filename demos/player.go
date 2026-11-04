package demos

import "github.com/orayew2002/gtui"

var (
	playerAt      = 65
	playerPlaying = true
	playerQueue   = []string{"Garagum gijeleri", "Yaz yagsy", "Altyn asyr"}
)

// Player is a music player: the component draws the chrome, the buttons under
// it do the work.
func Player(a *gtui.App, tick int) gtui.Component {
	host = a

	if playerPlaying && playerAt < 300 {
		playerAt = 65 + tick%235
	}

	return gtui.Center("full",
		gtui.Column("w-60 gap-1",
			gtui.Player(gtui.PlayerOptions{
				Class:    "key-p",
				Title:    playerQueue[0],
				Subtitle: "tk · 4k · 2.1 GB",
				Position: playerAt,
				Duration: 300,
				Playing:  playerPlaying,
				Volume:   70,
				Width:    44,
				Queue:    playerQueue,
			}),
			gtui.Row("gap-2 px-1 items-center",
				gtui.Button("text-muted focus:text-ink nowrap", "⏮", func(e *gtui.Event) {
					playerAt = max(0, playerAt-10)
					host.Refresh()
				}),
				gtui.Button("text-accent focus:bold nowrap", mark(), func(e *gtui.Event) {
					playerPlaying = !playerPlaying
					host.Refresh()
				}),
				gtui.Button("text-muted focus:text-ink nowrap", "⏭", func(e *gtui.Event) {
					playerAt = min(300, playerAt+10)
					host.Refresh()
				}),
				gtui.Spacer(),
				gtui.Equalizer("text-success nowrap", 10, tick),
			),
		),
	)
}

func mark() string {
	if playerPlaying {
		return "⏸"
	}
	return "▶"
}
