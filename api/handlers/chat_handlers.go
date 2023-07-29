package handlers

import (
	"net/http"
	"time"
)

///////////////////////////////////
///////// chats handlers //////////
///////////////////////////////////

type ChatsHandlers struct{}

func NewChatsHandlers() func(http.ResponseWriter, *http.Request) {
	ch := &ChatsHandlers{}
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
	UserName    string    `json:"userName"`
	ID          int       `json:"id"`
	Content     string    `json:"content"`
	CreateAt    time.Time `json:"createAt"`
	UpdateAt    time.Time `json:"updateAt"`
	Deleted     bool      `json:"deleted"`
}

func (ch *ChatsHandlers) GetChatsHandler(w http.ResponseWriter, r *http.Request) {
}

///////////////////////////////////
//////// chats at handlers ////////
///////////////////////////////////

type ChatsAtHandlers struct{}

func NewChatsAtHandlers() func(http.ResponseWriter, *http.Request) {
	cah := &ChatsAtHandlers{}
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

type GetChatsAtHandlerResp struct {
	ID        int       `json:"id"`
	SessionID int       `json:"sessionID"`
	UserID    int       `json:"userID"`
	Content   string    `json:"content"`
	CreateAt  time.Time `json:"createAt"`
	UpdateAt  time.Time `json:"updateAt"`
	Deleted   bool      `json:"deleted"`
}

func (cah *ChatsAtHandlers) GetChatsAtHandler(w http.ResponseWriter, r *http.Request) {
}

// send chat

type PostChatAtHandlerBody struct {
	UserID  int    `validate:"required"`
	Content string `validate:"required"`
}

func (cah *ChatsAtHandlers) PostChatsAtHandler(w http.ResponseWriter, r *http.Request) {
}
