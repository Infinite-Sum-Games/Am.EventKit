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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/segmentio/ksuid"
)

func NewEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
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
		"message":         "New event created successfully",
		"id":              result.ID.String(),
		"name":            result.Name,
		"blurb":           result.Blurb,
		"description":     result.Description,
		"poster_url":      result.CoverImageUrl.String,
		"price":           result.Price.Int,
		"is_per_head":     result.IsPerHead,
		"rules":           result.Rules,
		"is_group":        result.IsGroup,
		"min_teamsize":    result.MinTeamsize.Int32,
		"max_teamsize":    result.MaxTeamsize.Int32,
		"seat_count":      result.TotalSeats,
		"event_type":      result.EventType,
		"attendance_mode": result.AttendanceMode,
		"is_technical":    result.IsTechnical,
		"is_offline":      eventIsOffline,
		"is_published":    eventIsActive || eventIsCompleted,
		"is_completed":    eventIsCompleted,
		"people":          []string{},
		"organizers":      []string{},
		"tags":            []string{},
		"schedules":       []string{},
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
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}

	q := db.New()
	result, err := q.AddEventDetailsQuery(ctx, conn, db.AddEventDetailsQueryParams{
		ID:          eventId,
		Name:        req.Name,
		Blurb:       req.Blurb,
		Description: req.Description,
		Rules:       req.Rules,
		Price: pgtype.Numeric{
			Int:   big.NewInt(int64(req.Price)),
			Valid: true,
		},
		IsPerHead: req.IsPerHead,
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
		"updated_at":  result.UpdatedAt,
	})
	pkg.Log.SuccessCtx(c)
}

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
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.AddEventPosterQuery(ctx, conn, db.AddEventPosterQueryParams{
		ID:            eventId,
		CoverImageUrl: pgtype.Text{String: req.PosterUrl, Valid: true},
	})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to update event poster", err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to update event poster", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Successfully updated event poster",
		"id":         result.ID,
		"poster_url": result.CoverImageUrl.String,
		"updated_at": result.UpdatedAt,
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
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	_, err = q.DeleteEventPosterQuery(ctx, conn, eventId)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to delete event poster", err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to delete event poster", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted event poster",
	})
	pkg.Log.SuccessCtx(c)
}

func AddEventDimension(c *gin.Context) {
	eventId, ok := pkg.GrabUuid(c, c.Param("eventId"), "ADMIN-EVENT", "Event")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.AddEventDimensionRequest](c)
	if !ok {
		return
	}

	if !req.IsGroup {
		req.MinTeamSize = 1
		req.MaxTeamSize = 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.AddEventDimensionQuery(ctx, conn, db.AddEventDimensionQueryParams{
		ID:         eventId,
		IsGroup:    req.IsGroup,
		TotalSeats: int32(req.TotalSeats),
		MinTeamsize: pgtype.Int4{
			Int32: int32(req.MinTeamSize),
			Valid: true,
		},
		MaxTeamsize: pgtype.Int4{
			Int32: int32(req.MaxTeamSize),
			Valid: true,
		},
	})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to add event dimensions", err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to add event dimensions", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Event dimensions updated",
		"total_seats":  result.TotalSeats,
		"is_group":     result.IsGroup,
		"min_teamsize": result.MinTeamsize.Int32,
		"max_teamsize": result.MaxTeamsize.Int32,
		"updated_at":   result.UpdatedAt,
	})
	pkg.Log.SuccessCtx(c)
}

