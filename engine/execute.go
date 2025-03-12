package engine

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/query"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/playwright-community/playwright-go"
	"golang.org/x/exp/slog"
)

type executionEngine func(playwright.Page, config.Step, Vars, map[string]any) error

var executors = map[string]executionEngine{
	"nop":        nop,
	"sleep":      sleep,
	"select":     selectInput,
	"fill":       fillInput,
	"click":      click,
	"exec":       executeJs,
	"print":      elementSelector,
	"element":    elementSelector,
	"table":      table,
	"goto":       gotoPage,
	"screenshot": screenshot,
}

func executeStep(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	ok, err := evaluateExpression(page, step, vars)
	if err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("condition didn't pass for step: %v", step)
	}

	for key, executor := range executors {
		if step[key] != nil {
			return executor(page, step, vars, result)
		}
	}
	return fmt.Errorf("unknown step action: %v", step)
}

func evaluateExpression(page playwright.Page, step config.Step, vars Vars) (bool, error) {
	cond, ok := step["if"].(string)
	cond, err := execTemplate(cond, vars, page)
	if err != nil {
		return false, err
	}
	if !ok {
		return true, nil
	}
	query, err := query.ParseQuery(cond)
	if err != nil {
		return false, err
	}
	return query.EvaluateQuery(vars.Snapshot())
}

func sleep(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	waitTime := step["sleep"].(string)
	waitTime, err := execTemplate(waitTime, vars, page)
	if err != nil {
		return err
	}
	value, err := time.ParseDuration(waitTime)
	if err != nil {
		return nil
	}
	time.Sleep(value)
	return nil
}

func selectInput(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	selector := step["select"].(string)
	selector, err := execTemplate(selector, vars, page)
	if err != nil {
		return err
	}
	value := ""
	if step["var"] != nil {
		var err error
		value, err = vars.GetOrFail(step["var"].(string))
		if err != nil {
			return err
		}
	} else if step["value"] != nil {
		value = step["value"].(string)
	}
	value, err = execTemplate(value, vars, page)
	if err != nil {
		return err
	}

	if _, err := page.Locator(selector).SelectOption(playwright.SelectOptionValues{
		Values: &[]string{value},
	}); err != nil {
		return err
	}
	return nil
}

func fillInput(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	selector := step["fill"].(string)
	selector, err := execTemplate(selector, vars, page)
	if err != nil {
		return err
	}
	value := ""
	if step["var"] != nil {
		var err error
		value, err = vars.GetOrFail(step["var"].(string))
		if err != nil {
			return err
		}
	} else if step["value"] != nil {
		value = step["value"].(string)
	}
	value, err = execTemplate(value, vars, page)
	if err != nil {
		return err
	}
	return page.Locator(selector).Fill(value)
}

func click(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	selector := step["click"].(string)
	selector, err := execTemplate(selector, vars, page)
	if err != nil {
		return err
	}
	return page.Locator(selector).Click()
}

func executeJs(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	script := step["exec"].(string)

	script, err := execTemplate(script, vars, page)
	if err != nil {
		return err
	}
	value, err := page.Evaluate(script)
	if err != nil {
		return err
	}
	return setVar(step, value, vars, result)
}

func gotoPage(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	url := step["goto"].(string)
	_, err := page.Goto(url)
	return err
}

func elementSelector(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	selector := step["element"].(string)
	locator := page.Locator(selector)
	slog.Warn("element", slog.AnyValue(locator))
	isInput, _ := step["is-input"].(bool)
	asHTML, _ := step["as-html"].(bool)
	var value string
	var err error
	switch {
	case isInput:
		value, err = locator.InputValue()
	case asHTML:
		value, err = locator.InnerHTML()
	default:
		value, err = locator.TextContent()
	}

	if err != nil {
		return err
	}

	return setVar(step, value, vars, result)
}

func table(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	selector := step["table"].(string)
	locator := page.Locator(selector)

	value, err := locator.InnerHTML()
	if err != nil {
		return err
	}
	table := fmt.Sprintf("<table> %s </table>", value)
	if step["parse-mode"] == nil {
		step["parse-mode"] = "table"
	}
	return setVar(step, table, vars, result)
}

func setVar(step config.Step, value interface{}, vars Vars, result map[string]any) error {
	if setVar, ok := step["set-var"].(string); ok {
		parseMode, setVarIsSet := step["parse-mode"].(string)
		if !setVarIsSet {
			parseMode = "text"
		}
		val := value.(string)
		vars.SetOnce(setVar, val)
		switch parseMode {
		case "text":
			result[setVar] = val
		case "table":
			table, err := utils.ParseTable(val)
			if err != nil {
				slog.Error("failed to parse table", slog.Any("err", err))
				return err
			}
			result[setVar] = table
		case "json":
			var jsonVal map[string]interface{}
			if err := json.Unmarshal([]byte(val), &jsonVal); err != nil {
				slog.Error("error parsing JSON", slog.Any("err", err))
				return err
			}
			result[setVar] = jsonVal
		}
	} else {
		fmt.Println(value)
	}
	return nil
}

func nop(p playwright.Page, s config.Step, v Vars, r map[string]any) error { return nil }

func screenshot(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	selector := step["screenshot"].(string)

	lso, err := readParams[playwright.LocatorScreenshotOptions](step)
	if err != nil {
		return err
	}

	if lso.Path == nil {
		loc := "./screenshot.png"
		lso.Path = &loc
	}

	locator := page.Locator(selector)

	_, err = locator.Screenshot(*lso)
	return err
}

func readParams[T any](step config.Step) (*T, error) {
	params, _ := step["params"].(map[string]any)
	if params == nil {
		params = make(map[string]any)
	}
	var item T
	if err := mapstructure.Decode(params, &item); err != nil {
		return nil, err
	}
	return &item, nil
}
