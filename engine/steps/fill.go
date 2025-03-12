package steps

import (
	"fmt"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

type Fill struct {
	locator string
	Value   string
	params  playwright.LocatorFillOptions
}

// Execute implements Step.
func (f *Fill) Execute(page playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	locator, err := utils.EvaluateTemplate(f.locator, vars, page)
	if err != nil {
		return nil, err
	}
	fill, err := utils.EvaluateTemplate(f.locator, vars, page)
	if err != nil {
		return nil, err
	}
	return nil, page.Locator(locator).Fill(fill, f.params)
}

func BuildFill(step config.Step) (Step, error) {
	r := new(Fill)
	if locator, ok := step["fill"].(string); ok {
		r.locator = locator
	} else {
		return nil, fmt.Errorf("fill must have a string input got: %v", step)
	}

	if value, ok := step["value"].(string); ok {
		r.Value = value
	} else {
		return nil, fmt.Errorf("fill must have a string input got: %v", step)
	}

	r.params = playwright.LocatorFillOptions{}
	if params, err := utils.LoadParams[playwright.LocatorFillOptions](step); err != nil {
		slog.Error("failed to read params", slog.Any("err", err), slog.Any("step", step))
		return nil, err
	} else {
		r.params = *params
	}
	return r, nil
}
