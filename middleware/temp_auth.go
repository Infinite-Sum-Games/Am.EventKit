package mw

import (
	"net/http"

	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
	"github.com/gin-gonic/gin"
)

func TempAuth(c *gin.Context) {
	tempToken, err := c.Cookie("temp_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Access denied.",
		})
		pkg.Log.ErrorCtx(c, "[REQ-ERROR]: Missing temporary auth token", err)
		return
	}

	if !pkg.VerifyTempToken(c, tempToken) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "User is forbidden",
		})
		pkg.Log.ErrorCtx(c, "[REQ-ERROR]: Temporary token could not be verified", err)
		return
	}

	c.Next()
}
