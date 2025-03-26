package endpoints

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/engine"
	"github.com/fmotalleb/scrapper-go/log"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
)

func init() {
	registerEndpoint(
		endpoint{
			method:  "POST",
			path:    "/process",
			handler: processPipeline,
		},
	)
}

func processPipeline(c echo.Context) error {
	cfgMap := make(map[string]any)

	err := json.NewDecoder(c.Request().Body).Decode(&cfgMap)
	if err != nil {
		slog.Error("failed to body", log.ErrVal(err))
		return c.String(http.StatusBadRequest, "cannot unmarshal the given json body")
	}
	var cfg config.ExecutionConfig
	err = mapstructure.Decode(cfgMap, &cfg)
	if err != nil {
		slog.Error("failed to read config from body", log.ErrVal(err))
		return c.String(http.StatusBadRequest, "cannot unmarshal the given json body")
	}
	res, err := engine.ExecuteConfig(c.Request().Context(), cfg)
	if err != nil {
		slog.Error("failed to execute config", log.ErrVal(err))
		return c.String(http.StatusBadRequest, "failed to execute config. make sure the config is compatible with service")
	}
	return c.JSON(http.StatusOK, res)
}
