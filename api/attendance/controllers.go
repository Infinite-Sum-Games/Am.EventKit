package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func FetchEventsByOrganizer(c *gin.Context) {

	organizerIdStr, ok := pkg.GrabUserId(c, "ATTENDANCE")
	if !ok {
		return
	}
	organizerId, ok := pkg.GrabUuid(c, organizerIdStr, "ATTENDANCE", "organizerId")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ATTENDANCE") {
		return
	}
	defer conn.Release()

	q := db.New()

	events, err := q.FetchEventsByOrganizerQuery(ctx, conn, organizerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Failed to fetch events by organizer ID", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Events fetched successfully",
		"events":  events,
	})
	pkg.Log.SuccessCtx(c)
}

func FetchParticipantsByEvent(c *gin.Context) {

	eventIdStr := c.Param("eventId")

	eventId, ok := pkg.GrabUuid(c, eventIdStr, "ATTENDANCE", "eventId")

	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ATTENDANCE") {
		return
	}
	defer conn.Release()

	q := db.New()

	participants, err := q.FetchParticipantsByEventQuery(ctx, conn, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Failed to fetch events by organizer ID", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Participants list fetched successfully",
		"participants": participants,
	})
	pkg.Log.SuccessCtx(c)
}

func MarkSoloCheckInOutBoth(c *gin.Context) {
	key := c.Param("key")
	studentIdStr := c.Param("studentId")
	scheduleIdStr := c.Param("scheduleId")

	studentId, ok := pkg.GrabUuid(c, studentIdStr, "ATTENDANCE", "studentId")
	if !ok {
		return
	}

	scheduleId, ok := pkg.GrabUuid(c, scheduleIdStr, "ATTENDANCE", "scheduleId")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ATTENDANCE") {
		return
	}
	defer conn.Release()

	q := db.New()

	switch key {
	case "IN":
		rowsAffected, err := q.MarkSoloCheckInQuery(ctx, conn, db.MarkSoloCheckInQueryParams{
			StudentID:       studentId,
			EventScheduleID: scheduleId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Failed to mark check-in", err)
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Check-in could not be marked. Possible reasons: already checked in or invalid student/schedule ID",
			})
			pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Check-in could not be marked", nil)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Check-in marked successfully",
		})
		pkg.Log.SuccessCtx(c)
		return

	case "OUT":
		rowsAffected, err := q.MarkSoloCheckOutQuery(ctx, conn, db.MarkSoloCheckOutQueryParams{
			StudentID:       studentId,
			EventScheduleID: scheduleId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Failed to mark check-out", err)
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Check-out could not be marked. Possible reasons: not checked in yet or invalid student/schedule ID",
			})
			pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Check-out could not be marked", nil)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Check-out marked successfully",
		})
		pkg.Log.SuccessCtx(c)
		return

	case "BOTH":
		rowsAffected, err := q.MarkSoloBothQuery(ctx, conn, db.MarkSoloBothQueryParams{
			StudentID:       studentId,
			EventScheduleID: scheduleId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Failed to mark check-in and check-out", err)
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Check-in and Check-out could not be marked. Possible reasons: already checked in/out or invalid student/schedule ID",
			})
			pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Check-in and Check-out could not be marked", nil)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Check-in and Check-out marked successfully",
		})
		pkg.Log.SuccessCtx(c)
		return

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid key provided. Use 'IN', 'OUT', or 'BOTH'.",
		})
		pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Invalid key provided", nil)
		return
	}
}

