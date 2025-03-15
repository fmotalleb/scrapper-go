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
			_, ok := s["goto"].(string)
			return ok
		},
		Generator: buildGoto,
	})
}

type gotoStep struct {
	url    string
	params playwright.PageGotoOptions
	conf   config.Step
}

func (s *gotoStep) GetConfig() config.Step {
	return s.conf
}

// Execute implements Step.
func (g *gotoStep) Execute(p playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	// Evaluate the URL template with the given variables
	url, err := utils.EvaluateTemplate(g.url, vars, p)
	if err != nil {
		slog.Error("failed to evaluate URL template", slog.String("url", g.url), slog.Any("error", err))
		return nil, err
	}

	// Ensure the URL is not empty
	if url == "" {
		err := fmt.Errorf("evaluated URL is empty for goto step")
		slog.Error("empty URL in goto step", slog.String("url", url))
		return nil, err
	}

	slog.Debug("navigating to URL", slog.String("url", url))
	// Navigate to the evaluated URL
	return p.Goto(url, g.params)
}

func buildGoto(step config.Step) (Step, error) {
	r := new(gotoStep)
	r.conf = step
	// Extract the URL from the step
	var ok bool
	if r.url, ok = step["goto"].(string); !ok {
		return nil, fmt.Errorf("goto must have a string url field, got: %v", step)
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
