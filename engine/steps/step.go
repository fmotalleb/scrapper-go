package steps

import (
	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

type StepGenerator func(config.Step) (Step, error)

type Step interface {
	Execute(playwright.Page, utils.Vars, map[string]any) (interface{}, error)
}
