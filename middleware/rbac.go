package mw

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckUser(c *gin.Context) {
	if c.GetBool("STUDENT-ROLE") {
		c.Next()
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Access denied.",
	})
}

func CheckAdmin(c *gin.Context) {
	if c.GetBool("ADMIN-ROLE") {
		c.Next()
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Access denied.",
	})
}

func CheckOrganizer(c *gin.Context) {
	if c.GetBool("ORGANIZER-ROLE") {
		c.Next()
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Access denied.",
	})
}
