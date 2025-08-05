package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"log/slog"
)

type SimpleHandler struct {
	out io.Writer
}

func NewSimpleHandler(w io.Writer) slog.Handler {
	return &SimpleHandler{out: w}
}

func (h *SimpleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *SimpleHandler) Handle(ctx context.Context, r slog.Record) error {
	var caller string
	if _, file, line, ok := runtime.Caller(4); ok {
		caller = fmt.Sprintf("%s:%d", file, line)
	}

	var parts []string
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == slog.TimeKey || a.Key == slog.LevelKey || a.Key == slog.MessageKey {
			return true
		}
		parts = append(parts, fmt.Sprintf("%s=%v", a.Key, a.Value))
		return true
	})

	ts := r.Time.Format(time.RFC3339Nano)
	lvl := strings.ToUpper(r.Level.String())
	msg := r.Message

	line := fmt.Sprintf("%s %s: %s", ts, lvl, msg)
	if caller != "" {
		line += " caller=" + caller
	}
	if len(parts) > 0 {
		line += " " + strings.Join(parts, " ")
	}

	_, err := fmt.Fprintln(h.out, line)
	return err
}

func (h *SimpleHandler) WithAttrs(as []slog.Attr) slog.Handler {
	return h
}

func (h *SimpleHandler) WithGroup(name string) slog.Handler {
	return h
}

type Logger struct {
	handler slog.Handler
}

func CreateLogger() *Logger {
	base := NewSimpleHandler(os.Stdout)
	return &Logger{handler: base}
}

func (l *Logger) raw(level slog.Level, msg string) {
	rec := slog.NewRecord(time.Now(), level, msg, 0)
	_ = l.handler.Handle(context.Background(), rec)
}

func (l *Logger) Info(args ...interface{}) {
	l.raw(slog.LevelInfo, fmt.Sprint(args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.raw(slog.LevelInfo, fmt.Sprintf(format, args...))
}

func (l *Logger) Warn(args ...interface{}) {
	l.raw(slog.LevelWarn, fmt.Sprint(args...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.raw(slog.LevelWarn, fmt.Sprintf(format, args...))
}

func (l *Logger) Error(args ...interface{}) {
	l.raw(slog.LevelError, fmt.Sprint(args...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.raw(slog.LevelError, fmt.Sprintf(format, args...))
}
