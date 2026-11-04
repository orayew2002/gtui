<p align="center">
  <img src="docs/logo.svg" width="104" alt="gtui">
</p>

<h1 align="center">gtui</h1>

<p align="center"><b>Terminal user interfaces for Go, written like HTML.</b></p>

<p align="center">
  A div is a box · classes are utilities · links swap a target div<br>
  block &amp; flex layout · diffed cell rendering · themes · i18n · 60+ components
</p>

```go
package main

import "github.com/orayew2002/gtui"

func main() {
	gtui.Init().HTML(
		gtui.Center("full bg-canvas",
			gtui.Column("rounded border-accent bg-panel px-3 py-1 gap-0",
				gtui.Text("text-accent bold", "hello from gtui"),
				gtui.Text("text-muted", "press q to quit"),
			),
		),
	).Render()
}
```

```sh
go get github.com/orayew2002/gtui
```

## Why

Most TUI code places every cell by hand. gtui gives you the browser model instead — nested
boxes, utility classes, flexbox, wrapping text, focus and keyboard routing — on top of a renderer
that writes only the cells that actually changed.

- **Simple.** `Div(class, children...)` is the whole idea. Ids live in the class string: `#main`.
- **Fast.** A 120×40 frame with 220 nodes builds, lays out, paints and flushes in **~0.7 ms**;
  painting and flushing allocate **nothing** per frame.
- **Complete.** Tables, charts, forms, modals, images, animation, themes, translations.

```
 ╭───────────────┤n├╮ ╭──────────────────────────────────────────────────────────────╮
 │  ▣ panel  admin  │ │  movies                       en ru tk   ◐ theme   + create  │
 │ ╭──────────────╮ │ │  rows, status and encoding state                             │
 │ │ ⌕ search     │ │ ╰──────────────────────────────────────────────────────────────╯
 │ ╰──────────────╯ │
 │ MAIN             │ ╭───────────────────────────────────────────────────────────┤c├╮
 │  1 home          │ │╭──────────────────────────────────────────────────────────┤f├╮│
 │  y dynamic       │ ││ NAME             STORAGE   DURATION   SIZE     STATUS      ││
 │ CONTENT          │ ││──────────────────────────────────────────────────────────  ││
 │  2 movies        │ ││ first row        disk-01   1:52:10    2.1 GB   • ready     ││
 │  3 serials       │ ││ second row       disk-01   1:34:02    1.7 GB   • ready     ││
 │  4 books         │ ││ third row        disk-02   2:04:55    2.8 GB   • encoding  ││
 │  i widgets       │ │╰──────────────────────────────────────────────────────────  ││
 ╰──────────────────╯ ╰──────────────────────────────────────────────────────────────╯
```

Every box can own a jump key. `key-c` draws a `┤c├` tab in its top border, and pressing `c` moves
focus straight into it — a dense screen stays one keystroke deep.

## Three things to learn

**1. Boxes.** Every component takes its classes first, children after.

```go
gtui.Row("gap-1 p-1",
	gtui.Column("w-1/3 rounded border-line bg-panel px-2",
		gtui.Text("text-ink bold", "sidebar"),
	),
	gtui.Div("grow scroll #main"),
)
```

**2. Routing.** `Href` swaps the contents of a target div — that is a whole SPA.

```go
gtui.Href("px-1 current:bold key-1", "home", "#main", homePage())
```

**3. State.** Keep it in your own variables and re-render the part that changed.

```go
app.Open("#main", moviesPage())   // one box
app.Refresh()                     // the whole view
```

## Components

| | |
| --- | --- |
| structure | `Div` `Row` `Column` `Center` `Spacer` `Text` `Card` `Grid` `Split` `List` `ListItem` `Divider` `Toolbar` |
| navigation | `Href` `Tabs` `Breadcrumb` `Pagination` `Dropdown` `MenuItem` `Tree` `Accordion` `Stepper` `app.Palette` |
| data | `DataTable` (toolbar, zebra, selectable rows, `RowActions`, footer, pager) `Table` `Col` `Cells` `Badge` `Pill` `Toggle` `Avatar` `KeyValue` `StatCard` `Chips` `Rating` `ScrollList` |
| forms | `Input` `Password` `Number` `Email` `URL` `Phone` `Date` `Search` `TextArea` `Checkbox` `Radios` `Select` `Selector` `Range` `Switch` `Calendar` `Button` |
| feedback | `Alert` `Spinner` `Progress` `Gauge` `Skeleton` `Hint` `StatusBar` `app.Modal` `app.Toast` |
| charts | `BarChart` `Sparkline` `LineChart` `Donut` `Gauge` |
| views | `CodeView` `Chat` `Logs` `Player` `SeekBar` `AddressBar` `Image` `ImageFrom` `Video` |
| motion | `Marquee` `Typewriter` `Blink` `Pulse` `Wave` `Dots` `app.Animate` `app.Every` `app.After` |

## Classes

- layout `row` `col` `center` `block` `inline` `none` `flex-wrap` `absolute` `z-5`
- size `w-20` `h-3` `full` `w-1/3` `w-50%` `min-w-10` `max-h-8` `grow` `shrink-0` `basis-10`
- spacing `p-1` `px-2` `py-1` `pt-1` `m-1` `mx-2` `gap-1`
- align `justify-center|between|around|evenly` · `items-*` · `self-*` · `text-left|center|right`
- border `border` `rounded` `thick` `double` `dashed` `soft` `tight` `border-t|r|b|l|x|y` `border-accent`
- color `bg-panel` `text-ink` `border-line` `text-danger` `bg-[#1c1c1c]`
- text `bold` `dim` `italic` `underline` `strike` `nowrap` `pre`
- overflow `clip` `scroll` `scroll-x` `visible` `focusable`
- state `hover:` `focus:` `active:` `current:`

