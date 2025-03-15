package steps

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

func init() {
	stepSelectors = append(stepSelectors, stepSelector{
		CanHandle: func(s config.Step) bool {
			_, ok := s["select"].(string)
			return ok
		},
		Generator: BuildSelect,
	})
}

type Select struct {
	locator string

	vluesOrLabels []string
	values        []string
	indexes       []string
	labels        []string

	params playwright.LocatorSelectOptionOptions
	conf   config.Step
}

func (s *Select) GetConfig() config.Step {
	return s.conf
}

// Execute implements Step.
func (s *Select) Execute(page playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	locator, err := utils.EvaluateTemplate(s.locator, vars, page)
	if err != nil {
		return nil, err
	}
	selectOpt := new(playwright.SelectOptionValues)
	if values, err := utils.EvaluateTemplates(s.values, vars, page); err == nil {
		selectOpt.Values = &values
	} else {
		slog.Error("failed to execute template on select Values", slog.Any("err", err), slog.Any("Values", s.values))
		return nil, err
	}

	if values, err := utils.EvaluateTemplates(s.vluesOrLabels, vars, page); err == nil {
		selectOpt.ValuesOrLabels = &values
	} else {
		slog.Error("failed to execute template on select ValuesOrLabels", slog.Any("err", err), slog.Any("ValuesOrLabels", s.vluesOrLabels))
		return nil, err
	}

	if values, err := utils.EvaluateTemplates(s.labels, vars, page); err == nil {
		selectOpt.Labels = &values
	} else {
		slog.Error("failed to execute template on select Labels", slog.Any("err", err), slog.Any("Labels", s.labels))
		return nil, err
	}

	if values, err := utils.EvaluateTemplates(s.indexes, vars, page); err == nil {
		if values, err := utils.MapItems(values, strconv.Atoi); err == nil {
			selectOpt.Indexes = &values
		} else {
			slog.Error("failed to convert indexes to integer", slog.Any("err", err), slog.Any("Indexes(AfterEval)", values))
			return nil, err
		}
	} else {
		slog.Error("failed to execute template on select Indexes", slog.Any("err", err), slog.Any("Indexes", s.labels))
		return nil, err
	}

	return page.Locator(locator).SelectOption(*selectOpt, s.params)
}

func BuildSelect(step config.Step) (Step, error) {
	r := new(Select)
	r.conf = step
	if locator, ok := step["select"].(string); ok {
		r.locator = locator
	} else {
		return nil, fmt.Errorf("select must have a string input got: %v", step)
	}
	r.values = utils.SingleOrMulti[string](step, "value")
	r.vluesOrLabels = utils.SingleOrMulti[string](step, "value_or_label")
	r.vluesOrLabels = append(r.vluesOrLabels, utils.SingleOrMulti[string](step, "values_or_label")...)
	r.labels = utils.SingleOrMulti[string](step, "label")
	r.indexes = utils.SingleOrMulti[string](step, "index")
	if len(r.values)+len(r.vluesOrLabels)+len(r.labels)+len(r.indexes) == 0 {
		return nil, fmt.Errorf("cannot find any value to select, step: %v", step)
	}
	r.params = playwright.LocatorSelectOptionOptions{}
	if params, err := utils.LoadParams[playwright.LocatorSelectOptionOptions](step); err != nil {
		slog.Error("failed to read params", slog.Any("err", err), slog.Any("step", step))
		return nil, err
	} else {
		r.params = *params
	}
	return r, nil
}
