package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func BookingRoutes(r *gin.RouterGroup) {
	// CSRF Endpoints
	r.GET("/:eventId/book", mw.Auth, BookEventCsrf)

	// Booking endpoints
	r.POST("/:eventId/book", mw.VerifyCsrf, mw.Auth, BookEvent)
	r.POST("/verify", mw.Auth, VerifyTransaction)

	// Fetch transactions
	r.GET("/transactions", mw.Auth, mw.CheckAdmin, FetchAdminTransactions)
}
