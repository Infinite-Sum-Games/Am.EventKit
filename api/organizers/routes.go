package organizers

import "github.com/gin-gonic/gin"

func OrganizerRoutes(r *gin.RouterGroup) {
	r.GET("/all", GetAllOrganizers);
}