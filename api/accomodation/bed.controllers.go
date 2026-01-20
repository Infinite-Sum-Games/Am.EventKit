package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Infinite-Sum-Games/Am.EventKit/cmd"
	db "github.com/Infinite-Sum-Games/Am.EventKit/db/gen"
	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
	"github.com/gin-gonic/gin"
)

func FetchUnclaimedBeds(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "BED") {
		return
	}
	defer conn.Release()

	q := db.New()
	beds, err := q.FetchUnclaimedBedsQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[BED-ERROR]: Failed to fetch unclaimed beds", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fetched all unclaimed beds",
		"beds":    beds,
	})
	pkg.Log.SuccessCtx(c)
}

func DeleteUnclaimedBed(c *gin.Context) {
	accId, ok := pkg.GrabUuid(c, c.Param("accId"), "BED", "Accommodation")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbAcquireErr(c, err, "BED") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "BED")

	q := db.New()
	rowsAffected, err := q.DeleteUnclaimedBedQuery(ctx, tx, accId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[BED-ERROR]: Failed to unclaim bed", err)
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Bed not found",
		})
		pkg.Log.WarnCtx(c, "[BED-WARN]: No unclaimed bed with given accommodation ID")
		return
	}

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "BED") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Deleted unclaimed bed successfully",
	})
	pkg.Log.SuccessCtx(c)
}
