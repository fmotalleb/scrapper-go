package steps

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/log"
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
func (s *sleep) Execute(p playwright.Page, v utils.Vars, r map[string]any) (interface{}, error) {
	// Evaluate the sleep duration template
	waitTime, err := utils.EvaluateTemplate(s.sleep, v, p)
	if err != nil {
		return nil, err
	}

	// Parse the evaluated string to duration
	value, err := time.ParseDuration(waitTime)
	if err != nil {
		slog.Error("failed to parse the sleep duration", slog.String("input", waitTime), log.ErrVal(err))
		return nil, err
	}

	// Log the sleep duration for confirmation
	slog.Info("sleeping for duration", slog.String("duration", value.String()))

	// Sleep for the evaluated duration
	time.Sleep(value)
	return nil, nil
}

func buildSleep(step config.Step) (Step, error) {
	r := new(sleep)
	r.conf = step

	// Retrieve the sleep duration string from the config step
	if sleep, ok := step["sleep"].(string); ok {
		r.sleep = sleep
	} else {
		return nil, fmt.Errorf("sleep must have a string input, got: %v", step)
	}

	return r, nil
}
