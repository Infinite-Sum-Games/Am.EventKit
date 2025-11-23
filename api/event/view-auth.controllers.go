package api

import (
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

// All events should be fetched along with their "favourite" status and
// registration status
func FetchAllEventsWithAuth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "All events fetched successfully",
	})
	pkg.Log.SuccessCtx(c)
}

// An event should be fetched with it's registration status
func FetchEventByEventIdWithAuth(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Event successfully fetched",
	})
	pkg.Log.SuccessCtx(c)
}
