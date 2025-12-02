package api

import (
	"github.com/gin-gonic/gin"
)

func OrganizerRoutes(r *gin.RouterGroup) {
	r.GET("/", GetAllOrganizers)
	r.POST("/", CreateOrganizer)
	r.PUT("/:organizerId", EditOrganizer)
	r.DELETE("/:organizerId", DeleteOrganizer)
}
