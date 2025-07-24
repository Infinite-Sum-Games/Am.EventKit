package api

import (
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func SendCsrfToken(c *gin.Context) {
	// This will be repeated very often, so just created a function to use in routes
	csrfToken := pkg.CreateCsrfToken()
	pkg.SetCsrfCookie(c, csrfToken)
	c.JSON(http.StatusOK, gin.H{
		"message": "CSRF token created successfully",
		"token":   csrfToken,
	})
	pkg.Log.SuccessCtx(c)
}
