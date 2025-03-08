package engine

import (
	"fmt"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/playwright-community/playwright-go"
)

func ExecuteConfig(config config.ExecutionConfig) error {

	vars := generateVariables(config.Pipeline.Vars)

	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("could not start Playwright: %v", err)
	}
	defer pw.Stop()

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
	return nil
}
