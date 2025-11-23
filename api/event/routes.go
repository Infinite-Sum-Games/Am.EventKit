package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func EventRoutes(r *gin.RouterGroup) {
	r.GET("/", FetchAllEvents)
	r.GET("/:eventId", FetchEventById)

	r.GET("/auth/", FetchAllEventsWithAuth)
	r.GET("/auth/:eventId", FetchEventByEventIdWithAuth)

	r.POST("/register/:eventId")
	r.GET("/register/:eventId")

	r.GET("/tags", GetAllEventTags)
	r.GET("/organizers", GetAllOrganizers)

	r.PUT("/favourite/:eventId", mw.Auth, StarEvent)
	r.DELETE("/favourite/:eventId", mw.Auth, UnstarEvent)
}
