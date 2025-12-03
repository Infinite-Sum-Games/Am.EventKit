package api

import (
	"context"
	"math/big"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func CreateEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, ok := pkg.ValidateRequest[models.CreateEventRequest](c)
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	tx, err := conn.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to begin transaction", err)
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	eventID, err := q.CreateEventQuery(ctx, tx, db.CreateEventQueryParams{
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
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to create event", err)
		return
	}

	// 2) schedules
	for _, sch := range req.Schedules {
		if err := sch.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid schedule", err)
			return
		}

		eventDate, err := time.Parse("2006-01-02", sch.EventDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid event_date format", err)
			return
		}

		startTime, err := pkg.ParseTime(sch.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid start_time format", err)
			return
		}

		endTime, err := pkg.ParseTime(sch.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid end_time format", err)
			return
		}

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

	// 3) tag mappings
	for _, tagIDStr := range req.TagIDs {
		tagID, err := uuid.Parse(tagIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid tag ID", err)
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

	// 4) organizer mappings
	for _, orgIDStr := range req.OrganizerIDs {
		orgID, err := uuid.Parse(orgIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid organizer ID", err)
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

	// 5) people mappings
	for _, personIDStr := range req.PeopleIDs {
		personID, err := uuid.Parse(personIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid person ID", err)
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

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to commit transaction", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Event created successfully",
		"event_id": eventID,
	})
	pkg.Log.SuccessCtx(c)
}

func EditEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	eventIDStr := c.Param("eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Request is malformed",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid event ID format", err)
		return
	}

	req, ok := pkg.ValidateRequest[models.UpdateEventRequest](c)
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()
	tx, err := conn.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to begin transaction", err)
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

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
		if err := sch.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid schedule", err)
			return
		}

		eventDate, err := time.Parse("2006-01-02", sch.EventDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid event_date format", err)
			return
		}

		startTime, err := pkg.ParseTime(sch.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid start_time format", err)
			return
		}

		endTime, err := pkg.ParseTime(sch.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid end_time format", err)
			return
		}

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
		tagID, err := uuid.Parse(tagIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid tag ID", err)
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
		orgID, err := uuid.Parse(orgIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid organizer ID", err)
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
		personID, err := uuid.Parse(personIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid person ID", err)
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

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to commit transaction", err)
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

	eventIDStr := c.Param("eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Request is malformed",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid event ID format", err)
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()
	tx, err := conn.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to begin transaction", err)
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

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

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to commit transaction", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully",
	})
	pkg.Log.SuccessCtx(c)
}
