package api

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func FetchAllEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()

	events, err := q.GetEventsQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to fetch events", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Events list fetched successfully",
		"events":  events,
	})
	pkg.Log.SuccessCtx(c)
}

func FetchEventById(c *gin.Context) {
	eventId, ok := pkg.GrabUuid(c, c.Param("eventId"), "EVENT", "event")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	event, err := q.GetEventByIdQuery(ctx, conn, eventId)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Attempted to fetch event which does not exist", err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Request event does not exist", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event fetched successfully",
		"event":   event,
	})
	pkg.Log.SuccessCtx(c)
}
