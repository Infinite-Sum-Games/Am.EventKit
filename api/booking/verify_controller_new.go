package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/mail"
	messagequeue "github.com/Thanus-Kumaar/anokha-2025-backend/message-queue"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func VerifyTransactionNew(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.VerifyTransactionRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "VERIFY") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "VERIFY")

	q := db.New()

	email, err := q.GetEmailByTxnId(ctx, tx, req.TxnID)
	if err != nil {
		pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to get email by txn id", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		return
	}

	booking, err := q.GetBookingByTxnID(ctx, tx, req.TxnID)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Transaction not found",
		})
		pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Requested transaction id is not found", err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to fetch booking", err)
		return
	}

	// If booking was already verified earlier, return current status
	if booking.TxnStatus != models.PaymentPending {
		c.JSON(http.StatusOK, gin.H{
			"message": "Already verified",
			"status":  booking.TxnStatus,
		})
		pkg.Log.SuccessCtx(c)
		return
	}

	// TODO: Call the PayU verify API here.
	// Possible values:
	// "success"
	// "failure"
	// "pending"
	// "not_found"
	formBody := pkg.BuildVerifyPayUForm(req.TxnID)
	httpReq, err := http.NewRequest(
		"POST",
		cmd.Env.PayUVerifyURL,
		strings.NewReader(formBody),
	)
	if err != nil {
		pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to build PayU request", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		return
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to contact PayU", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Payment Gateway Failed. Try again later", // Different msg
		})
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to close body", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			return
		}
	}()

	// Parsing the response
	var payURes models.PayUVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&payURes); err != nil {
		pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to decode PayU verify response", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Payment Gateway Failed. Try again later",
		})
		return
	}

	gatewayStatus := models.MapPayUStatus(payURes, req.TxnID)

	event, err := q.GetEventByIdQuery(ctx, tx, booking.EventID)
	if err != nil {
		pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to get event details", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		return
	}

	if gatewayStatus == models.PaymentFailed || gatewayStatus == models.PaymentNotFound {
		// Restoring the seats
		err = q.UpdateEventSeats(ctx, tx, db.UpdateEventSeatsParams{
			SeatsFilled: -booking.SeatsReleased,
			ID:          booking.EventID,
		})
		if err != nil {
			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to update event seats", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			return
		}

		// // If team event, then remove the team details
		// if event.IsGroup {
		// 	teamID, err := q.GetTeamIDByBooking(ctx, tx, booking.ID)
		// 	if err != nil && err != sql.ErrNoRows {
		// 		pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to get team ID", err)
		// 		c.JSON(http.StatusInternalServerError, gin.H{
		// 			"message": "Oops! Something happened. Please try again later",
		// 		})
		// 		return
		// 	}
		// 	if err == nil {
		// 		// delete members first due to foreign key relation
		// 		if err := q.DeleteTeamDetailsOfTeam(ctx, tx, teamID); err != nil {
		// 			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to delete team members", err)
		// 			c.JSON(http.StatusInternalServerError, gin.H{
		// 				"message": "Oops! Something happened. Please try again later",
		// 			})
		// 			return
		// 		}
		// 		// delete team
		// 		if err := q.DeleteTeam(ctx, tx, booking.ID); err != nil {
		// 			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to delete team", err)
		// 			c.JSON(http.StatusInternalServerError, gin.H{
		// 				"message": "Oops! Something happened. Please try again later",
		// 			})
		// 			return
		// 		}
		// 	}
		// }

		// Update booking status → failed
		err = q.UpdateBookingStatus(ctx, tx, db.UpdateBookingStatusParams{
			TxnStatus: models.PaymentFailed,
			ID:        booking.ID,
		})
		if err != nil {
			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to update booking status", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			return
		}
		// failure case
		err = tx.Commit(ctx)
		if pkg.HandleDbTxnCommitErr(c, err, "VERIFY") {
			return
		}

		pkg.Log.SuccessCtx(c)
		c.JSON(http.StatusOK, gin.H{
			"message": "Payment failed",
			"status":  models.PaymentFailed,
		})
		return
	}
	if gatewayStatus == models.PaymentSuccess {
		// Getting the schedule ids of the selected event
		schedules, err := q.GetSchedulesByEventID(ctx, tx, event.ID)
		if err != nil {
			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to get schedules in success", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			return
		}

		// Getting details of the student - needed for mail and attendance
		student, err := q.GetStudentByEmail(ctx, tx, email)
		if err != nil {
			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to retrive student data", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
		}

		if !event.IsGroup {
			for _, s := range schedules {
				_, err := q.CreateSoloEventParticipant(ctx, tx, db.CreateSoloEventParticipantParams{
					StudentID:       booking.StudentID,
					EventID:         booking.EventID,
					EventScheduleID: s,
					BookingID:       booking.ID,
					StudentName:     student.Name,
					StudentEmail:    student.Email,
				})
				if err != nil {
					pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to insert in solo participant", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "Oops! Something happened. Please try again later",
					})
					return
				}
			}
		} else {
			teamId, err := q.GetTeamIDByBooking(ctx, tx, booking.ID)
			if err != nil {
				pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to get team details in verify", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				return
			}
			members, err := q.GetTeamMembersByTeamID(ctx, tx, teamId)
			if err != nil {
				pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to fetch team members.", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				return
			}
			// Inserting into team attendance table
			for _, s := range schedules {
				for _, m := range members {
					_, err := q.CreateTeamAttendance(ctx, tx, db.CreateTeamAttendanceParams{
						StudentID:       m.StudentID,
						EventScheduleID: s,
					})
					if err != nil {
						pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to insert into team attd.", err)
						c.JSON(http.StatusInternalServerError, gin.H{
							"message": "Oops! Something happened. Please try again later",
						})
						return
					}
				}
			}

		}

		err = q.UpdateBookingStatus(ctx, tx, db.UpdateBookingStatusParams{
			TxnStatus: models.PaymentSuccess,
			ID:        booking.ID,
		})
		if err != nil {
			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to update booking status", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			return
		}

		// success case
		err = tx.Commit(ctx)
		if pkg.HandleDbTxnCommitErr(c, err, "VERIFY") {
			return
		}

		// If there is metadata, read and publish
		if len(booking.Metadata) > 0 {
			var metadataMap map[string]any
			if err := json.Unmarshal(booking.Metadata, &metadataMap); err != nil {
				pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to unmarshal booking metadata", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				return
			}

			// Hackathon payload
			if raw, ok := metadataMap["hackathon_payload"]; ok {
				payloadStr, ok := raw.(string)
				if !ok {
					pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: hackathon_payload is not string", nil)
					return
				}

				if err := messagequeue.Rabbit.Publish(
					ctx,
					messagequeue.QueueHackathonRegistrations,
					[]byte(payloadStr),
				); err != nil {
					pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to publish hackathon payload", err)
					return
				}
				// WOC payload
			} else if raw, ok := metadataMap["woc_payload"]; ok {
				payloadStr, ok := raw.(string)
				if !ok {
					pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: woc_payload is not string", nil)
					return
				}

				if err := messagequeue.Rabbit.Publish(
					ctx,
					messagequeue.QueueWocRegistrations,
					[]byte(payloadStr),
				); err != nil {
					pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to publish WOC payload", err)
					return
				}
			}
		}

		var completeSchedules []models.EventScheduleInput
		schedulesBytes, _ := json.Marshal(event.Schedules)
		_ = json.Unmarshal(schedulesBytes, &completeSchedules)

		var selected models.EventScheduleInput
		if len(completeSchedules) > 0 {
			selected = completeSchedules[0]
			for _, s := range completeSchedules {
				d1, _ := time.Parse("2006-01-02", s.EventDate)
				d2, _ := time.Parse("2006-01-02", selected.EventDate)
				if d1.Before(d2) {
					selected = s
				}
			}
		}

		err = mail.Mail.Enqueue(&mail.EmailRequest{
			To:      []string{},
			Subject: "Event Registration - Anokha 2026",
			Type:    "event-reg",
			Data: &mail.RegistrationData{
				UserName:      student.Name,
				EventName:     event.EventName,
				EventDate:     selected.EventDate,
				EventTime:     selected.StartTime + " - " + selected.EndTime,
				EventLocation: selected.Venue,
			},
			Retries: mail.MaxRetryCount,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[MAIL-ERROR]: Failed to add request to email queue", err)
			return
		}
		pkg.Log.SuccessCtx(c)
		c.JSON(http.StatusOK, gin.H{
			"message": "Payment verified successfully",
			"status":  models.PaymentSuccess,
		})
		return
	}

	// Still pending case will be handled here

	// Warning so that filtering the pending trasactions from logs will be easier
	pkg.Log.WarnCtx(c, "[VERIFY-WARN]: Verification successful, but status still pending for "+req.TxnID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Payment still pending",
		"status":  models.PaymentPending,
	})
}

