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
			_, ok := s["omit"].(string)
			return ok
		},
		Generator: buildOmit,
	})
}

type omit struct {
	variable string
	conf     config.Step
}

func (s *omit) GetConfig() config.Step {
	return s.conf
}

// Execute implements Step.
func (c *omit) Execute(page playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	// Evaluate locator using template
	variable, err := utils.EvaluateTemplate(c.variable, vars, page)
	if err != nil {
		slog.Error("failed to evaluate locator template", slog.Any("variable", c.variable), log.ErrVal(err))
		return nil, err
	}
	delete(vars, variable)
	delete(result, variable)
	return nil, nil
}

func buildOmit(step config.Step) (Step, error) {
	r := &omit{}

	// Extract the locator for the click action
	if variable, ok := step["omit"].(string); ok {
		r.variable = variable
	} else {
		return nil, fmt.Errorf("expected 'click' key to be a string, got: %T", step["click"])
	}

	return r, nil
}