## Keyboard

| key | action |
| --- | --- |
| `tab` `shift+tab` | move focus anywhere |
| arrows | move inside the current box; scroll a `scroll` or `scroll-x` box |
| `enter` `space` | activate a link, button, switch, row |
| `key-x` class | press `x` to jump into that box (`┤x├` badge in its border, `key-hide` hides it) |
| `pgup` `pgdn` `home` `end`, wheel | scroll vertically or sideways |
| `← →` `home` `end` in a field | move the caret; `delete` and `backspace` edit in place |
| paste | arrives as one `paste` event and lands at the caret |
| `esc` | close a modal |
| `app.Bind("ctrl+p", fn)` | your own shortcut |

## Themes

Tokens resolve while styles compute, so switching costs one frame and no rebuild.

```go
app.Theme(colors.Light)       // or colors.Dark
app.ToggleTheme()
app.ThemeByName("midnight")

app.Tokens(func(t *colors.Theme) {    // your own palette
	t.Name = "midnight"
	t.Canvas = colors.RGB(8, 10, 20)
	t.Accent = colors.RGB(120, 200, 255)
})
```

Tokens: `canvas` `panel` `raised` `line` `muted` `ink` `accent` `success` `warning` `danger` `info`.

## Languages

One flat or nested JSON file per language.

```go
gtui.LoadLanguages("locales")      // en.json, tk.json, ru.json
gtui.SetFallback("en")

app.View(shell)                    // register the builder once
app.Language("tk")                 // rebuilds, keeps focus

gtui.T("movies")                   // "filmler"
gtui.T("action.create", "film")    // "+ film goş"
```

Nothing is stored per node, so a switch is one map swap plus a single rebuild.

## Animation

```go
app.Animate(8, func(tick int) {        // 8 frames per second
	app.Open("#main", dashboard(tick))
})

app.After(2*time.Second, func() { app.Toast("success", "done") })
```

The diffing writer means a moving UI still sends only the cells that changed.

## Examples

```sh
go run ./examples/hello        # one centered card
go run ./examples/dashboard    # stat cards, charts, a compact table
go run ./examples/form         # inputs, selector, range, switch, toasts
go run ./examples/browser      # tab bar, address bar, page, status bar
go run ./examples/player       # media player with seek, volume, queue
go run ./examples/editor       # file tree, code view with gutter, status bar
go run ./examples/chat         # assistant thread with a prompt box
go run ./app                   # the full admin panel
```

`app/` is a complete **admin panel** and the reference for structuring a real app:

```
app/main.go          wiring: languages, theme, routes, run
app/routes.go        the route table and the command palette
app/shell.go         sidebar, header, language and theme switches
app/page_*.go        one file per page
app/data.go          sample rows
app/helpers.go       small local helpers
```

Adding a page is one file plus one line in `routes.go`. The panel covers the whole component set
so you can try each one by hand:

| page | components on it |
| --- | --- |
| home | `StatCard` `Card` `KeyValue` `Gauge` `Sparkline` `Divider` |
| dynamic | `Progress` `Spinner` `Sparkline` `LineChart` `Wave` `Marquee` `Typewriter` `Pulse` `Dots` `Blink` `Range` `Switch` `Button` `app.Toast` |
| movies · serials · books · songs · storages | `DataTable` `Pill` `Toggle` `Pagination` with translated columns |
| editor | `Input` `Number` `Selector` `Range` `TextArea` `Radios` `Switch` `Chips` `Alert` `Button` `app.Modal` |
| uikit | `Button` `Dropdown` `Badge` `Pill` `Avatar` `Rating` `Chips` `Spinner` `Donut` `Hint` `Alert` `BarChart` `LineChart` `Stepper` `Tree` `Player` `Skeleton` |
| widgets | `Input` `Number` `Password` `Email` `Search` `Selector` `Range` `Switch` `Checkbox` `Radios` `Calendar` `CodeView` `Chat` `Tabs` `ScrollList` `Accordion` `Breadcrumb` `Pagination` |
| logs | `Logs` with level colors |
| settings | `KeyValue` `Switch` language and theme pickers |

```sh
go run ./app snapshot 110 30                      # print a frame as text
LANG_CODE=tk MODE=light PAGE=editor go run ./app snapshot 110 30
```

## API

```go
app := gtui.Init()

app.HTML(components...)         // mount a document
app.View(build)                 // mount a builder; Refresh / Language re-run it
app.Open(sel, components...)    // replace the children of one box
app.Add(sel, components...)
app.Query(sel) / app.QueryAll(sel)
app.Text(sel, text) / app.Class(sel, class, on)
app.Value(sel) / app.SetValue(sel, v) / app.Values()
app.Modal(title, body...) / app.CloseModal() / app.Toast(kind, msg) / app.Palette(title, actions...)
app.Focus(node) / app.Select(sel) / app.Bind(key, fn) / app.QuitOn(keys...)
app.Theme(t) / app.ToggleTheme() / app.Tokens(edit) / app.Language(code)
app.Animate(fps, fn) / app.Every(d, fn) / app.After(d, fn) / app.Do(fn)
app.Size() / app.OnResize(fn) / app.Copy(text)
app.UseMouse(on) / app.UseCursor(on)   // mouse capture and the real caret
app.Render()                    // open the terminal and run
app.Snapshot(w, h)              // render to plain text
```

## Requirements

Go 1.22+ and a terminal that speaks ANSI with truecolor. Mouse, resize and the alternate screen
are handled for you. `Image` decodes PNG/JPEG/GIF into half-block cells.

## License

MIT