func MarkTeamCheckInOutBoth(c *gin.Context) {
	key := c.Param("key")
	studentIdStr := c.Param("studentId")
	scheduleIdStr := c.Param("scheduleId")

	studentId, ok := pkg.GrabUuid(c, studentIdStr, "ATTENDANCE", "studentId")
	if !ok {
		return
	}

	scheduleId, ok := pkg.GrabUuid(c, scheduleIdStr, "ATTENDANCE", "scheduleId")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ATTENDANCE") {
		return
	}
	defer conn.Release()

	q := db.New()

	switch key {
	case "IN":
		rowAffected, err := q.MarkTeamCheckInQuery(ctx, conn, db.MarkTeamCheckInQueryParams{
			StudentID:       studentId,
			EventScheduleID: scheduleId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Failed to mark student check-in", err)
			return
		}
		if rowAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Student Check-in could not be marked. Possible reasons: already checked in or invalid student/schedule ID",
			})
			pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Student Check-in could not be marked", nil)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Student Check-in marked successfully",
		})
		pkg.Log.SuccessCtx(c)
		return

	case "OUT":
		rowAffected, err := q.MarkTeamCheckOutQuery(ctx, conn, db.MarkTeamCheckOutQueryParams{
			StudentID:       studentId,
			EventScheduleID: scheduleId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Failed to mark student check-out", err)
			return
		}
		if rowAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Student Check-out could not be marked. Possible reasons: not checked in yet or invalid student/schedule ID",
			})
			pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Student Check-out could not be marked", nil)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Student Check-out marked successfully",
		})
		pkg.Log.SuccessCtx(c)
		return

	case "BOTH":
		rowAffected, err := q.MarkTeamBothQuery(ctx, conn, db.MarkTeamBothQueryParams{
			StudentID:       studentId,
			EventScheduleID: scheduleId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Failed to mark student check-in and check-out", err)
			return
		}
		if rowAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Student Check-in and Check-out could not be marked. Possible reasons: already checked in/out or invalid student/schedule ID",
			})
			pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Student Check-in and Check-out could not be marked", nil)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Student Check-in and Check-out marked successfully",
		})
		pkg.Log.SuccessCtx(c)
		return

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid key provided. Use 'IN', 'OUT', or 'BOTH'.",
		})
		pkg.Log.ErrorCtx(c, "[ATTENDANCE-ERROR]: Invalid key provided", nil)
		return
	}
}

func UnMarkSoloCheckInOutBoth(c *gin.Context) {

}

func FetchEventDetailsByDateAndOrganizer(c *gin.Context) {
	// dateStr := c.Query("date")
	// organizer := c.Query("organizer")
	//
	// if dateStr == "" || organizer == "" {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	// 		"message": "date and organizer are required",
	// 	})
	//
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Missing date or organizer", nil)
	// 	return
	// }
	//
	// eventDate, err := time.Parse("2006-01-02", dateStr)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid date format, expected YYYY-MM-DD"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid date format", err)
	// 	return
	// }
	//
	// eventPgDate := pgtype.Date{}
	// _ = eventPgDate.Scan(eventDate)
	//
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	//
	// conn, err := cmd.DBPool.Acquire(ctx)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to acquire DB connection", err)
	// 	return
	// }
	// defer conn.Release()
	//
	// q := db.New()
	// events, err := q.GetEventsByDateAndOrganizer(ctx, conn, db.GetEventsByDateAndOrganizerParams{
	// 	EventDate: eventPgDate,
	// 	Abbr:      organizer,
	// })
	//
	// if err == pgx.ErrNoRows {
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"message": "No events found",
	// 	})
	//
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: No events with for date and organizer", err)
	// 	return
	// }
	//
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to fetch events", err)
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Events fetched successfully",
		// "events":  events,
	})

	pkg.Log.SuccessCtx(c)
}

// If both check-in and check-out is to be marked in one go, then
// utilize this controller, otherwise there are specific controllers
func MarkOneTimeAttendance(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "One-time attendance marked successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func MarkCheckIn(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Check-in marked successfully",
	})
	pkg.Log.SuccessCtx(c)

}

func MarkCheckOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Check-out marked successfully",
	})
	pkg.Log.SuccessCtx(c)
}

/*
This function
	Assumes:- Email, EventID, ScheduleID, Flag(checkin/checkout) from frontend
	Checks:- IsStudentRegisteredForAnokha, IsEventExists, IsStudentRegisteredForEvent, IsEventStartedAlready/IsEventEndedAlready
*/

