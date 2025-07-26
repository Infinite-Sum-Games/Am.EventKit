package api

import (
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

type Validatable interface {
	Validate() error
}

func ValidateRequest[T Validatable](c *gin.Context) (*T, bool) {
	var req T
	if err := c.BindJSON(&req); err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to bind JSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return nil, false
	}
	if err := req.Validate(); err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Validation failed", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, false
	}
	return &req, true
}
