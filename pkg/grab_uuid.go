package pkg

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GrabUuid(c *gin.Context, uuidStr string, route string, entity string) (uuid.UUID, bool) {
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The request is malformed.",
		})
		msg := fmt.Sprintf("[%s-ERROR]: Failed to parse %s UUID", route, entity)
		Log.ErrorCtx(c, msg, err)

		return uuid, false
	}
	return uuid, true
}
