package mw

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/Infinite-Sum-Games/Am.EventKit/internal/helpers"
	"github.com/Infinite-Sum-Games/Am.EventKit/logger"
	"github.com/gin-gonic/gin"
)

var (
	gateEmailRegex     = regexp.MustCompile(`^[a-z]+\.gate@amrita\.edu$`)
	financeEmailRegex  = regexp.MustCompile(`^[a-z]+\.finance@amrita\.edu$`)
	securityEmailRegex = regexp.MustCompile(`^[a-z]+\.security@amrita\.edu$`)
	hostelEmailRegex   = regexp.MustCompile(`^[a-z]+\.hostel@amrita\.edu$`)
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Any one of the roles are present in the gin.Context
		for _, role := range roles {
			if c.GetBool(role) {
				c.Next()
				return
			}
		}

		// Error handling path
		message := "Access denied."
		if len(roles) == 1 {
			switch roles[0] {
			case "STUDENT-ROLE":
				message = "Access denied."
				logger.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have user role")
			case "ADMIN-ROLE":
				message = "Admin access denied."
				logger.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have admin role")
			case "ORGANIZER-ROLE":
				message = "Organizer access denied."
				logger.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have organizer role")
			case "HOSPITALITY-ROLE":
				message = "Hospitality access denied."
				logger.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have hospitality role")
			}
		} else {
			message = "Access denied. Required role not found."
			logger.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have required role")
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": message,
		})
	}

	return fn
}

// Special function for handling roles that come up in hospitality division
func SplRole(role string) gin.HandlerFunc {
	var emailRegex *regexp.Regexp
	switch role {
	case "gate":
		emailRegex = gateEmailRegex
	case "finance":
		emailRegex = financeEmailRegex
	case "security":
		emailRegex = securityEmailRegex
	case "hostel":
		emailRegex = hostelEmailRegex
	default:
		emailRegex = regexp.MustCompile(`^[a-z]+\.` + role + `@amrita\.edu$`)
	}

	authKey := strings.ToUpper(role) + "-AUTH"
	message := "Unauthorized " + role + " email."

	fn := func(c *gin.Context) {
		email, ok := helpers.GrabEmail(c, authKey)
		if !ok {
			return
		}

		email = strings.ToLower(email)

		if !emailRegex.MatchString(email) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": message,
			})
			logger.Log.WarnCtx(c, "[RBAC-FAIL]: Does not have "+role+" role")
		}
	}
	return fn
}
