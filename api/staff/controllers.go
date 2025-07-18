package api

import (
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func FetchEventParticipantList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Participants list fetched successfully",
	})
	pkg.Log.SuccessCtx(c)
}

// If both check-in and check-out is to be marked in one go, then
// utilize this controller, otherwise there are specific controllers
func MarkOneTimeAttendance(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "One-time attendance marked successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func MarkCheckIn(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Check-in marked successfully",
	})
	pkg.Log.SuccessCtx(c)

}

func MarkCheckOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Check-out marked successfully",
	})
	pkg.Log.SuccessCtx(c)

}
