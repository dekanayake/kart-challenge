package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitLogger() {
	var level zerolog.Level
	switch AppConfig.LogLevel {
	case "debug":
		level = zerolog.DebugLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	case "fatal":
		level = zerolog.FatalLevel
	default:
		level = zerolog.InfoLevel
	}

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	Logger = zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Str("service", "kart-service").
		Logger()
}
