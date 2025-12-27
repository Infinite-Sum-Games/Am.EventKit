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

func GetDisputeByID(c *gin.Context) {
	disputeIdStr := c.Param("dispute_id")

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

	dispute, err := q.GetDisputeByIDQuery(ctx, conn, disputeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to fetch dispute by ID", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messaage": "Dispute fetched successfully",
		"dispute":  dispute,
	})
	pkg.Log.SuccessCtx(c)
}

func CreateDispute(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.CreateDisputeInput](c)
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

	err = q.CreateDisputeQuery(ctx, tx, db.CreateDisputeQueryParams{
		EventID: req.EventId,
		TxnID:   req.TransactionId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to create dispute", err)
		return
	}

	rows, err := q.IncrementSeatFilledCountQuery(ctx, tx, req.EventId)
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

	dispute, err := q.GetDisputeByIDQuery(ctx, conn, disputeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to fetch dispute by ID", err)
		return
	}

	event, err := q.GetEventByIdQuery(ctx, conn, dispute.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to fetch event by ID", err)
		return
	}

	if event.IsGroup {
		teamMemberBytes, err := pkg.MarshalTeamMemberDetails(req.TeamMemberDetails)
		if err != nil {
			return
		}

		row, err := q.UpdateDisputeGroupQuery(ctx, conn, db.UpdateDisputeGroupQueryParams{
			ID:                disputeId,
			StudentEmail:      pkg.ToPgText(req.StudentEmail),
			Description:       pkg.ToPgText(req.Description),
			TeamMemberDatails: teamMemberBytes,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: Failed to update group dispute", err)
			return
		}
		if row == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later",
			})
			pkg.Log.ErrorCtx(c, "[DISPUTE-ERROR]: No rows affected while updating group dispute", nil)
			return
		}
	} else {
		row, err := q.UpdateDisputeSoloQuery(ctx, conn, db.UpdateDisputeSoloQueryParams{
			ID:           disputeId,
			StudentEmail: pkg.ToPgText(req.StudentEmail),
			Description:  pkg.ToPgText(req.Description),
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

	row, err := q.DecrementSeatFilledCountQuery(ctx, tx, disputeId)
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
