package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func EventRoutes(r *gin.RouterGroup) {
	// User event controllers
	r.GET("/", FetchAllEvents)
	r.GET("/:eventId", FetchEventById)
	r.GET("/auth/", mw.Auth, mw.CheckUser, FetchAllEventsWithAuth)
	r.GET("/auth/:eventId", mw.Auth, mw.CheckUser, FetchEventByEventIdWithAuth)
	r.PUT("/favourite/:eventId", mw.Auth, mw.CheckUser, StarEvent)
	r.DELETE("/favourite/:eventId", mw.Auth, mw.CheckUser, UnstarEvent)

	// Admin event controllers
	r.GET("/admin", mw.Auth, mw.CheckAdmin, GetAllAdminEvents)
	r.GET("/admin/:eventId", mw.Auth, mw.CheckAdmin, GetAdminEventsById)
	r.GET("/admin/new", mw.Auth, mw.CheckAdmin, NewEvent)
	r.POST("/admin/details/:eventId", mw.Auth, mw.CheckAdmin, AddEventDetails)
	r.POST("/admin/poster/:eventId", mw.Auth, mw.CheckAdmin, AddEventPoster)
	r.DELETE("/admin/poster/:eventId", mw.Auth, mw.CheckAdmin, DeleteEventPoster)
	r.POST("/admin/size/:eventId", AddEventDimension)
	r.POST("/admin/toggle/:eventId", mw.Auth, mw.CheckAdmin, AddEventToggles)
	r.POST("/admin/organizer", mw.Auth, mw.CheckAdmin, ConnectEventAndOrganizer)
	r.DELETE("/admin/organizer", mw.Auth, mw.CheckAdmin, DisconnectEventAndOrganizer)
	r.POST("/admin/tag", mw.Auth, mw.CheckAdmin, ConnectEventAndTags)
	r.DELETE("/admin/tag", mw.Auth, mw.CheckAdmin, DisonnectEventAndTags)
	r.POST("/admin/people", mw.Auth, mw.CheckAdmin, ConnectEventAndPeople)
	r.DELETE("/admin/people", mw.Auth, mw.CheckAdmin, DisconnectEventAndPeople)
	r.POST("/admin/schedule/:eventId", mw.Auth, mw.CheckAdmin, AddEventSchedule)
	r.PUT("/admin/schedule/:scheduleId", mw.Auth, mw.CheckAdmin, EditEventSchedule)
	r.DELETE("/admin/schedule/:scheduleId", mw.Auth, mw.CheckAdmin, DeleteEventSchedule)
	r.POST("/admin/publish/:eventId", mw.Auth, mw.CheckAdmin, PublishEvent)
	r.DELETE("/admin/publish/:eventId", mw.Auth, mw.CheckAdmin, UnpublishEvent)
	r.POST("/admin/completed/:eventId", mw.Auth, mw.CheckAdmin, MarkEventAsCompleted)
	r.DELETE("/admin/completed/:eventId", mw.Auth, mw.CheckAdmin, UnmarkEventAsCompleted)
}
