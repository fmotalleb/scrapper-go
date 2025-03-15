package steps

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

type getTextMode string

func init() {
	stepSelectors = append(stepSelectors, stepSelector{
		CanHandle: func(s config.Step) bool {
			_, ok := s["element"].(string)
			return ok
		},
		Generator: buildElementSelector,
	})
}

const (
	GET_HTML  = getTextMode("html")
	GET_VALUE = getTextMode("value")
	GET_TEXT  = getTextMode("text")
	GET_TABLE = getTextMode("table")
)

var validModes = map[string]getTextMode{
	"html":  GET_HTML,
	"value": GET_VALUE,
	"text":  GET_TEXT,
	"table": GET_TABLE,
}

type getText struct {
	locator string
	mode    getTextMode
	params  playwright.LocatorEvaluateOptions
	conf    config.Step
}

func (s *getText) GetConfig() config.Step {
	return s.conf
}

// Execute implements Step.
func (g *getText) Execute(page playwright.Page, vars utils.Vars, result map[string]any) (interface{}, error) {
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
		if err == nil {
			r, err = utils.ParseTable(fmt.Sprintf("<table>%s</table>", body))
		}
	}
	return r, err
}

func buildElementSelector(step config.Step) (Step, error) {
	r := new(getText)
	r.conf = step
	if locator, ok := step["element"].(string); ok {
		r.locator = locator
	} else {
		r.locator = ""
	}
	if mode, ok := step["mode"].(string); ok && mode != "" {
		slog.Debug("selected mode", slog.String("mode", mode))
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
