package tag

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func TagRoutes(r *gin.RouterGroup) {
	r.GET("/", mw.Auth, mw.CheckAdmin, GetAllEventTags)
	r.POST("/", mw.Auth, mw.CheckAdmin, CreateEventTag)
	r.PUT("/:tagId", mw.Auth, mw.CheckAdmin, EditEventTag)
	r.DELETE("/:tagId", mw.Auth, mw.CheckAdmin, DeleteEventTag)
}
