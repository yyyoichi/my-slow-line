package handlers

import (
	"encoding/json"
	"errors"
	"himakiwa/auth"
	"himakiwa/database"
	"himakiwa/handlers/decode"
	"net/http"
)

var (
	ErrExistEmail = errors.New("already exist email")
)

type SigninBody struct {
	Name     string `validate:"omitempty,max=20"`
	Email    string `validate:"required,email,max=50"`
	Password string `validate:"required,alphanumary,min=8,max=24"`
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &SigninBody{}
	if err := decode.DecodeBody(r, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// connect db
	db := database.DB

	code := auth.GenerateRandomSixNumber()

	// write db
	u := &database.SignInUser{Name: b.Name, Email: b.Email, Password: b.Password, VerificationCode: code}
	userId, err := u.SignIn(db)
	if err != nil {
		result, err := database.ExistEmail(db, u.Email)
		if result {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ErrExistEmail)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send codes
	err = auth.SendCode(u.Email, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// send id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userId)
}
