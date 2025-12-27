package api

import "github.com/gin-gonic/gin"

func DisputeRoutes(r *gin.RouterGroup) {
	r.GET("/", GetAllDisputes)
	r.GET("/:disputeId", GetDisputeByID)
	r.POST("/", CreateDispute)
	r.PUT("/:disputeId", UpdateDispute)
	r.POST("/:disputeId/closeAsTrue", CloseAsTrueDispute)
	r.POST("/:disputeId/closeAsFalse", CloseAsFalseDispute)
}
