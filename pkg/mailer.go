package pkg

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

type OTPTemplateData struct {
	UserName string
	OTP      []string
}

type WelcomeTemplateData struct {
	UserName string
}

type RegistrationData struct {
	UserName      string
	EventName     string
	EventDate     string
	EventTime     string
	EventLocation string
}

func SendMail(toAddresses []string, subject string, emailType string, data any) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "thanuskumaara@gmail.com") //TODO: Should change this to official email of anokha from config
	m.SetHeader("To", toAddresses...)
	m.SetHeader("Subject", subject)
	body, err := getTemplate(emailType, data)
	if err != nil {
		return err
	}
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		//TODO: should setup from config file
		"",
		25,
		"",
		"",
	)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("cannot send email: %w", err)
	}
	//TODO: Log email sent successful message with logger
	return nil
}

func getTemplate(emailType string, data any) (string, error) {
	basePath := filepath.Join("pkg", "template")
	templateFiles := map[string]string{
		"otp":     "otp.html",
		"welcome": "welcome.html",
		"event-reg": "event-registration.html",
	}
	fileName, ok := templateFiles[emailType]
	if !ok {
		return "", fmt.Errorf("unsupported email type: %s", emailType)
	}
	templatePath := filepath.Join(basePath, fileName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("unable to parse html file %s: %w", fileName, err)
	}
	var rendered bytes.Buffer
	switch emailType {
	case "otp":
		otpData, ok := data.(OTPTemplateData)
		if !ok {
			return "", fmt.Errorf("invalid data type for 'otp' email template")
		}
		err = tmpl.Execute(&rendered, otpData)

	case "welcome":
		welcomeData, ok := data.(WelcomeTemplateData)
		if !ok {
			return "", fmt.Errorf("invalid data type for 'welcome' email template")
		}
		err = tmpl.Execute(&rendered, welcomeData)

	case "event-reg":
		eventData, ok := data.(RegistrationData)
		if !ok {
			return "", fmt.Errorf("invalid data type for 'event-reg' email template")
		}
		err = tmpl.Execute(&rendered, eventData)
	}
	
	if err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", fileName, err)
	}
	return rendered.String(), nil
}
