package mw

import (
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func TempTokenAuth(c *gin.Context) {
	tempToken, err := c.Cookie("temp_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Missing temporary token",
		})
		return
	}
	if !pkg.VerifyTempToken(c, tempToken) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid or expired temp token",
		})
		return
	}
	c.Next()
}
