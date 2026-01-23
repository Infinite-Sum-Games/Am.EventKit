package mw

import (
	"fmt"
	"net/http"

	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
	"github.com/gin-gonic/gin"
)

func VerifyCsrf(c *gin.Context) {
	csrfToken := c.Request.Header["X-Csrf-Token"]
	if len(csrfToken) != 1 {
		pkg.Log.WarnCtx(c, "[CSRF-WARN]: Could not find CSRF header.")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Access denied.",
		})
		return
	}

	csrfFromHeader := csrfToken[0] // extracting into string var for reuse
	if csrfFromHeader == "" {
		pkg.Log.WarnCtx(c, "[CSRF-WARN]: Found empty CSRF header.")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Access denied.",
		})
		return
	}

	csrfFromCookie, csrfErr := c.Cookie("csrf_token")
	if csrfErr == http.ErrNoCookie {
		pkg.Log.WarnCtx(c, "[CSRF-WARN]: Could not find CSRF token in cookie.")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Access denied.",
		})
		return
	}

	if csrfFromCookie != csrfFromHeader {
		pkg.Log.ErrorCtx(
			c,
			"[CSRF-ERROR]: Mismatching CSRF tokens",
			fmt.Errorf("ERR: CSRF tokens' do not match"))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Access denied.",
		})
		return
	}

	// Check whether the CSRF token is parsable
	ok, token := pkg.ParseToken(csrfFromCookie, "csrf_token")
	if !ok {
		pkg.Log.ErrorCtx(c,
			"[CSRF-ERROR]: Failed to parse token",
			fmt.Errorf("given string is not a CSRF token"),
		)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Access denied",
		})
		return
	}

	// After parsing, extract the subject
	tokenAud, tokenAudErr := token.GetAudience()
	if tokenAudErr != nil {
		pkg.Log.ErrorCtx(c,
			"[CSRF-ERROR]: Could not find audience; bad token",
			fmt.Errorf("given string is not a CSRF token"),
		)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "User is forbidden.",
		})
		return
	}

	// After extraction, check the path
	if tokenAud != pkg.CsrfRoutes[c.FullPath()] {
		pkg.Log.ErrorCtx(c, "[CSRF-ERROR]: Invalid token for submitted form", nil)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "User is forbidden.",
		})
		return
	}

	c.Next()
}