// EventType, Mark As Completed, EventMode, AttendanceType
func AddEventToggles(c *gin.Context) {
	eventId, ok := pkg.GrabUuid(c, c.Param("eventId"), "ADMIN-EVENT", "Event")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.AddEventTogglesRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	// Logic
	var eventMode = db.EventModeEnumONLINE
	var attendanceMode = db.AttendanceModeEnumSOLO
	var eventType = db.EventTypeEnumEVENT
	var eventStatus = db.EventStatusEnumCLOSED

	if req.IsOffline {
		eventMode = db.EventModeEnumOFFLINE
	}
	if req.AttendanceMode == "DUO" {
		attendanceMode = db.AttendanceModeEnumDUO
	}
	if req.EventType == "WORKSHOP" {
		eventType = db.EventTypeEnumWORKSHOP
	}
	// Check for publishing before completed. Publishing means both
	// ACTIVE state and COMPLETED state
	if req.IsPublished {
		eventStatus = db.EventStatusEnumACTIVE
	}
	if req.IsCompleted {
		eventStatus = db.EventStatusEnumCOMPLETED
	}

	q := db.New()
	result, err := q.AddEventTogglesQuery(ctx, conn, db.AddEventTogglesQueryParams{
		ID:             eventId,
		EventType:      eventType, // EVENT | WORKSHOP
		EventMode:      eventMode, // IsOffline or not
		AttendanceMode: attendanceMode,
		IsTechnical:    pgtype.Bool{Bool: req.IsTechnical, Valid: true},
		EventStatus:    eventStatus, // PUBLISHED | COMPLETED | CLOSED
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to add event togglables", err)
		return
	}

	eventIsCompleted := result.EventStatus == db.EventStatusEnumCOMPLETED
	eventIsActive := result.EventStatus == db.EventStatusEnumACTIVE
	eventIsOffline := result.EventMode == db.EventModeEnumOFFLINE

	c.JSON(http.StatusOK, gin.H{
		"message":         "Event toggle metadata added successfully",
		"event_type":      result.EventType,      // WORKSHOP | EVENT
		"attendance_mode": result.AttendanceMode, // SOLO | DUO
		"is_offline":      eventIsOffline,
		"is_technical":    result.IsTechnical,
		"is_published":    eventIsCompleted || eventIsActive,
		"is_completed":    eventIsCompleted,
	})
	pkg.Log.SuccessCtx(c)
}

func ConnectEventAndOrganizer(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.ConnectEventAndOrganizerRequest](c)
	if !ok {
		return
	}

	eventId, ok := pkg.GrabUuid(c, req.EventId, "ADMIN-EVENT", "Event")
	if !ok {
		return
	}
	organizerId, ok := pkg.GrabUuid(c, req.OrganizerId, "ADMIN-EVENT", "Organizer")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.ConnectEventAndOrganizerQuery(ctx, conn,
		db.ConnectEventAndOrganizerQueryParams{
			EventID:     eventId,
			OrganizerID: organizerId,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to connect event and organizer", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Successfully connected event with organizer",
		"id":           result.EventID.String(),
		"organizer_id": result.OrganizerID.String(),
	})
	pkg.Log.SuccessCtx(c)
}

func DisconnectEventAndOrganizer(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.DisconnectEventAndOrganizerRequest](c)
	if !ok {
		return
	}

	eventId, ok := pkg.GrabUuid(c, req.EventId, "ADMIN-EVENT", "Event")
	if !ok {
		return
	}
	organizerId, ok := pkg.GrabUuid(c, req.OrganizerId, "ADMIN-EVENT", "Organizer")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	_, err = q.DisconnectEventAndOrganizerQuery(ctx, conn,
		db.DisconnectEventAndOrganizerQueryParams{
			EventID:     eventId,
			OrganizerID: organizerId,
		})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event and organizer mapping not found",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT]: Could not delete event and organizer mapping", err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to delete event and organizer mapping", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Removed event and organizer mapping",
	})
	pkg.Log.SuccessCtx(c)
}

func ConnectEventAndTags(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.ConnectEventAndTagsRequest](c)
	if !ok {
		return
	}

	eventId, ok := pkg.GrabUuid(c, req.EventId, "ADMIN-EVENT", "Event")
	if !ok {
		return
	}
	tagId, ok := pkg.GrabUuid(c, req.EventId, "ADMIN-EVENT", "Tag")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	fmt.Printf("%s")

	q := db.New()
	result, err := q.ConnectEventAndTagsQuery(ctx, conn, db.ConnectEventAndTagsQueryParams{
		EventID: eventId,
		TagID:   tagId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"messsage": "Oops! Something package. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to add tag for event", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tag added for event",
		"id":      result.EventID.String(),
		"tag_id":  result.TagID.String(),
	})
	pkg.Log.SuccessCtx(c)
}

func DisonnectEventAndTags(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.DisconnectEventAndTagsRequest](c)
	if !ok {
		return
	}

	eventId, ok := pkg.GrabUuid(c, req.EventId, "ADMIN-EVENT", "Event")
	if !ok {
		return
	}
	tagId, ok := pkg.GrabUuid(c, req.TagId, "ADMIN-EVENT", "Tag")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	_, err = q.DisconnectEventAndTagsQuery(ctx, conn,
		db.DisconnectEventAndTagsQueryParams{
			EventID: eventId,
			TagID:   tagId,
		})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No mapping found for given event and tag",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to disconnect event and tag", err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to disconnect event and tag", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully removed tag for event",
	})
	pkg.Log.SuccessCtx(c)
}

func ConnectEventAndPeople(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.ConnectEventAndPeopleRequest](c)
	if !ok {
		return
	}

	eventId, ok := pkg.GrabUuid(c, req.EventId, "ADMIN-EVENT", "Event")
	if !ok {
		return
	}
	personId, ok := pkg.GrabUuid(c, req.PersonId, "ADMIN-EVENT", "Person")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.ConnectEventAndPeopleQuery(ctx, conn,
		db.ConnectEventAndPeopleQueryParams{
			EventID:  eventId,
			PersonID: personId,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT]: Failed to connect event and people", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messasge":  "Successfully connected event and person",
		"id":        result.EventID.String(),
		"person_id": result.PersonID.String(),
	})
	pkg.Log.SuccessCtx(c)
}

