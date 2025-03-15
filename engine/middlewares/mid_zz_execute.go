package middlewares

import (
	"fmt"

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
		return fmt.Errorf("found next middleware where its impossible to have one, no middleware are allowed after execute middleware")
	}
	result, err := s.Execute(p, v, r)
	if err != nil {
		return err
	} else {
		if key, ok := s.GetConfig()["set-var"]; ok {
			if key, ok := key.(string); ok {
				r[key] = result
			} else {
				return fmt.Errorf("expected set-var to be string got: %v", s.GetConfig())
			}
		}
	}
	return nil
}
