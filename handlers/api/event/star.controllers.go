package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Infinite-Sum-Games/Am.EventKit/cmd"
	db "github.com/Infinite-Sum-Games/Am.EventKit/db/gen"
	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func StarEvent(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "EVENT")
	if !ok {
		return
	}

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
