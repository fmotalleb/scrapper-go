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

func (filler *fill) GetConfig() config.Step {
	return filler.conf
}

// Execute implements Step.
func (filler *fill) Execute(p playwright.Page, v utils.Vars, r map[string]any) (interface{}, error) {
	// Evaluate locator and value templates
	locator, err := utils.EvaluateTemplate(filler.locator, v, p)
	if err != nil {
		slog.Error("failed to evaluate locator template", slog.Any("locator", filler.locator), log.ErrVal(err))
		return nil, err
	}

	fillValue, err := utils.EvaluateTemplate(filler.value, v, p)
	if err != nil {
		slog.Error("failed to evaluate value template", slog.Any("value", filler.value), log.ErrVal(err))
		return nil, err
	}

	// Log the final locator and value
	slog.Debug("filling input", slog.Any("locator", locator), slog.Any("value", fillValue))

	// Fill the locator with the evaluated value
	return nil, p.Locator(locator).Fill(fillValue, filler.params)
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
