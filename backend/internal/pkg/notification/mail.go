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

type ErrSendEmail struct {
	Message string
}

func (e ErrSendEmail) Error() string {
	return e.Message
}

type EmailNotificationService struct {
	Config *config.Config
}

func (n EmailNotificationService) SendHtmlNotification(recepient, subject, message string) error {
	err := n.sendNotification(recepient, subject, message, "text/html")

	if err != nil {
		return ErrSendEmail{Message: fmt.Sprintf("failed to send email notification: %v", err)}
	}

	return nil
}

func (n *EmailNotificationService) sendNotification(recepient, subject, message, bType string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", n.Config.AmsEmailAddress)
	m.SetHeader("To", recepient)
	m.SetHeader("Subject", subject)
	m.SetBody(bType, message)

	d := gomail.NewDialer(n.Config.SMTPHost, n.Config.SMTPPort, n.Config.AmsEmailAddress, n.Config.AmsEmailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}

func (n EmailNotificationService) CreateAndSendHtmlNotification(name, tplPath, email, subject string, data interface{}) error {
	message, err := n.createNotification(name, tplPath, subject, data)

	if err != nil {
		return err
	}

	return n.SendHtmlNotification(email, subject, message)
}

func (n EmailNotificationService) createNotification(name, tplPath, subject string, data interface{}) (string, error) {
	tplContent, err := n.readTemplateFile(tplPath)

	if err != nil {
		return "", err
	}

	report, err := template.New(name).Parse(string(tplContent))

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err = report.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func (n EmailNotificationService) readTemplateFile(subpath string) (string, error) {
	tplPath := filepath.Join(n.Config.GetVarDirAbsPath(), subpath)
	content, err := os.ReadFile(tplPath)

	if err != nil {
		return "", err
	}

	return string(content), nil
}
