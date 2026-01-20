package mw

import (
	"net/http"
	"runtime/debug"

	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
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
