package steps

import (
	"fmt"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

func init() {
	stepSelectors = append(stepSelectors, stepSelector{
		CanHandle: func(s config.Step) bool {
			_, ok := s["screenshot"].(string)
			return ok
		},
		Generator: BuildScreenShot,
	})
}

type ScreenShot struct {
	locator string
	params  playwright.LocatorScreenshotOptions
	conf    config.Step
}

func (s *ScreenShot) GetConfig() config.Step {
	return s.conf
}

// Execute implements Step.
func (sc *ScreenShot) Execute(page playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	locator, err := utils.EvaluateTemplate(sc.locator, vars, page)
	if err != nil {
		return nil, err
	}
	_, err = page.Locator(locator).Screenshot(sc.params)
	return nil, err
}

func BuildScreenShot(step config.Step) (Step, error) {
	r := new(ScreenShot)
	r.conf = step
	if locator, ok := step["screenshot"].(string); ok {
		r.locator = locator
	} else {
		return nil, fmt.Errorf("fill must have a string input got: %v", step)
	}

	r.params = playwright.LocatorScreenshotOptions{}
	if params, err := utils.LoadParams[playwright.LocatorScreenshotOptions](step); err != nil {
		slog.Error("failed to read params", slog.Any("err", err), slog.Any("step", step))
		return nil, err
	} else {
		r.params = *params
	}
	return r, nil
}
