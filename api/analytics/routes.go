package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func AnalyticsRoutes(r *gin.RouterGroup) {
	r.GET("/quick", mw.Auth, mw.CheckAdmin, GetQuickDashboard)
	r.GET("/revenue", mw.Auth, mw.CheckAdmin, GetRevenueAnalytics)
	r.GET("/registrations", mw.Auth, mw.CheckAdmin, GetEventRegistrationAnalytics)
	// r.GET("/people", mw.Auth, mw.CheckAdmin, GetPeopleAnalytics)
	r.GET("/transactions", mw.Auth, mw.CheckAdmin, GetTransactionAnalytics)

	// Hospitality Analytics
	r.GET("/hospitality/inside", GetInsideCampusAnalytics)
}
