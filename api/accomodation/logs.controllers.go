package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func GateLogsSink(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "GATE-LOGS") {
		return
	}
	defer conn.Release()

	q := db.New()
	logs, err := q.GetLiveGateLogsQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[GATE-LOGS-ERROR]: Failed to fetch gate logs", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gate logs send successfully",
		"logs":    logs,
	})
	pkg.Log.SuccessCtx(c)
}

func HostelLogsSink(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "HOSTEL-LOGS") {
		return
	}
	defer conn.Release()

	q := db.New()
	logs, err := q.GetLiveHostelCheckInQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[HOSTEL-LOGS-ERROR]: Failed to fetch hostel logs", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hostel logs sent successfully",
		"logs":    logs,
	})
	pkg.Log.SuccessCtx(c)
}
