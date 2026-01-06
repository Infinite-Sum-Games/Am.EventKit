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
	"github.com/jackc/pgx/v5"
)

func GetAllAccomodationRequests(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ACCOMODATION") {
		return
	}
	defer conn.Release()

	q := db.New()
	requests, err := q.GetAllAccommodationRequestsQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[ACCOMMODATION-ERROR]: Failed to fetch accomodation requests", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Accommodation requests returned successfully",
		"requests": requests,
	})
	pkg.Log.SuccessCtx(c)
}

func AddHostel(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.AddHostelRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "ADD-HOSTEL") {
		return
	}
	defer conn.Release()

	q := db.New()

	hostelId, err := q.AddHostelQuery(ctx, conn, db.AddHostelQueryParams{
		RoomCount:             req.RoomCount,
		IsMale:                req.IsMale,
		WardenEmail:           pkg.ToPgTextPtr(&req.WardenEmail),
		Latitude:              pkg.ToPgTextPtr(&req.Latitude),
		Longtitude:            pkg.ToPgTextPtr(&req.Longtitude),
		MapUrl:                pkg.ToPgTextPtr(&req.MapUrl),
		HostelName:            req.HostelName,
		AmritaDayscholarPrice: req.DayScholarPrice,
		NonAmritaPrice:        req.OutsiderPrice,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[ADD-HOSTEL-ERROR]: Failed to add new hostel", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Hostel added successfully",
		"hostel_id": hostelId,
	})
	pkg.Log.SuccessCtx(c)
}

