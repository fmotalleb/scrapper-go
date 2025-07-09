package shell

import (
	"fmt"

	"github.com/fmotalleb/scrapper-go/log"
	"github.com/rivo/tview"
)

// createLogView initializes the log view for displaying logs in the UI.
func createLogView() *tview.TextView {
	logBox := tview.NewTextView()
	logBox.SetText("waiting for logs to arrive ...").
		SetTextStyle(logOutputTheme).
		SetBackgroundColor(backgroundColor).
		SetBorder(true).
		SetBorderStyle(logOutputTheme)

	logOutput := &logBuffer{
		viewPort: logBox,
		buffer:   "",
	}

	// Setup the logger to write to the log view.
	_ = log.SetupLogger("DEBUG", logOutput)

	return logBox
}

// logBuffer is a custom type that implements io.Writer to capture log output.
type logBuffer struct {
	viewPort *tview.TextView
	buffer   string
}

// Write implements the io.Writer interface for logBuffer.
func (l *logBuffer) Write(p []byte) (n int, err error) {
	l.buffer = fmt.Sprintf("%s%s", l.buffer, string(p))
	l.viewPort.SetText(l.buffer) // Update the log view with the current buffer.
	l.viewPort.ScrollToEnd()     // Scroll to the end of the log view.
	return len(p), nil
}
