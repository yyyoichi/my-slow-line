package database

import (
	"database/sql"
	"time"
)

type TSessionStatus string

const (
	TActiveSession   TSessionStatus = "active"
	TArchivedSession TSessionStatus = "archived"
	TBreakupSession  TSessionStatus = "breakup"
)

type TParticipantStatus string

const (
	TInvitedParty  ParticipantStatus = "invited"
	TJoinedParty   ParticipantStatus = "joined"
	TRejectedParty ParticipantStatus = "rejected"
)

type SessionRepository struct{}
type SessionRepositoryInterface interface {
	QueryByUserID(tx *sql.Tx, userID int, options TQuerySessionsOptions) ([]*TQuerySessions, error)
	QueryBySessionUserID(tx *sql.Tx, sessionID, userID int) (*TQuerySessions, error)
	HasStatusAt(tx *sql.Tx, sessionID, userID int, inStatus []TParticipantStatus) (bool, error)
	Create(tx *sql.Tx, userID int, publicKey string, name string) (int, error)
	UpdateName(tx *sql.Tx, id int, name string) error
	HardDelete(tx *sql.Tx, sessionID int) error
	SoftDelete(tx *sql.Tx, sessionID int) error
	HardDeleteAll(tx *sql.Tx, sessionID int) error
}

// query sessions by userID

type TQuerySessions struct {
	ID            int
	Name          string
	PublicKey     string
	SessionStatus TSessionStatus
	Status        TParticipantStatus
	CreateAt      time.Time
	UpdateAt      time.Time
	Deleted       bool
}
type TQuerySessionsOptions struct {
	InPartyStatus   []TParticipantStatus
	InSessionStatus []TSessionStatus
}

