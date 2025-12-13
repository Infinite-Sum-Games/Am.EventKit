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
	result, err := q.NewUntitledEventQuery(ctx, conn, db.NewUntitledEventQueryParams{
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
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to create new event", err)
		return
	}

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
		"is_offline":       result.EventMode == db.EventModeEnumOFFLINE,
		"attendance_mode":  result.AttendanceMode,
		"is_published":     result.EventStatus == db.EventStatusEnumACTIVE,
		"people":           []string{},
		"organizers":       []string{},
		"tags":             []string{},
		"schedules":        []string{},
	})
	pkg.Log.SuccessCtx(c)
}

// Name, Blurb, Description, Rules
func AddEventDetails(c *gin.Context) {

}

// Poster URL
// func AddEventPoster(c *gin.Context) {
// 	req, ok := pkg.ValidateRequest[AddEventPosterRequest](c)
// 	if !ok {
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"event_id":   result.ID,
// 		"poster_url": result.ConverImageUrl.String,
// 	})
// }

// IsTeam, MinSize, MaxSize, Seats
func AddEventDimension(c *gin.Context) {

}

// EventType, Mark As Completed, EventMode, AttendanceType
func AddEventToggles(c *gin.Context) {
}

func ConnectEventAndOrganizer(c *gin.Context) {

}

func DisonnectEventAndOrganizer(c *gin.Context) {

}

func ConnectEventAndTags(c *gin.Context) {

}

func DisonnectEventAndTags(c *gin.Context) {

}

func AttachNewEventSchedule(c *gin.Context) {

}

func EditEventSchedule(c *gin.Context) {

}

func DeleteEventSchedule(c *gin.Context) {

}

func EditEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	eventID, ok := pkg.GrabUuid(c, c.Param("eventId"), "EVENT", "event")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.UpdateEventRequest](c)
	if !ok {
		return
	}

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "EVENT") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "EVENT")

	q := db.New()
	// 1) update base event
	rows, err := q.UpdateEventQuery(ctx, tx, db.UpdateEventQueryParams{
		ID:             eventID,
		Name:           req.Name,
		Blurb:          req.Blurb,
		Description:    req.Description,
		CoverImageUrl:  pkg.ToPgText(req.CoverImageURL),
		Price:          pgtype.Numeric{Valid: true, Int: big.NewInt(int64(req.Price))},
		IsPerHead:      req.IsPerHead,
		Rules:          req.Rules,
		EventType:      db.EventTypeEnum(req.EventType),
		IsGroup:        req.IsGroup,
		MaxTeamsize:    pkg.ToPgInt4(req.MaxTeamSize),
		MinTeamsize:    pkg.ToPgInt4(req.MinTeamSize),
		TotalSeats:     req.TotalSeats,
		SeatsFilled:    req.SeatsFilled,
		EventStatus:    db.EventStatusEnum(req.EventStatus),
		EventMode:      db.EventModeEnum(req.EventMode),
		AttendanceMode: db.AttendanceModeEnum(req.AttendanceMode),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to update event", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Event does not exist"})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Event does not exist for update", nil)
		return
	}

	// 2) clear old mappings
	clearFuncs := []func(context.Context) error{
		func(ctx context.Context) error { return q.DeleteEventSchedulesByEventIDQuery(ctx, tx, eventID) },
		func(ctx context.Context) error { return q.DeleteEventTagMappingsByEventIDQuery(ctx, tx, eventID) },
		func(ctx context.Context) error { return q.DeleteEventOrganizerMappingsByEventIDQuery(ctx, tx, eventID) },
		func(ctx context.Context) error { return q.DeletePeopleToEventMappingsByEventIDQuery(ctx, tx, eventID) },
	}
	for _, fn := range clearFuncs {
		if err := fn(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to clear event mappings", err)
			return
		}
	}

	// 3) re-insert schedules
	for _, sch := range req.Schedules {
		eventDate, _ := time.Parse("2006-01-02", sch.EventDate)
		startTime, _ := time.Parse(time.RFC3339, sch.StartTime)
		endTime, _ := time.Parse(time.RFC3339, sch.EndTime)

		if err := q.InsertEventScheduleQuery(ctx, tx, db.InsertEventScheduleQueryParams{
			EventID:   eventID,
			EventDate: pgtype.Date{Time: eventDate, Valid: true},
			StartTime: pgtype.Timestamp{Time: startTime, Valid: true},
			EndTime:   pgtype.Timestamp{Time: endTime, Valid: true},
			Venue:     sch.Venue,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to insert event schedule", err)
			return
		}
	}

	// 4) re-insert tag mappings
	for _, tagIDStr := range req.TagIDs {
		tagID, ok := pkg.GrabUuid(c, tagIDStr, "EVENT", "tag")
		if !ok {
			return
		}

		if err := q.InsertEventTagMappingQuery(ctx, tx, db.InsertEventTagMappingQueryParams{
			TagID:   tagID,
			EventID: eventID,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to insert tag mapping", err)
			return
		}
	}

	// 5) re-insert organizer mappings
	for _, orgIDStr := range req.OrganizerIDs {
		orgID, ok := pkg.GrabUuid(c, orgIDStr, "EVENT", "organizer")
		if !ok {
			return
		}

		if err := q.InsertEventOrganizerMappingQuery(ctx, tx, db.InsertEventOrganizerMappingQueryParams{
			EventID:     eventID,
			OrganizerID: orgID,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to insert organizer mapping", err)
			return
		}
	}

	// 6) re-insert people mappings
	for _, personIDStr := range req.PeopleIDs {
		personID, ok := pkg.GrabUuid(c, personIDStr, "EVENT", "person")
		if !ok {
			return
		}

		if err := q.InsertPeopleToEventMappingQuery(ctx, tx, db.InsertPeopleToEventMappingQueryParams{
			EventID:  eventID,
			PersonID: personID,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to insert people mapping", err)
			return
		}
	}

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "EVENT") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event updated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func DeleteEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	eventID, ok := pkg.GrabUuid(c, c.Param("eventId"), "EVENT", "event")
	if !ok {
		return
	}

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "EVENT") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "EVENT")

	q := db.New()
	// delete all mappings first
	clearFuncs := []func(context.Context) error{
		func(ctx context.Context) error { return q.DeleteEventSchedulesByEventIDQuery(ctx, tx, eventID) },
		func(ctx context.Context) error { return q.DeleteEventTagMappingsByEventIDQuery(ctx, tx, eventID) },
		func(ctx context.Context) error { return q.DeleteEventOrganizerMappingsByEventIDQuery(ctx, tx, eventID) },
		func(ctx context.Context) error { return q.DeletePeopleToEventMappingsByEventIDQuery(ctx, tx, eventID) },
	}
	for _, fn := range clearFuncs {
		if err := fn(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to clear event mappings for delete", err)
			return
		}
	}

	rows, err := q.DeleteEventQuery(ctx, tx, eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to delete event", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event does not exist",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Event does not exist for delete", nil)
		return
	}

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "EVENT") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func ToggleEventStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	eventID, ok := pkg.GrabUuid(c, c.Param("eventId"), "EVENT", "event")
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	rows, err := q.ToggleEventStatusQuery(ctx, conn, eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to toggle event status", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event does not exist",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Event does not exist for toggle", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event status toggled successfully",
	})
	pkg.Log.SuccessCtx(c)
}
