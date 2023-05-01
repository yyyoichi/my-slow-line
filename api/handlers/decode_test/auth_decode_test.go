package decode_test

import (
	"bytes"
	"fmt"
	"himakiwa/handlers"
	"himakiwa/handlers/utils"
	"net/http/httptest"
	"testing"
)

func TestLoginBody(t *testing.T) {
	email := "testbody@example.com"
	pass := "pa55word"
	test := []struct {
		email string
		pass  string
		err   error
	}{
		{
			email: email,
			pass:  pass,
			err:   nil,
		},
		{
			email: "demo@",
			pass:  pass,
			err:   utils.ErrInvalidBody,
		},
		{
			email: fmt.Sprintf("%s%s0123456789+", email, email),
			pass:  pass,
			err:   utils.ErrInvalidBody,
		},
		{
			email: "",
			pass:  pass,
			err:   utils.ErrInvalidBody,
		},
		{
			email: email,
			pass:  "password",
			err:   utils.ErrInvalidBody,
		},
		{
			email: email,
			pass:  "01234567",
			err:   utils.ErrInvalidBody,
		},
		{
			email: email,
			pass:  "pa55w",
			err:   utils.ErrInvalidBody,
		},
		{
			email: email,
			pass:  fmt.Sprintf("%s%s%s+", pass, pass, pass),
			err:   utils.ErrInvalidBody,
		},
		{
			email: email,
			pass:  "",
			err:   utils.ErrInvalidBody,
		},
	}

	for i, tt := range test {
		buf := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, tt.email, tt.pass)
		body := bytes.NewBufferString(buf)
		req := httptest.NewRequest("POST", "/", body)

		b := &handlers.LogininBody{}
		if err := utils.DecodeBody(req, b); err != tt.err {
			t.Errorf("%d: expected err is '%v' but got='%v'", i, tt.err, err)
		}
	}
}

func TestSigninBody(t *testing.T) {
	email := "testbody@example.com"
	pass := "pa55word"
	name := "test-user1"
	test := []struct {
		name string
		err  error
	}{
		{
			name: name,
			err:  nil,
		},
		{
			err: nil,
		},
		{
			err: nil,
		},
		{
			name: fmt.Sprintf("%s%s%s+", name, name, name),
			err:  utils.ErrInvalidBody,
		},
	}
	for i, tt := range test {
		buf := fmt.Sprintf(`{"email":"%s", "password":"%s", "name":"%s"}`, email, pass, tt.name)
		body := bytes.NewBufferString(buf)
		req := httptest.NewRequest("POST", "/", body)

		b := &handlers.SigninBody{}
		if err := utils.DecodeBody(req, b); err != tt.err {
			t.Errorf("%d: expected err is '%v' but got='%v'", i, tt.err, err)
		}
	}
}

func TestVerificateBody(t *testing.T) {
	code := "012345"
	jwt := "jwtoken"
	test := []struct {
		code string
		jwt  string
		err  error
	}{
		{
			code: code,
			jwt:  jwt,
			err:  nil,
		},
		{
			code: "01234",
			jwt:  jwt,
			err:  utils.ErrInvalidBody,
		},
		{
			jwt: jwt,
			err: utils.ErrInvalidBody,
		},
		{
			code: code,
			err:  utils.ErrInvalidBody,
		},
	}
	for i, tt := range test {
		buf := fmt.Sprintf(`{"code":"%s", "jwt":"%s"}`, tt.code, tt.jwt)
		body := bytes.NewBufferString(buf)
		req := httptest.NewRequest("POST", "/", body)

		b := &handlers.VerificateBody{}
		if err := utils.DecodeBody(req, b); err != tt.err {
			t.Errorf("%d: expected err is '%v' but got='%v'", i, tt.err, err)
		}
	}
}
