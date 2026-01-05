package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func AccomodationFormRoutes(r *gin.RouterGroup) {
	r.GET("/check", mw.Auth, mw.CheckUser, AccomodationExists)
	r.GET("/", mw.Auth, mw.CheckUser, AccomodationFormCsrf)
	r.POST("/", mw.Auth, mw.VerifyCsrf, mw.CheckUser, AccomodationFormSubmission)
}

func AccomodationAuthRoutes(r *gin.RouterGroup) {
	r.POST("/panel/login", AccomodationLogin)
	r.POST("/panel/logout", mw.Auth, mw.CheckHospitality, AccomodationLogout)
	r.GET("/panel/session", mw.Auth, mw.CheckHospitality, AccomodationSession)
}

func AccomodationPanelRoutes(r *gin.RouterGroup) {
	r.GET("/panel", mw.Auth, mw.CheckHospitality, GetAllAccomodationRequests)
	r.GET("/panel/hostel", mw.Auth, mw.CheckHospitality, GetAllHostelDetails)
	r.POST("/panel/hostel", mw.Auth, mw.CheckHospitality, AddHostel)
	r.PUT("/panel/hostel", mw.Auth, mw.CheckHospitality, UpdateHostel)
	r.DELETE("/panel/hostel/:id", mw.Auth, mw.CheckHospitality, DeleteHostel)
	r.POST("/panel/allot/:accommodationId", mw.Auth, mw.CheckHospitality, AllotHostel)
}

func FinanceRoutes(r *gin.RouterGroup) {
	r.GET("/app/pay/:hospitalityId", mw.Auth, mw.CheckHospitality, GetFinanceDetailsByHospitalityId)
	r.POST("/app/pay/confirm/:accommodationId", mw.Auth, mw.CheckHospitality, mw.CheckFinance, AffirmAccommodationPayment)
}

func GateRoutes(r *gin.RouterGroup) {
	r.GET("/app/hostels", mw.Auth, mw.CheckHospitality, mw.CheckGate, GetAllHostels)
	r.POST("/app/map", mw.Auth, mw.CheckHospitality, mw.CheckGate, MapQrStudentId)
	r.GET("/app/:accommodationId", mw.Auth, mw.CheckHospitality, mw.CheckGate, GetAccommodationById)
	r.PUT("/app/:accommodationId", mw.Auth, mw.CheckHospitality, mw.CheckGate, UpdateAccommodationById)

	r.POST("/app/gate/check-in/:accomodationId")
	r.POST("/app/gate/status/:accomodationId")
	r.POST("/app/gate/check-out/:hospitalityId")
}

func SecurityRoutes(r *gin.RouterGroup) {
	r.GET("/app/security/:hospitalityId", mw.Auth, mw.CheckHospitality, mw.CheckSecurity, GetStudentDetailsForSecurity)
}
