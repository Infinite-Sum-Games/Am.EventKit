package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Infinite-Sum-Games/Am.EventKit/cmd"
	db "github.com/Infinite-Sum-Games/Am.EventKit/db/gen"
	"github.com/Infinite-Sum-Games/Am.EventKit/mail"
	"github.com/Infinite-Sum-Games/Am.EventKit/models"
	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func ForgotUserPasswordCsrf(c *gin.Context) {
	csrfToken, err := pkg.CreateCsrfToken("forgot.password@anokha.edu", c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to generate CSRF token", err)
		return
	}

	pkg.SetCsrfCookie(c, csrfToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset action initiated successfully",
		"key":     csrfToken,
	})
	pkg.Log.SuccessCtx(c)
}

func ConfirmPasswordChangeCsrf(c *gin.Context) {
	csrfToken, err := pkg.CreateCsrfToken("forgot.password@anokha.edu", c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to generate CSRF token", err)
		return
	}

	pkg.SetCsrfCookie(c, csrfToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset confirm action initiated successfully.",
		"key":     csrfToken,
	})
	pkg.Log.SuccessCtx(c)
}

func ForgotUserPassword(c *gin.Context) {
	req, ok := pkg.ValidateRequest[models.ForgetPasswordRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "AUTH") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "AUTH")

	otpStr, otpSlice, err := pkg.GenerateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to generate OTP for forgot password", err)
		return
	}

	password, err := pkg.Hash(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.FatalCtx(c, "[AUTH-ERROR]: Failed to hash password", err)
		return
	}

	q := db.New()
	result, err := q.PasswordChangeOtpQuery(ctx, tx, db.PasswordChangeOtpQueryParams{
		Email:    req.Email,
		Password: password,
		Otp:      otpStr,
		ExpiryAt: pgtype.Timestamp{
			Time:  time.Now().Add(5 * time.Minute),
			Valid: true,
		},
	})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No registered user with given email",
		})
		pkg.Log.WarnCtx(c, "[AUTH-WARN]: Unregistered user attempts to change password")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to acquire DB connection", err)
		return
	}

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "AUTH") {
		return
	}

	tempToken := pkg.CreateTempToken(result.Email)
	pkg.SetTempCookie(c, tempToken)

	emailReq := mail.EmailRequest{
		To:      []string{result.Email},
		Subject: "Password Reset OTP - Pragati 2026",
		Type:    "otp",
		Data: &mail.OTPTemplateData{
			UserName: result.Name,
			OTP:      otpSlice,
		},
		Retries: mail.MaxRetryCount,
	}
	err = mail.Mail.Enqueue(&emailReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to add mail to mailing queue", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset in progress. Please check email for OTP.",
	})
	pkg.Log.SuccessCtx(c)
}

func ConfirmPasswordChange(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "AUTH")
	if !ok {
		return
	}

	req, ok := pkg.ValidateRequest[models.OtpRequest](c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := cmd.DBPool.Begin(ctx)
	if pkg.HandleDbTxnErr(c, err, "AUTH") {
		return
	}
	defer pkg.RollbackTx(c, tx, ctx, "AUTH")

	q := db.New()

	_, err = q.ConfirmPasswordChangeOtpQuery(ctx, tx,
		db.ConfirmPasswordChangeOtpQueryParams{
			Email: email,
			Otp:   req.Otp,
		})
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Given OTP is invalid or expired",
		})
		pkg.Log.WarnCtx(c, "[AUTH-WARN]: Time expired for password change")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.FatalCtx(c, "[AUTH-FATAL]: Failed to update password after OTP", err)
		return
	}

	err = tx.Commit(ctx)
	if pkg.HandleDbTxnCommitErr(c, err, "AUTH") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully. Proceed to login.",
	})
	pkg.Log.SuccessCtx(c)
}

func ResendPasswordChangeOtp(c *gin.Context) {
	email, ok := pkg.GrabEmail(c, "AUTH")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if pkg.HandleDbAcquireErr(c, err, "AUTH") {
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.ResendPasswordChangeOtpQuery(ctx, conn, email)
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
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to retrive OTP for resending", err)
		return
	}

	emailReq := mail.EmailRequest{
		To:      []string{email},
		Subject: fmt.Sprintf("Resend Password Reset OTP - Pragati 2026 - %d", time.Now().UnixMilli()),
		Type:    "otp",
		Data: &mail.OTPTemplateData{
			UserName: result.Name,
			OTP:      strings.Split(result.Otp, ""),
		},
		Retries: mail.MaxRetryCount,
	}
	if err := mail.Mail.Enqueue(&emailReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to add mail to mailing queue", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Password reset OTP resent. Please check email for OTP.",
		"expiry_at": result.ExpiryAt.Time,
	})
	pkg.Log.SuccessCtx(c)
}
