package database

import (
	"database/sql"
	"testing"
	"time"
)

func init() {
	Connect()
}

type tQuerySession struct {
	Name   string
	Status TSessionStatus
}

func querySessionByID(sessionID int) (*tQuerySession, error) {
	query := `SELECT name, status FROM chat_sessions WHERE id = ?`
	row := DB.QueryRow(query, sessionID)

	// result
	s := &tQuerySession{}
	err := row.Scan(&s.Name, &s.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found, return nil without an error.
			return nil, nil
		}
		return nil, err
	}

	return s, nil
}

func testQuerySession(t *testing.T, sessionID int, expName string, expStatus TSessionStatus) {
	session, err := querySessionByID(sessionID)
	if err != nil {
		t.Error(err)
	}
	if session.Name != expName {
		t.Errorf("Expected Name '%s', but got='%s'", expName, session.Name)
	}
	if session.Status != expStatus {
		t.Errorf("Expected Status '%s', but got='%s'", expStatus, session.Status)
	}
}
func TestSession(t *testing.T) {
	// Create a mock user
	usersR := &UserRepository{}
	mockUser := createMockUser()
	testUser, close := userMock(t, usersR, mockUser)
	userID := testUser.Id
	defer close()

	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	sr := SessionRepository{}

	// create
	sessionID, err := sr.Create(tx, userID, "", "Test")
	if err != nil {
		t.Error(err)
	}
	defer sr.HardDelete(tx, sessionID)

	// test
	testQuerySession(t, sessionID, "Test", TActiveSession)

	// update
	err = sr.UpdateName(tx, sessionID, "Update")
	if err != nil {
		t.Error(err)
	}
	testQuerySession(t, sessionID, "Update", TActiveSession)

	err = sr.UpdateStatus(tx, sessionID, TArchivedSession)
	if err != nil {
		t.Error(err)
	}
	testQuerySession(t, sessionID, "Update", TArchivedSession)
}

func testQueryParticipant(t *testing.T, tx *sql.Tx, sessionID int, expStatus TParticipantStatus, expUserID int) {
	spr := SessionParticipantRepository{}
	participants, err := spr.QueryBySessionID(tx, sessionID)
	if err != nil {
		t.Error(err)
	}
	if len(participants) != 1 {
		t.Errorf("Expected length of participants is 1, but got='%d'", len(participants))
	}
	participant := participants[0]

	if participant.Status != expStatus {
		t.Errorf("Expected Status '%s', but got='%s'", expStatus, participant.Status)
	}
	if participant.UserID != expUserID {
		t.Errorf("Expected UserID '%d', but got='%d'", expUserID, participant.UserID)
	}
}

func TestSessionParticipant(t *testing.T) {
	// Create a mock user
	usersR := &UserRepository{}
	mockUser := createMockUser()
	testUser, close := userMock(t, usersR, mockUser)
	userID := testUser.Id
	defer close()

	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		tx.Rollback()
	}()

	sr := SessionRepository{}

	// create
	sessionID, err := sr.Create(tx, userID, "", "Test")
	if err != nil {
		t.Error(err)
	}
	defer sr.HardDelete(tx, sessionID)

	spr := SessionParticipantRepository{}

	// create
	participantID, err := spr.Create(tx, sessionID, userID, userID, TInvitedParty)
	if err != nil {
		t.Error(err)
	}
	defer spr.HardDelete(tx, participantID)
	testQueryParticipant(t, tx, sessionID, TInvitedParty, userID)

	// update
	err = spr.UpdateStatus(tx, participantID, TJoinedParty)
	if err != nil {
		t.Error(err)
	}
	testQueryParticipant(t, tx, sessionID, TJoinedParty, userID)

	// joined test
	ok, err := sr.HasStatusAt(tx, sessionID, userID, []TParticipantStatus{TJoinedParty})
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("Expected HasStatusAt is true, but got=false")
	}

	ok, err = sr.HasStatusAt(tx, sessionID, userID, []TParticipantStatus{TInvitedParty})
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Error("Expected HasStatusAt is false, but got=true")
	}

	session := &TQuerySessions{}

	session, err = sr.QueryBySessionUserID(tx, sessionID, userID)
	if err != nil {
		t.Error(err)
	}
	if session.ID != sessionID {
		t.Errorf("Expected SessionID is '%d', but got='%d'", sessionID, session.ID)
	}
	if session.Status != TJoinedParty {
		t.Errorf("Expected Status is '%s', but got='%s'", TJoinedParty, session.Status)
	}

	sessions, err := sr.QueryByUserID(tx, userID, TQuerySessionsOptions{[]TParticipantStatus{TJoinedParty}, []TSessionStatus{TActiveSession}})
	if err != nil {
		t.Error(err)
	}
	if len(sessions) != 1 {
		t.Errorf("Expected len(sessions) is 1, but got='%d'", len(sessions))
	}
}

func TestSessionChat(t *testing.T) {
	// Create a mock user
	usersR := &UserRepository{}
	mockUser := createMockUser()
	testUser, close := userMock(t, usersR, mockUser)
	userID := testUser.Id
	defer close()

	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		tx.Rollback()
	}()

	sr := SessionRepository{}

	// create
	sessionID, err := sr.Create(tx, userID, "", "Test")
	if err != nil {
		t.Error(err)
	}
	defer sr.HardDelete(tx, sessionID)

	spr := SessionParticipantRepository{}

	// create
	participantID, err := spr.Create(tx, sessionID, userID, userID, TJoinedParty)
	if err != nil {
		t.Error(err)
	}
	defer spr.HardDelete(tx, participantID)

	scr := SessionChatRepository{}

	// create
	chatID, err := scr.Create(tx, sessionID, userID, "Test Chat")
	if err != nil {
		t.Error(err)
	}
	defer scr.HardDelete(tx, chatID)

	row := tx.QueryRow("SELECT content, create_at FROM chats where id = ?", chatID)
	var content string
	var at time.Time
	if err = row.Scan(&content, &at); err != nil {
		t.Error(err)
	}
	if content != "Test Chat" {
		t.Errorf("Expected Content is 'Test Chat', but got='%s'", content)
	}
	t.Log(at.String())

	chats, err := scr.QueryByUserIDInRange(tx, userID, struct {
		startDate time.Time
		endDate   time.Time
	}{
		startDate: time.Now().Add(-1 * time.Hour),
		endDate:   time.Now().Add(1 * time.Hour),
	})
	if err != nil {
		t.Error(err)
	}

	if len(chats) != 1 {
		t.Errorf("Expected len(chats) is 1, but got='%d'", len(chats))
	}
	if chats[0].Content != "Test Chat" {
		t.Errorf("Expected Content is 'Test Chat', but got='%s'", chats[0].Content)
	}
}
