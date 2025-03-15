package server

import (
	"log/slog"

	"github.com/fmotalleb/scrapper-go/server/endpoints"
	"github.com/labstack/echo/v4"
)

func StartServer(address string) error {
	e := echo.New()
	endpoints.PopulateEndpoints(e)
	if err := e.Start(address); err != nil {
		slog.Error("failed to start server", slog.Any("err", err))
		return err
	}
	return nil
}
