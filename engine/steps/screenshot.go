package steps

import (
	"encoding/base64"
	"fmt"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/log"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

func init() {
	stepSelectors = append(stepSelectors, stepSelector{
		CanHandle: func(s config.Step) bool {
			_, ok := s["screenshot"].(string)
			return ok
		},
		Generator: buildScreenShot,
	})
}

type screenShot struct {
	locator string
	params  playwright.LocatorScreenshotOptions
	conf    config.Step
}

func (sc *screenShot) GetConfig() config.Step {
	return sc.conf
}

// Execute implements Step.
func (sc *screenShot) Execute(p playwright.Page, v utils.Vars, r map[string]any) (interface{}, error) {
	// Evaluate the locator template
	locator, err := utils.EvaluateTemplate(sc.locator, v, p)
	if err != nil {
		slog.Error("failed to evaluate locator template", slog.String("locator", sc.locator), log.ErrVal(err))
		return nil, err
	}

	// Check if the locator is empty
	if locator == "" {
		err := fmt.Errorf("evaluated locator is empty")
		slog.Error("empty locator in screenshot step", slog.String("locator", locator))
		return nil, err
	}

	// Take a screenshot of the element identified by the locator
	slog.Debug("taking screenshot for locator", slog.String("locator", locator))
	data, err := p.Locator(locator).Screenshot(sc.params)
	if err != nil {
		slog.Error("failed to take screenshot", log.ErrVal(err))
		return nil, err
	}
	b64 := base64.StdEncoding.EncodeToString(data)
	return b64, nil
}

func buildScreenShot(step config.Step) (Step, error) {
	r := new(screenShot)
	r.conf = step

	// Extract the locator for the screenshot step
	if locator, ok := step["screenshot"].(string); ok {
		r.locator = locator
	} else {
		return nil, fmt.Errorf("screenshot must have a string input, got: %v", step)
	}

	// Load additional parameters for the screenshot
	r.params = playwright.LocatorScreenshotOptions{}
	if params, err := utils.LoadParams[playwright.LocatorScreenshotOptions](step); err != nil {
		slog.Error("failed to read screenshot params", log.ErrVal(err), slog.Any("step", step))
		return nil, err
	} else {
		r.params = *params
	}

	return r, nil
}
