package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/mail"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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
	var req models.StudentOnboardingRequest
	if err := c.BindJSON(&req); err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to bind JSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	if err := req.Validate(); err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Validation failed", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	otp, err := pkg.GenerateOTP()
	if err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Unable to generate OTP", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		return
	}
	otpStr := strconv.Itoa(otp)
	otpSlice := make([]string, len(otpStr))
	for i, ch := range otpStr {
		otpSlice[i] = string(ch)
	}

	var expiry pgtype.Timestamp
	expiry.Time = time.Now().Add(10 * time.Minute)
	expiry.Valid = true

	//TODO: hashing the password or will the password sent from frontend be hashes?
	hashedPassword, err := pkg.Hash(req.Password)
	if err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Unable to hash password", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		return
	}

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to acquire DB connection"})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()
	err = q.UpsertStudentOnboarding(ctx, conn, db.UpsertStudentOnboardingParams{
		Name:            req.Name,
		DepartmentName:  req.DepartmentName,
		Email:           req.Email,
		Password:        hashedPassword,
		PhoneNumber:     req.PhoneNumber,
		IsAmritaStudent: req.IsAmritaStudent,
		AmritaRollNumber: pgtype.Text{
			String: req.AmritaRollNumber,
			Valid:  req.AmritaRollNumber != "",
		},
		CollegeName:  req.CollegeName,
		CollegeCity:  req.CollegeCity,
		AcademicYear: req.AcademicYear,
		Otp:          otpStr,
		ExpiryAt:     expiry,
	})
	if err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to insert onboarding data", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	// QUESTION: CreateToken functions are not returning any errors, is that fine?
	tempToken := pkg.CreateTempToken(req.Name, req.Email)
	// also there is no function to set temp token, so i wrote a new one
	pkg.SetTempCookie(c, tempToken)

	mail.Mail.Enqueue(mail.EmailRequest{
		To:      []string{req.Email},
		Subject: "Welcome to Anokha 2025",
		Type:    "otp",
		// should pass this as pointer (IMPORTANT)
		Data: &mail.OTPTemplateData{
			UserName: req.Name,
			OTP:      otpSlice,
		},
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully!\nkindly check mail for OTP - Check SPAM too :)",
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