func DisconnectEventAndPeople(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.DisconnectEventAndPeopleRequest](c)
	if !ok {
		return
	}

	eventId, ok := pkg.GrabUuid(c, req.EventId, "ADMIN-EVENT", "Event")
	if !ok {
		return
	}
	personId, ok := pkg.GrabUuid(c, req.PersonId, "ADMIN-EVENT", "Person")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-ERROR") {
		return
	}
	defer conn.Release()

	q := db.New()
	_, err = q.DisconnectEventAndPeopleQuery(ctx, conn,
		db.DisconnectEventAndPeopleQueryParams{
			EventID:  eventId,
			PersonID: personId,
		})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No mapping found for given event and person",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to disconnect event and people", err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to disconnect event and people", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully disconnected event and person",
	})
	pkg.Log.SuccessCtx(c)
}

func AddEventSchedule(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.AddEventScheduleRequest](c)
	if !ok {
		return
	}

	eventId, ok := pkg.GrabUuid(c, c.Param("eventId"), "ADMIN-EVENT", "Event")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.AddEventScheduleQuery(ctx, conn, db.AddEventScheduleQueryParams{
		EventID: eventId,
		EventDate: pgtype.Date{
			Valid: true,
		},
		StartTime: pgtype.Timestamp{
			Valid: true,
		},
		EndTime: pgtype.Timestamp{
			Valid: true,
		},
		Venue: req.Venue,
	})

	c.JSON(http.StatusOK, gin.H{
		"message":     "Successfully added event schedules",
		"schedule_id": result.ID.String(),
		"event_date":  result.EventDate,
		"start_time":  result.StartTime,
		"end_time":    result.EndTime,
		"venue":       result.Venue,
	})
	pkg.Log.SuccessCtx(c)
}

func EditEventSchedule(c *gin.Context) {
	scheduleId, ok := pkg.GrabUuid(c, c.Param("scheduleId"), "ADMIN-EVENT", "Schedule")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.EditEventScheduleRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.EditEventScheduleQuery(ctx, conn, db.EditEventScheduleQueryParams{
		ID: scheduleId,
		EventDate: pgtype.Date{
			Valid: true,
		},
		StartTime: pgtype.Timestamp{
			Valid: true,
		},
		EndTime: pgtype.Timestamp{
			Valid: true,
		},
		Venue: req.Venue,
	})

	c.JSON(http.StatusOK, gin.H{
		"message":     "Successfully updated event schedule",
		"schedule_id": scheduleId,
		"event_date":  result.EventDate,
		"start_time":  result.StartTime,
		"end_time":    result.EndTime,
		"venue":       result.Venue,
		"updated_at":  result.UpdatedAt,
	})
	pkg.Log.SuccessCtx(c)
}

func DeleteEventSchedule(c *gin.Context) {
	scheduleId, ok := pkg.GrabUuid(c, c.Param("scheduleId"), "ADMIN-EVENT", "Schedule")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	_, err = q.DeleteEventScheduleByIdQuery(ctx, conn, scheduleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT]: Failed to delete event schedule", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully removed event schedule",
	})
	pkg.Log.SuccessCtx(c)
}

func PublishEvent(c *gin.Context) {
	eventId, ok := pkg.GrabUuid(c, c.Param("eventId"), "ADMIN-EVENT", "Event")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.PublishEventQuery(ctx, conn, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to publish event", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Successfully published event",
		"event_status": result.EventStatus,
		"updated_at":   result.UpdatedAt,
	})
	pkg.Log.SuccessCtx(c)
}

func UnpublishEvent(c *gin.Context) {
	eventId, ok := pkg.GrabUuid(c, c.Param("eventId"), "ADMIN-EVENT", "Event")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.UnpublishEventQuery(ctx, conn, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Unpublished event successfully",
		"event_status": result.EventStatus,
		"updated_at":   result.UpdatedAt,
	})
	pkg.Log.SuccessCtx(c)
}

func MarkEventAsCompleted(c *gin.Context) {
	eventId, ok := pkg.GrabUuid(c, c.Param("eventId"), "ADMIN-EVENT", "Event")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.MarkEventAsCompletedQuery(ctx, conn, eventId)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to mark event as completed", err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to mark event as completed", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Event marked as completed",
		"event_status": result.EventStatus,
		"updated_at":   result.UpdatedAt,
	})
	pkg.Log.SuccessCtx(c)
}

func UnmarkEventAsCompleted(c *gin.Context) {
	eventId, ok := pkg.GrabUuid(c, c.Param("eventId"), "ADMIN-EVENT", "Event")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADMIN-EVENT") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.UnmarkEventsAsCompletedQuery(ctx, conn, eventId)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to unmark event as completed", err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ADMIN-EVENT-ERROR]: Failed to unmark event as completed", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Event unmarked but is still in published state",
		"event_status": result.EventStatus,
		"updated_at":   result.UpdatedAt,
	})
	pkg.Log.SuccessCtx(c)
}
