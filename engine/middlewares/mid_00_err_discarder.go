package middlewares

import (
	"log/slog"

	"github.com/fmotalleb/scrapper-go/engine/steps"
	"github.com/fmotalleb/scrapper-go/utils"
	playwright "github.com/playwright-community/playwright-go"
)

func init() {
	registerMiddleware(errorDiscard)
}

// errorDiscard implements Middleware.
func errorDiscard(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any, next execFunc) error {
	err := next(p, s, v, r)
	if errMode, ok := s.GetConfig()["on-error"].(string); ok {
		switch errMode {
		case "ignore":
			return nil
		case "print":
			slog.Error("error discarded but displayed because of on-error: print", slog.Any("err", err), slog.Any("step", s.GetConfig()))
			return nil
		case "panic":
			panic(err)
		}
	}
	return err
}
