package mw

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
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
	pkg.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have user role")
}

func CheckAdmin(c *gin.Context) {
	if c.GetBool("ADMIN-ROLE") {
		c.Next()
		return
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Access denied.",
	})
	pkg.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have admin role")
}

func CheckOrganizer(c *gin.Context) {
	if c.GetBool("ORGANIZER-ROLE") {
		c.Next()
		return
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Organizer access denied.",
	})
	pkg.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have organizer role")
}

func CheckOrgAndAdmin(c *gin.Context) {
	if c.GetBool("ADMIN-ROLE") || c.GetBool("ORGANIZER-ROLE") {
		c.Next()
		return
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Admin access denied.",
	})
	pkg.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have organizer and admin role")
}

func CheckHospitality(c *gin.Context) {
	if c.GetBool("HOSPITALITY-ROLE") {
		c.Next()
		return
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": "Hospitality access denied.",
	})
	pkg.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have hospitality role")
}

func CheckGate(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "GATE-AUTH")
	if !ok {
		return
	}

	// Normalize just in case
	email = strings.ToLower(email)

	gateEmailRegex := regexp.MustCompile(`^[a-z0-9.]+\.gate@(cb\.students\.)?amrita\.edu$`)

	if !gateEmailRegex.MatchString(email) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Unauthorized gate email.",
		})
		pkg.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have hospitality-gate role")
	}
}

func CheckFinance(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "FINANCE-AUTH")
	if !ok {
		return
	}

	// Normalize just in case
	email = strings.ToLower(email)

	financeEmailRegex := regexp.MustCompile(`^[a-z0-9.]+\.finance@(cb\.students\.)?amrita\.edu$`)

	if !financeEmailRegex.MatchString(email) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Unauthorized finance email.",
		})
		pkg.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have hospitality-finance role")
	}
}

func CheckSecurity(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "SECURITY-AUTH")
	if !ok {
		return
	}

	// Normalize just in case
	email = strings.ToLower(email)

	securityEmailRegex := regexp.MustCompile(`^[a-z0-9.]+\.security@(cb\.students\.)?amrita\.edu$`)

	if !securityEmailRegex.MatchString(email) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Unauthorized security email.",
		})
		pkg.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have hospitality-security role")
	}
}

func CheckHostel(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "HOSTEL-AUTH")
	if !ok {
		return
	}

	// Normalize just in case
	email = strings.ToLower(email)

	hostelEmailRegex := regexp.MustCompile(`^[a-z0-9.]+\.hostel@(cb\.students\.)?amrita\.edu$`)

	if !hostelEmailRegex.MatchString(email) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Unauthorized hostel email.",
		})
		pkg.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have hospitality-hostel role")
	}
}
