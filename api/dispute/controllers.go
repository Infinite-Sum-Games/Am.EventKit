package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Infinite-Sum-Games/Am.EventKit/cmd"
	db "github.com/Infinite-Sum-Games/Am.EventKit/db/gen"
	"github.com/Infinite-Sum-Games/Am.EventKit/models"
	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
	"github.com/gin-gonic/gin"
)

func GetAllDisputes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "DISPUTE") {
		return
	}
	defer conn.Release()

	q := db.New()

	disputes, err := q.GetAllDisputesQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to fetch disputes", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messaage": "Disputes fetched successfully",
		"disputes": disputes,
	})
	pkg.Log.SuccessCtx(c)
}

func CreateDispute(c *gin.Context) {
	txnId := c.Param("txnId")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "DISPUTE") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "DISPUTE")

	q := db.New()

	event, err := q.GetEventIdByTxnIdQuery(ctx, tx, txnId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to fetch event ID by transaction ID", err)
		return
	}

	disputeCount, err := q.CheckDisputeExistsByTxnIdQuery(ctx, tx, txnId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to check if dispute exists by transaction ID", err)
		return
	}
	if disputeCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Dispute already exists for this transaction",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Attempted to create duplicate dispute for transaction ID", nil)
		return
	}

	studentEmail, err := q.GetEmailByTxnIdQuery(ctx, tx, txnId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to fetch student email by transaction ID", err)
		return
	}

	if event.EventStatus == "ACTIVE" {

		err = q.CreateDisputeQuery(ctx, tx, db.CreateDisputeQueryParams{
			EventID:      event.EventID,
			TxnID:        txnId,
			StudentEmail: pkg.ToPgText(studentEmail),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to create dispute", err)
			return
		}

		rows, err := q.IncrementSeatFilledCountQuery(ctx, tx, event.EventID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to increment seat filled count", err)
			return
		}
		if rows == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: No rows affected while incrementing seat filled count", nil)
			return
		}

		err = tx.Commit(ctx)
		if pkg.HandleDbTxnCommitErr(c, err, "DISPUTE") {
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Dispute created successfully",
		})
		pkg.Log.SuccessCtx(c)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot create dispute for inactive event",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Attempted to create dispute for inactive event", nil)
	}

}

func UpdateDispute(c *gin.Context) {
	disputeIdStr := c.Param("disputeId")

	disputeId, ok := pkg.GrabUuid(c, disputeIdStr, "DISPUTE", "disputeId")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.UpdateDisputeStatusInput](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "DISPUTE") {
		return
	}
	defer conn.Release()

	q := db.New()

	row, err := q.UpdateDisputeQuery(ctx, conn, db.UpdateDisputeQueryParams{
		ID:          disputeId,
		Description: pkg.ToPgText(req.Description),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to update individual dispute", err)
		return
	}
	if row == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: No rows affected while updating individual dispute", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dispute updated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func CloseAsTrueDispute(c *gin.Context) {
	disputeIdStr := c.Param("disputeId")

	disputeId, ok := pkg.GrabUuid(c, disputeIdStr, "DISPUTE", "disputeId")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "DISPUTE") {
		return
	}
	defer conn.Release()

	q := db.New()

	rows, err := q.CloseAsTrueDisputeQuery(ctx, conn, disputeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to close dispute as true", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: No rows affected while closing dispute as true", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dispute closed as true successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func CloseAsFalseDispute(c *gin.Context) {
	disputeIdStr := c.Param("disputeId")

	disputeId, ok := pkg.GrabUuid(c, disputeIdStr, "DISPUTE", "disputeId")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "DISPUTE") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "DISPUTE")

	q := db.New()

	eventId, err := q.GetEventIdByDisputeIDQuery(ctx, tx, disputeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to fetch event ID by dispute ID", err)
		return
	}

	rows, err := q.CloseAsFalseDisputeQuery(ctx, tx, disputeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to close dispute as false", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: No rows affected while closing dispute as false", nil)
		return
	}

	row, err := q.DecrementSeatFilledCountQuery(ctx, tx, eventId)
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

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "DISPUTE") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dispute closed as false successfully",
	})
	pkg.Log.SuccessCtx(c)

}
