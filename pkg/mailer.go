package pkg

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

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
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", fileName, err)
	}
	return rendered.String(), nil
}