func HandleEventCheckInOut(c *gin.Context) {
	// email := c.Query("email")
	// eventIdStr := c.Query("event_id")
	// scheduleIdStr := c.Query("schedule_id")
	// flag := c.Query("flag")
	//
	// if email == "" || eventIdStr == "" || scheduleIdStr == "" || flag == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "email, event_id, schedule_id and flag are required"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Missing required parameters", nil)
	// 	return
	// }
	//
	// if flag != "checkin" && flag != "checkout" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "flag must be either 'checkin' or 'checkout'"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid flag provided", nil)
	// 	return
	// }
	//
	// eventId, err := uuid.Parse(eventIdStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "invalid event_id"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid event_id", err)
	// 	return
	// }
	//
	// scheduleId, err := uuid.Parse(scheduleIdStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "invalid event_id"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Invalid schedule_id", err)
	// 	return
	// }
	//
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	//
	// conn, err := cmd.DBPool.Acquire(ctx)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to acquire DB connection"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to acquire DB connection", err)
	// 	return
	// }
	// defer conn.Release()
	//
	// q := db.New()
	//
	// student, err := q.GetStudentByEmail(ctx, conn, email)
	// if err == pgx.ErrNoRows {
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "Student not registered"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Student not registered", err)
	// 	return
	// } else if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error in student lookup"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Database error in student lookup", err)
	// 	return
	// }
	//
	// _, err = q.GetEventByIdQuery(ctx, conn, eventId)
	// if err == pgx.ErrNoRows {
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "Event does not exist"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Event does not exist", err)
	// 	return
	// } else if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error while fetching event"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Error fetching event", err)
	// 	return
	// }
	//
	// booking, err := q.CheckStudentRegisteredForEvent(ctx, conn,
	// 	db.CheckStudentRegisteredForEventParams{
	// 		StudentID: student.ID,
	// 		EventID:   eventId,
	// 	},
	// )
	//
	// if err == pgx.ErrNoRows {
	// 	c.JSON(http.StatusForbidden, gin.H{"message": "student is not registered for this event"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Student not registered for event", err)
	// 	return
	// }
	//
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "error checking registration"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Registration lookup failed", err)
	// 	return
	// }
	//
	// if booking.TxnStatus != "SUCCESS" {
	// 	c.JSON(http.StatusForbidden, gin.H{"message": "booking not confirmed"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Booking not confirmed", nil)
	// 	return
	// }
	//
	// schedule, err := q.GetScheduleById(ctx, conn, scheduleId)
	// if err == pgx.ErrNoRows {
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "schedule not found"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Schedule not found", err)
	// 	return
	// }
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "error fetching schedule"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Error fetching schedule", err)
	// 	return
	// }
	//
	// now := time.Now()
	// startTime := schedule.StartTime.Time
	// endTime := schedule.EndTime.Time
	//
	// if now.Before(startTime) {
	// 	c.JSON(http.StatusForbidden, gin.H{"message": "event not started yet"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Event not started", nil)
	// 	return
	// }
	//
	// if now.After(endTime) {
	// 	c.JSON(http.StatusForbidden, gin.H{"message": "event has already ended"})
	// 	pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Event already ended", nil)
	// 	return
	// }
	//
	// attendance, err := q.GetAttendanceRecord(ctx, conn, db.GetAttendanceRecordParams{StudentID: student.ID, EventScheduleID: scheduleId})
	//
	// if flag == "checkin" {
	//
	// 	if err == nil && attendance.CheckIn.Valid {
	// 		c.JSON(http.StatusBadRequest, gin.H{"message": "already checked in"})
	// 		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Already checked in", nil)
	// 		return
	// 	}
	//
	// 	_, err = q.InsertCheckIn(ctx, conn, db.InsertCheckInParams{
	// 		StudentID:       student.ID,
	// 		EventID:         eventId,
	// 		EventScheduleID: scheduleId,
	// 		BookingID:       booking.ID,
	// 		StudentName:     student.Name,
	// 		StudentEmail:    student.Email,
	// 	})
	//
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to check in"})
	// 		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to check in", err)
	// 		return
	// 	}
	//
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "check-in successful",
	// 	})
	//
	// 	pkg.Log.SuccessCtx(c)
	// 	return
	// }
	//
	// if flag == "checkout" {
	//
	// 	if err == pgx.ErrNoRows {
	// 		c.JSON(http.StatusBadRequest, gin.H{"message": "cannot check out without checking in"})
	// 		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Checkout attempted without check-in", err)
	// 		return
	// 	}
	//
	// 	if attendance.CheckOut.Valid {
	// 		c.JSON(http.StatusBadRequest, gin.H{"message": "already checked out"})
	// 		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Already checked out", nil)
	// 		return
	// 	}
	//
	// 	_, err = q.UpdateCheckOut(ctx, conn, attendance.ID)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to check out"})
	// 		pkg.Log.ErrorCtx(c, "[EVENT-ERROR]: Failed to check out", err)
	// 		return
	// 	}
	//
	// 	c.JSON(http.StatusOK, gin.H{"message": "checkout successful"})
	// 	pkg.Log.SuccessCtx(c)
	// 	return
	// }
}
