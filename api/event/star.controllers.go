package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func StarEvent(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "EVENT")
	if !ok {
		return
	}

	eventIdParam := c.Query("eventId")

	eventId, err := uuid.Parse(eventIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The request is malformed",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid UUID in request", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	eventId, err = q.MarkFavouriteEventQuery(ctx, conn, db.MarkFavouriteEventQueryParams{
		Email:   email,
		EventID: eventId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to mark event as favourite", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully marked favourite event",
		"eventId": eventId,
	})
	pkg.Log.SuccessCtx(c)
}

func UnstarEvent(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "EVENT")
	if !ok {
		return
	}

	eventIdParam := c.Query("eventId")

	eventId, err := uuid.Parse(eventIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The request is malformed",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid UUID in request", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()
	eventId, err = q.UnmarkFavouriteEventQuery(ctx, conn, db.UnmarkFavouriteEventQueryParams{
		Email:   email,
		EventID: eventId,
	})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Could not unmark event not already a favourite",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Attempt to unmark event not previous marked", err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to unmark favourite", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully unmarked favourite event",
		"eventId": eventId,
	})
	pkg.Log.SuccessCtx(c)
}
