package api

import (
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func CheckEmailExist(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Email existence check initiated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func LoginUserCsrf(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Login action initiated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func RegisterUserAccountCsrf(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Register account action initiated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func VerifyUserOtpCsrf(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verification action initiated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func ResendUserOtpCsrf(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Resend OTP action initiated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func LoginUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func RegisterUserAccount(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func VerifyUserOtp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verification completed successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func ResendUserOtp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP resent to user email successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func LogoutUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged out successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func LoginStaffCsrf(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Staff login action initiated successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func LoginStaff(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Staff logged in successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func LogoutStaff(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Staff logged out successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func FetchStaffSession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Staff session obtained successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func FetchUserSession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "User session obtained successfully",
	})
	pkg.Log.SuccessCtx(c)
}
