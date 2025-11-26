package api

import (
	"context"
	"net/http"
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to acquire DB connection", err)
		return
	}
	defer tx.Rollback(ctx)

	otpStr, otpSlice, err := pkg.GenerateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.ErrorCtx(c, "[AUTH-ERROR]: Failed to generate OTP for forgot password", err)
		return
	}

	q := db.New()
	result, err := q.PasswordChangeOtpQuery(ctx, tx, db.PasswordChangeOtpQueryParams{
		Email:    req.Email,
		Password: req.NewPassword,
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

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		pkg.Log.FatalCtx(c, "[AUTH-FATAL]: Failed to commit DB transaction", err)
		return
	}

	emailReq := mail.EmailRequest{
		To:      []string{result.Email},
		Subject: "Password Reset - Anokha 2025",
		Type:    "otp",
		Data: mail.OTPTemplateData{
			UserName: result.Name,
			OTP:      otpSlice,
		},
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully. Proceed to login.",
	})
	pkg.Log.SuccessCtx(c)
}
