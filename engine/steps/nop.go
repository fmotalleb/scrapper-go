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

func (n *nop) GetConfig() config.Step {
	return n.conf
}

// Execute implements Step.
func (n *nop) Execute(p playwright.Page, v utils.Vars, r map[string]any) (interface{}, error) {
	// Evaluate the URL template with the given variables
	text, err := utils.EvaluateTemplate(n.text, v, p)
	if err != nil {
		slog.Error("failed to evaluate URL template", slog.String("url", n.text), log.ErrVal(err))
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
