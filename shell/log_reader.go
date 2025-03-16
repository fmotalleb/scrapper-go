package shell

import (
	"fmt"

	"github.com/fmotalleb/scrapper-go/log"
	"github.com/rivo/tview"
)

func createLogView() *tview.TextView {
	logBox := tview.NewTextView()
	logBox.SetText("waiting for logs to arrive ...")
	logBox.SetBackgroundColor(backgroundColor).
		SetBorder(true)
	logBox.SetBorderStyle(outputTheme)
	logOutput := &logBuffer{
		viewPort: logBox,
		buffer:   "",
	}
	log.SetupLogger("DEBUG", logOutput)
	return logBox
}

type logBuffer struct {
	viewPort *tview.TextView
	buffer   string
}

// Write implements io.Writer.
func (l *logBuffer) Write(p []byte) (n int, err error) {
	l.buffer = fmt.Sprintf("%s%s", l.buffer, string(p))
	l.viewPort.SetText(l.buffer)
	l.viewPort.ScrollToEnd()
	return len(p), nil
}
