package pkg

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SetAuthCookie(c *gin.Context, authTokenString string) {
	if viper.GetString("env") == "PRODUCTION" {
		c.SetSameSite(http.SameSiteStrictMode)
	} else {
		c.SetSameSite(http.SameSiteNoneMode)
	}
	c.SetCookie(
		"access_token",       // key
		authTokenString,      // value
		3600,                 // maxAge (1 hour)
		"/",                  // path
		cmd.Env.CookieDomain, // domain
		cmd.Env.CookieSecure, // secure
		true,                 // httpOnly
	)
}

func SetRefreshCookie(c *gin.Context, refreshTokenString string) {
	if viper.GetString("env") == "PRODUCTION" {
		c.SetSameSite(http.SameSiteStrictMode)
	} else {
		c.SetSameSite(http.SameSiteNoneMode)
	}
	c.SetCookie(
		"refresh_token",      // key
		refreshTokenString,   // value
		3600*24*90,           // maxAge (90 days)
		"/",                  // path
		cmd.Env.CookieDomain, // domain
		cmd.Env.CookieSecure, // secure
		true,                 // httpOnly
	)
}

func SetTempCookie(c *gin.Context, tempTokenString string) {
	if viper.GetString("env") == "PRODUCTION" {
		c.SetSameSite(http.SameSiteStrictMode)
	} else {
		c.SetSameSite(http.SameSiteNoneMode)
	}
	c.SetCookie(
		"temp_token",         // key
		tempTokenString,      // value
		5*60,                 // maxAge (5 mins)
		"/",                  // path
		cmd.Env.CookieDomain, // domain
		cmd.Env.CookieSecure, // secure
		true,                 // httpOnly
	)
}

func SetCsrfCookie(c *gin.Context, csrfTokenString string) {
	if viper.GetString("env") == "PRODUCTION" {
		c.SetSameSite(http.SameSiteStrictMode)
	} else {
		c.SetSameSite(http.SameSiteNoneMode)
	}
	c.SetCookie(
		"csrf_token",         // key
		csrfTokenString,      // value
		5*60,                 // maxAge (5 minutes)
		c.Request.URL.Path,   // path to be constructed for restriction
		cmd.Env.CookieDomain, // domain
		cmd.Env.CookieSecure, // secure
		true,                 // httpOnly
	)
}

/*
* Nullify cookies during LogOut and ForbiddenAccess situations
 */
func NullifyCookies(c *gin.Context) {
	if viper.GetString("env") == "PRODUCTION" {
		c.SetSameSite(http.SameSiteStrictMode)
	} else {
		c.SetSameSite(http.SameSiteNoneMode)
	}

	c.SetCookie("access_token", "", -1, "/", cmd.Env.CookieDomain, cmd.Env.CookieSecure, true)
	c.SetCookie("refresh_token", "", -1, "/", cmd.Env.CookieDomain, cmd.Env.CookieSecure, true)
	c.SetCookie("csrf_token", "", -1, "/", cmd.Env.CookieDomain, cmd.Env.CookieSecure, true)

	email, exists := c.Get("email")
	if !exists {
		return
	}
	RevokeRefreshToken(c, email.(string))
}

/*
 * We are revoking the refresh-token so that you cannot use it to get any more
 * Auth-Rokens in-case you have managed to steal the token from the browser
 * and kept it somewhere
 */
func RevokeRefreshToken(c *gin.Context, email string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if HandleDbAcquireErr(c, err, "AUTH") {
		return
	}
	defer conn.Release()

	isStudent := c.GetBool("STUDENT-ROLE")
	isOrganizer := c.GetBool("ORGANIZER-ROLE")
	isAdmin := c.GetBool("ADMIN-ROLE")

	q := db.New()

	if isStudent {
		_, err = q.RevokeRefreshTokenQuery(ctx, conn, email)
	} else if isOrganizer {
		_, err = q.RevokeOrganizerRefreshTokenQuery(ctx, conn, email)
	} else if isAdmin {
		_, err = q.RevokeAdminRefreshTokenQuery(ctx, conn, email)
	}

	if err != nil {
		Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to revoke Refresh Token in DB", err)
		return
	}
	Log.InfoCtx(c, "[AUTH-INFO]: Successfully revoked Refresh Token in DB")
}
