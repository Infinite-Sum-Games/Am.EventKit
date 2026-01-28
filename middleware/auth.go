package mw

import (
	"net/http"

	"github.com/Infinite-Sum-Games/Am.EventKit/internal/helpers"
	"github.com/Infinite-Sum-Games/Am.EventKit/logger"
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {

	// Extract refresh token
	refreshToken, refErr := c.Cookie("refresh_token")
	if refErr == http.ErrNoCookie {
		helpers.NullifyCookies(c)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Access denied.",
		})
		logger.Log.ErrorCtx(c, "[AUTH-ERROR]: Refresh Cookie is missing", refErr)
		return
	}

	// FLOW: Extract AuthToken
	// 1. If AuthToken available then verify it and set the gin.Context
	// 2. If not available then check against DB and see if a refresh token
	// exists there and if it is a valid one or not
	// 3. If the token is valid then new authToken can be minted, added to
	// the cookie and the gin.Context be populate as well

	accessToken, accessErr := c.Cookie("access_token")
	if accessErr == nil && helpers.VerifyTokens(c, accessToken, refreshToken) {
		c.Next()
		return
	}

	if accessErr == http.ErrNoCookie {
		// Check if refresh token is valid. If yes, only then check for access token
		// validity. If access token is valid then setup gin.Context map otherwise
		// mint new token and then setup gin.Context map
		validToken, err := helpers.VerifyRefreshToken(c, refreshToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Access denied.",
			})
			logger.Log.ErrorCtx(c, "[COOKIE-ERROR]: Failed to verify refresh token", err)
			return
		}

		refreshTokenClaims := validToken.Claims()
		userId, _ := refreshTokenClaims["aud"].(string)
		email, _ := refreshTokenClaims["jti"].(string)
		isStudent, _ := refreshTokenClaims["STUDENT-ROLE"].(bool)
		isOrganizer, _ := refreshTokenClaims["ORGANIZER-ROLE"].(bool)
		isAdmin, _ := refreshTokenClaims["ADMIN-ROLE"].(bool)
		isHospitality, _ := refreshTokenClaims["HOSPITALITY-ROLE"].(bool)

		// Creating and setting auth token, so it can be used for future requests
		authToken, err := helpers.CreateAuthToken(userId, email, helpers.Roles{
			IsUser:        isStudent,
			IsOrganizer:   isOrganizer,
			IsAdmin:       isAdmin,
			IsHospitality: isHospitality,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})
			logger.Log.FatalCtx(c, "[COOKIE-ERROR]: Failed to mint new auth token", err)
			return
		}
		helpers.SetAuthCookie(c, authToken)
	}

	c.Next()
}
