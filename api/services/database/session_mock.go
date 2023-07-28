package database

import (
	"database/sql"
	"time"
)

type SessionDataMock struct {
	sessionByID            map[int]*tSessionTable
	participantBySessionID map[int][]*TSessionParticipant
	chatByUserID           map[int][]*TQuerySessionChat
}

func (dm *SessionDataMock) genSessions() <-chan *tSessionTable {
	ch := make(chan *tSessionTable)
	go func() {
		for _, session := range dm.sessionByID {
			ch <- session
		}
		close(ch)
	}()

	return ch
}

func (dm *SessionDataMock) genParticipants() <-chan *TSessionParticipant {
	ch := make(chan *TSessionParticipant)
	go func() {
		for _, participants := range dm.participantBySessionID {
			for _, party := range participants {
				ch <- party
			}
		}
		close(ch)
	}()

	return ch
}

func (dm *SessionDataMock) genChats() <-chan *TQuerySessionChat {
	ch := make(chan *TQuerySessionChat)
	go func() {
		for _, chats := range dm.chatByUserID {
			for _, chat := range chats {
				ch <- chat
			}
		}
		close(ch)
	}()

	return ch
}

// delete
func (dm *SessionDataMock) hardParticipantsDeleteBySessionID(tx *sql.Tx, sessionID int) error {
	delete(dm.participantBySessionID, sessionID)
	return nil
}

type tSessionTable struct {
	id         int
	user_id    int
	public_key string
	name       string
	status     TSessionStatus
	create_at  time.Time
	update_at  time.Time
	deleted    bool
}

func NewSessionRepositoriesMock() *SessionRepositories {
	sessionByID := make(map[int]*tSessionTable)
	participantBySessionID := make(map[int][]*TSessionParticipant)
	chatByUserID := make(map[int][]*TQuerySessionChat)
	mock := &SessionDataMock{sessionByID, participantBySessionID, chatByUserID}
	return &SessionRepositories{
		&SessionRepositoryMock{mock},
		&SessionParticipantRepositoryMock{mock},
		&SessionChatRepositoryMock{mock},
	}
}

type SessionRepositoryMock struct {
	mock *SessionDataMock
}

func (sr *SessionRepositoryMock) QueryByUserID(tx *sql.Tx, userID int, options TQuerySessionsOptions) ([]*TQuerySessions, error) {
	var tmpResults []struct {
		id     int
		status TParticipantStatus
	}
	for party := range sr.mock.genParticipants() {
		// eq user
		if party.UserID != userID {
			continue
		}
		// has status
		hasStatus := false
		for _, status := range options.InPartyStatus {
			if status == party.Status {
				hasStatus = true
				break
			}
		}
		if !hasStatus {
			continue
		}

		// add results
		session := struct {
			id     int
			status TParticipantStatus
		}{
			party.SessionID,
			party.Status,
		}
		tmpResults = append(tmpResults, session)
	}

	var results []*TQuerySessions
	// map session
	for _, rlt := range tmpResults {
		for session := range sr.mock.genSessions() {
			if rlt.id != session.id {
				continue
			}
			hasStatus := false
			for _, status := range options.InSessionStatus {
				if status == session.status {
					hasStatus = true
					break
				}
			}
			if !hasStatus {
				continue
			}
			result := &TQuerySessions{
				session.id,
				session.name,
				session.public_key,
				session.status,
				rlt.status,
				session.create_at,
				session.update_at,
				session.deleted,
			}
			results = append(results, result)
		}
	}
	return results, nil
}

// query a session
func (sr *SessionRepositoryMock) QueryBySessionUserID(tx *sql.Tx, sessionID, userID int) (*TQuerySessions, error) {
	results, err := sr.QueryByUserID(tx, userID, TQuerySessionsOptions{
		[]TParticipantStatus{TInvitedParty, TJoinedParty, TRejectedParty},
		[]TSessionStatus{TActiveSession, TArchivedSession, TBreakupSession},
	})
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		if r.ID == sessionID {
			return r, nil
		}
	}
	return nil, sql.ErrNoRows
}

// has status at session for user
func (sr *SessionRepositoryMock) HasStatusAt(tx *sql.Tx, sessionID, userID int, inStatus []TParticipantStatus) (bool, error) {
	result, err := sr.QueryBySessionUserID(tx, sessionID, userID)
	if err != nil {
		return false, err
	}
	hasStatus := false
	for _, s := range inStatus {
		if s == result.Status {
			hasStatus = true
			break
		}
	}
	return hasStatus, nil
}

// create
func (sr *SessionRepositoryMock) Create(tx *sql.Tx, userID int, publicKey string, name string) (int, error) {
	id := len(sr.mock.sessionByID) + 1
	session := &tSessionTable{id, userID, publicKey, name, TActiveSession, time.Now(), time.Now(), false}
	sr.mock.sessionByID[id] = session
	return id, nil
}

