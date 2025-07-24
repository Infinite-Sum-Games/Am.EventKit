package event

import (
	"github.com/gin-gonic/gin"
)

func EventRoutes(r *gin.RouterGroup) {
	r.GET("/events/", FetchAllEvents)
	r.GET("/events/:eventId", FetchEventById)
}
