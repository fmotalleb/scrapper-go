package shell

import (
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/rivo/tview"
)

func buildInput() (*tview.Flex, <-chan map[string]any) {
	queryChannel := make(chan map[string]any)
	inputBox := tview.NewTextArea()
	inputBox.SetTextStyle(inputTheme)
	inputBox.SetBorder(true)
	inputBox.SetBorderStyle(inputTheme)
	inputBox.SetText(`{
    "pipeline": {
        "browser": "chromium",
        "browser_params": {
            "headless": false
        }
    }
}`, false)

	// inputBox.SetBackgroundColor(inputBg)
	execBtn := tview.NewButton("Execute")
	execBtn.SetLabelColor(buttonTextColor)
	execBtn.SetBackgroundColor(inputBg)
	execBtn.SetStyle(inputTheme)
	execBtn.SetSelectedFunc(func() {
		text := inputBox.GetText()
		slog.Info("parsing", slog.String("query", text))
		if res, err := parseJSONToMap(text); err != nil {
			slog.Error("failed to parse query", slog.Any("error", err))
		} else {
			inputBox.SetText("", true)
			queryChannel <- res
		}
	})
	execBtn.
		SetBackgroundColor(buttonBackgroundColor)
	execBtn.SetBorder(true)

	inputFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	inputFlex.AddItem(inputBox, 0, 6, false)
	inputFlex.AddItem(execBtn, 0, 1, false)
	inputFlex.SetBackgroundColor(backgroundColor)
	return inputFlex, queryChannel
}

func parseJSONToMap(text string) (map[string]any, error) {
	var result map[string]any
	text = strings.ReplaceAll(text, `\"`, `"`)
	err := json.Unmarshal([]byte(text), &result)
	if err != nil {
		slog.Error("Error parsing JSON:", slog.Any("error", err))
		return nil, err
	}
	return result, nil
}
