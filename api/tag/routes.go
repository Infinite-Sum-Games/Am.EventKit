package tag

import (
	mw "github.com/Infinite-Sum-Games/Am.EventKit/middleware"
	"github.com/gin-gonic/gin"
)

func TagRoutes(r *gin.RouterGroup) {
	r.GET("/", mw.Auth, mw.CheckAdmin, FetchEventTags)
	r.POST("/", mw.Auth, mw.CheckAdmin, CreateEventTag)
	r.PUT("/:tagId", mw.Auth, mw.CheckAdmin, EditEventTag)
	r.DELETE("/:tagId", mw.Auth, mw.CheckAdmin, DeleteEventTag)
}
