package endpoints

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/engine"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
)

func init() {
	registerEndpoint(
		endpoint{
			method:  "GET",
			path:    "/live_stream",
			handler: liveStream,
		},
	)
}

var wsUpgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func liveStream(c echo.Context) error {
	sendChan := make(chan map[string]any)
	recvChan := handleWebSocket(c, sendChan)

	cfgMap := <-recvChan

	var cfg config.ExecutionConfig
	err := mapstructure.Decode(cfgMap, &cfg)
	if err != nil {
		slog.Error("failed to read config from body", slog.Any("err", err))
		c.String(http.StatusBadRequest, "cannot unmarshal the given json body")
		return err
	}
	pipe := make(chan []config.Step)
	resultChan, err := engine.ExecuteStream(c.Request().Context(), cfg, pipe)
	go func() {
		for i := range resultChan {
			sendChan <- i
		}
	}()
	for i := range recvChan {
		var cfg config.Step
		err = mapstructure.Decode(i, &cfg)
		if err != nil {
			slog.Error("failed to read config from body", slog.Any("err", err))
			sendChan <- map[string]any{
				"error": err.Error(),
			}
		} else {
			pipe <- []config.Step{cfg}
		}
	}

	// if err != nil {
	// 	slog.Error("failed to execute config", slog.Any("err", err))
	// 	c.String(http.StatusBadRequest, "failed to execute config. make sure the config is compatible with service")
	// 	return err
	// }
	// c.Request().Body.Read()
	// c.Response().Write()
	return nil
}

func handleWebSocket(c echo.Context, sendChan chan map[string]any) chan map[string]any {
	conn, err := wsUpgrade.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return nil
	}

	recvChan := make(chan map[string]any)

	// Goroutine to read messages from the client
	go func() {
		defer conn.Close()
		for {
			var msg map[string]any
			if err := conn.ReadJSON(&msg); err != nil {
				slog.Error("read error:", slog.Any("error", err))
				close(recvChan)
				return
			}
			recvChan <- msg
		}
	}()

	// Goroutine to send messages to the client
	go func() {
		defer conn.Close()
		for msg := range sendChan {
			if err := conn.WriteJSON(msg); err != nil {
				slog.Error("write error:", slog.Any("error", err))
				return
			}
		}
	}()

	return recvChan
}
