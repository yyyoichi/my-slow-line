package handlers

import (
	"encoding/json"
	"himakiwa/auth"
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
	db := database.DB

	// auth
	u := database.LoginUser{Email: b.Email, Password: b.Password}
	du, err := u.Login(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId := du.Id

	// update login stamp and twostep code
	code := auth.GenerateRandomSixNumber()
	err = database.LogEntryStamp(db, userId, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send code
	err = auth.SendCode(u.Email, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// send id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userId)
}
