package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Infinite-Sum-Games/Am.EventKit/cmd"
	db "github.com/Infinite-Sum-Games/Am.EventKit/db/gen"
	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
	"github.com/gin-gonic/gin"
)

func GetQuickDashboard(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ANALYTICS") {
		return
	}
	defer conn.Release()

	q := db.New()
	events, err := q.GetQuickDashboardQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get analytics", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quick draw!",
		"events":  events,
	})
	pkg.Log.SuccessCtx(c)
}
