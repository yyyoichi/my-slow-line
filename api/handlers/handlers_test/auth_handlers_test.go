package handlers_test

import (
	"fmt"
	"himakiwa/handlers"
	jwttoken "himakiwa/handlers/jwt"
	"himakiwa/services"
	"himakiwa/services/database"
	"net/http"
	"os"
	"strconv"
	"testing"
)

var SendVCodeMock = func(tu *database.TUser) error {
	fmt.Printf("send %s to %s\n", tu.VCode, tu.Email)
	return nil
}

func init() {
	database.Connect()
}

func TestLoginHandler(t *testing.T) {
	// create user
	email := "handlertest@example.com"
	pass := "pa55word"
	us := services.NewRepositoryServices().GetUser()
	tu, err := us.Signin(email, pass, "")
	if err != nil {
		t.Error(err)
	}
	defer us.HardDelete(tu.Id)

	// create server
	ah := &handlers.AutenticatehHandlers{SendVCode: SendVCodeMock}
	server := NewHttpMock(ah.LoginHandler)
	defer server.Close()

	// request
	body := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, pass)
	resp := server.Post(t, body)

	defer resp.BodyClose()

	verificateJwt(t, resp, email)
}
func TestSigninHandler(t *testing.T) {
	// basic detail
	email := "handlertest@example.com"
	pass := "pa55word"

	ur := &database.UserRepository{}
	defer ur.HardDeleteByEmail(email)

	// create server
	ah := &handlers.AutenticatehHandlers{SendVCode: SendVCodeMock}
	server := NewHttpMock(ah.SigninHandler)
	defer server.Close()

	// request
	body := fmt.Sprintf(`{"email":"%s", "password":"%s", "name":""}`, email, pass)
	resp := server.Post(t, body)

	defer resp.BodyClose()

	verificateJwt(t, resp, email)
}

func TestVerificateHandler(t *testing.T) {
	// create user
	email := "handlertest@example.com"
	pass := "pa55word"
	us := services.NewRepositoryServices().GetUser()
	tu, err := us.Signin(email, pass, "")
	if err != nil {
		t.Error(err)
	}
	defer us.HardDelete(tu.Id)

	// create jwt
	jt := jwttoken.New10minJwt(os.Getenv("JWT_SECRET"))
	token, err := jt.Generate(strconv.Itoa(tu.Id))
	if err != nil {
		t.Error(err)
	}

	// create server
	ah := &handlers.AutenticatehHandlers{SendVCode: SendVCodeMock}
	server := NewHttpMock(ah.VerificateHandler)
	defer server.Close()

	// request
	body := fmt.Sprintf(`{"code":"%s", "jwt":"%s"}`, tu.VCode, token)
	resp := server.Post(t, body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status-code is 200 but got='%d'", resp.StatusCode)
	}
}

func verificateJwt(t *testing.T, resp Resp, expEmail string) {
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status-code is 200 but got='%d'", resp.StatusCode)
	}

	basisResp := &handlers.BasicResp{}
	if err := resp.BodyJson(basisResp); err != nil {
		t.Error(err)
	}

	token := basisResp.Jwt
	if token == "" {
		t.Errorf("token is ''")
	}

	secret := os.Getenv("JWT_SECRET")
	jt := jwttoken.New10minJwt(secret)

	claim, err := jt.ParseToken(token)
	if err != nil {
		t.Error(err)
	}

	userId, err := strconv.Atoi(claim.ID)
	if err != nil {
		t.Error(err)
	}

	userServices := services.NewRepositoryServices().GetUser()
	tu, err := userServices.Query(userId)
	if err != nil {
		t.Error(err)
	}

	if tu.Email != expEmail {
		t.Errorf("expected email is %s but got='%s'", expEmail, tu.Email)
	}

	if err = userServices.HardDelete(userId); err != nil {
		t.Error(err)
	}
}
