package steps

import (
	"errors"
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
			// This enables the branching capabilities like for-loops
			_, steps := s["loop"].(string)
			return ok || steps
		},
		Generator: buildNop,
	})
}

type nop struct {
	text string
	conf config.Step
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
	r := new(nop)
	r.conf = step
	// Extract the URL from the step
	var ok bool
	if r.text, ok = step["nop"].(string); !ok {
		if r.text, ok = step["loop"].(string); !ok {
			return nil, errors.New("field to build nop node")
		}
	}

	return r, nil
}
