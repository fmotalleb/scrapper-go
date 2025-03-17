package shell

import "github.com/gdamore/tcell/v2"

var (
	// Darker background and more vibrant colors for contrast and readability.
	backgroundColor       = tcell.NewRGBColor(24, 24, 26)    // Very dark gray, almost black.
	inputBg               = tcell.NewRGBColor(40, 40, 45)    // Slightly lighter dark for input areas.
	outputBg              = tcell.NewRGBColor(28, 28, 30)    // Dark gray for output.
	foreGround            = tcell.NewRGBColor(230, 230, 240) // Light grayish text for good contrast.
	buttonBackgroundColor = tcell.NewRGBColor(60, 60, 65)    // Subtle dark for buttons.
	buttonTextColor       = tcell.NewRGBColor(255, 255, 255) // White text for buttons for clear visibility.

	// Styles for input and output areas, with enhanced contrast.
	inputTheme = tcell.StyleDefault.
			Background(inputBg).
			Foreground(foreGround)

	outputTheme = tcell.StyleDefault.
			Background(outputBg).
			Foreground(foreGround)
	logOutputTheme = outputTheme.
			Dim(true)
)
