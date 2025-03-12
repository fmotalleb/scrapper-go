package steps

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

type Select struct {
	locator string

	ValuesOrLabels []string
	Values         []string
	Indexes        []string
	Labels         []string

	params playwright.LocatorSelectOptionOptions
}

// Execute implements Step.
func (s *Select) Execute(page playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	locator, err := utils.EvaluateTemplate(s.locator, vars, page)
	if err != nil {
		return nil, err
	}
	selectOpt := new(playwright.SelectOptionValues)
	if values, err := utils.EvaluateTemplates(s.Values, vars, page); err == nil {
		selectOpt.Values = &values
	} else {
		slog.Error("failed to execute template on select Values", slog.Any("err", err), slog.Any("Values", s.Values))
		return nil, err
	}

	if values, err := utils.EvaluateTemplates(s.ValuesOrLabels, vars, page); err == nil {
		selectOpt.ValuesOrLabels = &values
	} else {
		slog.Error("failed to execute template on select ValuesOrLabels", slog.Any("err", err), slog.Any("ValuesOrLabels", s.ValuesOrLabels))
		return nil, err
	}

	if values, err := utils.EvaluateTemplates(s.Labels, vars, page); err == nil {
		selectOpt.Labels = &values
	} else {
		slog.Error("failed to execute template on select Labels", slog.Any("err", err), slog.Any("Labels", s.Labels))
		return nil, err
	}

	if values, err := utils.EvaluateTemplates(s.Indexes, vars, page); err == nil {
		if values, err := utils.MapItems(values, strconv.Atoi); err == nil {
			selectOpt.Indexes = &values
		} else {
			slog.Error("failed to convert indexes to integer", slog.Any("err", err), slog.Any("Indexes(AfterEval)", values))
			return nil, err
		}
	} else {
		slog.Error("failed to execute template on select Indexes", slog.Any("err", err), slog.Any("Indexes", s.Labels))
		return nil, err
	}

	return page.Locator(locator).SelectOption(*selectOpt, s.params)
}

func BuildSelect(step config.Step) (Step, error) {
	r := new(Select)
	if locator, ok := step["select"].(string); ok {
		r.locator = locator
	} else {
		return nil, fmt.Errorf("select must have a string input got: %v", step)
	}
	r.Values = utils.SingleOrMulti[string](step, "value")
	r.ValuesOrLabels = utils.SingleOrMulti[string](step, "value_or_label")
	r.ValuesOrLabels = append(r.ValuesOrLabels, utils.SingleOrMulti[string](step, "values_or_label")...)
	r.Labels = utils.SingleOrMulti[string](step, "label")
	r.Indexes = utils.SingleOrMulti[string](step, "index")
	if len(r.Values)+len(r.ValuesOrLabels)+len(r.Labels)+len(r.Indexes) == 0 {
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
