package shell

import (
	"time"

	"github.com/rivo/tview"
)

func refresh(app *tview.Application, frequency int64) {
	sleepTime := time.Second / time.Duration(frequency)
	for {
		app.Draw()
		time.Sleep(sleepTime)
	}
}
