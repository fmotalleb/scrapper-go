package middlewares

import (
	"fmt"

	"github.com/fmotalleb/scrapper-go/engine/steps"
	"github.com/fmotalleb/scrapper-go/utils"
	playwright "github.com/playwright-community/playwright-go"
)

func init() {
	registerMiddleware(new(execute))
}

type execute struct{}

// exec implements Middleware.
func (e *execute) exec(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any) error {
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
