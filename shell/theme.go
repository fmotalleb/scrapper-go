package shell

import "github.com/gdamore/tcell/v2"

var (
	backgroundColor       = tcell.NewRGBColor(30, 30, 30)
	inputBg               = tcell.NewRGBColor(45, 45, 48)
	outputBg              = tcell.NewRGBColor(37, 37, 38)
	foreGround            = tcell.NewRGBColor(220, 220, 220)
	buttonBackgroundColor = tcell.NewRGBColor(70, 70, 70)
	buttonTextColor       = tcell.NewRGBColor(220, 220, 220)

	inputTheme = tcell.StyleDefault.
			Background(inputBg).
			Foreground(foreGround)

	outputTheme = tcell.StyleDefault.
			Background(outputBg).
			Foreground(foreGround)
)
