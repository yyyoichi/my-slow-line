package handlers

import (
	"himakiwa/services/database"
	"net/http"
	"time"
)

///////////////////////////////////
/////// sessions handlers /////////
///////////////////////////////////

type SessionsHandlers struct{}

func NewSessionsHandlers() func(http.ResponseWriter, *http.Request) {
	sh := &SessionsHandlers{}
	return sh.SessionsHandlers
}

func (sh *SessionsHandlers) SessionsHandlers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		sh.GetSessionsHandler(w, r)
	case "POST":
		sh.PostSessionsHandler(w, r)
	default:
		http.Error(w, ErrUnExpcetedMethod, http.StatusBadRequest)
	}
}

// get

type GetSessionsResp struct {
	ID            int                         `json:"id"`
	Name          string                      `json:"name"`
	PublicKey     string                      `json:"publicKey"`
	SessionStatus database.TSessionStatus     `json:"sessionStatus"`
	Status        database.TParticipantStatus `json:"status"`
	CreateAt      time.Time                   `json:"createAt"`
	UpdateAt      time.Time                   `json:"updateAt"`
	Deleted       bool                        `json:"deleted"`
}

func (sh *SessionsHandlers) GetSessionsHandler(w http.ResponseWriter, r *http.Request) {

}

// create

type PostSessionsBody struct {
	RecruitUUID string `validate:"required"`
	SessionName string `validate:"required"`
	PublicKey   string `validate:"required"`
}

func (sh *SessionsHandlers) PostSessionsHandler(w http.ResponseWriter, r *http.Request) {

	// get user by recruitUUID

	// create session and invite user
}

////////////////////////////////////////////////////
/// session at handlers (sessions/{sessionID}) /////
////////////////////////////////////////////////////

type SessionAtHandlers struct{}

func NewSessionAtHandlers() func(http.ResponseWriter, *http.Request) {
	sah := &SessionAtHandlers{}
	return sah.SessionAtHandlers
}

func (sh *SessionAtHandlers) SessionAtHandlers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		sh.GetSessionAtHandler(w, r)
	case "PUT":
		sh.PutSessionAtHandler(w, r)
	default:
		http.Error(w, ErrUnExpcetedMethod, http.StatusBadRequest)
	}
}

// get

type SessionAtResp struct {
	ID            int                         `json:"id"`
	Name          string                      `json:"name"`
	PublicKey     string                      `json:"publicKey"`
	SessionStatus database.TSessionStatus     `json:"sessionStatus"`
	Status        database.TParticipantStatus `json:"status"`
	Participants  TSeesionAtParticipantResp   `json:"participants"`
	CreateAt      time.Time                   `json:"createAt"`
	UpdateAt      time.Time                   `json:"updateAt"`
	Deleted       bool                        `json:"deleted"`
}
type TSeesionAtParticipantResp struct {
	ID       int                         `json:"id"`
	UserID   int                         `json:"userID"`
	Status   database.TParticipantStatus `json:"status"`
	CreateAt time.Time                   `json:"createAt"`
	UpdateAt time.Time                   `json:"updateAt"`
	Deleted  bool                        `json:"deleted"`
}

func (sh *SessionAtHandlers) GetSessionAtHandler(w http.ResponseWriter, r *http.Request) {
}

// edit

type PostSessionAtBody struct {
	SessionName string `validate:"required"`
}

func (sh *SessionAtHandlers) PutSessionAtHandler(w http.ResponseWriter, r *http.Request) {
}
