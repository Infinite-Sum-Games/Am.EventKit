package api

import (
	"context"
	"math"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func BookEvent(c *gin.Context) {
	// if it is group event, the email is considered as leader's email, we can keep the same naming convention for solo event too
	leaderEmail := c.GetString("email")
	if leaderEmail == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Email is not found after auth middleware", nil)
		return
	}

	eventIdStr := c.Param("eventId")
	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Request not processed due to invalid parameters",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Invalid event ID", err)
		return
	}

	// longer context timeout for booking
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != pgx.ErrTxClosed {
			pkg.Log.FatalCtx(c, "[BOOKING-FATAL]: Failed to rollback", rbErr)
		}
	}()

	q := db.New()

	// Get event details
	event, err := q.GetEventForBooking(ctx, tx, eventId)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Event not found",
			})
			pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Booking request sent for invalid event", nil)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to get event details", err)
		return
	}

	// Check if event is active
	if event.EventStatus != "ACTIVE" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Booking for this event is not allowed as event is inactive",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to book seat due to event unavailability", err)
		return
	}

	var req models.TeamBookingRequest
	isGroupEvent := event.IsGroup

	// TODO: string{leaderEmail} is written asuming that frontend doesnt add the leader details in the team array
	allMembers := []string{leaderEmail}

	if isGroupEvent {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Team details are not proper", err)
			return
		}
		// Creating the group list
		for _, m := range req.TeamMembers {
			allMembers = append(allMembers, m.StudentEmail)
		}
		// Validating team size
		if len(allMembers) < int(event.MinTeamsize.Int32) || len(allMembers) > int(event.MaxTeamsize.Int32) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Team size does not meet event requirements.",
			})
			pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Team details are not proper", err)
			return
		}
	}

	// Check for seats
	if event.SeatsFilled+int32(len(allMembers)) > event.TotalSeats {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not enough seats available.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Not enough seats for team", err)
		return
	}

	priceVal, err := event.Price.Float64Value()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Unable to conver price to float", err)
		return
	}
	price := priceVal.Float64

	var totalFee float64
	if event.IsPerHead {
		totalFee = price * float64(len(allMembers))
	} else {
		totalFee = price
	}

	totalFee = totalFee + (totalFee * 0.18)

	// TODO: Should I take ceil for each multiplication or final number or finally after adding tax?
	totalFee = math.Ceil(totalFee)

	// Checking if all users are registered and not already booked
	for _, memberEmail := range allMembers {
		student, err := q.GetStudentByEmail(ctx, tx, memberEmail)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Student not registered in anokha: " + memberEmail,
				})
				pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Student not registered in anokha", err)
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})
			pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Unable to check student existence", err)
			return
		}

		// Checking if each user is leader for existing booking / same query used for solo event registration
		_, err = q.GetBookingByUserAndEvent(ctx, tx, db.GetBookingByUserAndEventParams{
			StudentID: student.ID,
			EventID:   eventId,
		})
		if err != nil && err != pgx.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})
			pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Cannot check if member present in booking", err)
			return
		}
		// not pgx.ErrNoRows, which means there is an entry
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"message": "You are already registered for this event: " + memberEmail,
			})
			pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Student already registered for event", nil)
			return
		}

		// Check for existing team booking
		_, err = q.GetTeamBookingByUserAndEvent(ctx, tx,
			db.GetTeamBookingByUserAndEventParams{
				StudentID: student.ID,
				EventID:   eventId,
			})
		if err != nil && err != pgx.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})
			pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Cannot check if member present in some team", err)
			return
		}
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"message": "You are already registered for this event in a team: " + memberEmail,
			})
			pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Student already present in some team of this event", nil)
			return
		}
	}

	leaderIdString := c.GetString("userId")
	leaderId, err := uuid.Parse(leaderIdString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Invalid event ID", err)
		return
	}

	registrationFee, err := pkg.ToPgNumericFromFloat(totalFee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to convert float to numeric", err)
		return
	}
	bookingID, err := q.CreateBooking(ctx, tx, db.CreateBookingParams{
		EventID:         eventId,
		StudentID:       leaderId,
		TxnID:           pkg.GenerateTxnID(leaderId, eventId),
		RegistrationFee: registrationFee,
		TxnStatus:       "PENDING",
		ProductInfo:     "Booking for " + eventId.String(), //TODO: I forgot what to put here?
		SeatsReleased:   int32(len(allMembers)),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to create booking", err)
		return
	}

	if isGroupEvent {
		teamID, err := q.CreateTeam(ctx, tx, db.CreateTeamParams{
			TeamName:   req.TeamName,
			EventID:    eventId,
			LeaderName: leaderEmail, // TODO: Should i fetch leader details using query again?
			BookingID:  bookingID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})
			pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to create team", err)
			return
		}

		for _, team_members := range req.TeamMembers {
			member, err := q.GetStudentByEmail(ctx, tx, team_members.StudentEmail)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later.",
				})
				pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to fetch team member during team creation", err)
				return
			}
			_, err = q.CreateTeamMember(ctx, tx, db.CreateTeamMemberParams{
				TeamID:       teamID,
				StudentID:    member.ID,
				StudentRole:  team_members.StudentRole,
				StudentName:  member.Name,
				StudentEmail: team_members.StudentEmail,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later.",
				})
				pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to create team member", err)
				return
			}
		}
	}

	err = q.UpdateEventSeats(ctx, tx, db.UpdateEventSeatsParams{
		SeatsFilled: int32(len(allMembers)),
		ID:          eventId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to update event seats", err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.FatalCtx(c, "[BOOKING-FATAL]: Failed to commit transaction", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking successful! Please complete the payment."})
}
