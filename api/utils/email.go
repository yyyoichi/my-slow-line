package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

type Email struct {
	To      string
	Subject string
	Message string
}

func (e *Email) Send() error {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("APP_PASSWORD")
	smtphost := os.Getenv("SMTP_HOST")
	smtpport := os.Getenv("SMTP_PORT")

	smtporigin := fmt.Sprintf("%s:%s", smtphost, smtpport)

	auth := smtp.PlainAuth(
		"",
		from,
		password,
		smtphost,
	)

	body := fmt.Sprintf("TO: %s\r\nSubject: %s\r\n\n\n%s", e.To, e.Subject, e.Message)

	return smtp.SendMail(
		smtporigin,
		auth,
		from,
		[]string{e.To},
		[]byte(body),
	)
}
