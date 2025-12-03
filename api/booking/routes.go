package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func BookingRoutes(r *gin.RouterGroup) {
	// CSRF Endpoints
	// TODO: Should write these endpoints
	r.GET("/:eventId/book", mw.Auth)
	r.GET("/verify", mw.Auth)

	// Booking endpoints
	r.POST("/:eventId/book", mw.VerifyCsrf, mw.Auth, BookEvent)
	r.POST("/verify", mw.VerifyCsrf, mw.Auth, VerifyTransaction)
}
