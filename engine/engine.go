// Package engine contains core functionality of scrapper machine
package engine

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/engine/middlewares"
	"github.com/fmotalleb/scrapper-go/engine/steps"
	"github.com/playwright-community/playwright-go"
)

// ExecuteConfig loaded from cli or api
func ExecuteConfig(ctx context.Context, config config.ExecutionConfig) (map[string]any, error) {
	// Initialize Variables
	vars, err := initializeVariables(config.Pipeline.Vars)
	if err != nil {
		slog.Error("failed to load variables", slog.Any("err", err))
		return nil, fmt.Errorf("preflight check failed: %w", err)
	}
	if len(config.Pipeline.Steps) == 0 {
		return nil, fmt.Errorf("pipeline has no steps, preflight check failed")
	}

	// Start Playwright
	pw, err := playwright.Run()
	if err != nil {
		slog.Error("could not start Playwright", slog.Any("err", err))
		return nil, fmt.Errorf("playwright startup failed: %w", err)
	}
	defer func() {
		if err := pw.Stop(); err != nil {
			slog.Warn("failed to stop Playwright session", slog.Any("err", err))
		}
	}()

	// Stop Playwright when context is canceled
	killWithContext(ctx, pw)

	slog.Info("Playwright initialized")

	// Launch Browser
	browser, err := launchBrowser(pw, config.Pipeline.Browser, config.Pipeline.BrowserParams)
	if err != nil {
		slog.Error("could not launch browser", slog.Any("err", err))
		return nil, err
	}
	defer func() {
		if err := browser.Close(); err != nil {
			slog.Error("failed to close browser", slog.Any("err", err))
		}
	}()

	// Handle KeepRunning at the end
	defer handleKeepRunning(config.Pipeline.KeepRunning)

	// Create Page
	page, err := browser.NewPage(config.Pipeline.BrowserOptions)
	if err != nil {
		slog.Error("could not create page", slog.Any("err", err))
		return nil, fmt.Errorf("page creation failed: %w", err)
	}

	// Build Steps
	stepList, err := steps.BuildSteps(config.Pipeline.Steps)
	if err != nil {
		return nil, err
	}

	// Execute Steps
	result := make(map[string]any)
	for _, step := range stepList {
		if err = middlewares.HandleStep(page, step, vars, result); err != nil {
			return nil, err
		}
	}

	slog.Info("Execution finished", slog.Any("vars_snapshot", vars.Snapshot()), slog.Any("result", result))
	return result, nil
}

func ExecuteStream(ctx context.Context, config config.ExecutionConfig, pipeline <-chan []config.Step) (<-chan map[string]any, error) {
	vars, err := initializeVariables(config.Pipeline.Vars)
	if err != nil {
		slog.Error("failed to load variables", slog.Any("err", err))
		return nil, fmt.Errorf("preflight check failed: %w", err)
	}

	// Start Playwright
	pw, err := playwright.Run()
	if err != nil {
		slog.Error("could not start Playwright", slog.Any("err", err))
		return nil, fmt.Errorf("playwright startup failed: %w", err)
	}

	// defer func() {
	// 	if err := pw.Stop(); err != nil {
	// 		slog.Warn("failed to stop Playwright session", slog.Any("err", err))
	// 	}
	// }()

	// Stop Playwright when context is canceled
	killWithContext(ctx, pw)

	slog.Info("Playwright initialized")

	// Launch Browser
	browser, err := launchBrowser(pw, config.Pipeline.Browser, config.Pipeline.BrowserParams)
	if err != nil {
		slog.Error("could not launch browser", slog.Any("err", err))
		return nil, err
	}
	// defer func() {
	// 	if err := browser.Close(); err != nil {
	// 		slog.Error("failed to close browser", slog.Any("err", err))
	// 	}
	// }()

	// Handle KeepRunning at the end
	defer handleKeepRunning(config.Pipeline.KeepRunning)

	// Create Page
	page, err := browser.NewPage(config.Pipeline.BrowserOptions)
	if err != nil {
		slog.Error("could not create page", slog.Any("err", err))
		return nil, fmt.Errorf("page creation failed: %w", err)
	}

	resultChan := make(chan map[string]any)

	go func() {
		for i := range pipeline {
			result := make(map[string]any)
			stepList, err := steps.BuildSteps(i)
			if err != nil {
				slog.Error("failed to build step", slog.Any("step", i))
				continue
			}
			for _, step := range stepList {
				if err = middlewares.HandleStep(page, step, vars, result); err != nil {
					slog.Error("failed to handle step", slog.Any("step", i))
					continue
				}
			}
			resultChan <- result
		}
	}()

	return resultChan, nil
}

func killWithContext(ctx context.Context, pw *playwright.Playwright) {
	go func() {
		<-ctx.Done()
		_ = pw.Stop()
	}()
}

// launchBrowser initializes the correct browser based on config
func launchBrowser(pw *playwright.Playwright, browserType string, params playwright.BrowserTypeLaunchOptions) (playwright.Browser, error) {
	switch browserType {
	case "chromium":
		return pw.Chromium.Launch(params)
	case "firefox":
		return pw.Firefox.Launch(params)
	case "webkit":
		return pw.WebKit.Launch(params)
	default:
		return nil, fmt.Errorf("unsupported browser type: %s", browserType)
	}
}

// handleKeepRunning ensures the process stays alive for a set duration if needed
func handleKeepRunning(durationStr string) {
	if durationStr == "" {
		return
	}

	sleepTime, err := time.ParseDuration(durationStr)
	if err != nil {
		slog.Error("Invalid KeepRunning duration", slog.String("input", durationStr), slog.Any("err", err))
		return
	}

	slog.Info("Sleeping for duration", slog.Duration("duration", sleepTime))
	time.Sleep(sleepTime)
}
