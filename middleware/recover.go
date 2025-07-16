package mw

import (
	"fmt"
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func RecoveryPanics(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			pkg.Log.FatalCtx(
				c,
				"Panic recovered.",
				fmt.Errorf("%v\n", err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Oops! Something happened. Please try again later.",
			})
		}
	}()
	c.Next()
}
