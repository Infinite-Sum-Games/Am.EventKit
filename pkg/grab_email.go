package pkg

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GrabEmail(c *gin.Context, route string) (string, bool) {
	email := c.GetString("email")
	if email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		msg := fmt.Sprintf("[%s-FATAL]: No email after crossing auth middleware", route)
		Log.FatalCtx(c, msg, nil)

		return "", false
	}

	return email, true
}
