package engine

import (
	"fmt"
	"time"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/playwright-community/playwright-go"
)

type executionEngine func(playwright.Page, config.Step, Vars, map[string]any) error

var executors = map[string]executionEngine{
	"nop":     nop,
	"sleep":   sleep,
	"select":  selectInput,
	"fill":    fillInput,
	"click":   click,
	"exec":    executeJs,
	"print":   elementSelector,
	"element": elementSelector,
	"goto":    gotoPage,
}

func executeStep(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	ok, err := evaluateExpression(step, vars)
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

func evaluateExpression(step config.Step, vars Vars) (bool, error) {
	cond, ok := step["if"].(string)
	if !ok {
		return true, nil
	}
	query, err := ParseQuery(cond)
	if err != nil {
		return false, err
	}
	return query.EvaluateQuery(vars.Snapshot())
}

func sleep(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	waitTime := step["sleep"].(string)
	value, err := time.ParseDuration(waitTime)
	if err != nil {
		return nil
	}
	time.Sleep(value)
	return nil
}

func selectInput(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	selector := step["select"].(string)
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
	if _, err := page.Locator(selector).SelectOption(playwright.SelectOptionValues{
		Values: &[]string{value},
	}); err != nil {
		return err
	}
	return nil
}

func fillInput(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	selector := step["fill"].(string)
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
	return page.Locator(selector).Fill(value)
}

func click(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	selector := step["click"].(string)
	return page.Locator(selector).Click()
}

func executeJs(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	script := step["exec"].(string)
	value, err := page.Evaluate(script)
	if err != nil {
		return err
	}
	setVar(step, value, vars, result)
	return nil
}

func gotoPage(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	url := step["goto"].(string)
	_, err := page.Goto(url)
	return err
}

func elementSelector(page playwright.Page, step config.Step, vars Vars, result map[string]any) error {
	selector := step["element"].(string)
	locator := page.Locator(selector)
	isInput, _ := step["is-input"].(bool)
	var value string
	var err error
	if isInput {
		value, err = locator.InputValue()
	} else {
		value, err = locator.TextContent()
	}
	if err != nil {
		return err
	}
	setVar(step, value, vars, result)
	return nil
}

func setVar(step config.Step, value interface{}, vars Vars, result map[string]any) {
	if setVar, ok := step["set-var"].(string); ok {
		val := value.(string)
		vars.SetOnce(setVar, val)
		result[setVar] = val
	} else {
		fmt.Println(value)
	}
}

func nop(p playwright.Page, s config.Step, v Vars, r map[string]any) error { return nil }
