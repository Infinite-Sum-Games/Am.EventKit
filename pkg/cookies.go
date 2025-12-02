package pkg

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/gin-gonic/gin"
)

// TODO: Have conditional rendering of None and Lax
func SetAuthCookie(c *gin.Context, authTokenString string) {
	c.SetSameSite(http.SameSiteNoneMode)
	// c.SetSameSite(http.SameSiteLaxMode)
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
	c.SetSameSite(http.SameSiteNoneMode)
	// c.SetSameSite(http.SameSiteLaxMode)
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
	c.SetSameSite(http.SameSiteNoneMode)
	// c.SetSameSite(http.SameSiteLaxMode)
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
	c.SetSameSite(http.SameSiteNoneMode)
	// c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"csrf_token",         // key
		csrfTokenString,      // value
		5*60,                 // maxAge (5 minutes)
		c.FullPath(),         // path to be constructed for restriction
		cmd.Env.CookieDomain, // domain
		cmd.Env.CookieSecure, // secure
		true,                 // httpOnly
	)
}

/*
* Nullify cookies during LogOut and ForbiddenAccess situations
 */
func NullifyCookies(c *gin.Context) {

	c.SetCookie("access_token", "", -1, "/", cmd.Env.CookieDomain, false, true)
	c.SetCookie("refesh_token", "", -1, "/", cmd.Env.CookieDomain, false, true)
	c.SetCookie("csrf_token", "", -1, "/", cmd.Env.CookieDomain, false, true)

	email, exists := c.Get("email")
	if !exists {
		Log.WarnCtx(c, "[AUTH-ERROR]: No email in gin.Context, failed to revoke Refresh-Token in DB")
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
	tx, err := cmd.DBPool.Begin(ctx)
	if err != nil {
		Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to Revoke Refresh Token in DB", err)
		return
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			Log.ErrorCtx(c, "[AUTH-ERROR]: Error in checking database for refresh tokens", err)
		}
	}()

	q := db.New()
	result, err := q.RevokeRefreshTokenQuery(ctx, tx, email)
	if err != nil || result.String != "" {
		Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to revoke Refresh Token in DB", err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		Log.FatalCtx(c, "[AUTH-FATAL]: Failed to commit txn while revoking refresh-token", err)
		return
	}

	Log.InfoCtx(c, "[AUTH-INFO]: Successfully revoked Refresh Token in DB")
}
