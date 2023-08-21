package database

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func init() {
	Connect()
}

// //////////////////////////////////
// //////// test Session ////////////
// //////////////////////////////////

type TestingSession struct {
	Session      *TQuerySession
	ExpSession   *TQuerySession
	tx           *sql.Tx
	repositories *SessionRepositories
}

func (ts *TestingSession) Delete() error {
	var err error
	err = ts.repositories.SessionRepository.HardDelete(ts.tx, ts.Session.ID)
	if err != nil {
		return err
	}
	err = ts.repositories.SessionParticipantRepository.HardDelete(ts.tx, ts.Session.ID)
	return err
}

var TestSessionCount = 0

func CreateTestingSession(tx *sql.Tx, srs *SessionRepositories, userID int) (*TestingSession, error) {
	TestSessionCount += 1
	sessionName := fmt.Sprintf("Test Session %d", TestSessionCount)
	sessionID, err := srs.SessionRepository.Create(tx, userID, "key", sessionName)
	if err != nil {
		return nil, err
	}
	_, err = srs.SessionParticipantRepository.Create(tx, sessionID, userID, userID, TInvitedParty)
	if err != nil {
		return nil, err
	}

	session, err := srs.SessionRepository.QueryBySessionUserID(tx, sessionID, userID)
	if err != nil {
		return nil, err
	}
	expSession := &TQuerySession{
		sessionID,
		sessionName,
		"key",
		TActiveSession,
		TInvitedParty,
		time.Now(),
		time.Now(),
		false,
	}

	return &TestingSession{session, expSession, tx, srs}, nil
}

func TestSession(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	usr := NewUserRepositories()
	user, err := CreateTestingUser(tx, usr)
	if err != nil {
		t.Error(err)
	}
	testSession(t, tx, NewSessionRepositories(), user.User.ID)
}

func TestSessionMock(t *testing.T) {
	testSession(t, &sql.Tx{}, NewSessionRepositoriesMock(), 1)
}

func testSession(t *testing.T, tx *sql.Tx, repos *SessionRepositories, userID int) {
	// create and test
	testingSession, err := CreateTestingSession(tx, repos, userID)
	if err != nil {
		t.Error(err)
	}
	defer testingSession.Delete()
	session := testingSession.Session
	expSession := testingSession.ExpSession
	sessionIsEqual(t, session, expSession)

	sr := repos.SessionRepository
	// status at
	if has, err := sr.HasStatusAt(tx, session.ID, userID, []TParticipantStatus{TInvitedParty}); err != nil {
		t.Error(err)
	} else if !has {
		t.Errorf("Expected has status is 'true', but got='false'")
	}
	if has, err := sr.HasStatusAt(tx, session.ID, userID, []TParticipantStatus{TJoinedParty, TRejectedParty}); err != nil {
		t.Error(err)
	} else if has {
		t.Errorf("Expected has status is 'false', but got='true'")
	}

	// update
	if err = sr.UpdateName(tx, session.ID, "Rename"); err != nil {
		t.Error(err)
	}
	expSession.Name = "Rename"
	session, err = sr.QueryBySessionUserID(tx, session.ID, userID)
	if err != nil {
		t.Error(err)
	}
	sessionIsEqual(t, session, expSession)

	// update session status
	err = sr.UpdateStatus(tx, session.ID, TArchivedSession)
	if err != nil {
		t.Error(err)
	}
	expSession.SessionStatus = TArchivedSession
	session, err = sr.QueryBySessionUserID(tx, session.ID, userID)
	if err != nil {
		t.Error(err)
	}
	sessionIsEqual(t, session, expSession)

	// soft delete session
	if err = sr.SoftDelete(tx, session.ID); err != nil {
		t.Error(err)
	}
	expSession.Deleted = true
	session, err = sr.QueryBySessionUserID(tx, session.ID, userID)
	if err != nil {
		t.Error(err)
	}
	sessionIsEqual(t, session, expSession)

	// test query sessions
	options := TQuerySessionsOptions{
		[]TParticipantStatus{TInvitedParty, TJoinedParty, TRejectedParty},
		[]TSessionStatus{TActiveSession, TArchivedSession, TBreakupSession},
	}
	if sessions, err := sr.QueryByUserID(tx, userID, options); err != nil {
		t.Error(err)
	} else if len(sessions) != 1 {
		t.Errorf("Expected len(sessions) is 1, but got=%d", len(sessions))
	}

	// add session
	if _, err = CreateTestingSession(tx, repos, userID); err != nil {
		t.Error(err)
	}

	if sessions, err := sr.QueryByUserID(tx, userID, options); err != nil {
		t.Error(err)
	} else if len(sessions) != 2 {
		t.Errorf("Expected len(sessions) is 2, but got=%d", len(sessions))
	}

	// filter
	options.InSessionStatus = []TSessionStatus{TArchivedSession}
	if sessions, err := sr.QueryByUserID(tx, userID, options); err != nil {
		t.Error(err)
	} else if len(sessions) != 1 {
		t.Errorf("Expected len(sessions) is 1, but got=%d", len(sessions))
	}

	fmt.Print("Done")
}

