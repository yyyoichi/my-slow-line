package sessions

import (
	"database/sql"
	"himakiwa/services/database"
)

////////////////////////////////////
//////./ mock services /////////////
////////////////////////////////////

func NewSessionServicesMock() UseSessionServicesFunc {

	tx := &sql.Tx{}
	userID1 := 1
	userID2 := 2

	useTransaction := func(fn func(tx *sql.Tx) error) error {
		return fn(tx)
	}
	mock := database.NewSessionRepositoriesMock()
	loginUserID := userID1
	ss := &SessionServices{useTransaction, mock, loginUserID}

	// insert data

	var sessionID int
	// session1 invite userID2
	sessionID = 1
	ss.CreateSession("", "Session1", userID2)

	// session2 reject userID2
	sessionID = 2
	ss.CreateSession("", "Session2", userID2)
	ss.repositories.SessionParticipantRepository.UpdateStatusBySessionUserID(tx, sessionID, userID2, database.TRejectedParty)

	// session3 joined userID2 and chats
	sessionID = 3
	ss.CreateSession("", "Session3", userID2)
	ss.repositories.SessionParticipantRepository.UpdateStatusBySessionUserID(tx, sessionID, userID2, database.TJoinedParty)
	ss.SendChatAt(3, "Hello, I am 1 in the session 3")
	ss.repositories.SessionChatRepository.Create(tx, sessionID, userID2, "Hello, I am 2 in the session 3")

	// session4 is same session3, session4 is arcived
	sessionID = 4
	ss.CreateSession("", "Session4", userID2)
	ss.repositories.SessionParticipantRepository.UpdateStatusBySessionUserID(tx, sessionID, userID2, database.TJoinedParty)
	ss.repositories.SessionChatRepository.Create(tx, sessionID, userID2, "Hello, I am 2 in the session 4")
	ss.SendChatAt(sessionID, "Hello, I am 1 in the session 4")
	ss.repositories.SessionRepository.UpdateStatus(tx, sessionID, database.TArchivedSession)

	// session5 is breakup
	sessionID = 5
	ss.CreateSession("", "Session5", userID2)
	ss.repositories.SessionRepository.UpdateStatus(tx, sessionID, database.TBreakupSession)

	// exp session6

	return func(loginID int) *SessionServices {
		ss.loginUserID = loginID
		return ss
	}
}
