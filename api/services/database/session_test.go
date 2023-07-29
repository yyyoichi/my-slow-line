package database

import (
	"database/sql"
	"testing"
	"time"
)

func init() {
	Connect()
}

// //////////////////////////////////
// //////// test Session ////////////
// //////////////////////////////////

func TestSession(t *testing.T) {
	// Create a mock user
	usersR := &UserRepository{}
	mockUser := createMockUser()
	testUser, close := userMock(t, usersR, mockUser)
	userID := testUser.Id
	defer close()
	testSession(t, NewSessionRepositories(), userID)
}

func TestSessionMock(t *testing.T) {
	testSession(t, NewSessionRepositoriesMock(), 1)
}

type tQuerySession struct {
	Name   string
	Status TSessionStatus
}

func testSession(t *testing.T, repos *SessionRepositories, userID int) {
	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	sr := repos.SessionRepository

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

////////////////////////////////////
///// test Session Participant//////
////////////////////////////////////

func TestSessionParticipant(t *testing.T) {
	// Create a mock user
	usersR := &UserRepository{}
	mockUser := createMockUser()
	testUser, close := userMock(t, usersR, mockUser)
	userID := testUser.Id
	defer close()
	testSessionParticipant(t, NewSessionRepositories(), userID)
}
func TestSessionParticipantMock(t *testing.T) {
	testSessionParticipant(t, NewSessionRepositoriesMock(), 1)
}

func testSessionParticipant(t *testing.T, repos *SessionRepositories, userID int) {
	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		tx.Rollback()
	}()

	sr := repos.SessionRepository

	// create
	sessionID, err := sr.Create(tx, userID, "", "Test")
	if err != nil {
		t.Error(err)
	}
	defer sr.HardDelete(tx, sessionID)

	spr := repos.SessionParticipantRepository

	// create
	participantID, err := spr.Create(tx, sessionID, userID, userID, TInvitedParty)
	if err != nil {
		t.Error(err)
	}
	defer spr.HardDelete(tx, participantID)
	var participants []*TSessionParticipant

	// query and check
	participants, err = spr.QueryBySessionID(tx, sessionID)
	if err != nil {
		t.Error(err)
	}
	testQueryParticipant(t, tx, participants, TInvitedParty, userID)

	// update
	err = spr.UpdateStatus(tx, participantID, TJoinedParty)
	if err != nil {
		t.Error(err)
	}
	// query and check
	participants, err = spr.QueryBySessionID(tx, sessionID)
	if err != nil {
		t.Error(err)
	}
	testQueryParticipant(t, tx, participants, TJoinedParty, userID)

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

func testQueryParticipant(t *testing.T, tx *sql.Tx, participants []*TSessionParticipant, expStatus TParticipantStatus, expUserID int) {
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

////////////////////////////////////
///////// test Session Chat/////////
////////////////////////////////////

func TestSessionChat(t *testing.T) {
	// Create a mock user
	usersR := &UserRepository{}
	mockUser := createMockUser()
	testUser, close := userMock(t, usersR, mockUser)
	userID := testUser.Id
	defer close()
	testSessionChat(t, NewSessionRepositories(), userID)
}
func TestSessionChatMock(t *testing.T) {
	testSessionChat(t, NewSessionRepositoriesMock(), 1)
}

func testSessionChat(t *testing.T, repos *SessionRepositories, userID int) {
	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		tx.Rollback()
	}()

	sr := repos.SessionRepository

	// create
	sessionID, err := sr.Create(tx, userID, "", "Test")
	if err != nil {
		t.Error(err)
	}
	sessionID2, err := sr.Create(tx, userID, "", "Test2")
	if err != nil {
		t.Error(err)
	}
	defer sr.HardDelete(tx, sessionID)
	defer sr.HardDelete(tx, sessionID2)

	spr := repos.SessionParticipantRepository

	// create
	participantID, err := spr.Create(tx, sessionID, userID, userID, TJoinedParty)
	if err != nil {
		t.Error(err)
	}
	participantID2, err := spr.Create(tx, sessionID2, userID, userID, TJoinedParty)
	if err != nil {
		t.Error(err)
	}
	defer spr.HardDelete(tx, participantID)
	defer spr.HardDelete(tx, participantID2)

	scr := repos.SessionChatRepository

	// create
	chatID, err := scr.Create(tx, sessionID, userID, "Test Chat")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(1 * time.Second)
	chatID2, err := scr.Create(tx, sessionID2, userID, "Test Chat")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(1 * time.Second)
	chatID3, err := scr.Create(tx, sessionID2, userID, "Test Chat2")
	if err != nil {
		t.Error(err)
	}
	defer scr.HardDelete(tx, chatID)
	defer scr.HardDelete(tx, chatID2)
	defer scr.HardDelete(tx, chatID3)

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

	if len(chats) != 3 {
		t.Errorf("Expected len(chats) is 3, but got='%d'", len(chats))
	}

	lastChats, err := scr.QueryLastChatInActiveSessions(tx, userID)
	if err != nil {
		t.Error(err)
	}
	if len(lastChats) != 2 {
		t.Errorf("Expected len(lastChats) is 2, but got='%d'", len(lastChats))
	}
	if lastChats[0].ID != chatID3 {
		t.Errorf("Expected lastChats[0].ID is '%d', but got='%d'", chatID3, lastChats[0].ID)
	}
	if lastChats[1].ID != chatID {
		t.Errorf("Expected lastChats[1].ID is '%d', but got='%d'", chatID, lastChats[1].ID)
	}
}
