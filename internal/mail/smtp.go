package mail

import (
	"fmt"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
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
	if !m.connected {
		s, err := m.dialer.Dial()
		if err != nil {
			return fmt.Errorf("SMTP dial failed: %w", err)
		}
		m.sender = s
		m.connected = true
	}

	body, err := getTemplate(emailType, data)
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.dialer.Username) //TODO: Should change this to official email of anokha from config
	msg.SetHeader("To", toAddresses...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	if err := gomail.Send(m.sender, msg); err != nil {
		m.connected = false
		_ = m.sender.Close()
		return fmt.Errorf("cannot send email: %w", err)
	}
	//TODO: Log email sent successful message with logger
	return nil
}
