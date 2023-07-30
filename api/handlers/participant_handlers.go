package handlers

import (
	"himakiwa/handlers/utils"
	"himakiwa/services/database"
	"himakiwa/services/sessions"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

///////////////////////////////////
//// participant at handlers //////
///////////////////////////////////

type ParticipantsAtHandlers struct {
	Use sessions.UseSessionServicesFunc
}

func NewParticipantsAtHandlers(use sessions.UseSessionServicesFunc) func(http.ResponseWriter, *http.Request) {
	pah := &ParticipantsAtHandlers{use}
	return pah.ParticipantsAtHandlers
}

func (pah *ParticipantsAtHandlers) ParticipantsAtHandlers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		pah.PostParticipantsAtHandler(w, r)
	case "PUT":
		pah.PutParticipantsAtHandler(w, r)
	default:
		http.Error(w, ErrUnExpcetedMethod, http.StatusBadRequest)
	}
}

// invite user

type PostParticipantsAtBody struct {
	UserID int `validate:"required"`
}

// TODO: write
func (pah *ParticipantsAtHandlers) PostParticipantsAtHandler(w http.ResponseWriter, r *http.Request) {
}

// update status

type PoutParticipantAtBody struct {
	UserID int                         `validate:"required"`
	Status database.TParticipantStatus `validate:"required"`
}

func (pah *ParticipantsAtHandlers) PutParticipantsAtHandler(w http.ResponseWriter, r *http.Request) {
	// read param
	vars := mux.Vars(r)
	sessionStrID := vars["sessionID"]
	if sessionStrID == "" {
		http.Error(w, ErrNotExistRecruit.Error(), http.StatusInternalServerError)
		return
	}
	sessionID, err := strconv.Atoi(sessionStrID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// parse body
	b := &PoutParticipantAtBody{}
	if err := utils.DecodeBody(r, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// read context
	userID, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// session services
	sessionServices := pah.Use(userID)

	// update
	err = sessionServices.UpdateParticipantStatusAt(sessionID, b.UserID, b.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
