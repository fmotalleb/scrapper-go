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
			_, ok := s["eval"].(string)
			return ok
		},
		Generator: buildEval,
	})
}

type eval struct {
	locator string
	jsCode  string
	params  playwright.LocatorEvaluateOptions
	conf    config.Step
}

func (s *eval) GetConfig() config.Step {
	return s.conf
}

// Execute implements Step.
func (e *eval) Execute(page playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	locator, err := utils.EvaluateTemplate(e.locator, vars, page)
	if err != nil {
		return nil, err
	}
	jsCode, err := utils.EvaluateTemplate(e.jsCode, vars, page)
	if err != nil {
		return nil, err
	}
	var r interface{}
	if locator == "" {
		r, err = page.Evaluate(jsCode, e.params)
	} else {
		r, err = page.Locator(locator).Evaluate(jsCode, e.params)
	}
	return r, err
}

func buildEval(step config.Step) (Step, error) {
	r := new(eval)
	r.conf = step
	if locator, ok := step["locator"].(string); ok {
		r.locator = locator
	} else {
		r.locator = ""
	}

	if value, ok := step["eval"].(string); ok {
		r.jsCode = value
	} else {
		return nil, fmt.Errorf("eval must have a string input got: %v", step)
	}

	r.params = playwright.LocatorEvaluateOptions{}
	if params, err := utils.LoadParams[playwright.LocatorEvaluateOptions](step); err != nil {
		slog.Error("failed to read params", slog.Any("err", err), slog.Any("step", step))
		return nil, err
	} else {
		r.params = *params
	}
	return r, nil
}
