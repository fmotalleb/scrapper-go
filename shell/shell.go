package shell

import (
	"context"

	"github.com/rivo/tview"
)

// RunShell initializes the tui application and binds it to browser.
func RunShell() {
	app := tview.NewApplication().
		EnableMouse(true).
		EnablePaste(true)

	// Start a goroutine to refresh the app at regular intervals.
	go refresh(app, 60)

	// Build input fields and context with cancel function.

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the main flex container and add the input fields.

	inputFlex, query := buildInput()
	mainFlex := tview.NewFlex()
	mainFlex.AddItem(inputFlex, 0, 1, true)

	// Start a goroutine to bind the query to a browser and build the output view.
	// this will result in a single input (no output) view until browser is initialized
	go func() {
		link := bindToBrowser(ctx, query)
		output := buildResultView(link)
		mainFlex.AddItem(output, 0, 1, true)
	}()

	// Create the log box view.
	logBox := createLogView()

	// Set up the final layout with mainFlex and logBox.
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 3, false).
		AddItem(logBox, 0, 1, false)

	// Run the tview application.
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
