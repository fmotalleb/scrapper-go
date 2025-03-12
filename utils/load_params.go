package utils

import (
	"github.com/fmotalleb/scrapper-go/config"
	"github.com/mitchellh/mapstructure"
)

func LoadParams[T any](step config.Step) (*T, error) {
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
