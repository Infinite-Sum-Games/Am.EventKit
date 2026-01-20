package mail

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/Infinite-Sum-Games/Am.EventKit/cmd"
	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
	gomail "gopkg.in/gomail.v2"
)

// there is no need for mutex or synchronization because each worker will have a mailer (SMTP connection)
// which implicitly means, we should not have more than 4 or 5 workers
type Mailer struct {
	dialer *gomail.Dialer
}

func NewMailer() *Mailer {
	return &Mailer{
		dialer: &gomail.Dialer{
			Host:     cmd.Env.SMTPHost,
			Port:     cmd.Env.SMTPPort,
			Username: cmd.Env.SMTPUsername,
			Password: cmd.Env.SMTPPassword,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS12,
			},
		},
	}
}

type EmailRequest struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Type    string   `json:"type"` // "otp" | "welcome" | "event-reg"
	Data    any      `json:"data"`
	Retries int
}

// function to send mail from to SMTP server
func (m *Mailer) Send(
	toAddresses []string,
	subject, emailType string,
	data any,
	retryCount int,
) error {
	var lastErr error

	if retryCount > 0 {
		body, err := getTemplate(emailType, data)
		if err != nil {
			return err
		}

		msg := gomail.NewMessage()
		msg.SetHeader("From", m.dialer.Username)
		msg.SetHeader("To", toAddresses...)
		msg.SetHeader("Subject", subject)
		msg.SetBody("text/html", body)

		// 🔑 Dial per send
		// sender, err := m.dialer.Dial()
		// if err != nil {
		// 	lastErr = err
		// 	time.Sleep(time.Duration(attempt+1) * 5 * time.Second)
		// 	continue
		// }
		//
		// err = gomail.Send(sender, msg)
		// _ = sender.Close()
		err = m.dialer.DialAndSend(msg)
		if err == nil {
			pkg.Log.Info(
				fmt.Sprintf("Email sent successfully: %s - Retry count: %d",
					strings.Join(toAddresses, ", "),
					3-retryCount,
				))
			return nil
		}

		lastErr = err
	} else {
		pkg.Log.Error("Max retries reached for email send", nil)
		return nil
	}
	return fmt.Errorf("email send failed after retries: %w", lastErr)
}