// query sessions
func (sr *SessionRepository) QueryByUserID(tx *sql.Tx, userID int, options TQuerySessionsOptions) ([]*TQuerySessions, error) {
	// query
	query := `
	SELECT s.id, s.public_key, s.name, s.status, s.create_at, s.update_at, s.deleted, p.status 
	FROM chat_sessions AS s 
		RIGHT JOIN chat_session_participants AS p ON s.id = p.chat_session_id
	WHERE p.user_id = ? 
		AND p.status IN (?)
		AND s.status IN (?)`
	rows, err := tx.Query(query, userID, options.InPartyStatus, options.InSessionStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*TQuerySessions
	for rows.Next() {
		s := &TQuerySessions{}
		err := rows.Scan(&s.ID, &s.PublicKey, &s.Name, &s.SessionStatus, &s.CreateAt, &s.UpdateAt, &s.Deleted, &s.Status)
		if err != nil {
			return nil, err
		}
		results = append(results, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// query a session
func (sr *SessionRepository) QueryBySessionUserID(tx *sql.Tx, sessionID, userID int) (*TQuerySessions, error) {
	// query
	query := `
	SELECT s.id, s.public_key, s.name, s.status, s.create_at, s.update_at, s.deleted, p.status 
	FROM chat_sessions AS s 
		RIGHT JOIN chat_session_participants AS p ON s.id = p.chat_session_id
	WHERE s.id = ? AND p.user_id = ?`
	row := tx.QueryRow(query, sessionID, userID)

	// result
	s := &TQuerySessions{}
	err := row.Scan(&s.ID, &s.PublicKey, &s.Name, &s.SessionStatus, &s.CreateAt, &s.UpdateAt, &s.Deleted, &s.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found, return nil without an error.
			return nil, nil
		}
		return nil, err
	}

	return s, nil
}

// has status at session for user
func (sr *SessionRepository) HasStatusAt(tx *sql.Tx, sessionID, userID int, inStatus []TParticipantStatus) (bool, error) {
	// query
	query := `
	SELECT s.id
	FROM chat_sessions AS s 
		RIGHT JOIN chat_session_participants AS p ON s.id = p.chat_session_id
	WHERE s.id = ? AND p.user_id = ? AND p.status IN (?)`
	row := tx.QueryRow(query, sessionID, userID)

	// result
	id := 0
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found, return nil without an error.
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// create
func (sr *SessionRepository) Create(tx *sql.Tx, userID int, publicKey string, name string) (int, error) {
	// query
	query := `INSERT INTO chat_sessions (user_id, public_key, name) VALUES (?, ?, ?)`
	result, err := tx.Exec(query, userID, publicKey, name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// update
func (sr *SessionRepository) UpdateName(tx *sql.Tx, id int, name string) error {
	// query
	query := `UPDATE chat_sessions SET name = ? WHERE id = ?`
	_, err := tx.Exec(query, name, id)
	if err != nil {
		return err
	}
	return nil
}

// delete row
func (sr *SessionRepository) HardDelete(tx *sql.Tx, sessionID int) error {
	// query
	query := `DELETE FROM chat_sessions WHERE id = ?`
	_, err := tx.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}

// up deleted flag
func (sr *SessionRepository) SoftDelete(tx *sql.Tx, sessionID int) error {
	// query
	query := `UPDATE chat_sessions SET deleted = 1 WHERE id = ?`
	_, err := tx.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}

// delete hard all session contents
func (sr *SessionRepository) HardDeleteAll(tx *sql.Tx, sessionID int) error {
	// query
	querys := []string{
		`DELETE FROM chats WHERE chat_session_id = ?`,
		`DELETE FROM chat_session_participants WHERE chat_session_id = ?`,
		`DELETE FROM chat_sessions WHERE id = ?`,
	}
	for _, query := range querys {
		_, err := tx.Exec(query, sessionID)
		if err != nil {
			return err
		}
	}
	return nil
}

// query session participants

type SessionParticipantRepository struct{}
type SessionParticipantRepositoryInterface interface {
	QueryBySessionID(tx *sql.Tx, sessionID int) ([]TSessionParticipant, error)
	Create(tx *sql.Tx, sessionID, userID, inviteUserID int, status TParticipantStatus) (int, error)
	UpdateStatus(tx *sql.Tx, id int, status TParticipantStatus) error
	HardDelete(tx *sql.Tx, sessionID int) error
}

type TSessionParticipant struct {
	ID        int
	SessionID int
	UserID    int
	Status    TParticipantStatus
	CreateAt  time.Time
	UpdateAt  time.Time
	Deleted   bool
}

func (spr *SessionParticipantRepository) QueryBySessionID(tx *sql.Tx, sessionID int) ([]TSessionParticipant, error) {
	// query
	query := `SELECT id, chat_session_id, user_id, status, create_at, update_at, deleted 
	FROM chat_session_participants 
	WHERE chat_session_id = ?`
	rows, err := tx.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []TSessionParticipant
	for rows.Next() {
		sp := TSessionParticipant{}
		err := rows.Scan(&sp.ID, &sp.SessionID, &sp.UserID, &sp.Status, &sp.CreateAt, &sp.UpdateAt, &sp.Deleted)
		if err != nil {
			return nil, err
		}
		results = append(results, sp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// [invite_user_id] is ID of the user exec to the session
func (spr *SessionParticipantRepository) Create(tx *sql.Tx, sessionID, userID, inviteUserID int, status TParticipantStatus) (int, error) {
	// query
	query := `INSERT INTO chat_session_participants (chat_session_id, user_id, invite_user_id, status) VALUES (?, ?, ?, ?)`
	result, err := tx.Exec(query, sessionID, userID, inviteUserID, status)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// update
func (spr *SessionParticipantRepository) UpdateStatus(tx *sql.Tx, id int, status TParticipantStatus) error {
	// query
	query := `UPDATE chat_session_participants SET status = ? WHERE id = ?`
	_, err := tx.Exec(query, status, id)
	if err != nil {
		return err
	}
	return nil
}

// delete
func (cspr *ChatSessionParticipantRepository) HardDelete(tx *sql.Tx, sessionID int) error {
	// query
	query := `DELETE FROM chat_session_participants WHERE chat_session_id = ?`
	_, err := tx.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}
