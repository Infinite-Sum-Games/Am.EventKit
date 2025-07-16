package mw

import (
	"net/http"
	"runtime/debug"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func RecoveryPanics(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			stackTrace := debug.Stack()
			pkg.Log.PanicCtx(c, "Panic recovered from HTTP request", r,
				string(stackTrace))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Oops! Something happened. Please try again later.",
			})
		}
	}()
	c.Next()
}
