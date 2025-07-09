package middlewares

import (
	"errors"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/engine/steps"
	"github.com/fmotalleb/scrapper-go/query"
	"github.com/fmotalleb/scrapper-go/utils"
	playwright "github.com/playwright-community/playwright-go"
)

var (
	errNoIf        = errors.New("no condition found in step")
	errTestFailed  = errors.New("condition check failed")
	errStepMissing = errors.New("step is missing")
)

func init() {
	registerMiddleware(conditionCheck)
}

// conditionCheck implements Middleware.
func conditionCheck(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any, next execFunc) error {
	if s == nil {
		return errStepMissing
	}

	exec, err := utils.ChainExec(s.GetConfig(),
		[]utils.ChainCallback{
			getCond,
			func(c any) (any, error) {
				slog.Debug("evaluating template", slog.Any("variables", v))
				str, ok := c.(string)
				if !ok {
					return nil, errors.New("invalid condition format")
				}
				return utils.TemplateEvalMapper(v, p)(str)
			},
			func(c any) (any, error) {
				str, ok := c.(string)
				if !ok {
					return nil, errors.New("invalid condition format after template evaluation")
				}
				q, err := query.ParseQuery(str)
				if err != nil {
					return nil, err
				}
				slog.Debug("testing condition", slog.Any("variables", v))
				return q.EvaluateQuery(v.Snapshot())
			},
		})
	if err != nil {
		switch err {
		case errNoIf:
			if next != nil {
				return next(p, s, v, r)
			}
			slog.Warn("no next middleware found, skipping execution")
			return nil
		default:
			return err
		}
	}

	if exec.(bool) {
		if next != nil {
			return next(p, s, v, r)
		}
		slog.Warn("no next middleware found, skipping execution")
		return nil
	}

	return errTestFailed
}

func getCond(c any) (any, error) {
	stepConfig, ok := c.(config.Step)
	if !ok {
		return nil, errors.New("invalid step configuration")
	}

	cond, exists := stepConfig["if"]
	if !exists {
		return nil, errNoIf
	}

	slog.Debug("found condition", slog.Any("condition", cond))
	return cond, nil
}
