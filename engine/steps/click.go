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
			_, ok := s["click"].(string)
			return ok
		},
		Generator: BuildClick,
	})
}

type Click struct {
	locator string
	params  playwright.LocatorClickOptions
	conf    config.Step
}

func (s *Click) GetConfig() config.Step {
	return s.conf
}

// Execute implements Step.
func (c *Click) Execute(page playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	locator, err := utils.EvaluateTemplate(c.locator, vars, page)
	if err != nil {
		return nil, err
	}

	return nil, page.Locator(locator).Click(c.params)
}

func BuildClick(step config.Step) (Step, error) {
	r := new(Click)
	r.conf = step
	if locator, ok := step["click"].(string); ok {
		r.locator = locator
	} else {
		return nil, fmt.Errorf("click must have a string input got: %v", step)
	}

	r.params = playwright.LocatorClickOptions{}
	if params, err := utils.LoadParams[playwright.LocatorClickOptions](step); err != nil {
		slog.Error("failed to read params", slog.Any("err", err), slog.Any("step", step))
		return nil, err
	} else {
		r.params = *params
	}
	return r, nil
}
