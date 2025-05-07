package steps

import (
	"fmt"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

func BuildSteps(steps []config.Step) ([]Step, error) {
	output := make([]Step, len(steps))
	var err error

	for index, step := range steps {
		var handled bool
		for _, selector := range stepSelectors {
			if selector.CanHandle(step) {
				// Log the selector being applied
				slog.Debug("applying selector", slog.Any("step", step), slog.Any("selector", selector))
				if output[index], err = selector.Generator(step); err != nil {
					return nil, fmt.Errorf("error generating step: %w", err)
				}
				handled = true
				break
			}
		}

		// Log failure if no selector could handle the step
		if !handled {
			slog.Error("no handler found for step", slog.Any("step", step))
			return nil, fmt.Errorf("no handler found for step: %v", step)
		}
	}

	return output, nil
}

var stepSelectors []stepSelector

type stepSelector struct {
	CanHandle func(config.Step) bool
	Generator stepGenerator
}

type stepGenerator func(config.Step) (Step, error)

type Step interface {
	Execute(playwright.Page, utils.Vars, map[string]any) (interface{}, error)
	GetConfig() config.Step
}
