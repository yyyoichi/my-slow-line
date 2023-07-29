package handlers

import (
	"himakiwa/services/database"
	"net/http"
)

///////////////////////////////////
//// participant at handlers //////
///////////////////////////////////

type ParticipantsAtHandlers struct{}

func NewParticipantsAtHandlers() func(http.ResponseWriter, *http.Request) {
	pah := &ParticipantsAtHandlers{}
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

type PostParticipantsAtHandlerBody struct {
	UserID int `validate:"required"`
}

// TODO: write
func (pah *ParticipantsAtHandlers) PostParticipantsAtHandler(w http.ResponseWriter, r *http.Request) {
}

// update status

type PoutParticipantAtHandlerBody struct {
	UserID int                         `validate:"required"`
	Status database.TParticipantStatus `validate:"required"`
}

func (pah *ParticipantsAtHandlers) PutParticipantsAtHandler(w http.ResponseWriter, r *http.Request) {
}
