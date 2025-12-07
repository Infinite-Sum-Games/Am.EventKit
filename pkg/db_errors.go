package pkg

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Use when acquiring connection might fail
func HandleDbAcquireErr(c *gin.Context, err error, path string) bool {
	if err == nil {
		return false
	}

	if err == context.DeadlineExceeded {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"message": "Server took too long to respond",
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
	}

	msg := fmt.Sprintf("[%s-FATAL]: Failed to acquire DB connection", path)
	Log.FatalCtx(c, msg, err)

	return true
}

// Use when transaction beginning might fail
func HandleDbTxnErr(c *gin.Context, err error, path string) bool {
	if err == nil {
		return false
	}

	if err == context.DeadlineExceeded {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"message": "Server took too long to respond",
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
	}

	msg := fmt.Sprintf("[%s-FATAL]: Failed to acquire DB txn", path)
	Log.FatalCtx(c, msg, err)

	return true
}

// Use when transaction commits might fail
func HandleDbTxnCommitErr(c *gin.Context, err error, path string) bool {
	if err == nil {
		return false
	}

	return true
}
