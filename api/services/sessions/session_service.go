package sessions

import (
	"database/sql"
	"errors"
	"himakiwa/services/database"
	"time"
)

var (
	ErrCannotAccessSession = errors.New("cannot access session")
	ErrCannotUpdateStatus  = errors.New("cannot access participant session")
)

type UseSessionServicesFunc func(loginID int) *SessionServices

type SessionServices struct {
	useTransaction database.TUseTransaction
	repositories   *database.SessionRepositories
	loginUserID    int
}

func NewSessionServices() UseSessionServicesFunc {
	ss := &SessionServices{useTransaction: database.UseTransaction, repositories: database.NewSessionRepositories()}
	return func(loginID int) *SessionServices {
		ss.loginUserID = loginID
		return ss
	}
}

// get active and archived sessions
func (ss *SessionServices) GetActiveOrArchivedSessions() ([]*database.TQuerySession, error) {
	options := database.TQuerySessionsOptions{
		InPartyStatus:   []database.TParticipantStatus{database.TJoinedParty, database.TInvitedParty},
		InSessionStatus: []database.TSessionStatus{database.TActiveSession, database.TArchivedSession},
	}
	var sessions []*database.TQuerySession
	err := ss.useTransaction(func(tx *sql.Tx) error {
		var err error
		sessions, err = ss.repositories.SessionRepository.QueryByUserID(tx, ss.loginUserID, options)
		return err
	})
	if err != nil {
		return nil, err
	}

	// to secure
	secureSessions := []*database.TQuerySession{}
	for _, session := range sessions {
		if session, err := ss.toSecureSession(session); err != nil {
			return nil, err
		} else if session != nil {
			secureSessions = append(secureSessions, session)
		}
	}
	return secureSessions, nil
}

// create session and loginUser invite userID, return sessionID
func (ss *SessionServices) CreateSession(publicKey, name string, userID int) (int, error) {
	var sessionID int
	err := ss.useTransaction(func(tx *sql.Tx) error {
		var err error
		sessionID, err = ss.repositories.SessionRepository.Create(tx, ss.loginUserID, publicKey, name)
		if err != nil {
			return err
		}
		_, err = ss.repositories.SessionParticipantRepository.Create(tx, sessionID, ss.loginUserID, ss.loginUserID, database.TJoinedParty)
		if err != nil {
			return err
		}
		_, err = ss.repositories.SessionParticipantRepository.Create(tx, sessionID, userID, ss.loginUserID, database.TInvitedParty)
		return err
	})
	if err != nil {
		return 0, err
	}
	return sessionID, nil
}

// get session and participants
func (ss *SessionServices) GetSessionAt(sessionID int) (*database.TQuerySession, []*database.TSessionParticipant, error) {
	var session *database.TQuerySession
	var participants []*database.TSessionParticipant
	err := ss.useTransaction(func(tx *sql.Tx) error {
		// get session
		var err error
		session, err = ss.repositories.SessionRepository.QueryBySessionUserID(tx, sessionID, ss.loginUserID)
		if err != nil {
			return err
		}
		// get participants
		participants, err = ss.repositories.SessionParticipantRepository.QueryBySessionID(tx, sessionID)
		return err
	})

	if err != nil {
		return nil, nil, err
	}

	// to secure
	if session, err = ss.toSecureSession(session); err != nil {
		return nil, nil, err
	}
	if participants, err = ss.toSecureParticipants(session.Status, participants); err != nil {
		return nil, nil, err
	}
	return session, participants, nil
}

func (ss *SessionServices) UpdateSessionNameAt(sessionID int, name string) error {
	return ss.useTransaction(func(tx *sql.Tx) error {
		ok, err := ss.repositories.SessionRepository.HasStatusAt(tx, sessionID, ss.loginUserID, []database.TParticipantStatus{database.TJoinedParty})
		if err != nil {
			return err
		}
		if !ok {
			return ErrCannotAccessSession
		}
		err = ss.repositories.SessionRepository.UpdateName(tx, sessionID, name)
		return err
	})
}

/*
secure session is
1. public key of session to that login user is not join should be hided
2. if login user is rejected, cannot access
*/
func (ss *SessionServices) toSecureSession(session *database.TQuerySession) (*database.TQuerySession, error) {
	if session.Status == database.TRejectedParty {
		return nil, ErrCannotAccessSession
	}
	if session.Status == database.TInvitedParty {
		session.PublicKey = ""
		return session, nil
	}
	// join
	return session, nil
}

