package helpers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GrabUserId(c *gin.Context, route string) (string, bool) {
	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		msg := fmt.Sprintf("[%s-FATAL]: Missing userId in context after auth middleware", route)
		Log.FatalCtx(c, msg, nil)

		return "", false
	}
	return userId, true
}
