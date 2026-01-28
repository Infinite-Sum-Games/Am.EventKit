package router

import (
	"github.com/Infinite-Sum-Games/Am.EventKit/configs"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, cfg *configs.Config) {
	apiV2 := r.Group("/api/v2")
	{
		WebRouter(apiV2)
		AdminRouter(apiV2)
		HospitalityWebRouter(apiV2)
		HospitalityAppRouter(apiV2)
		OrganizerAppRouter(apiV2)
		OrganizerWebRouter(apiV2)
	}
}
