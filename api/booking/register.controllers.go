package api

import (
	"context"
	"math"
	"net/http"
	"time"

	"fmt"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func BookEventCsrf(c *gin.Context) {
	csrfToken, err := pkg.CreateCsrfToken("event.book@amrita.edu", c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Cannot create csrf token", err)
		return
	}

	pkg.SetCsrfCookie(c, csrfToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Event booking action initiated successfully",
		"key":     csrfToken,
	})
	pkg.Log.SuccessCtx(c)
}

func BookEvent(c *gin.Context) {
	// if it is group event, the email is considered as leader's email,
	// we can keep the same naming convention for solo event too
	leaderEmail, ok := pkg.GrabEmail(c, "BOOKING")
	if !ok {
		return
	}

	eventIdStr := c.Param("eventId")
	eventId, ok := pkg.GrabUuid(c, eventIdStr, "BOOKING", "event")
	if !ok {
		return
	}

	// longer context timeout for booking
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "BOOKING") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "BOOKING")

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
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Event unavailable", err)
		return
	}

	var req models.TeamBookingRequest
	isGroupEvent := event.IsGroup

	// TODO: string{leaderEmail} is written asuming that frontend doesnt
	// add the leader details in the team array
	allMembers := []string{leaderEmail}
	emailCount := make(map[string]int)
	hasDuplicates := false

	if isGroupEvent {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Team details are not proper", err)
			return
		}
		if err := req.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request is malformed",
			})
			return
		}
		// Creating the group list
		for _, m := range req.TeamMembers {
			allMembers = append(allMembers, m.StudentEmail)
		}

		// Checking if duplicate emails are present in team details
		for _, email := range allMembers {
			emailCount[email]++
			if emailCount[email] > 1 {
				hasDuplicates = true
				break
			}
		}
		if hasDuplicates {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Duplicate team members found",
			})
			pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Duplicate team details found", err)
			return
		}
		// Validating team size
		lesser := len(allMembers) < int(event.MinTeamsize.Int32)
		greater := len(allMembers) > int(event.MaxTeamsize.Int32)
		if lesser || greater {
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
	totalFee = math.Ceil(totalFee)

	// Fetching all student details for the team
	students, err := q.GetStudentsByEmails(ctx, tx, allMembers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Unable to get student/member detail from email", err)
		return
	}
	studentMap := make(map[uuid.UUID]db.Student)
	for _, s := range students {
		studentMap[s.ID] = s
	}
	emailToId := make(map[string]uuid.UUID)
	for _, s := range students {
		emailToId[s.Email] = s.ID
	}

	// Checking if someone is not registered in anokha
	if len(students) != len(allMembers) {
		missing := ""
		for _, email := range allMembers {
			if _, ok := emailToId[email]; !ok {
				missing = email
				break
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Student not registered: " + missing,
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Student not registered in anokha"+missing, err)
		return
	}

	// Retriving the tag details of the event to check for specials
	// Convert interface{} → []string
	var specialTags []string
	if event.SpecialTags != nil {
		switch v := event.SpecialTags.(type) {
		case []any:
			for _, raw := range v {
				if s, ok := raw.(string); ok {
					specialTags = append(specialTags, s)
				}
			}
		case []string:
			specialTags = v
		}
	}

	var metadataJson []byte
	if len(specialTags) != 0 {
		metadataJson, err = pkg.BuildSpecialTags(specialTags, students)
		if err != nil {
			pkg.Log.ErrorCtx(c, "", err) // Empty message as the BuilSpecialTags sends proper errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})
			return
		}
	}

	var ids []uuid.UUID
	for _, s := range students {
		ids = append(ids, s.ID)
	}

	existing, err := q.GetAnyBookingByUsersAndEvent(ctx, tx, db.GetAnyBookingByUsersAndEventParams{
		Column1: ids,
		EventID: eventId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Unable to check bookings of user and event", err)
		return
	}
	if len(existing) > 0 {
		conflictId := existing[0]
		conflictEmail := studentMap[conflictId].Email

		c.JSON(http.StatusConflict, gin.H{
			"message": "User already registered: " + conflictEmail,
		})
		return
	}

	// Geting leader details
	leaderStrcut, err := q.GetStudentByEmail(ctx, tx, leaderEmail)

	leaderId := leaderStrcut.ID
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Invalid user ID", err)
		return
	}

	// Checking if the leader has any pending transactions across all events
	pendingBookings, err := q.GetAnyPendingBookingByUser(ctx, tx, leaderId)
	if err != nil && err != pgx.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Could not fetch pending events of leader", err)
		return
	}
	if err == nil && len(pendingBookings) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User has pending registrations for other events!",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: User has pending registrations for other events!", nil)
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
	txnId := pkg.GenerateTxnID(leaderId, eventId)
	prodInfo := fmt.Sprintf(
		"ERI-%s-%s-%d-%.2f",
		leaderId.String(),
		eventId.String(),
		len(allMembers),
		totalFee,
	)
	bookingID, err := q.CreateBooking(ctx, tx, db.CreateBookingParams{
		EventID:         eventId,
		StudentID:       leaderId,
		TxnID:           txnId,
		RegistrationFee: registrationFee,
		TxnStatus:       models.PaymentPending,
		ProductInfo:     prodInfo,
		SeatsReleased:   int32(len(allMembers)),
		Metadata:        metadataJson, // TODO: I need to set it as default data of the jsonb if not present
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to create booking", err)
		return
	}

	// Adding metadata for team if needed
	meta := pkg.NewJSONB()
	meta.Add("problem_stmt", req.ProblemStmt)
	metadata, err := meta.Bytes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to create metadata for team", err)
		return
	}

	if isGroupEvent {
		teamID, err := q.CreateTeam(ctx, tx, db.CreateTeamParams{
			TeamName:   req.TeamName,
			EventID:    eventId,
			LeaderName: leaderStrcut.Name,
			BookingID:  bookingID,
			Metadata:   metadata,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})
			pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to create team", err)
			return
		}

		// req.TeamMembers contains only details of team members, not leader
		for _, team_members := range req.TeamMembers {
			id := emailToId[team_members.StudentEmail]
			details := studentMap[id]
			_, err = q.CreateTeamMember(ctx, tx, db.CreateTeamMemberParams{
				TeamID:       teamID,
				StudentID:    id,
				StudentRole:  team_members.StudentRole,
				StudentName:  details.Name,
				StudentEmail: details.Email,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later.",
				})
				pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to create team member", err)
				return
			}
		}

		// Adding leader details into the table
		_, err = q.CreateTeamMember(ctx, tx, db.CreateTeamMemberParams{
			TeamID:       teamID,
			StudentID:    leaderId,
			StudentRole:  "leader",
			StudentName:  leaderStrcut.Name,
			StudentEmail: leaderStrcut.Email,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})
			pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to create leader team member", err)
			return
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

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "BOOKING") {
		return
	}

	// Generating the hash
	// TODO: Check salt
	hashedData := pkg.GenerateSHA512Hash(
		txnId,
		leaderEmail,
		fmt.Sprintf("%.2f", totalFee),
		prodInfo,
		leaderStrcut.Name,
	)

	c.JSON(http.StatusOK, gin.H{
		"message":         "Booking successful! Please complete the payment.",
		"txnId":           txnId,
		"name":            leaderStrcut.Name,
		"phone":           leaderStrcut.PhoneNumber,
		"registrationFee": totalFee,
		"productInfo":     prodInfo,
		"userEmail":       leaderEmail,
		"hash":            hashedData,
	})

}