func ReverifyFailedTransaction(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.ReverifyTransactionRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "REVERIFY") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "REVERIFY")

	q := db.New()

	email, err := q.GetEmailByTxnId(ctx, tx, req.TxnID)
	if err != nil {
		pkg.Log.ErrorCtx(c, "[REVERIFY-ERROR]: Failed to get email by txn id", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		return
	}

	booking, err := q.GetBookingByTxnID(ctx, tx, req.TxnID)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Transaction not found",
		})
		pkg.Log.ErrorCtx(c, "[REVERIFY-ERROR]: Requested transaction id is not found", err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[REVERIFY-ERROR]: Failed to fetch booking", err)
		return
	}

	// Only reverifying failed transactions
	if booking.TxnStatus == models.PaymentFailed {

		// Calling PayU verify API here.
		formBody := pkg.BuildVerifyPayUForm(req.TxnID)
		httpReq, err := http.NewRequest(
			"POST",
			cmd.Env.PayUVerifyURL,
			strings.NewReader(formBody),
		)
		if err != nil {
			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to build PayU request", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			return
		}
		httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to contact PayU", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Payment Gateway Failed. Try again later", // Different msg
			})
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to close body", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				return
			}
		}()

		// Parsing the response
		var payURes models.PayUVerifyResponse
		if err := json.NewDecoder(resp.Body).Decode(&payURes); err != nil {
			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to decode PayU verify response", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Payment Gateway Failed. Try again later",
			})
			return
		}

		gatewayStatus := models.MapPayUStatus(payURes, req.TxnID)

		event, err := q.GetEventByIdQuery(ctx, tx, booking.EventID)
		if err != nil {
			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to get event details", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			return
		}

		// if still failed
		if gatewayStatus == models.PaymentFailed || gatewayStatus == models.PaymentNotFound {
			// Restore released seats and close dispute
			disputeId, err := q.GetDisputeIdFromTxnIdQuery(ctx, tx, req.TxnID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to get dispute id from txn id", err)
				return
			}
			// disputeId, ok := pkg.GrabUuid(c, disputeIdStr, "REVERIFY", "disputeId")
			// if !ok {
			// 	return
			// }

			row, err := q.CloseAsFalseDisputeQuery(ctx, tx, disputeId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to close dispute as false", err)
				return
			}
			if row == 0 {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: No rows affected while closing dispute as false", nil)
				return
			}

			row, err = q.DecrementSeatFilledCountQuery(ctx, tx, booking.EventID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to decrement seat filled count", err)
				return
			}
			if row == 0 {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: No rows affected while decrementing seat filled count", nil)
				return
			}

			// booking status remains failed
			err = tx.Commit(ctx)
			if pkg.HandleDbTxnCommitErr(c, err, "REVERIFY") {
				return
			}

			pkg.Log.SuccessCtx(c)
			c.JSON(http.StatusOK, gin.H{
				"message": "Payment failed",
				"status":  models.PaymentFailed,
			})
			return
		}

		// if payment is success
		if gatewayStatus == models.PaymentSuccess {
			// Closing the dispute as valid payment
			disputeId, err := q.GetDisputeIdFromTxnIdQuery(ctx, tx, req.TxnID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to get dispute id from txn id", err)
				return
			}

			row, err := q.CloseAsTrueDisputeQuery(ctx, tx, disputeId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to close dispute as true", err)
				return
			}
			if row == 0 {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: No rows affected while closing dispute as true", nil)
				return
			}

			// Getting the schedule ids of the selected event
			schedules, err := q.GetSchedulesByEventID(ctx, tx, event.ID)
			if err != nil {
				pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to get schedules in success", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				return
			}

			// Getting details of the student - needed for mail and attendance
			student, err := q.GetStudentByEmail(ctx, tx, email)
			if err != nil {
				pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to retrive student data", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
			}

			if !event.IsGroup {
				for _, s := range schedules {
					_, err := q.CreateSoloEventParticipant(ctx, tx, db.CreateSoloEventParticipantParams{
						StudentID:       booking.StudentID,
						EventID:         booking.EventID,
						EventScheduleID: s,
						BookingID:       booking.ID,
						StudentName:     student.Name,
						StudentEmail:    student.Email,
					})
					if err != nil {
						pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to insert in solo participant", err)
						c.JSON(http.StatusInternalServerError, gin.H{
							"message": "Oops! Something happened. Please try again later",
						})
						return
					}
				}
			} else {
				teamId, err := q.GetTeamIDByBooking(ctx, tx, booking.ID)
				if err != nil {
					pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to get team details in verify", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "Oops! Something happened. Please try again later",
					})
					return
				}
				members, err := q.GetTeamMembersByTeamID(ctx, tx, teamId)
				if err != nil {
					pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to fetch team members.", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "Oops! Something happened. Please try again later",
					})
					return
				}
				// Inserting into team attendance table
				for _, s := range schedules {
					for _, m := range members {
						_, err := q.CreateTeamAttendance(ctx, tx, db.CreateTeamAttendanceParams{
							StudentID:       m.StudentID,
							EventScheduleID: s,
						})
						if err != nil {
							pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to insert into team attd.", err)
							c.JSON(http.StatusInternalServerError, gin.H{
								"message": "Oops! Something happened. Please try again later",
							})
							return
						}
					}
				}

			}

			err = q.UpdateBookingStatus(ctx, tx, db.UpdateBookingStatusParams{
				TxnStatus: models.PaymentSuccess,
				ID:        booking.ID,
			})
			if err != nil {
				pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to update booking status", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				return
			}

			// success case
			err = tx.Commit(ctx)
			if pkg.HandleDbTxnCommitErr(c, err, "VERIFY") {
				return
			}

			// If there is metadata, read and publish
			if len(booking.Metadata) > 0 {
				var metadataMap map[string]any
				if err := json.Unmarshal(booking.Metadata, &metadataMap); err != nil {
					pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to unmarshal booking metadata", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "Oops! Something happened. Please try again later",
					})
					return
				}

				// Hackathon payload
				if raw, ok := metadataMap["hackathon_payload"]; ok {
					payloadStr, ok := raw.(string)
					if !ok {
						pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: hackathon_payload is not string", nil)
						return
					}

					if err := messagequeue.Rabbit.Publish(
						ctx,
						messagequeue.QueueHackathonRegistrations,
						[]byte(payloadStr),
					); err != nil {
						pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to publish hackathon payload", err)
						return
					}
					// WOC payload
				} else if raw, ok := metadataMap["woc_payload"]; ok {
					payloadStr, ok := raw.(string)
					if !ok {
						pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: woc_payload is not string", nil)
						return
					}

					if err := messagequeue.Rabbit.Publish(
						ctx,
						messagequeue.QueueWocRegistrations,
						[]byte(payloadStr),
					); err != nil {
						pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to publish WOC payload", err)
						return
					}
				}
			}

			var completeSchedules []models.EventScheduleInput
			schedulesBytes, _ := json.Marshal(event.Schedules)
			_ = json.Unmarshal(schedulesBytes, &completeSchedules)

			var selected models.EventScheduleInput
			if len(completeSchedules) > 0 {
				selected = completeSchedules[0]
				for _, s := range completeSchedules {
					d1, _ := time.Parse("2006-01-02", s.EventDate)
					d2, _ := time.Parse("2006-01-02", selected.EventDate)
					if d1.Before(d2) {
						selected = s
					}
				}
			}

			err = mail.Mail.Enqueue(&mail.EmailRequest{
				To:      []string{},
				Subject: "Event Registration - Anokha 2026",
				Type:    "event-reg",
				Data: &mail.RegistrationData{
					UserName:      student.Name,
					EventName:     event.EventName,
					EventDate:     selected.EventDate,
					EventTime:     selected.StartTime + " - " + selected.EndTime,
					EventLocation: selected.Venue,
				},
				Retries: mail.MaxRetryCount,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				pkg.Log.ErrorCtx(c, "[MAIL-ERROR]: Failed to add request to email queue", err)
				return
			}
			pkg.Log.SuccessCtx(c)
			c.JSON(http.StatusOK, gin.H{
				"message": "Payment verified successfully",
				"status":  models.PaymentSuccess,
			})
			return
		}
		pkg.Log.WarnCtx(c, "[VERIFY-WARN]: Verification successful, but status still pending for "+req.TxnID)
		c.JSON(http.StatusOK, gin.H{
			"message": "Payment still pending",
			"status":  models.PaymentPending,
		})

	}
}
