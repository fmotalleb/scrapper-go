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
			_, ok := s["nop"].(string)
			return ok
		},
		Generator: buildNop,
	})
}

type nop struct {
	text   string
	params playwright.PageGotoOptions
	conf   config.Step
}

func (s *nop) GetConfig() config.Step {
	return s.conf
}

// Execute implements Step.
func (g *nop) Execute(p playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	// Evaluate the URL template with the given variables
	text, err := utils.EvaluateTemplate(g.text, vars, p)
	if err != nil {
		slog.Error("failed to evaluate URL template", slog.String("url", g.text), log.ErrVal(err))
		return nil, err
	}

	return text, nil
}

func buildNop(step config.Step) (Step, error) {
	r := new(gotoStep)
	r.conf = step
	// Extract the URL from the step
	var ok bool
	if r.url, ok = step["nop"].(string); !ok {
		return nil, fmt.Errorf("nop must have a string field, got: %v", step)
	}

	// Load additional parameters
	r.params = playwright.PageGotoOptions{}
	if params, err := utils.LoadParams[playwright.PageGotoOptions](step); err != nil {
		return nil, err
	} else {
		r.params = *params
	}

	return r, nil
}
