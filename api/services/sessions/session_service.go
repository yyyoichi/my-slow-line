package sessions

import (
	"database/sql"
	"errors"
	"himakiwa/services/database"
	"time"
)

var (
	ErrCannotAccessSession = errors.New("cannot access session")
	ErrNoFoundUUID         = errors.New("not found uuid")
	ErrCannotUpdateStatus  = errors.New("cannot access participant session")
)

type UseSessionServicesFunc func(loginID int) *SessionServices

type SessionServices struct {
	repositories          *database.SessionRepositories
	recruitmentRepository database.FRecruitmentRepositoryInterface
	loginUserID           int
}

func NewSessionServices() UseSessionServicesFunc {
	ss := &SessionServices{repositories: database.NewSessionRepositories(), recruitmentRepository: &database.FRecruitmentRepository{}}
	return func(loginID int) *SessionServices {
		ss.loginUserID = loginID
		return ss
	}
}

// get active and archived sessions
func (ss *SessionServices) GetActiveOrArchivedSessions() ([]*database.TQuerySessions, error) {
	options := database.TQuerySessionsOptions{
		InPartyStatus:   []database.TParticipantStatus{database.TJoinedParty, database.TInvitedParty},
		InSessionStatus: []database.TSessionStatus{database.TActiveSession, database.TArchivedSession},
	}
	var sessions []*database.TQuerySessions
	err := database.WithTransaction(func(tx *sql.Tx) error {
		var err error
		sessions, err = ss.repositories.SessionRepository.QueryByUserID(tx, ss.loginUserID, options)
		return err
	})
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// get userID data from uuid
func (ss *SessionServices) LookUpRecruitment(uuid string) (*database.TFRecruitment, error) {
	recruit, err := ss.recruitmentRepository.QueryByUUID(uuid)
	if err != nil {
		return nil, err
	}
	if recruit.Deleted {
		return nil, ErrNoFoundUUID
	}

	return recruit, nil
}

// create session and loginUser invite userID
func (ss *SessionServices) CreateSession(publicKey, name string, userID int) error {
	return database.WithTransaction(func(tx *sql.Tx) error {
		sessionID, err := ss.repositories.SessionRepository.Create(tx, ss.loginUserID, publicKey, name)
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
}

// get session and participants
func (ss *SessionServices) GetSessionAt(sessionID int) (*database.TQuerySessions, []*database.TSessionParticipant, error) {
	var session *database.TQuerySessions
	var participants []*database.TSessionParticipant
	err := database.WithTransaction(func(tx *sql.Tx) error {
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
	return session, participants, nil
}

func (ss *SessionServices) UpdateSessionNameAt(sessionID int, name string) error {
	return database.WithTransaction(func(tx *sql.Tx) error {
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

// chats

// get last chat in acvite sessions in whith the user joins
func (ss *SessionServices) GetLastChatInActiveSessions() ([]*database.TQueryLastChat, error) {
	var lastChats []*database.TQueryLastChat
	err := database.WithTransaction(func(tx *sql.Tx) error {
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
	err := database.WithTransaction(func(tx *sql.Tx) error {
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
	return database.WithTransaction(func(tx *sql.Tx) error {
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
	if ss.loginUserID == userID {
		return database.WithTransaction(func(tx *sql.Tx) error {
			err := ss.repositories.SessionParticipantRepository.UpdateStatusBySessionUserID(tx, sessionID, userID, status)
			return err
		})
	}
	if status != database.TRejectedParty {
		return ErrCannotUpdateStatus
	}
	return database.WithTransaction(func(tx *sql.Tx) error {
		ok, err := ss.repositories.SessionRepository.HasStatusAt(tx, sessionID, ss.loginUserID, []database.TParticipantStatus{database.TJoinedParty})
		if err != nil {
			return err
		}
		if !ok {
			return ErrCannotAccessSession
		}

		err = ss.repositories.SessionParticipantRepository.UpdateStatusBySessionUserID(tx, sessionID, userID, status)
		return err
	})
}