func sessionIsEqual(t *testing.T, act, exp *TQuerySession) {
	if act.Name != exp.Name {
		t.Errorf("Expected Name '%s', but got='%s'", exp.Name, act.Name)
	}
	if act.Status != exp.Status {
		t.Errorf("Expected Status '%s', but got='%s'", exp.Status, act.Status)
	}
	if act.SessionStatus != exp.SessionStatus {
		t.Errorf("Expected SessionStatus '%s', but got='%s'", exp.SessionStatus, act.SessionStatus)
	}
	if act.PublicKey != exp.PublicKey {
		t.Errorf("Expected PublicKey '%s', but got='%s'", exp.PublicKey, act.PublicKey)
	}
	if act.ID != exp.ID {
		t.Errorf("Expected ID '%d', but got='%d'", exp.ID, act.ID)
	}
	if act.ID != exp.ID {
		t.Errorf("Expected ID '%d', but got='%d'", exp.ID, act.ID)
	}
	if act.Deleted != exp.Deleted {
		t.Errorf("Expected Deleted '%v', but got='%v'", exp.Deleted, act.Deleted)
	}
	if act.CreateAt.IsZero() {
		t.Errorf("Expected CreateAt is not zero but got='zero'")
	}
	if act.UpdateAt.IsZero() {
		t.Errorf("Expected UpdateAt is not zero but got='zero'")
	}
}

////////////////////////////////////
///// test Session Participant//////
////////////////////////////////////

func TestSessionParticipant(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	usr := NewUserRepositories()
	user, err := CreateTestingUser(tx, usr)
	if err != nil {
		t.Error(err)
	}
	testSessionParticipant(t, tx, NewSessionRepositories(), user.User.ID)
}
func TestSessionParticipantMock(t *testing.T) {
	testSessionParticipant(t, &sql.Tx{}, NewSessionRepositoriesMock(), 1)
}

func testSessionParticipant(t *testing.T, tx *sql.Tx, repos *SessionRepositories, userID int) {
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
	err = spr.UpdateStatusBySessionUserID(tx, sessionID, userID, TJoinedParty)
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

	session := &TQuerySession{}

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
	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	usr := NewUserRepositories()
	user, err := CreateTestingUser(tx, usr)
	if err != nil {
		t.Error(err)
	}
	testSessionChat(t, tx, NewSessionRepositories(), user.User.ID)
}
func TestSessionChatMock(t *testing.T) {
	testSessionChat(t, &sql.Tx{}, NewSessionRepositoriesMock(), 1)
}

func testSessionChat(t *testing.T, tx *sql.Tx, repos *SessionRepositories, userID int) {
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

	inRange := TQuerySessionChatInRange{
		StartDate: time.Now().Add(-1 * time.Hour),
		EndDate:   time.Now().Add(1 * time.Hour),
	}
	chats, err := scr.QueryByUserIDInRange(tx, userID, inRange)
	if err != nil {
		t.Error(err)
	}

	if len(chats) != 3 {
		t.Errorf("Expected len(chats) is 3, but got='%d'", len(chats))
	}

	chats, err = scr.QueryBySessionIDInRange(tx, sessionID2, inRange)
	if err != nil {
		t.Error(err)
	}
	if len(chats) != 2 {
		t.Errorf("Expected len(chats) is 2, but got='%d'", len(chats))
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
