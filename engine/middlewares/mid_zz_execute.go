package middlewares

import (
	"fmt"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/engine/steps"
	"github.com/fmotalleb/scrapper-go/utils"
	playwright "github.com/playwright-community/playwright-go"
)

func init() {
	registerMiddleware(exec)
}

// exec implements Middleware.
func exec(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any, next execFunc) error {
	if next != nil {
		return fmt.Errorf("unexpected middleware after execute middleware: execution should be the final step")
	}

	result, err := s.Execute(p, v, r)
	if err != nil {
		slog.Error("step execution failed", slog.Any("step", s.GetConfig()), slog.Any("error", err))
		return err
	}

	if key, ok := s.GetConfig()["set-var"]; ok {
		strKey, valid := key.(string)
		if !valid {
			return fmt.Errorf("expected set-var to be a string, got: %T", key)
		}
		r[strKey] = result
		slog.Debug("stored result in set-var", slog.String("key", strKey), slog.Any("value", result))
	}

	return nil
}
