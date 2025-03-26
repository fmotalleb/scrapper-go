package steps

import (
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
			_, ok := s["click"].(string)
			return ok
		},
		Generator: buildClick,
	})
}

type click struct {
	locator string
	params  playwright.LocatorClickOptions
	conf    config.Step
}

func (s *click) GetConfig() config.Step {
	return s.conf
}

// Execute implements Step.
func (c *click) Execute(page playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	// Evaluate locator using template
	locator, err := utils.EvaluateTemplate(c.locator, vars, page)
	if err != nil {
		slog.Error("failed to evaluate locator template", slog.Any("locator", c.locator), log.ErrVal(err))
		return nil, err
	}

	// Perform click on the locator
	err = page.Locator(locator).Click(c.params)
	if err != nil {
		slog.Error("failed to click on locator", slog.Any("locator", locator), slog.Any("params", c.params), log.ErrVal(err))
	}
	return nil, err
}

func buildClick(step config.Step) (Step, error) {
	r := &click{
		conf: step,
	}

	// Extract the locator for the click action
	if locator, ok := step["click"].(string); ok {
		r.locator = locator
	} else {
		return nil, fmt.Errorf("expected 'click' key to be a string, got: %T", step["click"])
	}

	// Load additional parameters
	r.params = playwright.LocatorClickOptions{}
	if params, err := utils.LoadParams[playwright.LocatorClickOptions](step); err != nil {
		slog.Error("failed to load parameters for click", log.ErrVal(err), slog.Any("step", step))
		return nil, err
	} else {
		r.params = *params
	}

	return r, nil
}
