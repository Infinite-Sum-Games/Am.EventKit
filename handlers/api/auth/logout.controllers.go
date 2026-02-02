package api

import (
	"net/http"

	"github.com/Infinite-Sum-Games/Am.EventKit/logger"
	"github.com/gin-gonic/gin"
)

// The same logout handles for user, admin and organizer
func Logout(c *gin.Context) {
	pkg.NullifyCookies(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged out successfully",
	})
	logger.Log.SuccessCtx(c)
}
