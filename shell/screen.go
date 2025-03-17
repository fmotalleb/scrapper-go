package shell

import (
	"time"

	"github.com/rivo/tview"
)

// refresh periodically redraws the tview application at the specified frequency.
func refresh(app *tview.Application, frequency int64) {
	// Calculate sleep time between redraws based on the frequency.
	sleepTime := time.Second / time.Duration(frequency)

	// Infinite loop to repeatedly draw the application at the specified frequency.
	for {
		app.Draw()            // Redraw the application.
		time.Sleep(sleepTime) // Sleep for the calculated time before the next redraw.
	}
}
