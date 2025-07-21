package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func ProfileRoutes(r *gin.RouterGroup) {
	r.GET("/profile/:id", FetchUserProfile)
	r.GET("/profile/edit", mw.Auth, EditUserProfileCsrf)
	r.POST("/profile/edit", mw.Auth, mw.VerifyCsrf, EditUserProfile)
}
