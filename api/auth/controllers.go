package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	api "github.com/Thanus-Kumaar/anokha-2025-backend/api/util"
	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/mail"
	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

func LoginUser(c *gin.Context) {
	// use the ValidateRequest function if the models.Type has a Validate() function defined
	req, ok := api.ValidateRequest[models.LoginRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()
	user, err := q.CheckStudentVerifiedQuery(ctx, conn, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found. Please register first."})
			return
		}
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to fetch user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		return
	}
	if req.HashedPassword != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password for given email"})
		return
	}

	var refreshToken string
	isNewToken := false
	if user.RefreshToken.String == "" {
		// no token exists, create new refresh token and set it in database
		refreshToken = pkg.CreateRefreshToken(user.ID.String(), user.Name, user.Email, true, false, false)
		isNewToken = true
	} else {
		// Validate existing token
		valid, _ := pkg.ParseToken(user.RefreshToken.String, "refresh_token")
		if !valid {
			refreshToken = pkg.CreateRefreshToken(user.ID.String(), user.Name, user.Email, true, false, false)
			isNewToken = true
		} else {
			refreshToken = user.RefreshToken.String
		}
	}
	if isNewToken {
		err = q.UpdateRefreshTokenQuery(ctx, conn, db.UpdateRefreshTokenQueryParams{
			ID: user.ID,
			RefreshToken: pgtype.Text{
				String: refreshToken,
				Valid:  refreshToken != "",
			},
		})
		if err != nil {
			pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to update refresh token", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
			return
		}
	}
	// createing auth token and setting both tokens as cookies
	authToken := pkg.CreateAuthToken(user.ID.String(), user.Name, user.Email, true, false, false)
	pkg.SetAuthCookie(c, authToken)
	pkg.SetRefreshCookie(c, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func RegisterUserAccount(c *gin.Context) {
	req, ok := api.ValidateRequest[models.StudentOnboardingRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// checking if user already registered successfully!
	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()
	_, err = q.CheckStudentVerifiedQuery(ctx, conn, req.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: DB error while checking student", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		return
	}
	if err == nil {
		// Student already exists
		c.JSON(http.StatusConflict, gin.H{"message": "Student is already registered"})
		return
	}

	otpStr, otpSlice, err := pkg.GenerateOTP()
	if err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Unable to generate OTP", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		return
	}

	var expiry pgtype.Timestamp
	expiry.Time = time.Now().Add(10 * time.Minute)
	expiry.Valid = true

	err = q.UpsertStudentOnboardingQuery(ctx, conn, db.UpsertStudentOnboardingQueryParams{
		Name:            req.Name,
		DepartmentName:  req.DepartmentName,
		Email:           req.Email,
		Password:        req.Password, // assuming password is hashed in frontend
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		return
	}
	// QUESTION: CreateToken functions are not returning any errors, is that fine?
	tempToken := pkg.CreateTempToken(req.Name, req.Email)
	// also there is no function to set temp token, so i wrote a new one
	pkg.SetTempCookie(c, tempToken)

	err = mail.Mail.Enqueue(&mail.EmailRequest{
		To:      []string{req.Email},
		Subject: "Welcome to Anokha 2025",
		Type:    "otp",
		// should pass this as pointer (IMPORTANT)
		Data: &mail.OTPTemplateData{
			UserName: req.Name,
			OTP:      otpSlice,
		},
	})
	if err != nil {
		pkg.Log.ErrorCtx(c, "[MAIL-ERROR]: Failed to add request to email queue", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully!\nkindly check mail for OTP - Check SPAM too :)",
	})
	pkg.Log.SuccessCtx(c)
}

func VerifyUserOtp(c *gin.Context) {
	var req struct {
		Otp string `json:"otp" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid OTP request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()
	row, err := q.GetStudentOtpQuery(ctx, conn, c.GetString("email"))
	if err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to get otp from table", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		return
	}
	if req.Otp != row.Otp {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid OTP"})
		return
	}
	if !row.ExpiryAt.Valid || row.ExpiryAt.Time.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "OTP has expired"})
		return
	}

	// Migration of data from onboarding table to original table
	tx, err := conn.Begin(ctx)
	if err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Unable to start transaction", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		return
	}
	defer tx.Rollback(ctx)
	userID, err := q.FinalizeStudentSignUpQuery(ctx, tx, c.GetString("email"))
	if err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Migration of student from onboarding failed!", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		return
	}
	if err = q.DeleteOnboardingQuery(ctx, tx, c.GetString("email")); err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to delete onboarding record", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		return
	}

	// token generation and setting cookie
	authToken := pkg.CreateAuthToken(userID.String(), c.GetString("username"), c.GetString("email"), true, false, false)
	refreshToken := pkg.CreateRefreshToken(userID.String(), c.GetString("username"), c.GetString("email"), true, false, false)
	pkg.SetAuthCookie(c, authToken)
	pkg.SetRefreshCookie(c, refreshToken)

	// Adding token to database
	if err = q.UpdateRefreshTokenQuery(ctx, tx, db.UpdateRefreshTokenQueryParams{
		RefreshToken: pgtype.Text{
			String: refreshToken,
			Valid:  refreshToken != "",
		},
		ID: userID,
	}); err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to add refresh token", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to commit transaction", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something happened. Please try again later"})
		return
	}
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
