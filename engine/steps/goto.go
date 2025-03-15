package steps

import (
	"fmt"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

func init() {
	stepSelectors = append(stepSelectors, stepSelector{
		CanHandle: func(s config.Step) bool {
			_, ok := s["goto"].(string)
			return ok
		},
		Generator: buildGoto,
	})
}

type gotoStep struct {
	url    string
	params playwright.PageGotoOptions
	conf   config.Step
}

func (s *gotoStep) GetConfig() config.Step {
	return s.conf
}

// Execute implements Step.
func (g *gotoStep) Execute(p playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	if url, err := utils.EvaluateTemplate(g.url, vars, p); err != nil {
		return nil, err
	} else {
		return p.Goto(url, g.params)
	}
}

func buildGoto(step config.Step) (Step, error) {
	r := new(gotoStep)
	r.conf = step
	r.params = playwright.PageGotoOptions{}
	var ok bool
	if r.url, ok = step["goto"].(string); !ok {
		return nil, fmt.Errorf("goto must have a string url field got: %v", step)
	}
	if params, err := utils.LoadParams[playwright.PageGotoOptions](step); err != nil {
		return nil, err
	} else {
		r.params = *params
	}
	return r, nil
}
