package shell

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/engine"
	"github.com/mitchellh/mapstructure"
)

func bindToBrowser(ctx context.Context, recvChan <-chan map[string]any) <-chan map[string]any {

	cfgMap := <-recvChan

	var cfg config.ExecutionConfig
	err := mapstructure.Decode(cfgMap, &cfg)
	if err != nil {
		slog.Error("failed to read config from body", slog.Any("err", err))
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}
	sendChan := make(chan map[string]any)
	pipe := make(chan []config.Step)
	resultChan, err := engine.ExecuteStream(ctx, cfg, pipe)
	go func() {
		for i := range resultChan {
			sendChan <- i
		}
	}()
	go func() {
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
	}()
	return sendChan
}
