package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"maps"
)

var Log *LoggerService

type LoggerService struct {
	Logger zerolog.Logger
	Env    string
}

type LoggerServiceInterface interface {
	// Error when something breaks (validation, DB fail, etc.)
	LogError(err error, msg string, ctx *gin.Context)

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

func (l *LoggerService) LogError(err error, msg string, ctx *gin.Context) {
	event := l.Logger.WithLevel(zerolog.ErrorLevel).
		Str("route", ctx.FullPath()).
		Str("method", ctx.Request.Method).Err(err).Caller()
	event.Msg(msg)
}

func (l *LoggerService) LogWarn(msg string, ctx *gin.Context) {
	l.Logger.WithLevel(zerolog.WarnLevel).
		Str("route", ctx.FullPath()).
		Str("function", ctx.Request.Method).Msg(msg)
}

func (l *LoggerService) LogInfo(msg string) {
	l.Logger.WithLevel(zerolog.InfoLevel).Msg(msg)
}

func (l *LoggerService) LogDebug(msg string, route string, function string) {
	// Only log debug messages in non-production environments
	if l.Env == "PROD" {
		return
	}

	l.Logger.WithLevel(zerolog.DebugLevel).
		Str("route", route).
		Str("function", function).
		Caller().Msg(msg)
}

func (l *LoggerService) LogFatal(msg string, err error) {
	l.Logger.WithLevel(zerolog.FatalLevel).
		Err(err).
		Caller().Msg(msg)
}

// NOTE: Do not use this function in any module or file, this is being used in middleware
func (l *LoggerService) LogRequest(c *gin.Context) {
	// Path parameters
	pathParams := map[string]string{}
	for _, param := range c.Params {
		pathParams[param.Key] = param.Value
	}

	// Query parameters
	queryParams := map[string][]string{}
	maps.Copy(queryParams, c.Request.URL.Query())

	event := l.Logger.WithLevel(zerolog.InfoLevel).
		Str("route", c.FullPath()).
		Str("method", c.Request.Method).
		Interface("path-params", pathParams).
		Interface("query-params", queryParams).
		Str("ip", c.ClientIP()).
		Str("user-agent", c.Request.UserAgent()).
		Str("function", "HTTPMiddleware")

	event.Msg("incoming HTTP request")
}
