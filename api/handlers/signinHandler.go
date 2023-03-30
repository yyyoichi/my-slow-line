package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"himakiwa/database"
	"himakiwa/handlers/decode"
	"himakiwa/utils"
	"net/http"
	"os"
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
	db, err := database.GetDatabase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write db
	u := &database.SignInUser{Name: b.Name, Email: b.Email, Password: b.Password}
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

	// set-cookie jwt-token
	jt := utils.NewJwt(os.Getenv("JWT_SECRET"))
	token, err := jt.Generate(fmt.Sprintf("%d", userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SetJWTCookie(w, token)
	w.WriteHeader(http.StatusOK)
}
