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

func GetInsideCampusAnalytics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ANALYTICS") {
		return
	}
	defer conn.Release()

	q := db.New()

	insideCampusSummary, err := q.GetInsideCampusAnalyticsQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get inside campus summary", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":               "Successfully fetched inside campus analytics",
		"inside_campus_summary": insideCampusSummary,
	})
	pkg.Log.SuccessCtx(c)
}

func GetLiveBedsAnalytics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ANALYTICS") {
		return
	}
	defer conn.Release()

	q := db.New()

	liveBedsSummary, err := q.GetLiveBedsAnalyticsQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get live beds summary", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Successfully fetched live beds analytics",
		"live_beds_summary": liveBedsSummary,
	})
	pkg.Log.SuccessCtx(c)
}
