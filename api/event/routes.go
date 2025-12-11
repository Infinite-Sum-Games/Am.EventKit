package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func EventRoutes(r *gin.RouterGroup) {
	r.GET("/", FetchAllEvents)
	r.GET("/:eventId", FetchEventById)

	r.GET("/auth/", mw.Auth, mw.CheckUser, FetchAllEventsWithAuth)
	r.GET("/auth/:eventId", mw.Auth, mw.CheckUser, FetchEventByEventIdWithAuth)

	r.POST("/register/:eventId")
	r.GET("/register/:eventId")

	r.POST("/", mw.Auth, mw.CheckAdmin, CreateEvent)
	r.PUT("/:eventId", mw.Auth, mw.CheckAdmin, EditEvent)
	r.PUT("/:eventId/toggle-status", mw.Auth, mw.CheckAdmin, ToggleEventStatus)

	r.PUT("/favourite/:eventId", mw.Auth, mw.CheckUser, StarEvent)
	r.DELETE("/favourite/:eventId", mw.Auth, mw.CheckUser, UnstarEvent)

	r.GET("/dependency", mw.Auth, mw.CheckAdmin, GetAllDependencies)
	r.POST("/dependency", mw.Auth, mw.CheckAdmin, AddDependency)
	r.DELETE("/dependency", mw.Auth, mw.CheckAdmin, RemoveDependency)
}
