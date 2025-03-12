package steps

import (
	"fmt"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

type Eval struct {
	locator string
	JsCode  string
	params  playwright.LocatorEvaluateOptions
}

// Execute implements Step.
func (e *Eval) Execute(page playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	locator, err := utils.EvaluateTemplate(e.locator, vars, page)
	if err != nil {
		return nil, err
	}
	fill, err := utils.EvaluateTemplate(e.locator, vars, page)
	if err != nil {
		return nil, err
	}
	var r interface{}
	if locator == "" {
		r, err = page.Evaluate(fill, e.params)
	} else {
		r, err = page.Locator(locator).Evaluate(fill, e.params)
	}
	return r, err
}

func BuildEval(step config.Step) (Step, error) {
	r := new(Eval)
	if locator, ok := step["locator"].(string); ok {
		r.locator = locator
	} else {
		r.locator = ""
	}

	if value, ok := step["eval"].(string); ok {
		r.JsCode = value
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
