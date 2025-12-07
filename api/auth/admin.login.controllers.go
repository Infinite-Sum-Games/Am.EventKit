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

func LoginAdmin(c *gin.Context) {
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
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != pgx.ErrTxClosed {
			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to rollback", rbErr)
		}
	}()

	q := db.New()
	result, err := q.LoginAdminQuery(ctx, tx, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to fetch login details", err)
		return
	}
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.WarnCtx(c, "[AUTH-WARN]: Email does not exist")
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
			IsUser:      false,
			IsOrganizer: false,
			IsAdmin:     true,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to create refresh token", err)
			return
		}

		token, err := q.UpdateAdminRefreshTokenQuery(ctx, tx, db.UpdateAdminRefreshTokenQueryParams{
			Email: req.Email,
			RefreshToken: pgtype.Text{
				String: refreshToken,
				Valid:  true,
			},
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
		IsUser:      false,
		IsOrganizer: false,
		IsAdmin:     true,
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
		"message": "Admin logged in successful",
		"name":    result.Name,
		"email":   result.Email,
	})
}
