package mail

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"html/template"
	"path/filepath"
)

func init() {
	gob.Register(OTPTemplateData{})
	gob.Register(WelcomeTemplateData{})
	gob.Register(RegistrationData{})
}

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

// funtion to get template of html based on email type
func getTemplate(emailType string, data any) (string, error) {
	basePath := filepath.Join("pkg", "template")
	templateFiles := map[string]string{
		"otp":       "otp.html",
		"welcome":   "welcome.html",
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
