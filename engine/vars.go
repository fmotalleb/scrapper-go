package engine

import (
	"fmt"
	"log/slog"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
)

func initializeVariables(varsConfig []config.Variable) (utils.Vars, error) {
	vars := make(utils.Vars)

	for _, v := range varsConfig {
		// Log the variable being processed
		slog.Debug("Processing variable", slog.Any("var", v))

		switch {
		case v.Random == "once":
			value := v.Prefix + utils.RandomString(v.RandomChars, v.RandomLength) + v.Postfix
			vars.SetOnce(v.Name, value)
			slog.Debug("Set variable once", slog.Any("name", v.Name), slog.Any("value", value))
		case v.Random == "always":
			vars.SetGetter(
				v.Name,
				func() string {
					randomValue := v.Prefix + utils.RandomString(v.RandomChars, v.RandomLength) + v.Postfix
					slog.Debug("Generated random value", slog.Any("name", v.Name), slog.Any("random_value", randomValue))
					return randomValue
				},
			)
			slog.Debug("Set variable getter", slog.Any("name", v.Name))
		case v.Value != "":
			finalValue := v.Prefix + v.Value + v.Postfix
			vars.SetOnce(v.Name, finalValue)
			slog.Debug("Set variable once", slog.Any("name", v.Name), slog.Any("value", finalValue))
		default:
			slog.Error("Unknown variable type", slog.Any("variable", v))
			return nil, fmt.Errorf("unknown variable type: %v", v)
		}
	}

	slog.Info("Variables initialization completed")
	return vars, nil
}
