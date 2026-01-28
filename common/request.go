package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TagRequestWithId(c *gin.Context) {
	id := uuid.New()
	c.Set("request_id", id.String())
	c.Next()
}

func GrabRequestId(c *gin.Context) string {
	reqId, ok := c.Get("request_id")
	if !ok {
		return "missing-id"
	}
	return fmt.Sprintf("%v", reqId)
}
