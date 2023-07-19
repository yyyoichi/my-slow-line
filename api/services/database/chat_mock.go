package database

import (
	"database/sql"
	"time"
)

type MockChatSessionRepository struct {
	SessionsByID map[int]*ChatSession
}

func NewMockChatSessionRepository() *MockChatSessionRepository {
	return &MockChatSessionRepository{
		SessionsByID: make(map[int]*ChatSession),
	}
}

func (m *MockChatSessionRepository) Query(tx *sql.Tx, sessionID int) (*ChatSession, error) {
	if session, found := m.SessionsByID[sessionID]; found {
		return session, nil
	}
	return nil, sql.ErrNoRows
}

func (m *MockChatSessionRepository) QueryByUserID(tx *sql.Tx, userID int) ([]ChatSession, error) {
	var sessions []ChatSession
	for _, session := range m.SessionsByID {
		if session.UserID == userID {
			sessions = append(sessions, *session)
		}
	}
	return sessions, nil
}

func (m *MockChatSessionRepository) UpdateName(tx *sql.Tx, id int, name string) error {
	if session, found := m.SessionsByID[id]; found {
		session.Name = name
		session.UpdateAt = time.Now()
		return nil
	}
	return sql.ErrNoRows
}

func (m *MockChatSessionRepository) Create(tx *sql.Tx, userID int, publicKey string, name string) (int, error) {
	id := len(m.SessionsByID) + 1
	session := &ChatSession{
		ID:        id,
		UserID:    userID,
		PublicKey: publicKey,
		Name:      name,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
		Deleted:   false,
	}
	m.SessionsByID[id] = session
	return id, nil
}

func (m *MockChatSessionRepository) Delete(tx *sql.Tx, sessionID int) error {
	if _, found := m.SessionsByID[sessionID]; found {
		delete(m.SessionsByID, sessionID)
		return nil
	}
	return sql.ErrNoRows
}

/////////////////////////////
//chat session participants//
/////////////////////////////

type MockChatSessionParticipantRepository struct {
	ParticipantsBySessionID map[int][]ChatSessionParticipant
	ParticipantsByUserID    map[int]*ChatSessionParticipant
}

func NewMockChatSessionParticipantRepository() *MockChatSessionParticipantRepository {
	return &MockChatSessionParticipantRepository{
		ParticipantsBySessionID: make(map[int][]ChatSessionParticipant),
		ParticipantsByUserID:    make(map[int]*ChatSessionParticipant),
	}
}

func (m *MockChatSessionParticipantRepository) QueryBySessionID(tx *sql.Tx, sessionID int) ([]ChatSessionParticipant, error) {
	participants, found := m.ParticipantsBySessionID[sessionID]
	if !found {
		return nil, sql.ErrNoRows
	}
	return participants, nil
}

func (m *MockChatSessionParticipantRepository) QueryBySessionAndUser(tx *sql.Tx, sessionID, userID int) (*ChatSessionParticipant, error) {
	participants, found := m.ParticipantsBySessionID[sessionID]
	if !found {
		return nil, sql.ErrNoRows
	}

	for _, participant := range participants {
		if participant.UserID == userID {
			return &participant, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (m *MockChatSessionParticipantRepository) QueryInvitedByUserID(tx *sql.Tx, userID int) ([]ChatSessionParticipant, error) {
	participant, found := m.ParticipantsByUserID[userID]
	if !found {
		return nil, sql.ErrNoRows
	}
	return []ChatSessionParticipant{*participant}, nil
}

func (m *MockChatSessionParticipantRepository) Create(tx *sql.Tx, sessionID int, userID, inviteUserID int, status ParticipantStatus) error {
	participants := m.ParticipantsBySessionID[sessionID]
	participant := &ChatSessionParticipant{
		ID:        len(participants) + 1,
		SessionID: sessionID,
		UserID:    userID,
		Status:    status,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
		Deleted:   false,
	}
	m.ParticipantsBySessionID[sessionID] = append(participants, *participant)
	m.ParticipantsByUserID[userID] = participant
	return nil
}

func (m *MockChatSessionParticipantRepository) Delete(tx *sql.Tx, sessionID int) error {
	delete(m.ParticipantsBySessionID, sessionID)
	return nil
}

func (m *MockChatSessionParticipantRepository) UpdateStatus(tx *sql.Tx, id int, status ParticipantStatus) error {
	for _, participants := range m.ParticipantsBySessionID {
		for i, participant := range participants {
			if participant.ID == id {
				participants[i].Status = status
				participants[i].UpdateAt = time.Now()
				return nil
			}
		}
	}
	return sql.ErrNoRows
}

/////////////////////////////
//         chats           //
/////////////////////////////

type MockChatRepository struct {
	ChatsBySessionID map[int][]Chat
	ChatsByUserID    map[int][]Chat
}

func NewMockChatRepository() *MockChatRepository {
	return &MockChatRepository{
		ChatsBySessionID: make(map[int][]Chat),
		ChatsByUserID:    make(map[int][]Chat),
	}
}

func (m *MockChatRepository) QueryBySessionID(tx *sql.Tx, sessionID int) ([]Chat, error) {
	chats, found := m.ChatsBySessionID[sessionID]
	if !found {
		return nil, sql.ErrNoRows
	}
	return chats, nil
}

func (m *MockChatRepository) QueryByUserIDAndTimeRange(tx *sql.Tx, userID int, endTime time.Time) ([]Chat, error) {
	// Implement the mock behavior here (return predefined data for testing)
	return nil, nil
}

func (m *MockChatRepository) Create(tx *sql.Tx, sessionID int, userID int, content string) error {
	// Implement the mock behavior here (store the data in the mock)
	chat := Chat{
		ID:        len(m.ChatsBySessionID[sessionID]) + 1,
		SessionID: sessionID,
		UserID:    userID,
		Content:   content,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
		Deleted:   false,
	}
	m.ChatsBySessionID[sessionID] = append(m.ChatsBySessionID[sessionID], chat)
	m.ChatsByUserID[userID] = append(m.ChatsByUserID[userID], chat)
	return nil
}

func (m *MockChatRepository) Delete(tx *sql.Tx, sessionID int) error {
	// Implement the mock behavior here (delete data from the mock)
	delete(m.ChatsBySessionID, sessionID)
	return nil
}

func (m *MockChatRepository) DeleteBySessionAndUser(tx *sql.Tx, sessionID int, userID int) error {
	// Implement the mock behavior here (delete data from the mock)
	chats := m.ChatsBySessionID[sessionID]
	var updatedChats []Chat
	for _, chat := range chats {
		if chat.UserID != userID {
			updatedChats = append(updatedChats, chat)
		}
	}
	m.ChatsBySessionID[sessionID] = updatedChats

	chats = m.ChatsByUserID[userID]
	updatedChats = nil
	for _, chat := range chats {
		if chat.SessionID != sessionID {
			updatedChats = append(updatedChats, chat)
		}
	}
	m.ChatsByUserID[userID] = updatedChats

	return nil
}

func (m *MockChatRepository) QueryByUserID(tx *sql.Tx, userID int) ([]Chat, error) {
	chats, found := m.ChatsByUserID[userID]
	if !found {
		return nil, sql.ErrNoRows
	}
	return chats, nil
}
