package pkg

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

	if err == context.DeadlineExceeded {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"message": "Server took too long to respond",
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
	}
	msg := fmt.Sprintf("[%s-FATAL]: Failed to commit DB transaction", path)
	Log.FatalCtx(c, msg, err)

	return true
}

func RollbackTx(c *gin.Context, tx pgx.Tx, ctx context.Context, path string) {
	err := tx.Rollback(ctx)
	if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		msg := fmt.Sprintf("[%s-FATAL]: Failed to rollback DB txn", path)
		Log.FatalCtx(c, msg, err)
	}
}
