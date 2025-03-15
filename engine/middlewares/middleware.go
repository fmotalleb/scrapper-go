package middlewares

import (
	"errors"

	"github.com/fmotalleb/scrapper-go/engine/steps"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

var middlewares []middleware

func registerMiddleware(m middleware) {
	middlewares = append(middlewares, m)
}

type execFunc = func(playwright.Page, steps.Step, utils.Vars, map[string]any) error

type middleware = func(playwright.Page, steps.Step, utils.Vars, map[string]any, execFunc) error

func HandleStep(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any) error {
	if len(middlewares) == 0 {
		return errors.New("no middlewares registered")
	}
	return middlewareExec(0, p, s, v, r)
}

func middlewareExec(index int, p playwright.Page, s steps.Step, v utils.Vars, r map[string]any) error {
	if index >= len(middlewares) {
		return errors.New("middleware index out of range, internal bug")
	}

	current := middlewares[index]
	next := func(p playwright.Page, s steps.Step, v utils.Vars, r map[string]any) error {
		return middlewareExec(index+1, p, s, v, r)
	}

	if index+1 >= len(middlewares) {
		next = nil // Last middleware, no next function
	}

	return current(p, s, v, r, next)
}
