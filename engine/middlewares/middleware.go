package middlewares

import (
	"fmt"

	"github.com/fmotalleb/scrapper-go/engine/steps"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

var middlewares []middleware

func registerMiddleware(m middleware) {
	middlewares = append(middlewares, m)
}

type execFunc = func(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any) error

type middleware = func(playwright.Page, steps.Step, utils.Vars, map[string]any, execFunc) error

func HandleStep(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any) error {
	return middlewareExec(0, p, s, v, r)
}

func middlewareExec(index int, p playwright.Page, s steps.Step, v utils.Vars, r map[string]any) error {
	if len(middlewares) == 0 {
		return fmt.Errorf("no middlewares registered")
	}
	if index >= len(middlewares) {
		return fmt.Errorf("reached end of middleware stack but the call stack is not closed, internal bug")
	}
	current := middlewares[index]
	if index == len(middlewares)-1 {
		return current(p, s, v, r, nil)
	} else {
		return current(p, s, v, r, func(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any) error {
			return middlewareExec(index+1, p, s, v, r)
		})
	}
}
