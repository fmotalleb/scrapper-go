package steps

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

func init() {
	stepSelectors = append(stepSelectors, stepSelector{
		CanHandle: func(s config.Step) bool {
			_, ok := s["sleep"].(string)
			return ok
		},
		Generator: buildSleep,
	})
}

type sleep struct {
	sleep string
	conf  config.Step
}

func (s *sleep) GetConfig() config.Step {
	return s.conf
}

// Execute implements Step.
func (s *sleep) Execute(page playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	waitTime, err := utils.EvaluateTemplate(s.sleep, vars, page)
	if err != nil {
		return nil, err
	}
	value, err := time.ParseDuration(waitTime)
	if err != nil {
		slog.Error("was not able to parse given duration from string", slog.Any("err", err))
		return nil, err
	}
	time.Sleep(value)
	return nil, nil
}

func buildSleep(step config.Step) (Step, error) {
	r := new(sleep)
	r.conf = step
	if sleep, ok := step["sleep"].(string); ok {
		r.sleep = sleep
	} else {
		return nil, fmt.Errorf("sleep must have a string input got: %v", step)
	}
	return r, nil
}
