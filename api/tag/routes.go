package tag

import (
	"github.com/gin-gonic/gin"
)

func TagRoutes(r *gin.RouterGroup) {
	r.GET("/", GetAllEventTags)
}
