package mw

import (
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {

	// Extract refresh token
	refreshToken, refErr := c.Cookie("refresh_token")
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
	accessToken, authErr := c.Cookie("access_token")
	if authErr == nil && pkg.VerifyTokens(c, accessToken, refreshToken) {
		c.Next()
		return
	}
	if authErr == http.ErrNoCookie {
		validToken, err := pkg.VerifyRefreshToken(c, refreshToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Acess denied.",
			})
			return
		}
		refreshTokenClaims := validToken.Claims()
		userID, _ := refreshTokenClaims["USER-ID"].(string)
		username, _ := refreshTokenClaims["audience"].(string)
		email, _ := refreshTokenClaims["jti"].(string)
		isStudent, _ := refreshTokenClaims["STUDENT-ROLE"].(bool)
		isStaff, _ := refreshTokenClaims["STAFF-ROLE"].(bool)
		// Creating and setting auth token, so it can be used for future requests
		authToken := pkg.CreateAuthToken(userID, username, email, isStudent, false, isStaff)
		pkg.SetAuthCookie(c, authToken)
	}

	c.Next()
}
