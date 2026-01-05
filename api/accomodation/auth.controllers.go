package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func AccomodationLogin(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.LoginRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "AUTH") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "AUTH")

	q := db.New()
	result, err := q.LoginHospitalityQuery(ctx, tx, req.Email)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Hospitality user not found",
		})
		pkg.Log.WarnCtx(c, "[AUTH-WARN]: Email does not exist")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to fetch login details", err)
		return
	}

	err = pkg.CompareHash(result.Password, req.HashedPassword)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Invalid credentials",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Password verification failed", err)
		return
	}

	if result.RefreshToken.String == "" {
		refreshToken, err := pkg.CreateRefreshToken(result.ID.String(), req.Email, pkg.Roles{
			IsUser:        false,
			IsOrganizer:   false,
			IsAdmin:       false,
			IsHospitality: true,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to create refresh token", err)
			return
		}

		token, err := q.UpdateHospitalityRefreshTokenQuery(ctx, tx, db.UpdateHospitalityRefreshTokenQueryParams{
			RefreshToken: pgtype.Text{
				String: refreshToken,
				Valid:  true,
			},
			Email: req.Email,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Could not add refresh token to DB", err)
			return
		}
		pkg.SetRefreshCookie(c, token.String)
	} else {
		pkg.SetRefreshCookie(c, result.RefreshToken.String)
	}

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "AUTH") {
		return
	}

	authToken, err := pkg.CreateAuthToken(result.ID.String(), result.Email, pkg.Roles{
		IsUser:        false,
		IsOrganizer:   false,
		IsAdmin:       false,
		IsHospitality: true,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to create auth token", err)
		return
	}

	pkg.SetAuthCookie(c, authToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Accomodation personnel logged in successful",
		"name":    result.Name,
		"email":   result.Email,
	})
	pkg.Log.SuccessCtx(c)
}

func AccomodationLogout(c *gin.Context) {
	pkg.NullifyCookies(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "Accomodation personnel logged out successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func AccomodationSession(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "SESSION")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "SESSION") {
		return
	}
	defer conn.Release()

	q := db.New()

	result, err := q.FetchHospitalitySessionQuery(ctx, conn, email)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No session found for user",
		})
		pkg.Log.WarnCtx(c, "[SESSION-WARN]: User might deleted but cookies exist")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[SESSION-ERROR]: Failed to fetch user session", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User session obtained successfully",
		"name":    result.Name,
		"email":   result.Email,
	})
	pkg.Log.SuccessCtx(c)
}
