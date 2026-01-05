package mw

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
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

func CheckGate(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "GATE-AUTH")
	if !ok {
		return
	}

	// Normalize just in case
	email = strings.ToLower(email)

	gateEmailRegex := regexp.MustCompile(`^[a-z]+\.gate@amrita\.edu$`)

	if !gateEmailRegex.MatchString(email) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Unauthorized gate email.",
		})
	}
}

func CheckFinance(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "FINANCE-AUTH")
	if !ok {
		return
	}

	// Normalize just in case
	email = strings.ToLower(email)

	financeEmailRegex := regexp.MustCompile(`^[a-z]+\.finance@amrita\.edu$`)

	if !financeEmailRegex.MatchString(email) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Unauthorized finance email.",
		})
	}
}

func CheckSecurity(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "SECURITY-AUTH")
	if !ok {
		return
	}

	// Normalize just in case
	email = strings.ToLower(email)

	securityEmailRegex := regexp.MustCompile(`^[a-z]+\.security@amrita\.edu$`)

	if !securityEmailRegex.MatchString(email) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Unauthorized security email.",
		})
	}
}

func CheckHostel(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "HOSTEL-AUTH")
	if !ok {
		return
	}

	// Normalize just in case
	email = strings.ToLower(email)

	hostelEmailRegex := regexp.MustCompile(`^[a-z]+\.hostel@amrita\.edu$`)

	if !hostelEmailRegex.MatchString(email) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Unauthorized hostel email.",
		})
	}
}
