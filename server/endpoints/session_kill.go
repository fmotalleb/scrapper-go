package endpoints

import (
	"net/http"

	"github.com/fmotalleb/scrapper-go/session"
	"github.com/labstack/echo/v4"
)

func init() {
	registerEndpoint(
		endpoint{
			method:  "DELETE",
			path:    "/session/:id",
			handler: sessionKill,
		},
	)
}

func sessionKill(c echo.Context) error {
	id := c.Param("id")
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
