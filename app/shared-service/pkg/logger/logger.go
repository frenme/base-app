package logger

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger = logrus.Logger

type ctxKey string

const keyReqID ctxKey = "req_id"

func New() *Logger {
	env := os.Getenv("MODE_ENV")
	config := GetLoggerConfig(env)

	return createLoggerFromConfig(config)
}

func NewForEnv(env string) *Logger {
	config := GetLoggerConfig(env)
	return createLoggerFromConfig(config)
}

func createLoggerFromConfig(config LoggerConfig) *Logger {
	l := logrus.New()
	l.SetLevel(config.Level)
	l.SetFormatter(config.Formatter)
	l.SetOutput(config.Output)
	l.SetReportCaller(config.ReportCaller)

	return l
}

func WithReqID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyReqID, id)
}

func From(ctx context.Context, base *Logger) *logrus.Entry {
	e := logrus.NewEntry(base)
	if v := ctx.Value(keyReqID); v != nil {
		e = e.WithField("req_id", v)
	}
	return e
}

func callerPretty(f *runtime.Frame) (function string, file string) {
	fn := f.Function
	if i := strings.LastIndex(fn, "/"); i >= 0 {
		fn = fn[i+1:]
	}
	file = filepath.Base(f.File) + ":" + strconv.Itoa(f.Line)
	return fn, file
}
