package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func FetchUserProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	idStr := c.Param("id")
	if err := validation.Validate(idStr, validation.Required, is.UUID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID format"})
		pkg.Log.ErrorCtx(c, "[TAG-ERROR]: Invalid user ID format", err)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID format"})
		pkg.Log.ErrorCtx(c, "[TAG-ERROR]: Invalid user ID format", err)
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to acquire DB connection"})
		pkg.Log.ErrorCtx(c, "[TAG-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	profile, err := q.ListProfileInfo(ctx, conn, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "User profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch user profile"})
		pkg.Log.ErrorCtx(c, "[TAG-ERROR]: Failed to fetch user profile", err)
		return
	}

	if profile.AccountStatus != db.AccountStatusEnumVERIFIED {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User profile is not verified"})
		pkg.Log.ErrorCtx(c, "[TAG-ERROR]: User profile is not verified", nil)
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
