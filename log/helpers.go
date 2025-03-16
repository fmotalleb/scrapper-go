// Package log contains utilities for logging system
package log

import (
	"io"
	"log/slog"
	"os"
)

func SetupLogger(level string, w ...io.Writer) error {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		return err
	}
	var writer io.Writer = os.Stderr
	if len(w) != 0 {
		writer = w[0]
	}
	opts := &slog.HandlerOptions{}
	opts.Level = lvl
	newRoot := slog.New(slog.NewTextHandler(writer, opts))
	slog.SetDefault(newRoot)
	return nil
}
