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

func GetAllDiscounts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "DISCOUNT") {
		return
	}
	defer conn.Release()

	q := db.New()

	discounts, err := q.GetAllDiscounts(ctx, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISCOUNT-ERROR]: Failed to fetch discounts", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Discounts fetched successfully",
		"discounts": discounts,
	})
	pkg.Log.SuccessCtx(c)
}

func CreateDiscount(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, ok := pkg.ValidateRequest[models.CreateDiscountRequest](c)
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "DISCOUNT") {
		return
	}
	defer conn.Release()

	q := db.New()

	discount, err := q.CreateDiscount(ctx, conn, db.CreateDiscountParams{
		DiscountType:        req.DiscountType,
		DiscountedSoloSeats: req.DiscountedSoloSeats,
		DiscountedTeamSeats: req.DiscountedTeamSeats,
		StartTime:           req.StartTime,
		EndTime:             req.EndTime,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISCOUNT-ERROR]: Failed to create discount", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Discount created successfully",
		"discount": discount,
	})
	pkg.Log.SuccessCtx(c)
}

func EditDiscount(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	discountID, ok := pkg.GrabUuid(c, c.Param("discountId"), "DISCOUNT", "discount")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.CreateDiscountRequest](c)
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "DISCOUNT") {
		return
	}
	defer conn.Release()

	q := db.New()

	rowsAffected, err := q.EditDiscount(ctx, conn, db.EditDiscountParams{
		ID:                  discountID,
		DiscountType:        req.DiscountType,
		DiscountedSoloSeats: req.DiscountedSoloSeats,
		DiscountedTeamSeats: req.DiscountedTeamSeats,
		StartTime:           req.StartTime,
		EndTime:             req.EndTime,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISCOUNT-ERROR]: Failed to update discount", err)
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Discount does not exist",
		})
		pkg.Log.ErrorCtx(c, "[DISCOUNT-ERROR]: No rows updated", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Discount updated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func DeleteDiscount(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	discountID, ok := pkg.GrabUuid(c, c.Param("discountId"), "DISCOUNT", "discount")
	if !ok {
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "DISCOUNT") {
		return
	}
	defer conn.Release()

	q := db.New()

	rowsAffected, err := q.DeleteDiscount(ctx, conn, discountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[DISCOUNT-ERROR]: Failed to delete discount", err)
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Discount does not exist",
		})
		pkg.Log.ErrorCtx(c, "[DISCOUNT-ERROR]: No rows deleted", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Discount deleted successfully",
	})
	pkg.Log.SuccessCtx(c)
}
