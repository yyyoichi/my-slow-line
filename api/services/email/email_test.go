package email

import (
	"os"
	"testing"
)

func TestSendCode(t *testing.T) {
	email := os.Getenv("EMAIL_ADDRESS")
	code := "000000"
	if err := NewEmailServices().SendVCode(email, code); err != nil {
		t.Error(err)
	}
}
