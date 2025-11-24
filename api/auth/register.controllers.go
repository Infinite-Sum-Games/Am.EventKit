package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

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

	req, ok := pkg.ValidateRequest[models.CheckEmailRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})

		pkg.Log.FatalCtx(c, "[AUTH-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()
	_, err = q.FindEmailQuery(ctx, conn, req.Email)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{
			"message": "Email is available",
		})

		pkg.Log.SuccessCtx(c)
		return
	}

	c.JSON(http.StatusConflict, gin.H{
		"message": "Email already registered",
	})

	pkg.Log.WarnCtx(c, "[AUTH-WARN]: The email already exists")
}

func RegisterUserAccountCsrf(c *gin.Context) {
	csrfToken, tokenErr := pkg.CreateCsrfToken("register@amrita.edu", c)
	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to create CSRF token", tokenErr)
		return
	}

	pkg.SetCsrfCookie(c, csrfToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Register account action initiated successfully",
		"key":     csrfToken,
	})

	pkg.Log.SuccessCtx(c)
}

func RegisterUserAccount(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.StudentOnboardingRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})

		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer tx.Rollback(ctx)

	q := db.New()
	ok, err = q.FindEmailQuery(ctx, tx, req.Email)
	if err != nil && err != pgx.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: DB error while checking student", err)
		return
	}
	if ok {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Student is already registered",
		})
		pkg.Log.WarnCtx(c, "[AUTH-ERROR]: Re-attempt to register existing account")
		return
	}

	// If student has not registered
	otpStr, otpSlice, err := pkg.GenerateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})

		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Unable to generate OTP", err)
		return
	}

	expiry := pgtype.Timestamp{
		Time:  time.Now().Add(10 * time.Minute),
		Valid: true,
	}

	hashedPass, err := pkg.Hash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to hash password", err)
		return
	}

	err = q.UpsertStudentOnboardingQuery(ctx, tx, db.UpsertStudentOnboardingQueryParams{
		Name:            req.Name,
		DepartmentName:  req.DepartmentName,
		Email:           req.Email,
		Password:        hashedPass,
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

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})

		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to commit transaction", err)
		return
	}

	// Create a temporary token and dispatch it for the OTP request and forward
	// the OTP email. The sending of OTP email is not tied to the transaction
	// and is thus out of the transaction block. It's an async mailer process
	tempToken := pkg.CreateTempToken(req.Name, req.Email)
	pkg.SetTempCookie(c, tempToken)

	err = mail.Mail.Enqueue(&mail.EmailRequest{
		To:      []string{req.Email},
		Subject: "Welcome to Anokha 2025",
		Type:    "otp",
		Data: &mail.OTPTemplateData{
			UserName: req.Name,
			OTP:      otpSlice,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})

		pkg.Log.ErrorCtx(c, "[MAIL-ERROR]: Failed to add request to email queue", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User onboarding initiated. OTP sent to email.",
	})
	pkg.Log.SuccessCtx(c)
}

func VerifyUserOtpCsrf(c *gin.Context) {
	csrfToken, tokenErr := pkg.CreateCsrfToken("verify@otp", c)
	if tokenErr != nil {
		return
	}

	pkg.SetCsrfCookie(c, csrfToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Verify onboarding action initiated successfully",
		"key":     csrfToken,
	})

	pkg.Log.SuccessCtx(c)
}

func VerifyUserOtp(c *gin.Context) {

	email := c.GetString("email")
	if email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})

		pkg.Log.FatalCtx(
			c,
			"[AUTH-ERROR]: Could not find email in ctx",
			fmt.Errorf("BUG: Middleware did not add email in ctx"),
		)
		return
	}

	req, ok := pkg.ValidateRequest[models.OtpRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})

		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to begin transaction", err)
		return
	}
	defer tx.Rollback(ctx)

	q := db.New()

	row, err := q.GetStudentOtpQuery(ctx, tx, email)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "OTP is either invalid or expired.",
		})

		pkg.Log.WarnCtx(c, "[AUTH-WARN]: Could not find OTP")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})

		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to get otp from table", err)
		return
	}

	// token generation and setting cookie
	authToken, err := pkg.CreateAuthToken(row.ID, email, true, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		return
	}

	refreshToken, err := pkg.CreateRefreshToken(row.ID.String(), email, true, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})

		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to add refresh token", err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})

		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to commit transaction", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verification completed successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func ResendUserOtp(c *gin.Context) {
	email := c.GetString("email")
	if email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})

		pkg.Log.FatalCtx(c, "[AUTH-ERROR]: Email not in ctx after middleware", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})

		pkg.Log.FatalCtx(c, "[AUTH-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer conn.Release()

	q := db.New()

	results, err := q.GetStudentOtpQuery(ctx, conn, email)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No OTP available beyond expiry time.",
		})

		pkg.Log.WarnCtx(c, "[AUTH-ERROR]: All OTPs expired")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Could not fetch OTP for onboarding", err)
		return
	}

	// Resend OTP via Mail
	emailReq := mail.EmailRequest{
		To:      []string{email},
		Subject: fmt.Sprintf("Resend OTP - Anokha 2025 - %d", time.Now().UnixMilli()),
		Type:    "otp",
		Data: mail.OTPTemplateData{
			UserName: results.Name,
			OTP:      strings.Split(results.Otp, ""),
		},
	}
	if err := mail.Mail.Enqueue(&emailReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to add mail to mailing queue", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "OTP resent to user email successfully",
		"expiry_at": results.ExpiryAt.Time,
	})
	pkg.Log.SuccessCtx(c)
}
