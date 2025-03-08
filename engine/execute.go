package engine

import (
	"fmt"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/playwright-community/playwright-go"
)

func executeStep(page playwright.Page, step config.Step, vars map[string]func() string) error {
	switch {
	case step["goto"] != nil:
		url := step["goto"].(string)
		_, err := page.Goto(url)
		return err
	case step["click"] != nil:
		selector := step["click"].(string)
		return page.Locator(selector).Click()
	case step["fill"] != nil:
		selector := step["fill"].(string)
		value := ""
		if step["var"] != nil {
			value = vars[step["var"].(string)]()
		} else if step["value"] != nil {
			value = step["value"].(string)
		}
		return page.Locator(selector).Fill(value)
	case step["select"] != nil:
		selector := step["select"].(string)
		value := ""
		if step["var"] != nil {
			value = vars[step["var"].(string)]()
		} else if step["value"] != nil {
			value = step["value"].(string)
		}
		if _, err := page.Locator(selector).SelectOption(playwright.SelectOptionValues{
			Values: &[]string{value},
		}); err != nil {
			return err
		}
		return nil
	case step["print"] != nil:
		selector := step["print"].(string)
		value, err := page.Locator(selector).TextContent()
		if err != nil {
			return err
		}
		fmt.Println(value)
		return nil
	default:
		return fmt.Errorf("unknown step action: %v", step)
	}
}
