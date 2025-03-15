// Package engine contains core functionality of scrapper machine
package engine

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/engine/middlewares"
	"github.com/fmotalleb/scrapper-go/engine/steps"
	"github.com/playwright-community/playwright-go"
)

func ExecuteConfig(config config.ExecutionConfig) (map[string]any, error) {
	vars, err := initializeVariables(config.Pipeline.Vars)
	if err != nil {
		slog.Error("failed to load variables", slog.Any("err", err))
		return nil, fmt.Errorf("could not start Preflight check failed: %v", err)
	}
	if len(config.Pipeline.Steps) == 0 {
		return nil, fmt.Errorf("given pipeline does not have any steps, preflight check failed")
	}
	pw, err := playwright.Run()
	if err != nil {
		slog.Error("could not start Playwright", slog.Any("err", err))
		return nil, fmt.Errorf("could not start Playwright: %v", err)
	}
	defer func() {
		if err := pw.Stop(); err != nil {
			slog.Warn("Failed to stop Playwright session", slog.Any("err", err))
		}
	}()

	slog.Info("Playwright initialized")
	var browser playwright.Browser

	switch config.Pipeline.Browser {
	case "chromium":
		browser, err = pw.Chromium.Launch(config.Pipeline.BrowserParams)
	case "firefox":
		browser, err = pw.Firefox.Launch(config.Pipeline.BrowserParams)
	case "webkit":
		browser, err = pw.WebKit.Launch(config.Pipeline.BrowserParams)
	}
	result := make(map[string]any)

	if err != nil {
		slog.Error("could not launch browser", slog.Any("err", err))
		return nil, fmt.Errorf("could not launch browser: %v", err)
	}
	defer func() {
		if err := browser.Close(); err != nil {
			slog.Error("failed to close browser", slog.Any("err", err))
		}
	}()
	defer func() {
		if config.Pipeline.KeepRunning != "" {
			sleepTime, err := time.ParseDuration(config.Pipeline.KeepRunning)
			if err != nil {
				slog.Error("Cannot parse duration in KeepRunning", slog.Any("err", err))
				return
			}
			slog.Info("Sleeping for duration", slog.Duration("duration", sleepTime))
			time.Sleep(sleepTime)
		}
	}()
	page, err := browser.NewPage(config.Pipeline.BrowserOptions)
	if err != nil {
		slog.Error("could not create page", slog.Any("err", err))
		return nil, fmt.Errorf("could not create page: %v", err)
	}
	// steps.StepSelectors
	steps, err := steps.BuildSteps(config.Pipeline.Steps)
	if err != nil {
		return nil, err
	}

	for _, step := range steps {
		if err = middlewares.HandleStep(page, step, vars, result); err != nil {
			return nil, err
		}
	}

	slog.Info("Execution finished", slog.Any("vars_snapshot", vars.Snapshot()), slog.Any("result", result))

	return result, nil
}
