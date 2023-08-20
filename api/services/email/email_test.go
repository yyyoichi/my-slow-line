package email

import (
	"os"
	"testing"
)

func TestSendCode(t *testing.T) {
	email := os.Getenv("EMAIL_ADDRESS")
	code := "000000"
	es := NewEmailServices()(email)
	if err := es.SendVCode(code); err != nil {
		t.Error(err)
	}
}
