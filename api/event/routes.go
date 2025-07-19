package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func EventRoutes(r *gin.RouterGroup) {

	// Unauthenticated endpoints
	r.GET("/events/", FetchEventCatalog)
	r.GET("/events/:eventId", FetchEventById)

	// Authenticated endpoints
	r.GET("/events/registered", mw.Auth, FetchRegisteredEvents)
	r.POST("/events/:eventId/star", mw.Auth, ToggleFavouriteEvent)
}
