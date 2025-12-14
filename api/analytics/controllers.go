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
	if pkg.HandleDbAcquireErr(c, err, "ANALYTICS") {
		return
	}
	defer conn.Release()

	q := db.New()

	revenueSummary, err := q.GetRevenueAnalytics(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get revenue summary", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Successfully fetched revenue analytics",
		"revenue_summary": revenueSummary,
	})
	pkg.Log.SuccessCtx(c)
}

func GetEventRegistrationAnalytics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ANALYTICS") {
		return
	}
	defer conn.Release()

	q := db.New()

	registrationSummary, err := q.GetEventRegistrationAnalytics(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get registrations summary", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":              "Successfully fetched registration analytics",
		"registration_summary": registrationSummary,
	})
	pkg.Log.SuccessCtx(c)

}

func GetPeopleAnalytics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ANALYTICS") {
		return
	}
	defer conn.Release()

	q := db.New()

	peopleRegistrationSummary, err := q.GetPeopleRegistrationAnalytics(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get people registration summary", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Successfully fetched people analytics",
		"people_registration_summary": peopleRegistrationSummary,
	})
	pkg.Log.SuccessCtx(c)

}

func GetTransactionAnalytics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ANALYTICS") {
		return
	}
	defer conn.Release()

	q := db.New()

	transactionSummary, err := q.GetTransactionAnalytics(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get transactions summary", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Successfully fetched participants analytics",
		"transaction_summary":      transactionSummary,
	})
	pkg.Log.SuccessCtx(c)
}