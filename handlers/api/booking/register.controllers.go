package api

//
// import (
// 	"context"
// 	"fmt"
// 	"math"
// 	"net/http"
// 	"time"
//
// 	"github.com/Infinite-Sum-Games/Am.EventKit/cmd"
// 	db "github.com/Infinite-Sum-Games/Am.EventKit/internal/datarepo/gen"
// 	pkg "github.com/Infinite-Sum-Games/Am.EventKit/internal/helpers"
// 	mq "github.com/Infinite-Sum-Games/Am.EventKit/internal/services/messagequeue"
// 	"github.com/Infinite-Sum-Games/Am.EventKit/logger"
// 	"github.com/Infinite-Sum-Games/Am.EventKit/models"
// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/jackc/pgx/v5"
// )
//
// func BookEventCsrf(c *gin.Context) {
// 	csrfToken, err := pkg.CreateCsrfToken("event.book@amrita.edu", c)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later.",
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Cannot create csrf token", err)
// 		return
// 	}
//
// 	pkg.SetCsrfCookie(c, csrfToken)
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Event booking action initiated successfully",
// 		"key":     csrfToken,
// 	})
// 	logger.Log.SuccessCtx(c)
// }
//
// func BookEvent(c *gin.Context) {
// 	// if it is group event, the email is considered as leader's email,
// 	// we can keep the same naming convention for solo event too
// 	leaderEmail, ok := pkg.GrabEmail(c, "BOOKING")
// 	if !ok {
// 		return
// 	}
//
// 	eventIdStr := c.Param("eventId")
// 	eventId, ok := pkg.GrabUuid(c, eventIdStr, "BOOKING", "event")
// 	if !ok {
// 		return
// 	}
//
// 	// longer context timeout for booking
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
//
// 	tx, err := cmd.DBPool.Begin(ctx)
// 	if pkg.HandleDbTxnErr(c, err, "BOOKING") {
// 		return
// 	}
// 	defer pkg.RollbackTx(c, tx, ctx, "BOOKING")
//
// 	q := db.New()
//
// 	// Get event details
// 	event, err := q.GetEventForBooking(ctx, tx, eventId)
// 	if err == pgx.ErrNoRows {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"message": "Event not found",
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Booking request sent for invalid event", nil)
// 		return
// 	}
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later.",
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to get event details", err)
// 		return
// 	}
//
// 	// Check if event is active
// 	if event.EventStatus != "ACTIVE" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "Booking for this event is not allowed as event is inactive",
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Event unavailable", err)
// 		return
// 	}
//
// 	var req *models.TeamBookingRequest
// 	isGroupEvent := event.IsGroup
//
// 	// TODO: string{leaderEmail} is written asuming that frontend doesnt
// 	// add the leader details in the team array
// 	allMembers := []string{leaderEmail}
// 	emailCount := make(map[string]int)
// 	hasDuplicates := false
//
// 	if isGroupEvent {
// 		req, ok = pkg.ValidateRequest[models.TeamBookingRequest](c)
// 		if !ok {
// 			return
// 		}
//
// 		// Creating the group list and check for duplicates
// 		for _, m := range req.TeamMembers {
// 			allMembers = append(allMembers, m.StudentEmail)
// 			emailCount[m.StudentEmail]++
// 			if emailCount[m.StudentEmail] > 1 {
// 				hasDuplicates = true
// 				break
// 			}
// 		}
//
// 		if hasDuplicates {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"message": "Duplicate team members found",
// 			})
// 			logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Duplicate team details found", err)
// 			return
// 		}
//
// 		// Validating team size
// 		lesser := len(allMembers) < int(event.MinTeamsize.Int32)
// 		greater := len(allMembers) > int(event.MaxTeamsize.Int32)
// 		if lesser || greater {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"message": "Team size does not meet event requirements.",
// 			})
// 			logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Team details are not proper", err)
// 			return
// 		}
// 	}
//
// 	// Check for seats
// 	if event.SeatsFilled+1 > event.TotalSeats {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "Not enough seats available.",
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Not enough seats for team", err)
// 		return
// 	}
//
// 	totalFee := int32(0)
// 	// Calculating total fee
// 	if event.IsPerHead {
// 		totalFee = event.Price * int32(len(allMembers))
// 	} else {
// 		totalFee = event.Price
// 	}
//
// 	// Adding GST
// 	totalFeeWithGST := float64(totalFee) * 1.18
// 	totalFee = int32(math.Ceil(totalFeeWithGST))
//
// 	// Fetching all student details for the team
// 	students, err := q.GetStudentsByEmails(ctx, tx, allMembers)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later.",
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Unable to get student/member detail from email", err)
// 		return
// 	}
//
// 	studentMap := make(map[uuid.UUID]db.Student)
// 	emailToId := make(map[string]uuid.UUID)
// 	for _, s := range students {
// 		studentMap[s.ID] = s
// 		emailToId[s.Email] = s.ID
// 	}
//
// 	// Checking if someone is not registered in anokha
// 	if len(students) != len(allMembers) {
// 		missing := ""
// 		for _, email := range allMembers {
// 			if _, ok := emailToId[email]; !ok {
// 				missing = email
// 				break
// 			}
// 		}
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "Student not signed-up: " + missing,
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Student not registered in anokha"+missing, err)
// 		return
// 	}
//
// 	// Retriving the tag details of the event to check for specials
// 	// Convert interface{} → []string
// 	var specialTags []string
// 	if event.SpecialTags != nil {
// 		switch v := event.SpecialTags.(type) {
// 		case []any:
// 			for _, raw := range v {
// 				if s, ok := raw.(string); ok {
// 					specialTags = append(specialTags, s)
// 				}
// 			}
// 		case []string:
// 			specialTags = v
// 		}
// 	}
//
// 	// Log special tags
// 	if len(specialTags) > 0 {
// 		logger.Log.InfoCtx(c, fmt.Sprintf("[BOOKING-INFO]: Special tags fetched for event %s: %v", eventId, specialTags))
// 	} else {
// 		logger.Log.InfoCtx(c, fmt.Sprintf("[BOOKING-INFO]: No special tags for event %s", eventId))
// 	}
//
// 	// Switch case for special tags metadata
// 	metadataJson := []byte(`{}`)
// 	for _, tag := range specialTags {
// 		switch tag {
// 		case "!woc":
// 			// Create and send payload to message queue
// 			payload, err := mq.CreateWoCPayload(
// 				leaderEmail,
// 				studentMap[emailToId[leaderEmail]].Name,
// 				studentMap[emailToId[leaderEmail]].Password,
// 			)
// 			if err != nil {
// 				logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Unable to create WOC payload", err)
// 				c.JSON(http.StatusInternalServerError, gin.H{
// 					"message": "Oops! Something happened. Please try again later.",
// 				})
// 				return
// 			}
// 			logger.Log.InfoCtx(c, fmt.Sprintf("[BOOKING-INFO]: WOC payload constructed for %s", leaderEmail))
//
// 			// updating the metadata to include queue name and payload
// 			meta := pkg.NewJSONB()
// 			meta.Add("woc_payload", string(payload))
// 			metadataJson, err = meta.Bytes()
// 			if err != nil {
// 				logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Unable to create WOC metadata", err)
// 				c.JSON(http.StatusInternalServerError, gin.H{
// 					"message": "Oops! Something happened. Please try again later.",
// 				})
// 				return
// 			}
// 		case "!hackathon":
// 			// Create and send payload to message queue
// 			teamMembers, err := mq.BuildHackathonTeamMembers(students, leaderEmail)
// 			if err != nil {
// 				logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Unable to create team details payload", err)
// 				c.JSON(http.StatusInternalServerError, gin.H{
// 					"message": "Oops! Something happened. Please try again later.",
// 				})
// 				return
// 			}
// 			logger.Log.InfoCtx(c, "[BOOKING-INFO]: Hackathon team members constructed successfully")
//
// 			leader := studentMap[emailToId[leaderEmail]]
// 			problemStmt := ""
// 			if req.ProblemStmt != nil {
// 				problemStmt = *req.ProblemStmt
// 			}
// 			payloadBytes, err := mq.CreateHackathonPayload(
// 				mq.HackathonPayload{
// 					TeamName:          req.TeamName,
// 					LeaderName:        leader.Name,
// 					LeaderEmail:       leaderEmail,
// 					LeaderPhoneNumber: leader.PhoneNumber,
// 					LeaderCollegeName: leader.CollegeName,
// 					ProblemStatement:  problemStmt,
// 					TeamMembers:       teamMembers,
// 				},
// 			)
// 			if err != nil {
// 				logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Unable to create Hackathon payload", err)
// 				c.JSON(http.StatusInternalServerError, gin.H{
// 					"message": "Oops! Something happened. Please try again later.",
// 				})
// 				return
// 			}
// 			logger.Log.InfoCtx(c, "[BOOKING-INFO]: Hackathon payload constructed successfully")
//
// 			meta := pkg.NewJSONB()
// 			meta.Add("hackathon_payload", string(payloadBytes))
// 			metadataJson, err = meta.Bytes()
// 			if err != nil {
// 				logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Unable to create Hackathon metadata", err)
// 				c.JSON(http.StatusInternalServerError, gin.H{
// 					"message": "Oops! Something happened. Please try again later.",
// 				})
// 				return
// 			}
// 		}
// 	}
// 	logger.Log.InfoCtx(c, "[BOOKING-INFO]: Final metadata: "+string(metadataJson))
//
// 	var ids []uuid.UUID
// 	for _, s := range students {
// 		ids = append(ids, s.ID)
// 	}
//
// 	existing, err := q.GetAnyBookingByUsersAndEvent(ctx, tx, db.GetAnyBookingByUsersAndEventParams{
// 		Column1: ids,
// 		EventID: eventId,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later.",
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Unable to check bookings of user and event", err)
// 		return
// 	}
// 	if len(existing) > 0 {
// 		conflictId := existing[0]
// 		conflictEmail := studentMap[conflictId].Email
//
// 		c.JSON(http.StatusConflict, gin.H{
// 			"message": "User already registered: " + conflictEmail,
// 		})
// 		return
// 	}
//
// 	// Geting leader details
// 	leaderStrcut, err := q.GetStudentByEmail(ctx, tx, leaderEmail)
//
// 	leaderId := leaderStrcut.ID
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later.",
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Invalid user ID", err)
// 		return
// 	}
//
// 	// Checking if the leader has any pending transactions across all events
// 	pendingBookings, err := q.GetAnyPendingBookingByUser(ctx, tx, leaderId)
// 	if err != nil && err != pgx.ErrNoRows {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later.",
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Could not fetch pending events of leader", err)
// 		return
// 	}
// 	if err == nil && len(pendingBookings) > 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "User has pending registrations for other events!",
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: User has pending registrations for other events!", nil)
// 		return
// 	}
//
// 	txnId := pkg.GenerateTxnID()
// 	prodInfo := fmt.Sprintf(
// 		"ERI-%s-%s-%d-%d",
// 		leaderId.String(),
// 		eventId.String(),
// 		len(allMembers),
// 		totalFee,
// 	)
// 	bookingID, err := q.CreateBooking(ctx, tx, db.CreateBookingParams{
// 		EventID:         eventId,
// 		StudentID:       leaderId,
// 		TxnID:           txnId,
// 		RegistrationFee: totalFee,
// 		TxnStatus:       models.PaymentPending,
// 		ProductInfo:     prodInfo,
// 		SeatsReleased:   1,
// 		Metadata:        metadataJson, // TODO: I need to set it as default data of the jsonb if not present
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later.",
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to create booking", err)
// 		return
// 	}
//
// 	// Adding metadata for team if needed
// 	meta := pkg.NewJSONB()
// 	if req != nil {
// 		meta.Add("problem_stmt", req.ProblemStmt)
// 	}
// 	metadata, err := meta.Bytes()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later.",
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to create metadata for team", err)
// 		return
// 	}
//
// 	if isGroupEvent {
// 		teamID, err := q.CreateTeam(ctx, tx, db.CreateTeamParams{
// 			TeamName:   req.TeamName,
// 			EventID:    eventId,
// 			LeaderName: leaderStrcut.Name,
// 			BookingID:  bookingID,
// 			Metadata:   metadata,
// 		})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"message": "Oops! Something happened. Please try again later.",
// 			})
// 			logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to create team", err)
// 			return
// 		}
//
// 		// req.TeamMembers contains only details of team members, not leader
// 		for _, team_members := range req.TeamMembers {
// 			id := emailToId[team_members.StudentEmail]
// 			details := studentMap[id]
// 			_, err = q.CreateTeamMember(ctx, tx, db.CreateTeamMemberParams{
// 				TeamID:       teamID,
// 				StudentID:    id,
// 				StudentRole:  team_members.StudentRole,
// 				StudentName:  details.Name,
// 				StudentEmail: details.Email,
// 			})
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{
// 					"message": "Oops! Something happened. Please try again later.",
// 				})
// 				logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to create team member", err)
// 				return
// 			}
// 		}
//
// 		// Adding leader details into the table
// 		_, err = q.CreateTeamMember(ctx, tx, db.CreateTeamMemberParams{
// 			TeamID:       teamID,
// 			StudentID:    leaderId,
// 			StudentRole:  "leader",
// 			StudentName:  leaderStrcut.Name,
// 			StudentEmail: leaderStrcut.Email,
// 		})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"message": "Oops! Something happened. Please try again later.",
// 			})
// 			logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to create leader team member", err)
// 			return
// 		}
// 	}
//
// 	err = q.UpdateEventSeats(ctx, tx, db.UpdateEventSeatsParams{
// 		SeatsFilled: 1,
// 		ID:          eventId,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Oops! Something happened. Please try again later.",
// 		})
// 		logger.Log.ErrorCtx(c, "[BOOKING-ERROR]: Failed to update event seats", err)
// 		return
// 	}
//
// 	err = tx.Commit(ctx)
// 	if pkg.HandleDbTxnCommitErr(c, err, "BOOKING") {
// 		return
// 	}
//
// 	// Generating the hash
// 	hashedData := pkg.GenerateSHA512Hash(
// 		txnId,
// 		leaderEmail,
// 		fmt.Sprintf("%d", totalFee),
// 		prodInfo,
// 		leaderStrcut.Name,
// 	)
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"message":         "Booking successful! Please complete the payment.",
// 		"txnId":           txnId,
// 		"name":            leaderStrcut.Name,
// 		"phone":           leaderStrcut.PhoneNumber,
// 		"registrationFee": totalFee,
// 		"productInfo":     prodInfo,
// 		"userEmail":       leaderEmail,
// 		"hash":            hashedData,
// 	})
// 	logger.Log.SuccessCtx(c)
//
// }
