package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func StudentAuthRoutes(r *gin.RouterGroup) {
	r.POST("/user/check", CheckEmailExist)

	// CSRF requests
	r.GET("/user/login", LoginUserCsrf)
	r.GET("/user/register", RegisterUserAccountCsrf)
	r.GET("/user/register/otp/verify", VerifyUserOtpCsrf)
	r.GET("/user/forgot-password", ForgotUserPasswordCsrf)
	r.GET("/user/forgot-password/otp/verify", ConfirmPasswordChangeCsrf)

	// Actual requests
	r.POST("/user/login", mw.VerifyCsrf, LoginUser)
	r.POST("/user/register", mw.VerifyCsrf, RegisterUserAccount)
	r.POST("/user/register/otp/verify", mw.VerifyCsrf, mw.TempAuth, VerifyUserOtp)
	r.GET("/user/register/otp/resend", mw.TempAuth, ResendUserOtp)

	r.POST("/user/forgot-password", mw.VerifyCsrf, ForgotUserPassword)
	r.POST("/user/forgot-password/otp/verify", mw.VerifyCsrf, mw.TempAuth, ConfirmPasswordChange)
	r.GET("/user/session", mw.Auth, FetchUserSession)
	r.GET("/user/logout", mw.Auth, Logout)
}

func OrganizerAuthRoutes(r *gin.RouterGroup) {
	r.GET("/organizer/login", LoginOrganizerCsrf)
	r.POST("/organizer/login", LoginOrganizer)
	r.GET("/organizer/logout", mw.Auth, Logout)
}

func AdminAuthRoutes(r *gin.RouterGroup) {
	r.POST("/admin/login", LoginAdmin)
}
