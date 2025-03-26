// Package utils contains utilities functions
package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/log"
	"gopkg.in/yaml.v3"
)

type Output string

const (
	yamlFmt = Output("yaml")
	jsonFmt = Output("json")
)

func (f Output) Format(data map[string]any) (string, error) {
	switch f {
	case yamlFmt:
		result, err := yaml.Marshal(data)
		if err != nil {
			slog.Error("Failed to marshal YAML", log.ErrVal(err))
			return "", err
		}
		slog.Debug("Successfully formatted data as YAML")
		return string(result), nil

	case jsonFmt:
		result, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			slog.Error("Failed to marshal JSON", log.ErrVal(err))
			return "", err
		}
		slog.Debug("Successfully formatted data as JSON")
		return string(result), nil

	default:
		slog.Warn("Unsupported format requested", slog.Any("format", f))
		return "", fmt.Errorf("unsupported format: %s", f)
	}
}
