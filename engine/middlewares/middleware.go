package middlewares

import (
	"github.com/fmotalleb/scrapper-go/engine/steps"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

func StartExecution(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any) error {
	for _, node := range middlewares {
		if err := node.exec(p, s, v, r); err != nil {
			return err
		}
	}
	return nil
}

var middlewares []Middleware

func registerMiddleware(m Middleware) {
	middlewares = append(middlewares, m)
}

type Middleware interface {
	exec(playwright.Page, steps.Step, utils.Vars, map[string]any) error
}
