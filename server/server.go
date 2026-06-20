// Package server contains api server logic
package server

import (
	"log/slog"

	"github.com/labstack/echo/v4"

	"github.com/fmotalleb/scrapper-go/log"
	"github.com/fmotalleb/scrapper-go/server/endpoints"
)

func StartServer(address string) error {
	e := echo.New()
	endpoints.PopulateEndpoints(e)
	if err := e.Start(address); err != nil {
		slog.Error("failed to start server", log.ErrVal(err))
		return err
	}
	return nil
}
