package mw

import (
	"fmt"
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

/* CSRF Verification Process
*  1. Collect CSRF token from Header and Cookie
*  2.
 */

func VerifyCsrf(c *gin.Context) {
	// Check request header presence
	csrfToken := c.Request.Header["X-Csrf-Token"]
	if len(csrfToken) != 1 {
		pkg.Log.WarnCtx(c, "[CSRF-WARN]: Could not find CSRF header.")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "The request is malformed.",
		})
		return
	}

	csrfFromHeader := csrfToken[0] // extracting into string var for reuse
	if csrfFromHeader == "" {
		pkg.Log.WarnCtx(c, "[CSRF-WARN]: Found empty CSRF header.")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "The request is malformed.",
		})
		return
	}

	csrfFromCookie, csrfErr := c.Cookie("csrf_token")
	if csrfErr == http.ErrNoCookie {
		pkg.Log.WarnCtx(c, "[CSRF-WARN]: Could not find CSRF token in cookie.")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Missing security token.",
		})
		return
	}

	if csrfFromCookie != csrfFromHeader {
		pkg.Log.ErrorCtx(
			c,
			"[CSRF-ERROR]: Mismatching CSRF tokens",
			fmt.Errorf("ERR: CSRF tokens' do not match. Cookie: %s, Header: %s",
				csrfFromCookie,
				csrfFromHeader,
			))
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Security token compromised.",
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
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Security token compromised.",
		})
		return
	}

	// After parsing, extract the subject
	tokenSub, tokenSubErr := token.GetSubject()
	if tokenSubErr != nil {
		pkg.Log.ErrorCtx(c,
			"[CSRF-ERROR]: Could not find subject; bad token",
			fmt.Errorf("given string is not a CSRF token"),
		)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Security token compromised.",
		})
		return
	}

	// After extraction, check the path
	if tokenSub != pkg.CsrfRoutes[c.FullPath()] {
		pkg.Log.ErrorCtx(c, "[CSRF-ERROR]: ", fmt.Errorf(""))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "The request is malformed",
		})
		return
	}

	pkg.Log.InfoCtx(c, "[SUCCESS]: Verified CSRF token successfully.")
	c.Next()
}
