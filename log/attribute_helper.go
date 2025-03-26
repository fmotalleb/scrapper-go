package log

import "log/slog"

func Error(msg string, err error, attrs ...any) {
	slog.Error(msg,
		append(attrs, ErrVal(err))...,
	)
}

func ErrVal(err error) slog.Attr {
	return slog.Any("error", err)
}
