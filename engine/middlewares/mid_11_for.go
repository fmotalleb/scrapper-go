package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/engine/steps"
	"github.com/fmotalleb/scrapper-go/log"
	"github.com/fmotalleb/scrapper-go/utils"
	playwright "github.com/playwright-community/playwright-go"
	"golang.org/x/exp/slog"
)

func init() {
	registerMiddleware(forLoop)
}

// conditionCheck implements Middleware.
func forLoop(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any, next execFunc) error {
	if s == nil {
		return errStepMissing
	}

	var cond string
	if raw, ok := s.GetConfig()["loop"]; !ok {
		return next(p, s, v, r)
	} else {
		var ok bool
		if cond, ok = raw.(string); !ok {
			slog.Error("loop condition found but was unable to read it as an string")
			return errors.New("")
		}
	}
	var nextSteps []steps.Step

	if stepsConfig, exists := s.GetConfig()["steps"]; exists {
		stepsArray, valid := stepsConfig.([]any)
		if !valid {
			slog.Error("expected steps to be an array of maps")
			return errors.New("steps configuration must be of type []map[string]any")
		}
		var innerSteps []config.Step
		// Iterate over the array and validate each item as a map
		for i, stepConfig := range stepsArray {
			stepMap, ok := stepConfig.(map[string]any)

			if !ok {
				slog.Error("step configuration is not a map", slog.Any("step", stepConfig))
				return fmt.Errorf("each step must be a map, got: %T at index %d", stepConfig, i)
			}
			innerSteps = append(innerSteps, stepMap)
			slog.Debug("step configuration received", slog.Any("step", stepMap))
		}
		var err error
		nextSteps, err = steps.BuildSteps(innerSteps)
		if err != nil {
			slog.Error("failed to build steps from configuration", log.ErrVal(err))
			return err
		}
	} else {
		slog.Error("steps configuration missing")
		return errors.New("steps configuration must be provided")
	}

	slog.Debug("loop condition received", slog.String("condition", cond))
	items, err := evaluateLoop(cond, v, p)
	if err != nil {
		slog.Error("failed to evaluate loop", slog.Any("err", err))
		return err
	}
	for _, i := range items {
		v.SetOnce("item", i)
		for _, step := range nextSteps {
			if err := HandleStep(p, step, v, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func evaluateLoop(cond string, v utils.Vars, p playwright.Page) ([]string, error) {
	result, err := utils.EvaluateTemplate(cond, v, p)
	if err != nil {
		return nil, err
	}

	if num, err := strconv.ParseFloat(result, 64); err == nil {
		if num != float64(int(num)) { // Not an integer
			return nil, errors.New("result is a non-integer number")
		}
		n := int(num)
		items := make([]string, n)
		for i := 0; i < n; i++ {
			items[i] = strconv.Itoa(i)
		}
		return items, nil
	}

	// The result is not a number, parse it as JSON
	var arr []interface{}
	if err := json.Unmarshal([]byte(result), &arr); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert []interface{} to []string
	items := make([]string, len(arr))
	for i, item := range arr {
		switch val := item.(type) {
		case float64:
			items[i] = strconv.FormatFloat(val, 'f', -1, 64)
		case string:
			items[i] = val
		default:
			return nil, errors.New("invalid type in JSON array")
		}
	}
	return items, nil
}
