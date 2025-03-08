package engine

import (
	"fmt"
	"time"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/playwright-community/playwright-go"
)

type executionEngine func(config.Step, Vars, playwright.Page) error

var executors = map[string]executionEngine{
	"nop":     func(s config.Step, v Vars, p playwright.Page) error { return nil },
	"sleep":   sleep,
	"select":  selectInput,
	"fill":    fillInput,
	"click":   click,
	"exec":    executeJs,
	"print":   elementSelector,
	"element": elementSelector,
	"goto":    gotoPage,
}

func executeStep(page playwright.Page, step config.Step, vars Vars) error {
	ok, err := evaluateExpression(step, vars)
	if err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("condition didn't pass for step: %v", step)
	}

	for key, executor := range executors {
		if step[key] != nil {
			return executor(step, vars, page)
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

func sleep(step config.Step, vars Vars, page playwright.Page) error {
	waitTime := step["sleep"].(string)
	value, err := time.ParseDuration(waitTime)
	if err != nil {
		return nil
	}
	time.Sleep(value)
	return nil
}

func selectInput(step config.Step, vars Vars, page playwright.Page) error {
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

func fillInput(step config.Step, vars Vars, page playwright.Page) error {
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

func click(step config.Step, vars Vars, page playwright.Page) error {
	selector := step["click"].(string)
	return page.Locator(selector).Click()
}

func executeJs(step config.Step, vars Vars, page playwright.Page) error {
	script := step["exec"].(string)
	value, err := page.Evaluate(script)
	if err != nil {
		return err
	}
	if setVar, ok := step["set-var"].(string); ok {
		vars.SetOnce(setVar, value.(string))
	} else {
		fmt.Println(value)
	}
	return nil
}

func gotoPage(step config.Step, vars Vars, page playwright.Page) error {
	url := step["goto"].(string)
	_, err := page.Goto(url)
	return err
}

func elementSelector(step config.Step, vars Vars, page playwright.Page) error {
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
	if setVar, ok := step["set-var"].(string); ok {
		vars.SetOnce(setVar, value)
	} else {
		fmt.Println(value)
	}
	return nil

}
