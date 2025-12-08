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

func GetRevenueAnalytics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.FatalCtx(c, "[ANALYTICS-FATAL]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	revenueAnalytics, err := q.GetRevenueAnalyticsQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get revenue analytics", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Successfully fetched revenue analytics",
		"revenue-analytics": revenueAnalytics,
	})
	pkg.Log.SuccessCtx(c)
}

func GetParticipantAnalytics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.FatalCtx(c, "[ANALYTICS-FATAL]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	participantsAnalytics, err := q.GetParticipantAnalyticsQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get participants analytics", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":                "Successfully fetched participants analytics",
		"participants-analytics": participantsAnalytics,
	})
	pkg.Log.SuccessCtx(c)
}

func GetRegistrationAnalytics(c *gin.Context) {

}

func GetPeopleAnalytics(c *gin.Context) {

}
