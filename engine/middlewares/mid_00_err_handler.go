package middlewares

import (
	"log/slog"

	"github.com/fmotalleb/scrapper-go/engine/steps"
	"github.com/fmotalleb/scrapper-go/log"
	"github.com/fmotalleb/scrapper-go/utils"
	playwright "github.com/playwright-community/playwright-go"
)

func init() {
	registerMiddleware(errorHandler)
}

// errorHandler implements Middleware.
func errorHandler(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any, next execFunc) error {
	if next == nil {
		slog.Warn("no next middleware found, skipping execution", slog.Any("step", s.GetConfig()))
		return nil
	}

	err := next(p, s, v, r)
	if err == nil {
		return nil
	}

	if errMode, ok := s.GetConfig()["on-error"].(string); ok {
		switch errMode {
		case "ignore":
			return nil
		case "print":
			slog.Error("error discarded but displayed due to on-error: print", log.ErrVal(err), slog.Any("step", s.GetConfig()))
			return nil
		case "panic":
			panic(err)
		}
	}

	return err
}
