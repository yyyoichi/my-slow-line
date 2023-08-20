package email

import (
	"fmt"
)

func NewEmailServicesMock() UseEmailServicesFunc {
	email := &EmailServices{mail: &mailMock{}}
	return func(to string) *EmailServices {
		email.to = to
		return email
	}
}

type mailMock struct{}

func (m *mailMock) sendMail(to, subject, message string) error {
	fmt.Println("///////////////////////////////////////")
	fmt.Println("////// SEND WEBPUSH NOTIFICATION //////")
	fmt.Println("///////////////////////////////////////")
	fmt.Printf("// to: %s\n", to)
	fmt.Printf("// subject: %s\n", subject)
	fmt.Println("//-----------------------------------//")
	fmt.Printf("// message: %s\n", message)
	fmt.Println("//-----------------------------------//")
	fmt.Println("///////////////////////////////////////")
	return nil
}
