package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func AttendanceRoutes(r *gin.RouterGroup) {
	r.GET("/list/event/:organizerId", FetchEventsByOrganizer)
	r.GET("/list/:eventId", mw.Auth, mw.CheckOrganizer, FetchEventParticipantList)
	r.POST("/mark/:eventId/scan", mw.Auth, mw.CheckOrganizer, MarkOneTimeAttendance)
	r.POST("/mark/:eventId/scan/check-in", mw.Auth, mw.CheckOrganizer, MarkCheckIn)
	r.POST("/mark/:eventId/scan/check-out", mw.Auth, mw.CheckOrganizer, MarkCheckOut)
}
