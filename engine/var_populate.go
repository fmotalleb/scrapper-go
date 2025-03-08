package engine

import (
	"log"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/utils"
)

func generateVariables(varsConfig []config.Variable) map[string]func() string {
	vars := make(map[string]func() string)
	for _, v := range varsConfig {
		switch {
		case v.Random == "once":
			value := v.Prefix + utils.RandomString(v.RandomChars, v.RandomLength) + v.Postfix
			vars[v.Name] = func() string {
				return value
			}
		case v.Random == "always":
			vars[v.Name] = func() string {
				return v.Prefix + utils.RandomString(v.RandomChars, v.RandomLength) + v.Postfix
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