// update
func (sr *SessionRepositoryMock) UpdateName(tx *sql.Tx, id int, name string) error {
	if session, found := sr.mock.sessionByID[id]; found {
		session.name = name
		session.update_at = time.Now()
		return nil
	}
	return sql.ErrNoRows
}

func (sr *SessionRepositoryMock) UpdateStatus(tx *sql.Tx, id int, status TSessionStatus) error {
	if session, found := sr.mock.sessionByID[id]; found {
		session.status = status
		session.update_at = time.Now()
		return nil
	}
	return sql.ErrNoRows
}

// delete row
func (sr *SessionRepositoryMock) HardDelete(tx *sql.Tx, sessionID int) error {
	if _, found := sr.mock.sessionByID[sessionID]; found {
		delete(sr.mock.sessionByID, sessionID)
		return nil
	}
	return sql.ErrNoRows
}

// up deleted flag
func (sr *SessionRepositoryMock) SoftDelete(tx *sql.Tx, sessionID int) error {
	if session, found := sr.mock.sessionByID[sessionID]; found {
		session.deleted = true
		session.update_at = time.Now()
		return nil
	}
	return sql.ErrNoRows
}

// delete hard all session contents
func (sr *SessionRepositoryMock) HardDeleteAll(tx *sql.Tx, sessionID int) error {
	sr.HardDelete(tx, sessionID)
	sr.mock.hardParticipantsDeleteBySessionID(tx, sessionID)
	return nil
}

type SessionParticipantRepositoryMock struct {
	mock *SessionDataMock
}

func (spr *SessionParticipantRepositoryMock) QueryBySessionID(tx *sql.Tx, sessionID int) ([]*TSessionParticipant, error) {
	// results
	var results []*TSessionParticipant
	for party := range spr.mock.genParticipants() {
		if sessionID == party.SessionID {
			results = append(results, party)
		}
	}

	return results, nil
}

// [invite_user_id] is ID of the user exec to the session
func (spr *SessionParticipantRepositoryMock) Create(tx *sql.Tx, sessionID, userID, inviteUserID int, status TParticipantStatus) (int, error) {
	participants := spr.mock.participantBySessionID[sessionID]
	id := len(participants) + 1
	participant := &TSessionParticipant{
		ID:        id,
		SessionID: sessionID,
		UserID:    userID,
		Status:    status,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
		Deleted:   false,
	}
	spr.mock.participantBySessionID[sessionID] = append(participants, participant)
	return id, nil
}

// update
func (spr *SessionParticipantRepositoryMock) UpdateStatus(tx *sql.Tx, id int, status TParticipantStatus) error {
	for party := range spr.mock.genParticipants() {
		if id == party.ID {
			party.Status = status
			return nil
		}
	}

	return sql.ErrNoRows
}

// delete
func (spr *SessionParticipantRepositoryMock) HardDelete(tx *sql.Tx, participantID int) error {
	for sessionID, participants := range spr.mock.participantBySessionID {
		p := []*TSessionParticipant{}
		for _, party := range participants {
			if party.ID != participantID {
				continue
			}
			p = append(p, party)
		}
		spr.mock.participantBySessionID[sessionID] = p
	}

	return nil
}

type SessionChatRepositoryMock struct {
	mock *SessionDataMock
}

func (cr *SessionChatRepositoryMock) QueryByUserIDInRange(tx *sql.Tx, userID int, inRange struct{ startDate, endDate time.Time }) ([]*TQuerySessionChat, error) {
	targetSessionID := make(map[int]bool)
	for party := range cr.mock.genParticipants() {
		if party.UserID == userID && party.Status == TJoinedParty {
			targetSessionID[party.SessionID] = true
		}
	}
	var results []*TQuerySessionChat
	for chat := range cr.mock.genChats() {
		if _, found := targetSessionID[chat.SessionID]; found {
			results = append(results, chat)
		}
	}
	return results, nil
}

// create
func (scr *SessionChatRepositoryMock) Create(tx *sql.Tx, sessionID, userID int, content string) (int, error) {
	id := len(scr.mock.chatByUserID[userID]) + 1
	chat := &TQuerySessionChat{
		ID:        id,
		SessionID: sessionID,
		UserID:    userID,
		Content:   content,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
		Deleted:   false,
	}
	scr.mock.chatByUserID[userID] = append(scr.mock.chatByUserID[userID], chat)
	return id, nil
}

// delete
func (cr *SessionChatRepositoryMock) HardDelete(tx *sql.Tx, chatID int) error {
	for userID, chats := range cr.mock.chatByUserID {
		c := []*TQuerySessionChat{}
		for _, chat := range chats {
			if chat.ID != chatID {
				continue
			}
			c = append(c, chat)
		}
		cr.mock.chatByUserID[userID] = c
	}
	return nil
}
