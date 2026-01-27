package common

import (
	"net/http"

	"github.com/Infinite-Sum-Games/Am.EventKit/logger"
	"github.com/gin-gonic/gin"
)

type Validatable interface {
	Validate() error
}

func ValidateRequest[T Validatable](c *gin.Context) (*T, bool) {
	var req T

	if err := c.BindJSON(&req); err != nil {
		logger.Log.ErrorCtx(c, "[REQ-ERROR]: Failed to bind JSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "The request is malformed"})
		return nil, false
	}

	if err := req.Validate(); err != nil {
		logger.Log.ErrorCtx(c, "[REQ-ERROR]: Validation failed", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "The request is malformed"})
		return nil, false
	}

	return &req, true
}
