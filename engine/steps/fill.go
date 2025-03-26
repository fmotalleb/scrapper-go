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
			_, ok := s["fill"].(string)
			return ok
		},
		Generator: buildFill,
	})
}

type fill struct {
	locator string
	value   string
	params  playwright.LocatorFillOptions
	conf    config.Step
}

func (s *fill) GetConfig() config.Step {
	return s.conf
}

// Execute implements Step.
func (f *fill) Execute(page playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	// Evaluate locator and value templates
	locator, err := utils.EvaluateTemplate(f.locator, vars, page)
	if err != nil {
		slog.Error("failed to evaluate locator template", slog.Any("locator", f.locator), log.ErrVal(err))
		return nil, err
	}

	fillValue, err := utils.EvaluateTemplate(f.value, vars, page)
	if err != nil {
		slog.Error("failed to evaluate value template", slog.Any("value", f.value), log.ErrVal(err))
		return nil, err
	}

	// Log the final locator and value
	slog.Debug("filling input", slog.Any("locator", locator), slog.Any("value", fillValue))

	// Fill the locator with the evaluated value
	return nil, page.Locator(locator).Fill(fillValue, f.params)
}

func buildFill(step config.Step) (Step, error) {
	r := new(fill)
	r.conf = step

	// Extract locator
	if locator, ok := step["fill"].(string); ok {
		r.locator = locator
	} else {
		return nil, fmt.Errorf("expected 'fill' to be a string, got: %T", step["fill"])
	}

	// Extract value
	if value, ok := step["value"].(string); ok {
		r.value = value
	} else {
		return nil, fmt.Errorf("expected 'value' to be a string, got: %T", step["value"])
	}

	// Load additional parameters
	r.params = playwright.LocatorFillOptions{}
	if params, err := utils.LoadParams[playwright.LocatorFillOptions](step); err != nil {
		slog.Error("failed to load parameters for fill", log.ErrVal(err), slog.Any("step", step))
		return nil, err
	} else {
		r.params = *params
	}

	return r, nil
}
