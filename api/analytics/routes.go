package api

import "github.com/gin-gonic/gin"

func AnalyticsRoutes(r *gin.RouterGroup) {
	r.GET("/revenue", GetRevenueAnalytics)
	r.GET("/participants", GetParticipantAnalytics)
	r.GET("/registrations", GetRegistrationAnalytics)
	r.GET("/people", GetPeopleAnalytics)
}
