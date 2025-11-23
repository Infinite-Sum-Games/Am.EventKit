package pkg

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

func TagRequestWithId(c *gin.Context) {
	id := ksuid.New()
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
