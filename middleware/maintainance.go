package mw

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func MaintainanceMiddleware(c *gin.Context) {
	if viper.GetString("mode") == "MAINTAINANCE" {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"message": "Server in maintainance mode",
		})
		return
	}
	c.Next()
}
