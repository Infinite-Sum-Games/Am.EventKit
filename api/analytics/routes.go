package api

import (
	"github.com/gin-gonic/gin"
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
)

func AnalyticsRoutes(r *gin.RouterGroup) {
	r.GET("/revenue", mw.Auth, mw.CheckAdmin, GetRevenueAnalytics)
	r.GET("/registrations", mw.Auth, mw.CheckAdmin, GetEventRegistrationAnalytics)
	r.GET("/people", mw.Auth, mw.CheckAdmin, GetPeopleAnalytics)
	r.GET("/transactions", mw.Auth, mw.CheckAdmin, GetTransactionAnalytics)
}
