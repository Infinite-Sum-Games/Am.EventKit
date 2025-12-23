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
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quick draw!",
		"events":  events,
	})
	pkg.Log.SuccessCtx(c)
}
