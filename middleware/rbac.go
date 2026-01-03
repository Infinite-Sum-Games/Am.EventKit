package mw

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckUser(c *gin.Context) {
	if c.GetBool("STUDENT-ROLE") {
		c.Next()
		return
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Access denied.",
	})
}

func CheckAdmin(c *gin.Context) {
	if c.GetBool("ADMIN-ROLE") {
		c.Next()
		return
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Access denied.",
	})
}

func CheckOrganizer(c *gin.Context) {
	if c.GetBool("ORGANIZER-ROLE") {
		c.Next()
		return
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Organizer access denied.",
	})
}

func CheckOrgAndAdmin(c *gin.Context) {
	if c.GetBool("ADMIN-ROLE") || c.GetBool("ORGANIZER-ROLE") {
		c.Next()
		return
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Admin access denied.",
	})
}

func CheckHospitality(c *gin.Context) {
	if c.GetBool("HOSPITALITY-ROLE") {
		c.Next()
		return
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Hospitality access denied.",
	})
}
