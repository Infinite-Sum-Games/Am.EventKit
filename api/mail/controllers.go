package mail

import (
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/mail"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	mailerSvc *mail.MailerService
}

func NewController(mailerSvc *mail.MailerService) *Controller {
	return &Controller{
		mailerSvc: mailerSvc,
	}
}

type TestEmailRequest struct {
	To      []string `json:"to" binding:"required"`
	Subject string   `json:"subject" binding:"required"`
	Type    string   `json:"type" binding:"required"`
}

func (c *Controller) SendTestEmail(ctx *gin.Context) {
	var req TestEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var data any
	switch req.Type {
	case "otp":
		data = &mail.OTPTemplateData{
			UserName: "Test User",
			OTP:      []string{"1", "2", "3", "4", "5", "6"},
		}
	case "welcome":
		data = &mail.WelcomeTemplateData{
			UserName: "Test User",
		}
	case "event-reg":
		data = &mail.RegistrationData{
			UserName:      "Test User",
			EventName:     "Test Event",
			EventDate:     "Tomorrow",
			EventTime:     "10 AM",
			EventLocation: "Here",
		}
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email type"})
		return
	}

	emailReq := mail.EmailRequest{
		To:      req.To,
		Subject: req.Subject,
		Type:    req.Type,
		Data:    data,
	}

	if err := c.mailerSvc.Enqueue(&emailReq); err != nil {
		pkg.Log.Error("failed to enqueue email", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "email enqueued successfully"})
}
