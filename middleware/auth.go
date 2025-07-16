package mw

import (
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {

	// Extract refresh token
	RefreshToken, refErr := c.Cookie("refresh_token")
	if refErr == http.ErrNoCookie {
		pkg.NullifyCookies(c)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Access denied.",
		})
		return
	}

	// FLOW: Extract AuthToken
	// 1. If AuthToken available then verify it
	// 2. If not available then check against DB and see if a refresh token
	// exists there and if it is a valid one or not
	// 3. If the token is valid then new authToken can be minted, added to
	// the cookie
	_, authErr := c.Cookie("access_token")
	if authErr == http.ErrNoCookie {
		_, err := pkg.VerifyRefreshToken(c, RefreshToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Acess denied.",
			})
			return
		}
	}

	c.Next()
}
