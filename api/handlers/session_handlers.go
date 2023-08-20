package handlers

import (
	"encoding/json"
	"himakiwa/handlers/utils"
	"himakiwa/services"
	"himakiwa/services/database"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

///////////////////////////////////
/////// sessions handlers /////////
///////////////////////////////////

type SessionsHandlers struct {
	UseRepositoryServices services.UseRepositoryServices
}

func NewSessionsHandlers(use services.UseRepositoryServices) func(http.ResponseWriter, *http.Request) {
	sh := &SessionsHandlers{use}
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
	// read context
	userID, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// session services
	sessionServices := sh.UseRepositoryServices(userID).SessionServices
	sessions, err := sessionServices.GetActiveOrArchivedSessions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp = []GetSessionsResp{}
	for _, session := range sessions {
		rs := GetSessionsResp{
			session.ID,
			session.Name,
			session.PublicKey,
			session.SessionStatus,
			session.Status,
			session.CreateAt,
			session.UpdateAt,
			session.Deleted,
		}
		resp = append(resp, rs)
	}
	//resp
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// create

type PostSessionsBody struct {
	RecruitUUID string `validate:"required"`
	SessionName string `validate:"required"`
	PublicKey   string `validate:"required"`
}

func (sh *SessionsHandlers) PostSessionsHandler(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &PostSessionsBody{}
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
	repositoryServices := sh.UseRepositoryServices(userID)
	// user services
	userSerivices := repositoryServices.UserServices
	// session services
	sessionServices := repositoryServices.SessionServices

	// get recruit user data
	user, err := userSerivices.GetUserByRecruitUUID(b.RecruitUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sessionServices.CreateSession(b.PublicKey, b.SessionName, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

////////////////////////////////////////////////////
/// session at handlers (sessions/{sessionID}) /////
////////////////////////////////////////////////////

type SessionAtHandlers struct {
	UseRepositoryServices services.UseRepositoryServices
}

func NewSessionAtHandlers(use services.UseRepositoryServices) func(http.ResponseWriter, *http.Request) {
	sah := &SessionAtHandlers{use}
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
	Participants  []TSeesionAtParticipantResp `json:"participants"`
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

	// read context
	userID, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// session services
	sessionServices := sh.UseRepositoryServices(userID).SessionServices

	// get session
	session, participants, err := sessionServices.GetSessionAt(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// map participants
	respParticipants := []TSeesionAtParticipantResp{}
	for _, party := range participants {
		rspp := TSeesionAtParticipantResp{
			party.ID,
			party.UserID,
			party.Status,
			party.CreateAt,
			party.UpdateAt,
			party.Deleted,
		}
		respParticipants = append(respParticipants, rspp)
	}

	// map resp
	resp := SessionAtResp{
		session.ID,
		session.Name,
		session.PublicKey,
		session.SessionStatus,
		session.Status,
		respParticipants,
		session.CreateAt,
		session.UpdateAt,
		session.Deleted,
	}

	//resp
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// edit

type PutSessionAtBody struct {
	SessionName string `validate:"required"`
}

func (sh *SessionAtHandlers) PutSessionAtHandler(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &PutSessionAtBody{}
	if err := utils.DecodeBody(r, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	// read context
	userID, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// session services
	sessionServices := sh.UseRepositoryServices(userID).SessionServices

	err = sessionServices.UpdateSessionNameAt(sessionID, b.SessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
