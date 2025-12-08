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

	revenue, err := q.GetRevenueQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get revenue analytics", err)
		return
	}

	revenueSummary, err := q.GetRevenueSummaryQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get revenue summary", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Successfully fetched revenue analytics",
		"revenue":         revenue,
		"revenue-summary": revenueSummary,
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

	participants, err := q.GetParticipantQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get participants analytics", err)
		return
	}

	participantSummary, err := q.GetParticipantSummaryQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get participants summary", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Successfully fetched participants analytics",
		"participants": participants,
		"summary":      participantSummary,
	})
	pkg.Log.SuccessCtx(c)
}

func GetRegistrationAnalytics(c *gin.Context) {
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

	registration, err := q.GetRegistrationsQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get registrations analytics", err)
		return
	}

	registrationSummary, err := q.GetRegistrationSummaryQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get registrations summary", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":              "Successfully fetched registration analytics",
		"registration":         registration,
		"registration_summary": registrationSummary,
	})
	pkg.Log.SuccessCtx(c)

}

func GetPeopleAnalytics(c *gin.Context) {
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

	people, err := q.ListPeopleQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get people List", err)
		return
	}

	peopleCount, err := q.GetPeopleCountQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ANALYTICS-ERROR]: Failed to get people Count", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Successfully fetched people analytics",
		"people":       people,
		"people_count": peopleCount,
	})
	pkg.Log.SuccessCtx(c)

}
