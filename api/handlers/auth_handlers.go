package handlers

import (
	"encoding/json"
	jwttoken "himakiwa/handlers/jwt"
	"himakiwa/handlers/utils"
	"himakiwa/services"
	"himakiwa/services/database"
	"himakiwa/services/email"
	"net/http"
	"os"
	"strconv"
)

var (
	ErrInvalidAuth   = "invalid auth request"
	ErrIncorrectAuth = "incorrect auth request"
)

func NewAutenticateHandlers() AutenticatehHandlers {
	var sendVCode = func(tu *database.TUser) error {
		return email.NewEmailServices().SendVCode(tu.Email, tu.VCode)
	}
	return AutenticatehHandlers{SendVCode: sendVCode}
}

type AutenticatehHandlers struct {
	SendVCode func(*database.TUser) error
}

type BasicResp struct {
	Jwt string
}

type LogininBody struct {
	Email    string `validate:"required,email,max=50"`
	Password string `validate:"required,alphanumary,min=8,max=24"`
}

func (ah *AutenticatehHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &LogininBody{}
	if err := utils.DecodeBody(r, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userServices := services.NewRepositoryServices().GetUser()

	// log in
	tu, err := userServices.Login(b.Email, b.Password)
	if err != nil {
		http.Error(w, ErrInvalidAuth, http.StatusBadRequest)
		return
	}

	// send vcode
	if err = ah.SendVCode(tu); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set jwt
	ah.return10JWT(w, tu)
}

type SigninBody struct {
	Name     string `validate:"omitempty,max=20"`
	Email    string `validate:"required,email,max=50"`
	Password string `validate:"required,alphanumary,min=8,max=24"`
}

func (ah *AutenticatehHandlers) SigninHandler(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &SigninBody{}
	if err := utils.DecodeBody(r, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userServices := services.NewRepositoryServices().GetUser()

	// sign in
	tu, err := userServices.Signin(b.Email, b.Password, b.Name)
	if err != nil {
		http.Error(w, ErrInvalidAuth, http.StatusBadRequest)
		return
	}

	// send vcode
	if err = ah.SendVCode(tu); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set jwt
	ah.return10JWT(w, tu)
}

type VerificateBody struct {
	Code string `validate:"required,len=6"`
	Jwt  string `validate:"required"`
}

func (*AutenticatehHandlers) VerificateHandler(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &VerificateBody{}
	if err := utils.DecodeBody(r, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// parse jwt
	secret := os.Getenv("JWT_SECRET")
	jt := jwttoken.NewJwt(secret)
	claim, err := jt.ParseToken(b.Jwt)
	if err != nil {
		http.Error(w, ErrInvalidAuth, http.StatusBadRequest)
		return
	}

	// parse userId
	userId, err := strconv.Atoi(claim.ID)
	if err != nil {
		http.Error(w, ErrIncorrectAuth, http.StatusBadRequest)
		return
	}

	//verificate
	userServices := services.NewRepositoryServices().GetUser()

	_, err = userServices.Verificate(userId, b.Code)
	if err != nil {
		http.Error(w, ErrInvalidAuth, http.StatusBadRequest)
		return
	}
	// set jwt-token
	token, err := jt.Generate(claim.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SetJWTCookie(w, token)
	w.WriteHeader(http.StatusOK)
}

func (*AutenticatehHandlers) return10JWT(w http.ResponseWriter, tu *database.TUser) {
	// create jwt-token
	secret := os.Getenv("JWT_SECRET")
	jt := jwttoken.New10minJwt(secret)
	token, err := jt.Generate(strconv.Itoa(tu.Id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set jwt-token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&BasicResp{token})
}
