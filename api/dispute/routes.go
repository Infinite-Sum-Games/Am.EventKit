package api

import "github.com/gin-gonic/gin"

func DisputeRoutes(r *gin.RouterGroup) {
	r.GET("/", GetAllDisputes)
	r.POST("/:txnId", CreateDispute)
	r.PUT("/:disputeId", UpdateDispute)
	r.POST("/closeAsTrue/:disputeId", CloseAsTrueDispute)
	r.POST("/closeAsFalse/:disputeId", CloseAsFalseDispute)
}
