package handlers

import (
	"errors"
	"fmt"
	"himakiwa/handlers/utils"
	"himakiwa/services/sessions"
	"net/http"
	"strconv"
)

var (
	ErrCannotAccessSession = errors.New("cannot access session")
)

type SessionPublicKeyHandlers struct {
	Use sessions.UseSessionServicesFunc
}

func NewSessionPublicKeyHandlers(use sessions.UseSessionServicesFunc) func(http.ResponseWriter, *http.Request) {
	pkh := &SessionPublicKeyHandlers{use}
	return pkh.PostSessionPublicKey
}

type PostSessionPublicKeyBody struct {
	SessionID int    `validate:"required"`
	SenderID  int    `validate:"required"`
	PublicKey string `validate:"required"`
}

func (pkh *SessionPublicKeyHandlers) PostSessionPublicKey(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &PostSessionPublicKeyBody{}
	if err := utils.DecodeBody(r, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// read context
	userID, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// session services
	sessionServices := pkh.Use(userID)

	// validation
	userIDIsJoined, err := sessionServices.IsJoined(b.SessionID, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	senderIDIsJoined, err := sessionServices.IsJoined(b.SessionID, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !userIDIsJoined || !senderIDIsJoined {
		fmt.Println(err)
		http.Error(w, ErrCannotAccessSession.Error(), http.StatusBadRequest)
		return
	}

	// send push
}
