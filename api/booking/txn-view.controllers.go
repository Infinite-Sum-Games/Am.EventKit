package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

// Controller Response Structure
// -----------------------------
// txn_id,
// event_id,
// event_name,
// student_name,
// email,
// phone_number,
// college_name,
// college_city,
// is_amrita_student,
// is_group
// txn_status
func FetchAdminTransactions(c *gin.Context) {
	txnType := c.Query("status")
	if txnType != models.PaymentFailed && txnType != models.PaymentPending && txnType != models.PaymentSuccess {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No transaction status requested",
		})
		pkg.Log.WarnCtx(c, "[BOOKING-VIEW-WARN]: Missing txn status in query parameter")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "BOOKING-VIEW") {
		return
	}
	defer conn.Release()

	q := db.New()
	results, err := q.FetchAdminTransactionsQuery(ctx, conn, txnType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[BOOKING-VIEW-ERROR]: Failed to fetch transactions", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Transaction list fetched successfully",
		"transactions": results,
	})
	pkg.Log.SuccessCtx(c)
}
