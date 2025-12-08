package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func VerifyTransaction(c *gin.Context) {
	var req models.VerifyTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Request is malformed",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-ERROR]: Invalid request body for verify transaction", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != pgx.ErrTxClosed {
			pkg.Log.FatalCtx(c, "[VERIFY-FATAL]: Failed to rollback", rbErr)
		}
	}()

	q := db.New()

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

		// commiting transaction here
		if err := tx.Commit(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.FatalCtx(c, "[VERIFY-FATAL]: Failed to commit transaction", err)
			return
		}
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

		// If team event, then remove the team details
		event, err := q.GetEventByIdQuery(ctx, tx, booking.EventID)
		if err != nil {
			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR]: Failed to get event details", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			return
		}
		if event.IsGroup {
			teamID, err := q.GetTeamIDByBooking(ctx, tx, booking.ID)
			if err != nil && err != sql.ErrNoRows {
				pkg.Log.ErrorCtx(c, "[VERIFY-ERROR] Failed to get team ID", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Oops! Something happened. Please try again later",
				})
				return
			}
			if err == nil {
				// delete members first due to foreign key relation
				if err := q.DeleteTeamDetailsOfTeam(ctx, tx, teamID); err != nil {
					pkg.Log.ErrorCtx(c, "[VERIFY-ERROR] Failed to delete team members", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "Oops! Something happened. Please try again later",
					})
					return
				}
				// delete team
				if err := q.DeleteTeam(ctx, tx, booking.ID); err != nil {
					pkg.Log.ErrorCtx(c, "[VERIFY-ERROR] Failed to delete team", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "Oops! Something happened. Please try again later",
					})
					return
				}
			}
		}

		// Update booking status → failed
		err = q.UpdateBookingStatus(ctx, tx, db.UpdateBookingStatusParams{
			TxnStatus: models.PaymentFailed,
			ID:        booking.ID,
		})
		if err != nil {
			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR] Failed to update booking status", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			return
		}
		if err := tx.Commit(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.FatalCtx(c, "[VERIFY-FATAL]: Failed to commit transaction", err)
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

		err := q.UpdateBookingStatus(ctx, tx, db.UpdateBookingStatusParams{
			TxnStatus: models.PaymentSuccess,
			ID:        booking.ID,
		})
		if err != nil {
			pkg.Log.ErrorCtx(c, "[VERIFY-ERROR] Failed to update booking status", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			return
		}

		if err := tx.Commit(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.FatalCtx(c, "[VERIFY-FATAL]: Failed to commit transaction", err)
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
