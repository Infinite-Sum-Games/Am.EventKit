package api

import (
	"context"
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jackc/pgx/v5"
)

func FetchUserProfile(c *gin.Context) {
	ctx := context.Background()

	// TODO - Replace with email retireved from auth token
	email := "sample@gmail.com" // For testing purposes

	err := validation.Validate(email, validation.Required, is.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal server error"})
		pkg.Log.ErrorCtx(c, "[PROFILE-ERROR]: Invalid email format", err)
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch user profile"})
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
	c.JSON(http.StatusOK, gin.H{
		"message": "User profile update initiated",
		"key":     "",
	})
	pkg.Log.SuccessCtx(c)
}

func EditUserProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "User profile updated successfully",
	})
	pkg.Log.SuccessCtx(c)
}
