package endpoints

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/session"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
)

func init() {
	registerEndpoint(
		endpoint{
			method:  "POST",
			path:    "/sessions/:id",
			handler: sessionHandle,
		},
	)
}
func sessionHandle(c echo.Context) error {
	id := c.Param("id")
	slog.Info("session handle requested", slog.String("id", id))
	sess, ok := session.GetSession(id)
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]any{"id": id})
	}
	var cfgInterface any
	err := json.NewDecoder(c.Request().Body).Decode(&cfgInterface)
	if err != nil {
		slog.Error("failed to parse body", slog.Any("err", err))
		return c.String(http.StatusBadRequest, "cannot unmarshal the given json body")
	}
	switch v := cfgInterface.(type) {
	case map[string]any:
		var cfg config.Step
		err = mapstructure.Decode(v, &cfg)
		if err != nil {
			slog.Error("failed to decode config", slog.Any("err", err))
			return c.String(http.StatusBadRequest, "invalid config format")
		}
		res, err := sess.Handle(cfg)
		if err != nil {
			slog.Error("failed to execute config", slog.Any("err", err))
			return c.String(http.StatusBadRequest, "failed to execute config")
		}
		return c.JSON(http.StatusOK, res)

	case []any:
		var steps []config.Step
		err = mapstructure.Decode(v, &steps)
		if err != nil {
			slog.Error("failed to decode config array", slog.Any("err", err))
			return c.String(http.StatusBadRequest, "invalid config array format")
		}
		res, err := sess.Handle(steps...)
		if err != nil {
			slog.Error("failed to execute config steps", slog.Any("err", err))
			return c.String(http.StatusBadRequest, "failed to execute config steps")
		}
		return c.JSON(http.StatusOK, res)

	default:
		return c.String(http.StatusBadRequest, "invalid JSON format, expected object or array")
	}
}
