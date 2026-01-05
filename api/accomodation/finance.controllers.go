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

	roomType := pkg.ExtractRoomType(financeDetails.HostelName)

	if roomType == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[GET-FINANCE-DETAILS-ERROR]: Failed to extract room type from hostel name", err)
		return
	}

	var amount int32
	var paymentStatus bool

	if financeDetails.PaymentStatus == "COMPLETED" {
		paymentStatus = true
	} else {
		paymentStatus = false
	}

	if !financeDetails.IsAmritaCampus {
		if roomType == "SINGLE" {
			amount = 340
		} else {
			amount = 230
		}
	} else {
		if financeDetails.IsHosteller {
			paymentStatus = true
			amount = 0
		} else {
			if roomType == "SINGLE" {
				amount = 300
			} else {
				amount = 200
			}
		}
	}

	if paymentStatus {
		c.JSON(http.StatusOK, gin.H{
			"message": "Ammount already paid or is a hosteller",
			"amount":  0,
		})
		pkg.Log.SuccessCtx(c)
		return
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
