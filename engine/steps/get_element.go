package steps

import (
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
	GET_HTML       = getTextMode("html")
	GET_VALUE      = getTextMode("value")
	GET_TEXT       = getTextMode("text")
	GET_TABLE      = getTextMode("table")
	GET_TABLE_FLAT = getTextMode("table-flat")
)

var validModes = map[string]getTextMode{
	"html":       GET_HTML,
	"value":      GET_VALUE,
	"text":       GET_TEXT,
	"table":      GET_TABLE,
	"table-flat": GET_TABLE_FLAT,
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
		slog.Error("failed to evaluate locator template", slog.String("locator", g.locator), slog.Any("error", err))
		return nil, err
	}
	element := page.Locator(locator)
	var output interface{}

	slog.Debug("fetching element content", slog.String("locator", locator), slog.String("mode", string(g.mode)))

	switch g.mode {
	case GET_HTML:
		output, err = element.InnerHTML()

	case GET_VALUE:
		output, err = element.InputValue()

	case GET_TEXT:
		output, err = element.TextContent()

	case GET_TABLE:
		var body string
		body, err = element.InnerHTML()
		if err == nil {
			output, err = utils.ParseTable(fmt.Sprintf("<table>%s</table>", body))
		}

	case GET_TABLE_FLAT:
		var body string
		body, err = element.InnerHTML()
		if err == nil {
			output, err = utils.ParseTableFlat(fmt.Sprintf("<table>%s</table>", body))
		}
	}

	if err != nil {
		slog.Error("failed to fetch content from element", slog.String("locator", locator), slog.String("mode", string(g.mode)), slog.Any("error", err))
	}
	return output, err
}

func buildElementSelector(step config.Step) (Step, error) {
	r := new(getText)
	r.conf = step

	// Extract locator
	if locator, ok := step["element"].(string); ok {
		r.locator = locator
	} else {
		return nil, fmt.Errorf("expected 'element' to be a string, got: %T", step["element"])
	}

	// Extract and validate mode
	if mode, ok := step["mode"].(string); ok && mode != "" {
		slog.Debug("selected mode", slog.String("mode", mode))
		if validMode, exists := validModes[mode]; exists {
			r.mode = validMode
		} else {
			err := fmt.Errorf("invalid mode '%s' selected, valid modes are: %v", mode, validModes)
			slog.Error("invalid mode", slog.String("mode", mode), slog.Any("valid-modes", validModes))
			return nil, err
		}
	} else {
		r.mode = GET_HTML // Default mode if not provided
	}

	// Load optional parameters
	r.params = playwright.LocatorEvaluateOptions{}
	if params, err := utils.LoadParams[playwright.LocatorEvaluateOptions](step); err != nil {
		slog.Error("failed to read params", slog.Any("error", err), slog.Any("step", step))
		return nil, err
	} else {
		r.params = *params
	}

	return r, nil
}
