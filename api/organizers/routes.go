package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func OrganizerRoutes(r *gin.RouterGroup) {
	r.GET("/", mw.Auth, mw.CheckAdmin, GetAllOrganizers)
	r.POST("/", mw.Auth, mw.CheckAdmin, CreateOrganizer)
	r.PUT("/:organizerId", mw.Auth, mw.CheckAdmin, EditOrganizer)
	r.DELETE("/:organizerId", mw.Auth, mw.CheckAdmin, DeleteOrganizer)
}
