package handlers

import (
	"encoding/json"
	"himakiwa/handlers/utils"
	"himakiwa/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

///////////////////////////////////
///////// chats handlers //////////
///////////////////////////////////

type ChatsHandlers struct {
	UseRepositoryServices services.UseRepositoryServices
}

func NewChatsHandlers(use services.UseRepositoryServices) func(http.ResponseWriter, *http.Request) {
	ch := &ChatsHandlers{use}
	return ch.ChatsHandlers
}

func (ch *ChatsHandlers) ChatsHandlers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ch.GetChatsHandler(w, r)
	default:
		http.Error(w, ErrUnExpcetedMethod, http.StatusBadRequest)
	}
}

// get latest chat in sessions

type GetChatsResp struct {
	SessionName string    `json:"sessionName"`
	SessionID   int       `json:"sesseionID"`
	UserID      int       `json:"userID"`
	ID          int       `json:"id"`
	Content     string    `json:"content"`
	CreateAt    time.Time `json:"createAt"`
	UpdateAt    time.Time `json:"updateAt"`
	Deleted     bool      `json:"deleted"`
}

func (ch *ChatsHandlers) GetChatsHandler(w http.ResponseWriter, r *http.Request) {
	// read context
	userID, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// session services
	sessionServices := ch.UseRepositoryServices(userID).SessionServices
	lastChats, err := sessionServices.GetLastChatInActiveSessions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//map
	resp := []GetChatsResp{}
	for _, chat := range lastChats {
		rs := GetChatsResp{
			chat.SessionName,
			chat.SessionID,
			chat.UserID,
			chat.ID,
			chat.Content,
			chat.CreateAt,
			chat.UpdateAt,
			chat.Deleted,
		}
		resp = append(resp, rs)
	}
	// resp
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

///////////////////////////////////
//////// chats at handlers ////////
///////////////////////////////////

type ChatsAtHandlers struct {
	UseRepositoryServices services.UseRepositoryServices
}

func NewChatsAtHandlers(use services.UseRepositoryServices) func(http.ResponseWriter, *http.Request) {
	cah := &ChatsAtHandlers{use}
	return cah.ChatsAtHandlers
}

func (cah *ChatsAtHandlers) ChatsAtHandlers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		cah.GetChatsAtHandler(w, r)
	case "POST":
		cah.PostChatsAtHandler(w, r)
	default:
		http.Error(w, ErrUnExpcetedMethod, http.StatusBadRequest)
	}
}

// get chats in 48hours

type GetChatsAtResp struct {
	ID        int       `json:"id"`
	SessionID int       `json:"sessionID"`
	UserID    int       `json:"userID"`
	Content   string    `json:"content"`
	CreateAt  time.Time `json:"createAt"`
	UpdateAt  time.Time `json:"updateAt"`
	Deleted   bool      `json:"deleted"`
}

func (cah *ChatsAtHandlers) GetChatsAtHandler(w http.ResponseWriter, r *http.Request) {
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
	sessionServices := cah.UseRepositoryServices(userID).SessionServices
	chats, err := sessionServices.GetChatsAtIn48Hours(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// map
	resp := []GetChatsAtResp{}
	for _, chat := range chats {
		rs := GetChatsAtResp{
			chat.ID,
			chat.SessionID,
			chat.UserID,
			chat.Content,
			chat.CreateAt,
			chat.UpdateAt,
			chat.Deleted,
		}
		resp = append(resp, rs)
	}
	// resp
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// send chat

type PostChatAtBody struct {
	Content string `validate:"required"`
}

func (cah *ChatsAtHandlers) PostChatsAtHandler(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &PostChatAtBody{}
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
	sessionServices := cah.UseRepositoryServices(userID).SessionServices

	err = sessionServices.SendChatAt(sessionID, b.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
