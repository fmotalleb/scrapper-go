package steps

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/playwright-community/playwright-go"
)

func init() {
	stepSelectors = append(stepSelectors, stepSelector{
		CanHandle: func(s config.Step) bool {
			_, ok := s["mouse"].(string)
			return ok
		},
		Generator: buildMouse,
	})
}

type mouseAction string

const (
	mouseActionClick       mouseAction = "click"
	mouseActionDoubleClick mouseAction = "double-click"
	mouseActionScroll      mouseAction = "scroll"
	mouseActionMove        mouseAction = "move"

	mouseActionUp   mouseAction = "up"
	mouseActionDown mouseAction = "down"
)

type mouse struct {
	location []float64
	action   mouseAction
	// params   any
	conf config.Step
}

func (mp *mouse) GetConfig() config.Step {
	return mp.conf
}

// Execute implements Step.
func (mp *mouse) Execute(p playwright.Page, v utils.Vars, r map[string]any) (interface{}, error) {
	// Evaluate locator using template
	// locator, err := utils.EvaluateTemplate(c.location, vars, page)
	// if err != nil {
	// 	slog.Error("failed to evaluate locator template", slog.Any("locator", c.location), log.ErrVal(err))
	// 	return nil, err
	// }
	switch mp.action {
	case mouseActionClick:

		return nil, p.Mouse().Click(mp.location[0], mp.location[1])

	case mouseActionDoubleClick:
		return nil, p.Mouse().Dblclick(mp.location[0], mp.location[1])
	case mouseActionMove:
		return nil, p.Mouse().Move(mp.location[0], mp.location[1])
	case mouseActionScroll:
		return nil, p.Mouse().Wheel(mp.location[0], mp.location[1])
	case mouseActionDown:
		return nil, p.Mouse().Down()
	case mouseActionUp:
		return nil, p.Mouse().Up()
	}

	return nil, fmt.Errorf("unknown mouse action: %s", mp.action)
}

func buildMouse(step config.Step) (Step, error) {
	r := &mouse{
		conf: step,
	}

	// Extract the location for the click action
	if location, ok := step["mouse"].(string); ok {
		pointsStr := strings.Split(location, ",")
		if len(pointsStr) != 2 {
			return nil, fmt.Errorf("expected 'mouse' key to be a string with 2 points, got: %T", step["mouse"])
		}
		x, err := strconv.ParseFloat(pointsStr[0], 64)
		if err != nil {
			return nil, fmt.Errorf("expected 'mouse' key to be a string with 2 int points, got: %T", step["mouse"])
		}
		y, err := strconv.ParseFloat(pointsStr[1], 64)
		if err != nil {
			return nil, fmt.Errorf("expected 'mouse' key to be a string with 2 int points, got: %T", step["mouse"])
		}
		r.location = []float64{x, y}
	} else {
		return nil, fmt.Errorf("expected 'click' key to be a string, got: %T", step["click"])
	}
	if act, ok := step["action"]; ok {
		r.action = mouseAction(act.(string))
	} else {
		r.action = mouseActionClick
	}

	// Load additional parameters
	// r.params = playwright.LocatorClickOptions{}
	// if params, err := utils.LoadParams[playwright.LocatorClickOptions](step); err != nil {
	// 	slog.Error("failed to load parameters for click", log.ErrVal(err), slog.Any("step", step))
	// 	return nil, err
	// } else {
	// 	r.params = *params
	// }

	return r, nil
}
