package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	api "github.com/Thanus-Kumaar/anokha-2025-backend/api/util"
	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func LoginUserCsrf(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Login action initiated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func LoginUser(c *gin.Context) {
	// use the ValidateRequest function if the models.Type has a Validate() function defined
	req, ok := api.ValidateRequest[models.LoginRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})

		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()
	user, err := q.CheckStudentVerifiedQuery(ctx, conn, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found. Please register first."})
			return
		}
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to fetch user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		return
	}
	if req.HashedPassword != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password for given email"})
		return
	}

	var refreshToken string
	isNewToken := false
	if user.RefreshToken.String == "" {
		// no token exists, create new refresh token and set it in database
		refreshToken = pkg.CreateRefreshToken(user.ID.String(), user.Name, user.Email, true, false, false)
		isNewToken = true
	} else {
		// Validate existing token
		valid, _ := pkg.ParseToken(user.RefreshToken.String, "refresh_token")
		if !valid {
			refreshToken = pkg.CreateRefreshToken(user.ID.String(), user.Name, user.Email, true, false, false)
			isNewToken = true
		} else {
			refreshToken = user.RefreshToken.String
		}
	}
	if isNewToken {
		err = q.UpdateRefreshTokenQuery(ctx, conn, db.UpdateRefreshTokenQueryParams{
			ID: user.ID,
			RefreshToken: pgtype.Text{
				String: refreshToken,
				Valid:  refreshToken != "",
			},
		})
		if err != nil {
			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to update refresh token", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
			return
		}
	}
	// createing auth token and setting both tokens as cookies
	authToken := pkg.CreateAuthToken(user.ID.String(), user.Name, user.Email, true, false, false)
	pkg.SetAuthCookie(c, authToken)
	pkg.SetRefreshCookie(c, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func LoginOrganizerCsrf(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Organizer login action initiated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func LoginOrganizer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Organizer logged in successfully",
	})
	pkg.Log.SuccessCtx(c)
}
