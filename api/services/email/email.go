package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type build struct {
	To      string
	Subject string
	Message string
}

func (b *build) send() error {
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

	body := fmt.Sprintf("TO: %s\r\nSubject: %s\r\n\n\n%s", b.To, b.Subject, b.Message)

	return smtp.SendMail(
		smtporigin,
		auth,
		from,
		[]string{b.To},
		[]byte(body),
	)
}
