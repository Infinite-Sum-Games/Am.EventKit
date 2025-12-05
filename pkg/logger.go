package pkg

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var Log *LoggerService

type LoggerService struct {
	Logger zerolog.Logger
	Env    string
}

type ILoggerService interface {

	// Function to enrich each log with data
	enrich(c *gin.Context, e *zerolog.Event) *zerolog.Event

	// This set of functions is to be used in the context of the web-server
	// where there is a gin.Context (server context) involved
	DebugCtx(c *gin.Context, msg string)
	InfoCtx(c *gin.Context, msg string)
	WarnCtx(c *gin.Context, msg string)
	ErrorCtx(c *gin.Context, msg string, err error)
	FatalCtx(c *gin.Context, msg string, err error)
	PanicCtx(c *gin.Context, msg string, r any, trace string) // r = recover()
	SuccessCtx(c *gin.Context)

	// This set of functions can be used in scenarios where there is no
	// gin.Context (server context) involved
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string, err error)
	Fatal(msg string, err error)

	// Logging middleware to be used only as a global middleware during router
	// initialization
	LogMiddleware(c *gin.Context)
}

func InitLogger(env string) (*LoggerService, error) {
	var output io.Writer

	switch env {
	case "DEVELOPMENT":
		file, err := os.OpenFile("dev.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}
		fileWriter := zerolog.ConsoleWriter{
			Out:        file,
			TimeFormat: "",
			FormatFieldName: func(i any) string {
				return fmt.Sprintf("%s=", i)
			},
			FormatFieldValue: func(i any) string {
				s := fmt.Sprintf("%v", i)
				if strings.ContainsAny(s, " \t\n\r") {
					return fmt.Sprintf("%q", s)
				}
				return s
			},
			FormatTimestamp: func(i any) string {
				t, err := time.Parse(time.RFC3339, i.(string))
				if err != nil {
					return fmt.Sprintf("time=%q", i) // Fallback if parsing fails
				}
				return fmt.Sprintf("time=%d", t.UnixMilli())
			},
			FormatLevel: func(i any) string {
				return fmt.Sprintf("level=%q", i)
			},
			FormatMessage: func(i any) string {
				return fmt.Sprintf("msg=%q", i) // Quoting the message automatically
			},
			NoColor: true,
		}
		output = zerolog.MultiLevelWriter(consoleWriter, fileWriter)
	case "PRODUCTION":
		file, err := os.OpenFile("prod.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		output = zerolog.ConsoleWriter{
			Out:        file,
			TimeFormat: "",
			FormatFieldName: func(i any) string {
				return fmt.Sprintf("%s=", i)
			},
			FormatFieldValue: func(i any) string {
				s := fmt.Sprintf("%v", i)
				if strings.ContainsAny(s, " \t\n\r") {
					return fmt.Sprintf("%q", s)
				}
				return s
			},
			FormatTimestamp: func(i any) string {
				t, err := time.Parse(time.RFC3339, i.(string))
				if err != nil {
					return fmt.Sprintf("time=%q", i) // Fallback if parsing fails
				}
				return fmt.Sprintf("time=%d", t.UnixMilli())
			},
			FormatLevel: func(i interface{}) string {
				return fmt.Sprintf("level=%s", i)
			},
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("msg=%q", i) // Quoting the message automatically
			},
			NoColor: true,
		}
	default:
		return nil, errors.New("invalid environment for logger setup")
	}

	logger := zerolog.New(output).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = time.RFC3339Nano
	return &LoggerService{
		Logger: logger,
		Env:    env,
	}, nil
}

