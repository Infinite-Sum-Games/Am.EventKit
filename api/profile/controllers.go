package api

import (
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func FetchUserProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "User profile fetched successfully",
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
