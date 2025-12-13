package api

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/segmentio/ksuid"
)

func NewEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.NewUntitledEventQuery(ctx, conn,
		db.NewUntitledEventQueryParams{
			Name:        fmt.Sprintf("Untitled %s", ksuid.New().String()),
			Blurb:       "",
			Description: "",
			Price: pgtype.Numeric{
				Valid: true,
				Int:   big.NewInt(int64(0)),
			},
			IsPerHead: true,
			Rules:     "",
			EventType: db.EventTypeEnumEVENT,
			IsGroup:   false,
			MinTeamsize: pgtype.Int4{
				Valid: true,
				Int32: 1,
			},
			MaxTeamsize: pgtype.Int4{
				Valid: true,
				Int32: 1,
			},
			TotalSeats:     0,
			EventStatus:    db.EventStatusEnumCLOSED,
			EventMode:      db.EventModeEnumOFFLINE,
			AttendanceMode: db.AttendanceModeEnumSOLO,
			IsTechnical: pgtype.Bool{
				Valid: true,
				Bool:  false,
			},
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to create new event", err)
		return
	}

	eventIsCompleted := result.EventStatus == db.EventStatusEnumCOMPLETED
	eventIsActive := result.EventStatus == db.EventStatusEnumACTIVE
	eventIsOffline := result.EventMode == db.EventModeEnumOFFLINE

	c.JSON(http.StatusOK, gin.H{
		"message":          "New event created successfully",
		"event_id":         result.ID.String(),
		"event_name":       result.Name,
		"blurb":            result.Blurb,
		"description":      result.Description,
		"poster_url":       result.CoverImageUrl.String,
		"price":            result.Price.Int,
		"pricing_per_head": result.IsPerHead,
		"rules":            result.Rules,
		"is_group":         result.IsGroup,
		"min_teamsize":     result.MinTeamsize.Int32,
		"max_teamsize":     result.MaxTeamsize.Int32,
		"seat_count":       result.TotalSeats,
		"event_type":       result.EventType,
		"is_technical":     result.IsTechnical,
		"is_offline":       eventIsOffline,
		"attendance_mode":  result.AttendanceMode,
		"is_published":     eventIsActive || eventIsCompleted,
		"event_status":     eventIsCompleted,
		"people":           []string{},
		"organizers":       []string{},
		"tags":             []string{},
		"schedules":        []string{},
	})
	pkg.Log.SuccessCtx(c)
}

// Name, Blurb, Description, Rules, Price
func AddEventDetails(c *gin.Context) {
	eventId, ok := pkg.GrabUuid(c, c.Param("eventId"), "ADMIN-EVENT", "Event")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.AddEventDetailsRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if !pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}

	q := db.New()
	result, err := q.AddEventDetailsQuery(ctx, conn, db.AddEventDetailsQueryParams{
		ID:          eventId,
		Name:        req.Name,
		Blurb:       req.Blurb,
		Description: req.Description,
		Rules:       req.Rules,
		Price:       req.Price,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to add event details", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Event updated successfully",
		"id":          result.ID.String(),
		"name":        result.Name,
		"blurb":       result.Blurb,
		"description": result.Description,
		"rules":       result.Rules,
		"price":       result.Price,
	})
	pkg.Log.SuccessCtx(c)
}

// Poster URL
func AddEventPoster(c *gin.Context) {
	eventId, ok := pkg.GrabUuid(c, c.Param("eventId"), "ADMIN-EVENT", "Event")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.AddEventPosterRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if !pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.AddEventPosterQuery(ctx, conn, db.AddEventPosterQueryParams{
		EventId:       eventId,
		CoverImageUrl: req.PosterUrl,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT]: Failed to update event poster", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Successfully updated event poster",
		"event_id":   result.ID,
		"poster_url": result.CoverImageUrl.String,
	})
	pkg.Log.SuccessCtx(c)
}

func DeleteEventPoster(c *gin.Context) {
	eventId, ok := pkg.GrabUuid(c, c.Param("eventId"), "ADMIN-EVENT", "Event")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if !pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()

	_, err = q.DeleteEventPosterQuery(ctx, conn, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT]: Failed to delete event poster", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted event poster",
	})
}

// IsTeam, MinSize, MaxSize, Seats
func AddEventDimension(c *gin.Context) {

}

// EventType, Mark As Completed, EventMode, AttendanceType
func AddEventToggles(c *gin.Context) {

}

func ConnectEventAndOrganizer(c *gin.Context) {

}

func DisconnectEventAndOrganizer(c *gin.Context) {

}

func ConnectEventAndTags(c *gin.Context) {

}

func DisonnectEventAndTags(c *gin.Context) {

}

func AddEventSchedule(c *gin.Context) {

}

func EditEventSchedule(c *gin.Context) {

}

func DeleteEventSchedule(c *gin.Context) {

}

func PublishEvent(c *gin.Context) {

}

func UnpublishEvent(c *gin.Context) {

}

func MarkEventAsCompleted(c *gin.Context) {
	eventId, ok := pkg.GrabUuid(c, "eventId", "ADMIN-EVENT", "Event")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if !pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.MarkEventAsCompletedQuery(ctx, conn, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: ", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Event marked as completed",
		"event_status": result,
	})
	pkg.Log.SuccessCtx(c)
}

func UnmarkEventAsCompleted(c *gin.Context) {
	eventId, ok := pkg.GrabUuid(c, "eventId", "ADMIN-EVENT", "Event")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if !pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.UnmarkEventsAsCompletedQuery(ctx, conn, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: ", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Event unmarked but is still in published state",
		"event_status": result,
	})
	pkg.Log.SuccessCtx(c)
}
