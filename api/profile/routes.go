package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func ProfileRoutes(r *gin.RouterGroup) {
	r.GET("/profile", mw.Auth, FetchUserProfile)
	r.GET("/profile/edit", mw.Auth, EditUserProfileCsrf)
	r.POST("/profile/edit", mw.Auth, EditUserProfile)
}
