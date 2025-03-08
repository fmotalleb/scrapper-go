package engine

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/playwright-community/playwright-go"
)

func ExecuteConfig(config config.ExecutionConfig) error {
	vars := initializeVariables(config.Pipeline.Vars)
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("could not start Playwright: %v", err)
	}
	defer func() {

		if err := pw.Stop(); err != nil {
			fmt.Println("Failed to stop the session")
		}
	}()

	browser, err := pw.Chromium.Launch(config.Pipeline.BrowserParams)
	if err != nil {
		return fmt.Errorf("could not launch browser: %v", err)
	}

	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %v", err)
	}

	for _, step := range config.Pipeline.Steps {
		if err := executeStep(page, step, vars); err != nil {
			return fmt.Errorf("Error executing step: %v, step: %v", err, step)
		}
	}

	slog.Info(fmt.Sprintf("%v", vars.Snapshot()))
	if config.Pipeline.KeepRunning != "" {
		sleepTime, err := time.ParseDuration(config.Pipeline.KeepRunning)
		if err != nil {
			return fmt.Errorf("cannot parse given duration in keep running: %s", config.Pipeline.KeepRunning)
		}
		time.Sleep(sleepTime)
	}

	return nil
}
