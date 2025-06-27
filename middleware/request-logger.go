package middleware

import (
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func RequestLoggerMiddleware(logger *pkg.LoggerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogRequest(c)
		c.Next()
	}
}