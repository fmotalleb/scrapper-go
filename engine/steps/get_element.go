package steps

import (
	"fmt"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/log"
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
	ReadTypeHTML      = getTextMode("html")
	ReadTypeValue     = getTextMode("value")
	ReadTypeText      = getTextMode("text")
	ReadTypeTable     = getTextMode("table")
	ReadTypeTableFlat = getTextMode("table-flat")
)

var validModes = map[string]getTextMode{
	"html":       ReadTypeHTML,
	"value":      ReadTypeValue,
	"text":       ReadTypeText,
	"table":      ReadTypeTable,
	"table-flat": ReadTypeTableFlat,
}

type getText struct {
	locator string
	mode    getTextMode
	params  playwright.LocatorEvaluateOptions
	conf    config.Step
}

func (ge *getText) GetConfig() config.Step {
	return ge.conf
}

// Execute implements Step.
func (ge *getText) Execute(p playwright.Page, v utils.Vars, r map[string]any) (interface{}, error) {
	locator, err := utils.EvaluateTemplate(ge.locator, v, p)
	if err != nil {
		slog.Error("failed to evaluate locator template", slog.String("locator", ge.locator), log.ErrVal(err))
		return nil, err
	}
	element := p.Locator(locator)
	var output interface{}

	slog.Debug("fetching element content", slog.String("locator", locator), slog.String("mode", string(ge.mode)))

	switch ge.mode {
	case ReadTypeHTML:
		output, err = element.InnerHTML()

	case ReadTypeValue:
		output, err = element.InputValue()

	case ReadTypeText:
		output, err = element.TextContent()

	case ReadTypeTable:
		var body string
		body, err = element.InnerHTML()
		if err == nil {
			output, err = utils.ParseTable(fmt.Sprintf("<table>%s</table>", body))
		}

	case ReadTypeTableFlat:
		var body string
		body, err = element.InnerHTML()
		if err == nil {
			output, err = utils.ParseTableFlat(fmt.Sprintf("<table>%s</table>", body))
		}
	}

	if err != nil {
		slog.Error("failed to fetch content from element", slog.String("locator", locator), slog.String("mode", string(ge.mode)), log.ErrVal(err))
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
		r.mode = ReadTypeHTML // Default mode if not provided
	}

	// Load optional parameters
	r.params = playwright.LocatorEvaluateOptions{}
	if params, err := utils.LoadParams[playwright.LocatorEvaluateOptions](step); err != nil {
		slog.Error("failed to read params", log.ErrVal(err), slog.Any("step", step))
		return nil, err
	} else {
		r.params = *params
	}

	return r, nil
}
