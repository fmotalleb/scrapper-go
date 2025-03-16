package shell

import "github.com/rivo/tview"

func buildResultView(dataLink <-chan map[string]any) tview.Primitive {
	output := tview.NewTreeView()
	go func() {
		for i := range dataLink {
			output.SetRoot(buildTree(i))
		}
	}()
	output.SetBackgroundColor(outputBg)
	output.SetBorder(true).SetBorderStyle(outputTheme)

	return output
}
