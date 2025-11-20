package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func StudentAuthRoutes(r *gin.RouterGroup) {
	r.GET("/user/check", CheckEmailExist)

	// CSRF requests
	r.GET("/user/login", LoginUserCsrf)
	r.GET("/user/register", RegisterUserAccountCsrf)
	r.GET("/user/register/otp/verify")

	// Actual requests
	r.POST("/user/login", LoginUser)
	r.POST("/user/register", RegisterUserAccount)
	r.POST("/user/register/otp/verify", mw.VerifyCsrf, mw.TempTokenAuth, VerifyUserOtp)
	r.GET("/user/register/otp/resend", ResendUserOtp)

	r.GET("/user/session", mw.Auth, FetchUserSession)
	r.GET("/user/logout", mw.Auth, LogoutUser)
}

func OrganizerAuthRoutes(r *gin.RouterGroup) {
	r.GET("/organizer/login", LoginOrganizerCsrf)
	r.POST("/organizer/login", LoginOrganizer)
	r.GET("/organizer/logout", mw.Auth, LogoutOrganizer)
}
