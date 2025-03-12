package steps

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

type GetTextMode string

func init() {
	StepSelectors = append(StepSelectors, StepSelector{
		CanHandle: func(s config.Step) bool {
			_, ok := s["element"].(string)
			return ok
		},
		Generator: BuildElementSelector,
	})
}

const (
	GET_HTML  = GetTextMode("html")
	GET_VALUE = GetTextMode("value")
	GET_TEXT  = GetTextMode("text")
	GET_TABLE = GetTextMode("table")
)

var validModes = map[string]GetTextMode{
	"html":  GET_HTML,
	"value": GET_VALUE,
	"text":  GET_TEXT,
	"table": GET_TABLE,
}

type GetText struct {
	locator string
	mode    GetTextMode
	params  playwright.LocatorEvaluateOptions
}

// Execute implements Step.
func (g *GetText) Execute(page playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
	locator, err := utils.EvaluateTemplate(g.locator, vars, page)
	if err != nil {
		return nil, err
	}
	element := page.Locator(locator)
	var r interface{}
	switch g.mode {
	case GET_HTML:
		r, err = element.InnerHTML()

	case GET_VALUE:
		r, err = element.InputValue()

	case GET_TEXT:
		r, err = element.TextContent()

	case GET_TABLE:
		var body string
		body, err = element.InnerHTML()
		if err != nil {
			r, err = utils.ParseTable(fmt.Sprintf("<table>%s</table>", body))
		}
	}
	return r, err
}

func BuildElementSelector(step config.Step) (Step, error) {
	r := new(GetText)
	if locator, ok := step["element"].(string); ok {
		r.locator = locator
	} else {
		r.locator = ""
	}
	if mode, ok := step["mode"].(string); ok && mode != "" {
		if mode, ok := validModes[mode]; ok {
			r.mode = mode
		} else {
			slog.Error("cannot parse mode", slog.Any("step", step), slog.Any("valid-modes", validModes))
			return nil, errors.New("selected mode is not in valid modes")
		}
	} else {
		r.mode = GET_HTML
	}

	r.params = playwright.LocatorEvaluateOptions{}
	if params, err := utils.LoadParams[playwright.LocatorEvaluateOptions](step); err != nil {
		slog.Error("failed to read params", slog.Any("err", err), slog.Any("step", step))
		return nil, err
	} else {
		r.params = *params
	}
	return r, nil
}
