package api

import (
	api "github.com/Thanus-Kumaar/anokha-2025-backend/api/util"
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func StudentAuthRoutes(r *gin.RouterGroup) {
	r.GET("/user/check", CheckEmailExist)

	// CSRF requests
	// r.GET("/user/login", LoginUserCsrf)
	// r.GET("/user/register", RegisterUserAccountCsrf)
	r.GET("/user/register/otp/verify", api.SendCsrfToken)
	r.GET("/user/register/otp/resend", api.SendCsrfToken)

	// Actual requests
	r.POST("/user/login", LoginUser)
	r.POST("/user/register", RegisterUserAccount)
	r.POST("/user/register/otp/verify", mw.VerifyCsrf, VerifyUserOtp)
	r.POST("/user/register/otp/resend", mw.VerifyCsrf, ResendUserOtp)

	r.GET("/user/session", mw.Auth, FetchUserSession)
	r.GET("/user/logout", mw.Auth, LogoutUser)
}

func StaffAuthRoutes(r *gin.RouterGroup) {
	r.GET("/staff/login", LoginStaffCsrf)
	r.POST("/staff/login", LoginStaff)

	r.GET("/staff/session", mw.Auth, FetchStaffSession)
	r.GET("/staff/logout", mw.Auth, LogoutStaff)
}
