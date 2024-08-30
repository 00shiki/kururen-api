package mail

import (
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
)

type Service interface {
	SendMail(string, string, string, string) error
}

type MailClient struct {
	dialer *gomail.Dialer
}

func NewMailService() *MailClient {
	var (
		host     = os.Getenv("SMTP_HOST")
		portStr  = os.Getenv("SMTP_PORT")
		user     = os.Getenv("SMTP_AUTH_USERNAME")
		password = os.Getenv("SMTP_AUTH_PASSWORD")
	)
	port, _ := strconv.Atoi(portStr)
	return &MailClient{
		dialer: gomail.NewDialer(
			host,
			port,
			user,
			password,
		),
	}
}

func (mc *MailClient) SendMail(to, name, subject, content string) error {
	email := os.Getenv("SMTP_SENDER_EMAIL")
	m := gomail.NewMessage()
	m.SetAddressHeader("Cc", to, name)
	m.SetHeader("From", email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)
	if err := mc.dialer.DialAndSend(m); err != nil {
		log.Printf("error sending mail: %v", err)
		return err
	}
	return nil
}
