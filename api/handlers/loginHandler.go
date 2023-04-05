package handlers

import (
	"encoding/json"
	"himakiwa/database"
	"himakiwa/handlers/decode"
	"net/http"
)

type LogininBody struct {
	Email    string `validate:"required,email,max=50"`
	Password string `validate:"required,alphanumary,min=8,max=24"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &LogininBody{}
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

	// auth
	u := database.LoginUser{Email: b.Email, Password: b.Password}
	du, err := u.Login(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send code
	// send id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(du.Id)
}
