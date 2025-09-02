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
			_, ok := s["debug"].(string)
			return ok
		},
		Generator: buildDebug,
	})
}

type debug struct {
	text string
	conf config.Step
}

func (n *debug) GetConfig() config.Step {
	return n.conf
}

// Execute implements Step.
func (n *debug) Execute(p playwright.Page, v utils.Vars, r map[string]any) (interface{}, error) {
	// Evaluate the URL template with the given variables
	text, err := utils.EvaluateTemplate(n.text, v, p)
	if err != nil {
		slog.Error("failed to evaluate debug template, not causing a panic attack", slog.String("template", n.text), log.ErrVal(err))
		return nil, nil
	}

	slog.Info(text, slog.String("source", "debug"), slog.Any("step", n.GetConfig()))
	return text, nil
}

func buildDebug(step config.Step) (Step, error) {
	r := new(debug)
	r.conf = step
	// Extract the URL from the step
	var ok bool
	if r.text, ok = step["debug"].(string); !ok {
		return nil, errors.New("field to build debug node")
	}

	return r, nil
}
