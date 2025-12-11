package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func DiscountRoutes(r *gin.RouterGroup) {
	r.GET("/", mw.Auth, mw.CheckAdmin, GetAllDiscounts)
	// r.POST("/", mw.Auth, mw.CheckAdmin, CreateDiscount)
	// r.PUT("/:discountId", mw.Auth, mw.CheckAdmin, EditDiscount)
	// r.DELETE("/:organizerId", mw.Auth, mw.CheckAdmin, DeleteDiscount)
}
