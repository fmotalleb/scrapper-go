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
	no_if              = errors.New("no_if")
	test_failed        = errors.New("condition in if failed")
	step_missing_fatal = errors.New("step is missing")
)

func init() {
	registerMiddleware(conditionCheck)
}

// exec implements Middleware.
func conditionCheck(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any, next execFunc) error {
	if s == nil {
		return step_missing_fatal
	}
	exec, err := utils.ChainExec(s.GetConfig(),
		[]utils.ChainCallback{
			getCond,
			func(c any) (any, error) {
				slog.Debug("template ", slog.Any("variables", v))
				return utils.TemplateEvalMapper(v, p)(c.(string))
			},
			func(c any) (any, error) {
				if q, err := query.ParseQuery(c.(string)); err != nil {
					return nil, err
				} else {
					slog.Debug("testing against", slog.Any("variables", v))
					return q.EvaluateQuery(v.Snapshot())
				}
			},
		})
	switch err {
	case nil:
		if exec.(bool) {
			return next(p, s, v, r)
		}
		return test_failed
	case no_if:
		return next(p, s, v, r)
	default:
		return err
	}
}

func getCond(c any) (any, error) {

	if cond, ok := c.(config.Step)["if"]; ok {
		slog.Debug("found condition", slog.Any("condition", cond))
		return cond, nil
	} else {
		return nil, no_if
	}
}
