package api

import (
	"context"
	"net/http"

	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func FetchUserProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO - Replace with email retrieved from auth token
	email := "sample@gmail.com" // For testing purposes

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	profile, err := q.FetchUserProfileQuery(ctx, conn, email)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "User profile does not exist"})
		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: User profile does not exist", nil)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
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
	// TODO - Replace with email retireved from auth token
	email := "sample@gmail.com" // For testing purposes

	csrfToken := pkg.CreateCsrfToken(email, "edit_profile")

	pkg.SetCsrfCookie(c, csrfToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile update initiated",
		"key":     csrfToken,
	})
	pkg.Log.SuccessCtx(c)
}

func EditUserProfile(c *gin.Context) {

	mw.VerifyCsrf(c)

	// TODO - Replace with email retrieved from auth token
	email := "sample@gmail.com" // For testing purposes"

	var req models.EditProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: Invalid request payload", err)
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation error", "errors": err.Error()})
		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: Validation error", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: Failed to begin DB transaction", err)
		return
	}

	defer tx.Rollback(ctx)

	q := db.New()

	rowAffected, err := q.EditUserProfileQuery(ctx, tx, db.EditUserProfileQueryParams{
		Email:       email,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		CollegeName: req.CollegeName,
		CollegeCity: req.CollegeCity,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: Failed to edit user profile", err)
		return
	} else if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User profile does not exist"})
		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: User profile does not exist", nil)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: Failed to commit DB transaction", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile updated successfully",
	})
	pkg.Log.SuccessCtx(c)
}
