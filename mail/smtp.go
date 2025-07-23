package mail

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	gomail "gopkg.in/gomail.v2"
)

// there is no need for mutex or synchronization because each worker will have a mailer (SMTP connection)
// which implicitly means, we should not have more than 4 or 5 workers
type Mailer struct {
	dialer    *gomail.Dialer
	sender    gomail.SendCloser
	connected bool
}

func NewMailer() *Mailer {
	return &Mailer{
		dialer: &gomail.Dialer{
			Host:     cmd.Env.SMTPHost,
			Port:     cmd.Env.SMTPPort,
			Username: cmd.Env.SMTPUsername,
			Password: cmd.Env.SMTPPassword,
		},
	}
}

type EmailRequest struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Type    string   `json:"type"`
	Data    any      `json:"data"`
}

// function to send mail from to SMTP server
func (m *Mailer) Send(toAddresses []string, subject string, emailType string, data any) error {
	var err error

	for attempt := range 3 {
		if !m.connected {
			s, dialErr := m.dialer.Dial()
			if dialErr != nil {
				err = fmt.Errorf("SMTP dial failed: %w", dialErr)
			} else {
				m.sender = s
				m.connected = true
			}
		}

		if err == nil {
			var body string
			body, err = getTemplate(emailType, data)
			if err != nil {
				break
			}

			msg := gomail.NewMessage()
			msg.SetHeader("From", m.dialer.Username)
			msg.SetHeader("To", toAddresses...)
			msg.SetHeader("Subject", subject)
			msg.SetBody("text/html", body)

			err = gomail.Send(m.sender, msg)
			if err == nil {
				pkg.Log.Info(fmt.Sprintf("Email sent successfully:\tTYPE: %s\tTO: %s", emailType, strings.Join(toAddresses, ", ")))
				return nil
			}

			// Resetting connection if sending failed
			m.connected = false
			_ = m.sender.Close()
		}

		backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		pkg.Log.Warn(fmt.Sprintf("Retry %d for email send in %v: %v", attempt+1, backoff, err))
		time.Sleep(backoff)
	}

	return fmt.Errorf("email send failed after retries: %w", err)
}
