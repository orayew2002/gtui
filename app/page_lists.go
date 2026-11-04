package main

import "github.com/orayew2002/gtui"

func pageMovies() gtui.Component {
	data := make([][]gtui.Component, 0, len(movieRows))

	for _, item := range movieRows {
		state := "off"
		if item.Active {
			state = "on"
		}

		data = append(data, cells(
			item.Title, item.Storage, item.Duration, item.Size, item.Status, state, item.Tracks,
		))
	}

	cols := []gtui.Col{
		col("name", "w-1/4", false),
		col("storage", "w-10", false),
		col("duration", "w-8", false),
		col("size", "w-8", true),
		col("status", "w-12", false),
		col("active", "w-6", false),
		col("tracks", "grow", true),
	}

	titles := make([]string, 0, len(movieRows))
	for _, item := range movieRows {
		titles = append(titles, item.Title)
	}

	return listPage("f", cols, data, gtui.T("pager", 1284, 1, 92), titles)
}

func pageSerials() gtui.Component {
	cols := []gtui.Col{
		col("name", "w-1/4", false),
		{Label: gtui.T("col.seasons"), Width: "w-8", Right: true},
		{Label: gtui.T("col.episodes"), Width: "w-10", Right: true},
		col("size", "w-8", true),
		col("status", "w-12", false),
		col("active", "grow", true),
	}

	return listPage("f", cols, rows(serialRows), gtui.T("pager", 96, 1, 7), names(serialRows))
}

func pageBooks() gtui.Component {
	cols := []gtui.Col{
		col("name", "w-1/4", false),
		{Label: gtui.T("col.reader"), Width: "w-14", Right: false},
		{Label: gtui.T("col.chapters"), Width: "w-10", Right: true},
		col("size", "w-8", true),
		col("status", "w-12", false),
		col("active", "grow", true),
	}

	return listPage("f", cols, rows(bookRows), gtui.T("pager", 317, 1, 22), names(bookRows))
}

func pageSongs() gtui.Component {
	cols := []gtui.Col{
		col("name", "w-1/4", false),
		{Label: gtui.T("col.singer"), Width: "w-14"},
		{Label: gtui.T("col.album"), Width: "w-12"},
		col("duration", "w-8", true),
		{Label: gtui.T("col.plays"), Width: "w-10", Right: true},
		col("active", "grow", true),
	}

	return listPage("f", cols, rows(songRows), gtui.T("pager", 8942, 1, 596), names(songRows))
}

func pageStorages() gtui.Component {
	cols := []gtui.Col{
		col("name", "w-1/4", false),
		{Label: gtui.T("col.mount"), Width: "w-16"},
		col("size", "w-10", true),
		{Label: gtui.T("col.used"), Width: "w-10", Right: true},
		col("status", "grow", false),
	}

	return listPage("f", cols, rows(storageRows), gtui.T("pager", 4, 1, 1), names(storageRows))
}

func pageLogs() gtui.Component {
	lines := make([]string, 0, len(logLines)*4)
	stamps := make([]string, 0, len(logLines)*4)

	for range 4 {
		for _, line := range logLines {
			lines = append(lines, line.Text)
			stamps = append(stamps, line.Time)
		}
	}

	return gtui.Column("grow rounded border-line bg-panel px-1 key-f",
		gtui.Logs("", lines, stamps),
	)
}
