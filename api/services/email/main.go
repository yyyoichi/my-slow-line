package email

import "fmt"

type EmailServices struct{}

func NewEmailServices() *EmailServices {
	return &EmailServices{}
}

func (e *EmailServices) SendVCode(email, vcode string) error {
	subject := "two-step verification"
	message := fmt.Sprintf("Please entry '%s' in app.\nIf you do not recognize this e-mail address, please discard it.\n\nThank you.\n\n\nCtrl +", vcode)

	mailer := &build{
		To:      email,
		Subject: subject,
		Message: message,
	}
	return mailer.send()
}
