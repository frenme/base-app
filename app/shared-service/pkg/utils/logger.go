package utils

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
)

type LoggerHandler struct {
	Handler slog.Handler
}

func (ch LoggerHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return ch.Handler.Enabled(ctx, level)
}

func (ch LoggerHandler) Handle(ctx context.Context, r slog.Record) error {
	_, file, line, ok := runtime.Caller(3)
	if ok {
		r.AddAttrs(slog.String("caller", fmt.Sprintf("%s:%d", file, line)))
	}
	return ch.Handler.Handle(ctx, r)
}

func (ch LoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return LoggerHandler{Handler: ch.Handler.WithAttrs(attrs)}
}

func (ch LoggerHandler) WithGroup(name string) slog.Handler {
	return LoggerHandler{Handler: ch.Handler.WithGroup(name)}
}
