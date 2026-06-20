package endpoints

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/fmotalleb/scrapper-go/session"
)

func init() {
	registerEndpoint(
		endpoint{
			method:  "DELETE",
			path:    "/sessions/:id",
			handler: sessionKill,
		},
	)
}

func sessionKill(c echo.Context) error {
	id := c.Param("id")
	slog.Info("session kill requested", slog.String("id", id))

	if session, ok := session.GetSession(id); ok {
		session.Kill()
		return c.JSON(http.StatusOK, map[string]any{
			"id": id,
		})
	}
	return c.JSON(http.StatusNotFound, map[string]any{
		"id": id,
	})
}
