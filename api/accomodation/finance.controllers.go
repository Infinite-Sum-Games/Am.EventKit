package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func GetFinanceDetailsByHospitalityId(c *gin.Context) {
	hospitalityId := c.Param("hospitalityId")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "GET-FINANCE-DETAILS") {
		return
	}
	defer conn.Release()

	q := db.New()

	financeDetails, err := q.GetFinanceDetailsByHospitalityIdQuery(ctx, conn, pkg.ToPgText(hospitalityId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[GET-FINANCE-DETAILS-ERROR]: Failed to fetch finance details", err)
		return
	}

	if financeDetails.PaymentStatus == "COMPLETED" {
		c.JSON(http.StatusOK, gin.H{
			"message":         "Ammount already paid or is a hosteller",
			"accommodationId": financeDetails.AccommodationID,
			"name":            financeDetails.Name,
			"email":           financeDetails.Email,
			"num_days":        financeDetails.DayCount,
			"amount":          0,
			"paymentStatus":   true,
		})
		pkg.Log.SuccessCtx(c)
		return
	}

	var amount int32
	var paymentStatus = false

	if !financeDetails.IsAmritaCampus {
		amount = financeDetails.OutsiderPrice
	} else {
		amount = financeDetails.DayScholarPrice
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Finance details fetched successfully",
		"accommodationId": financeDetails.AccommodationID,
		"name":            financeDetails.Name,
		"email":           financeDetails.Email,
		"num_days":        financeDetails.DayCount,
		"amount":          amount * financeDetails.DayCount,
		"paymentStatus":   paymentStatus,
	})
	pkg.Log.SuccessCtx(c)

}
