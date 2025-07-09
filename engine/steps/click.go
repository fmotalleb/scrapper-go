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

func (mc *click) GetConfig() config.Step {
	return mc.conf
}

// Execute implements Step.
func (mc *click) Execute(p playwright.Page, v utils.Vars, r map[string]any) (interface{}, error) {
	// Evaluate locator using template
	locator, err := utils.EvaluateTemplate(mc.locator, v, p)
	if err != nil {
		slog.Error("failed to evaluate locator template", slog.Any("locator", mc.locator), log.ErrVal(err))
		return nil, err
	}

	// Perform click on the locator
	err = p.Locator(locator).Click(mc.params)
	if err != nil {
		slog.Error("failed to click on locator", slog.Any("locator", locator), slog.Any("params", mc.params), log.ErrVal(err))
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
