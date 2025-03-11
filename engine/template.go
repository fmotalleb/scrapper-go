package engine

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
func applyTemplate(text string, vars Vars, page playwright.Page) (string, error) {

	tmpl := template.New("template")

	tmpl = tmpl.Funcs(map[string]any{
		"eval": page.Evaluate,
	})

	tmpl, err := tmpl.Parse(text)

	if err != nil {
		return "", fmt.Errorf("failed to parse template: %s", err)
	}
	variables := vars.LiveSnapshot()
	if _, ok := variables["page"]; ok {
		slog.Error("found a page variable in live snapshot generated for template, renaming old page variable to _page")
		variables = unShadow(variables, "page")
	}
	variables["page"] = page
	output := bytes.NewBufferString("")
	err = tmpl.Execute(output, variables)
	if err != nil {
		return "", fmt.Errorf("failed to execute template using vars snapshot: %s", err)
	}
	return output.String(), nil
}
