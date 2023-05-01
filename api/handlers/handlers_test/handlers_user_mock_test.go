package handlers_test

import (
	jwttoken "himakiwa/handlers/jwt"
	"himakiwa/services"
	"himakiwa/services/database"
	"os"
	"strconv"
	"testing"
)

type AuthedMock struct {
	Jwt   string
	TUser *database.TUser
}

var (
	userServices = services.NewRepositoryServices().GetUser()
)

func NewAuthedMock(t *testing.T) AuthedMock {

	// create user
	email := "mock@example.com"
	pass := "pa55word"
	tu, err := userServices.Signin(email, pass, "")
	if err != nil {
		t.Error(err)
	}

	// create jwt
	jt := jwttoken.NewJwt(os.Getenv("JWT_SECRET"))
	token, err := jt.Generate(strconv.Itoa(tu.Id))
	if err != nil {
		t.Error(err)
	}

	return AuthedMock{token, tu}
}

func (m *AuthedMock) Delete(t *testing.T) {
	if err := userServices.HardDelete(m.TUser.Id); err != nil {
		t.Error(err)
	}
}
