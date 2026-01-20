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

func GetOrganizerEvents(c *gin.Context) {
	orgIdStr, ok := pkg.GrabUserId(c, "ORGANIZER")
	if !ok {
		return
	}
	orgId, ok := pkg.GrabUuid(c, orgIdStr, "ORGANIZER", "Organizer")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ORGANIZER") {
		return
	}
	defer conn.Release()

	q := db.New()
	events, err := q.GetEventsByOrganizerIdQuery(ctx, conn, orgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Failed to fetch event list", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All organizer events fetched successfully",
		"events":  events,
	})
	pkg.Log.SuccessCtx(c)
}

func GetOrganizerEventParticipantList(c *gin.Context) {
	eventId, ok := pkg.GrabUuid(c, c.Param("eventId"), "ORGANIZER", "Event")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ORGANIZER") {
		return
	}
	defer conn.Release()

	q := db.New()

	isGroup, err := q.CheckIfGroupEventQuery(ctx, conn, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Failed to check if event is group", err)
		return
	}

	if isGroup {
		participants, err := q.GetOrganizerGroupEventParticipantListQuery(ctx, conn, eventId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Failed to fetch participant list", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "Participant list fetched successfully",
			"participants": participants,
		})
		pkg.Log.SuccessCtx(c)
		return
	}

	participants, err := q.GetOrganizerSoloEventParticipantListQuery(ctx, conn, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ORGANIZER-ERROR]: Failed to fetch participant list", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Participant list fetched successfully",
		"participants": participants,
	})
	pkg.Log.SuccessCtx(c)
}
