package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func AttendanceRoutes(r *gin.RouterGroup) {
	r.GET("/list/:eventId", FetchEventParticipantList)
	r.POST("/mark/:eventId/scan", mw.Auth, MarkOneTimeAttendance)
	r.POST("/mark/:eventId/scan/check-in", mw.Auth, MarkCheckIn)
	r.POST("/mark/:eventId/scan/check-out", mw.Auth, MarkCheckOut)
}
