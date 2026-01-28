package router

import (
	accommodationApi "github.com/Infinite-Sum-Games/Am.EventKit/api/accomodation"
	mw "github.com/Infinite-Sum-Games/Am.EventKit/middleware"
	"github.com/gin-gonic/gin"
)

func HospitalityAppRouter(r *gin.RouterGroup) {
	// Authentication routes
	auth := r.Group("/hospitality")
	{
		auth.POST("/login", accommodationApi.AccomodationLogin)
		auth.POST("/logout", accommodationApi.AccomodationLogout)
		auth.GET("/session", accommodationApi.AccomodationSession)
	}

	// Hostel management routes
	hostel := r.Group("/hostel", mw.SplRole("hostel"))
	{
		hostel.GET("", accommodationApi.GetAllHostels)
		hostel.POST("/:hospId/checkin", accommodationApi.HostelCheckIn)
		hostel.GET("/logs", accommodationApi.HostelLogsSink)
	}

	// Gate management routes
	gate := r.Group("/gate", mw.SplRole("gate"))
	{
		gate.POST("/:hospId/checkin", accommodationApi.GateCheckIn)
		gate.POST("/:hospId/checkout", accommodationApi.GateCheckOut)
		gate.GET("/status/:hospId", accommodationApi.GateStatus)
		gate.GET("/checkin-status/:hospId", accommodationApi.GateCheckInStatus)
		gate.GET("/checkout-status/:hospId", accommodationApi.GateCheckOutStatus)
		gate.GET("/logs", accommodationApi.GateLogsSink)
	}

	// Security routes
	security := r.Group("/security", mw.SplRole("security"))
	{
		security.GET("/check/:hospId", accommodationApi.SecurityCheck)
	}

	// Finance routes
	finance := r.Group("/finance", mw.SplRole("finance"))
	{
		finance.GET("/:hospId", accommodationApi.GetFinanceDetailsByHospitalityId)
	}

	// Bed management routes
	beds := r.Group("/beds", mw.SplRole("hostel"))
	{
		beds.GET("/unclaimed", accommodationApi.FetchUnclaimedBeds)
		beds.DELETE("/unclaimed/:accId", accommodationApi.DeleteUnclaimedBed)
	}

	// Individual accommodation management routes
	accommodation := r.Group("/accommodation", mw.SplRole("hostel"))
	{
		accommodation.GET("/:accId", accommodationApi.GetAccommodationById)
		accommodation.PUT("/:accId", accommodationApi.UpdateAccommodationById)
		accommodation.DELETE("/:accId", accommodationApi.DeleteAccommodationById)
	}

	// QR management
	r.POST("/qr/map", mw.SplRole("gate"), accommodationApi.MapQrStudentId)
}

func HospitalityWebRouter(r *gin.RouterGroup) {
	// Panel management routes (Administrative web interface)
	panel := r.Group("/panel", mw.Auth, mw.RequireRoles("hospitality"))
	{
		panel.GET("/requests", accommodationApi.GetAllAccomodationRequests)
		panel.POST("/hostels", accommodationApi.AddHostel)
		panel.PUT("/hostels/:id", accommodationApi.UpdateHostel)
		panel.DELETE("/hostels/:id", accommodationApi.DeleteHostel)
		panel.POST("/allot/:accId", accommodationApi.AllotHostel)
		panel.POST("/payment/:accId", accommodationApi.AffirmAccommodationPayment)
		panel.GET("/hostels/details", accommodationApi.GetAllHostelDetails)
	}
}
