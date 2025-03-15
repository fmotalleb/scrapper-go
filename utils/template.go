package utils

import (
	"bytes"
	"fmt"
	"log/slog"
	"text/template"

	"github.com/playwright-community/playwright-go"
)

func unShadow(data map[string]any, key string) map[string]any {
	if _, exists := data[key]; exists {
		newKey := fmt.Sprintf("_%s", key)
		data = unShadow(data, newKey)
		data[newKey] = data[key]
	} else {
	}
	return data
}

func EvaluateTemplates(texts []string, vars Vars, page playwright.Page) ([]string, error) {
	return MapItems(texts, TemplateEvalMapper(vars, page))
}

func TemplateEvalMapper(vars Vars, page playwright.Page) func(string) (string, error) {
	return func(s string) (string, error) {
		return EvaluateTemplate(s, vars, page)
	}
}

func EvaluateTemplate(text string, vars Vars, page playwright.Page) (string, error) {
	templateObj := template.New("template")

	variables := vars.LiveSnapshot()
	if _, ok := variables["eval"]; ok {
		slog.Error("found a page variable in live snapshot generated for template, renaming old page variable to _page")
		variables = unShadow(variables, "eval")
	}
	variables["eval"] = page.Evaluate
	templateObj = templateObj.Funcs(variables)

	templateObj, err := templateObj.Parse(text)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %s", err)
	}

	variables["page"] = page
	output := bytes.NewBufferString("")
	err = templateObj.Execute(output, variables)
	if err != nil {
		return "", fmt.Errorf("failed to execute template using vars snapshot: %s", err)
	}
	return output.String(), nil
}
