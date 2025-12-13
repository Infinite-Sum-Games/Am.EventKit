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
	r.GET("/admin/new", NewEvent)
	r.POST("/admin/details", AddEventDetails)
	r.POST("/admin/poster", AddEventPoster)
	r.DELETE("/admin/poster/:eventId", DeleteEventPoster)
	r.POST("/admin/size", AddEventDimension)
	r.POST("/admin/toggle", AddEventToggles)
	r.POST("/admin/organizer", ConnectEventAndOrganizer)
	r.DELETE("/admin/organizer", DisconnectEventAndOrganizer)
	r.POST("/admin/tags", ConnectEventAndTags)
	r.DELETE("/admin/tags", DisonnectEventAndTags)
	r.POST("/admin/schedule", AddEventSchedule)
	r.PUT("/admin/schedule/:scheduleId", EditEventSchedule)
	r.DELETE("/admin/schedule/:scheduleId", DeleteEventSchedule)
	r.POST("/admin/publish/:eventId", PublishEvent)
	r.DELETE("/admin/publish/:eventId", UnpublishEvent)

	// Deprecated routes; under refactoring
	r.PUT("/:eventId", mw.Auth, mw.CheckAdmin, EditEvent)
	r.PUT("/:eventId/toggle-status", mw.Auth, mw.CheckAdmin, ToggleEventStatus)
}
