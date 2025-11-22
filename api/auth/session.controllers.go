package api

import (
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func FetchUserSession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "User session obtained successfully",
	})
	pkg.Log.SuccessCtx(c)
}
