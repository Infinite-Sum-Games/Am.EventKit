package cmd

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/rs/zerolog"
)

func InitLogger(env string) (*pkg.LoggerService, error) {
	var output io.Writer

	switch env {
	case "DEVELOPMENT":
		file, err := os.OpenFile("logs/dev.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		fileWriter := zerolog.ConsoleWriter{
			Out:        file,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		}
		output = zerolog.MultiLevelWriter(consoleWriter, fileWriter)
	case "PRODUCTION":
		file, err := os.OpenFile("logs/prod.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		output = file
	default:
		return nil, errors.New("invalid environment for logger setup")
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	return &pkg.LoggerService{
		Logger: logger,
		Env:    env,
	}, nil
}
