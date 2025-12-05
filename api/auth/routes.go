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
	r.POST("/user/register/otp/verify", mw.TempAuth, mw.VerifyCsrf, VerifyUserOtp)
	r.GET("/user/register/otp/resend", mw.TempAuth, ResendUserOtp)

	r.POST("/user/forgot-password", mw.VerifyCsrf, ForgotUserPassword)
	r.POST("/user/forgot-password/otp/verify", mw.TempAuth, mw.VerifyCsrf, ConfirmPasswordChange)
	r.GET("/user/session", mw.Auth, mw.CheckUser, FetchUserSession)
	r.GET("/user/logout", mw.Auth, mw.CheckUser, Logout)
}

func OrganizerAuthRoutes(r *gin.RouterGroup) {
	r.GET("/organizer/login", LoginOrganizerCsrf)
	r.POST("/organizer/login", LoginOrganizer)
	r.GET("/organizer/logout", mw.Auth, mw.CheckOrganizer, Logout)
}

func AdminAuthRoutes(r *gin.RouterGroup) {
	r.POST("/admin/login", LoginAdmin)
	r.GET("/admin/logout", mw.Auth, mw.CheckAdmin, Logout)
	r.GET("/admin/session", mw.Auth, mw.CheckAdmin, FetchAdminSession)
}
