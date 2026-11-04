package main

import (
	"time"

	"github.com/orayew2002/gtui"
)

func pageFeedback() gtui.Component {
	return gtui.Column("grow gap-1",
		gtui.Banner("", "warning", "disk almost full", "three pools are over 90%",
			gtui.Button("pill border-warning text-warning", "clean up", func(e *gtui.Event) {
				app.Toast("success", "cleaned 4 GB")
			}),
			gtui.Button("pill border-line text-muted", "later", func(e *gtui.Event) {
				app.Toast("info", "reminder set")
			}),
		),
		// basis-23 is the width one alert wants: four fit across a wide terminal
		// and the row wraps to two by two at eighty columns
		gtui.Row("gap-1 flex-wrap",
			gtui.Alert("basis-23 clip border-info", "info", "info", "3 updates ready"),
			gtui.Alert("basis-23 clip border-success", "success", "success", "backup finished"),
			gtui.Alert("basis-23 clip border-warning", "warning", "warning", "disk 91% full"),
			gtui.Alert("basis-23 clip border-danger", "danger", "danger", "upload failed"),
		),
		cards(
			card("a", "toasts",
				note("they appear in the corner and fade on their own"),
				gtui.Row("gap-x-2 flex-wrap",
					gtui.Button("pill border-success text-success", "success", func(e *gtui.Event) {
						app.Toast("success", "saved")
					}),
					gtui.Button("pill border-warning text-warning", "warning", func(e *gtui.Event) {
						app.Toast("warning", "disk almost full")
					}),
					gtui.Button("pill border-danger text-danger", "danger", func(e *gtui.Event) {
						app.Toast("danger", "upload failed")
					}),
				),
			),
			card("b", "dialogs",
				note("esc closes them, tab moves inside them"),
				gtui.Row("gap-x-2 flex-wrap",
					gtui.Button("pill border-accent text-accent", "alert", func(e *gtui.Event) {
						app.Alert("heads up", "modals take focus and give it back")
					}),
					gtui.Button("pill border-danger text-danger", "confirm", func(e *gtui.Event) {
						app.Confirm("delete the pool?", "every snapshot goes with it", func() {
							say("busy-out", "pool deleted")
						})
					}),
					gtui.Button("pill border-line text-muted", "busy", func(e *gtui.Event) {
						stop := app.Busy("uploading")
						app.After(2*time.Second, func() {
							stop()
							say("busy-out", "upload finished")
						})
					}),
				),
				out("busy-out", "—"),
			),
		),
		cards(
			card("d", "progress and waiting",
				note("one frame of the gallery's animation loop"),
				gtui.Progress("text-accent", 7, 10, 18),
				gtui.Gauge("text-success", 3, 4, 18),
				gtui.Loader("", "restoring 4 of 12 disks", frame),
				gtui.Skeleton("", 2, 18),
			),
			card("h", "a modal of your own",
				note("the dialog helpers wrap app.Modal; this is the bare one"),
				gtui.Button("pill border-accent text-accent #fb-modal", "open", func(e *gtui.Event) {
					app.Modal("disk-07",
						gtui.KeyValue("", [2]string{"pool", "archive"}, [2]string{"state", "failed"}),
						gtui.Row("gap-1 justify-end pt-1",
							gtui.Button("pill border-line text-muted", "close", func(e *gtui.Event) {
								app.CloseModal()
								say("fb-modal-out", "closed from inside")
							}),
						),
					)
					say("fb-modal-out", "modal open: %v", app.ModalOpen())
				}),
				out("fb-modal-out", "nothing opened yet"),
			),
			card("g", "tags, marks and faces",
				note("Toggle and Rating show state, Avatar initials a name"),
				gtui.Row("gap-x-2 items-center flex-wrap",
					gtui.Badge("", "new"),
					gtui.Badge("bg-danger text-canvas", "3"),
					gtui.Pill("", "stable"),
					gtui.Pill("bg-success text-canvas", "live"),
				),
				gtui.Row("gap-x-2 items-center flex-wrap mt-1",
					gtui.Toggle("", true),
					gtui.Toggle("", false),
					gtui.Rating("text-warning", 4, 5),
					gtui.Spinner("text-accent", frame),
				),
				gtui.Row("gap-x-1 items-center mt-1",
					gtui.Avatar("bg-accent text-canvas", "Ada Lovelace"),
					gtui.Avatar("bg-info text-canvas", "Alan Turing"),
					gtui.Avatar("bg-success text-canvas", "Grace Hopper"),
					gtui.Text("text-muted dim", "three members"),
				),
			),
			card("f", "status and emptiness",
				note("StatusDot takes its colour from a kind name"),
				gtui.Row("gap-x-2 items-center flex-wrap",
					gtui.Row("gap-1", gtui.StatusDot("", "ok"), gtui.Text("nowrap", "live")),
					gtui.Row("gap-1", gtui.StatusDot("", "busy"), gtui.Text("nowrap", "syncing")),
					gtui.Row("gap-1", gtui.StatusDot("", "down"), gtui.Text("nowrap", "offline")),
				),
				gtui.Separator("", "nothing here"),
				gtui.EmptyState("", "files", "no backups yet", "make one to see it listed here",
					gtui.Button("pill border-accent text-accent", "new backup", func(e *gtui.Event) {
						app.Toast("success", "backup started")
					}),
				),
			),
		),
	)
}
