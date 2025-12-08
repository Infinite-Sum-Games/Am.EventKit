package booking

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func BookingRoutes(r *gin.RouterGroup) {
	// CSRF Endpoints
	r.GET("/:eventId/book", mw.Auth)

	// Booking endpoints
	r.POST("/:eventId/book", mw.VerifyCsrf, mw.Auth, BookEvent)
}
