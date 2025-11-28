package tag

import (
	"github.com/gin-gonic/gin"
)

func TagRoutes(r *gin.RouterGroup) {
	r.GET("/", GetAllEventTags)
	r.GET("/:tagId", GetEventTagByID)
	r.POST("/", CreateEventTag)
	r.PATCH("/:tagId", EditEventTag)
	r.DELETE("/:tagId", DeleteEventTag)
}
