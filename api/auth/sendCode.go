package auth

import (
	"fmt"
	"himakiwa/utils"
)

// send email to verificate two-step code.
func SendCode(to, code string) error {
	subject := "two-step verification"
	message := fmt.Sprintf("please entry '%s' in app.", code)

	email := &utils.Email{
		To:      to,
		Subject: subject,
		Message: message,
	}
	return email.Send()
}
