package notification

import (
	"backend/config"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	gomail "gopkg.in/mail.v2"
)

type EmailNotificationService struct {
	Config *config.Config
}

func (n EmailNotificationService) SendPlainNotification(recepient, subject, message string) error {
	return n.sendNotification(recepient, subject, message, "text/plain")
}

func (n *EmailNotificationService) sendNotification(recepient, subject, message, bType string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", n.Config.AmsEmailAddress)
	m.SetHeader("To", recepient)
	m.SetHeader("Subject", subject)
	m.SetBody(bType, message)

	d := gomail.NewDialer(n.Config.SMTPHost, n.Config.SMTPPort, n.Config.AmsEmailAddress, n.Config.AmsEmailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send notification to recepient '%s': %v", recepient, err)
	}

	return nil
}

func (n EmailNotificationService) CreateAndSendPlainNotification(name, tplPath, email, subject string, data interface{}) error {
	tplContent, err := n.readTemplateFile(tplPath)

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

func (n EmailNotificationService) readTemplateFile(subpath string) (string, error) {
	tplPath := filepath.Join(n.Config.GetVarDirAbsPath(), subpath)
	content, err := os.ReadFile(tplPath)

	if err != nil {
		return "", fmt.Errorf("could not read file '%s': %v", tplPath, err)
	}

	return string(content), nil
}
