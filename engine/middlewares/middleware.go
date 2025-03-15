package middlewares

import (
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
	current := middlewares[index]
	if index == len(middlewares)-1 {
		return current(p, s, v, r, nil)
	} else {
		return current(p, s, v, r, func(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any) error {
			return middlewareExec(index+1, p, s, v, r)
		})
	}
}
