package shell

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/engine"
	"github.com/mitchellh/mapstructure"
	"github.com/rivo/tview"
)

func RunShell() {
	app := tview.NewApplication().
		EnableMouse(true).
		EnablePaste(true)
	go refresh(app, 3)

	inputFlex, query := buildInput()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mainFlex := tview.NewFlex()
	mainFlex.AddItem(inputFlex, 0, 1, true)

	go func() {
		link := bindToBrowser(ctx, query)
		output := buildResultView(link)
		mainFlex.AddItem(output, 0, 1, true)
	}()
	logBox := createLogView()

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 3, false).
		AddItem(logBox, 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func bindToBrowser(ctx context.Context, recvChan <-chan map[string]any) <-chan map[string]any {

	cfgMap := <-recvChan

	var cfg config.ExecutionConfig
	err := mapstructure.Decode(cfgMap, &cfg)
	if err != nil {
		slog.Error("failed to read config from body", slog.Any("err", err))
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}
	sendChan := make(chan map[string]any)
	pipe := make(chan []config.Step)
	resultChan, err := engine.ExecuteStream(ctx, cfg, pipe)
	go func() {
		for i := range resultChan {
			sendChan <- i
		}
	}()
	go func() {
		for i := range recvChan {
			var cfg config.Step
			err = mapstructure.Decode(i, &cfg)
			if err != nil {
				slog.Error("failed to read config from body", slog.Any("err", err))
				sendChan <- map[string]any{
					"error": err.Error(),
				}
			} else {
				pipe <- []config.Step{cfg}
			}
		}
	}()
	return sendChan
}

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