func UpdateHostel(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.UpdateHostelRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "UPDATE-HOSTEL") {
		return
	}
	defer conn.Release()

	q := db.New()

	HostelUuid, ok := pkg.GrabUuid(c, req.HostelID, "UPDATE-HOSTEL", "hostelID")
	if !ok {
		return
	}

	rows, err := q.UpdateHostelQuery(ctx, conn, db.UpdateHostelQueryParams{
		ID:                    HostelUuid,
		RoomCount:             req.RoomCount,
		WardenEmail:           pkg.ToPgTextPtr(&req.WardenEmail),
		Latitude:              pkg.ToPgTextPtr(&req.Latitude),
		Longtitude:            pkg.ToPgTextPtr(&req.Longtitude),
		MapUrl:                pkg.ToPgTextPtr(&req.MapUrl),
		IsMale:                req.IsMale,
		AmritaDayscholarPrice: req.DayScholarPrice,
		NonAmritaPrice:        req.OutsiderPrice,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[UPDATE-HOSTEL-ERROR]: Failed to update hostel details", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Hostel not found",
		})
		pkg.Log.WarnCtx(c, "[UPDATE-HOSTEL-WARN]: Hostel ID does not exist")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hostel updated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func DeleteHostel(c *gin.Context) {
	hostelIdStr := c.Param("id")
	hostelId, ok := pkg.GrabUuid(c, hostelIdStr, "DELETE-HOSTEL", "hostelID")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "DELETE-HOSTEL") {
		return
	}
	defer conn.Release()

	q := db.New()

	rows, err := q.DeleteHostelQuery(ctx, conn, hostelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[DELETE-HOSTEL-ERROR]: Failed to delete hostel", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Hostel not found",
		})
		pkg.Log.WarnCtx(c, "[DELETE-HOSTEL-WARN]: Hostel ID does not exist")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hostel deleted successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func AllotHostel(c *gin.Context) {
	accIdStr := c.Param("accId")
	accommodationId, ok := pkg.GrabUuid(c, accIdStr, "ALLOT-HOSTEL", "Accommodation")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.AllotHostelRequest](c)
	if !ok {
		return
	}

	hostleIdPgUuid, err := pkg.ToPgUuidPtr(&req.HostelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid hostel ID format",
		})
		pkg.Log.ErrorCtx(c, "[ALLOT-HOSTEL-WARN]: Invalid hostel ID format", err)
		return
	}

	hostelIdUuid, ok := pkg.GrabUuid(c, req.HostelID, "ALLOT-HOSTEL", "Hostel")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "ALLOT-HOSTEL") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "ALLOT-HOSTEL")

	q := db.New()

	accommodation, err := q.GetAccommodationByIdQuery(ctx, tx, accommodationId)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Accommodation request not found",
		})
		pkg.Log.WarnCtx(c, "[ALLOT-HOSTEL-WARN]: Accommodation ID does not exist")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[ALLOT-HOSTEL-ERROR]: Failed to fetch accommodation request", err)
		return
	}

	hostel, err := q.GetHostelQuery(ctx, tx, hostelIdUuid)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Hostel not found",
		})
		pkg.Log.WarnCtx(c, "[ALLOT-HOSTEL-WARN]: Hostel ID does not exist")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[ALLOT-HOSTEL-ERROR]: Failed to fetch hostel details", err)
		return
	}

	if accommodation.IsMale != hostel.IsMale {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Hostel gender does not match student gender",
		})
		pkg.Log.WarnCtx(c,
			"[ALLOT-HOSTEL-WARN]: Hostel gender does not match student gender")
		return
	}

	if hostel.RoomCount <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No rooms available in the selected hostel",
		})
		pkg.Log.WarnCtx(c,
			"[ALLOT-HOSTEL-WARN]: No rooms available in the selected hostel")
		return
	}

	if accommodation.IsHosteller && accommodation.IsAmritaCampus {
		row, err := q.AffirmAccommodationAndPaymentQuery(ctx, tx,
			db.AffirmAccommodationAndPaymentQueryParams{
				ID:       accommodationId,
				ID_2:     hostelIdUuid,
				DayCount: req.DayCount,
			})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Oops! Something happened. Please try again later.",
			})
			pkg.Log.ErrorCtx(c,
				"[ALLOT-HOSTEL-ERROR]: Failed to affirm accommodation and payment for amrita hostellers", err)
			return
		}
		if row == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Accommodation request not found",
			})
			pkg.Log.ErrorCtx(c,
				"[ALLOT-HOSTEL-ERROR]: Accommodation ID does not exist", err)
			return
		}

		err = tx.Commit(ctx)
		if pkg.HandleDbTxnCommitErr(c, err, "ALLOT-HOSTEL") {
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Hostel allotted and payment confirmed for Amrita hosteller successfully",
		})
		pkg.Log.SuccessCtx(c)
		return
	}

	row, err := q.AllotHostelQuery(ctx, tx, db.AllotHostelQueryParams{
		ID:       accommodationId,
		HostelID: hostleIdPgUuid,
		DayCount: req.DayCount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[ALLOT-HOSTEL-ERROR]: Failed to allot hostel", err)
		return
	}
	if row == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Accommodation request or Hostel not found",
		})
		pkg.Log.ErrorCtx(c,
			"[ALLOT-HOSTEL-ERROR]: Accommodation ID or Hostel ID does not exist", err)
		return
	}

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "ALLOT-HOSTEL") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hostel allotted successfully",
	})
	pkg.Log.SuccessCtx(c)

}

func AffirmAccommodationPayment(c *gin.Context) {
	accommodationIdStr := c.Param("accId")
	accommodationId, ok := pkg.GrabUuid(c, accommodationIdStr, "AFFIRM-PAYMENT", "Accommodation")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "AFFIRM-PAYMENT") {
		return
	}
	defer conn.Release()

	q := db.New()

	rows, err := q.AffirmAccommodationPaymentQuery(ctx, conn, accommodationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AFFIRM-PAYMENT-ERROR]: Failed to affirm payment", err)
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Accommodation request not found",
		})
		pkg.Log.WarnCtx(c, "[AFFIRM-PAYMENT-WARN]: Accommodation ID does not exist")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Accommodation payment affirmed successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func GetAllHostelDetails(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "GET-HOSTELS") {
		return
	}
	defer conn.Release()

	q := db.New()

	hostels, err := q.GetAllHostelDetailsQuery(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[GET-HOSTELS-ERROR]: Failed to fetch hostel details", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hostel details fetched successfully",
		"hostels": hostels,
	})
	pkg.Log.SuccessCtx(c)
}
