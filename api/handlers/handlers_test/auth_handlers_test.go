package handlers_test

import (
	"fmt"
	"himakiwa/handlers"
	jwttoken "himakiwa/handlers/jwt"
	"himakiwa/services"
	"himakiwa/services/email"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	var emailAddr = "test1@example.com"
	var password = "pa55word"
	var userID = 1
	var useEmailServices = email.NewEmailServicesMock()
	var useRepositoryServices = services.NewRepositoryServicesMock()
	var ah = handlers.NewAutenticateHandlers(useEmailServices, useRepositoryServices)

	// create server
	server := NewHttpMock(ah.LoginHandler)
	defer server.Close()

	// request
	body := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, emailAddr, password)
	resp := server.Post(t, body)

	defer resp.BodyClose()

	verificateJwt(t, resp, userID, emailAddr)
}
func TestSigninHandler(t *testing.T) {
	var password = "pa55word"
	var useEmailServices = email.NewEmailServicesMock()
	var useRepositoryServices = services.NewRepositoryServicesMock()
	var ah = handlers.NewAutenticateHandlers(useEmailServices, useRepositoryServices)
	// create server
	server := NewHttpMock(ah.SigninHandler)
	defer server.Close()

	// request
	body := fmt.Sprintf(`{"email":"%s", "password":"%s", "name":""}`, "test999@example.com", password)
	resp := server.Post(t, body)

	defer resp.BodyClose()

	verificateJwt(t, resp, 3, "test999@example.com")
}

func TestVerificateHandler(t *testing.T) {
	var emailAddr = "test1@example.com"
	var password = "pa55word"
	var useEmailServices = email.NewEmailServicesMock()
	var useRepositoryServices = services.NewRepositoryServicesMock()
	var ah = handlers.NewAutenticateHandlers(useEmailServices, useRepositoryServices)

	// login At
	useRepositoryServices(1).UserServices.Login(emailAddr, password)
	// create jwt
	jt := jwttoken.New10minJwt(os.Getenv("JWT_SECRET"))
	token, err := jt.Generate("1")
	if err != nil {
		t.Error(err)
	}

	// create server
	vcode, err := useRepositoryServices(1).UserServices.RefreshVCode(1)
	if err != nil {
		t.Error(err)
	}
	server := NewHttpMock(ah.VerificateHandler)
	defer server.Close()

	// request
	body := fmt.Sprintf(`{"code":"%s", "jwt":"%s"}`, vcode, token)
	resp := server.Post(t, body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status-code is 200 but got='%d'", resp.StatusCode)
	}
}

func verificateJwt(t *testing.T, resp Resp, expID int, expEmail string) {
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

	if actUserID, err := strconv.Atoi(claim.ID); err != nil {
		t.Error(err)
	} else if actUserID != expID {
		t.Errorf("Expected userID is '%d', but got='%d'", expID, actUserID)
	}
}
