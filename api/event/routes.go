package event

import (
	"github.com/gin-gonic/gin"
)

func EventRoutes(r *gin.RouterGroup) {
	r.GET("/", FetchAllEvents)
	r.GET("/:eventId", FetchEventById)
	r.POST("/register/:eventId")
	r.GET("/register/:eventId")
}
