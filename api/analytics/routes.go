package api

import (
	mw "github.com/Infinite-Sum-Games/Am.EventKit/middleware"
	"github.com/gin-gonic/gin"
)

func AnalyticsRoutes(r *gin.RouterGroup) {
	r.GET("/quick", mw.Auth, mw.CheckAdmin, GetQuickDashboard)
	r.GET("/revenue", mw.Auth, mw.CheckAdmin, GetRevenueAnalytics)
	r.GET("/registrations", mw.Auth, mw.CheckAdmin, GetEventRegistrationAnalytics)
	// r.GET("/people", mw.Auth, mw.CheckAdmin, GetPeopleAnalytics)
	r.GET("/transactions", mw.Auth, mw.CheckAdmin, GetTransactionAnalytics)

	// Hospitality Analytics
	r.GET("/hospitality/inside", mw.Auth, mw.CheckHospitality, GetInsideCampusAnalytics)
	r.GET("/hospitality/beds", mw.Auth, mw.CheckHospitality, GetLiveBedsAnalytics)
}
