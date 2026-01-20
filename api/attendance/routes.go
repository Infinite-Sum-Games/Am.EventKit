package api

import (
	mw "github.com/Infinite-Sum-Games/Am.EventKit/middleware"
	"github.com/gin-gonic/gin"
)

func AttendanceRoutes(r *gin.RouterGroup) {
	r.GET("/list/event", mw.Auth, mw.CheckOrganizer, FetchEventsByOrganizer)
	r.GET("/list/:eventId/:scheduleId", mw.Auth, mw.CheckOrganizer, FetchParticipantsByEvent)
	r.POST("/solo/mark/:key/:studentId/:scheduleId", mw.Auth, mw.CheckOrganizer, MarkSoloCheckInOutBoth)
	r.POST("/team/mark/:key/:studentId/:scheduleId", mw.Auth, mw.CheckOrganizer, MarkTeamCheckInOutBoth)
	r.DELETE("/solo/unMark/:key/:studentId/:scheduleId", mw.Auth, mw.CheckOrganizer, UnMarkSoloCheckInOutBoth)
	r.DELETE("/team/unMark/:key/:studentId/:scheduleId", mw.Auth, mw.CheckOrganizer, UnMarkTeamCheckInOutBoth)
}
