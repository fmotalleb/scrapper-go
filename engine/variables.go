package engine

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
)

type (
	getter = func() string
	Vars   map[string]varValue
)

type varValue struct {
	isGenerative bool
	value        string
	get          getter
}

func (v *varValue) getValue() string {
	if v.isGenerative {
		return v.get()
	} else {
		return v.value
	}
}

func (v Vars) Snapshot() map[string]string {
	data := make(map[string]string)
	for k, g := range v {
		data[k] = g.getValue()
	}
	return data
}

func (v Vars) LiveSnapshot() map[string]interface{} {
	data := make(map[string]interface{})
	for k, g := range v {
		if g.isGenerative {
			data[k] = g.get
		} else {
			data[k] = g.value
		}

	}
	return data
}

func (v Vars) SetOnce(key string, value string) {
	v[key] = varValue{
		isGenerative: false,
		value:        value,
	}
}

func (v Vars) SetGetter(key string, getter getter) {
	v[key] = varValue{
		isGenerative: true,
		get:          getter,
	}
}

func (v Vars) Get(key string) (string, bool) {
	item, ok := v[key]

	return item.getValue(), ok
}

func (v Vars) GetOr(key string, def string) string {
	value, ok := v.Get(key)
	if ok {
		return value
	}
	return def
}

func (v Vars) GetOrFail(key string) (string, error) {
	value, ok := v.Get(key)
	if ok {
		return value, nil
	}
	return "", fmt.Errorf("use of undefined variable: %s", key)
}

func initializeVariables(varsConfig []config.Variable) Vars {
	vars := make(Vars)

	for _, v := range varsConfig {
		// Log the variable being processed
		slog.Debug("Processing variable", slog.Any("variable_name", v.Name), slog.Any("random", v.Random), slog.Any("value", v.Value))

		switch {
		case v.Random == "once":
			value := v.Prefix + utils.RandomString(v.RandomChars, v.RandomLength) + v.Postfix
			vars.SetOnce(v.Name, value)
			slog.Info("Set variable once", slog.Any("name", v.Name), slog.Any("value", value))

		case v.Random == "always":
			vars.SetGetter(
				v.Name,
				func() string {
					randomValue := v.Prefix + utils.RandomString(v.RandomChars, v.RandomLength) + v.Postfix
					slog.Debug("Generated random value", slog.Any("name", v.Name), slog.Any("random_value", randomValue))
					return randomValue
				},
			)
			slog.Info("Set variable getter", slog.Any("name", v.Name))

		case v.Value != "":
			finalValue := v.Prefix + v.Value + v.Postfix
			vars.SetOnce(v.Name, finalValue)
			slog.Info("Set variable once", slog.Any("name", v.Name), slog.Any("value", finalValue))

		default:
			slog.Error("Unknown variable type", slog.Any("variable", v))
			log.Fatalf("unknown variable type: %v", v)
		}
	}

	slog.Info("Variables initialization completed")
	return vars
}
