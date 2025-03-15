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
		if v.Name == "" {
			slog.Warn("skipping variable with empty name", slog.Any("variable", v))
			continue
		}

		var value string
		switch v.Random {
		case "once":
			value = v.Prefix + utils.RandomString(v.RandomChars, v.RandomLength) + v.Postfix
			vars.SetOnce(v.Name, value)
			slog.Debug("initialized 'once' random variable", slog.String("name", v.Name), slog.String("value", value))

		case "always":
			vars.SetGetter(v.Name, func() string {
				return v.Prefix + utils.RandomString(v.RandomChars, v.RandomLength) + v.Postfix
			})
			slog.Debug("initialized 'always' random variable getter", slog.String("name", v.Name))

		default:
			if v.Value != "" {
				value = v.Prefix + v.Value + v.Postfix
				vars.SetOnce(v.Name, value)
				slog.Debug("initialized fixed variable", slog.String("name", v.Name), slog.String("value", value))
			} else {
				slog.Error("unknown variable configuration", slog.Any("variable", v))
				return nil, fmt.Errorf("invalid variable configuration: %v", v)
			}
		}
	}

	slog.Info("variables initialization completed", slog.Int("total", len(varsConfig)))
	return vars, nil
}
