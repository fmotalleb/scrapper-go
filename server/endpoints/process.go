package endpoints

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/engine"
	"github.com/labstack/echo/v4"
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
	slog.Info("test")
	var cfg config.ExecutionConfig
	err := json.NewDecoder(c.Request().Body).Decode(&cfg)
	if err != nil {
		slog.Error("failed to read config from body", slog.Any("err", err))
		c.String(http.StatusBadRequest, "cannot unmarshal the given json body")
		return err
	}
	res, err := engine.ExecuteConfig(cfg)
	if err != nil {
		slog.Error("failed to execute config", slog.Any("err", err))
		c.String(http.StatusBadRequest, "failed to execute config. make sure the config is compatible with service")
		return err
	}
	return c.JSON(http.StatusOK, res)
}
