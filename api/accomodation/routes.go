package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func AccomodationRoutes(r *gin.RouterGroup) {
	// Admin Panel
	// r.GET("/panel", GetAllAccomodationRequests)

	// Website
	r.GET("/check", mw.Auth, mw.CheckUser, AccomodationExists)
	r.GET("/", mw.Auth, mw.CheckUser, AccomodationFormCsrf)
	r.POST("/", mw.Auth, mw.VerifyCsrf, mw.CheckUser, AccomodationFormSubmission)
}
