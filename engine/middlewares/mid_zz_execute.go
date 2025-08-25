package middlewares

import (
	"fmt"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/engine/steps"
	"github.com/fmotalleb/scrapper-go/log"
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
		slog.Error("step execution failed", slog.Any("step", s.GetConfig()), log.ErrVal(err))
		return err
	}

	if key, ok := s.GetConfig()["set-var"]; ok {
		strKey, valid := key.(string)
		if !valid {
			return fmt.Errorf("expected set-var to be a string, got: %T", key)
		}
		if nkey, err := utils.EvaluateTemplate(strKey, v, p); err != nil {
			slog.Error("failed to evaluate template for set-var key",
				slog.String("key", strKey),
				log.ErrVal(err),
			)
		} else {
			strKey = nkey
		}
		if err := setOrAppendWithMeta(r, strKey, result); err != nil {
			slog.Error("failed to store data in variable",
				slog.String("key", strKey),
				slog.Any("value", result),
				slog.Any("table", result),
				log.ErrVal(err),
			)
		} else {
			slog.Debug("stored result in set-var", slog.String("key", strKey), slog.Any("value", r[strKey]))
		}
	}
	return nil
}

func setOrAppendWithMeta(r map[string]any, key string, value any) error {
	metaKey := "__$" + key
	var isFirstTime bool
	var hasMeta bool
	if isFirstTime, hasMeta = r[metaKey].(bool); !hasMeta {
		// No meta field: set value as-is
		r[key] = value
		r[metaKey] = true
		return nil
	}
	existing := r[key]
	if isFirstTime {
		r[metaKey] = false
		r[key] = []any{existing, value}
		return nil
	}

	switch v := existing.(type) {
	case []any:
		// Already a slice
		r[key] = append(v, value)
	default:
		return fmt.Errorf("unsupported type for key '%s': %T", key, existing)
	}

	return nil
}
