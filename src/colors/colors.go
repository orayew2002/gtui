package colors

type Color struct {
	R, G, B uint8
}

func RGB(r, g, b uint8) Color {
	return Color{r, g, b}
}

var (
	Black  = RGB(0, 0, 0)
	White  = RGB(255, 255, 255)
	Gray   = RGB(108, 117, 125)
	Red    = RGB(220, 53, 69)
	Green  = RGB(25, 135, 84)
	Blue   = RGB(13, 110, 253)
	Yellow = RGB(255, 193, 7)
	Orange = RGB(253, 126, 20)
	Purple = RGB(111, 66, 193)
	Pink   = RGB(214, 51, 132)
	Cyan   = RGB(13, 202, 240)
	Teal   = RGB(32, 201, 151)
	Indigo = RGB(102, 16, 242)
)

var (
	Primary   = Blue
	Secondary = Gray
	Success   = Green
	Danger    = Red
	Warning   = Yellow
	Info      = Cyan
)

var (
	Base      = RGB(30, 30, 46)
	Surface   = RGB(49, 50, 68)
	Overlay   = RGB(108, 112, 134)
	Text      = RGB(205, 214, 244)
	Lavender  = RGB(180, 190, 254)
	Sky       = RGB(137, 220, 235)
	Mauve     = RGB(203, 166, 247)
	Peach     = RGB(250, 179, 135)
	Rosewater = RGB(245, 224, 220)
)
