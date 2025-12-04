package mw

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckUser(c *gin.Context) {
	if c.GetBool("USER-ROLE") {
		c.Next()
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Access denied.",
	})
}

func CheckAdmin(c *gin.Context) {
	if c.GetBool("USER-ROLE") {
		c.Next()
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Access denied.",
	})
}

func CheckOrganizer(c *gin.Context) {
	if c.GetBool("USER-ROLE") {
		c.Next()
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Access denied.",
	})
}
