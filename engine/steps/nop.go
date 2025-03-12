package steps

import (
	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

type Nop struct {
}

// Execute implements Step.
func (g *Nop) Execute(p playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	return nil, nil
}

func BuildNop(step config.Step) (Step, error) {
	r := new(Nop)
	return r, nil
}
