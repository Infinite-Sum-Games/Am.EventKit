package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// All events should be fetched along with their "favourite" status and registration status
func FetchAllEventsWithAuth(c *gin.Context) {
	email, ok1 := pkg.GrabEmail(c, "EVENT-AUTH")
	userIDStr, ok2 := pkg.GrabUserId(c, "EVENT-AUTH")
	if !ok1 || !ok2 {
		return
	}

	studentID, ok := pkg.GrabUuid(c, userIDStr, "EVENT-AUTH", "student")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "EVENT-AUTH") {
		return
	}
	defer conn.Release()

	q := db.New()
	events, err := q.GetEventsWithAuthQuery(ctx, conn, db.GetEventsWithAuthQueryParams{
		StudentID: studentID,
		Email:     email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-AUTH-ERROR]: Failed to fetch events with auth", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All events fetched successfully",
		"events":  events,
	})
	pkg.Log.SuccessCtx(c)
}

// An event should be fetched with it's registration status
func FetchEventByEventIdWithAuth(c *gin.Context) {
	email, ok1 := pkg.GrabEmail(c, "EVENT-AUTH")
	userIDStr, ok2 := pkg.GrabUserId(c, "EVENT-AUTH")
	if !ok1 || !ok2 {
		return
	}

	studentID, ok := pkg.GrabUuid(c, userIDStr, "EVENT-AUTH", "student")
	if !ok {
		return
	}

	eventID, ok := pkg.GrabUuid(c, c.Param("eventId"), "EVENT-AUTH", "event")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "EVENT-AUTH") {
		return
	}
	defer conn.Release()

	q := db.New()
	event, err := q.GetEventByIdWithAuthQuery(ctx, conn, db.GetEventByIdWithAuthQueryParams{
		ID:        eventID,
		StudentID: studentID,
		Email:     email,
	})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found.",
		})
		pkg.Log.WarnCtx(c, "[EVENT-AUTH-WARN]: Event not found")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-AUTH-ERROR]: Failed to fetch event with auth", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event successfully fetched",
		"event":   event,
	})
	pkg.Log.SuccessCtx(c)
}
