package shell

import "github.com/rivo/tview"

// buildResultView creates and returns a TreeView to display the data from the provided channel.
func buildResultView(dataLink <-chan map[string]any) tview.Primitive {
	output := tview.NewTreeView()

	// Goroutine to listen for data from the channel and update the TreeView.
	go func() {
		for data := range dataLink {
			output.SetRoot(buildTree(data)) // Update the root of the TreeView with the new data.
		}
	}()

	// Set the background color and border style for the TreeView.
	output.SetBackgroundColor(outputBg).
		SetBorder(true).
		SetBorderStyle(outputTheme)

	return output
}