func (l *LoggerService) enrich(c *gin.Context, e *zerolog.Event) *zerolog.Event {
	// Path parameters
	pathParams := make(map[string]string)
	for _, param := range c.Params {
		pathParams[param.Key] = param.Value
	}
	pathParamsStr := fmt.Sprintf("%v", pathParams)

	// Convert query params from map[string][]string → map[string]string
	queryParams := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) == 1 {
			queryParams[k] = v[0]
		} else {
			queryParams[k] = strings.Join(v, ",")
		}
	}
	queryParamsStr := fmt.Sprintf("%v", queryParams)

	return e.
		Str("req-id", GrabRequestId(c)).
		Str("route", c.FullPath()).
		Str("method", c.Request.Method).
		Str("path-params", pathParamsStr).
		Str("query-params", queryParamsStr).
		Str("ip", c.ClientIP()).
		Str("user-agent", c.Request.UserAgent())
}

func (l *LoggerService) DebugCtx(c *gin.Context, msg string) {
	if l.Env == "PRODUCTION" {
		return
	}
	event := l.Logger.WithLevel(zerolog.DebugLevel)
	l.enrich(c, event).Msg(msg)
}

func (l *LoggerService) InfoCtx(c *gin.Context, msg string) {
	event := l.Logger.WithLevel(zerolog.InfoLevel)
	l.enrich(c, event).Msg(msg)
}

func (l *LoggerService) WarnCtx(c *gin.Context, msg string) {
	event := l.Logger.WithLevel(zerolog.WarnLevel)
	l.enrich(c, event).Msg(msg)
}

func (l *LoggerService) ErrorCtx(c *gin.Context, msg string, err error) {
	event := l.Logger.WithLevel(zerolog.ErrorLevel).Err(err)
	l.enrich(c, event).Msg(msg)
}

func (l *LoggerService) FatalCtx(c *gin.Context, msg string, err error) {
	event := l.Logger.WithLevel(zerolog.FatalLevel).Err(err)
	l.enrich(c, event).Msg(msg)
}

func (l *LoggerService) PanicCtx(c *gin.Context, msg string, r any, trace string) {
	event := l.Logger.WithLevel(zerolog.InfoLevel).
		Str("panic_value", fmt.Sprintf("%v", r)).
		Str("trace", trace)
	l.enrich(c, event).Msg(msg)
}

func (l *LoggerService) SuccessCtx(c *gin.Context) {
	event := l.Logger.WithLevel(zerolog.InfoLevel)
	l.enrich(c, event).Msg("request successful")
}

func (l *LoggerService) Debug(msg string) {
	if l.Env == "PRODUCTION" {
		return
	}
	l.Logger.WithLevel(zerolog.DebugLevel).Msg(msg)
}

func (l *LoggerService) Info(msg string) {
	l.Logger.WithLevel(zerolog.InfoLevel).Msg(msg)
}

func (l *LoggerService) Warn(msg string) {
	l.Logger.WithLevel(zerolog.InfoLevel).Msg(msg)
}

func (l *LoggerService) Error(msg string, err error) {
	l.Logger.WithLevel(zerolog.InfoLevel).Err(err).Msg(msg)
}

func (l *LoggerService) Fatal(msg string, err error) {
	l.Logger.WithLevel(zerolog.FatalLevel).Err(err).Msg(msg)
}

func (l *LoggerService) LogMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()

	pathParams := map[string]string{}
	for _, param := range c.Params {
		pathParams[param.Key] = param.Value
	}
	pathParamsStr := fmt.Sprintf("%v", pathParams)

	queryParams := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) == 1 {
			queryParams[k] = v[0]
		} else {
			queryParams[k] = strings.Join(v, ",")
		}
	}
	queryParamsStr := fmt.Sprintf("%v", queryParams)

	l.Logger.WithLevel(zerolog.InfoLevel).
		Str("req-id", GrabRequestId(c)).
		Str("route", c.FullPath()).
		Str("method", c.Request.Method).
		Int("status", c.Writer.Status()).
		Int("response-size", c.Writer.Size()).
		Dur("duration", time.Since(start)).
		Str("path-params", pathParamsStr).
		Str("query-params", queryParamsStr).
		Str("ip", c.ClientIP()).
		Str("user-agent", c.Request.UserAgent())
}
