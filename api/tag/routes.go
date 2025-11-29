package tag

import (
	"github.com/gin-gonic/gin"
)

func TagRoutes(r *gin.RouterGroup) {
	r.GET("/", GetAllEventTags)
	r.POST("/", CreateEventTag)
	r.PUT("/:tagId", EditEventTag)
	r.DELETE("/:tagId", DeleteEventTag)
}
