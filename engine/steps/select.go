package steps

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/log"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

func init() {
	stepSelectors = append(stepSelectors, stepSelector{
		CanHandle: func(s config.Step) bool {
			_, ok := s["select"].(string)
			return ok
		},
		Generator: buildSelect,
	})
}

type selectStep struct {
	locator string

	valuesOrLabels []string
	values         []string
	indexes        []string
	labels         []string

	params playwright.LocatorSelectOptionOptions
	conf   config.Step
}

func (ss *selectStep) GetConfig() config.Step {
	return ss.conf
}

// Execute implements Step.
func (ss *selectStep) Execute(p playwright.Page, v utils.Vars, r map[string]any) (interface{}, error) {
	locator, err := utils.EvaluateTemplate(ss.locator, v, p)
	if err != nil {
		return nil, err
	}

	selectOpt := new(playwright.SelectOptionValues)

	// Consolidate the template evaluation and assignment
	// Evaluating and setting values or labels
	setValuesOrLabels := func(fieldName string, source []string) (*[]string, error) {
		values, err := utils.EvaluateTemplates(source, v, p)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to execute template on %s", fieldName), log.ErrVal(err), slog.Any(fieldName, source))
			return nil, err
		}
		return &values, nil
	}

	// Evaluate all options for select
	if selectOpt.Values, err = setValuesOrLabels("values", ss.values); err != nil {
		return nil, err
	}
	if selectOpt.ValuesOrLabels, err = setValuesOrLabels("values_or_labels", ss.valuesOrLabels); err != nil {
		return nil, err
	}
	if selectOpt.Labels, err = setValuesOrLabels("labels", ss.labels); err != nil {
		return nil, err
	}
	// if err := setValuesOrLabels("indexes", s.indexes, selectOpt.Indexes); err != nil {
	// 	return nil, err
	// }

	// Convert index strings to integers
	if values, err := utils.EvaluateTemplates(ss.indexes, v, p); err == nil {
		if values, err := utils.MapItems(values, strconv.Atoi); err == nil {
			selectOpt.Indexes = &values
		} else {
			slog.Error("failed to convert indexes to integer", log.ErrVal(err), slog.Any("Indexes(AfterEval)", values))
			return nil, err
		}
	} else {
		slog.Error("failed to execute template on select Indexes", log.ErrVal(err), slog.Any("Indexes", ss.indexes))
		return nil, err
	}

	// If none of the selectors were populated, return an error
	if len(*selectOpt.Values) == 0 &&
		len(*selectOpt.ValuesOrLabels) == 0 &&
		len(*selectOpt.Labels) == 0 &&
		len(*selectOpt.Indexes) == 0 {
		return nil, fmt.Errorf("no valid select option values found for step: %v", ss.conf)
	}

	return p.Locator(locator).SelectOption(*selectOpt, ss.params)
}

func buildSelect(step config.Step) (Step, error) {
	r := new(selectStep)
	r.conf = step

	// Extract the locator for the select step
	if locator, ok := step["select"].(string); ok {
		r.locator = locator
	} else {
		return nil, fmt.Errorf("select must have a string input, got: %v", step)
	}

	// Load possible select options
	r.values = utils.SingleOrMulti[string](step, "value")
	r.valuesOrLabels = utils.SingleOrMulti[string](step, "value_or_label")
	r.valuesOrLabels = append(r.valuesOrLabels, utils.SingleOrMulti[string](step, "values_or_label")...)
	r.labels = utils.SingleOrMulti[string](step, "label")
	r.indexes = utils.SingleOrMulti[string](step, "index")

	// If no selection options are provided, return an error
	if len(r.values)+len(r.valuesOrLabels)+len(r.labels)+len(r.indexes) == 0 {
		return nil, fmt.Errorf("no valid selection options found, step: %v", step)
	}

	// Load additional parameters for the select action
	r.params = playwright.LocatorSelectOptionOptions{}
	if params, err := utils.LoadParams[playwright.LocatorSelectOptionOptions](step); err != nil {
		slog.Error("failed to read params", log.ErrVal(err), slog.Any("step", step))
		return nil, err
	} else {
		r.params = *params
	}

	return r, nil
}
