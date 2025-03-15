package steps

import (
	"fmt"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

func BuildSteps(steps []config.Step) ([]Step, error) {
	output := make([]Step, len(steps))
	var err error
	for index, step := range steps {
		for _, selector := range stepSelectors {
			if selector.CanHandle(step) {
				if output[index], err = selector.Generator(step); err != nil {
					return nil, err
				}
			}
		}
		if output[index] == nil {
			return nil, fmt.Errorf("failed to handle step: %v", step)
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
