package handlers

import (
	"himakiwa/services/database"
	"net/http"
)

///////////////////////////////////
//// participant at handlers //////
///////////////////////////////////

type ParticipantAtHandlers struct{}

func NewParticipantAtHandlers() func(http.ResponseWriter, *http.Request) {
	cah := &ChatAtHandlers{}
	return cah.ChatAtHandlers
}

func (pah *ParticipantAtHandlers) ParticipantAtHandlers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		pah.PostParticipantAtHandler(w, r)
	case "PUT":
		pah.PutParticipantAtHandler(w, r)
	default:
		http.Error(w, ErrUnExpcetedMethod, http.StatusBadRequest)
	}
}

// invite user

type PostParticipantAtHandlerBody struct {
	UserID int `validate:"required"`
}

// TODO: write
func (pah *ParticipantAtHandlers) PostParticipantAtHandler(w http.ResponseWriter, r *http.Request) {
}

// update status

type PoutParticipantAtHandlerBody struct {
	UserID int                         `validate:"required"`
	Status database.TParticipantStatus `validate:"required"`
}

func (pah *ParticipantAtHandlers) PutParticipantAtHandler(w http.ResponseWriter, r *http.Request) {
}
