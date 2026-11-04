package main

import "github.com/orayew2002/gtui"

func pageSettings() gtui.Component {
	return gtui.Row("gap-1 grow",
		gtui.Column("grow rounded border-line bg-panel px-2 gap-0 key-f",
			gtui.Text("text-ink bold nowrap", gtui.T("settings")),
			gtui.KeyValue("",
				[2]string{gtui.T("language"), gtui.T(gtui.Language())},
				[2]string{"fallback", "en"},
				[2]string{"client timeout", "30 s"},
				[2]string{"max upload", "8 GB"},
				[2]string{"encode workers", "4"},
			),
			gtui.Divider("border-line"),
			gtui.Switch("#s-encode", "autoEncode", "auto encode", true),
			gtui.Switch("#s-release", "autoRelease", "auto release", false),
			gtui.Switch("#s-maint", "maintenance", "maintenance", false),
		),
		gtui.Column("w-1/3 rounded border-line bg-panel px-2 gap-0 key-k",
			gtui.Text("text-ink bold nowrap", gtui.T("language")),
			languagePicker(),
			gtui.Divider("border-line"),
			gtui.Text("text-ink bold nowrap", gtui.T("theme")),
			themeSwitch(),
		),
	)
}

func languagePicker() gtui.Component {
	codes := gtui.Languages()
	items := make([]gtui.Component, 0, len(codes))

	for _, code := range codes {
		pick := code
		mark := "○ "
		class := "row gap-1 text-muted focus:text-ink focus:bold #pick-" + code

		if code == gtui.Language() {
			mark = "◉ "
			class = "row gap-1 text-ink bold focus:underline #pick-" + code
		}

		items = append(items, gtui.Button(class, mark+gtui.T(code), func(e *gtui.Event) {
			app.Language(pick)
		}))
	}

	return gtui.Column("", items...)
}
