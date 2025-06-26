package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type LoggerService struct {
	Logger zerolog.Logger
	Env    string
}

type LoggerServiceInterface interface {
	// Error when something breaks (validation, DB fail, etc.)
	LogError(err error, msg string, route string, function string)

	// Warn about something unusual but not fatal (e.g., cache miss)
	LogWarn(msg string, route string, function string)

	// Info log for successful operations (e.g., user signed up)
	LogInfo(msg string)

	// Debug logs (dev only, hide in prod)
	LogDebug(msg string, route string, function string)

	// Log an HTTP request's metadata (used in middleware)
	LogRequest(c *gin.Context)

	// For startup-time fatal crashes (optional)
	LogFatal(msg string, err error)
}

func (l *LoggerService) LogError(err error, msg string, route string, function string) {
	event := l.Logger.Error().
		Str("route", route).
		Str("function", function).Err(err).Caller()
	event.Msg(msg)
}

func (l *LoggerService) LogWarn(msg string, route string, function string) {
	l.Logger.Warn().
		Str("route", route).
		Str("function", function).Msg(msg)
}

func (l *LoggerService) LogInfo(msg string) {
	l.Logger.Info().Msg(msg)
}

func (l *LoggerService) LogDebug(msg string, route string, function string) {
	// Only log debug messages in non-production environments
	if l.Env == "PROD" {
		return
	}

	l.Logger.Debug().
		Str("route", route).
		Str("function", function).
		Caller().Msg(msg)
}

func (l *LoggerService) LogFatal(msg string, err error) {
	l.Logger.Fatal().
		Err(err).
		Caller().Msg(msg)
}

func (l *LoggerService) LogRequest(c *gin.Context) {
	// Extract query + path params into a map
	params := map[string]string{}
	for _, param := range c.Params {
		params[param.Key] = param.Value
	}
	for key, value := range c.Request.URL.Query() {
		if len(value) > 0 {
			params[key] = value[0]
		}
	}

	event := l.Logger.Info().
		Str("route", c.FullPath()).
		Str("method", c.Request.Method).
		Interface("params", params).
		Str("ip", c.ClientIP()).
		Str("user-agent", c.Request.UserAgent()).
		Str("function", "HTTPMiddleware")

	event.Msg("incoming HTTP request")
}

