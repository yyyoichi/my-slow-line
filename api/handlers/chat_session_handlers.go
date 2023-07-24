package handlers

import (
	"encoding/json"
	"fmt"
	"himakiwa/handlers/utils"
	chat_services "himakiwa/services/chats"
	"himakiwa/services/database"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ChatSessionsHandlers struct {
	GetChatService func(loginUserID int) *chat_services.ChatService
}

func NewChatSessionsHandlers() *ChatSessionsHandlers {
	return &ChatSessionsHandlers{func(loginUserID int) *chat_services.ChatService {
		return chat_services.NewChatService(loginUserID)
	}}
}

func (csh *ChatSessionsHandlers) ChatSessionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		csh.GetChatSessionsHandler(w, r)
	default:
		http.Error(w, ErrUnExpcetedMethod, http.StatusBadRequest)
	}
}

type GetSessionsHandlerRespChat struct {
	UserID  int       `json:"userID"`
	Content string    `json:"content"`
	At      time.Time `json:"at"`
}

type GetSessionsHandlerResp struct {
	ID       int                          `json:"id"`
	Name     string                       `json:"name"`
	Status   database.ParticipantStatus   `json:"status"`
	CreateAt time.Time                    `json:"createAt"`
	Chats    []GetSessionsHandlerRespChat `json:"chats"`
}

func (csh *ChatSessionsHandlers) GetChatSessionsHandler(w http.ResponseWriter, r *http.Request) {
	// read context
	userID, err := strconv.Atoi(utils.ReadUserContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get service
	sessionService := csh.GetChatService(userID)

	// go wait
	var wg sync.WaitGroup
	wg.Add(2)

	var sessions []GetSessionsHandlerResp
	var chats []database.Chat
	var err1, err2 error

	// get all sessions
	go func() {
		defer wg.Done()
		// get join
		sessionParticipants, err := sessionService.GetInvitedJoinedByUserID(userID)
		if err != nil {
			err1 = err
			return
		}

		var sessionWg sync.WaitGroup
		sessionWg.Add(len(sessionParticipants))

		// get session data
		for _, sp := range sessionParticipants {
			sessionID := sp.SessionID
			status := sp.Status
			createAt := sp.CreateAt
			go func() {
				defer sessionWg.Done()
				ses, err := sessionService.GetChatSessionByID(sessionID)
				if err != nil {
					err1 = err
					return
				}
				sessions = append(sessions, GetSessionsHandlerResp{
					ID: sessionID, Name: ses.Name, Status: status, CreateAt: createAt,
				})
			}()
		}
		sessionWg.Wait()
	}()

	// get chats in timerange
	go func() {
		defer wg.Done()
		// expected sync user schedule
		endTime := time.Now()
		chats, err2 = sessionService.GetChatMessagesByTimeRange(endTime)
	}()

	wg.Wait()

	if err1 != nil || err2 != nil {
		err := fmt.Sprintf("%s,%s", err1.Error(), err2.Error())
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	// map chat
	for i, ses := range sessions {
		for _, chat := range chats {
			if chat.SessionID != ses.ID {
				continue
			}
			// append chats
			sessions[i].Chats = append(ses.Chats, GetSessionsHandlerRespChat{
				UserID: chat.UserID, Content: chat.Content, At: chat.CreateAt,
			})
		}
	}

	//resp
	err = json.NewEncoder(w).Encode(sessions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
