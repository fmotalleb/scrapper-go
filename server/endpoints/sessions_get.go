package endpoints

import (
	"net/http"

	"github.com/fmotalleb/scrapper-go/session"
	"github.com/labstack/echo/v4"
)

func init() {
	registerEndpoint(
		endpoint{
			method:  "GET",
			path:    "/sessions",
			handler: sessionsGet,
		},
	)
}

func sessionsGet(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		session.GetSessions(),
	)
}
