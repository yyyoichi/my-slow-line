package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"himakiwa/handlers/utils"
	"himakiwa/services"
	"himakiwa/services/database"
	"himakiwa/services/webpush"
	"net/http"
	"strconv"
)

var (
	ErrCannotAccessSession      = errors.New("cannot access session")
	ErrInviteeCannotBeInvted    = errors.New("invitee cannot be invited")
	ErrCannotPermissionToInvite = errors.New("cannnot invite the user not to have permission")
)

type SessionKeyHandlers struct {
	services.UseRepositoryServices
	webpush.UserWebpushServices
}

func NewSessionKeyHandlers(useRepository services.UseRepositoryServices, useWebpush webpush.UserWebpushServices) func(http.ResponseWriter, *http.Request) {
	kh := &SessionKeyHandlers{useRepository, useWebpush}
	return kh.PostSessionKey
}

type PostSessionKeyBody struct {
	SessionID int    `validate:"required"`
	InviteeID int    `validate:"required"`
	Key       string `validate:"required"`
}

// send sessionkey to invitee with webpush
func (kh *SessionKeyHandlers) PostSessionKey(w http.ResponseWriter, r *http.Request) {
	// parse body
	b := &PostSessionKeyBody{}
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

	reposiotryServices := kh.UseRepositoryServices(userID)

	// the offer user joins the session and the invitee user invite it
	session, participants, err := reposiotryServices.SessionServices.GetSessionAt(b.SessionID)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session == nil {
		http.Error(w, ErrCannotAccessSession.Error(), http.StatusBadRequest)
		return
	}
	inviteeIsInvite := false
	for _, party := range participants {
		if party.UserID == b.InviteeID && party.Status == database.TInvitedParty {
			inviteeIsInvite = true
			break
		}
	}

	offerIsJoined := session.Status == database.TJoinedParty
	if !offerIsJoined || !inviteeIsInvite {
		fmt.Println(ErrCannotPermissionToInvite)
		http.Error(w, ErrCannotPermissionToInvite.Error(), http.StatusBadRequest)
		return
	}

	// get invitee webpush  subscription data
	subscriptions, err := reposiotryServices.UserServices.GetWebpushSubscriptions(b.InviteeID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(subscriptions) == 0 {
		fmt.Println(ErrInviteeCannotBeInvted)
		http.Error(w, ErrInviteeCannotBeInvted.Error(), http.StatusBadRequest)
		return
	}

	// send key with webpush

	// create args to send key
	exchSessionKeyArgs := webpush.TExchSessionKeyArgs{
		SessionID:         session.ID,
		SessionName:       session.Name,
		NumOfParticipants: 0,  // require
		OfferUserName:     "", // require
		Key:               b.Key,
	}
	// numOfParticipants
	for _, party := range participants {
		if party.Status == database.TJoinedParty {
			exchSessionKeyArgs.NumOfParticipants += 1
		}
	}
	// offerUserName
	user, err := reposiotryServices.UserServices.GetUser(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	exchSessionKeyArgs.OfferUserName = user.Name

	// send webpush
	ss := subscriptions[0]
	if err := kh.UserWebpushServices(ss.Endpoint, ss.Auth, ss.P256dh).SendExchSessionKeyMessage(exchSessionKeyArgs); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
