package router

import (
	attendanceApi "github.com/Infinite-Sum-Games/Am.EventKit/api/attendance"
	authApi "github.com/Infinite-Sum-Games/Am.EventKit/api/auth"
	organizersApi "github.com/Infinite-Sum-Games/Am.EventKit/api/organizers"
	mw "github.com/Infinite-Sum-Games/Am.EventKit/middleware"
	"github.com/gin-gonic/gin"
)

func OrganizerAppRouter(r *gin.RouterGroup) {
	// Attendance management routes
	org := r.Group("/org/app", mw.Auth)
	{
		// Event and schedule management
		org.GET("/events", attendanceApi.FetchEventsByOrganizer)
		org.GET("/events/:eventId/schedules/:scheduleId/participants", attendanceApi.FetchParticipantsByEvent)

		// Solo attendance management
		org.POST("/attendance/solo/:key/:studentId/:scheduleId", attendanceApi.MarkSoloCheckInOutBoth)
		org.POST("/attendance/solo/unmark/:key/:studentId/:scheduleId", attendanceApi.UnMarkSoloCheckInOutBoth)

		// Team attendance management
		org.POST("/attendance/team/:key/:studentId/:scheduleId", attendanceApi.MarkTeamCheckInOutBoth)
		org.POST("/attendance/team/unmark/:key/:studentId/:scheduleId", attendanceApi.UnMarkTeamCheckInOutBoth)

		// Organizer session management
		org.GET("/session", authApi.FetchOrganizerSession)
	}
}

func OrganizerWebRouter(r *gin.RouterGroup) {
	// Dashboard and web interface routes
	org := r.Group("/org/web", mw.Auth)
	{
		// Organizer dashboard
		org.GET("/dashboard", organizersApi.GetOrganizerEvents)

		// Event participant management
		org.GET("/events/:eventId/participants", organizersApi.GetOrganizerEventParticipantList)
		org.GET("/session", authApi.FetchOrganizerSession)
	}
}
