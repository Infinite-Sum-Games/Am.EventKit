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
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Refresh Cookie is missing", refErr)
		return
	}

	// FLOW: Extract AuthToken
	// 1. If AuthToken available then verify it
	// 2. If not available then check against DB and see if a refresh token
	// exists there and if it is a valid one or not
	// 3. If the token is valid then new authToken can be minted, added to
	// the cookie
	accessToken, accessErr := c.Cookie("access_token")
	if accessErr == nil && pkg.VerifyTokens(c, accessToken, refreshToken) {
		c.Next()
		return
	}

	if accessErr == http.ErrNoCookie {

		// Check if refresh token is valid. If yes, mint a new access token
		// otherwise go back to old token
		validToken, err := pkg.VerifyRefreshToken(c, refreshToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Acess denied.",
			})
			return
		}

		refreshTokenClaims := validToken.Claims()
		userId, _ := refreshTokenClaims["audience"].(string)
		email, _ := refreshTokenClaims["jti"].(string)
		isStudent, _ := refreshTokenClaims["STUDENT-ROLE"].(bool)
		isOrganizer, _ := refreshTokenClaims["ORGANIZER-ROLE"].(bool)

		// Creating and setting auth token, so it can be used for future requests
		authToken, err := pkg.CreateAuthToken(userId, email, isStudent, isOrganizer)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})
			pkg.Log.FatalCtx(c, "[COOKIE-ERROR]: Failed to mint new auth token", err)
			return
		}

		pkg.SetAuthCookie(c, authToken)
	}

	c.Next()
}
