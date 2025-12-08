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
)

func FetchUserProfile(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "PROFILE")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "PROFILE") {
		return
	}
	defer conn.Release()

	q := db.New()

	profile, err := q.FetchUserProfileQuery(ctx, conn, email)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User profile does not exist",
		})

		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: User profile does not exist", nil)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})

		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: Failed to fetch user profile", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile fetched successfully",
		"profile": profile,
	})
	pkg.Log.SuccessCtx(c)
}

func EditUserProfileCsrf(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "PROFILE")
	if !ok {
		return
	}

	csrfToken, err := pkg.CreateCsrfToken(email, c)
	if err != nil {
		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: Failed to create CSRF token", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		return
	}

	pkg.SetCsrfCookie(c, csrfToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile update initiated",
		"key":     csrfToken,
	})
	pkg.Log.SuccessCtx(c)
}

func EditUserProfile(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "PROFILE")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.EditProfileRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "PROFILE") {
		return
	}

	q := db.New()
	rowAffected, err := q.EditUserProfileQuery(ctx, conn, db.EditUserProfileQueryParams{
		Email:       email,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		CollegeName: req.CollegeName,
		CollegeCity: req.CollegeCity,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: Failed to edit user profile", err)
		return
	}
	if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User profile does not exist",
		})
		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: User profile does not exist despite auth", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile updated successfully",
	})
	pkg.Log.SuccessCtx(c)
}
