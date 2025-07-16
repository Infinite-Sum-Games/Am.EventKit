package api

import "github.com/gin-gonic/gin"

func AttendanceRoutes(r *gin.RouterGroup) {
	r.GET("/list/:eventId", FetchEventParticipantList)
	r.POST("/mark/:eventId/scan", MarkOneTimeAttendance)
	r.POST("/mark/:eventId/scan/check-in", MarkCheckIn)
	r.POST("/mark/:eventId/scan/check-out", MarkCheckOut)
}
