package api

import (
	mw "github.com/Infinite-Sum-Games/Am.EventKit/middleware"
	"github.com/gin-gonic/gin"
)

func ProfileRoutes(r *gin.RouterGroup) {
	r.GET("/profile", mw.Auth, mw.CheckUser, FetchUserProfile)

	r.GET("/profile/edit", mw.Auth, mw.CheckUser, EditUserProfileCsrf)
	r.POST("/profile/edit", mw.Auth, mw.VerifyCsrf, mw.CheckUser, EditUserProfile)

	r.GET("/profile/transactions", mw.Auth, mw.CheckUser, GetAllUserTransactions)
	r.GET("/profile/tickets", mw.Auth, mw.CheckUser, GetTickets)
}
