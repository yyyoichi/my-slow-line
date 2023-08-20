package email

import "fmt"

type UseEmailServicesFunc func(to string) *EmailServices
type EmailServices struct {
	mail mailInterface
	to   string
}

func NewEmailServices() UseEmailServicesFunc {
	email := &EmailServices{mail: &mail{}}
	return func(to string) *EmailServices {
		email.to = to
		return email
	}
}

func (es *EmailServices) SendVCode(vcode string) error {
	subject := "two-step verification"
	message := fmt.Sprintf("Please entry '%s' in app.\nIf you do not recognize this e-mail address, please discard it.\n\nThank you.\n\n\nCtrl +", vcode)
	return es.mail.sendMail(es.to, subject, message)
}
