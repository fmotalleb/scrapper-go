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

func (o *omit) GetConfig() config.Step {
	return o.conf
}

// Execute implements Step.
func (o *omit) Execute(p playwright.Page, v utils.Vars, r map[string]any) (interface{}, error) {
	// Evaluate locator using template
	variable, err := utils.EvaluateTemplate(o.variable, v, p)
	if err != nil {
		slog.Error("failed to evaluate locator template", slog.Any("variable", o.variable), log.ErrVal(err))
		return nil, err
	}
	delete(v, variable)
	delete(r, variable)
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
