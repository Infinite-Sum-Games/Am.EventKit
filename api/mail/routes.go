package mail

import "github.com/gin-gonic/gin"

func SetRoutes(router *gin.Engine, controller *Controller) {
	mailGroup := router.Group("/test")
	{
		mailGroup.POST("/mail", controller.SendTestEmail)
	}
}
