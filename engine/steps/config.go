package steps

import (
	"errors"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
	"github.com/spf13/cast"
)

func init() {
	stepSelectors = append(stepSelectors, stepSelector{
		CanHandle: func(s config.Step) bool {
			_, ok := s["config"].(map[string]any)
			return ok
		},
		Generator: buildEngineConfig,
	})
}

type engineConfig struct {
	params map[string]any
	conf   config.Step
}

func (n *engineConfig) GetConfig() config.Step {
	return n.conf
}

// Execute implements Step.
func (n *engineConfig) Execute(p playwright.Page, v utils.Vars, r map[string]any) (interface{}, error) {
	if t, ok := n.params["timeout"]; ok {
		timeOut := cast.ToFloat64(t)
		p.SetDefaultTimeout(timeOut)
		slog.Debug("config updated", slog.String("key", "timeout"), slog.Float64("value", timeOut))
	}
	if t, ok := n.params["nav_timeout"]; ok {
		timeOut := cast.ToFloat64(t)
		p.SetDefaultNavigationTimeout(timeOut)
		slog.Debug("config updated", slog.String("key", "nav_timeout"), slog.Float64("value", timeOut))
	}
	return nil, nil
}

func buildEngineConfig(step config.Step) (Step, error) {
	r := new(engineConfig)
	r.conf = step
	// Extract the URL from the step
	var ok bool
	if r.params, ok = step["config"].(map[string]any); !ok {
		return nil, errors.New("field to build config node")
	}

	return r, nil
}
