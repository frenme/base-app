package logger

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggerConfig struct {
	Level        logrus.Level
	Formatter    logrus.Formatter
	Output       io.Writer
	ReportCaller bool
}

var (
	textFormatter = &logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        time.RFC3339Nano,
		DisableLevelTruncation: true,
		PadLevelText:           true,
		CallerPrettyfier:       callerPretty,
	}

	jsonFormatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	}
)

var environments = map[string]LoggerConfig{
	"dev": {
		Level:        logrus.DebugLevel,
		Formatter:    textFormatter,
		Output:       os.Stdout,
		ReportCaller: true,
	},
	"prod": {
		Level:        logrus.InfoLevel,
		Formatter:    jsonFormatter,
		Output:       os.Stdout,
		ReportCaller: false,
	},
}

var defaultConfig = LoggerConfig{
	Level:        logrus.DebugLevel,
	Formatter:    textFormatter,
	Output:       os.Stdout,
	ReportCaller: true,
}

func GetLoggerConfig(env string) LoggerConfig {
	if env == "" {
		return defaultConfig
	}

	if config, exists := environments[env]; exists {
		return config
	}

	return defaultConfig
}
