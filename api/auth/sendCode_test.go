package auth

import (
	"os"
	"testing"
)

func TestSendCode(t *testing.T) {
	email := os.Getenv("EMAIL_ADDRESS")
	code := "000000"
	if err := SendCode(email, code); err != nil {
		t.Error(err)
	}
}
