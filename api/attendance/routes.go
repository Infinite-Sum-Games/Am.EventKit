package api

import (
	"github.com/gin-gonic/gin"
)

func AttendanceRoutes(r *gin.RouterGroup) {
	r.GET("/list/event", FetchEventsByOrganizer)
	r.GET("/list/:eventId", FetchParticipantsByEvent)
	r.POST("/solo/mark/:key/:studentId/:scheduleId", MarkSoloCheckInOutBoth)
	r.POST("/team/mark/:key/:studentId/:scheduleId", MarkTeamCheckInOutBoth)
	// r.GET("/list/:eventId", mw.Auth, mw.CheckOrganizer, FetchEventParticipantList)
	// r.POST("/mark/:eventId/scan", mw.Auth, mw.CheckOrganizer, MarkOneTimeAttendance)
	// r.POST("/mark/:eventId/scan/check-in", mw.Auth, mw.CheckOrganizer, MarkCheckIn)
	// r.POST("/mark/:eventId/scan/check-out", mw.Auth, mw.CheckOrganizer, MarkCheckOut)
}
