package shell

import (
	"fmt"
	"log"

	"github.com/rivo/tview"
)

func RunShell() {
	opt := ""
	test := ""
	app := tview.NewApplication().EnableMouse(true)

	opts := []string{"Option 1", "Option 2", "Option 3"}
	// Create a form
	form := tview.NewForm().
		AddDropDown("Choose an option", opts, 0, nil).
		AddInputField("Enter text", "", 20, nil, nil)
	form = form.AddButton("Submit", func() {
		_, opt = form.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
		test = form.GetFormItem(1).(*tview.InputField).GetText()
		app.Stop()
	}).
		AddButton("Quit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Interactive Form").SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(form, true).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
	fmt.Printf("Selected Option: %s\nEntered Text: %s\n", opt, test)
}