/*
secure sesseion participants is
1. participants of session is hided from invited, rejected user
2. if login user is rejected, cannot access
*/
func (ss *SessionServices) toSecureParticipants(loginUserStatus database.TParticipantStatus, participants []*database.TSessionParticipant) ([]*database.TSessionParticipant, error) {
	if loginUserStatus == database.TRejectedParty {
		return nil, ErrCannotAccessSession
	}

	if loginUserStatus == database.TJoinedParty {
		return participants, nil
	}

	// invite
	var joinedAndMeParticipants []*database.TSessionParticipant
	for _, party := range participants {
		if party.Status == database.TJoinedParty || party.UserID == ss.loginUserID {
			joinedAndMeParticipants = append(joinedAndMeParticipants, party)
		}
	}
	return joinedAndMeParticipants, nil
}

// chats

// get last chat in acvite sessions in whith the user joins
func (ss *SessionServices) GetLastChatInActiveSessions() ([]*database.TQueryLastChat, error) {
	var lastChats []*database.TQueryLastChat
	err := ss.useTransaction(func(tx *sql.Tx) error {
		var err error
		lastChats, err = ss.repositories.SessionChatRepository.QueryLastChatInActiveSessions(tx, ss.loginUserID)
		return err
	})

	if err != nil {
		return nil, err
	}
	return lastChats, nil
}

// get chats in session of [sessionID] in 48 hours
func (ss *SessionServices) GetChatsAtIn48Hours(sessionID int) ([]*database.TQuerySessionChat, error) {
	var chats []*database.TQuerySessionChat
	inRange := database.TQuerySessionChatInRange{
		StartDate: time.Now().Add(-48 * time.Hour),
		EndDate:   time.Now(),
	}
	err := ss.useTransaction(func(tx *sql.Tx) error {
		ok, err := ss.repositories.SessionRepository.HasStatusAt(tx, sessionID, ss.loginUserID, []database.TParticipantStatus{database.TJoinedParty})
		if err != nil {
			return err
		}
		if !ok {
			return ErrCannotAccessSession
		}
		chats, err = ss.repositories.SessionChatRepository.QueryBySessionIDInRange(tx, sessionID, inRange)
		return err
	})

	if err != nil {
		return nil, err
	}
	return chats, nil
}

// send chat content
func (ss *SessionServices) SendChatAt(sessionID int, content string) error {
	return ss.useTransaction(func(tx *sql.Tx) error {
		ok, err := ss.repositories.SessionRepository.HasStatusAt(tx, sessionID, ss.loginUserID, []database.TParticipantStatus{database.TJoinedParty})
		if err != nil {
			return err
		}
		if !ok {
			return ErrCannotAccessSession
		}
		_, err = ss.repositories.SessionChatRepository.Create(tx, sessionID, ss.loginUserID, content)
		return err
	})
}

// update participant status
func (ss *SessionServices) UpdateParticipantStatusAt(sessionID, userID int, status database.TParticipantStatus) error {
	// update status me
	if ss.loginUserID == userID {
		return ss.useTransaction(func(tx *sql.Tx) error {
			// updated user must have invited status
			if ok, err := ss.repositories.SessionRepository.HasStatusAt(tx, sessionID, userID, []database.TParticipantStatus{database.TInvitedParty}); err != nil {
				return err
			} else if !ok {
				return ErrCannotAccessSession
			}
			return ss.repositories.SessionParticipantRepository.UpdateStatusBySessionUserID(tx, sessionID, userID, status)
		})
	}
	// status that login user can change of userID in session is only reject
	// and updated user(userID) must have invited status.
	if status != database.TRejectedParty {
		return ErrCannotUpdateStatus
	}
	return ss.useTransaction(func(tx *sql.Tx) error {
		// login user must have joined status
		if ok, err := ss.repositories.SessionRepository.HasStatusAt(tx, sessionID, ss.loginUserID, []database.TParticipantStatus{database.TJoinedParty}); err != nil {
			return err
		} else if !ok {
			return ErrCannotAccessSession
		}

		// updated user must have invited status
		if ok, err := ss.repositories.SessionRepository.HasStatusAt(tx, sessionID, userID, []database.TParticipantStatus{database.TInvitedParty}); err != nil {
			return err
		} else if !ok {
			return ErrCannotAccessSession
		}

		return ss.repositories.SessionParticipantRepository.UpdateStatusBySessionUserID(tx, sessionID, userID, status)
	})
}

func (ss *SessionServices) IsJoined(sessionID, userID int) (bool, error) {
	joined := []database.TParticipantStatus{database.TJoinedParty}
	ok := false
	err := ss.useTransaction(func(tx *sql.Tx) error {
		var err error
		ok, err = ss.repositories.SessionRepository.HasStatusAt(tx, sessionID, userID, joined)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return false, err
	}
	return ok, nil
}
