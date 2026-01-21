package router

import (
	attendApi "github.com/Infinite-Sum-Games/Am.EventKit/api/attendance"
	orgApi "github.com/Infinite-Sum-Games/Am.EventKit/api/organizers"
	mw "github.com/Infinite-Sum-Games/Am.EventKit/middleware"
	"github.com/gin-gonic/gin"
)

func OrganizerAppRouter(r *gin.RouterGroup) {
	// Attendance routes
	org := r.Group("/org/app", mw.Auth)
	{
		org.GET("/event", attendApi.FetchEventsByOrganizer)
		org.GET("/event/:eventId/:scheduleId", attendApi.FetchParticipantsByEvent)
		org.POST("/solo/mark", attendApi.MarkSoloCheckInOutBoth)
		org.POST("/team/mark", attendApi.MarkTeamCheckInOutBoth)
		org.POST("/solo/unmark", attendApi.UnMarkSoloCheckInOutBoth)
		org.POST("/team/unmark", attendApi.UnMarkTeamCheckInOutBoth)
	}
}

func OrganizerWebRouter(r *gin.RouterGroup) {
	// Dashboard routes
	org := r.Group("/org/web", mw.Auth)
	{
		org.GET("/dashboard", orgApi.GetOrganizerEvents)
		org.GET("/dashboard/:eventId", orgApi.GetOrganizerEvents)
	}
}
