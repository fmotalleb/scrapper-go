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

func (ev *eval) GetConfig() config.Step {
	return ev.conf
}

// Execute implements Step.
func (ev *eval) Execute(p playwright.Page, v utils.Vars, r map[string]any) (interface{}, error) {
	// Evaluate locator and JS code templates
	locator, err := utils.EvaluateTemplate(ev.locator, v, p)
	if err != nil {
		slog.Error("failed to evaluate locator template", slog.Any("locator", ev.locator), log.ErrVal(err))
		return nil, err
	}

	jsCode, err := utils.EvaluateTemplate(ev.jsCode, v, p)
	if err != nil {
		slog.Error("failed to evaluate JS code template", slog.Any("jsCode", ev.jsCode), log.ErrVal(err))
		return nil, err
	}

	// Execute JS code on the page or locator
	var res interface{}
	if locator == "" {
		slog.Debug("evaluating JS code on page", slog.Any("jsCode", jsCode), slog.Any("params", ev.params))
		res, err = p.Evaluate(jsCode, ev.params)
	} else {
		slog.Debug("evaluating JS code on locator", slog.Any("locator", locator), slog.Any("jsCode", jsCode), slog.Any("params", ev.params))
		res, err = p.Locator(locator).Evaluate(jsCode, ev.params)
	}

	if err != nil {
		slog.Error("failed to evaluate JS code", log.ErrVal(err))
	}
	return res, err
}

func buildEval(step config.Step) (Step, error) {
	r := new(eval)
	r.conf = step

	// Extract locator
	if locator, ok := step["locator"].(string); ok {
		r.locator = locator
	} else {
		r.locator = "" // default empty string for no locator
	}

	// Extract JS code
	if jsCode, ok := step["eval"].(string); ok {
		r.jsCode = jsCode
	} else {
		return nil, fmt.Errorf("eval step must have a string input for 'eval' key, got: %T", step["eval"])
	}

	// Load additional parameters
	r.params = playwright.LocatorEvaluateOptions{}
	if params, err := utils.LoadParams[playwright.LocatorEvaluateOptions](step); err != nil {
		slog.Error("failed to load parameters for eval", log.ErrVal(err), slog.Any("step", step))
		return nil, err
	} else {
		r.params = *params
	}

	return r, nil
}
