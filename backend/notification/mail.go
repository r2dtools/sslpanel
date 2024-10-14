package notification

import (
	"backend/config"
	"backend/utils"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"

	gomail "gopkg.in/mail.v2"
)

// EmailNotification sends notifications
type EmailNotification struct{}

// SendPlainNotification sends plain notification
func (n *EmailNotification) SendPlainNotification(recepient, subject, message string) error {
	return n.sendNotification(recepient, subject, message, "text/plain")
}

// SendEmailNotification sends notification
func (n *EmailNotification) sendNotification(recepient, subject, message, bType string) error {
	aConfig := config.GetConfig()
	m := gomail.NewMessage()
	m.SetHeader("From", aConfig.AmsEmailAddress)
	m.SetHeader("To", recepient)
	m.SetHeader("Subject", subject)
	m.SetBody(bType, message)

	d := gomail.NewDialer(aConfig.SMTPHost, aConfig.SMTPPort, aConfig.AmsEmailAddress, aConfig.AmsEmailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send notification to recepient '%s': %v", recepient, err)
	}

	return nil
}

// SendNotification creates and sends notification to email address
// tplPath is a template relative to var directory
func (n *EmailNotification) CreateAndSendPlainNotification(name, tplPath, email, subject string, data interface{}) error {
	tplContent, err := utils.ReadVarFile(tplPath)

	if err != nil {
		return err
	}

	report, err := template.New(name).Parse(string(tplContent))

	if err != nil {
		return fmt.Errorf("could not parse notification template: %v", err)
	}

	var tpl bytes.Buffer

	if err = report.Execute(&tpl, data); err != nil {
		return err
	}

	if err = n.SendPlainNotification(email, subject, tpl.String()); err != nil {
		return err
	}

	return nil
}
