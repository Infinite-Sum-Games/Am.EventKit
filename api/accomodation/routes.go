package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func AccomodationRoutes(r *gin.RouterGroup) {
	// Accomodation Panel
	r.POST("/panel/login", AccomodationLogin)
	r.POST("/panel/logout", mw.Auth, mw.CheckHospitality, AccomodationLogout)
	r.GET("/panel/session", mw.Auth, mw.CheckHospitality, AccomodationSession)
	r.GET("/panel", mw.Auth, mw.CheckHospitality, GetAllAccomodationRequests)
	r.POST("panel/mapQrStudent", MapQrStudentId)
	r.GET("/panel/:accommodationId", mw.Auth, mw.CheckHospitality, mw.CheckGate, GetAccommodationById)
	r.PUT("/panel/:accommodationId", mw.Auth, mw.CheckHospitality, mw.CheckGate, UpdateAccommodationById)
	r.GET("panel/financeDetails/:hospitalityId", GetFinanceDetailsByHospitalityId)
	r.GET("panel/security/:hospitalityId", GetStudentDetailsForSecurity)
	r.GET("/panel/hostels", GetAllHostels)
	r.POST("/panel/hostel", AddHostel)
	r.PUT("/panel/hostel", UpdateHostel)
	r.DELETE("/panel/hostel/:id", DeleteHostel)
	r.POST("/panel/allot/:accommodationId", AllotHostel)
	r.POST("/panel/affirmPayment/:accommodationId", AffirmAccommodationPayment)

	// Website
	r.GET("/check", mw.Auth, mw.CheckUser, AccomodationExists)
	r.GET("/", mw.Auth, mw.CheckUser, AccomodationFormCsrf)
	r.POST("/", mw.Auth, mw.VerifyCsrf, mw.CheckUser, AccomodationFormSubmission)
}
