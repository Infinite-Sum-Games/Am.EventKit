package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/rs/zerolog"
)

func InitLogger(env string) *pkg.LoggerService {
	var output io.Writer
	
	switch env {
		case "development", "DEV", "dev" :
			file, err := os.OpenFile("logs/dev.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				panic("could not create log file: " + err.Error())
			}
			consoleWriter := zerolog.ConsoleWriter{
				Out: os.Stdout,
				TimeFormat: time.RFC3339,
			}
			fileWriter := zerolog.ConsoleWriter{
				Out: file,
				TimeFormat: time.RFC3339,
				NoColor: true,
			}
			output = zerolog.MultiLevelWriter(consoleWriter, fileWriter)
		case "production", "prod", "PROD":
			file, err := os.OpenFile("logs/prod.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err!=nil {
				panic("could not create log file: "+err.Error())
			}
			output=file
		default:
			panic(errors.New("invalid environment for logger setup"))
	}
	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(output).With().Timestamp().Logger()

	fmt.Print("hello")

	return &pkg.LoggerService{
		Logger: logger,
		Env: env,
	}
}