package event

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FetchAllEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to acquire DB connection"})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	// Optional query param for filtering by organizerId
	organizerId := c.Query("organizerId")
	var filter interface{} = nil
	if organizerId != "" {
		filter = organizerId
	}

	events, err := q.ListEventsQuery(ctx, conn, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch events"})
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
	eventIdStr := c.Param("eventId")
	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID"})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid event ID", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to acquire DB connection"})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()
	event, err := q.GetEventByIdQuery(ctx, conn, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch event"})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to fetch event", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Event fetched successfully",
		"event":   event,
	})
	pkg.Log.SuccessCtx(c)
}
