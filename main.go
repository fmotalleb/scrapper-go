package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/playwright-community/playwright-go"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config config.ExecutionConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	vars := generateVariables(config.Pipeline.Vars)

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start Playwright: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(config.Pipeline.BrowserParams)
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	for _, step := range config.Pipeline.Steps {
		if err := executeStep(page, step, vars); err != nil {
			log.Fatalf("Error executing step: %v, %v", err, step)
		}
	}
}

func generateVariables(varsConfig []config.Variable) map[string]func() string {
	vars := make(map[string]func() string)
	for _, v := range varsConfig {
		switch {
		case v.Random == "once":
			value := v.Prefix + randomString(v.RandomChars, v.RandomLength) + v.Postfix
			vars[v.Name] = func() string {
				return value
			}
		case v.Random == "always":
			vars[v.Name] = func() string {
				return v.Prefix + randomString(v.RandomChars, v.RandomLength) + v.Postfix
			}
		case v.Value != "":
			vars[v.Name] = func() string {
				return v.Prefix + v.Value + v.Postfix
			}
		default:
			log.Fatalf("unknown variable type: %v", v)
		}

	}
	return vars
}

func randomString(charset string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

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
		value, err := page.Locator(selector).InputValue()
		if err != nil {
			return err
		}
		fmt.Println(value)
		return nil
	default:
		return fmt.Errorf("unknown step action: %v", step)
	}
}
